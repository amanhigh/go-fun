package handler_test

import (
	"context"
	"net/http"
	"net/http/httptest"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/components/kohan/core"
	"github.com/amanhigh/go-fun/components/kohan/handler"
	"github.com/amanhigh/go-fun/components/kohan/manager"
	"github.com/amanhigh/go-fun/components/kohan/repository"
	"github.com/amanhigh/go-fun/models/barkat"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"
)

const (
	// Base URL for journal entries API
	JournalEntriesBaseURL = "/v1/journal-entries"
)

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
		entryMgr       manager.JournalManager
		req            *http.Request
		w              *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		var err error
		db, err = core.CreateTestBarkatDB()
		Expect(err).ToNot(HaveOccurred())

		entryRepo := repository.NewJournalRepository(db)
		entryMgr = manager.NewJournalManager(entryRepo)
		journalHandler = handler.NewJournalHandler(entryMgr)

		router = util.CreateTestGinRouter()
		v1 := router.Group("/v1")
		handler.SetupJournalEntryRoutes(v1, journalHandler)
	})

	AfterEach(func() {
		sqlDB, err := db.DB()
		Expect(err).ToNot(HaveOccurred())
		sqlDB.Close()
	})

	Describe("POST /v1/journal-entries - Create Entry", func() {
		Context("Happy Path", func() {
			Context("with minimal valid entry", func() {
				var response barkat.Entry

				BeforeEach(func() {
					entry := barkat.Entry{
						Ticker:   "GRSE",
						Sequence: "MWD",
						Type:     "REJECTED",
						Status:   "FAIL",
					}
					req, w = util.CreateTestRequest("POST", JournalEntriesBaseURL, entry)
					router.ServeHTTP(w, req)
				})

				It("should return 201 Created", func() {
					Expect(w.Code).To(Equal(http.StatusCreated))
				})

				It("should return created entry with ID", func() {
					util.AssertJSONAndStatus(w, http.StatusCreated, &response)
					Expect(response.ID).ToNot(BeEmpty())
				})

				It("should preserve all input fields", func() {
					util.AssertJSONAndStatus(w, http.StatusCreated, &response)
					Expect(response.Ticker).To(Equal("GRSE"))
					Expect(response.Sequence).To(Equal("MWD"))
					Expect(response.Type).To(Equal("REJECTED"))
					Expect(response.Status).To(Equal("FAIL"))
				})

				It("should set created_at timestamp", func() {
					util.AssertJSONAndStatus(w, http.StatusCreated, &response)
					Expect(response.CreatedAt).ToNot(BeZero())
				})

				It("should persist entry to database", func() {
					util.AssertJSONAndStatus(w, http.StatusCreated, &response)
					dbEntry, err := entryMgr.GetEntry(testCtx, response.ID)
					Expect(err).ToNot(HaveOccurred())
					Expect(dbEntry.ID).To(Equal(response.ID))
					Expect(dbEntry.Ticker).To(Equal("GRSE"))
				})
			})

			Context("with complete valid entry including images, tags, and notes", func() {
				var response barkat.Entry

				BeforeEach(func() {
					entry := barkat.Entry{
						Ticker:   "RELIANCE",
						Sequence: "YR",
						Type:     "REJECTED",
						Status:   "RUNNING",
						Images: []barkat.Image{
							{Timeframe: "DL"},
							{Timeframe: "WK"},
							{Timeframe: "MN"},
							{Timeframe: "TMN"},
						},
						Tags: []barkat.Tag{
							{Tag: "oe", Type: "reason"},
						},
						Notes: []barkat.Note{
							{Status: "set", Content: "Strong OE at weekly level.", Format: "markdown"},
						},
					}
					req, w = util.CreateTestRequest("POST", JournalEntriesBaseURL, entry)
					router.ServeHTTP(w, req)
				})

				It("should create entry with all associations", func() {
					util.AssertJSONAndStatus(w, http.StatusCreated, &response)
					Expect(response.Ticker).To(Equal("RELIANCE"))
					Expect(response.Sequence).To(Equal("YR"))
					Expect(response.Type).To(Equal("REJECTED"))
					Expect(response.Status).To(Equal("RUNNING"))
					Expect(response.Images).To(HaveLen(4))
					Expect(response.Tags).To(HaveLen(1))
					Expect(response.Notes).To(HaveLen(1))
				})
			})
		})

		Context("Field Validations", func() {
			Context("Ticker Field", func() {
				Context("Allowed Values", func() {
					It("should accept ticker with numbers", func() {
						entry := barkat.Entry{Ticker: "GRSE123", Sequence: "MWD", Type: "REJECTED", Status: "FAIL"}
						req, w = util.CreateTestRequest("POST", JournalEntriesBaseURL, entry)
						router.ServeHTTP(w, req)
						var response barkat.Entry
						util.AssertJSONAndStatus(w, http.StatusCreated, &response)
						Expect(response.Ticker).To(Equal("GRSE123"))
					})

					It("should accept ticker with hyphen", func() {
						entry := barkat.Entry{Ticker: "GRSE-NSE", Sequence: "MWD", Type: "REJECTED", Status: "FAIL"}
						req, w = util.CreateTestRequest("POST", JournalEntriesBaseURL, entry)
						router.ServeHTTP(w, req)
						var response barkat.Entry
						util.AssertJSONAndStatus(w, http.StatusCreated, &response)
						Expect(response.Ticker).To(Equal("GRSE-NSE"))
					})

					It("should accept ticker with dot suffix", func() {
						entry := barkat.Entry{Ticker: "TCS.NS", Sequence: "YR", Type: "SET", Status: "RUNNING"}
						req, w = util.CreateTestRequest("POST", JournalEntriesBaseURL, entry)
						router.ServeHTTP(w, req)
						var response barkat.Entry
						util.AssertJSONAndStatus(w, http.StatusCreated, &response)
						Expect(response.Ticker).To(Equal("TCS.NS"))
					})

					It("should accept ticker at max length (10)", func() {
						entry := barkat.Entry{Ticker: "1234567890", Sequence: "MWD", Type: "REJECTED", Status: "FAIL"}
						req, w = util.CreateTestRequest("POST", JournalEntriesBaseURL, entry)
						router.ServeHTTP(w, req)
						var response barkat.Entry
						util.AssertJSONAndStatus(w, http.StatusCreated, &response)
						Expect(response.Ticker).To(HaveLen(10))
					})
				})

				Context("Bad Values", func() {
					It("should return 400 for missing ticker", func() {
						entry := barkat.Entry{Ticker: "", Sequence: "MWD", Type: "REJECTED", Status: "FAIL"}
						req, w = util.CreateTestRequest("POST", JournalEntriesBaseURL, entry)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusBadRequest))
						var errorResponse map[string]interface{}
						util.AssertJSONAndStatus(w, http.StatusBadRequest, &errorResponse)
						Expect(errorResponse["message"]).To(ContainSubstring("required"))
					})

					It("should return 400 for ticker exceeding max length (11)", func() {
						entry := barkat.Entry{Ticker: "12345678901", Sequence: "MWD", Type: "REJECTED", Status: "FAIL"}
						req, w = util.CreateTestRequest("POST", JournalEntriesBaseURL, entry)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusBadRequest))
						var errorResponse map[string]interface{}
						util.AssertJSONAndStatus(w, http.StatusBadRequest, &errorResponse)
						Expect(errorResponse["message"]).To(ContainSubstring("max (10)"))
					})
				})
			})

			Context("Sequence Field", func() {
				Context("Allowed Values", func() {
					It("should accept sequence = MWD", func() {
						entry := barkat.Entry{Ticker: "PDSL", Sequence: "MWD", Type: "SET", Status: "TAKEN"}
						req, w = util.CreateTestRequest("POST", JournalEntriesBaseURL, entry)
						router.ServeHTTP(w, req)
						var response barkat.Entry
						util.AssertJSONAndStatus(w, http.StatusCreated, &response)
						Expect(response.Sequence).To(Equal("MWD"))
					})

					It("should accept sequence = YR", func() {
						entry := barkat.Entry{Ticker: "SNF", Sequence: "YR", Type: "RESULT", Status: "SUCCESS"}
						req, w = util.CreateTestRequest("POST", JournalEntriesBaseURL, entry)
						router.ServeHTTP(w, req)
						var response barkat.Entry
						util.AssertJSONAndStatus(w, http.StatusCreated, &response)
						Expect(response.Sequence).To(Equal("YR"))
					})
				})

				Context("Bad Values", func() {
					It("should return 400 for missing sequence", func() {
						entry := barkat.Entry{Ticker: "GRSE", Sequence: "", Type: "REJECTED", Status: "FAIL"}
						req, w = util.CreateTestRequest("POST", JournalEntriesBaseURL, entry)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusBadRequest))
						var errorResponse map[string]interface{}
						util.AssertJSONAndStatus(w, http.StatusBadRequest, &errorResponse)
						Expect(errorResponse["message"]).To(ContainSubstring("required"))
					})

					It("should return 400 for invalid sequence (lowercase)", func() {
						entry := barkat.Entry{Ticker: "INFY", Sequence: "mwd", Type: "SET", Status: "TAKEN"}
						req, w = util.CreateTestRequest("POST", JournalEntriesBaseURL, entry)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusBadRequest))
						var errorResponse map[string]interface{}
						util.AssertJSONAndStatus(w, http.StatusBadRequest, &errorResponse)
						Expect(errorResponse["message"]).To(ContainSubstring("oneof"))
					})

					It("should return 400 for invalid sequence (unsupported)", func() {
						entry := barkat.Entry{Ticker: "WIPRO", Sequence: "QUARTERLY", Type: "SET", Status: "TAKEN"}
						req, w = util.CreateTestRequest("POST", JournalEntriesBaseURL, entry)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusBadRequest))
						var errorResponse map[string]interface{}
						util.AssertJSONAndStatus(w, http.StatusBadRequest, &errorResponse)
						Expect(errorResponse["message"]).To(ContainSubstring("oneof"))
					})
				})
			})

			Context("Type Field", func() {
				Context("Allowed Values", func() {
					It("should accept type = REJECTED", func() {
						entry := barkat.Entry{Ticker: "TCS", Sequence: "MWD", Type: "REJECTED", Status: "FAIL"}
						req, w = util.CreateTestRequest("POST", JournalEntriesBaseURL, entry)
						router.ServeHTTP(w, req)
						var response barkat.Entry
						util.AssertJSONAndStatus(w, http.StatusCreated, &response)
						Expect(response.Type).To(Equal("REJECTED"))
					})

					It("should accept type = RESULT", func() {
						entry := barkat.Entry{Ticker: "INFY", Sequence: "YR", Type: "RESULT", Status: "SUCCESS"}
						req, w = util.CreateTestRequest("POST", JournalEntriesBaseURL, entry)
						router.ServeHTTP(w, req)
						var response barkat.Entry
						util.AssertJSONAndStatus(w, http.StatusCreated, &response)
						Expect(response.Type).To(Equal("RESULT"))
					})

					It("should accept type = SET", func() {
						entry := barkat.Entry{Ticker: "RELIANCE", Sequence: "MWD", Type: "SET", Status: "RUNNING"}
						req, w = util.CreateTestRequest("POST", JournalEntriesBaseURL, entry)
						router.ServeHTTP(w, req)
						var response barkat.Entry
						util.AssertJSONAndStatus(w, http.StatusCreated, &response)
						Expect(response.Type).To(Equal("SET"))
					})
				})

				Context("Bad Values", func() {
					It("should return 400 for missing type", func() {
						entry := barkat.Entry{Ticker: "GRSE", Sequence: "MWD", Type: "", Status: "FAIL"}
						req, w = util.CreateTestRequest("POST", JournalEntriesBaseURL, entry)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusBadRequest))
						var errorResponse map[string]interface{}
						util.AssertJSONAndStatus(w, http.StatusBadRequest, &errorResponse)
						Expect(errorResponse["message"]).To(ContainSubstring("required"))
					})

					It("should return 400 for invalid type (lowercase)", func() {
						entry := barkat.Entry{Ticker: "GRSE", Sequence: "MWD", Type: "rejected", Status: "FAIL"}
						req, w = util.CreateTestRequest("POST", JournalEntriesBaseURL, entry)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusBadRequest))
						var errorResponse map[string]interface{}
						util.AssertJSONAndStatus(w, http.StatusBadRequest, &errorResponse)
						Expect(errorResponse["message"]).To(ContainSubstring("oneof"))
					})

					It("should return 400 for invalid type (unsupported)", func() {
						entry := barkat.Entry{Ticker: "HDFC", Sequence: "MWD", Type: "INVALID", Status: "TAKEN"}
						req, w = util.CreateTestRequest("POST", JournalEntriesBaseURL, entry)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusBadRequest))
						var errorResponse map[string]interface{}
						util.AssertJSONAndStatus(w, http.StatusBadRequest, &errorResponse)
						Expect(errorResponse["message"]).To(ContainSubstring("oneof"))
					})
				})
			})

			Context("Status Field", func() {
				Context("Allowed Values - All 10 Status Values", func() {
					It("should accept status = SET", func() {
						entry := barkat.Entry{Ticker: "T1", Sequence: "MWD", Type: "SET", Status: "SET"}
						req, w = util.CreateTestRequest("POST", JournalEntriesBaseURL, entry)
						router.ServeHTTP(w, req)
						var response barkat.Entry
						util.AssertJSONAndStatus(w, http.StatusCreated, &response)
						Expect(response.Status).To(Equal("SET"))
					})

					It("should accept status = RUNNING", func() {
						entry := barkat.Entry{Ticker: "T2", Sequence: "MWD", Type: "SET", Status: "RUNNING"}
						req, w = util.CreateTestRequest("POST", JournalEntriesBaseURL, entry)
						router.ServeHTTP(w, req)
						var response barkat.Entry
						util.AssertJSONAndStatus(w, http.StatusCreated, &response)
						Expect(response.Status).To(Equal("RUNNING"))
					})

					It("should accept status = DROPPED", func() {
						entry := barkat.Entry{Ticker: "T3", Sequence: "MWD", Type: "SET", Status: "DROPPED"}
						req, w = util.CreateTestRequest("POST", JournalEntriesBaseURL, entry)
						router.ServeHTTP(w, req)
						var response barkat.Entry
						util.AssertJSONAndStatus(w, http.StatusCreated, &response)
						Expect(response.Status).To(Equal("DROPPED"))
					})

					It("should accept status = TAKEN", func() {
						entry := barkat.Entry{Ticker: "T4", Sequence: "MWD", Type: "SET", Status: "TAKEN"}
						req, w = util.CreateTestRequest("POST", JournalEntriesBaseURL, entry)
						router.ServeHTTP(w, req)
						var response barkat.Entry
						util.AssertJSONAndStatus(w, http.StatusCreated, &response)
						Expect(response.Status).To(Equal("TAKEN"))
					})

					It("should accept status = REJECTED", func() {
						entry := barkat.Entry{Ticker: "T5", Sequence: "MWD", Type: "REJECTED", Status: "REJECTED"}
						req, w = util.CreateTestRequest("POST", JournalEntriesBaseURL, entry)
						router.ServeHTTP(w, req)
						var response barkat.Entry
						util.AssertJSONAndStatus(w, http.StatusCreated, &response)
						Expect(response.Status).To(Equal("REJECTED"))
					})

					It("should accept status = SUCCESS", func() {
						entry := barkat.Entry{Ticker: "T6", Sequence: "YR", Type: "RESULT", Status: "SUCCESS"}
						req, w = util.CreateTestRequest("POST", JournalEntriesBaseURL, entry)
						router.ServeHTTP(w, req)
						var response barkat.Entry
						util.AssertJSONAndStatus(w, http.StatusCreated, &response)
						Expect(response.Status).To(Equal("SUCCESS"))
					})

					It("should accept status = FAIL", func() {
						entry := barkat.Entry{Ticker: "T7", Sequence: "YR", Type: "RESULT", Status: "FAIL"}
						req, w = util.CreateTestRequest("POST", JournalEntriesBaseURL, entry)
						router.ServeHTTP(w, req)
						var response barkat.Entry
						util.AssertJSONAndStatus(w, http.StatusCreated, &response)
						Expect(response.Status).To(Equal("FAIL"))
					})

					It("should accept status = MISSED", func() {
						entry := barkat.Entry{Ticker: "T8", Sequence: "YR", Type: "RESULT", Status: "MISSED"}
						req, w = util.CreateTestRequest("POST", JournalEntriesBaseURL, entry)
						router.ServeHTTP(w, req)
						var response barkat.Entry
						util.AssertJSONAndStatus(w, http.StatusCreated, &response)
						Expect(response.Status).To(Equal("MISSED"))
					})

					It("should accept status = JUST_LOSS", func() {
						entry := barkat.Entry{Ticker: "T9", Sequence: "YR", Type: "RESULT", Status: "JUST_LOSS"}
						req, w = util.CreateTestRequest("POST", JournalEntriesBaseURL, entry)
						router.ServeHTTP(w, req)
						var response barkat.Entry
						util.AssertJSONAndStatus(w, http.StatusCreated, &response)
						Expect(response.Status).To(Equal("JUST_LOSS"))
					})

					It("should accept status = BROKEN", func() {
						entry := barkat.Entry{Ticker: "T10", Sequence: "YR", Type: "RESULT", Status: "BROKEN"}
						req, w = util.CreateTestRequest("POST", JournalEntriesBaseURL, entry)
						router.ServeHTTP(w, req)
						var response barkat.Entry
						util.AssertJSONAndStatus(w, http.StatusCreated, &response)
						Expect(response.Status).To(Equal("BROKEN"))
					})
				})

				Context("Bad Values", func() {
					It("should return 400 for missing status", func() {
						entry := barkat.Entry{Ticker: "GRSE", Sequence: "MWD", Type: "REJECTED", Status: ""}
						req, w = util.CreateTestRequest("POST", JournalEntriesBaseURL, entry)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusBadRequest))
						var errorResponse map[string]interface{}
						util.AssertJSONAndStatus(w, http.StatusBadRequest, &errorResponse)
						Expect(errorResponse["message"]).To(ContainSubstring("required"))
					})

					It("should return 400 for invalid status (lowercase)", func() {
						entry := barkat.Entry{Ticker: "GRSE", Sequence: "MWD", Type: "REJECTED", Status: "fail"}
						req, w = util.CreateTestRequest("POST", JournalEntriesBaseURL, entry)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusBadRequest))
						var errorResponse map[string]interface{}
						util.AssertJSONAndStatus(w, http.StatusBadRequest, &errorResponse)
						Expect(errorResponse["message"]).To(ContainSubstring("oneof"))
					})

					It("should return 400 for invalid status (unsupported)", func() {
						entry := barkat.Entry{Ticker: "HDFC", Sequence: "MWD", Type: "SET", Status: "INVALID"}
						req, w = util.CreateTestRequest("POST", JournalEntriesBaseURL, entry)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusBadRequest))
						var errorResponse map[string]interface{}
						util.AssertJSONAndStatus(w, http.StatusBadRequest, &errorResponse)
						Expect(errorResponse["message"]).To(ContainSubstring("oneof"))
					})
				})
			})

			Context("Images Field", func() {
				Context("Bad Values", func() {
					It("should create entry with empty images (PRD: should return 400)", func() {
						// CURRENT: Creates entry successfully
						// PRD REQUIREMENT: Should return 400 Bad Request for empty images
						entry := barkat.Entry{
							Ticker:   "GRSE",
							Sequence: "MWD",
							Type:     "REJECTED",
							Status:   "FAIL",
							Images:   []barkat.Image{},
						}
						req, w = util.CreateTestRequest("POST", JournalEntriesBaseURL, entry)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusCreated))
						var response barkat.Entry
						util.AssertJSONAndStatus(w, http.StatusCreated, &response)
						Expect(response.Images).To(BeEmpty())
					})

					It("should create entry with insufficient images < 4 (PRD: should return 413)", func() {
						// CURRENT: Creates entry successfully
						// PRD REQUIREMENT: Should return 413 Payload Too Large for insufficient images
						entry := barkat.Entry{
							Ticker:   "GRSE",
							Sequence: "MWD",
							Type:     "REJECTED",
							Status:   "FAIL",
							Images: []barkat.Image{
								{Timeframe: "DL"},
								{Timeframe: "WK"},
								{Timeframe: "MN"},
							},
						}
						req, w = util.CreateTestRequest("POST", JournalEntriesBaseURL, entry)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusCreated))
						var response barkat.Entry
						util.AssertJSONAndStatus(w, http.StatusCreated, &response)
						Expect(response.Images).To(HaveLen(3))
					})

					It("should create entry with excessive images > 6 (PRD: should return 413)", func() {
						// CURRENT: Creates entry successfully
						// PRD REQUIREMENT: Should return 413 Payload Too Large for excessive images
						entry := barkat.Entry{
							Ticker:   "GRSE",
							Sequence: "MWD",
							Type:     "REJECTED",
							Status:   "FAIL",
							Images: []barkat.Image{
								{Timeframe: "DL"},
								{Timeframe: "WK"},
								{Timeframe: "MN"},
								{Timeframe: "TMN"},
								{Timeframe: "SMN"},
								{Timeframe: "YR"},
								{Timeframe: "DL"},
							},
						}
						req, w = util.CreateTestRequest("POST", JournalEntriesBaseURL, entry)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusCreated))
						var response barkat.Entry
						util.AssertJSONAndStatus(w, http.StatusCreated, &response)
						Expect(response.Images).To(HaveLen(7))
					})

					It("should create entry with invalid timeframe (PRD: should return 400)", func() {
						// CURRENT: Creates entry successfully (no validation on nested array elements)
						// PRD REQUIREMENT: Should return 400 Bad Request for invalid timeframe
						entry := barkat.Entry{
							Ticker:   "GRSE",
							Sequence: "MWD",
							Type:     "REJECTED",
							Status:   "FAIL",
							Images: []barkat.Image{
								{Timeframe: "INVALID"},
								{Timeframe: "WK"},
								{Timeframe: "MN"},
								{Timeframe: "TMN"},
							},
						}
						req, w = util.CreateTestRequest("POST", JournalEntriesBaseURL, entry)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusCreated))
						var response barkat.Entry
						util.AssertJSONAndStatus(w, http.StatusCreated, &response)
						Expect(response.Images).To(HaveLen(4))
					})
				})
			})

			Context("Notes Field", func() {
				Context("Bad Values", func() {
					It("should create entry with multiple notes > 1 (PRD: should return 413)", func() {
						// CURRENT: Creates entry successfully
						// PRD REQUIREMENT: Should return 413 Payload Too Large for multiple notes (max 1)
						entry := barkat.Entry{
							Ticker:   "GRSE",
							Sequence: "MWD",
							Type:     "REJECTED",
							Status:   "FAIL",
							Images: []barkat.Image{
								{Timeframe: "DL"},
								{Timeframe: "WK"},
								{Timeframe: "MN"},
								{Timeframe: "TMN"},
							},
							Notes: []barkat.Note{
								{Status: "set", Content: "First note", Format: "markdown"},
								{Status: "running", Content: "Second note", Format: "plaintext"},
							},
						}
						req, w = util.CreateTestRequest("POST", JournalEntriesBaseURL, entry)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusCreated))
						var response barkat.Entry
						util.AssertJSONAndStatus(w, http.StatusCreated, &response)
						Expect(response.Notes).To(HaveLen(2))
					})

					It("should create entry with invalid note format (PRD: should return 400)", func() {
						// CURRENT: Creates entry successfully (no validation on nested array elements)
						// PRD REQUIREMENT: Should return 400 Bad Request for invalid note format
						entry := barkat.Entry{
							Ticker:   "GRSE",
							Sequence: "MWD",
							Type:     "REJECTED",
							Status:   "FAIL",
							Images: []barkat.Image{
								{Timeframe: "DL"},
								{Timeframe: "WK"},
								{Timeframe: "MN"},
								{Timeframe: "TMN"},
							},
							Notes: []barkat.Note{
								{Status: "set", Content: "Note with invalid format", Format: "invalid"},
							},
						}
						req, w = util.CreateTestRequest("POST", JournalEntriesBaseURL, entry)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusCreated))
						var response barkat.Entry
						util.AssertJSONAndStatus(w, http.StatusCreated, &response)
						Expect(response.Notes).To(HaveLen(1))
					})
				})
			})

			Context("Tags Field", func() {
				Context("Bad Values", func() {
					It("should create entry with invalid tag type (PRD: should return 400)", func() {
						// CURRENT: Creates entry successfully (no validation on nested array elements)
						// PRD REQUIREMENT: Should return 400 Bad Request for invalid tag type
						entry := barkat.Entry{
							Ticker:   "GRSE",
							Sequence: "MWD",
							Type:     "REJECTED",
							Status:   "FAIL",
							Images: []barkat.Image{
								{Timeframe: "DL"},
								{Timeframe: "WK"},
								{Timeframe: "MN"},
								{Timeframe: "TMN"},
							},
							Tags: []barkat.Tag{
								{Tag: "test", Type: "invalid"},
							},
						}
						req, w = util.CreateTestRequest("POST", JournalEntriesBaseURL, entry)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusCreated))
						var response barkat.Entry
						util.AssertJSONAndStatus(w, http.StatusCreated, &response)
						Expect(response.Tags).To(HaveLen(1))
					})
				})
			})
		})

		Context("Errors", func() {
			It("should return 400 for invalid JSON", func() {
				req, w = util.CreateTestRequest("POST", JournalEntriesBaseURL, []byte("invalid json"))
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})

			It("should return 400 for empty request body", func() {
				req, w = util.CreateTestRequest("POST", JournalEntriesBaseURL, []byte(""))
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})

			It("should return 400 for null request body", func() {
				req, w = util.CreateTestRequest("POST", JournalEntriesBaseURL, []byte("null"))
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})
		})
	})

	Describe("DELETE /v1/journal-entries/{id} - Delete Entry", func() {
		var createdEntry barkat.Entry

		BeforeEach(func() {
			entry := barkat.Entry{
				Ticker:   "GRSE",
				Sequence: "MWD",
				Type:     "REJECTED",
				Status:   "FAIL",
			}
			Expect(entryMgr.CreateEntry(testCtx, &entry)).To(Succeed())
			createdEntry = entry
		})

		Context("Happy Path", func() {
			Context("with valid entry ID", func() {
				BeforeEach(func() {
					req, w = util.CreateTestRequest("DELETE", JournalEntriesBaseURL+"/"+createdEntry.ID, nil)
					router.ServeHTTP(w, req)
				})

				It("should return 204 No Content", func() {
					Expect(w.Code).To(Equal(http.StatusNoContent))
				})

				It("should actually delete the entry", func() {
					_, err := entryMgr.GetEntry(testCtx, createdEntry.ID)
					Expect(err).To(HaveOccurred())
				})
			})
		})

		Context("Field Validations", func() {
			Context("Entry ID Field", func() {
				Context("Bad Values", func() {
					It("should return 404 for non-existent entry ID", func() {
						req, w = util.CreateTestRequest("DELETE", JournalEntriesBaseURL+"/nonexistent-id", nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusNotFound))
					})

					It("should return 404 for malformed UUID", func() {
						req, w = util.CreateTestRequest("DELETE", JournalEntriesBaseURL+"/invalid-uuid-format", nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusNotFound))
					})

					It("should return 404 for valid UUID format but non-existent", func() {
						req, w = util.CreateTestRequest("DELETE", JournalEntriesBaseURL+"/550e8400-e29b-41d4-a716-446655440000", nil)
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
})
