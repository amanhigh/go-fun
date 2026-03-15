package core_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/components/kohan/core"
	"github.com/amanhigh/go-fun/components/kohan/handler"
	"github.com/amanhigh/go-fun/components/kohan/manager"
	"github.com/amanhigh/go-fun/components/kohan/manager/mocks"
	"github.com/amanhigh/go-fun/components/kohan/repository"
	"github.com/amanhigh/go-fun/models/config"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
)

var _ = Describe("KohanServer", func() {
	var (
		mockManager    *mocks.AutoManagerInterface
		testPath       = "/tmp/test-capture"
		monitorHandler handler.MonitorHandler
		journalHandler handler.JournalHandler
		imageHandler   handler.ImageHandler
		noteHandler    handler.NoteHandler
		tagHandler     handler.TagHandler
		server         *util.HttpServerImpl
		shutdown       util.Shutdown
	)

	BeforeEach(func() {
		gin.SetMode(gin.TestMode)
		mockManager = mocks.NewAutoManagerInterface(GinkgoT())

		// Setup handlers
		monitorHandler = handler.NewMonitorHandler(testPath, mockManager)
		db, err := core.CreateTestBarkatDB()
		Expect(err).ToNot(HaveOccurred())
		journalMgr := manager.NewJournalManager(repository.NewJournalRepository(db))
		journalHandler = handler.NewJournalHandler(journalMgr)
		imageHandler = handler.NewImageHandler(manager.NewImageManager(journalMgr, repository.NewImageRepository(db)))
		noteHandler = handler.NewNoteHandler(manager.NewNoteManager(journalMgr, repository.NewNoteRepository(db)))
		tagHandler = handler.NewTagHandler(manager.NewTagManager(journalMgr, repository.NewTagRepository(db)))

		// Setup server
		shutdown = util.NewGracefulShutdown()
		server = util.NewHttpServer(config.HttpServerConfig{Name: "kohan", Port: 0}, gin.Default(), shutdown)
		lifecycle := core.NewKohanServerLifecycle(monitorHandler, journalHandler, imageHandler, noteHandler, tagHandler)
		server.SetLifecycle(lifecycle)
	})

	Context("Kohan Integration Smoke Test", func() {
		It("should create server with all Kohan handlers", func() {
			Expect(server).ToNot(BeNil())
			Expect(server.Name).To(Equal("kohan"))
		})
	})

	Context("MonitorHandler Integration", func() {
		var recorder *httptest.ResponseRecorder

		BeforeEach(func() {
			mockManager.EXPECT().
				RecordTicker(mock.Anything, "AAPL", testPath).
				Return(nil)
			recorder = httptest.NewRecorder()
		})

		It("should handle ticker recording endpoint", func() {
			req := httptest.NewRequest("GET", "/v1/ticker/AAPL/record", nil)
			c, _ := gin.CreateTestContext(recorder)
			c.Request = req
			c.Params = gin.Params{{Key: "ticker", Value: "AAPL"}}

			monitorHandler.HandleRecordTicker(c)

			Expect(recorder.Code).To(Equal(http.StatusOK))
			Expect(recorder.Body.String()).To(ContainSubstring("Success"))
		})
	})

	Context("JournalHandler Integration", func() {
		var recorder *httptest.ResponseRecorder

		BeforeEach(func() {
			recorder = httptest.NewRecorder()
		})

		It("should handle journal creation", func() {
			body := `{"ticker":"RELIANCE","sequence":"MWD","type":"REJECTED","status":"FAIL","images":[` +
				`{"timeframe":"DL","file_name":"RELIANCE.mwd.test.png"},` +
				`{"timeframe":"WK","file_name":"RELIANCE.mwd.test.png"},` +
				`{"timeframe":"MN","file_name":"RELIANCE.mwd.test.png"},` +
				`{"timeframe":"TMN","file_name":"RELIANCE.mwd.test.png"}]}`
			req := httptest.NewRequest("POST", "/v1/journals", bytes.NewBufferString(body))
			req.Header.Set("Content-Type", "application/json")

			// Use the server's engine which has validators registered
			c, _ := gin.CreateTestContext(recorder)
			c.Request = req
			// Register validators for this test context
			core.RegisterJournalValidators()

			journalHandler.HandleCreateJournal(c)

			Expect(recorder.Code).To(Equal(http.StatusCreated))
			Expect(recorder.Body.String()).To(ContainSubstring("RELIANCE"))
		})
	})
})
