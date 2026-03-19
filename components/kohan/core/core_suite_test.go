package core_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/components/kohan/core"
	"github.com/amanhigh/go-fun/components/kohan/handler"
	handlerMocks "github.com/amanhigh/go-fun/components/kohan/handler/mocks"
	"github.com/amanhigh/go-fun/components/kohan/manager"
	"github.com/amanhigh/go-fun/components/kohan/repository"
	"github.com/amanhigh/go-fun/models/config"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
)

const testPort = 19020

var (
	server util.HttpServer
)

// configureMockOSHandler creates and configures a mock OS handler for safe testing
func configureMockOSHandler() *handlerMocks.OSHandler {
	osHandler := handlerMocks.NewOSHandler(GinkgoT())

	// Configure mock to return safe responses instead of executing real OS operations
	osHandler.On("HandleSubmapControl", mock.Anything).Run(func(args mock.Arguments) {
		ctx, ok := args.Get(0).(*gin.Context)
		if !ok {
			return
		}
		action := ctx.Param("action")
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "action": action})
	})
	osHandler.On("HandleRecordTicker", mock.Anything).Run(func(args mock.Arguments) {
		ctx, ok := args.Get(0).(*gin.Context)
		if !ok {
			return
		}
		ctx.JSON(http.StatusOK, "Success")
	})
	osHandler.On("HandleReadClip", mock.Anything).Run(func(args mock.Arguments) {
		ctx, ok := args.Get(0).(*gin.Context)
		if !ok {
			return
		}
		ctx.JSON(http.StatusOK, "mock_clipboard_content")
	}).Maybe()

	return osHandler
}

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

	// Create mock OS handler for testing (safe - no real OS operations)
	osHandler := configureMockOSHandler()

	// Use provider function to create lifecycle (in sync with production)
	lifecycle := core.ProvideKohanLifecycle(
		osHandler,
		journalHandler,
		imageHandler,
		noteHandler,
		tagHandler,
	)

	shutdown := util.NewGracefulShutdown()
	engine := gin.Default()
	core.RegisterJournalValidators()
	server = util.NewHttpServer(config.HttpServerConfig{Name: "kohan-e2e", Port: testPort}, engine, shutdown)
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
