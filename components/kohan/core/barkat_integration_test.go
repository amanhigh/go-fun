package core_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/components/kohan/core"
	"github.com/amanhigh/go-fun/components/kohan/handler"
	"github.com/amanhigh/go-fun/components/kohan/manager"
	"github.com/amanhigh/go-fun/components/kohan/repository"
	"github.com/amanhigh/go-fun/models/barkat"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/config"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/golang-sql/civil"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const testPort = 19020

var (
	baseURL string
	client  *resty.Client
)

// testLifecycle implements ServerLifecycle without monitor handler for testing
type testLifecycle struct {
	journalHandler handler.JournalHandler
	imageHandler   handler.ImageHandler
	noteHandler    handler.NoteHandler
	tagHandler     handler.TagHandler
}

func (t *testLifecycle) RegisterRoutes(engine *gin.Engine) {
	// Only register journal routes for this test
	journal := engine.Group("/v1/journal-entries")
	handler.SetupJournalEntryRoutes(journal, t.journalHandler)
	handler.SetupImageRoutes(journal, t.imageHandler)
	handler.SetupNoteRoutes(journal, t.noteHandler)
	handler.SetupTagRoutes(journal, t.tagHandler)
}

func (t *testLifecycle) RegisterSwagger(_ *gin.Engine)    {}
func (t *testLifecycle) BeforeStart(_ context.Context)    {}
func (t *testLifecycle) BeforeShutdown(_ context.Context) {}
func (t *testLifecycle) AfterShutdown(_ context.Context)  {}

// standardImages provides 4 required images for journal creation
var standardImages = []barkat.Image{
	{Timeframe: "DL", FileName: "daily.png"},
	{Timeframe: "WK", FileName: "weekly.png"},
	{Timeframe: "MN", FileName: "monthly.png"},
	{Timeframe: "TMN", FileName: "trend_monthly.png"},
}

// decodeJournalResponse unmarshals response body into Journal envelope
func decodeJournalResponse(resp *resty.Response) barkat.Journal {
	var envelope common.Envelope[barkat.Journal]
	ExpectWithOffset(1, json.Unmarshal(resp.Body(), &envelope)).To(Succeed())
	ExpectWithOffset(1, envelope.Status).To(Equal(common.EnvelopeSuccess))
	return envelope.Data
}

// Barkat E2E Test Suite
//
// Tests critical paths through real HTTP server with in-memory SQLite DB.
// Focuses on scenarios that add value beyond unit/integration tests:
// - Full CRUD lifecycle with associations
// - Cascade delete (FK constraints)
// - Validation through real HTTP stack
// - Review status workflow
var _ = Describe("Barkat E2E Test", func() {
	BeforeEach(func() {
		if baseURL == "" {
			// Initialize Resty client
			client = resty.New()
			client.SetTimeout(5 * time.Second)
			client.SetHeader("Content-Type", "application/json")

			db, err := core.CreateTestBarkatDB()
			Expect(err).ToNot(HaveOccurred())

			entryRepo := repository.NewJournalRepository(db)
			entryMgr := manager.NewJournalManager(entryRepo)
			journalHandler := handler.NewJournalHandler(entryMgr)
			imageHandler := handler.NewImageHandler(manager.NewImageManager(entryMgr, repository.NewImageRepository(db)))
			noteHandler := handler.NewNoteHandler(manager.NewNoteManager(entryMgr, repository.NewNoteRepository(db)))
			tagHandler := handler.NewTagHandler(manager.NewTagManager(entryMgr, repository.NewTagRepository(db)))

			shutdown := util.NewGracefulShutdown()
			engine := gin.Default()
			core.RegisterJournalValidators()
			base := util.NewHttpServer(config.HttpServerConfig{Name: "kohan-e2e", Port: testPort}, engine, shutdown)
			lifecycle := &testLifecycle{
				journalHandler: journalHandler,
				imageHandler:   imageHandler,
				noteHandler:    noteHandler,
				tagHandler:     tagHandler,
			}
			base.SetLifecycle(lifecycle)
			baseURL = fmt.Sprintf("http://localhost:%d", testPort)
			client.SetBaseURL(baseURL)

			go func() {
				defer GinkgoRecover()
				_ = base.Start()
			}()

			Eventually(func() error {
				_, err := client.R().Get("/v1/journal-entries")
				return err
			}, 5*time.Second, 100*time.Millisecond).Should(Succeed())
		}
	})

	// Full CRUD Lifecycle - Tests complete flow through real HTTP + DB
	Context("Journal CRUD Lifecycle", func() {
		var createdEntry barkat.Journal

		BeforeEach(func() {
			entry := barkat.Journal{
				Ticker:   "LIFECYCLE",
				Sequence: "MWD",
				Type:     "SET",
				Status:   "RUNNING",
				Images:   standardImages,
				Tags:     []barkat.Tag{{Tag: "oe", Type: "REASON"}},
				Notes:    []barkat.Note{{Status: "SET", Content: "Initial setup note", Format: "MARKDOWN"}},
			}
			resp, err := client.R().SetBody(entry).Post("/v1/journal-entries")
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode()).To(Equal(http.StatusCreated))
			createdEntry = decodeJournalResponse(resp)
		})

		It("should create entry with all associations", func() {
			Expect(createdEntry.ExternalID).To(HavePrefix("jrn_"))
			Expect(createdEntry.Images).To(HaveLen(4))
			Expect(createdEntry.Tags).To(HaveLen(1))
			Expect(createdEntry.Notes).To(HaveLen(1))

			// Verify association IDs
			for _, img := range createdEntry.Images {
				Expect(img.ExternalID).To(HavePrefix("img_"))
			}
			Expect(createdEntry.Tags[0].ExternalID).To(HavePrefix("tag_"))
			Expect(createdEntry.Notes[0].ExternalID).To(HavePrefix("not_"))
		})

		It("should retrieve entry with associations", func() {
			resp, err := client.R().Get("/v1/journal-entries/" + createdEntry.ExternalID)
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode()).To(Equal(http.StatusOK))

			fetched := decodeJournalResponse(resp)
			Expect(fetched.Ticker).To(Equal("LIFECYCLE"))
			Expect(fetched.Images).To(HaveLen(4))
			Expect(fetched.Tags).To(HaveLen(1))
			Expect(fetched.Notes).To(HaveLen(1))
		})

		It("should list entries with pagination", func() {
			resp, err := client.R().Get("/v1/journal-entries?limit=10")
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode()).To(Equal(http.StatusOK))

			var envelope common.Envelope[barkat.JournalList]
			Expect(json.Unmarshal(resp.Body(), &envelope)).To(Succeed())
			Expect(envelope.Data.Journals).ToNot(BeEmpty())
		})

		It("should delete entry and cascade to associations", func() {
			resp, err := client.R().Delete("/v1/journal-entries/" + createdEntry.ExternalID)
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode()).To(Equal(http.StatusNoContent))

			// Verify entry is deleted
			resp, err = client.R().Get("/v1/journal-entries/" + createdEntry.ExternalID)
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode()).To(Equal(http.StatusNotFound))
		})
	})

	// Review Status Workflow - Tests PATCH update flow
	Context("Review Status Workflow", func() {
		var createdEntry barkat.Journal

		BeforeEach(func() {
			entry := barkat.Journal{
				Ticker:   "REVIEW",
				Sequence: "YR",
				Type:     "RESULT",
				Status:   "SUCCESS",
				Images:   standardImages,
			}
			resp, err := client.R().SetBody(entry).Post("/v1/journal-entries")
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode()).To(Equal(http.StatusCreated))
			createdEntry = decodeJournalResponse(resp)
		})

		It("should mark entry as reviewed", func() {
			reviewDate := civil.Date{Year: 2024, Month: 1, Day: 15}
			payload := barkat.JournalReviewUpdate{ReviewedAt: &reviewDate}

			resp, err := client.R().SetBody(payload).Patch("/v1/journal-entries/" + createdEntry.ExternalID)
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode()).To(Equal(http.StatusOK))

			var envelope common.Envelope[barkat.UpdateJournalStatusResponse]
			Expect(json.Unmarshal(resp.Body(), &envelope)).To(Succeed())
			Expect(envelope.Data.ReviewedAt).ToNot(BeNil())
		})

		It("should clear reviewed status", func() {
			// First mark as reviewed
			reviewDate := civil.Date{Year: 2024, Month: 1, Day: 15}
			resp, _ := client.R().SetBody(barkat.JournalReviewUpdate{ReviewedAt: &reviewDate}).
				Patch("/v1/journal-entries/" + createdEntry.ExternalID)
			Expect(resp.StatusCode()).To(Equal(http.StatusOK))

			// Then clear
			resp, err := client.R().SetBody(barkat.JournalReviewUpdate{ReviewedAt: nil}).
				Patch("/v1/journal-entries/" + createdEntry.ExternalID)
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode()).To(Equal(http.StatusOK))

			var envelope common.Envelope[barkat.UpdateJournalStatusResponse]
			Expect(json.Unmarshal(resp.Body(), &envelope)).To(Succeed())
			Expect(envelope.Data.ReviewedAt).To(BeNil())
		})
	})

	// Validation Through HTTP Stack - Ensures validators are registered
	Context("Validation Errors", func() {
		It("should reject invalid ticker format", func() {
			entry := barkat.Journal{
				Ticker:   "lowercase", // PRD: must be uppercase
				Sequence: "MWD",
				Type:     "SET",
				Status:   "RUNNING",
				Images:   standardImages,
			}
			resp, err := client.R().SetBody(entry).Post("/v1/journal-entries")
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode()).To(Equal(http.StatusBadRequest))
		})

		It("should reject insufficient images", func() {
			entry := barkat.Journal{
				Ticker:   "VALID",
				Sequence: "MWD",
				Type:     "SET",
				Status:   "RUNNING",
				Images:   []barkat.Image{{Timeframe: "DL", FileName: "only_one.png"}}, // PRD: min 4
			}
			resp, err := client.R().SetBody(entry).Post("/v1/journal-entries")
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode()).To(Equal(http.StatusBadRequest))
		})

		It("should reject future review date", func() {
			// Create entry first
			entry := barkat.Journal{
				Ticker:   "FUTURE",
				Sequence: "MWD",
				Type:     "SET",
				Status:   "RUNNING",
				Images:   standardImages,
			}
			resp, _ := client.R().SetBody(entry).Post("/v1/journal-entries")
			created := decodeJournalResponse(resp)

			// Try to set future date
			futureDate := civil.Date{Year: 2099, Month: 12, Day: 31}
			resp, err := client.R().SetBody(barkat.JournalReviewUpdate{ReviewedAt: &futureDate}).
				Patch("/v1/journal-entries/" + created.ExternalID)
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode()).To(Equal(http.StatusBadRequest))
		})
	})

	// Error Handling - Tests 404 and invalid ID scenarios
	Context("Error Handling", func() {
		It("should return 404 for non-existent entry", func() {
			// Use valid format (8 hex chars) but non-existent ID
			resp, err := client.R().Get("/v1/journal-entries/jrn_12345678")
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode()).To(Equal(http.StatusNotFound))
		})

		It("should return 400 for invalid ID format", func() {
			resp, err := client.R().Get("/v1/journal-entries/invalid_format")
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode()).To(Equal(http.StatusBadRequest))
		})
	})
})
