//nolint:dupl
package handler_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/components/kohan/core"
	"github.com/amanhigh/go-fun/components/kohan/handler"
	"github.com/amanhigh/go-fun/components/kohan/manager"
	"github.com/amanhigh/go-fun/components/kohan/repository"
	"github.com/amanhigh/go-fun/models/barkat"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/gin-gonic/gin"
	"github.com/golang-sql/civil"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"
)

func decodeCreateJournalResponse(w *httptest.ResponseRecorder) barkat.Journal {
	var envelope common.Envelope[barkat.Journal]
	util.AssertSuccess(w, http.StatusCreated, &envelope)
	return envelope.Data
}

// JournalHandler Integration Tests - Comprehensive Master Specification
// Tests complete HTTP → Handler → Manager → Repository → Database flow
// Covers all PRD validations, enum values, edge cases, pagination, sorting, and error scenarios
//
// TEST STRUCTURE FORMAT:
// ====================
// Describe(API)
// -> Context(Happy Path): 2xx Success Cases
// -> Context(Field Validations): All 4xx Validation Cases
//    -> Context(Field Name): One Context for Each Field
//       -> Context(Allowed Values): All Varitions of Valid Values (2xx) - If Applicable
//       -> Context(Bad Values): All Varitions of Missing,Regex,Min,Max Edge Cases (4xx)
// -> Context(Errors): 5xx Server Error Cases
//
// This structure ensures comprehensive coverage of:
// - Every field with all its validation rules
// - All enum values and case sensitivity
// - Missing required fields and optional fields
// - Boundary conditions and edge cases
// - Error scenarios and business logic

var _ = Describe("JournalHandler Integration - CUD Tests", func() {
	var (
		journalHandler *handler.JournalHandlerImpl
		router         *gin.Engine
		testCtx        = context.Background()
		db             *gorm.DB
		journalMgr     manager.JournalManager
		req            *http.Request
		w              *httptest.ResponseRecorder
		// Common image objects to reduce duplication
		standardImages = []barkat.Image{
			{Timeframe: "DL", FileName: "test-dl.png"},
			{Timeframe: "WK", FileName: "test-wk.png"},
			{Timeframe: "MN", FileName: "test-mn.png"},
			{Timeframe: "TMN", FileName: "test-tmn.png"},
		}
	)

	BeforeEach(func() {
		// FIXME: #A Sqlite Snapshots, Comparison and Revert.
		var err error
		db, err = core.CreateTestBarkatDB()
		Expect(err).ToNot(HaveOccurred())

		journalRepo := repository.NewJournalRepository(db)
		journalMgr = manager.NewJournalManager(journalRepo)
		journalHandler = handler.NewJournalHandler(journalMgr)

		router = util.CreateTestGinRouter()
		v1 := router.Group("/v1")
		journal := v1.Group("/journals")
		handler.SetupJournalEntryRoutes(journal, journalHandler)
	})

	AfterEach(func() {
		sqlDB, err := db.DB()
		Expect(err).ToNot(HaveOccurred())
		sqlDB.Close()
	})

	Describe("POST /v1/journal - Create Entry", func() {
		defaultTimeframes := []string{"DL", "WK", "MN", "TMN"}

		createTestImages := func() []barkat.Image {
			images := make([]barkat.Image, len(defaultTimeframes))
			for i, timeframe := range defaultTimeframes {
				images[i] = barkat.Image{
					Timeframe: timeframe,
					FileName:  fmt.Sprintf("TEST.%s.rejected.oe__20240115_132138.png", timeframe),
				}
			}
			return images
		}

		Context("Happy Path", func() {
			Context("with minimal valid entry (required fields + min 4 images)", func() {
				var envelopeResponse common.Envelope[barkat.Journal]

				BeforeEach(func() {
					journal := barkat.Journal{
						Ticker:   "GRSE",
						Sequence: "MWD",
						Type:     "REJECTED",
						Status:   "FAIL",
						Images:   createTestImages(),
					}
					req, w = util.CreateTestRequest("POST", barkat.JournalEntries, journal)
					router.ServeHTTP(w, req)
				})

				It("should return 201 Created", func() {
					util.AssertSuccess(w, http.StatusCreated, &envelopeResponse)
				})

				It("should return Envelope success", func() {
					util.AssertSuccess(w, http.StatusCreated, &envelopeResponse)
					Expect(envelopeResponse.Status).To(Equal(common.EnvelopeSuccess))
				})

				It("should return created entry with ID in data", func() {
					util.AssertSuccess(w, http.StatusCreated, &envelopeResponse)
					Expect(envelopeResponse.Data.ExternalID).To(HavePrefix("jrn_"))
				})

				It("should preserve all input fields", func() {
					util.AssertSuccess(w, http.StatusCreated, &envelopeResponse)
					Expect(envelopeResponse.Data.Ticker).To(Equal("GRSE"))
					Expect(envelopeResponse.Data.Sequence).To(Equal("MWD"))
					Expect(envelopeResponse.Data.Type).To(Equal("REJECTED"))
					Expect(envelopeResponse.Data.Status).To(Equal("FAIL"))
				})

				It("should set created_at timestamp", func() {
					util.AssertSuccess(w, http.StatusCreated, &envelopeResponse)
					Expect(envelopeResponse.Data.CreatedAt).ToNot(BeZero())
				})

				It("should include default timeframe images in response", func() {
					util.AssertSuccess(w, http.StatusCreated, &envelopeResponse)
					Expect(envelopeResponse.Data.Images).To(HaveLen(len(defaultTimeframes)))

					recorded := make([]string, 0, len(envelopeResponse.Data.Images))
					for _, img := range envelopeResponse.Data.Images {
						Expect(img.ExternalID).To(HavePrefix("img_"))
						recorded = append(recorded, img.Timeframe)
					}
					Expect(recorded).To(ConsistOf(defaultTimeframes))
				})

				It("should persist journal to database", func() {
					util.AssertSuccess(w, http.StatusCreated, &envelopeResponse)
					dbJournal, err := journalMgr.GetJournal(testCtx, envelopeResponse.Data.ExternalID)
					Expect(err).ToNot(HaveOccurred())
					Expect(dbJournal.ExternalID).To(Equal(envelopeResponse.Data.ExternalID))
					Expect(dbJournal.Ticker).To(Equal("GRSE"))
				})
			})

			Context("with complete valid journal including images, tags, and notes", func() {
				BeforeEach(func() {
					journal := barkat.Journal{
						Ticker:   "RELIANCE",
						Sequence: "YR",
						Type:     "REJECTED",
						Status:   "RUNNING",
						Images:   createTestImages(),
						Tags: []barkat.Tag{
							{Tag: "oe", Type: "REASON"},
						},
						Notes: []barkat.Note{
							{Status: "SET", Content: "Strong OE at weekly level.", Format: "MARKDOWN"},
						},
					}
					req, w = util.CreateTestRequest("POST", barkat.JournalEntries, journal)
					router.ServeHTTP(w, req)
				})

				It("should create entry with all associations", func() {
					response := decodeCreateJournalResponse(w)
					Expect(response.ExternalID).To(HavePrefix("jrn_"))
					Expect(response.Ticker).To(Equal("RELIANCE"))
					Expect(response.Sequence).To(Equal("YR"))
					Expect(response.Type).To(Equal("REJECTED"))
					Expect(response.Status).To(Equal("RUNNING"))
					Expect(response.Images).To(HaveLen(4))
					Expect(response.Tags).To(HaveLen(1))
					Expect(response.Notes).To(HaveLen(1))

					// Validate associated entity IDs
					for _, img := range response.Images {
						Expect(img.ExternalID).To(HavePrefix("img_"))
					}
					for _, tag := range response.Tags {
						Expect(tag.ExternalID).To(HavePrefix("tag_"))
					}
					for _, note := range response.Notes {
						Expect(note.ExternalID).To(HavePrefix("not_"))
					}
				})
			})
		})

		Context("Field Validations", func() {
			Context("Ticker Field", func() {
				Context("Allowed Values", func() {
					It("should accept ticker with numbers", func() {
						journal := barkat.Journal{
							Ticker:   "GRSE123",
							Sequence: "MWD",
							Type:     "REJECTED",
							Status:   "FAIL",
							Images:   standardImages,
						}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, journal)
						router.ServeHTTP(w, req)
						response := decodeCreateJournalResponse(w)
						Expect(response.Ticker).To(Equal("GRSE123"))
					})

					It("should accept ticker with hyphen", func() {
						journal := barkat.Journal{
							Ticker:   "GRSE-NSE",
							Sequence: "MWD",
							Type:     "REJECTED",
							Status:   "FAIL",
							Images:   standardImages,
						}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, journal)
						router.ServeHTTP(w, req)
						response := decodeCreateJournalResponse(w)
						Expect(response.Ticker).To(Equal("GRSE-NSE"))
					})

					It("should accept ticker with dot suffix", func() {
						journal := barkat.Journal{
							Ticker:   "TCS.NS",
							Sequence: "YR",
							Type:     "SET",
							Status:   "RUNNING",
							Images:   standardImages,
						}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, journal)
						router.ServeHTTP(w, req)
						response := decodeCreateJournalResponse(w)
						Expect(response.Ticker).To(Equal("TCS.NS"))
					})

					It("should accept ticker at max length (10)", func() {
						journal := barkat.Journal{
							Ticker: "1234567890", Sequence: "MWD", Type: "REJECTED", Status: "FAIL",
							Images: standardImages,
						}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, journal)
						router.ServeHTTP(w, req)
						response := decodeCreateJournalResponse(w)
						Expect(response.Ticker).To(HaveLen(10))
					})
				})

				Context("Bad Values", func() {
					It("should return 400 for missing ticker", func() {
						journal := barkat.Journal{Ticker: "", Sequence: "MWD", Type: "REJECTED", Status: "FAIL"}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, journal)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Ticker", "required")
					})

					It("should return 400 for ticker exceeding max length (11)", func() {
						journal := barkat.Journal{Ticker: "12345678901", Sequence: "MWD", Type: "REJECTED", Status: "FAIL"}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, journal)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Ticker", "max (10)")
					})

					It("should return 400 for lowercase ticker (PRD: uppercase only)", func() {
						journal := barkat.Journal{Ticker: "grse", Sequence: "MWD", Type: "REJECTED", Status: "FAIL", Images: standardImages}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, journal)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Ticker", "ticker")
					})
				})
			})

			Context("Sequence Field", func() {
				Context("Allowed Values", func() {
					It("should accept sequence = MWD", func() {
						journal := barkat.Journal{Ticker: "PDSL", Sequence: "MWD", Type: "SET", Status: "TAKEN", Images: standardImages}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, journal)
						router.ServeHTTP(w, req)
						response := decodeCreateJournalResponse(w)
						Expect(response.Sequence).To(Equal("MWD"))
					})

					It("should accept sequence = YR", func() {
						journal := barkat.Journal{Ticker: "SNF", Sequence: "YR", Type: "RESULT", Status: "SUCCESS", Images: standardImages}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, journal)
						router.ServeHTTP(w, req)
						response := decodeCreateJournalResponse(w)
						Expect(response.Sequence).To(Equal("YR"))
					})
				})

				Context("Bad Values", func() {
					It("should return 400 for missing sequence", func() {
						journal := barkat.Journal{Ticker: "GRSE", Sequence: "", Type: "REJECTED", Status: "FAIL"}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, journal)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Sequence", "required")
					})

					It("should return 400 for invalid sequence (lowercase)", func() {
						journal := barkat.Journal{Ticker: "INFY", Sequence: "mwd", Type: "SET", Status: "TAKEN"}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, journal)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Sequence", "oneof")
					})

					It("should return 400 for invalid sequence (unsupported)", func() {
						journal := barkat.Journal{Ticker: "WIPRO", Sequence: "QUARTERLY", Type: "SET", Status: "TAKEN"}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, journal)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Sequence", "oneof")
					})
				})
			})

			Context("Type Field", func() {
				Context("Allowed Values", func() {
					It("should accept type = REJECTED", func() {
						journal := barkat.Journal{Ticker: "TCS", Sequence: "MWD", Type: "REJECTED", Status: "FAIL", Images: standardImages}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, journal)
						router.ServeHTTP(w, req)
						response := decodeCreateJournalResponse(w)
						Expect(response.Type).To(Equal("REJECTED"))
					})

					It("should accept type = RESULT", func() {
						journal := barkat.Journal{Ticker: "INFY", Sequence: "YR", Type: "RESULT", Status: "SUCCESS", Images: standardImages}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, journal)
						router.ServeHTTP(w, req)
						response := decodeCreateJournalResponse(w)
						Expect(response.Type).To(Equal("RESULT"))
					})

					It("should accept type = SET", func() {
						journal := barkat.Journal{
							Ticker:   "RELIANCE",
							Sequence: "MWD",
							Type:     "SET",
							Status:   "RUNNING",
							Images:   standardImages,
						}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, journal)
						router.ServeHTTP(w, req)
						response := decodeCreateJournalResponse(w)
						Expect(response.Type).To(Equal("SET"))
					})
				})

				Context("Bad Values", func() {
					It("should return 400 for missing type", func() {
						journal := barkat.Journal{Ticker: "GRSE", Sequence: "MWD", Type: "", Status: "FAIL"}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, journal)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Type", "required")
					})

					It("should return 400 for invalid type (lowercase)", func() {
						journal := barkat.Journal{Ticker: "GRSE", Sequence: "MWD", Type: "rejected", Status: "FAIL"}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, journal)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Type", "oneof")
					})

					It("should return 400 for invalid type (unsupported)", func() {
						journal := barkat.Journal{Ticker: "HDFC", Sequence: "MWD", Type: "INVALID", Status: "TAKEN"}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, journal)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Type", "oneof")
					})
				})
			})

			Context("Status Field", func() {
				Context("Allowed Values - All 10 Status Values", func() {
					It("should accept status = SET", func() {
						journal := barkat.Journal{
							Ticker:   "T1",
							Sequence: "MWD",
							Type:     "SET",
							Status:   "SET",
							Images:   standardImages,
						}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, journal)
						router.ServeHTTP(w, req)
						response := decodeCreateJournalResponse(w)
						Expect(response.Status).To(Equal("SET"))
					})

					It("should accept status = RUNNING", func() {
						journal := barkat.Journal{Ticker: "T2", Sequence: "MWD", Type: "SET", Status: "RUNNING", Images: standardImages}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, journal)
						router.ServeHTTP(w, req)
						response := decodeCreateJournalResponse(w)
						Expect(response.Status).To(Equal("RUNNING"))
					})

					It("should accept status = DROPPED", func() {
						journal := barkat.Journal{Ticker: "T3", Sequence: "MWD", Type: "SET", Status: "DROPPED", Images: standardImages}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, journal)
						router.ServeHTTP(w, req)
						response := decodeCreateJournalResponse(w)
						Expect(response.Status).To(Equal("DROPPED"))
					})

					It("should accept status = TAKEN", func() {
						journal := barkat.Journal{Ticker: "T4", Sequence: "MWD", Type: "SET", Status: "TAKEN", Images: standardImages}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, journal)
						router.ServeHTTP(w, req)
						response := decodeCreateJournalResponse(w)
						Expect(response.Status).To(Equal("TAKEN"))
					})

					It("should accept status = REJECTED", func() {
						journal := barkat.Journal{
							Ticker:   "T5",
							Sequence: "MWD",
							Type:     "REJECTED",
							Status:   "REJECTED",
							Images:   standardImages,
						}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, journal)
						router.ServeHTTP(w, req)
						response := decodeCreateJournalResponse(w)
						Expect(response.Status).To(Equal("REJECTED"))
					})

					It("should accept status = SUCCESS", func() {
						journal := barkat.Journal{Ticker: "T6", Sequence: "YR", Type: "RESULT", Status: "SUCCESS", Images: standardImages}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, journal)
						router.ServeHTTP(w, req)
						response := decodeCreateJournalResponse(w)
						Expect(response.Status).To(Equal("SUCCESS"))
					})

					It("should accept status = FAIL", func() {
						journal := barkat.Journal{Ticker: "T7", Sequence: "YR", Type: "RESULT", Status: "FAIL", Images: standardImages}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, journal)
						router.ServeHTTP(w, req)
						response := decodeCreateJournalResponse(w)
						Expect(response.Status).To(Equal("FAIL"))
					})

					It("should accept status = MISSED", func() {
						journal := barkat.Journal{Ticker: "T8", Sequence: "YR", Type: "RESULT", Status: "MISSED", Images: standardImages}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, journal)
						router.ServeHTTP(w, req)
						response := decodeCreateJournalResponse(w)
						Expect(response.Status).To(Equal("MISSED"))
					})

					It("should accept status = JUST_LOSS", func() {
						journal := barkat.Journal{Ticker: "T9", Sequence: "YR", Type: "RESULT", Status: "JUST_LOSS", Images: standardImages}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, journal)
						router.ServeHTTP(w, req)
						response := decodeCreateJournalResponse(w)
						Expect(response.Status).To(Equal("JUST_LOSS"))
					})

					It("should accept status = BROKEN", func() {
						journal := barkat.Journal{Ticker: "T10", Sequence: "YR", Type: "RESULT", Status: "BROKEN", Images: standardImages}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, journal)
						router.ServeHTTP(w, req)
						response := decodeCreateJournalResponse(w)
						Expect(response.Status).To(Equal("BROKEN"))
					})
				})

				Context("Bad Values", func() {
					It("should return 400 for missing status", func() {
						journal := barkat.Journal{Ticker: "GRSE", Sequence: "MWD", Type: "REJECTED", Status: ""}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, journal)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Status", "required")
					})

					It("should return 400 for invalid status (lowercase)", func() {
						journal := barkat.Journal{Ticker: "GRSE", Sequence: "MWD", Type: "REJECTED", Status: "fail"}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, journal)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Status", "oneof")
					})

					It("should return 400 for invalid status (unsupported)", func() {
						journal := barkat.Journal{Ticker: "HDFC", Sequence: "MWD", Type: "SET", Status: "INVALID"}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, journal)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Status", "oneof")
					})
				})
			})

			Context("Images Field", func() {
				Context("Bad Values", func() {
					It("should return 400 for empty images (PRD: min 4 required)", func() {
						journal := barkat.Journal{
							Ticker:   "GRSE",
							Sequence: "MWD",
							Type:     "REJECTED",
							Status:   "FAIL",
							Images:   []barkat.Image{},
						}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, journal)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Images", "required")
					})

					It("should return 400 for insufficient images < 4 (PRD: min 4 required)", func() {
						journal := barkat.Journal{
							Ticker:   "GRSE",
							Sequence: "MWD",
							Type:     "REJECTED",
							Status:   "FAIL",
							Images: []barkat.Image{
								{Timeframe: "DL", FileName: "test-dl.png"},
								{Timeframe: "WK", FileName: "test-wk.png"},
								{Timeframe: "MN", FileName: "test-mn.png"},
							},
						}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, journal)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Images", "min")
					})

					It("should return 400 for excessive images > 6 (PRD: max 6 allowed)", func() {
						journal := barkat.Journal{
							Ticker:   "GRSE",
							Sequence: "MWD",
							Type:     "REJECTED",
							Status:   "FAIL",
							Images: []barkat.Image{
								{Timeframe: "DL", FileName: "test-dl.png"},
								{Timeframe: "WK", FileName: "test-wk.png"},
								{Timeframe: "MN", FileName: "test-mn.png"},
								{Timeframe: "TMN", FileName: "test-tmn.png"},
								{Timeframe: "SMN", FileName: "test-smn.png"},
								{Timeframe: "YR", FileName: "test-yr.png"},
								{Timeframe: "DL", FileName: "test-dl.png"},
							},
						}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, journal)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Images", "max")
					})

					It("should return 400 for invalid timeframe (PRD: must be DL,WK,MN,TMN,SMN,YR)", func() {
						journal := barkat.Journal{
							Ticker:   "GRSE",
							Sequence: "MWD",
							Type:     "REJECTED",
							Status:   "FAIL",
							Images: []barkat.Image{
								{Timeframe: "INVALID", FileName: "test-invalid.png"},
								{Timeframe: "WK", FileName: "test-wk.png"},
								{Timeframe: "MN", FileName: "test-mn.png"},
								{Timeframe: "TMN", FileName: "test-tmn.png"},
							},
						}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, journal)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Timeframe", "oneof")
					})

					It("should return 400 for duplicate timeframes (PRD: unique timeframes required)", func() {
						journal := barkat.Journal{
							Ticker:   "GRSE",
							Sequence: "MWD",
							Type:     "REJECTED",
							Status:   "FAIL",
							Images: []barkat.Image{
								{Timeframe: "DL", FileName: "test-dl.png"},
								{Timeframe: "DL", FileName: "test-dl.png"},
								{Timeframe: "MN", FileName: "test-mn.png"},
								{Timeframe: "TMN", FileName: "test-tmn.png"},
							},
						}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, journal)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Images", "Duplicate")
					})
				})
			})

			Context("Notes Field", func() {
				var (
					journal barkat.Journal
				)

				BeforeEach(func() {
					journal = barkat.Journal{
						Ticker:   "GRSE",
						Sequence: "MWD",
						Type:     "REJECTED",
						Status:   "FAIL",
						Images:   standardImages,
					}
				})

				Context("Bad Values", func() {
					It("should return 400 for missing note status (PRD: status required)", func() {
						journal.Notes = []barkat.Note{
							{Status: "", Content: "Note without status", Format: "MARKDOWN"},
						}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, journal)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Status", "required")
					})

					It("should return 400 for missing note content (PRD: content required)", func() {
						journal.Notes = []barkat.Note{
							{Status: "SET", Content: "", Format: "MARKDOWN"},
						}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, journal)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Content", "required")
					})

					It("should return 400 for multiple notes > 1 (PRD: max 1 at create)", func() {
						journal.Notes = []barkat.Note{
							{Status: "SET", Content: "First note", Format: "MARKDOWN"},
							{Status: "RUNNING", Content: "Second note", Format: "PLAINTEXT"},
						}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, journal)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Notes", "max")
					})

					It("should return 400 for invalid note format (PRD: must be MARKDOWN or PLAINTEXT)", func() {
						journal.Notes = []barkat.Note{
							{Status: "SET", Content: "Note with invalid format", Format: "invalid"},
						}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, journal)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Format", "oneof")
					})

					It("should return 400 for invalid note status (PRD: must match entry status enum)", func() {
						journal.Notes = []barkat.Note{
							{Status: "INVALID", Content: "Note with invalid status", Format: "MARKDOWN"},
						}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, journal)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Status", "oneof")
					})

					It("should return 400 for note content exceeding max length (PRD: max 2000 chars)", func() {
						longContent := string(make([]byte, 2001))
						for i := range longContent {
							longContent = longContent[:i] + "a" + longContent[i+1:]
						}
						journal.Notes = []barkat.Note{
							{Status: "SET", Content: longContent, Format: "MARKDOWN"},
						}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, journal)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Content", "max")
					})
				})
			})

			Context("Tags Field", func() {
				Context("Allowed Values", func() {
					It("should accept tag with override field (PRD: optional override max 50 chars)", func() {
						override := "loc"
						journal := barkat.Journal{
							Ticker:   "GRSE",
							Sequence: "MWD",
							Type:     "REJECTED",
							Status:   "FAIL",
							Images:   standardImages,
							Tags: []barkat.Tag{
								{Tag: "dep", Type: "REASON", Override: &override},
							},
						}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, journal)
						router.ServeHTTP(w, req)
						var envelopeResponse common.Envelope[barkat.Journal]
						util.AssertSuccess(w, http.StatusCreated, &envelopeResponse)
						Expect(envelopeResponse.Status).To(Equal(common.EnvelopeSuccess))
						Expect(envelopeResponse.Data.Tags).To(HaveLen(1))
						Expect(envelopeResponse.Data.Tags[0].Override).ToNot(BeNil())
						Expect(*envelopeResponse.Data.Tags[0].Override).To(Equal("loc"))
					})

					It("should accept tag type = REASON (PRD: uppercase)", func() {
						journal := barkat.Journal{
							Ticker:   "GRSE",
							Sequence: "MWD",
							Type:     "REJECTED",
							Status:   "FAIL",
							Images:   standardImages,
							Tags: []barkat.Tag{
								{Tag: "oe", Type: "REASON"},
							},
						}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, journal)
						router.ServeHTTP(w, req)
						var envelopeResponse common.Envelope[barkat.Journal]
						util.AssertSuccess(w, http.StatusCreated, &envelopeResponse)
					})

					It("should accept tag type = MANAGEMENT (PRD: uppercase)", func() {
						journal := barkat.Journal{
							Ticker:   "GRSE",
							Sequence: "MWD",
							Type:     "REJECTED",
							Status:   "FAIL",
							Images:   standardImages,
							Tags: []barkat.Tag{
								{Tag: "sl", Type: "MANAGEMENT"},
							},
						}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, journal)
						router.ServeHTTP(w, req)
						var envelopeResponse common.Envelope[barkat.Journal]
						util.AssertSuccess(w, http.StatusCreated, &envelopeResponse)
					})
				})

				Context("Bad Values", func() {
					It("should return 400 for missing tag name (PRD: tag required)", func() {
						journal := barkat.Journal{
							Ticker:   "GRSE",
							Sequence: "MWD",
							Type:     "REJECTED",
							Status:   "FAIL",
							Images:   standardImages,
							Tags: []barkat.Tag{
								{Tag: "", Type: "REASON"},
							},
						}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, journal)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Tag", "required")
					})

					It("should return 400 for missing tag type (PRD: type required)", func() {
						journal := barkat.Journal{
							Ticker:   "GRSE",
							Sequence: "MWD",
							Type:     "REJECTED",
							Status:   "FAIL",
							Images:   standardImages,
							Tags: []barkat.Tag{
								{Tag: "oe", Type: ""},
							},
						}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, journal)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Type", "required")
					})

					It("should return 400 for invalid tag type (PRD: must be REASON or MANAGEMENT)", func() {
						journal := barkat.Journal{
							Ticker:   "GRSE",
							Sequence: "MWD",
							Type:     "REJECTED",
							Status:   "FAIL",
							Images:   standardImages,
							Tags: []barkat.Tag{
								{Tag: "test", Type: "invalid"},
							},
						}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, journal)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Type", "oneof")
					})

					It("should return 400 for tag exceeding max length (PRD: max 10 chars)", func() {
						journal := barkat.Journal{
							Ticker:   "GRSE",
							Sequence: "MWD",
							Type:     "REJECTED",
							Status:   "FAIL",
							Images:   standardImages,
							Tags: []barkat.Tag{
								{Tag: "verylongtag1", Type: "REASON"},
							},
						}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, journal)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Tag", "max")
					})

					It("should return 400 for override exceeding max length (PRD: max 5 chars)", func() {
						longOverride := "toolong"
						journal := barkat.Journal{
							Ticker:   "GRSE",
							Sequence: "MWD",
							Type:     "REJECTED",
							Status:   "FAIL",
							Images:   standardImages,
							Tags: []barkat.Tag{
								{Tag: "dep", Type: "REASON", Override: &longOverride},
							},
						}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, journal)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Override", "max")
					})

					It("should return 400 for invalid tag format (PRD: alphanumeric with hyphens)", func() {
						journal := barkat.Journal{
							Ticker:   "GRSE",
							Sequence: "MWD",
							Type:     "REJECTED",
							Status:   "FAIL",
							Images:   standardImages,
							Tags: []barkat.Tag{
								{Tag: "bad@tag", Type: "REASON"},
							},
						}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, journal)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Tag", "tag")
					})

					It("should return 400 for invalid override format (PRD: letters only)", func() {
						invalidOverride := "a-b"
						journal := barkat.Journal{
							Ticker:   "GRSE",
							Sequence: "MWD",
							Type:     "REJECTED",
							Status:   "FAIL",
							Images:   standardImages,
							Tags: []barkat.Tag{
								{Tag: "dep", Type: "REASON", Override: &invalidOverride},
							},
						}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, journal)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Override", "override")
					})

					It("should return 400 for exceeding max tags (PRD: max 10)", func() {
						journal := barkat.Journal{
							Ticker:   "GRSE",
							Sequence: "MWD",
							Type:     "REJECTED",
							Status:   "FAIL",
							Images:   standardImages,
							Tags: []barkat.Tag{
								{Tag: "t1", Type: "REASON"},
								{Tag: "t2", Type: "REASON"},
								{Tag: "t3", Type: "REASON"},
								{Tag: "t4", Type: "REASON"},
								{Tag: "t5", Type: "REASON"},
								{Tag: "t6", Type: "REASON"},
								{Tag: "t7", Type: "REASON"},
								{Tag: "t8", Type: "REASON"},
								{Tag: "t9", Type: "REASON"},
								{Tag: "t10", Type: "REASON"},
								{Tag: "t11", Type: "REASON"},
							},
						}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, journal)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Tags", "max")
					})
				})
			})
		})

		Context("Errors", func() {
			It("should return 400 for invalid JSON", func() {
				req, w = util.CreateTestRequest("POST", barkat.JournalEntries, []byte("invalid json"))
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})

			It("should return 400 for empty request body", func() {
				req, w = util.CreateTestRequest("POST", barkat.JournalEntries, []byte(""))
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})

			It("should return 400 for null request body", func() {
				req, w = util.CreateTestRequest("POST", barkat.JournalEntries, []byte("null"))
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})
		})
	})

	Describe("DELETE /v1/journal/{id} - Delete Entry", func() {
		var createdJournal barkat.Journal

		BeforeEach(func() {
			journal := barkat.Journal{
				Ticker:   "GRSE",
				Sequence: "MWD",
				Type:     "REJECTED",
				Status:   "FAIL",
				Images:   standardImages,
			}
			Expect(journalMgr.CreateJournal(testCtx, &journal)).To(Succeed())
			createdJournal = journal
		})

		Context("Happy Path", func() {
			Context("with valid entry ID", func() {
				BeforeEach(func() {
					req, w = util.CreateTestRequest("DELETE", barkat.JournalEntries+"/"+createdJournal.ExternalID, nil)
					router.ServeHTTP(w, req)
				})

				It("should return 204 No Content", func() {
					Expect(w.Code).To(Equal(http.StatusNoContent))
				})

				It("should actually delete the journal", func() {
					_, err := journalMgr.GetJournal(testCtx, createdJournal.ExternalID)
					Expect(err).To(HaveOccurred())
				})
			})
		})

		Context("Field Validations", func() {
			Context("Entry ID Field", func() {
				Context("Bad Values", func() {
					It("should return 404 for non-existent entry ID", func() {
						req, w = util.CreateTestRequest("DELETE", barkat.JournalEntries+"/nonexistent-id", nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusNotFound))
					})

					It("should return 404 for malformed UUID", func() {
						req, w = util.CreateTestRequest("DELETE", barkat.JournalEntries+"/invalid-uuid-format", nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusNotFound))
					})

					It("should return 404 for valid UUID format but non-existent", func() {
						req, w = util.CreateTestRequest("DELETE", barkat.JournalEntries+"/550e8400-e29b-41d4-a716-446655440000", nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusNotFound))
					})
				})
			})
		})

		Context("Errors", func() {
			// No server error scenarios for DELETE currently
		})
	})

	Describe("PATCH /v1/journal/{id} - Update Review Status", func() {
		// Common test payloads for JournalReviewUpdate to reduce duplication
		// Kept close to the API tests that use them per AGENTS.md guidelines
		var (
			validReviewDate  = civil.Date{Year: 2024, Month: 1, Day: 16}
			futureReviewDate = civil.Date{Year: 2099, Month: 12, Day: 31}

			markReviewedPayload   = barkat.JournalReviewUpdate{ReviewedAt: &validReviewDate}
			markUnreviewedPayload = barkat.JournalReviewUpdate{ReviewedAt: nil}
			futureDatePayload     = barkat.JournalReviewUpdate{ReviewedAt: &futureReviewDate}
		)

		var createdJournal barkat.Journal

		BeforeEach(func() {
			journal := barkat.Journal{
				Ticker:   "TCS",
				Sequence: "MWD",
				Type:     "REJECTED",
				Status:   "FAIL",
				Images:   standardImages,
			}
			Expect(journalMgr.CreateJournal(testCtx, &journal)).To(Succeed())
			createdJournal = journal
		})

		Context("Happy Path", func() {
			Context("mark as reviewed", func() {
				BeforeEach(func() {
					req, w = util.CreateTestRequest("PATCH", barkat.JournalEntries+"/"+createdJournal.ExternalID, markReviewedPayload)
					router.ServeHTTP(w, req)
				})

				It("should return 200 OK", func() {
					Expect(w.Code).To(Equal(http.StatusOK))
				})

				It("should return success envelope", func() {
					var envelope common.Envelope[barkat.UpdateJournalStatusResponse]
					util.AssertSuccess(w, http.StatusOK, &envelope)
					Expect(envelope.Data).ToNot(BeNil())
				})

				It("should return correct response fields", func() {
					var envelope common.Envelope[barkat.UpdateJournalStatusResponse]
					util.AssertSuccess(w, http.StatusOK, &envelope)

					data := envelope.Data
					Expect(data.ID).To(Equal(createdJournal.ExternalID))
					Expect(data.ReviewedAt).ToNot(BeNil())
				})

				It("should actually update the journal in database", func() {
					updatedJournal, err := journalMgr.GetJournal(testCtx, createdJournal.ExternalID)
					Expect(err).ToNot(HaveOccurred())
					Expect(updatedJournal.ReviewedAt).ToNot(BeNil())
				})
			})

			Context("mark as unreviewed", func() {
				BeforeEach(func() {
					// First mark as reviewed
					req, w = util.CreateTestRequest("PATCH", barkat.JournalEntries+"/"+createdJournal.ExternalID, markReviewedPayload)
					router.ServeHTTP(w, req)

					// Then mark as unreviewed
					req, w = util.CreateTestRequest("PATCH", barkat.JournalEntries+"/"+createdJournal.ExternalID, markUnreviewedPayload)
					router.ServeHTTP(w, req)
				})

				It("should return 200 OK", func() {
					Expect(w.Code).To(Equal(http.StatusOK))
				})

				It("should return reviewed_at as null", func() {
					var envelope common.Envelope[barkat.UpdateJournalStatusResponse]
					util.AssertSuccess(w, http.StatusOK, &envelope)

					data := envelope.Data
					Expect(data.ReviewedAt).To(BeNil())
				})

				It("should actually clear reviewed_at in database", func() {
					updatedJournal, err := journalMgr.GetJournal(testCtx, createdJournal.ExternalID)
					Expect(err).ToNot(HaveOccurred())
					Expect(updatedJournal.ReviewedAt).To(BeNil())
				})
			})

			Context("idempotent operations", func() {
				BeforeEach(func() {
					// Mark as reviewed twice
					req, w = util.CreateTestRequest("PATCH", barkat.JournalEntries+"/"+createdJournal.ExternalID, markReviewedPayload)
					router.ServeHTTP(w, req)

					// Second time should be idempotent
					req, w = util.CreateTestRequest("PATCH", barkat.JournalEntries+"/"+createdJournal.ExternalID, markReviewedPayload)
					router.ServeHTTP(w, req)
				})

				It("should return 200 OK for idempotent operation", func() {
					Expect(w.Code).To(Equal(http.StatusOK))
				})

				It("should maintain reviewed_at timestamp", func() {
					var envelope common.Envelope[barkat.UpdateJournalStatusResponse]
					util.AssertSuccess(w, http.StatusOK, &envelope)

					data := envelope.Data
					Expect(data.ReviewedAt).ToNot(BeNil())
				})
			})
		})

		Context("Field Validations", func() {
			Context("reviewed_at field", func() {
				Context("Allowed Values", func() {
					It("should accept valid date string", func() {
						req, w = util.CreateTestRequest("PATCH", barkat.JournalEntries+"/"+createdJournal.ExternalID, markReviewedPayload)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusOK))
					})

					It("should accept empty string for unreviewed", func() {
						req, w = util.CreateTestRequest("PATCH", barkat.JournalEntries+"/"+createdJournal.ExternalID, markUnreviewedPayload)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusOK))
					})
				})

				Context("Bad Values", func() {
					It("should return 400 for missing reviewed_at field", func() {
						payload := `{}`
						req, w = util.CreateTestRequest("PATCH", barkat.JournalEntries+"/"+createdJournal.ExternalID, payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "reviewed_at", "Violates 'required'")
					})

					It("should return 400 for null reviewed_at field", func() {
						payload := `{"reviewed_at": null}`
						req, w = util.CreateTestRequest("PATCH", barkat.JournalEntries+"/"+createdJournal.ExternalID, payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "reviewed_at", "Violates 'required'")
					})

					It("should return 400 for invalid date format", func() {
						payload := `{"reviewed_at": "2024-01-16T14:30:00Z"}`
						req, w = util.CreateTestRequest("PATCH", barkat.JournalEntries+"/"+createdJournal.ExternalID, payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "reviewed_at", "Violates")
					})

					It("should return 400 for future date (PRD 3.1 business rule)", func() {
						req, w = util.CreateTestRequest("PATCH", barkat.JournalEntries+"/"+createdJournal.ExternalID, futureDatePayload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "reviewed_at", "not_future")
					})

					It("should return 400 for non-string types", func() {
						payload := `{"reviewed_at": true}`
						req, w = util.CreateTestRequest("PATCH", barkat.JournalEntries+"/"+createdJournal.ExternalID, payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "reviewed_at", "Violates")
					})
				})
			})

			Context("Entry ID Field", func() {
				Context("Bad Values", func() {
					It("should return 404 for non-existent entry ID", func() {
						req, w = util.CreateTestRequest("PATCH", barkat.JournalEntries+"/nonexistent-id", markReviewedPayload)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusNotFound))
					})

					It("should return 404 for invalid UUID format", func() {
						req, w = util.CreateTestRequest("PATCH", barkat.JournalEntries+"/invalid-uuid-format", markReviewedPayload)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusNotFound))
					})

					It("should return 404 for valid UUID format but non-existent", func() {
						req, w = util.CreateTestRequest("PATCH", barkat.JournalEntries+"/550e8400-e29b-41d4-a716-446655440000", markReviewedPayload)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusNotFound))
					})
				})
			})
		})

		Context("Errors", func() {
			Context("Malformed JSON", func() {
				It("should return 400 for empty body", func() {
					req, w = util.CreateTestRequest("PATCH", barkat.JournalEntries+"/"+createdJournal.ExternalID, "")
					router.ServeHTTP(w, req)
					Expect(w.Code).To(Equal(http.StatusBadRequest))
				})

				It("should return 400 for invalid JSON syntax", func() {
					payload := `{"reviewed_at": "2024-01-16"` // Missing closing brace
					req, w = util.CreateTestRequest("PATCH", barkat.JournalEntries+"/"+createdJournal.ExternalID, payload)
					router.ServeHTTP(w, req)
					Expect(w.Code).To(Equal(http.StatusBadRequest))
				})
			})

			Context("Content-Type Issues", func() {
				It("should return 400 for wrong content type", func() {
					req, w = util.CreateTestRequest("PATCH", barkat.JournalEntries+"/"+createdJournal.ExternalID, markReviewedPayload)
					req.Header.Set("Content-Type", "text/plain")
					router.ServeHTTP(w, req)
					Expect(w.Code).To(Equal(http.StatusBadRequest))
				})
			})
		})
	})
})
