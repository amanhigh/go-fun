package handler_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"time"

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

func decodeEntry(w *httptest.ResponseRecorder, expectedStatus int) barkat.Entry {
	return util.UnenvelopeAndAssertStatus[barkat.Entry](w, expectedStatus)
}

func decodeEntryList(w *httptest.ResponseRecorder, expectedStatus int) barkat.EntryList {
	return util.UnenvelopeAndAssertStatus[barkat.EntryList](w, expectedStatus)
}

var _ = Describe("JournalHandler Integration - GET Tests", func() {
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
		journal := v1.Group("/journal")
		handler.SetupJournalEntryRoutes(journal, journalHandler)
	})

	AfterEach(func() {
		sqlDB, err := db.DB()
		Expect(err).ToNot(HaveOccurred())
		sqlDB.Close()
	})

	Describe("GET /v1/journal-entries/{id} - Retrieve Entry", func() {
		var createdEntry barkat.Entry

		BeforeEach(func() {
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
			}
			Expect(entryMgr.CreateEntry(testCtx, &entry)).To(Succeed())
			createdEntry = entry
		})

		Context("Happy Path", func() {
			Context("with valid entry ID", func() {
				var response barkat.Entry

				BeforeEach(func() {
					req, w = util.CreateTestRequest("GET", JournalEntriesBaseURL+"/"+createdEntry.ID, nil)
					router.ServeHTTP(w, req)
				})

				It("should return 200 OK", func() {
					Expect(w.Code).To(Equal(http.StatusOK))
				})

				It("should return entry with correct ID", func() {
					response = decodeEntry(w, http.StatusOK)
					Expect(response.ID).To(Equal(createdEntry.ID))
				})

				It("should return all entry fields", func() {
					response = decodeEntry(w, http.StatusOK)
					Expect(response.Ticker).To(Equal("GRSE"))
					Expect(response.Sequence).To(Equal("MWD"))
					Expect(response.Type).To(Equal("REJECTED"))
					Expect(response.Status).To(Equal("FAIL"))
					Expect(response.CreatedAt).ToNot(BeZero())
				})
			})
		})

		Context("Field Validations", func() {
			Context("Entry ID Field", func() {
				Context("Bad Values", func() {
					It("should return 404 for non-existent entry ID", func() {
						req, w = util.CreateTestRequest("GET", JournalEntriesBaseURL+"/nonexistent-id", nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusNotFound))
					})

					It("should return 404 for malformed UUID", func() {
						req, w = util.CreateTestRequest("GET", JournalEntriesBaseURL+"/invalid-uuid-format", nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusNotFound))
					})

					It("should return 404 for valid UUID format but non-existent", func() {
						req, w = util.CreateTestRequest("GET", JournalEntriesBaseURL+"/550e8400-e29b-41d4-a716-446655440000", nil)
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

	Describe("GET /v1/journal-entries - List Entries", func() {
		var createdEntries []barkat.Entry

		BeforeEach(func() {
			defaultImages := []barkat.Image{
				{Timeframe: "DL"},
				{Timeframe: "WK"},
				{Timeframe: "MN"},
				{Timeframe: "TMN"},
			}
			entries := []barkat.Entry{
				{Ticker: "GRSE", Sequence: "MWD", Type: "REJECTED", Status: "FAIL", Images: defaultImages},
				{Ticker: "PDSL", Sequence: "YR", Type: "SET", Status: "TAKEN", Images: defaultImages},
				{Ticker: "SNF", Sequence: "MWD", Type: "RESULT", Status: "SUCCESS", Images: defaultImages},
				{Ticker: "TCS", Sequence: "YR", Type: "REJECTED", Status: "REJECTED", Images: defaultImages},
				{Ticker: "INFY", Sequence: "MWD", Type: "SET", Status: "RUNNING", Images: defaultImages},
			}
			for _, entry := range entries {
				Expect(entryMgr.CreateEntry(testCtx, &entry)).To(Succeed())
				createdEntries = append(createdEntries, entry)
			}
		})

		Context("Happy Path", func() {
			Context("default pagination (no filters)", func() {
				var response barkat.EntryList

				BeforeEach(func() {
					req, w = util.CreateTestRequest("GET", JournalEntriesBaseURL, nil)
					router.ServeHTTP(w, req)
				})

				It("should return 200 OK", func() {
					Expect(w.Code).To(Equal(http.StatusOK))
				})

				It("should return all entries", func() {
					response = decodeEntryList(w, http.StatusOK)
					Expect(response.Records).To(HaveLen(5))
				})

				It("should return correct total count", func() {
					response = decodeEntryList(w, http.StatusOK)
					Expect(response.Metadata.Total).To(Equal(int64(5)))
				})

				It("should return entries in reverse chronological order by default", func() {
					response = decodeEntryList(w, http.StatusOK)
					for i := 1; i < len(response.Records); i++ {
						prevTime := response.Records[i-1].CreatedAt
						currTime := response.Records[i].CreatedAt
						Expect(prevTime).To(BeTemporally(">=", currTime))
					}
				})

				It("should include all required fields in each entry", func() {
					response = decodeEntryList(w, http.StatusOK)
					for _, entry := range response.Records {
						Expect(entry.ID).ToNot(BeEmpty())
						Expect(entry.Ticker).ToNot(BeEmpty())
						Expect(entry.Sequence).ToNot(BeEmpty())
						Expect(entry.Type).ToNot(BeEmpty())
						Expect(entry.Status).ToNot(BeEmpty())
						Expect(entry.CreatedAt).ToNot(BeZero())
					}
				})
			})
		})

		Context("Field Validations", func() {
			Context("Ticker Filter", func() {
				Context("Allowed Values", func() {
					It("should filter by exact ticker match", func() {
						req, w = util.CreateTestRequest("GET", JournalEntriesBaseURL+"?ticker=GRSE", nil)
						router.ServeHTTP(w, req)
						response := decodeEntryList(w, http.StatusOK)
						Expect(response.Records).To(HaveLen(1))
						Expect(response.Records[0].Ticker).To(Equal("GRSE"))
						Expect(response.Metadata.Total).To(Equal(int64(1)))
					})

					It("should return empty list for ticker with no matches", func() {
						req, w = util.CreateTestRequest("GET", JournalEntriesBaseURL+"?ticker=NOTFOUND", nil)
						router.ServeHTTP(w, req)
						response := decodeEntryList(w, http.StatusOK)
						Expect(response.Records).To(BeEmpty())
						Expect(response.Metadata.Total).To(Equal(int64(0)))
					})
				})

				Context("Bad Values", func() {
					It("should return 400 for invalid ticker length", func() {
						req, w = util.CreateTestRequest("GET", JournalEntriesBaseURL+"?ticker=1234567890123456789012345678901", nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusBadRequest))
					})
				})
			})

			Context("Type Filter", func() {
				Context("Allowed Values", func() {
					It("should filter by type = REJECTED", func() {
						req, w = util.CreateTestRequest("GET", JournalEntriesBaseURL+"?type=REJECTED", nil)
						router.ServeHTTP(w, req)
						response := decodeEntryList(w, http.StatusOK)
						Expect(response.Records).To(HaveLen(2))
						for _, entry := range response.Records {
							Expect(entry.Type).To(Equal("REJECTED"))
						}
					})

					It("should filter by type = SET", func() {
						req, w = util.CreateTestRequest("GET", JournalEntriesBaseURL+"?type=SET", nil)
						router.ServeHTTP(w, req)
						response := decodeEntryList(w, http.StatusOK)
						Expect(response.Records).To(HaveLen(2))
						for _, entry := range response.Records {
							Expect(entry.Type).To(Equal("SET"))
						}
					})

					It("should filter by type = RESULT", func() {
						req, w = util.CreateTestRequest("GET", JournalEntriesBaseURL+"?type=RESULT", nil)
						router.ServeHTTP(w, req)
						response := decodeEntryList(w, http.StatusOK)
						Expect(response.Records).To(HaveLen(1))
						Expect(response.Records[0].Type).To(Equal("RESULT"))
					})
				})

				Context("Bad Values", func() {
					It("should return 400 for invalid type enum", func() {
						req, w = util.CreateTestRequest("GET", JournalEntriesBaseURL+"?type=invalid", nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusBadRequest))
					})
				})
			})

			Context("Status Filter", func() {
				Context("Allowed Values", func() {
					It("should filter by status = FAIL", func() {
						req, w = util.CreateTestRequest("GET", JournalEntriesBaseURL+"?status=FAIL", nil)
						router.ServeHTTP(w, req)
						response := decodeEntryList(w, http.StatusOK)
						Expect(response.Records).To(HaveLen(1))
						Expect(response.Records[0].Status).To(Equal("FAIL"))
					})

					It("should filter by status = TAKEN", func() {
						req, w = util.CreateTestRequest("GET", JournalEntriesBaseURL+"?status=TAKEN", nil)
						router.ServeHTTP(w, req)
						response := decodeEntryList(w, http.StatusOK)
						Expect(response.Records).To(HaveLen(1))
						Expect(response.Records[0].Status).To(Equal("TAKEN"))
					})

					It("should filter by status = SUCCESS", func() {
						req, w = util.CreateTestRequest("GET", JournalEntriesBaseURL+"?status=SUCCESS", nil)
						router.ServeHTTP(w, req)
						response := decodeEntryList(w, http.StatusOK)
						Expect(response.Records).To(HaveLen(1))
						Expect(response.Records[0].Status).To(Equal("SUCCESS"))
					})

					It("should filter by status = RUNNING", func() {
						req, w = util.CreateTestRequest("GET", JournalEntriesBaseURL+"?status=RUNNING", nil)
						router.ServeHTTP(w, req)
						response := decodeEntryList(w, http.StatusOK)
						Expect(response.Records).To(HaveLen(1))
						Expect(response.Records[0].Status).To(Equal("RUNNING"))
					})

					It("should filter by status = REJECTED", func() {
						req, w = util.CreateTestRequest("GET", JournalEntriesBaseURL+"?status=REJECTED", nil)
						router.ServeHTTP(w, req)
						response := decodeEntryList(w, http.StatusOK)
						Expect(response.Records).To(HaveLen(1))
						Expect(response.Records[0].Status).To(Equal("REJECTED"))
					})
				})

				Context("Bad Values", func() {
					It("should return 400 for invalid status enum", func() {
						req, w = util.CreateTestRequest("GET", JournalEntriesBaseURL+"?status=invalid", nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusBadRequest))
					})
				})
			})

			Context("Sequence Filter", func() {
				Context("Allowed Values", func() {
					It("should filter by sequence = MWD", func() {
						req, w = util.CreateTestRequest("GET", JournalEntriesBaseURL+"?sequence=MWD", nil)
						router.ServeHTTP(w, req)
						response := decodeEntryList(w, http.StatusOK)
						Expect(response.Records).To(HaveLen(3))
						for _, entry := range response.Records {
							Expect(entry.Sequence).To(Equal("MWD"))
						}
					})

					It("should filter by sequence = YR", func() {
						req, w = util.CreateTestRequest("GET", JournalEntriesBaseURL+"?sequence=YR", nil)
						router.ServeHTTP(w, req)
						response := decodeEntryList(w, http.StatusOK)
						Expect(response.Records).To(HaveLen(2))
						for _, entry := range response.Records {
							Expect(entry.Sequence).To(Equal("YR"))
						}
					})
				})

				Context("Bad Values", func() {
					It("should return 400 for invalid sequence enum", func() {
						req, w = util.CreateTestRequest("GET", JournalEntriesBaseURL+"?sequence=invalid", nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusBadRequest))
					})
				})
			})

			Context("Combined Filters", func() {
				It("should apply ticker + type filters", func() {
					req, w = util.CreateTestRequest("GET", JournalEntriesBaseURL+"?ticker=GRSE&type=REJECTED", nil)
					router.ServeHTTP(w, req)
					response := decodeEntryList(w, http.StatusOK)
					Expect(response.Records).To(HaveLen(1))
					Expect(response.Records[0].Ticker).To(Equal("GRSE"))
					Expect(response.Records[0].Type).To(Equal("REJECTED"))
				})

				It("should apply sequence + status filters", func() {
					req, w = util.CreateTestRequest("GET", JournalEntriesBaseURL+"?sequence=YR&status=TAKEN", nil)
					router.ServeHTTP(w, req)
					response := decodeEntryList(w, http.StatusOK)
					Expect(response.Records).To(HaveLen(1))
					Expect(response.Records[0].Sequence).To(Equal("YR"))
					Expect(response.Records[0].Status).To(Equal("TAKEN"))
				})

				It("should apply type + status + sequence filters", func() {
					req, w = util.CreateTestRequest("GET", JournalEntriesBaseURL+"?type=SET&status=RUNNING&sequence=MWD", nil)
					router.ServeHTTP(w, req)
					response := decodeEntryList(w, http.StatusOK)
					Expect(response.Records).To(HaveLen(1))
					Expect(response.Records[0].Type).To(Equal("SET"))
					Expect(response.Records[0].Status).To(Equal("RUNNING"))
					Expect(response.Records[0].Sequence).To(Equal("MWD"))
				})
			})

			PContext("Date Fields", func() {
				Context("Created-After Field", func() {
					Context("Allowed Values", func() {
						It("should accept valid ISO 8601 datetime and filter entries", func() {
							// Get current time and format as ISO 8601
							afterTime := time.Now().Add(-1 * time.Hour).Format(time.RFC3339)
							req, w = util.CreateTestRequest("GET", JournalEntriesBaseURL+"?created-after="+url.QueryEscape(afterTime), nil)
							router.ServeHTTP(w, req)
							response := decodeEntryList(w, http.StatusOK)
							// All entries created in this test should be returned
							Expect(response.Records).To(HaveLen(5))
						})

						It("should return empty list for future date", func() {
							futureTime := time.Now().Add(24 * time.Hour).Format(time.RFC3339)
							req, w = util.CreateTestRequest("GET", JournalEntriesBaseURL+"?created-after="+url.QueryEscape(futureTime), nil)
							router.ServeHTTP(w, req)
							response := decodeEntryList(w, http.StatusOK)
							Expect(response.Records).To(BeEmpty())
						})

						It("should work with created-before combined filter", func() {
							afterTime := time.Now().Add(-1 * time.Hour).Format(time.RFC3339)
							beforeTime := time.Now().Add(1 * time.Hour).Format(time.RFC3339)
							req, w = util.CreateTestRequest("GET", JournalEntriesBaseURL+"?created-after="+url.QueryEscape(afterTime)+"&created-before="+url.QueryEscape(beforeTime), nil)
							router.ServeHTTP(w, req)
							response := decodeEntryList(w, http.StatusOK)
							Expect(response.Records).To(HaveLen(5))
						})
					})

					Context("Bad Values", func() {
						It("should return 400 for invalid date format", func() {
							req, w = util.CreateTestRequest("GET", JournalEntriesBaseURL+"?created-after=invalid-date", nil)
							router.ServeHTTP(w, req)
							Expect(w.Code).To(Equal(http.StatusBadRequest))
						})

						It("should return 400 for empty date", func() {
							req, w = util.CreateTestRequest("GET", JournalEntriesBaseURL+"?created-after=", nil)
							router.ServeHTTP(w, req)
							Expect(w.Code).To(Equal(http.StatusBadRequest))
						})

						It("should return 400 for non-ISO format", func() {
							req, w = util.CreateTestRequest("GET", JournalEntriesBaseURL+"?created-after=2024-02-15", nil)
							router.ServeHTTP(w, req)
							Expect(w.Code).To(Equal(http.StatusBadRequest))
						})
					})
				})

				Context("Created-Before Field", func() {
					Context("Allowed Values", func() {
						It("should accept valid ISO 8601 datetime and filter entries", func() {
							beforeTime := time.Now().Add(1 * time.Hour).Format(time.RFC3339)
							req, w = util.CreateTestRequest("GET", JournalEntriesBaseURL+"?created-before="+url.QueryEscape(beforeTime), nil)
							router.ServeHTTP(w, req)
							var response barkat.EntryList
							util.AssertJSONAndStatus(w, http.StatusOK, &response)
							// All entries created in this test should be returned
							Expect(response.Records).To(HaveLen(5))
						})

						It("should return empty list for past date", func() {
							pastTime := time.Now().Add(-24 * time.Hour).Format(time.RFC3339)
							req, w = util.CreateTestRequest("GET", JournalEntriesBaseURL+"?created-before="+url.QueryEscape(pastTime), nil)
							router.ServeHTTP(w, req)
							response := decodeEntryList(w, http.StatusOK)
							Expect(response.Records).To(BeEmpty())
						})

						It("should work with created-after combined filter", func() {
							afterTime := time.Now().Add(-1 * time.Hour).Format(time.RFC3339)
							beforeTime := time.Now().Add(1 * time.Hour).Format(time.RFC3339)
							req, w = util.CreateTestRequest("GET", JournalEntriesBaseURL+"?created-after="+url.QueryEscape(afterTime)+"&created-before="+url.QueryEscape(beforeTime), nil)
							router.ServeHTTP(w, req)
							response := decodeEntryList(w, http.StatusOK)
							Expect(response.Records).To(HaveLen(5))
						})
					})

					Context("Bad Values", func() {
						It("should return 400 for invalid date format", func() {
							req, w = util.CreateTestRequest("GET", JournalEntriesBaseURL+"?created-before=not-a-date", nil)
							router.ServeHTTP(w, req)
							Expect(w.Code).To(Equal(http.StatusBadRequest))
						})

						It("should return 400 for empty date", func() {
							req, w = util.CreateTestRequest("GET", JournalEntriesBaseURL+"?created-before=", nil)
							router.ServeHTTP(w, req)
							Expect(w.Code).To(Equal(http.StatusBadRequest))
						})

						It("should return 400 for non-ISO format", func() {
							req, w = util.CreateTestRequest("GET", JournalEntriesBaseURL+"?created-before=15-02-2024", nil)
							router.ServeHTTP(w, req)
							Expect(w.Code).To(Equal(http.StatusBadRequest))
						})
					})
				})
			})

			Context("Sorting", func() {
				Context("Allowed Values", func() {
					It("should sort by ticker ascending", func() {
						req, w = util.CreateTestRequest("GET", JournalEntriesBaseURL+"?sort-by=ticker&sort-order=asc", nil)
						router.ServeHTTP(w, req)
						response := decodeEntryList(w, http.StatusOK)
						Expect(response.Records).To(HaveLen(5))
						Expect(response.Records[0].Ticker).To(Equal("GRSE"))
						Expect(response.Records[1].Ticker).To(Equal("INFY"))
						Expect(response.Records[2].Ticker).To(Equal("PDSL"))
						Expect(response.Records[3].Ticker).To(Equal("SNF"))
						Expect(response.Records[4].Ticker).To(Equal("TCS"))
					})

					It("should sort by ticker descending", func() {
						req, w = util.CreateTestRequest("GET", JournalEntriesBaseURL+"?sort-by=ticker&sort-order=desc", nil)
						router.ServeHTTP(w, req)
						response := decodeEntryList(w, http.StatusOK)
						Expect(response.Records).To(HaveLen(5))
						Expect(response.Records[0].Ticker).To(Equal("TCS"))
						Expect(response.Records[1].Ticker).To(Equal("SNF"))
						Expect(response.Records[2].Ticker).To(Equal("PDSL"))
						Expect(response.Records[3].Ticker).To(Equal("INFY"))
						Expect(response.Records[4].Ticker).To(Equal("GRSE"))
					})

					It("should sort by sequence ascending", func() {
						req, w = util.CreateTestRequest("GET", JournalEntriesBaseURL+"?sort-by=sequence&sort-order=asc", nil)
						router.ServeHTTP(w, req)
						response := decodeEntryList(w, http.StatusOK)
						Expect(response.Records).To(HaveLen(5))
						for i := 0; i < 3; i++ {
							Expect(response.Records[i].Sequence).To(Equal("MWD"))
						}
						for i := 3; i < 5; i++ {
							Expect(response.Records[i].Sequence).To(Equal("YR"))
						}
					})

					It("should sort by created_at ascending", func() {
						req, w = util.CreateTestRequest("GET", JournalEntriesBaseURL+"?sort-by=created_at&sort-order=asc", nil)
						router.ServeHTTP(w, req)
						response := decodeEntryList(w, http.StatusOK)
						for i := 1; i < len(response.Records); i++ {
							prevTime := response.Records[i-1].CreatedAt
							currTime := response.Records[i].CreatedAt
							Expect(prevTime).To(BeTemporally("<=", currTime))
						}
					})

					It("should sort by created_at descending (default)", func() {
						req, w = util.CreateTestRequest("GET", JournalEntriesBaseURL+"?sort-by=created_at&sort-order=desc", nil)
						router.ServeHTTP(w, req)
						response := decodeEntryList(w, http.StatusOK)
						for i := 1; i < len(response.Records); i++ {
							prevTime := response.Records[i-1].CreatedAt
							currTime := response.Records[i].CreatedAt
							Expect(prevTime).To(BeTemporally(">=", currTime))
						}
					})
				})

				Context("Bad Values", func() {
					It("should return 400 for invalid sort-by field", func() {
						req, w = util.CreateTestRequest("GET", JournalEntriesBaseURL+"?sort-by=invalid", nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusBadRequest))
					})

					It("should return 400 for invalid sort-order value", func() {
						req, w = util.CreateTestRequest("GET", JournalEntriesBaseURL+"?sort-order=invalid", nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusBadRequest))
					})
				})
			})

			Context("Pagination", func() {
				Context("Allowed Values", func() {
					It("should limit results with limit = 2", func() {
						req, w = util.CreateTestRequest("GET", JournalEntriesBaseURL+"?limit=2", nil)
						router.ServeHTTP(w, req)
						response := decodeEntryList(w, http.StatusOK)
						Expect(response.Records).To(HaveLen(2))
						Expect(response.Metadata.Total).To(Equal(int64(5)))
					})

					It("should skip entries with offset = 2, limit = 2", func() {
						req, w = util.CreateTestRequest("GET", JournalEntriesBaseURL+"?offset=2&limit=2", nil)
						router.ServeHTTP(w, req)
						response := decodeEntryList(w, http.StatusOK)
						Expect(response.Records).To(HaveLen(2))
						Expect(response.Metadata.Total).To(Equal(int64(5)))
					})

					It("should return last entry with offset = 4, limit = 2", func() {
						req, w = util.CreateTestRequest("GET", JournalEntriesBaseURL+"?offset=4&limit=2", nil)
						router.ServeHTTP(w, req)
						response := decodeEntryList(w, http.StatusOK)
						Expect(response.Records).To(HaveLen(1))
						Expect(response.Metadata.Total).To(Equal(int64(5)))
					})

					It("should return empty list for offset beyond total", func() {
						req, w = util.CreateTestRequest("GET", JournalEntriesBaseURL+"?offset=10", nil)
						router.ServeHTTP(w, req)
						response := decodeEntryList(w, http.StatusOK)
						Expect(response.Records).To(BeEmpty())
						Expect(response.Metadata.Total).To(Equal(int64(5)))
					})

					It("should accept limit = 1 (minimum)", func() {
						req, w = util.CreateTestRequest("GET", JournalEntriesBaseURL+"?limit=1", nil)
						router.ServeHTTP(w, req)
						response := decodeEntryList(w, http.StatusOK)
						Expect(response.Records).To(HaveLen(1))
					})

					It("should accept limit = 100 (maximum)", func() {
						req, w = util.CreateTestRequest("GET", JournalEntriesBaseURL+"?limit=100", nil)
						router.ServeHTTP(w, req)
						response := decodeEntryList(w, http.StatusOK)
						Expect(response.Records).To(HaveLen(5))
					})
				})

				Context("Bad Values", func() {
					It("should return 400 for limit exceeds maximum (101)", func() {
						req, w = util.CreateTestRequest("GET", JournalEntriesBaseURL+"?limit=101", nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusBadRequest))
					})

					It("should return 400 for limit = 0", func() {
						req, w = util.CreateTestRequest("GET", JournalEntriesBaseURL+"?limit=0", nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusBadRequest))
					})

					It("should return 400 for negative limit", func() {
						req, w = util.CreateTestRequest("GET", JournalEntriesBaseURL+"?limit=-1", nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusBadRequest))
					})

					It("should return 400 for negative offset", func() {
						req, w = util.CreateTestRequest("GET", JournalEntriesBaseURL+"?offset=-1", nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusBadRequest))
					})

					It("should return 400 for non-numeric limit", func() {
						req, w = util.CreateTestRequest("GET", JournalEntriesBaseURL+"?limit=abc", nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusBadRequest))
					})

					It("should return 400 for non-numeric offset", func() {
						req, w = util.CreateTestRequest("GET", JournalEntriesBaseURL+"?offset=xyz", nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusBadRequest))
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

					entryRepo := repository.NewJournalRepository(db)
					entryMgr = manager.NewJournalManager(entryRepo)
					journalHandler = handler.NewJournalHandler(entryMgr)

					router = util.CreateTestGinRouter()
					v1 := router.Group("/v1")
					handler.SetupJournalEntryRoutes(v1, journalHandler)

					req, w = util.CreateTestRequest("GET", JournalEntriesBaseURL, nil)
					router.ServeHTTP(w, req)
				})

				It("should return empty list", func() {
					response := decodeEntryList(w, http.StatusOK)
					Expect(response.Records).To(BeEmpty())
					Expect(response.Metadata.Total).To(Equal(int64(0)))
				})
			})
		})
	})
})
