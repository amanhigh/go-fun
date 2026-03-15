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

var _ = Describe("Barkat Integration Test", func() {
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
			base := util.NewHttpServer(config.HttpServerConfig{Name: "kohan", Port: testPort}, engine, shutdown)
			// Create custom lifecycle without monitor handler to avoid nil pointer
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
				time.Sleep(2 * time.Second) // wait for server to start
			}()

			Eventually(func() error {
				_, err := client.R().Get("/v1/journal-entries")
				return err
			}, 5*time.Second, 100*time.Millisecond).Should(Succeed())
		}
	})

	Context("Basic Smoke Test", func() {
		var createdEntry barkat.Journal

		BeforeEach(func() {
			entry := barkat.Journal{
				Ticker:   "TEST",
				Sequence: "MWD",
				Type:     "SET",
				Status:   "RUNNING",
				Images: []barkat.Image{
					{Timeframe: "DL", FileName: "daily_chart.png"},
					{Timeframe: "WK", FileName: "weekly_chart.png"},
					{Timeframe: "MN", FileName: "monthly_chart.png"},
					{Timeframe: "TMN", FileName: "trend_monthly_chart.png"},
				},
				Tags: []barkat.Tag{
					{Tag: "test", Type: "REASON"},
				},
			}
			resp, err := client.R().SetBody(entry).Post("/v1/journal-entries")
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode()).To(Equal(http.StatusCreated))

			var envelope common.Envelope[barkat.Journal]
			Expect(json.Unmarshal(resp.Body(), &envelope)).To(Succeed())
			Expect(envelope.Status).To(Equal(common.EnvelopeSuccess))
			createdEntry = envelope.Data
		})

		It("should create and retrieve entry", func() {
			Expect(createdEntry.ExternalID).ToNot(BeEmpty())
			Expect(createdEntry.Ticker).To(Equal("TEST"))
			Expect(createdEntry.Status).To(Equal("RUNNING"))

			// Verify we can retrieve the entry
			resp, err := client.R().Get("/v1/journal-entries/" + createdEntry.ExternalID)
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode()).To(Equal(200))

			var envelope common.Envelope[barkat.Journal]
			Expect(json.Unmarshal(resp.Body(), &envelope)).To(Succeed())
			Expect(envelope.Status).To(Equal(common.EnvelopeSuccess))
			fetchedEntry := envelope.Data
			Expect(fetchedEntry.ID).To(Equal(createdEntry.ID))
			Expect(fetchedEntry.Ticker).To(Equal("TEST"))
		})

		It("should list entries", func() {
			resp, err := client.R().Get("/v1/journal-entries?limit=10")
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode()).To(Equal(http.StatusOK))
			var envelope common.Envelope[barkat.JournalList]
			Expect(json.Unmarshal(resp.Body(), &envelope)).To(Succeed())
			Expect(envelope.Status).To(Equal(common.EnvelopeSuccess))
			Expect(envelope.Data.Journals).ToNot(BeEmpty())
		})
	})
})
