//nolint:dupl
package handler_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/components/kohan/core"
	"github.com/amanhigh/go-fun/components/kohan/handler"
	"github.com/amanhigh/go-fun/components/kohan/manager"
	"github.com/amanhigh/go-fun/components/kohan/repository"
	"github.com/amanhigh/go-fun/models/barkat"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"
)

func decodeNoteResponse(w *httptest.ResponseRecorder) barkat.Note {
	var envelope common.Envelope[barkat.Note]
	util.AssertSuccess(w, http.StatusCreated, &envelope)
	return envelope.Data
}

func decodeNoteListResponse(w *httptest.ResponseRecorder) []barkat.Note {
	var envelope common.Envelope[barkat.NoteList]
	util.AssertSuccess(w, http.StatusOK, &envelope)
	return envelope.Data.Notes
}

// NoteHandler Integration Tests - Comprehensive Master Specification
// Tests complete HTTP → Handler → Manager → Repository → Database flow
// Covers all PRD validations for Section 2.3 JournalNote APIs
//
// TEST STRUCTURE FORMAT:
// ====================
// Describe(API)
// -> Context(Happy Path): 2xx Success Cases
// -> Context(Field Validations): All 4xx Validation Cases
//    -> Context(Field Name): One Context for Each Field
//       -> Context(Allowed Values): All Variations of Valid Values (2xx) - If Applicable
//       -> Context(Bad Values): All Variations of Missing,Regex,Min,Max Edge Cases (4xx)
// -> Context(Errors): 5xx Server Error Cases

var _ = Describe("NoteHandler Integration - Section 2.3 JournalNote APIs", func() {
	var (
		noteHandler *handler.NoteHandlerImpl
		router      *gin.Engine
		testCtx     = context.Background()
		db          *gorm.DB
		journalMgr  manager.JournalManager
		noteMgr     manager.NoteManager
		journal     barkat.Journal
		req         *http.Request
		w           *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		var err error
		db, err = core.CreateTestBarkatDB()
		Expect(err).ToNot(HaveOccurred())

		journalRepo := repository.NewJournalRepository(db)
		journalMgr = manager.NewJournalManager(journalRepo)
		noteMgr = manager.NewNoteManager(journalMgr, repository.NewNoteRepository(db))
		noteHandler = handler.NewNoteHandler(noteMgr)

		router = util.CreateTestGinRouter()
		v1 := router.Group("/v1/api")
		journalGroup := v1.Group("/journals")
		handler.SetupNoteRoutes(journalGroup, noteHandler)

		// Create base journal for note operations
		journal = barkat.Journal{
			Ticker:   "GRSE",
			Sequence: "MWD",
			Type:     "REJECTED",
			Status:   "FAIL",
			Images: []barkat.Image{
				{Timeframe: "DL", FileName: "test-dl.png"},
				{Timeframe: "WK", FileName: "test-wk.png"},
				{Timeframe: "MN", FileName: "test-mn.png"},
				{Timeframe: "TMN", FileName: "test-tmn.png"},
			},
		}
		Expect(journalMgr.CreateJournal(testCtx, &journal)).To(Succeed())
	})

	AfterEach(func() {
		sqlDB, err := db.DB()
		Expect(err).ToNot(HaveOccurred())
		sqlDB.Close()
	})

	// ============================================================================
	// 2.3.1 POST /v1/journals/{journal-id}/notes - Add Note
	// ============================================================================
	Describe("POST /v1/journals/{journal-id}/notes - Add Note (2.3.1)", func() {
		Context("Happy Path", func() {
			Context("with valid note data", func() {
				var response barkat.Note

				BeforeEach(func() {
					note := barkat.Note{
						Status:  "SET",
						Content: "Strong OE at weekly level, watching for confirmation on daily.",
						Format:  "MARKDOWN",
					}
					req, w = util.CreateTestRequest("POST", barkat.JournalBase+"/"+journal.ExternalID+"/notes", note)
					router.ServeHTTP(w, req)
				})

				It("should return 201 Created", func() {
					Expect(w.Code).To(Equal(http.StatusCreated))
				})

				It("should return Envelope success", func() {
					var envelope common.Envelope[barkat.Note]
					util.AssertSuccess(w, http.StatusCreated, &envelope)
					Expect(envelope.Status).To(Equal(common.EnvelopeSuccess))
				})

				It("should return created note with external ID", func() {
					response = decodeNoteResponse(w)
					Expect(response.ExternalID).To(HavePrefix("not_"))
				})

				It("should preserve status field", func() {
					response = decodeNoteResponse(w)
					Expect(response.Status).To(Equal("SET"))
				})

				It("should preserve content field", func() {
					response = decodeNoteResponse(w)
					Expect(response.Content).To(Equal("Strong OE at weekly level, watching for confirmation on daily."))
				})

				It("should preserve format field", func() {
					response = decodeNoteResponse(w)
					Expect(response.Format).To(Equal("MARKDOWN"))
				})

				It("should set created_at timestamp", func() {
					response = decodeNoteResponse(w)
					Expect(response.CreatedAt).ToNot(BeZero())
				})

				It("should persist note to database", func() {
					noteList, err := noteMgr.ListNotes(testCtx, journal.ExternalID, "")
					Expect(err).ToNot(HaveOccurred())
					Expect(noteList.Notes).To(HaveLen(1))
				})
			})

			Context("with format defaulting to MARKDOWN", func() {
				It("should default format to MARKDOWN when omitted", func() {
					note := barkat.Note{
						Status:  "SET",
						Content: "Note without explicit format.",
					}
					req, w = util.CreateTestRequest("POST", barkat.JournalBase+"/"+journal.ExternalID+"/notes", note)
					router.ServeHTTP(w, req)
					response := decodeNoteResponse(w)
					Expect(response.Format).To(Equal("MARKDOWN"))
				})
			})
		})

		Context("Field Validations", func() {
			Context("Status Field", func() {
				Context("Allowed Values", func() {
					It("should accept status = SET", func() {
						note := barkat.Note{Status: "SET", Content: "Test content", Format: "MARKDOWN"}
						req, w = util.CreateTestRequest("POST", barkat.JournalBase+"/"+journal.ExternalID+"/notes", note)
						router.ServeHTTP(w, req)
						response := decodeNoteResponse(w)
						Expect(response.Status).To(Equal("SET"))
					})

					It("should accept status = RUNNING", func() {
						note := barkat.Note{Status: "RUNNING", Content: "Test content", Format: "MARKDOWN"}
						req, w = util.CreateTestRequest("POST", barkat.JournalBase+"/"+journal.ExternalID+"/notes", note)
						router.ServeHTTP(w, req)
						response := decodeNoteResponse(w)
						Expect(response.Status).To(Equal("RUNNING"))
					})

					It("should accept status = DROPPED", func() {
						note := barkat.Note{Status: "DROPPED", Content: "Test content", Format: "MARKDOWN"}
						req, w = util.CreateTestRequest("POST", barkat.JournalBase+"/"+journal.ExternalID+"/notes", note)
						router.ServeHTTP(w, req)
						response := decodeNoteResponse(w)
						Expect(response.Status).To(Equal("DROPPED"))
					})

					It("should accept status = TAKEN", func() {
						note := barkat.Note{Status: "TAKEN", Content: "Test content", Format: "MARKDOWN"}
						req, w = util.CreateTestRequest("POST", barkat.JournalBase+"/"+journal.ExternalID+"/notes", note)
						router.ServeHTTP(w, req)
						response := decodeNoteResponse(w)
						Expect(response.Status).To(Equal("TAKEN"))
					})

					It("should accept status = REJECTED", func() {
						note := barkat.Note{Status: "REJECTED", Content: "Test content", Format: "MARKDOWN"}
						req, w = util.CreateTestRequest("POST", barkat.JournalBase+"/"+journal.ExternalID+"/notes", note)
						router.ServeHTTP(w, req)
						response := decodeNoteResponse(w)
						Expect(response.Status).To(Equal("REJECTED"))
					})

					It("should accept status = SUCCESS", func() {
						note := barkat.Note{Status: "SUCCESS", Content: "Test content", Format: "MARKDOWN"}
						req, w = util.CreateTestRequest("POST", barkat.JournalBase+"/"+journal.ExternalID+"/notes", note)
						router.ServeHTTP(w, req)
						response := decodeNoteResponse(w)
						Expect(response.Status).To(Equal("SUCCESS"))
					})

					It("should accept status = FAIL", func() {
						note := barkat.Note{Status: "FAIL", Content: "Test content", Format: "MARKDOWN"}
						req, w = util.CreateTestRequest("POST", barkat.JournalBase+"/"+journal.ExternalID+"/notes", note)
						router.ServeHTTP(w, req)
						response := decodeNoteResponse(w)
						Expect(response.Status).To(Equal("FAIL"))
					})

					It("should accept status = MISSED", func() {
						note := barkat.Note{Status: "MISSED", Content: "Test content", Format: "MARKDOWN"}
						req, w = util.CreateTestRequest("POST", barkat.JournalBase+"/"+journal.ExternalID+"/notes", note)
						router.ServeHTTP(w, req)
						response := decodeNoteResponse(w)
						Expect(response.Status).To(Equal("MISSED"))
					})

					It("should accept status = JUST_LOSS", func() {
						note := barkat.Note{Status: "JUST_LOSS", Content: "Test content", Format: "MARKDOWN"}
						req, w = util.CreateTestRequest("POST", barkat.JournalBase+"/"+journal.ExternalID+"/notes", note)
						router.ServeHTTP(w, req)
						response := decodeNoteResponse(w)
						Expect(response.Status).To(Equal("JUST_LOSS"))
					})

					It("should accept status = BROKEN", func() {
						note := barkat.Note{Status: "BROKEN", Content: "Test content", Format: "MARKDOWN"}
						req, w = util.CreateTestRequest("POST", barkat.JournalBase+"/"+journal.ExternalID+"/notes", note)
						router.ServeHTTP(w, req)
						response := decodeNoteResponse(w)
						Expect(response.Status).To(Equal("BROKEN"))
					})
				})

				Context("Bad Values", func() {
					It("should return 400 for missing status (PRD: required)", func() {
						note := barkat.Note{Status: "", Content: "Test content", Format: "MARKDOWN"}
						req, w = util.CreateTestRequest("POST", barkat.JournalBase+"/"+journal.ExternalID+"/notes", note)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Status", "required")
					})

					It("should return 400 for invalid status enum (PRD: must be valid status)", func() {
						note := barkat.Note{Status: "INVALID", Content: "Test content", Format: "MARKDOWN"}
						req, w = util.CreateTestRequest("POST", barkat.JournalBase+"/"+journal.ExternalID+"/notes", note)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Status", "oneof")
					})

					It("should return 400 for lowercase status (PRD: case-sensitive)", func() {
						note := barkat.Note{Status: "set", Content: "Test content", Format: "MARKDOWN"}
						req, w = util.CreateTestRequest("POST", barkat.JournalBase+"/"+journal.ExternalID+"/notes", note)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Status", "oneof")
					})
				})
			})

			Context("Content Field", func() {
				Context("Allowed Values", func() {
					It("should accept minimum content length (1 char)", func() {
						note := barkat.Note{Status: "SET", Content: "X", Format: "MARKDOWN"}
						req, w = util.CreateTestRequest("POST", barkat.JournalBase+"/"+journal.ExternalID+"/notes", note)
						router.ServeHTTP(w, req)
						response := decodeNoteResponse(w)
						Expect(response.Content).To(Equal("X"))
					})

					It("should accept content with special characters", func() {
						note := barkat.Note{Status: "SET", Content: "Test with special chars: @#$%^&*()", Format: "MARKDOWN"}
						req, w = util.CreateTestRequest("POST", barkat.JournalBase+"/"+journal.ExternalID+"/notes", note)
						router.ServeHTTP(w, req)
						response := decodeNoteResponse(w)
						Expect(response.Content).To(Equal("Test with special chars: @#$%^&*()"))
					})

					It("should accept content with newlines", func() {
						note := barkat.Note{Status: "SET", Content: "Line 1\nLine 2\nLine 3", Format: "MARKDOWN"}
						req, w = util.CreateTestRequest("POST", barkat.JournalBase+"/"+journal.ExternalID+"/notes", note)
						router.ServeHTTP(w, req)
						response := decodeNoteResponse(w)
						Expect(response.Content).To(ContainSubstring("\n"))
					})

					It("should accept content with markdown formatting", func() {
						note := barkat.Note{Status: "SET", Content: "# Header\n- Item 1\n- Item 2\n**bold**", Format: "MARKDOWN"}
						req, w = util.CreateTestRequest("POST", barkat.JournalBase+"/"+journal.ExternalID+"/notes", note)
						router.ServeHTTP(w, req)
						response := decodeNoteResponse(w)
						Expect(response.Content).To(ContainSubstring("# Header"))
					})
				})

				Context("Bad Values", func() {
					It("should return 400 for missing content (PRD: required)", func() {
						note := barkat.Note{Status: "SET", Content: "", Format: "MARKDOWN"}
						req, w = util.CreateTestRequest("POST", barkat.JournalBase+"/"+journal.ExternalID+"/notes", note)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Content", "required")
					})

					It("should return 400 for content exceeding max length (PRD: max 2000 chars)", func() {
						var longContent strings.Builder
						for range barkat.NoteContentMaxLength + 1 {
							longContent.WriteString("X")
						}
						note := barkat.Note{Status: "SET", Content: longContent.String(), Format: "MARKDOWN"}
						req, w = util.CreateTestRequest("POST", barkat.JournalBase+"/"+journal.ExternalID+"/notes", note)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Content", "max")
					})
				})
			})

			Context("Format Field", func() {
				Context("Allowed Values", func() {
					It("should accept format = MARKDOWN", func() {
						note := barkat.Note{Status: "SET", Content: "Test content", Format: "MARKDOWN"}
						req, w = util.CreateTestRequest("POST", barkat.JournalBase+"/"+journal.ExternalID+"/notes", note)
						router.ServeHTTP(w, req)
						response := decodeNoteResponse(w)
						Expect(response.Format).To(Equal("MARKDOWN"))
					})

					It("should accept format = PLAINTEXT", func() {
						note := barkat.Note{Status: "SET", Content: "Test content", Format: "PLAINTEXT"}
						req, w = util.CreateTestRequest("POST", barkat.JournalBase+"/"+journal.ExternalID+"/notes", note)
						router.ServeHTTP(w, req)
						response := decodeNoteResponse(w)
						Expect(response.Format).To(Equal("PLAINTEXT"))
					})
				})

				Context("Bad Values", func() {
					It("should return 400 for invalid format enum (PRD: must be MARKDOWN or PLAINTEXT)", func() {
						note := barkat.Note{Status: "SET", Content: "Test content", Format: "HTML"}
						req, w = util.CreateTestRequest("POST", barkat.JournalBase+"/"+journal.ExternalID+"/notes", note)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Format", "oneof")
					})

					It("should return 400 for lowercase format (PRD: case-sensitive)", func() {
						note := barkat.Note{Status: "SET", Content: "Test content", Format: "markdown"}
						req, w = util.CreateTestRequest("POST", barkat.JournalBase+"/"+journal.ExternalID+"/notes", note)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Format", "oneof")
					})
				})
			})

			Context("CreatedAt Field", func() {
				Context("Allowed Values", func() {
					It("should accept valid ISO 8601 datetime (PRD: optional for migration)", func() {
						historicalTime := time.Date(2023, 6, 15, 14, 30, 0, 0, time.UTC)
						note := barkat.Note{
							Status:    "SET",
							Content:   "Strong OE at weekly level, watching for confirmation on daily.",
							Format:    "MARKDOWN",
							CreatedAt: historicalTime,
						}
						req, w = util.CreateTestRequest("POST", barkat.JournalBase+"/"+journal.ExternalID+"/notes", note)
						router.ServeHTTP(w, req)
						response := decodeNoteResponse(w)
						Expect(response.CreatedAt).To(Equal(historicalTime))
					})

					It("should accept nil CreatedAt (PRD: system sets current timestamp)", func() {
						note := barkat.Note{
							Status:  "SET",
							Content: "Strong OE at weekly level, watching for confirmation on daily.",
							Format:  "MARKDOWN",
							// CreatedAt left as zero value
						}
						req, w = util.CreateTestRequest("POST", barkat.JournalBase+"/"+journal.ExternalID+"/notes", note)
						router.ServeHTTP(w, req)
						response := decodeNoteResponse(w)
						Expect(response.CreatedAt).ToNot(BeZero())
						Expect(response.CreatedAt).To(BeTemporally("~", time.Now(), 5*time.Second))
					})

					It("should accept zero CreatedAt (PRD: BeforeCreate hook sets current time)", func() {
						note := barkat.Note{
							Status:    "SET",
							Content:   "Strong OE at weekly level, watching for confirmation on daily.",
							Format:    "MARKDOWN",
							CreatedAt: time.Time{}, // Zero time - BeforeCreate will set current time
						}
						req, w = util.CreateTestRequest("POST", barkat.JournalBase+"/"+journal.ExternalID+"/notes", note)
						router.ServeHTTP(w, req)
						response := decodeNoteResponse(w)
						Expect(response.CreatedAt).ToNot(BeZero())
						Expect(response.CreatedAt).To(BeTemporally("~", time.Now(), 5*time.Second))
					})
				})

				Context("Bad Values", func() {
					// No bad values for CreatedAt - it's optional and zero values are handled by BeforeCreate
				})
			})

			Context("Journal ID Path Parameter", func() {
				Context("Bad Values", func() {
					It("should return 400 for malformed journal ID", func() {
						note := barkat.Note{Status: "SET", Content: "Test content", Format: "MARKDOWN"}
						req, w = util.CreateTestRequest("POST", barkat.JournalBase+"/invalid-uuid-format/notes", note)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusBadRequest))
					})

					It("should return 404 for valid ID format but non-existent", func() {
						note := barkat.Note{Status: "SET", Content: "Test content", Format: "MARKDOWN"}
						req, w = util.CreateTestRequest("POST", barkat.JournalBase+"/jrn_12345678/notes", note)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusNotFound))
					})
				})
			})
		})

		Context("Errors", func() {
			It("should return 400 for invalid JSON", func() {
				req, w = util.CreateTestRequest("POST", barkat.JournalBase+"/"+journal.ExternalID+"/notes", []byte("invalid json"))
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})

			It("should return 400 for empty request body", func() {
				req, w = util.CreateTestRequest("POST", barkat.JournalBase+"/"+journal.ExternalID+"/notes", []byte(""))
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})

			It("should return 400 for null request body", func() {
				req, w = util.CreateTestRequest("POST", barkat.JournalBase+"/"+journal.ExternalID+"/notes", []byte("null"))
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})
		})
	})

	// ============================================================================
	// 2.3.2 GET /v1/journals/{journal-id}/notes - List Notes
	// ============================================================================
	Describe("GET /v1/journals/{journal-id}/notes - List Notes (2.3.2)", func() {
		Context("Happy Path", func() {
			Context("with journal having notes", func() {
				var notes []barkat.Note

				BeforeEach(func() {
					// Create multiple notes for testing
					note1 := barkat.Note{Status: "SET", Content: "First note", Format: "MARKDOWN"}
					_, err := noteMgr.CreateNote(testCtx, journal.ExternalID, note1)
					Expect(err).ToNot(HaveOccurred())

					note2 := barkat.Note{Status: "TAKEN", Content: "Second note", Format: "PLAINTEXT"}
					_, err = noteMgr.CreateNote(testCtx, journal.ExternalID, note2)
					Expect(err).ToNot(HaveOccurred())

					req, w = util.CreateTestRequest("GET", barkat.JournalBase+"/"+journal.ExternalID+"/notes", nil)
					router.ServeHTTP(w, req)
				})

				It("should return 200 OK", func() {
					Expect(w.Code).To(Equal(http.StatusOK))
				})

				It("should return all notes for journal", func() {
					notes = decodeNoteListResponse(w)
					Expect(notes).To(HaveLen(2))
				})

				It("should return notes with correct statuses", func() {
					notes = decodeNoteListResponse(w)
					statuses := []string{}
					for _, note := range notes {
						statuses = append(statuses, note.Status)
					}
					Expect(statuses).To(ContainElements("SET", "TAKEN"))
				})

				It("should return notes with external IDs", func() {
					notes = decodeNoteListResponse(w)
					for _, note := range notes {
						Expect(note.ExternalID).To(HavePrefix("not_"))
					}
				})

				It("should return notes with created_at timestamps", func() {
					notes = decodeNoteListResponse(w)
					for _, note := range notes {
						Expect(note.CreatedAt).ToNot(BeZero())
					}
				})
			})

			Context("with journal having no notes", func() {
				BeforeEach(func() {
					req, w = util.CreateTestRequest("GET", barkat.JournalBase+"/"+journal.ExternalID+"/notes", nil)
					router.ServeHTTP(w, req)
				})

				It("should return 200 OK with empty array", func() {
					Expect(w.Code).To(Equal(http.StatusOK))
					notes := decodeNoteListResponse(w)
					Expect(notes).To(BeEmpty())
				})
			})
		})

		Context("Field Validations", func() {
			Context("Journal ID Path Parameter", func() {
				Context("Bad Values", func() {
					It("should return 400 for malformed journal ID", func() {
						note := barkat.Note{Status: "SET", Content: "Test content", Format: "MARKDOWN"}
						req, w = util.CreateTestRequest("GET", barkat.JournalBase+"/invalid_format/notes", note)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusBadRequest))
					})

					It("should return 404 for valid ID format but non-existent", func() {
						req, w = util.CreateTestRequest("GET", barkat.JournalBase+"/jrn_12345678/notes", nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusNotFound))
					})
				})
			})
		})

		Context("Errors", func() {
			// No server error scenarios for GET currently
		})
	})

	// ============================================================================
	// 2.3.3 DELETE /v1/journals/{journal-id}/notes/{note-id} - Remove Note
	// ============================================================================
	Describe("DELETE /v1/journals/{journal-id}/notes/{note-id} - Remove Note (2.3.3)", func() {
		var noteToDelete *barkat.Note

		BeforeEach(func() {
			// Create a note to delete
			note := barkat.Note{Status: "SET", Content: "Note to delete", Format: "MARKDOWN"}
			var err error
			noteToDelete, err = noteMgr.CreateNote(testCtx, journal.ExternalID, note)
			Expect(err).ToNot(HaveOccurred())
		})

		Context("Happy Path", func() {
			Context("with valid journal and note IDs", func() {
				BeforeEach(func() {
					req, w = util.CreateTestRequest("DELETE", barkat.JournalBase+"/"+journal.ExternalID+"/notes/"+noteToDelete.ExternalID, nil)
					router.ServeHTTP(w, req)
				})

				It("should return 204 No Content", func() {
					Expect(w.Code).To(Equal(http.StatusNoContent))
				})

				It("should return empty body", func() {
					Expect(w.Body.String()).To(BeEmpty())
				})

				It("should actually delete the note from database", func() {
					noteList, err := noteMgr.ListNotes(testCtx, journal.ExternalID, "")
					Expect(err).ToNot(HaveOccurred())
					Expect(noteList.Notes).To(BeEmpty())
				})
			})
		})

		Context("Field Validations", func() {
			Context("Journal ID Path Parameter", func() {
				Context("Bad Values", func() {
					It("should return 400 for malformed journal ID", func() {
						req, w = util.CreateTestRequest("DELETE", barkat.JournalBase+"/invalid-uuid-format/notes/"+noteToDelete.ExternalID, nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusBadRequest))
					})

					It("should return 404 for valid ID format but non-existent", func() {
						req, w = util.CreateTestRequest("DELETE", barkat.JournalBase+"/jrn_12345678/notes/"+noteToDelete.ExternalID, nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusNotFound))
					})
				})
			})

			Context("Note ID Path Parameter", func() {
				Context("Bad Values", func() {
					It("should return 400 for invalid note ID format", func() {
						req, w = util.CreateTestRequest("DELETE", barkat.JournalBase+"/"+journal.ExternalID+"/notes/invalid-note-format", nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusBadRequest))
					})

					It("should return 404 for non-existent note ID", func() {
						req, w = util.CreateTestRequest("DELETE", barkat.JournalBase+"/"+journal.ExternalID+"/notes/not_99999999", nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusNotFound))
					})
				})
			})
		})

		Context("Errors", func() {
			It("should return 404 on second delete (idempotency check)", func() {
				// First delete
				req1, w1 := util.CreateTestRequest("DELETE", barkat.JournalBase+"/"+journal.ExternalID+"/notes/"+noteToDelete.ExternalID, nil)
				router.ServeHTTP(w1, req1)
				Expect(w1.Code).To(Equal(http.StatusNoContent))

				// Second delete should return 404
				req2, w2 := util.CreateTestRequest("DELETE", barkat.JournalBase+"/"+journal.ExternalID+"/notes/"+noteToDelete.ExternalID, nil)
				router.ServeHTTP(w2, req2)
				Expect(w2.Code).To(Equal(http.StatusNotFound))
			})
		})
	})
})
