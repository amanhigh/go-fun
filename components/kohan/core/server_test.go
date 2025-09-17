package core_test

import (
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/amanhigh/go-fun/components/kohan/core"
	"github.com/amanhigh/go-fun/components/kohan/manager/mocks"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
)

var _ = Describe("MonitorServer", func() {
	var (
		server      *core.MonitorServer
		mockManager *mocks.AutoManagerInterface
		testPath    = "/tmp/test-capture"
	)

	BeforeEach(func() {
		gin.SetMode(gin.TestMode)
		mockManager = mocks.NewAutoManagerInterface(GinkgoT())
		server = core.NewMonitorServer(testPath, mockManager)
	})

	Context("Constructor and Configuration", func() {
		It("should create server with correct configuration", func() {
			Expect(server).ToNot(BeNil())
		})

		It("should accept nil manager for constructor", func() {
			nilServer := core.NewMonitorServer(testPath, nil)
			Expect(nilServer).ToNot(BeNil())
		})
	})

	Context("HandleRecordTicker", func() {
		var (
			recorder *httptest.ResponseRecorder
			req      *http.Request
		)

		BeforeEach(func() {
			recorder = httptest.NewRecorder()
		})

		Context("when recording succeeds", func() {
			BeforeEach(func() {
				mockManager.EXPECT().
					RecordTicker(mock.Anything, "AAPL", testPath).
					Return(nil)
				req = httptest.NewRequest("GET", "/v1/ticker/AAPL/record", nil)
			})

			It("should return success response", func() {
				c, _ := gin.CreateTestContext(recorder)
				c.Request = req
				c.Params = gin.Params{
					{Key: "ticker", Value: "AAPL"},
				}

				server.HandleRecordTicker(c)

				Expect(recorder.Code).To(Equal(http.StatusOK))
				Expect(recorder.Body.String()).To(ContainSubstring("Success"))
			})
		})

		Context("when recording fails", func() {
			BeforeEach(func() {
				mockManager.EXPECT().
					RecordTicker(mock.Anything, "MSFT", testPath).
					Return(errors.New("screenshot failed"))
				req = httptest.NewRequest("GET", "/v1/ticker/MSFT/record", nil)
			})

			It("should return error response", func() {
				c, _ := gin.CreateTestContext(recorder)
				c.Request = req
				c.Params = gin.Params{
					{Key: "ticker", Value: "MSFT"},
				}

				server.HandleRecordTicker(c)

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

				server.HandleRecordTicker(c)

				Expect(recorder.Code).To(Equal(http.StatusOK))
			})
		})
	})

})
