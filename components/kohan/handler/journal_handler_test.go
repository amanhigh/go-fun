//nolint:dupl
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

// JournalHandler Integration Tests - Tests behavior with real SQLite DB, managers, and repositories
// This tests the complete HTTP → Handler → Manager → Repository → Database flow
var _ = Describe("JournalHandler Integration", func() {
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
		// Create real SQLite database for testing using proper migrations
		db, err = core.CreateTestBarkatDB()
		Expect(err).ToNot(HaveOccurred())

		// Create real managers and repositories (no mocks)
		entryRepo := repository.NewJournalRepository(db)
		entryMgr = manager.NewJournalManager(entryRepo)
		journalHandler = handler.NewJournalHandler(entryMgr)

		// Setup Gin router using helper
		router = util.CreateTestGinRouter()
		v1 := router.Group("/v1")
		handler.SetupJournalEntryRoutes(v1, journalHandler)
	})

	AfterEach(func() {
		sqlDB, err := db.DB()
		Expect(err).ToNot(HaveOccurred())
		sqlDB.Close()
	})

	Context("HandleCreateEntry", func() {
		Context("with valid request", func() {
			BeforeEach(func() {
				// Setup HTTP request with real entry data
				entry := barkat.Entry{
					Ticker:   "GRSE",
					Sequence: "mwd",
					Type:     "rejected",
					Status:   "fail",
				}
				req, w = util.CreateTestRequest("POST", "/v1/journal-entries", entry)
			})

			It("should create entry and return 201", func() {
				router.ServeHTTP(w, req)

				var response barkat.Entry
				util.AssertJSONAndStatus(w, http.StatusCreated, &response)

				// Verify the created entry has proper data
				Expect(response.Ticker).To(Equal("GRSE"))
				Expect(response.Sequence).To(Equal("mwd"))
				Expect(response.Type).To(Equal("rejected"))
				Expect(response.Status).To(Equal("fail"))
				Expect(response.ID).ToNot(BeEmpty())
				Expect(response.CreatedAt).ToNot(BeZero())

				// Verify entry is actually in database
				entry, err := entryMgr.GetEntry(testCtx, response.ID)
				Expect(err).ToNot(HaveOccurred())
				Expect(entry.ID).To(Equal(response.ID))
				Expect(entry.Ticker).To(Equal("GRSE"))
			})
		})

		Context("with valid field variations", func() {
			It("should accept mwd sequence", func() {
				entry := barkat.Entry{
					Ticker:   "PDSL",
					Sequence: "mwd",
					Type:     "set",
					Status:   "taken",
				}
				req, w = util.CreateTestRequest("POST", "/v1/journal-entries", entry)
				router.ServeHTTP(w, req)
				var response barkat.Entry
				util.AssertJSONAndStatus(w, http.StatusCreated, &response)
				Expect(response.Sequence).To(Equal("mwd"))
			})

			It("should accept yr sequence", func() {
				entry := barkat.Entry{
					Ticker:   "SNF",
					Sequence: "yr",
					Type:     "result",
					Status:   "success",
				}
				req, w = util.CreateTestRequest("POST", "/v1/journal-entries", entry)
				router.ServeHTTP(w, req)
				var response barkat.Entry
				util.AssertJSONAndStatus(w, http.StatusCreated, &response)
				Expect(response.Sequence).To(Equal("yr"))
			})

			It("should accept rejected type", func() {
				entry := barkat.Entry{
					Ticker:   "TCS",
					Sequence: "mwd",
					Type:     "rejected",
					Status:   "fail",
				}
				req, w = util.CreateTestRequest("POST", "/v1/journal-entries", entry)
				router.ServeHTTP(w, req)
				var response barkat.Entry
				util.AssertJSONAndStatus(w, http.StatusCreated, &response)
				Expect(response.Type).To(Equal("rejected"))
			})

			It("should accept result type", func() {
				entry := barkat.Entry{
					Ticker:   "INFY",
					Sequence: "yr",
					Type:     "result",
					Status:   "success",
				}
				req, w = util.CreateTestRequest("POST", "/v1/journal-entries", entry)
				router.ServeHTTP(w, req)
				var response barkat.Entry
				util.AssertJSONAndStatus(w, http.StatusCreated, &response)
				Expect(response.Type).To(Equal("result"))
			})

			It("should accept set type", func() {
				entry := barkat.Entry{
					Ticker:   "RELIANCE",
					Sequence: "mwd",
					Type:     "set",
					Status:   "running",
				}
				req, w = util.CreateTestRequest("POST", "/v1/journal-entries", entry)
				router.ServeHTTP(w, req)
				var response barkat.Entry
				util.AssertJSONAndStatus(w, http.StatusCreated, &response)
				Expect(response.Type).To(Equal("set"))
			})

			It("should accept all valid status values", func() {
				statuses := []string{"set", "running", "dropped", "taken", "rejected", "success", "fail", "missed", "just_loss", "broken"}
				for _, status := range statuses {
					entry := barkat.Entry{
						Ticker:   "TEST",
						Sequence: "mwd",
						Type:     "rejected",
						Status:   status,
					}
					req, w = util.CreateTestRequest("POST", "/v1/journal-entries", entry)
					router.ServeHTTP(w, req)
					var response barkat.Entry
					util.AssertJSONAndStatus(w, http.StatusCreated, &response)
					Expect(response.Status).To(Equal(status))
				}
			})
		})

		Context("field validation", func() {
			It("should reject empty ticker", func() {
				entry := barkat.Entry{
					Ticker:   "",
					Sequence: "mwd",
					Type:     "rejected",
					Status:   "fail",
				}
				req, w = util.CreateTestRequest("POST", "/v1/journal-entries", entry)
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusBadRequest))
				var errorResponse map[string]interface{}
				util.AssertJSONAndStatus(w, http.StatusBadRequest, &errorResponse)
				Expect(errorResponse["error"]).To(ContainSubstring("required"))
			})

			It("should reject empty sequence", func() {
				entry := barkat.Entry{
					Ticker:   "GRSE",
					Sequence: "",
					Type:     "rejected",
					Status:   "fail",
				}
				req, w = util.CreateTestRequest("POST", "/v1/journal-entries", entry)
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusBadRequest))
				var errorResponse map[string]interface{}
				util.AssertJSONAndStatus(w, http.StatusBadRequest, &errorResponse)
				Expect(errorResponse["error"]).To(ContainSubstring("required"))
			})

			It("should reject empty type", func() {
				entry := barkat.Entry{
					Ticker:   "GRSE",
					Sequence: "mwd",
					Type:     "",
					Status:   "fail",
				}
				req, w = util.CreateTestRequest("POST", "/v1/journal-entries", entry)
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusBadRequest))
				var errorResponse map[string]interface{}
				util.AssertJSONAndStatus(w, http.StatusBadRequest, &errorResponse)
				Expect(errorResponse["error"]).To(ContainSubstring("required"))
			})

			It("should reject empty status", func() {
				entry := barkat.Entry{
					Ticker:   "GRSE",
					Sequence: "mwd",
					Type:     "rejected",
					Status:   "",
				}
				req, w = util.CreateTestRequest("POST", "/v1/journal-entries", entry)
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusBadRequest))
				var errorResponse map[string]interface{}
				util.AssertJSONAndStatus(w, http.StatusBadRequest, &errorResponse)
				Expect(errorResponse["error"]).To(ContainSubstring("required"))
			})

			It("should reject invalid sequence", func() {
				entry := barkat.Entry{
					Ticker:   "GRSE",
					Sequence: "invalid",
					Type:     "rejected",
					Status:   "fail",
				}
				req, w = util.CreateTestRequest("POST", "/v1/journal-entries", entry)
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusBadRequest))
				var errorResponse map[string]interface{}
				util.AssertJSONAndStatus(w, http.StatusBadRequest, &errorResponse)
				Expect(errorResponse["error"]).To(ContainSubstring("oneof"))
			})

			It("should reject invalid type", func() {
				entry := barkat.Entry{
					Ticker:   "GRSE",
					Sequence: "mwd",
					Type:     "invalid",
					Status:   "fail",
				}
				req, w = util.CreateTestRequest("POST", "/v1/journal-entries", entry)
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusBadRequest))
				var errorResponse map[string]interface{}
				util.AssertJSONAndStatus(w, http.StatusBadRequest, &errorResponse)
				Expect(errorResponse["error"]).To(ContainSubstring("oneof"))
			})

			It("should reject invalid status", func() {
				entry := barkat.Entry{
					Ticker:   "GRSE",
					Sequence: "mwd",
					Type:     "rejected",
					Status:   "invalid",
				}
				req, w = util.CreateTestRequest("POST", "/v1/journal-entries", entry)
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusBadRequest))
				var errorResponse map[string]interface{}
				util.AssertJSONAndStatus(w, http.StatusBadRequest, &errorResponse)
				Expect(errorResponse["error"]).To(ContainSubstring("oneof"))
			})

			It("should reject ticker too long", func() {
				entry := barkat.Entry{
					Ticker:   "VERYLONGTICKERNAMETHATEXCEEDS30CHARS",
					Sequence: "mwd",
					Type:     "rejected",
					Status:   "fail",
				}
				req, w = util.CreateTestRequest("POST", "/v1/journal-entries", entry)
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusBadRequest))
				var errorResponse map[string]interface{}
				util.AssertJSONAndStatus(w, http.StatusBadRequest, &errorResponse)
				Expect(errorResponse["error"]).To(ContainSubstring("max"))
			})

			It("should reject ticker too short", func() {
				entry := barkat.Entry{
					Ticker:   "",
					Sequence: "mwd",
					Type:     "rejected",
					Status:   "fail",
				}
				req, w = util.CreateTestRequest("POST", "/v1/journal-entries", entry)
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusBadRequest))
				var errorResponse map[string]interface{}
				util.AssertJSONAndStatus(w, http.StatusBadRequest, &errorResponse)
				Expect(errorResponse["error"]).To(ContainSubstring("min"))
			})
		})

		Context("with invalid JSON", func() {
			BeforeEach(func() {
				req, w = util.CreateTestRequest("POST", "/v1/journal-entries", []byte("invalid json"))
			})

			It("should return 400 error", func() {
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})
		})

		Context("edge cases", func() {
			It("should accept ticker at max length", func() {
				entry := barkat.Entry{
					Ticker:   "123456789012345678901234567890", // exactly 30 chars
					Sequence: "mwd",
					Type:     "rejected",
					Status:   "fail",
				}
				req, w = util.CreateTestRequest("POST", "/v1/journal-entries", entry)
				router.ServeHTTP(w, req)
				var response barkat.Entry
				util.AssertJSONAndStatus(w, http.StatusCreated, &response)
				Expect(response.Ticker).To(HaveLen(30))
			})

			It("should accept ticker with numbers", func() {
				entry := barkat.Entry{
					Ticker:   "GRSE123",
					Sequence: "mwd",
					Type:     "rejected",
					Status:   "fail",
				}
				req, w = util.CreateTestRequest("POST", "/v1/journal-entries", entry)
				router.ServeHTTP(w, req)
				var response barkat.Entry
				util.AssertJSONAndStatus(w, http.StatusCreated, &response)
				Expect(response.Ticker).To(Equal("GRSE123"))
			})

			It("should accept ticker with special chars", func() {
				entry := barkat.Entry{
					Ticker:   "GRSE-NSE",
					Sequence: "mwd",
					Type:     "rejected",
					Status:   "fail",
				}
				req, w = util.CreateTestRequest("POST", "/v1/journal-entries", entry)
				router.ServeHTTP(w, req)
				var response barkat.Entry
				util.AssertJSONAndStatus(w, http.StatusCreated, &response)
				Expect(response.Ticker).To(Equal("GRSE-NSE"))
			})
		})
	})

	Context("HandleGetEntry", func() {
		var createdEntry barkat.Entry

		BeforeEach(func() {
			// Create an entry to retrieve
			entry := barkat.Entry{
				Ticker:   "GRSE",
				Sequence: "mwd",
				Type:     "rejected",
				Status:   "fail",
			}
			Expect(entryMgr.CreateEntry(testCtx, &entry)).To(Succeed())
			createdEntry = entry
		})

		Context("with valid entry ID", func() {
			BeforeEach(func() {
				req, w = util.CreateTestRequest("GET", "/v1/journal-entries/"+createdEntry.ID, nil)
			})

			It("should return entry and 200", func() {
				router.ServeHTTP(w, req)

				var response barkat.Entry
				util.AssertJSONAndStatus(w, http.StatusOK, &response)

				// Verify all fields are present and correct
				Expect(response.ID).To(Equal(createdEntry.ID))
				Expect(response.Ticker).To(Equal("GRSE"))
				Expect(response.Sequence).To(Equal("mwd"))
				Expect(response.Type).To(Equal("rejected"))
				Expect(response.Status).To(Equal("fail"))
				Expect(response.CreatedAt).ToNot(BeZero())
			})
		})

		Context("with non-existent entry", func() {
			BeforeEach(func() {
				req, w = util.CreateTestRequest("GET", "/v1/journal-entries/nonexistent", nil)
			})

			It("should return 404 for non-existent entry", func() {
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusNotFound))
			})
		})

		Context("with malformed entry ID", func() {
			BeforeEach(func() {
				req, w = util.CreateTestRequest("GET", "/v1/journal-entries/invalid-id", nil)
			})

			It("should return 404 for malformed entry ID", func() {
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusNotFound))
			})
		})

		Context("with empty entry ID", func() {
			BeforeEach(func() {
				req, w = util.CreateTestRequest("GET", "/v1/journal-entries/", nil)
			})

			It("should return 400 for empty entry ID (route not found)", func() {
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})
		})
	})

	Context("HandleListEntries", func() {
		var createdEntries []barkat.Entry

		BeforeEach(func() {
			// Create multiple entries for testing
			entries := []barkat.Entry{
				{Ticker: "GRSE", Sequence: "mwd", Type: "rejected", Status: "fail"},
				{Ticker: "PDSL", Sequence: "yr", Type: "set", Status: "taken"},
				{Ticker: "SNF", Sequence: "mwd", Type: "result", Status: "success"},
			}
			for _, entry := range entries {
				Expect(entryMgr.CreateEntry(testCtx, &entry)).To(Succeed())
				createdEntries = append(createdEntries, entry)
			}
		})

		Context("with no filters", func() {
			BeforeEach(func() {
				req, w = util.CreateTestRequest("GET", "/v1/journal-entries", nil)
			})

			It("should list entries and return 200", func() {
				router.ServeHTTP(w, req)

				var response barkat.EntryList
				util.AssertJSONAndStatus(w, http.StatusOK, &response)

				// Verify entries are returned
				Expect(response.Records).To(HaveLen(3))
				Expect(response.Metadata.Total).To(Equal(int64(3)))

				// Verify entries have proper metadata
				for _, entry := range response.Records {
					Expect(entry.ID).ToNot(BeEmpty())
					Expect(entry.Ticker).ToNot(BeEmpty())
					Expect(entry.Sequence).ToNot(BeEmpty())
					Expect(entry.Type).ToNot(BeEmpty())
					Expect(entry.Status).ToNot(BeEmpty())
					Expect(entry.CreatedAt).ToNot(BeZero())
				}
			})

			It("should return entries in reverse chronological order", func() {
				router.ServeHTTP(w, req)

				var response barkat.EntryList
				util.AssertJSONAndStatus(w, http.StatusOK, &response)

				// Verify entries are sorted by created_at desc (newest first)
				for i := 1; i < len(response.Records); i++ {
					prevTime := response.Records[i-1].CreatedAt
					currTime := response.Records[i].CreatedAt
					Expect(prevTime).To(BeTemporally(">", currTime))
				}
			})
		})

		Context("with ticker filter", func() {
			BeforeEach(func() {
				req, w = util.CreateTestRequest("GET", "/v1/journal-entries?ticker=GRSE", nil)
			})

			It("should filter by ticker", func() {
				router.ServeHTTP(w, req)

				var response barkat.EntryList
				util.AssertJSONAndStatus(w, http.StatusOK, &response)

				Expect(response.Records).To(HaveLen(1))
				Expect(response.Records[0].Ticker).To(Equal("GRSE"))
				Expect(response.Metadata.Total).To(Equal(int64(1)))
			})
		})

		Context("with type filter", func() {
			BeforeEach(func() {
				req, w = util.CreateTestRequest("GET", "/v1/journal-entries?type=rejected", nil)
			})

			It("should filter by type", func() {
				router.ServeHTTP(w, req)

				var response barkat.EntryList
				util.AssertJSONAndStatus(w, http.StatusOK, &response)

				Expect(response.Records).To(HaveLen(1))
				Expect(response.Records[0].Type).To(Equal("rejected"))
				Expect(response.Metadata.Total).To(Equal(int64(1)))
			})
		})

		Context("with status filter", func() {
			BeforeEach(func() {
				req, w = util.CreateTestRequest("GET", "/v1/journal-entries?status=success", nil)
			})

			It("should filter by status", func() {
				router.ServeHTTP(w, req)

				var response barkat.EntryList
				util.AssertJSONAndStatus(w, http.StatusOK, &response)

				Expect(response.Records).To(HaveLen(1))
				Expect(response.Records[0].Status).To(Equal("success"))
				Expect(response.Metadata.Total).To(Equal(int64(1)))
			})
		})

		Context("with sequence filter", func() {
			BeforeEach(func() {
				req, w = util.CreateTestRequest("GET", "/v1/journal-entries?sequence=mwd", nil)
			})

			It("should filter by sequence", func() {
				router.ServeHTTP(w, req)

				var response barkat.EntryList
				util.AssertJSONAndStatus(w, http.StatusOK, &response)

				Expect(response.Records).To(HaveLen(2))
				Expect(response.Metadata.Total).To(Equal(int64(2)))
				for _, entry := range response.Records {
					Expect(entry.Sequence).To(Equal("mwd"))
				}
			})
		})

		Context("with pagination", func() {
			BeforeEach(func() {
				req, w = util.CreateTestRequest("GET", "/v1/journal-entries?limit=2&offset=1", nil)
			})

			It("should paginate results", func() {
				router.ServeHTTP(w, req)

				var response barkat.EntryList
				util.AssertJSONAndStatus(w, http.StatusOK, &response)

				Expect(response.Records).To(HaveLen(2))
				Expect(response.Metadata.Total).To(Equal(int64(3)))
			})
		})

		Context("with sorting", func() {
			BeforeEach(func() {
				req, w = util.CreateTestRequest("GET", "/v1/journal-entries?sort-by=ticker&sort-order=asc", nil)
			})

			It("should sort by ticker ascending", func() {
				router.ServeHTTP(w, req)

				var response barkat.EntryList
				util.AssertJSONAndStatus(w, http.StatusOK, &response)

				Expect(response.Records).To(HaveLen(3))
				Expect(response.Records[0].Ticker).To(Equal("GRSE"))
				Expect(response.Records[1].Ticker).To(Equal("PDSL"))
				Expect(response.Records[2].Ticker).To(Equal("SNF"))
			})
		})

		Context("with invalid query parameters", func() {
			It("should reject invalid ticker length", func() {
				req, w = util.CreateTestRequest("GET", "/v1/journal-entries?ticker=verylongtickernamethatexceedslimit", nil)
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})

			It("should reject invalid type", func() {
				req, w = util.CreateTestRequest("GET", "/v1/journal-entries?type=invalid", nil)
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})

			It("should reject invalid status", func() {
				req, w = util.CreateTestRequest("GET", "/v1/journal-entries?status=invalid", nil)
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})

			It("should reject invalid sequence", func() {
				req, w = util.CreateTestRequest("GET", "/v1/journal-entries?sequence=invalid", nil)
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})

			It("should reject invalid sort field", func() {
				req, w = util.CreateTestRequest("GET", "/v1/journal-entries?sort-by=invalid", nil)
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})

			It("should reject invalid sort order", func() {
				req, w = util.CreateTestRequest("GET", "/v1/journal-entries?sort-order=invalid", nil)
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})
		})

		Context("with no entries", func() {
			BeforeEach(func() {
				// Delete all entries
				for _, _ = range createdEntries {
					// Note: We don't have a delete method for entries, so we'll create a fresh DB
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
				}
				req, w = util.CreateTestRequest("GET", "/v1/journal-entries", nil)
			})

			It("should return empty list for no entries", func() {
				router.ServeHTTP(w, req)

				var response barkat.EntryList
				util.AssertJSONAndStatus(w, http.StatusOK, &response)
				Expect(response.Records).To(HaveLen(0))
				Expect(response.Metadata.Total).To(Equal(int64(0)))
			})
		})
	})
})
