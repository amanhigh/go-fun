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

	Describe("GET /v1/journal/{id} - Retrieve Entry", func() {
		var createdJournal barkat.Journal

		BeforeEach(func() {
			journal := barkat.Journal{
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
			Expect(journalMgr.CreateJournal(testCtx, &journal)).To(Succeed())
			createdJournal = journal
		})

		Context("Happy Path", func() {
			Context("with valid entry ID", func() {
				var response barkat.Journal

				BeforeEach(func() {
					req, w = util.CreateTestRequest("GET", barkat.JournalEntries+"/"+createdJournal.ExternalID, nil)
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
					Expect(response.Sequence).To(Equal("MWD"))
					Expect(response.Type).To(Equal("REJECTED"))
					Expect(response.Status).To(Equal("FAIL"))
					Expect(response.CreatedAt).ToNot(BeZero())
					Expect(response.Images).To(HaveLen(4))
					Expect(response.Images[0].Timeframe).To(Equal("DL"))
					Expect(response.Images[1].Timeframe).To(Equal("WK"))
					Expect(response.Images[2].Timeframe).To(Equal("MN"))
					Expect(response.Images[3].Timeframe).To(Equal("TMN"))
				})
			})
		})

		Context("Field Validations", func() {
			Context("Entry ID Field", func() {
				Context("Bad Values", func() {
					It("should return 404 for non-existent entry ID", func() {
						req, w = util.CreateTestRequest("GET", barkat.JournalEntries+"/nonexistent-id", nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusNotFound))
					})

					It("should return 404 for malformed UUID", func() {
						req, w = util.CreateTestRequest("GET", barkat.JournalEntries+"/invalid-uuid-format", nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusNotFound))
					})

					It("should return 404 for valid UUID format but non-existent", func() {
						req, w = util.CreateTestRequest("GET", barkat.JournalEntries+"/550e8400-e29b-41d4-a716-446655440000", nil)
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
				{Timeframe: "DL"},
				{Timeframe: "WK"},
				{Timeframe: "MN"},
				{Timeframe: "TMN"},
			}

			journals := []barkat.Journal{
				{Ticker: "GRSE", Sequence: "MWD", Type: "REJECTED", Status: "FAIL"},
				{Ticker: "PDSL", Sequence: "YR", Type: "SET", Status: "TAKEN"},
				{Ticker: "SNF", Sequence: "MWD", Type: "RESULT", Status: "SUCCESS"},
				{Ticker: "TCS", Sequence: "YR", Type: "REJECTED", Status: "REJECTED"},
				{Ticker: "INFY", Sequence: "MWD", Type: "SET", Status: "RUNNING"},
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
					req, w = util.CreateTestRequest("GET", barkat.JournalEntries, nil)
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
						Expect(journal.Sequence).ToNot(BeEmpty())
						Expect(journal.Type).ToNot(BeEmpty())
						Expect(journal.Status).ToNot(BeEmpty())
						Expect(journal.CreatedAt).ToNot(BeZero())
						Expect(journal.Images).To(HaveLen(4))
						Expect(journal.Images[0].Timeframe).To(Equal("DL"))
						Expect(journal.Images[1].Timeframe).To(Equal("WK"))
						Expect(journal.Images[2].Timeframe).To(Equal("MN"))
						Expect(journal.Images[3].Timeframe).To(Equal("TMN"))
					}
				})
			})
		})

		Context("Field Validations", func() {
			Context("Ticker Filter", func() {
				Context("Allowed Values", func() {
					It("should filter by exact ticker match", func() {
						req, w = util.CreateTestRequest("GET", barkat.JournalEntries+"?ticker=GRSE", nil)
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
						req, w = util.CreateTestRequest("GET", barkat.JournalEntries+"?ticker=NOTFOUND", nil)
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
						req, w = util.CreateTestRequest("GET", barkat.JournalEntries+"?ticker=1234567890123456789012345678901", nil)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Ticker", "max")
					})

					It("should return 400 for lowercase ticker (PRD: uppercase only)", func() {
						req, w = util.CreateTestRequest("GET", barkat.JournalEntries+"?ticker=grse", nil)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Ticker", "ticker")
					})
				})
			})

			Context("Type Filter", func() {
				Context("Allowed Values", func() {
					It("should filter by type = REJECTED", func() {
						req, w = util.CreateTestRequest("GET", barkat.JournalEntries+"?type=REJECTED", nil)
						router.ServeHTTP(w, req)
						response := decodeJournalList(w, http.StatusOK)
						Expect(response.Journals).To(HaveLen(2))
						for _, journal := range response.Journals {
							Expect(journal.Type).To(Equal("REJECTED"))
						}
					})

					It("should filter by type = SET", func() {
						req, w = util.CreateTestRequest("GET", barkat.JournalEntries+"?type=SET", nil)
						router.ServeHTTP(w, req)
						response := decodeJournalList(w, http.StatusOK)
						Expect(response.Journals).To(HaveLen(2))
						for _, journal := range response.Journals {
							Expect(journal.Type).To(Equal("SET"))
						}
					})

					It("should filter by type = RESULT", func() {
						req, w = util.CreateTestRequest("GET", barkat.JournalEntries+"?type=RESULT", nil)
						router.ServeHTTP(w, req)
						response := decodeJournalList(w, http.StatusOK)
						Expect(response.Journals).To(HaveLen(1))
						journals := response.Journals
						Expect(journals[0].Type).To(Equal("RESULT"))
					})
				})

				Context("Bad Values", func() {
					It("should return 400 for invalid type enum", func() {
						req, w = util.CreateTestRequest("GET", barkat.JournalEntries+"?type=invalid", nil)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Type", "oneof")
					})
				})
			})

			Context("Status Filter", func() {
				Context("Allowed Values", func() {
					It("should filter by status = FAIL", func() {
						req, w = util.CreateTestRequest("GET", barkat.JournalEntries+"?status=FAIL", nil)
						router.ServeHTTP(w, req)
						response := decodeJournalList(w, http.StatusOK)
						Expect(response.Journals).To(HaveLen(1))
						journals := response.Journals
						Expect(journals[0].Status).To(Equal("FAIL"))
					})

					It("should filter by status = TAKEN", func() {
						req, w = util.CreateTestRequest("GET", barkat.JournalEntries+"?status=TAKEN", nil)
						router.ServeHTTP(w, req)
						response := decodeJournalList(w, http.StatusOK)
						Expect(response.Journals).To(HaveLen(1))
						Expect(response.Journals[0].Status).To(Equal("TAKEN"))
					})

					It("should filter by status = SUCCESS", func() {
						req, w = util.CreateTestRequest("GET", barkat.JournalEntries+"?status=SUCCESS", nil)
						router.ServeHTTP(w, req)
						response := decodeJournalList(w, http.StatusOK)
						Expect(response.Journals).To(HaveLen(1))
						Expect(response.Journals[0].Status).To(Equal("SUCCESS"))
					})

					It("should filter by status = RUNNING", func() {
						req, w = util.CreateTestRequest("GET", barkat.JournalEntries+"?status=RUNNING", nil)
						router.ServeHTTP(w, req)
						response := decodeJournalList(w, http.StatusOK)
						Expect(response.Journals).To(HaveLen(1))
						Expect(response.Journals[0].Status).To(Equal("RUNNING"))
					})

					It("should filter by status = REJECTED", func() {
						req, w = util.CreateTestRequest("GET", barkat.JournalEntries+"?status=REJECTED", nil)
						router.ServeHTTP(w, req)
						response := decodeJournalList(w, http.StatusOK)
						Expect(response.Journals).To(HaveLen(1))
						Expect(response.Journals[0].Status).To(Equal("REJECTED"))
					})
				})

				Context("Bad Values", func() {
					It("should return 400 for invalid status enum", func() {
						req, w = util.CreateTestRequest("GET", barkat.JournalEntries+"?status=invalid", nil)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Status", "oneof")
					})
				})
			})

			Context("Sequence Filter", func() {
				Context("Allowed Values", func() {
					It("should filter by sequence = MWD", func() {
						req, w = util.CreateTestRequest("GET", barkat.JournalEntries+"?sequence=MWD", nil)
						router.ServeHTTP(w, req)
						response := decodeJournalList(w, http.StatusOK)
						Expect(response.Journals).To(HaveLen(3))
						for _, journal := range response.Journals {
							Expect(journal.Sequence).To(Equal("MWD"))
						}
					})

					It("should filter by sequence = YR", func() {
						req, w = util.CreateTestRequest("GET", barkat.JournalEntries+"?sequence=YR", nil)
						router.ServeHTTP(w, req)
						response := decodeJournalList(w, http.StatusOK)
						Expect(response.Journals).To(HaveLen(2))
						for _, journal := range response.Journals {
							Expect(journal.Sequence).To(Equal("YR"))
						}
					})
				})

				Context("Bad Values", func() {
					It("should return 400 for invalid sequence enum", func() {
						req, w = util.CreateTestRequest("GET", barkat.JournalEntries+"?sequence=invalid", nil)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Sequence", "oneof")
					})
				})
			})

			Context("Combined Filters", func() {
				It("should apply ticker + type filters", func() {
					req, w = util.CreateTestRequest("GET", barkat.JournalEntries+"?ticker=GRSE&type=REJECTED", nil)
					router.ServeHTTP(w, req)
					response := decodeJournalList(w, http.StatusOK)
					Expect(response.Journals).To(HaveLen(1))
					journals := response.Journals
					Expect(journals[0].Ticker).To(Equal("GRSE"))
					Expect(journals[0].Type).To(Equal("REJECTED"))
				})

				It("should apply sequence + status filters", func() {
					req, w = util.CreateTestRequest("GET", barkat.JournalEntries+"?sequence=YR&status=TAKEN", nil)
					router.ServeHTTP(w, req)
					response := decodeJournalList(w, http.StatusOK)
					Expect(response.Journals).To(HaveLen(1))
					Expect(response.Journals[0].Sequence).To(Equal("YR"))
					Expect(response.Journals[0].Status).To(Equal("TAKEN"))
				})

				It("should apply type + status + sequence filters", func() {
					req, w = util.CreateTestRequest("GET", barkat.JournalEntries+"?type=SET&status=RUNNING&sequence=MWD", nil)
					router.ServeHTTP(w, req)
					response := decodeJournalList(w, http.StatusOK)
					Expect(response.Journals).To(HaveLen(1))
					Expect(response.Journals[0].Type).To(Equal("SET"))
					Expect(response.Journals[0].Status).To(Equal("RUNNING"))
					Expect(response.Journals[0].Sequence).To(Equal("MWD"))
				})
			})

			Context("Date Fields", func() {
				Context("Created-After Field", func() {
					Context("Allowed Values", func() {
						It("should accept valid YYYY-MM-DD date and filter entries", func() {
							// Use yesterday's date to ensure all entries are captured
							afterDate := time.Now().Add(-24 * time.Hour).Format("2006-01-02")
							req, w = util.CreateTestRequest("GET", barkat.JournalEntries+"?created-after="+afterDate, nil)
							router.ServeHTTP(w, req)
							response := decodeJournalList(w, http.StatusOK)
							// All entries created in this test should be returned
							Expect(response.Journals).To(HaveLen(5))
						})

						It("should return empty list for future date", func() {
							futureDate := time.Now().Add(24 * time.Hour).Format("2006-01-02")
							req, w = util.CreateTestRequest("GET", barkat.JournalEntries+"?created-after="+futureDate, nil)
							router.ServeHTTP(w, req)
							response := decodeJournalList(w, http.StatusOK)
							Expect(response.Journals).To(BeEmpty())
						})

						It("should work with created-before combined filter", func() {
							afterDate := time.Now().Add(-24 * time.Hour).Format("2006-01-02")
							beforeDate := time.Now().Add(24 * time.Hour).Format("2006-01-02")
							req, w = util.CreateTestRequest("GET", barkat.JournalEntries+"?created-after="+afterDate+"&created-before="+beforeDate, nil)
							router.ServeHTTP(w, req)
							response := decodeJournalList(w, http.StatusOK)
							Expect(response.Journals).To(HaveLen(5))
						})
					})

					Context("Bad Values", func() {
						It("should return 400 for invalid date format", func() {
							req, w = util.CreateTestRequest("GET", barkat.JournalEntries+"?created-after=invalid-date", nil)
							router.ServeHTTP(w, req)
							util.AssertError(w, "CreatedAfter", "Violates 'datetime")
						})

						It("should ignore empty date and return all entries", func() {
							req, w = util.CreateTestRequest("GET", barkat.JournalEntries+"?created-after=", nil)
							router.ServeHTTP(w, req)
							response := decodeJournalList(w, http.StatusOK)
							// Empty date is ignored, all entries returned
							Expect(response.Journals).To(HaveLen(5))
						})

						It("should return 400 for wrong format (DD-MM-YYYY)", func() {
							req, w = util.CreateTestRequest("GET", barkat.JournalEntries+"?created-after=15-02-2024", nil)
							router.ServeHTTP(w, req)
							util.AssertError(w, "CreatedAfter", "Violates 'datetime")
						})
					})
				})

				Context("Created-Before Field", func() {
					Context("Allowed Values", func() {
						It("should accept valid YYYY-MM-DD date and filter entries", func() {
							beforeDate := time.Now().Add(24 * time.Hour).Format("2006-01-02")
							req, w = util.CreateTestRequest("GET", barkat.JournalEntries+"?created-before="+beforeDate, nil)
							router.ServeHTTP(w, req)
							response := decodeJournalList(w, http.StatusOK)
							// All entries created in this test should be returned
							Expect(response.Journals).To(HaveLen(5))
						})

						It("should return empty list for past date", func() {
							pastDate := time.Now().Add(-48 * time.Hour).Format("2006-01-02")
							req, w = util.CreateTestRequest("GET", barkat.JournalEntries+"?created-before="+pastDate, nil)
							router.ServeHTTP(w, req)
							response := decodeJournalList(w, http.StatusOK)
							Expect(response.Journals).To(BeEmpty())
						})

						It("should work with created-after combined filter", func() {
							afterDate := time.Now().Add(-24 * time.Hour).Format("2006-01-02")
							beforeDate := time.Now().Add(24 * time.Hour).Format("2006-01-02")
							req, w = util.CreateTestRequest("GET", barkat.JournalEntries+"?created-after="+afterDate+"&created-before="+beforeDate, nil)
							router.ServeHTTP(w, req)
							response := decodeJournalList(w, http.StatusOK)
							Expect(response.Journals).To(HaveLen(5))
						})
					})

					Context("Bad Values", func() {
						It("should return 400 for invalid date format", func() {
							req, w = util.CreateTestRequest("GET", barkat.JournalEntries+"?created-before=not-a-date", nil)
							router.ServeHTTP(w, req)
							util.AssertError(w, "CreatedBefore", "Violates 'datetime")
						})

						It("should ignore empty date and return all entries", func() {
							req, w = util.CreateTestRequest("GET", barkat.JournalEntries+"?created-before=", nil)
							router.ServeHTTP(w, req)
							response := decodeJournalList(w, http.StatusOK)
							// Empty date is ignored, all entries returned
							Expect(response.Journals).To(HaveLen(5))
						})

						It("should return 400 for wrong format (DD-MM-YYYY)", func() {
							req, w = util.CreateTestRequest("GET", barkat.JournalEntries+"?created-before=15-02-2024", nil)
							router.ServeHTTP(w, req)
							util.AssertError(w, "CreatedBefore", "Violates 'datetime")
						})
					})
				})
			})

			Context("Sorting", func() {
				Context("Allowed Values", func() {
					It("should sort by ticker ascending", func() {
						req, w = util.CreateTestRequest("GET", barkat.JournalEntries+"?sort-by=ticker&sort-order=asc", nil)
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
						req, w = util.CreateTestRequest("GET", barkat.JournalEntries+"?sort-by=ticker&sort-order=desc", nil)
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

					It("should sort by sequence ascending", func() {
						req, w = util.CreateTestRequest("GET", barkat.JournalEntries+"?sort-by=sequence&sort-order=asc", nil)
						router.ServeHTTP(w, req)
						response := decodeJournalList(w, http.StatusOK)
						Expect(response.Journals).To(HaveLen(5))
						journals := response.Journals
						for i := range 3 {
							Expect(journals[i].Sequence).To(Equal("MWD"))
						}
						for i := 3; i < 5; i++ {
							Expect(journals[i].Sequence).To(Equal("YR"))
						}
					})

					It("should sort by created_at ascending", func() {
						req, w = util.CreateTestRequest("GET", barkat.JournalEntries+"?sort-by=created_at&sort-order=asc", nil)
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
						req, w = util.CreateTestRequest("GET", barkat.JournalEntries+"?sort-by=created_at&sort-order=desc", nil)
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
						req, w = util.CreateTestRequest("GET", barkat.JournalEntries+"?sort-by=invalid", nil)
						router.ServeHTTP(w, req)
						util.AssertError(w, "SortBy", "oneof")
					})

					It("should return 400 for invalid sort-order value", func() {
						req, w = util.CreateTestRequest("GET", barkat.JournalEntries+"?sort-order=invalid", nil)
						router.ServeHTTP(w, req)
						util.AssertError(w, "SortOrder", "oneof")
					})
				})
			})

			Context("Pagination", func() {
				Context("Allowed Values", func() {
					It("should limit results with limit = 2", func() {
						req, w = util.CreateTestRequest("GET", barkat.JournalEntries+"?limit=2", nil)
						router.ServeHTTP(w, req)
						response := decodeJournalList(w, http.StatusOK)
						Expect(response.Journals).To(HaveLen(2))
						Expect(response.Metadata.Total).To(Equal(int64(5)))
						Expect(response.Metadata.Offset).To(Equal(0))
						Expect(response.Metadata.Limit).To(Equal(2))
					})

					It("should skip entries with offset = 2, limit = 2", func() {
						req, w = util.CreateTestRequest("GET", barkat.JournalEntries+"?offset=2&limit=2", nil)
						router.ServeHTTP(w, req)
						response := decodeJournalList(w, http.StatusOK)
						Expect(response.Journals).To(HaveLen(2))
						Expect(response.Metadata.Total).To(Equal(int64(5)))
						Expect(response.Metadata.Offset).To(Equal(2))
						Expect(response.Metadata.Limit).To(Equal(2))
					})

					It("should return last entry with offset = 4, limit = 2", func() {
						req, w = util.CreateTestRequest("GET", barkat.JournalEntries+"?offset=4&limit=2", nil)
						router.ServeHTTP(w, req)
						response := decodeJournalList(w, http.StatusOK)
						Expect(response.Journals).To(HaveLen(1))
						Expect(response.Metadata.Total).To(Equal(int64(5)))
						Expect(response.Metadata.Offset).To(Equal(4))
						Expect(response.Metadata.Limit).To(Equal(2))
					})

					It("should return empty list for offset beyond total", func() {
						req, w = util.CreateTestRequest("GET", barkat.JournalEntries+"?offset=10", nil)
						router.ServeHTTP(w, req)
						response := decodeJournalList(w, http.StatusOK)
						Expect(response.Journals).To(BeEmpty())
						Expect(response.Metadata.Total).To(Equal(int64(5)))
						Expect(response.Metadata.Offset).To(Equal(10))
						Expect(response.Metadata.Limit).To(Equal(20))
					})

					It("should accept limit = 1 (minimum)", func() {
						req, w = util.CreateTestRequest("GET", barkat.JournalEntries+"?limit=1", nil)
						router.ServeHTTP(w, req)
						response := decodeJournalList(w, http.StatusOK)
						Expect(response.Journals).To(HaveLen(1))
						Expect(response.Metadata.Total).To(Equal(int64(5)))
						Expect(response.Metadata.Offset).To(Equal(0))
						Expect(response.Metadata.Limit).To(Equal(1))
					})

					It("should accept limit = 100 (maximum)", func() {
						req, w = util.CreateTestRequest("GET", barkat.JournalEntries+"?limit=100", nil)
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
						req, w = util.CreateTestRequest("GET", barkat.JournalEntries+"?limit=101", nil)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Limit", "max")
					})

					It("should return 400 for limit = 0 (fails min=1 validation)", func() {
						req, w = util.CreateTestRequest("GET", barkat.JournalEntries+"?limit=0", nil)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Limit", "min")
					})

					It("should return 400 for negative limit", func() {
						req, w = util.CreateTestRequest("GET", barkat.JournalEntries+"?limit=-1", nil)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Limit", "min")
					})

					It("should return 400 for negative offset", func() {
						req, w = util.CreateTestRequest("GET", barkat.JournalEntries+"?offset=-1", nil)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Offset", "min")
					})

					It("should return 400 for non-numeric limit", func() {
						req, w = util.CreateTestRequest("GET", barkat.JournalEntries+"?limit=abc", nil)
						router.ServeHTTP(w, req)
						// BUG: Gin returns strconv.NumError without field context, so check generic message
						util.AssertError(w, "message", "numeric")
					})

					It("should return 400 for non-numeric offset", func() {
						req, w = util.CreateTestRequest("GET", barkat.JournalEntries+"?offset=xyz", nil)
						router.ServeHTTP(w, req)
						// Gin returns strconv.NumError without field context, so check generic message
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

					journalRepo := repository.NewJournalRepository(db)
					journalMgr = manager.NewJournalManager(journalRepo)
					journalHandler = handler.NewJournalHandler(journalMgr)

					router = util.CreateTestGinRouter()
					v1 := router.Group("/v1")
					journal := v1.Group("/journals")
					handler.SetupJournalEntryRoutes(journal, journalHandler)

					req, w = util.CreateTestRequest("GET", barkat.JournalEntries, nil)
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
