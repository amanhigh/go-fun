package core_test

import (
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/components/kohan/core"
	"github.com/amanhigh/go-fun/components/kohan/handler"
	"github.com/amanhigh/go-fun/components/kohan/manager/mocks"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
)

var _ = Describe("KohanServer", func() {
	var (
		server      *core.KohanServer
		mockManager *mocks.AutoManagerInterface
		testPath    = "/tmp/test-capture"
	)

	BeforeEach(func() {
		gin.SetMode(gin.TestMode)
		mockManager = mocks.NewAutoManagerInterface(GinkgoT())
	})

	Context("Constructor", func() {
		It("should create server with handlers", func() {
			monitorHandler := handler.NewMonitorHandler(testPath, mockManager)
			server = core.NewKohanServer(0, monitorHandler, nil, util.NewGracefulShutdown())
			Expect(server).ToNot(BeNil())
		})

		It("should accept nil handlers", func() {
			// FIXME: Remove this test.
			server = core.NewKohanServer(0, nil, nil, util.NewGracefulShutdown())
			Expect(server).ToNot(BeNil())
		})
	})

	// FIXME: Add Content for Journal Handler Tests. (Only one test case is enough)
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
})
