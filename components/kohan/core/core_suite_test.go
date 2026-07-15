package core_test

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/components/kohan/core"
	"github.com/amanhigh/go-fun/components/kohan/handler"
	handlerMocks "github.com/amanhigh/go-fun/components/kohan/handler/mocks"
	"github.com/amanhigh/go-fun/components/kohan/manager"
	"github.com/amanhigh/go-fun/components/kohan/manager/audit"
	"github.com/amanhigh/go-fun/components/kohan/repository"
	"github.com/amanhigh/go-fun/models/config"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
)

const testPort = 19020

var (
	server       util.HttpServer
	testImageDir string
)

func createTestImageDir() string {
	imageDir, err := os.MkdirTemp("", "kohan-e2e-images-*")
	Expect(err).ToNot(HaveOccurred())

	sampleImageDir := filepath.Join(imageDir, "2024", "01")
	Expect(os.MkdirAll(sampleImageDir, 0o755)).To(Succeed())
	Expect(os.WriteFile(filepath.Join(sampleImageDir, "sample.png"), []byte("sample-image"), 0o600)).To(Succeed())

	return imageDir
}

// configureMockOSHandler creates and configures a mock OS handler for safe testing
func configureMockOSHandler() *handlerMocks.OSHandler {
	osHandler := handlerMocks.NewOSHandler(GinkgoT())

	// Configure mock to return safe responses instead of executing real OS operations
	osHandler.EXPECT().HandleScreenshot(mock.Anything).Run(func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{
			"file_name":     "mock.png",
			"relative_path": "mock.png",
			"full_path":     "/mock/path/mock.png",
		}})
	}).Maybe()
	osHandler.EXPECT().HandleSubmapControl(mock.Anything).Run(func(ctx *gin.Context) {
		action := ctx.Param("action")
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "action": action})
	}).Maybe()
	osHandler.EXPECT().HandleRecordTicker(mock.Anything).Run(func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "Success")
	}).Maybe()
	osHandler.EXPECT().HandleReadClip(mock.Anything).Run(func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "mock_clipboard_content")
	}).Maybe()

	return osHandler
}

// verifyServerUp checks that the test server is accepting TCP connections.
func verifyServerUp() error {
	conn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", testPort))
	if conn != nil {
		conn.Close()
	}
	return err
}

func TestCore(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Core Suite")
}

var _ = BeforeSuite(func() {
	testImageDir = createTestImageDir()

	var err error
	db, err := core.CreateTestBarkatDB()
	Expect(err).ToNot(HaveOccurred())

	journalRepo := repository.NewJournalRepository(db)
	journalMgr := manager.NewJournalManager(journalRepo)
	journalHandler := handler.NewJournalHandler(journalMgr)
	imageHandler := handler.NewImageHandler(manager.NewImageManager(journalMgr, repository.NewImageRepository(db)))
	noteHandler := handler.NewNoteHandler(manager.NewNoteManager(journalMgr, repository.NewNoteRepository(db)))
	tagHandler := handler.NewTagHandler(manager.NewTagManager(journalMgr, repository.NewTagRepository(db)))
	tickerRepo := repository.NewTickerRepository(db)
	tickerMgr := manager.NewBarkatTickerManager(tickerRepo)
	tickerHandler := handler.NewTickerHandler(tickerMgr)
	alertTickerRepo := repository.NewAlertTickerRepository(db)
	alertTickerMgr := manager.NewAlertTickerManager(alertTickerRepo)
	alertTickerHandler := handler.NewAlertTickerHandler(alertTickerMgr)
	priceAlertRepo := repository.NewPriceAlertRepository(db)
	priceAlertMgr := manager.NewPriceAlertManager(priceAlertRepo)
	priceAlertHandler := handler.NewPriceAlertHandler(priceAlertMgr)
	auditRepo := repository.NewAuditRepository(db)
	auditRegistry := audit.NewPluginRegistry()
	Expect(auditRegistry.RegisterPlugin(audit.NewAlertCoveragePlugin(auditRepo))).ToNot(HaveOccurred())
	Expect(auditRegistry.RegisterPlugin(audit.NewStaleReviewPlugin(auditRepo))).ToNot(HaveOccurred())
	auditMgr := manager.NewAuditManager(auditRegistry)
	auditHandler := handler.NewAuditHandler(auditMgr)
	indexPortal := handler.NewIndexPortal()
	journalPortal := handler.NewJournalPortal(testImageDir)

	osHandler := configureMockOSHandler()

	lifecycle := core.ProvideKohanLifecycle(
		osHandler,
		journalHandler,
		imageHandler,
		noteHandler,
		tagHandler,
		tickerHandler,
		alertTickerHandler,
		priceAlertHandler,
		auditHandler,
		core.PortalHandlers{
			IndexPortal:   indexPortal,
			JournalPortal: journalPortal,
		},
	)

	shutdown := util.NewGracefulShutdown()
	engine := gin.Default()
	engine.UseRawPath = true         // treat %2F as single path segment for composite tickers like NIFTY/USDINR
	engine.UnescapePathValues = true // decode encoded path segments back to original ticker value
	core.RegisterJournalValidators()
	server = util.NewHttpServer(config.HttpServerConfig{Name: "kohan-e2e", Port: testPort}, engine, shutdown)
	server.SetLifecycle(lifecycle)

	go func() {
		defer GinkgoRecover()
		_ = server.Start()
	}()

	Eventually(verifyServerUp, 5*time.Second, 10*time.Millisecond).Should(Succeed())
})

var _ = AfterSuite(func() {
	if server != nil {
		server.Stop(context.Background())
	}
	if testImageDir != "" {
		Expect(os.RemoveAll(testImageDir)).To(Succeed())
	}
})
