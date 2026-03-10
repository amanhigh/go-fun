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
	"github.com/amanhigh/go-fun/models/common"
	"github.com/gin-gonic/gin"
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
		entryMgr       manager.JournalManager
		req            *http.Request
		w              *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		// FIXME: #A Sqlite Snapshots, Comparison and Revert.
		var err error
		db, err = core.CreateTestBarkatDB()
		Expect(err).ToNot(HaveOccurred())

		entryRepo := repository.NewJournalRepository(db)
		entryMgr = manager.NewJournalManager(entryRepo)
		journalHandler = handler.NewJournalHandler(entryMgr)

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
		Context("Happy Path", func() {
			Context("with minimal valid entry (required fields + min 4 images)", func() {
				var envelopeResponse common.Envelope[barkat.Journal]

				BeforeEach(func() {
					entry := barkat.Journal{
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
					}
					req, w = util.CreateTestRequest("POST", barkat.JournalEntries, entry)
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
					Expect(envelopeResponse.Data.ID).ToNot(BeEmpty())
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

				It("should persist entry to database", func() {
					util.AssertSuccess(w, http.StatusCreated, &envelopeResponse)
					dbEntry, err := entryMgr.GetJournal(testCtx, envelopeResponse.Data.ID)
					Expect(err).ToNot(HaveOccurred())
					Expect(dbEntry.ID).To(Equal(envelopeResponse.Data.ID))
					Expect(dbEntry.Ticker).To(Equal("GRSE"))
				})
			})

			Context("with complete valid entry including images, tags, and notes", func() {
				BeforeEach(func() {
					entry := barkat.Journal{
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
							{Tag: "oe", Type: "REASON"},
						},
						Notes: []barkat.Note{
							{Status: "SET", Content: "Strong OE at weekly level.", Format: "markdown"},
						},
					}
					req, w = util.CreateTestRequest("POST", barkat.JournalEntries, entry)
					router.ServeHTTP(w, req)
				})

				It("should create entry with all associations", func() {
					response := decodeCreateJournalResponse(w)
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
						entry := barkat.Journal{Ticker: "GRSE123", Sequence: "MWD", Type: "REJECTED", Status: "FAIL", Images: []barkat.Image{{Timeframe: "DL"}, {Timeframe: "WK"}, {Timeframe: "MN"}, {Timeframe: "TMN"}}}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, entry)
						router.ServeHTTP(w, req)
						response := decodeCreateJournalResponse(w)
						Expect(response.Ticker).To(Equal("GRSE123"))
					})

					It("should accept ticker with hyphen", func() {
						entry := barkat.Journal{
							Ticker: "GRSE-NSE", Sequence: "MWD", Type: "REJECTED", Status: "FAIL",
							Images: []barkat.Image{{Timeframe: "DL"}, {Timeframe: "WK"}, {Timeframe: "MN"}, {Timeframe: "TMN"}},
						}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, entry)
						router.ServeHTTP(w, req)
						response := decodeCreateJournalResponse(w)
						Expect(response.Ticker).To(Equal("GRSE-NSE"))
					})

					It("should accept ticker with dot suffix", func() {
						entry := barkat.Journal{Ticker: "TCS.NS", Sequence: "YR", Type: "SET", Status: "RUNNING", Images: []barkat.Image{{Timeframe: "DL"}, {Timeframe: "WK"}, {Timeframe: "MN"}, {Timeframe: "TMN"}}}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, entry)
						router.ServeHTTP(w, req)
						response := decodeCreateJournalResponse(w)
						Expect(response.Ticker).To(Equal("TCS.NS"))
					})

					It("should accept ticker at max length (10)", func() {
						entry := barkat.Journal{
							Ticker: "1234567890", Sequence: "MWD", Type: "REJECTED", Status: "FAIL",
							Images: []barkat.Image{{Timeframe: "DL"}, {Timeframe: "WK"}, {Timeframe: "MN"}, {Timeframe: "TMN"}},
						}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, entry)
						router.ServeHTTP(w, req)
						response := decodeCreateJournalResponse(w)
						Expect(response.Ticker).To(HaveLen(10))
					})
				})

				Context("Bad Values", func() {
					It("should return 400 for missing ticker", func() {
						entry := barkat.Journal{Ticker: "", Sequence: "MWD", Type: "REJECTED", Status: "FAIL"}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, entry)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Ticker", "required")
					})

					It("should return 400 for ticker exceeding max length (11)", func() {
						entry := barkat.Journal{Ticker: "12345678901", Sequence: "MWD", Type: "REJECTED", Status: "FAIL"}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, entry)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Ticker", "max (10)")
					})

					It("should return 400 for lowercase ticker (PRD: uppercase only)", func() {
						entry := barkat.Journal{Ticker: "grse", Sequence: "MWD", Type: "REJECTED", Status: "FAIL", Images: []barkat.Image{{Timeframe: "DL"}, {Timeframe: "WK"}, {Timeframe: "MN"}, {Timeframe: "TMN"}}}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, entry)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Ticker", "ticker_format")
					})
				})
			})

			Context("Sequence Field", func() {
				Context("Allowed Values", func() {
					It("should accept sequence = MWD", func() {
						entry := barkat.Journal{Ticker: "PDSL", Sequence: "MWD", Type: "SET", Status: "TAKEN", Images: []barkat.Image{{Timeframe: "DL"}, {Timeframe: "WK"}, {Timeframe: "MN"}, {Timeframe: "TMN"}}}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, entry)
						router.ServeHTTP(w, req)
						response := decodeCreateJournalResponse(w)
						Expect(response.Sequence).To(Equal("MWD"))
					})

					It("should accept sequence = YR", func() {
						entry := barkat.Journal{Ticker: "SNF", Sequence: "YR", Type: "RESULT", Status: "SUCCESS", Images: []barkat.Image{{Timeframe: "DL"}, {Timeframe: "WK"}, {Timeframe: "MN"}, {Timeframe: "TMN"}}}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, entry)
						router.ServeHTTP(w, req)
						response := decodeCreateJournalResponse(w)
						Expect(response.Sequence).To(Equal("YR"))
					})
				})

				Context("Bad Values", func() {
					It("should return 400 for missing sequence", func() {
						entry := barkat.Journal{Ticker: "GRSE", Sequence: "", Type: "REJECTED", Status: "FAIL"}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, entry)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Sequence", "required")
					})

					It("should return 400 for invalid sequence (lowercase)", func() {
						entry := barkat.Journal{Ticker: "INFY", Sequence: "mwd", Type: "SET", Status: "TAKEN"}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, entry)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Sequence", "oneof")
					})

					It("should return 400 for invalid sequence (unsupported)", func() {
						entry := barkat.Journal{Ticker: "WIPRO", Sequence: "QUARTERLY", Type: "SET", Status: "TAKEN"}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, entry)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Sequence", "oneof")
					})
				})
			})

			Context("Type Field", func() {
				Context("Allowed Values", func() {
					It("should accept type = REJECTED", func() {
						entry := barkat.Journal{Ticker: "TCS", Sequence: "MWD", Type: "REJECTED", Status: "FAIL", Images: []barkat.Image{{Timeframe: "DL"}, {Timeframe: "WK"}, {Timeframe: "MN"}, {Timeframe: "TMN"}}}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, entry)
						router.ServeHTTP(w, req)
						response := decodeCreateJournalResponse(w)
						Expect(response.Type).To(Equal("REJECTED"))
					})

					It("should accept type = RESULT", func() {
						entry := barkat.Journal{Ticker: "INFY", Sequence: "YR", Type: "RESULT", Status: "SUCCESS", Images: []barkat.Image{{Timeframe: "DL"}, {Timeframe: "WK"}, {Timeframe: "MN"}, {Timeframe: "TMN"}}}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, entry)
						router.ServeHTTP(w, req)
						response := decodeCreateJournalResponse(w)
						Expect(response.Type).To(Equal("RESULT"))
					})

					It("should accept type = SET", func() {
						entry := barkat.Journal{Ticker: "RELIANCE", Sequence: "MWD", Type: "SET", Status: "RUNNING", Images: []barkat.Image{{Timeframe: "DL"}, {Timeframe: "WK"}, {Timeframe: "MN"}, {Timeframe: "TMN"}}}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, entry)
						router.ServeHTTP(w, req)
						response := decodeCreateJournalResponse(w)
						Expect(response.Type).To(Equal("SET"))
					})
				})

				Context("Bad Values", func() {
					It("should return 400 for missing type", func() {
						entry := barkat.Journal{Ticker: "GRSE", Sequence: "MWD", Type: "", Status: "FAIL"}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, entry)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Type", "required")
					})

					It("should return 400 for invalid type (lowercase)", func() {
						entry := barkat.Journal{Ticker: "GRSE", Sequence: "MWD", Type: "rejected", Status: "FAIL"}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, entry)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Type", "oneof")
					})

					It("should return 400 for invalid type (unsupported)", func() {
						entry := barkat.Journal{Ticker: "HDFC", Sequence: "MWD", Type: "INVALID", Status: "TAKEN"}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, entry)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Type", "oneof")
					})
				})
			})

			Context("Status Field", func() {
				Context("Allowed Values - All 10 Status Values", func() {
					It("should accept status = SET", func() {
						entry := barkat.Journal{Ticker: "T1", Sequence: "MWD", Type: "SET", Status: "SET", Images: []barkat.Image{{Timeframe: "DL"}, {Timeframe: "WK"}, {Timeframe: "MN"}, {Timeframe: "TMN"}}}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, entry)
						router.ServeHTTP(w, req)
						response := decodeCreateJournalResponse(w)
						Expect(response.Status).To(Equal("SET"))
					})

					It("should accept status = RUNNING", func() {
						entry := barkat.Journal{Ticker: "T2", Sequence: "MWD", Type: "SET", Status: "RUNNING", Images: []barkat.Image{{Timeframe: "DL"}, {Timeframe: "WK"}, {Timeframe: "MN"}, {Timeframe: "TMN"}}}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, entry)
						router.ServeHTTP(w, req)
						response := decodeCreateJournalResponse(w)
						Expect(response.Status).To(Equal("RUNNING"))
					})

					It("should accept status = DROPPED", func() {
						entry := barkat.Journal{Ticker: "T3", Sequence: "MWD", Type: "SET", Status: "DROPPED", Images: []barkat.Image{{Timeframe: "DL"}, {Timeframe: "WK"}, {Timeframe: "MN"}, {Timeframe: "TMN"}}}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, entry)
						router.ServeHTTP(w, req)
						response := decodeCreateJournalResponse(w)
						Expect(response.Status).To(Equal("DROPPED"))
					})

					It("should accept status = TAKEN", func() {
						entry := barkat.Journal{Ticker: "T4", Sequence: "MWD", Type: "SET", Status: "TAKEN", Images: []barkat.Image{{Timeframe: "DL"}, {Timeframe: "WK"}, {Timeframe: "MN"}, {Timeframe: "TMN"}}}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, entry)
						router.ServeHTTP(w, req)
						response := decodeCreateJournalResponse(w)
						Expect(response.Status).To(Equal("TAKEN"))
					})

					It("should accept status = REJECTED", func() {
						entry := barkat.Journal{Ticker: "T5", Sequence: "MWD", Type: "REJECTED", Status: "REJECTED", Images: []barkat.Image{{Timeframe: "DL"}, {Timeframe: "WK"}, {Timeframe: "MN"}, {Timeframe: "TMN"}}}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, entry)
						router.ServeHTTP(w, req)
						response := decodeCreateJournalResponse(w)
						Expect(response.Status).To(Equal("REJECTED"))
					})

					It("should accept status = SUCCESS", func() {
						entry := barkat.Journal{Ticker: "T6", Sequence: "YR", Type: "RESULT", Status: "SUCCESS", Images: []barkat.Image{{Timeframe: "DL"}, {Timeframe: "WK"}, {Timeframe: "MN"}, {Timeframe: "TMN"}}}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, entry)
						router.ServeHTTP(w, req)
						response := decodeCreateJournalResponse(w)
						Expect(response.Status).To(Equal("SUCCESS"))
					})

					It("should accept status = FAIL", func() {
						entry := barkat.Journal{Ticker: "T7", Sequence: "YR", Type: "RESULT", Status: "FAIL", Images: []barkat.Image{{Timeframe: "DL"}, {Timeframe: "WK"}, {Timeframe: "MN"}, {Timeframe: "TMN"}}}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, entry)
						router.ServeHTTP(w, req)
						response := decodeCreateJournalResponse(w)
						Expect(response.Status).To(Equal("FAIL"))
					})

					It("should accept status = MISSED", func() {
						entry := barkat.Journal{Ticker: "T8", Sequence: "YR", Type: "RESULT", Status: "MISSED", Images: []barkat.Image{{Timeframe: "DL"}, {Timeframe: "WK"}, {Timeframe: "MN"}, {Timeframe: "TMN"}}}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, entry)
						router.ServeHTTP(w, req)
						response := decodeCreateJournalResponse(w)
						Expect(response.Status).To(Equal("MISSED"))
					})

					It("should accept status = JUST_LOSS", func() {
						entry := barkat.Journal{Ticker: "T9", Sequence: "YR", Type: "RESULT", Status: "JUST_LOSS", Images: []barkat.Image{{Timeframe: "DL"}, {Timeframe: "WK"}, {Timeframe: "MN"}, {Timeframe: "TMN"}}}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, entry)
						router.ServeHTTP(w, req)
						response := decodeCreateJournalResponse(w)
						Expect(response.Status).To(Equal("JUST_LOSS"))
					})

					It("should accept status = BROKEN", func() {
						entry := barkat.Journal{Ticker: "T10", Sequence: "YR", Type: "RESULT", Status: "BROKEN", Images: []barkat.Image{{Timeframe: "DL"}, {Timeframe: "WK"}, {Timeframe: "MN"}, {Timeframe: "TMN"}}}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, entry)
						router.ServeHTTP(w, req)
						response := decodeCreateJournalResponse(w)
						Expect(response.Status).To(Equal("BROKEN"))
					})
				})

				Context("Bad Values", func() {
					It("should return 400 for missing status", func() {
						entry := barkat.Journal{Ticker: "GRSE", Sequence: "MWD", Type: "REJECTED", Status: ""}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, entry)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Status", "required")
					})

					It("should return 400 for invalid status (lowercase)", func() {
						entry := barkat.Journal{Ticker: "GRSE", Sequence: "MWD", Type: "REJECTED", Status: "fail"}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, entry)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Status", "oneof")
					})

					It("should return 400 for invalid status (unsupported)", func() {
						entry := barkat.Journal{Ticker: "HDFC", Sequence: "MWD", Type: "SET", Status: "INVALID"}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, entry)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Status", "oneof")
					})
				})
			})

			Context("Images Field", func() {
				Context("Bad Values", func() {
					It("should return 400 for empty images (PRD: min 4 required)", func() {
						entry := barkat.Journal{
							Ticker:   "GRSE",
							Sequence: "MWD",
							Type:     "REJECTED",
							Status:   "FAIL",
							Images:   []barkat.Image{},
						}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, entry)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Images", "required")
					})

					It("should return 400 for insufficient images < 4 (PRD: min 4 required)", func() {
						entry := barkat.Journal{
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
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, entry)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Images", "min")
					})

					It("should return 400 for excessive images > 6 (PRD: max 6 allowed)", func() {
						entry := barkat.Journal{
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
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, entry)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Images", "max")
					})

					It("should return 400 for invalid timeframe (PRD: must be DL,WK,MN,TMN,SMN,YR)", func() {
						entry := barkat.Journal{
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
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, entry)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Timeframe", "oneof")
					})

					It("should return 400 for duplicate timeframes (PRD: unique timeframes required)", func() {
						entry := barkat.Journal{
							Ticker:   "GRSE",
							Sequence: "MWD",
							Type:     "REJECTED",
							Status:   "FAIL",
							Images: []barkat.Image{
								{Timeframe: "DL"},
								{Timeframe: "DL"},
								{Timeframe: "MN"},
								{Timeframe: "TMN"},
							},
						}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, entry)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Images", "Duplicate")
					})
				})
			})

			Context("Notes Field", func() {
				Context("Bad Values", func() {
					It("should return 400 for missing note status (PRD: status required)", func() {
						entry := barkat.Journal{
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
								{Status: "", Content: "Note without status", Format: "MARKDOWN"},
							},
						}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, entry)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Status", "required")
					})

					It("should return 400 for missing note content (PRD: content required)", func() {
						entry := barkat.Journal{
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
								{Status: "SET", Content: "", Format: "MARKDOWN"},
							},
						}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, entry)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Content", "required")
					})

					It("should return 400 for multiple notes > 1 (PRD: max 1 at create)", func() {
						entry := barkat.Journal{
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
								{Status: "SET", Content: "First note", Format: "MARKDOWN"},
								{Status: "RUNNING", Content: "Second note", Format: "PLAINTEXT"},
							},
						}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, entry)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Notes", "max")
					})

					It("should return 400 for invalid note format (PRD: must be MARKDOWN or PLAINTEXT)", func() {
						// HACK: Extract Created in Before each and override field for Bad request.
						entry := barkat.Journal{
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
								{Status: "SET", Content: "Note with invalid format", Format: "invalid"},
							},
						}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, entry)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Format", "oneof")
					})

					It("should return 400 for invalid note status (PRD: must match entry status enum)", func() {
						entry := barkat.Journal{
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
								{Status: "INVALID", Content: "Note with invalid status", Format: "MARKDOWN"},
							},
						}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, entry)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Status", "oneof")
					})

					It("should return 400 for note content exceeding max length (PRD: max 2000 chars)", func() {
						longContent := string(make([]byte, 2001))
						for i := range longContent {
							longContent = longContent[:i] + "a" + longContent[i+1:]
						}
						entry := barkat.Journal{
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
								{Status: "SET", Content: longContent, Format: "MARKDOWN"},
							},
						}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, entry)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Content", "max")
					})
				})
			})

			Context("Tags Field", func() {
				Context("Allowed Values", func() {
					It("should accept tag with override field (PRD: optional override max 50 chars)", func() {
						override := "loc"
						entry := barkat.Journal{
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
								{Tag: "dep", Type: "REASON", Override: &override},
							},
						}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, entry)
						router.ServeHTTP(w, req)
						var envelopeResponse common.Envelope[barkat.Journal]
						util.AssertSuccess(w, http.StatusCreated, &envelopeResponse)
						Expect(envelopeResponse.Status).To(Equal(common.EnvelopeSuccess))
						Expect(envelopeResponse.Data.Tags).To(HaveLen(1))
						Expect(envelopeResponse.Data.Tags[0].Override).ToNot(BeNil())
						Expect(*envelopeResponse.Data.Tags[0].Override).To(Equal("loc"))
					})

					It("should accept tag type = REASON (PRD: uppercase)", func() {
						entry := barkat.Journal{
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
								{Tag: "oe", Type: "REASON"},
							},
						}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, entry)
						router.ServeHTTP(w, req)
						var envelopeResponse common.Envelope[barkat.Journal]
						util.AssertSuccess(w, http.StatusCreated, &envelopeResponse)
					})

					It("should accept tag type = MANAGEMENT (PRD: uppercase)", func() {
						entry := barkat.Journal{
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
								{Tag: "sl", Type: "MANAGEMENT"},
							},
						}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, entry)
						router.ServeHTTP(w, req)
						var envelopeResponse common.Envelope[barkat.Journal]
						util.AssertSuccess(w, http.StatusCreated, &envelopeResponse)
					})
				})

				Context("Bad Values", func() {
					It("should return 400 for missing tag name (PRD: tag required)", func() {
						entry := barkat.Journal{
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
								{Tag: "", Type: "REASON"},
							},
						}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, entry)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Tag", "required")
					})

					It("should return 400 for missing tag type (PRD: type required)", func() {
						entry := barkat.Journal{
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
								{Tag: "oe", Type: ""},
							},
						}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, entry)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Type", "required")
					})

					It("should return 400 for invalid tag type (PRD: must be REASON or MANAGEMENT)", func() {
						entry := barkat.Journal{
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
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, entry)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Type", "oneof")
					})

					It("should return 400 for tag exceeding max length (PRD: max 10 chars)", func() {
						entry := barkat.Journal{
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
								{Tag: "verylongtag1", Type: "REASON"},
							},
						}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, entry)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Tag", "max")
					})

					It("should return 400 for override exceeding max length (PRD: max 5 chars)", func() {
						longOverride := "toolong"
						entry := barkat.Journal{
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
								{Tag: "dep", Type: "REASON", Override: &longOverride},
							},
						}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, entry)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Override", "max")
					})

					It("should return 400 for invalid tag format (PRD: alphanumeric with hyphens)", func() {
						entry := barkat.Journal{
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
								{Tag: "bad@tag", Type: "REASON"},
							},
						}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, entry)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Tag", "tag_format")
					})

					It("should return 400 for invalid override format (PRD: alphanumeric only)", func() {
						invalidOverride := "a-b"
						entry := barkat.Journal{
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
								{Tag: "dep", Type: "REASON", Override: &invalidOverride},
							},
						}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, entry)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Override", "override_format")
					})

					It("should return 400 for exceeding max tags (PRD: max 10)", func() {
						entry := barkat.Journal{
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
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries, entry)
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
		var createdEntry barkat.Journal

		BeforeEach(func() {
			entry := barkat.Journal{
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
			}
			Expect(entryMgr.CreateJournal(testCtx, &entry)).To(Succeed())
			createdEntry = entry
		})

		Context("Happy Path", func() {
			Context("with valid entry ID", func() {
				BeforeEach(func() {
					req, w = util.CreateTestRequest("DELETE", barkat.JournalEntries+"/"+createdEntry.ID, nil)
					router.ServeHTTP(w, req)
				})

				It("should return 204 No Content", func() {
					Expect(w.Code).To(Equal(http.StatusNoContent))
				})

				It("should actually delete the entry", func() {
					_, err := entryMgr.GetJournal(testCtx, createdEntry.ID)
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
})
