package core_test

import (
	"bytes"
	"errors"
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
		mockManager *mocks.AutoManagerInterface
		testPath    = "/tmp/test-capture"
	)

	BeforeEach(func() {
		gin.SetMode(gin.TestMode)
		mockManager = mocks.NewAutoManagerInterface(GinkgoT())
	})

	Context("Constructor", func() {
		// TODO: #B Remove this Test after writing extensive Base Server Test.
		It("should create server with handlers", func() {
			monitorHandler := handler.NewMonitorHandler(testPath, mockManager)
			db, err := core.CreateTestBarkatDB()
			Expect(err).ToNot(HaveOccurred())
			entryMgr := manager.NewJournalManager(repository.NewJournalRepository(db))
			journalHandler := handler.NewJournalHandler(entryMgr)
			// TODO: #C Wire Up Test using DI Framework Module Override as well.
			imageHandler := handler.NewImageHandler(manager.NewImageManager(entryMgr, repository.NewImageRepository(db)))
			noteHandler := handler.NewNoteHandler(manager.NewNoteManager(entryMgr, repository.NewNoteRepository(db)))
			tagHandler := handler.NewTagHandler(manager.NewTagManager(entryMgr, repository.NewTagRepository(db)))
			shutdown := util.NewGracefulShutdown()
			base := util.NewHttpServer(config.HttpServerConfig{Name: "kohan", Port: 0}, gin.Default(), shutdown)
			lifecycle := core.NewKohanServerLifecycle(monitorHandler, journalHandler, imageHandler, noteHandler, tagHandler)
			base.SetLifecycle(lifecycle)
			Expect(base).ToNot(BeNil())
		})
	})

	Context("MonitorHandler", func() {
		var (
			monitorHandler handler.MonitorHandler
			recorder       *httptest.ResponseRecorder
		)

		BeforeEach(func() {
			monitorHandler = handler.NewMonitorHandler(testPath, mockManager)
			recorder = httptest.NewRecorder()
		})

		Context("HandleRecordTicker", func() {
			Context("when recording succeeds", func() {
				BeforeEach(func() {
					mockManager.EXPECT().
						RecordTicker(mock.Anything, "AAPL", testPath).
						Return(nil)
				})

				It("should return success response", func() {
					req := httptest.NewRequest("GET", "/v1/ticker/AAPL/record", nil)
					c, _ := gin.CreateTestContext(recorder)
					c.Request = req
					c.Params = gin.Params{
						{Key: "ticker", Value: "AAPL"},
					}

					monitorHandler.HandleRecordTicker(c)

					Expect(recorder.Code).To(Equal(http.StatusOK))
					Expect(recorder.Body.String()).To(ContainSubstring("Success"))
				})
			})

			Context("when recording fails", func() {
				BeforeEach(func() {
					mockManager.EXPECT().
						RecordTicker(mock.Anything, "MSFT", testPath).
						Return(errors.New("screenshot failed"))
				})

				It("should return error response", func() {
					req := httptest.NewRequest("GET", "/v1/ticker/MSFT/record", nil)
					c, _ := gin.CreateTestContext(recorder)
					c.Request = req
					c.Params = gin.Params{
						{Key: "ticker", Value: "MSFT"},
					}

					monitorHandler.HandleRecordTicker(c)

					Expect(recorder.Code).To(Equal(http.StatusInternalServerError))
					Expect(recorder.Body.String()).To(ContainSubstring("screenshot failed"))
				})
			})

			Context("with different ticker values", func() {
				It("should pass correct ticker to manager", func() {
					testTicker := "GOOGL"
					mockManager.EXPECT().
						RecordTicker(mock.Anything, testTicker, testPath).
						Return(nil)

					req := httptest.NewRequest("GET", "/v1/ticker/"+testTicker+"/record", nil)
					c, _ := gin.CreateTestContext(recorder)
					c.Request = req
					c.Params = gin.Params{
						{Key: "ticker", Value: testTicker},
					}

					monitorHandler.HandleRecordTicker(c)

					Expect(recorder.Code).To(Equal(http.StatusOK))
				})
			})
		})
	})

	Context("JournalHandler", func() {
		var (
			journalHandler handler.JournalHandler
			recorder       *httptest.ResponseRecorder
		)

		BeforeEach(func() {
			db, err := core.CreateTestBarkatDB()
			Expect(err).ToNot(HaveOccurred())
			repo := repository.NewJournalRepository(db)
			mgr := manager.NewJournalManager(repo)
			journalHandler = handler.NewJournalHandler(mgr)
			recorder = httptest.NewRecorder()
		})

		Context("HandleCreateEntry", func() {
			It("should create entry and return 201", func() {
				body := `{"ticker":"RELIANCE","sequence":"MWD","type":"REJECTED","status":"FAIL","images":[{"timeframe":"DL"},{"timeframe":"WK"},{"timeframe":"MN"},{"timeframe":"TMN"}]}`
				req := httptest.NewRequest("POST", "/v1/journal-entries", bytes.NewBufferString(body))
				req.Header.Set("Content-Type", "application/json")
				c, _ := gin.CreateTestContext(recorder)
				c.Request = req

				journalHandler.HandleCreateEntry(c)

				Expect(recorder.Code).To(Equal(http.StatusCreated))
				Expect(recorder.Body.String()).To(ContainSubstring("RELIANCE"))
			})
		})
	})
})
