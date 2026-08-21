//nolint:dupl
package handler_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/components/kohan/core"
	"github.com/amanhigh/go-fun/components/kohan/handler"
	"github.com/amanhigh/go-fun/components/kohan/manager"
	"github.com/amanhigh/go-fun/components/kohan/repository"
	"github.com/amanhigh/go-fun/models/barkat"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"
)

func decodeJournal(w *httptest.ResponseRecorder, expectedStatus int) barkat.Journal {
	var envelope common.Envelope[barkat.Journal]
	util.AssertSuccess(w, expectedStatus, &envelope)
	return envelope.Data
}

func decodeJournalList(w *httptest.ResponseRecorder, expectedStatus int) barkat.JournalList {
	var envelope common.Envelope[barkat.JournalList]
	util.AssertSuccess(w, expectedStatus, &envelope)
	return envelope.Data
}

var _ = Describe("JournalHandler Integration - GET Tests", func() {
	var (
		journalHandler *handler.JournalHandlerImpl
		router         *gin.Engine
		testCtx        = context.Background()
		db             *gorm.DB
		journalMgr     manager.JournalManager
		req            *http.Request
		w              *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		var err error
		db, err = core.CreateTestBarkatDB()
		Expect(err).ToNot(HaveOccurred())

		journalRepo := repository.NewJournalRepository(util.NewBaseDbRepository(db))
		journalMgr = manager.NewJournalManager(journalRepo)
		journalHandler = handler.NewJournalHandler(journalMgr)

		router = util.CreateTestGinRouter()
		v1 := router.Group("/v1/api")
		journal := v1.Group("/journals")
		handler.SetupJournalRoutes(journal, journalHandler)
	})

	AfterEach(func() {
		sqlDB, err := db.DB()
		Expect(err).ToNot(HaveOccurred())
		sqlDB.Close()
	})

	Describe("GET /v1/journal/{id} - Retrieve Journal", func() {
		var createdJournal barkat.Journal

		BeforeEach(func() {
			journal := barkat.Journal{
				Ticker:       "GRSE",
				TopTimeframe: "TMN",
				Type:         "REJECTED",
				Status:       "FAIL",
				Images: []barkat.Image{
					{Timeframe: "WK", CreatedAt: time.Date(2023, time.June, 1, 10, 0, 0, 0, time.UTC), ImageType: "SET"},
					{Timeframe: "DL", CreatedAt: time.Date(2023, time.June, 2, 10, 0, 0, 0, time.UTC), ImageType: "SET"},
					{Timeframe: "MN", CreatedAt: time.Date(2023, time.June, 2, 10, 0, 0, 0, time.UTC), ImageType: "SET"},
					{Timeframe: "TMN", CreatedAt: time.Date(2023, time.June, 2, 10, 0, 0, 0, time.UTC), ImageType: "SET"},
					{Timeframe: "SMN", CreatedAt: time.Date(2023, time.June, 2, 10, 0, 0, 0, time.UTC), ImageType: "SET"},
					{Timeframe: "YR", CreatedAt: time.Date(2023, time.June, 2, 10, 0, 0, 0, time.UTC), ImageType: "SET"},
				},
			}
			Expect(journalMgr.CreateJournal(testCtx, &journal)).To(Succeed())
			createdJournal = journal
		})

		Context("Happy Path", func() {
			Context("with valid journal ID", func() {
				var response barkat.Journal

				BeforeEach(func() {
					req, w = util.CreateTestRequest("GET", barkat.JournalBase+"/"+createdJournal.ExternalID, nil)
					router.ServeHTTP(w, req)
				})

				It("should return 200 OK", func() {
					Expect(w.Code).To(Equal(http.StatusOK))
				})

				It("should return journal with correct ID", func() {
					response = decodeJournal(w, http.StatusOK)
					Expect(response.ExternalID).To(Equal(createdJournal.ExternalID))
				})

				It("should return all journal fields including images", func() {
					response = decodeJournal(w, http.StatusOK)
					Expect(response.Ticker).To(Equal("GRSE"))
					Expect(response.TopTimeframe).To(Equal("TMN"))
					Expect(response.Type).To(Equal("REJECTED"))
					Expect(response.Status).To(Equal("FAIL"))
					Expect(response.CreatedAt).ToNot(BeZero())
					Expect(response.Images).To(HaveLen(6))
				})

				It("should return images sorted by date then timeframe", func() {
					response = decodeJournal(w, http.StatusOK)

					Expect(response.Images[0].CreatedAt).To(Equal(time.Date(2023, time.June, 1, 10, 0, 0, 0, time.UTC)))
					Expect(response.Images[0].Timeframe).To(Equal("WK"))

					Expect(response.Images[1].Timeframe).To(Equal("YR"))
					Expect(response.Images[2].Timeframe).To(Equal("SMN"))
					Expect(response.Images[3].Timeframe).To(Equal("TMN"))
					Expect(response.Images[4].Timeframe).To(Equal("MN"))
					Expect(response.Images[5].Timeframe).To(Equal("DL"))
				})
			})
		})

		Context("Field Validations", func() {
			Context("Journal ID Field", func() {
				Context("Bad Values", func() {
					It("should return 400 for malformed UUID", func() {
						req, w = util.CreateTestRequest("GET", barkat.JournalBase+"/invalid-uuid-format", nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusBadRequest))
					})

					It("should return 404 for valid journal ID format but non-existent", func() {
						req, w = util.CreateTestRequest("GET", barkat.JournalBase+"/jrn_12345678", nil)
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

	Describe("GET /v1/journal - List Entries", func() {
		var createdJournals []barkat.Journal

		BeforeEach(func() {
			// Define default images template
			defaultImages := []barkat.Image{
				{Timeframe: "DL", ImageType: "SET"},
				{Timeframe: "WK", ImageType: "SET"},
				{Timeframe: "MN", ImageType: "SET"},
				{Timeframe: "TMN", ImageType: "SET"},
				{Timeframe: "SMN", ImageType: "SET"},
				{Timeframe: "YR", ImageType: "SET"},
			}

			journals := []barkat.Journal{
				{Ticker: "GRSE", TopTimeframe: "TMN", Type: "REJECTED", Status: "FAIL"},
				{Ticker: "PDSL", TopTimeframe: "SMN", Type: "TAKEN", Status: "SET"},
				{Ticker: "SNF", TopTimeframe: "TMN", Type: "TAKEN", Status: "SUCCESS"},
				{Ticker: "TCS", TopTimeframe: "SMN", Type: "REJECTED", Status: "BROKEN"},
				{Ticker: "INFY", TopTimeframe: "TMN", Type: "TAKEN", Status: "RUNNING"},
			}

			// Copy default images for each journal to avoid shared slice mutation
			for i := range journals {
				var copiedImages []barkat.Image
				err := copier.Copy(&copiedImages, &defaultImages)
				Expect(err).ToNot(HaveOccurred())
				journals[i].Images = copiedImages
			}

			for _, journal := range journals {
				Expect(journalMgr.CreateJournal(testCtx, &journal)).To(Succeed())
				createdJournals = append(createdJournals, journal)
			}
		})

		Context("Happy Path", func() {
			Context("default pagination (no filters)", func() {
				var response barkat.JournalList

				BeforeEach(func() {
					req, w = util.CreateTestRequest("GET", barkat.JournalBase, nil)
					router.ServeHTTP(w, req)
				})

				It("should return 200 OK", func() {
					Expect(w.Code).To(Equal(http.StatusOK))
				})

				It("should return all entries", func() {
					response = decodeJournalList(w, http.StatusOK)
					Expect(response.Journals).To(HaveLen(5))
				})

				It("should return correct total count", func() {
					response = decodeJournalList(w, http.StatusOK)
					Expect(response.Metadata.Total).To(Equal(int64(5)))
					Expect(response.Metadata.Offset).To(Equal(0))
					Expect(response.Metadata.Limit).To(Equal(20))
				})

				It("should return entries in reverse chronological order by default", func() {
					response = decodeJournalList(w, http.StatusOK)
					journals := response.Journals
					for i := 1; i < len(journals); i++ {
						prevTime := journals[i-1].CreatedAt
						currTime := journals[i].CreatedAt
						Expect(prevTime).To(BeTemporally(">=", currTime))
					}
				})

				It("should include all required fields and images in each journal", func() {
					response = decodeJournalList(w, http.StatusOK)
					for _, journal := range response.Journals {
						Expect(journal.ExternalID).To(HavePrefix("jrn_"))
						Expect(journal.Ticker).ToNot(BeEmpty())
						Expect(journal.TopTimeframe).ToNot(BeEmpty())
						Expect(journal.Type).ToNot(BeEmpty())
						Expect(journal.Status).ToNot(BeEmpty())
						Expect(journal.CreatedAt).ToNot(BeZero())
						Expect(journal.Images).To(HaveLen(6))
						Expect(journal.Images[0].Timeframe).To(Equal("DL"))
						Expect(journal.Images[1].Timeframe).To(Equal("WK"))
						Expect(journal.Images[2].Timeframe).To(Equal("MN"))
						Expect(journal.Images[3].Timeframe).To(Equal("TMN"))
						Expect(journal.Images[4].Timeframe).To(Equal("SMN"))
						Expect(journal.Images[5].Timeframe).To(Equal("YR"))
					}
				})
			})
		})

		Context("Field Validations", func() {
			Context("Ticker Filter", func() {
				Context("Allowed Values", func() {
					It("should filter by exact ticker match", func() {
						req, w = util.CreateTestRequest("GET", barkat.JournalBase+"?ticker=GRSE", nil)
						router.ServeHTTP(w, req)
						response := decodeJournalList(w, http.StatusOK)
						Expect(response.Journals).To(HaveLen(1))
						journals := response.Journals
						Expect(journals[0].Ticker).To(Equal("GRSE"))
						Expect(response.Metadata.Total).To(Equal(int64(1)))
						Expect(response.Metadata.Offset).To(Equal(0))
						Expect(response.Metadata.Limit).To(Equal(20))
					})

					It("should return empty list for ticker with no matches", func() {
						req, w = util.CreateTestRequest("GET", barkat.JournalBase+"?ticker=NOTFOUND", nil)
						router.ServeHTTP(w, req)
						response := decodeJournalList(w, http.StatusOK)
						Expect(response.Journals).To(BeEmpty())
						Expect(response.Metadata.Total).To(Equal(int64(0)))
						Expect(response.Metadata.Offset).To(Equal(0))
						Expect(response.Metadata.Limit).To(Equal(20))
					})
				})

				Context("Bad Values", func() {
					It("should return 400 for invalid ticker length", func() {
						req, w = util.CreateTestRequest("GET", barkat.JournalBase+"?ticker=1234567890123456789012345678901", nil)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Ticker", "max")
					})

					It("should return 400 for lowercase ticker (PRD: uppercase only)", func() {
						req, w = util.CreateTestRequest("GET", barkat.JournalBase+"?ticker=grse", nil)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Ticker", "ticker")
					})

					It("should return 400 for ticker starting with dot", func() {
						req, w = util.CreateTestRequest("GET", barkat.JournalBase+"?ticker=.MCX", nil)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Ticker", "ticker")
					})

					It("should return 400 for ticker with hyphen", func() {
						req, w = util.CreateTestRequest("GET", barkat.JournalBase+"?ticker=ABC-DEF", nil)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Ticker", "ticker")
					})

					It("should return 400 for ticker with unsupported special character", func() {
						req, w = util.CreateTestRequest("GET", barkat.JournalBase+"?ticker=MCX@", nil)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Ticker", "ticker")
					})

					It("should return 400 for ticker with whitespace", func() {
						req, w = util.CreateTestRequest("GET", barkat.JournalBase+"?ticker=MC%20X", nil)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Ticker", "ticker")
					})
				})
			})

			Context("Search Filter", func() {
				Context("Allowed Values", func() {
					It("should filter by case-insensitive ticker substring", func() {
						req, w = util.CreateTestRequest("GET", barkat.JournalBase+"?search=rs", nil)
						router.ServeHTTP(w, req)
						response := decodeJournalList(w, http.StatusOK)
						Expect(response.Journals).To(HaveLen(1))
						Expect(response.Journals[0].Ticker).To(Equal("GRSE"))
					})

					It("should combine search with exact ticker filter", func() {
						req, w = util.CreateTestRequest("GET", barkat.JournalBase+"?search=rs&ticker=GRSE", nil)
						router.ServeHTTP(w, req)
						response := decodeJournalList(w, http.StatusOK)
						Expect(response.Journals).To(HaveLen(1))
						Expect(response.Journals[0].Ticker).To(Equal("GRSE"))
					})
				})

				Context("Bad Values", func() {
					It("should return 400 for invalid search format", func() {
						req, w = util.CreateTestRequest("GET", barkat.JournalBase+"?search=RE*", nil)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Search", "alphanum")
					})

					It("should return 400 for invalid search length", func() {
						req, w = util.CreateTestRequest("GET", barkat.JournalBase+"?search=ABCDEFGHIJK", nil)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Search", "max")
					})
				})
			})

			Context("Type Filter", func() {
				Context("Allowed Values", func() {
					It("should filter by type = REJECTED", func() {
						req, w = util.CreateTestRequest("GET", barkat.JournalBase+"?type=REJECTED", nil)
						router.ServeHTTP(w, req)
						response := decodeJournalList(w, http.StatusOK)
						Expect(response.Journals).To(HaveLen(2))
						for _, journal := range response.Journals {
							Expect(journal.Type).To(Equal("REJECTED"))
						}
					})

					It("should filter by type = TAKEN", func() {
						req, w = util.CreateTestRequest("GET", barkat.JournalBase+"?type=TAKEN", nil)
						router.ServeHTTP(w, req)
						response := decodeJournalList(w, http.StatusOK)
						Expect(response.Journals).To(HaveLen(3))
						Expect(response.Journals[0].Type).To(Equal("TAKEN"))
						Expect(response.Journals[1].Type).To(Equal("TAKEN"))
						Expect(response.Journals[2].Type).To(Equal("TAKEN"))
					})
				})

				Context("Bad Values", func() {
					It("should return 400 for invalid type enum", func() {
						req, w = util.CreateTestRequest("GET", barkat.JournalBase+"?type=invalid", nil)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Type", "oneof")
					})
				})
			})

			Context("Status Filter", func() {
				Context("Allowed Values", func() {
					It("should filter by status = FAIL", func() {
						req, w = util.CreateTestRequest("GET", barkat.JournalBase+"?status=FAIL", nil)
						router.ServeHTTP(w, req)
						response := decodeJournalList(w, http.StatusOK)
						Expect(response.Journals).To(HaveLen(1))
						journals := response.Journals
						Expect(journals[0].Status).To(Equal("FAIL"))
					})

					It("should filter by status = SET", func() {
						req, w = util.CreateTestRequest("GET", barkat.JournalBase+"?status=SET", nil)
						router.ServeHTTP(w, req)
						response := decodeJournalList(w, http.StatusOK)
						Expect(response.Journals).To(HaveLen(1))
						Expect(response.Journals[0].Status).To(Equal("SET"))
					})

					It("should filter by status = SUCCESS", func() {
						req, w = util.CreateTestRequest("GET", barkat.JournalBase+"?status=SUCCESS", nil)
						router.ServeHTTP(w, req)
						response := decodeJournalList(w, http.StatusOK)
						Expect(response.Journals).To(HaveLen(1))
						Expect(response.Journals[0].Status).To(Equal("SUCCESS"))
					})

					It("should filter by status = RUNNING", func() {
						req, w = util.CreateTestRequest("GET", barkat.JournalBase+"?status=RUNNING", nil)
						router.ServeHTTP(w, req)
						response := decodeJournalList(w, http.StatusOK)
						Expect(response.Journals).To(HaveLen(1))
						Expect(response.Journals[0].Status).To(Equal("RUNNING"))
					})

					It("should filter by status = BROKEN", func() {
						req, w = util.CreateTestRequest("GET", barkat.JournalBase+"?status=BROKEN", nil)
						router.ServeHTTP(w, req)
						response := decodeJournalList(w, http.StatusOK)
						Expect(response.Journals).To(HaveLen(1))
						Expect(response.Journals[0].Status).To(Equal("BROKEN"))
					})
				})

				Context("Bad Values", func() {
					It("should return 400 for invalid status enum", func() {
						req, w = util.CreateTestRequest("GET", barkat.JournalBase+"?status=invalid", nil)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Status", "oneof")
					})
				})
			})

			Context("TopTimeframe Filter", func() {
				Context("Allowed Values", func() {
					It("should filter by top_timeframe = TMN", func() {
						req, w = util.CreateTestRequest("GET", barkat.JournalBase+"?top_timeframe=TMN", nil)
						router.ServeHTTP(w, req)
						response := decodeJournalList(w, http.StatusOK)
						Expect(response.Journals).To(HaveLen(3))
						for _, journal := range response.Journals {
							Expect(journal.TopTimeframe).To(Equal("TMN"))
						}
					})

					It("should filter by top_timeframe = SMN", func() {
						req, w = util.CreateTestRequest("GET", barkat.JournalBase+"?top_timeframe=SMN", nil)
						router.ServeHTTP(w, req)
						response := decodeJournalList(w, http.StatusOK)
						Expect(response.Journals).To(HaveLen(2))
						for _, journal := range response.Journals {
							Expect(journal.TopTimeframe).To(Equal("SMN"))
						}
					})
				})

				Context("Bad Values", func() {
					It("should return 400 for invalid top_timeframe enum", func() {
						req, w = util.CreateTestRequest("GET", barkat.JournalBase+"?top_timeframe=invalid", nil)
						router.ServeHTTP(w, req)
						util.AssertError(w, "TopTimeframe", "oneof")
					})

					It("should return 400 for top_timeframe = MN (migration-only value)", func() {
						req, w = util.CreateTestRequest("GET", barkat.JournalBase+"?top_timeframe=MN", nil)
						router.ServeHTTP(w, req)
						util.AssertError(w, "TopTimeframe", "oneof")
					})

					It("should return 400 for lowercase top_timeframe", func() {
						req, w = util.CreateTestRequest("GET", barkat.JournalBase+"?top_timeframe=tmn", nil)
						router.ServeHTTP(w, req)
						util.AssertError(w, "TopTimeframe", "oneof")
					})
				})
			})

			Context("Combined Filters", func() {
				It("should apply ticker + type filters", func() {
					req, w = util.CreateTestRequest("GET", barkat.JournalBase+"?ticker=GRSE&type=REJECTED", nil)
					router.ServeHTTP(w, req)
					response := decodeJournalList(w, http.StatusOK)
					Expect(response.Journals).To(HaveLen(1))
					journals := response.Journals
					Expect(journals[0].Ticker).To(Equal("GRSE"))
					Expect(journals[0].Type).To(Equal("REJECTED"))
				})

				It("should apply top_timeframe + status filters", func() {
					req, w = util.CreateTestRequest("GET", barkat.JournalBase+"?top_timeframe=SMN&status=SET", nil)
					router.ServeHTTP(w, req)
					response := decodeJournalList(w, http.StatusOK)
					Expect(response.Journals).To(HaveLen(1))
					Expect(response.Journals[0].TopTimeframe).To(Equal("SMN"))
					Expect(response.Journals[0].Status).To(Equal("SET"))
				})

				It("should apply type + status + top_timeframe filters", func() {
					req, w = util.CreateTestRequest("GET", barkat.JournalBase+"?type=TAKEN&status=RUNNING&top_timeframe=TMN", nil)
					router.ServeHTTP(w, req)
					response := decodeJournalList(w, http.StatusOK)
					Expect(response.Journals).To(HaveLen(1))
					Expect(response.Journals[0].Type).To(Equal("TAKEN"))
					Expect(response.Journals[0].Status).To(Equal("RUNNING"))
					Expect(response.Journals[0].TopTimeframe).To(Equal("TMN"))
				})
			})

			Context("Date Fields", func() {
				Context("Created-After Field", func() {
					Context("Allowed Values", func() {
						It("should accept valid YYYY-MM-DD date and filter entries", func() {
							// Use yesterday's date to ensure all entries are captured
							afterDate := time.Now().Add(-24 * time.Hour).Format("2006-01-02")
							req, w = util.CreateTestRequest("GET", barkat.JournalBase+"?created-after="+afterDate, nil)
							router.ServeHTTP(w, req)
							response := decodeJournalList(w, http.StatusOK)
							// All entries created in this test should be returned
							Expect(response.Journals).To(HaveLen(5))
						})

						It("should return empty list for future date", func() {
							futureDate := time.Now().Add(24 * time.Hour).Format("2006-01-02")
							req, w = util.CreateTestRequest("GET", barkat.JournalBase+"?created-after="+futureDate, nil)
							router.ServeHTTP(w, req)
							response := decodeJournalList(w, http.StatusOK)
							Expect(response.Journals).To(BeEmpty())
						})

						It("should work with created-before combined filter", func() {
							afterDate := time.Now().Add(-24 * time.Hour).Format("2006-01-02")
							beforeDate := time.Now().Add(24 * time.Hour).Format("2006-01-02")
							req, w = util.CreateTestRequest("GET", barkat.JournalBase+"?created-after="+afterDate+"&created-before="+beforeDate, nil)
							router.ServeHTTP(w, req)
							response := decodeJournalList(w, http.StatusOK)
							Expect(response.Journals).To(HaveLen(5))
						})
					})

					Context("Bad Values", func() {
						It("should return 400 for invalid date format", func() {
							req, w = util.CreateTestRequest("GET", barkat.JournalBase+"?created-after=invalid-date", nil)
							router.ServeHTTP(w, req)
							util.AssertError(w, "CreatedAfter", "Violates 'datetime")
						})

						It("should ignore empty date and return all entries", func() {
							req, w = util.CreateTestRequest("GET", barkat.JournalBase+"?created-after=", nil)
							router.ServeHTTP(w, req)
							response := decodeJournalList(w, http.StatusOK)
							// Empty date is ignored, all entries returned
							Expect(response.Journals).To(HaveLen(5))
						})

						It("should return 400 for wrong format (DD-MM-YYYY)", func() {
							req, w = util.CreateTestRequest("GET", barkat.JournalBase+"?created-after=15-02-2024", nil)
							router.ServeHTTP(w, req)
							util.AssertError(w, "CreatedAfter", "Violates 'datetime")
						})
					})
				})

				Context("Created-Before Field", func() {
					Context("Allowed Values", func() {
						It("should accept valid YYYY-MM-DD date and filter entries", func() {
							beforeDate := time.Now().Add(24 * time.Hour).Format("2006-01-02")
							req, w = util.CreateTestRequest("GET", barkat.JournalBase+"?created-before="+beforeDate, nil)
							router.ServeHTTP(w, req)
							response := decodeJournalList(w, http.StatusOK)
							// All entries created in this test should be returned
							Expect(response.Journals).To(HaveLen(5))
						})

						It("should return empty list for past date", func() {
							pastDate := time.Now().Add(-48 * time.Hour).Format("2006-01-02")
							req, w = util.CreateTestRequest("GET", barkat.JournalBase+"?created-before="+pastDate, nil)
							router.ServeHTTP(w, req)
							response := decodeJournalList(w, http.StatusOK)
							Expect(response.Journals).To(BeEmpty())
						})

						It("should work with created-after combined filter", func() {
							afterDate := time.Now().Add(-24 * time.Hour).Format("2006-01-02")
							beforeDate := time.Now().Add(24 * time.Hour).Format("2006-01-02")
							req, w = util.CreateTestRequest("GET", barkat.JournalBase+"?created-after="+afterDate+"&created-before="+beforeDate, nil)
							router.ServeHTTP(w, req)
							response := decodeJournalList(w, http.StatusOK)
							Expect(response.Journals).To(HaveLen(5))
						})
					})

					Context("Bad Values", func() {
						It("should return 400 for invalid date format", func() {
							req, w = util.CreateTestRequest("GET", barkat.JournalBase+"?created-before=not-a-date", nil)
							router.ServeHTTP(w, req)
							util.AssertError(w, "CreatedBefore", "Violates 'datetime")
						})

						It("should ignore empty date and return all entries", func() {
							req, w = util.CreateTestRequest("GET", barkat.JournalBase+"?created-before=", nil)
							router.ServeHTTP(w, req)
							response := decodeJournalList(w, http.StatusOK)
							// Empty date is ignored, all entries returned
							Expect(response.Journals).To(HaveLen(5))
						})

						It("should return 400 for wrong format (DD-MM-YYYY)", func() {
							req, w = util.CreateTestRequest("GET", barkat.JournalBase+"?created-before=15-02-2024", nil)
							router.ServeHTTP(w, req)
							util.AssertError(w, "CreatedBefore", "Violates 'datetime")
						})
					})
				})
			})

			Context("Sorting", func() {
				Context("Allowed Values", func() {
					It("should sort by ticker ascending", func() {
						req, w = util.CreateTestRequest("GET", barkat.JournalBase+"?sort-by=ticker&sort-order=asc", nil)
						router.ServeHTTP(w, req)
						response := decodeJournalList(w, http.StatusOK)
						Expect(response.Journals).To(HaveLen(5))
						journals := response.Journals
						Expect(journals[0].Ticker).To(Equal("GRSE"))
						Expect(journals[1].Ticker).To(Equal("INFY"))
						Expect(journals[2].Ticker).To(Equal("PDSL"))
						Expect(journals[3].Ticker).To(Equal("SNF"))
						Expect(journals[4].Ticker).To(Equal("TCS"))
					})

					It("should sort by ticker descending", func() {
						req, w = util.CreateTestRequest("GET", barkat.JournalBase+"?sort-by=ticker&sort-order=desc", nil)
						router.ServeHTTP(w, req)
						response := decodeJournalList(w, http.StatusOK)
						Expect(response.Journals).To(HaveLen(5))
						journals := response.Journals
						Expect(journals[0].Ticker).To(Equal("TCS"))
						Expect(journals[1].Ticker).To(Equal("SNF"))
						Expect(journals[2].Ticker).To(Equal("PDSL"))
						Expect(journals[3].Ticker).To(Equal("INFY"))
						Expect(journals[4].Ticker).To(Equal("GRSE"))
					})

					It("should sort by top_timeframe ascending", func() {
						req, w = util.CreateTestRequest("GET", barkat.JournalBase+"?sort-by=top_timeframe&sort-order=asc", nil)
						router.ServeHTTP(w, req)
						response := decodeJournalList(w, http.StatusOK)
						Expect(response.Journals).To(HaveLen(5))
						journals := response.Journals
						for i := range 2 {
							Expect(journals[i].TopTimeframe).To(Equal("SMN"))
						}
						for i := 2; i < 5; i++ {
							Expect(journals[i].TopTimeframe).To(Equal("TMN"))
						}
					})

					It("should sort by created_at ascending", func() {
						req, w = util.CreateTestRequest("GET", barkat.JournalBase+"?sort-by=created_at&sort-order=asc", nil)
						router.ServeHTTP(w, req)
						response := decodeJournalList(w, http.StatusOK)
						journals := response.Journals
						for i := 1; i < len(journals); i++ {
							prevTime := journals[i-1].CreatedAt
							currTime := journals[i].CreatedAt
							Expect(prevTime).To(BeTemporally("<=", currTime))
						}
					})

					It("should sort by created_at descending (default)", func() {
						req, w = util.CreateTestRequest("GET", barkat.JournalBase+"?sort-by=created_at&sort-order=desc", nil)
						router.ServeHTTP(w, req)
						response := decodeJournalList(w, http.StatusOK)
						journals := response.Journals
						for i := 1; i < len(journals); i++ {
							prevTime := journals[i-1].CreatedAt
							currTime := journals[i].CreatedAt
							Expect(prevTime).To(BeTemporally(">=", currTime))
						}
					})
				})

				Context("Bad Values", func() {
					It("should return 400 for invalid sort-by field", func() {
						req, w = util.CreateTestRequest("GET", barkat.JournalBase+"?sort-by=invalid", nil)
						router.ServeHTTP(w, req)
						util.AssertError(w, "SortBy", "oneof")
					})

					It("should return 400 for invalid sort-order value", func() {
						req, w = util.CreateTestRequest("GET", barkat.JournalBase+"?sort-order=invalid", nil)
						router.ServeHTTP(w, req)
						util.AssertError(w, "SortOrder", "oneof")
					})
				})
			})

			Context("Pagination", func() {
				Context("Allowed Values", func() {
					It("should limit results with limit = 2", func() {
						req, w = util.CreateTestRequest("GET", barkat.JournalBase+"?limit=2", nil)
						router.ServeHTTP(w, req)
						response := decodeJournalList(w, http.StatusOK)
						Expect(response.Journals).To(HaveLen(2))
						Expect(response.Metadata.Total).To(Equal(int64(5)))
						Expect(response.Metadata.Offset).To(Equal(0))
						Expect(response.Metadata.Limit).To(Equal(2))
					})

					It("should skip entries with offset = 2, limit = 2", func() {
						req, w = util.CreateTestRequest("GET", barkat.JournalBase+"?offset=2&limit=2", nil)
						router.ServeHTTP(w, req)
						response := decodeJournalList(w, http.StatusOK)
						Expect(response.Journals).To(HaveLen(2))
						Expect(response.Metadata.Total).To(Equal(int64(5)))
						Expect(response.Metadata.Offset).To(Equal(2))
						Expect(response.Metadata.Limit).To(Equal(2))
					})

					It("should return last journal with offset = 4, limit = 2", func() {
						req, w = util.CreateTestRequest("GET", barkat.JournalBase+"?offset=4&limit=2", nil)
						router.ServeHTTP(w, req)
						response := decodeJournalList(w, http.StatusOK)
						Expect(response.Journals).To(HaveLen(1))
						Expect(response.Metadata.Total).To(Equal(int64(5)))
						Expect(response.Metadata.Offset).To(Equal(4))
						Expect(response.Metadata.Limit).To(Equal(2))
					})

					It("should return empty list for offset beyond total", func() {
						req, w = util.CreateTestRequest("GET", barkat.JournalBase+"?offset=10", nil)
						router.ServeHTTP(w, req)
						response := decodeJournalList(w, http.StatusOK)
						Expect(response.Journals).To(BeEmpty())
						Expect(response.Metadata.Total).To(Equal(int64(5)))
						Expect(response.Metadata.Offset).To(Equal(10))
						Expect(response.Metadata.Limit).To(Equal(20))
					})

					It("should accept limit = 1 (minimum)", func() {
						req, w = util.CreateTestRequest("GET", barkat.JournalBase+"?limit=1", nil)
						router.ServeHTTP(w, req)
						response := decodeJournalList(w, http.StatusOK)
						Expect(response.Journals).To(HaveLen(1))
						Expect(response.Metadata.Total).To(Equal(int64(5)))
						Expect(response.Metadata.Offset).To(Equal(0))
						Expect(response.Metadata.Limit).To(Equal(1))
					})

					It("should accept limit = 100 (maximum)", func() {
						req, w = util.CreateTestRequest("GET", barkat.JournalBase+"?limit=100", nil)
						router.ServeHTTP(w, req)
						response := decodeJournalList(w, http.StatusOK)
						Expect(response.Journals).To(HaveLen(5))
						Expect(response.Metadata.Total).To(Equal(int64(5)))
						Expect(response.Metadata.Offset).To(Equal(0))
						Expect(response.Metadata.Limit).To(Equal(100))
					})
				})

				Context("Bad Values", func() {
					It("should return 400 for limit exceeds maximum (101)", func() {
						req, w = util.CreateTestRequest("GET", barkat.JournalBase+"?limit=101", nil)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Limit", "max")
					})

					It("should return 400 for limit = 0 (fails min=1 validation)", func() {
						req, w = util.CreateTestRequest("GET", barkat.JournalBase+"?limit=0", nil)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Limit", "min")
					})

					It("should return 400 for negative limit", func() {
						req, w = util.CreateTestRequest("GET", barkat.JournalBase+"?limit=-1", nil)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Limit", "min")
					})

					It("should return 400 for negative offset", func() {
						req, w = util.CreateTestRequest("GET", barkat.JournalBase+"?offset=-1", nil)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Offset", "min")
					})

					It("should return 400 for non-numeric limit", func() {
						req, w = util.CreateTestRequest("GET", barkat.JournalBase+"?limit=abc", nil)
						router.ServeHTTP(w, req)
						util.AssertError(w, "message", "numeric")
					})

					It("should return 400 for non-numeric offset", func() {
						req, w = util.CreateTestRequest("GET", barkat.JournalBase+"?offset=xyz", nil)
						router.ServeHTTP(w, req)
						util.AssertError(w, "message", "numeric")
					})
				})
			})
		})

		Context("Errors", func() {
			Context("empty database", func() {
				BeforeEach(func() {
					sqlDB, _ := db.DB()
					sqlDB.Close()

					var err error
					db, err = core.CreateTestBarkatDB()
					Expect(err).ToNot(HaveOccurred())

					journalRepo := repository.NewJournalRepository(util.NewBaseDbRepository(db))
					journalMgr = manager.NewJournalManager(journalRepo)
					journalHandler = handler.NewJournalHandler(journalMgr)

					router = util.CreateTestGinRouter()
					v1 := router.Group("/v1/api")
					journal := v1.Group("/journals")
					handler.SetupJournalRoutes(journal, journalHandler)

					req, w = util.CreateTestRequest("GET", barkat.JournalBase, nil)
					router.ServeHTTP(w, req)
				})

				It("should return empty list", func() {
					response := decodeJournalList(w, http.StatusOK)
					Expect(response.Journals).To(BeEmpty())
					Expect(response.Metadata.Total).To(Equal(int64(0)))
					Expect(response.Metadata.Offset).To(Equal(0))
					Expect(response.Metadata.Limit).To(Equal(20))
				})
			})
		})
	})
})
