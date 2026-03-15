package core_test

import (
	"context"
	"testing"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/components/kohan/core"
	"github.com/amanhigh/go-fun/components/kohan/handler"
	"github.com/amanhigh/go-fun/components/kohan/manager"
	"github.com/amanhigh/go-fun/components/kohan/repository"
	"github.com/amanhigh/go-fun/models/barkat"
	"github.com/amanhigh/go-fun/models/config"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const testPort = 19020

var (
	server util.HttpServer
)

// testLifecycle implements ServerLifecycle without monitor handler for testing
type testLifecycle struct {
	journalHandler handler.JournalHandler
	imageHandler   handler.ImageHandler
	noteHandler    handler.NoteHandler
	tagHandler     handler.TagHandler
}

func (t *testLifecycle) RegisterRoutes(engine *gin.Engine) {
	// Using same path as real server to avoid bugs
	journal := engine.Group(barkat.JournalBase)
	handler.SetupJournalRoutes(journal, t.journalHandler)
	handler.SetupImageRoutes(journal, t.imageHandler)
	handler.SetupNoteRoutes(journal, t.noteHandler)
	handler.SetupTagRoutes(journal, t.tagHandler)
}

func (t *testLifecycle) RegisterSwagger(_ *gin.Engine)    {}
func (t *testLifecycle) BeforeStart(_ context.Context)    {}
func (t *testLifecycle) BeforeShutdown(_ context.Context) {}
func (t *testLifecycle) AfterShutdown(_ context.Context)  {}

func TestCore(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Core Suite")
}

var _ = BeforeSuite(func() {
	db, err := core.CreateTestBarkatDB()
	Expect(err).ToNot(HaveOccurred())

	journalRepo := repository.NewJournalRepository(db)
	journalMgr := manager.NewJournalManager(journalRepo)
	journalHandler := handler.NewJournalHandler(journalMgr)
	imageHandler := handler.NewImageHandler(manager.NewImageManager(journalMgr, repository.NewImageRepository(db)))
	noteHandler := handler.NewNoteHandler(manager.NewNoteManager(journalMgr, repository.NewNoteRepository(db)))
	tagHandler := handler.NewTagHandler(manager.NewTagManager(journalMgr, repository.NewTagRepository(db)))

	shutdown := util.NewGracefulShutdown()
	engine := gin.Default()
	core.RegisterJournalValidators()
	server = util.NewHttpServer(config.HttpServerConfig{Name: "kohan-e2e", Port: testPort}, engine, shutdown)
	lifecycle := &testLifecycle{
		journalHandler: journalHandler,
		imageHandler:   imageHandler,
		noteHandler:    noteHandler,
		tagHandler:     tagHandler,
	}
	server.SetLifecycle(lifecycle)

	go func() {
		defer GinkgoRecover()
		_ = server.Start()
	}()
})

var _ = AfterSuite(func() {
	if server != nil {
		server.Stop(context.Background())
	}
})
