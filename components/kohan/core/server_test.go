package core_test

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/amanhigh/go-fun/common/util"
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

	Context("Graceful Shutdown Integration", func() {
		var (
			realShutdown          util.Shutdown
			serverDone            chan error
			freePort              int
			startAndWaitForServer func(port int)
		)

		BeforeEach(func() {
			// Use REAL util.Shutdown - no mocks!
			realShutdown = util.NewGracefulShutdown()
			serverDone = make(chan error, 1)

			// Get free port
			listener, err := net.Listen("tcp", ":0")
			Expect(err).ToNot(HaveOccurred())
			freePort = listener.Addr().(*net.TCPAddr).Port
			listener.Close()

			// Helper function to start server and wait for readiness
			startAndWaitForServer = func(port int) {
				go func() {
					err := server.StartWithShutdownHandler(port, realShutdown)
					serverDone <- err
				}()

				Eventually(func() error {
					conn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", port))
					if conn != nil {
						conn.Close()
					}
					return err
				}, 1*time.Second, 50*time.Millisecond).Should(Succeed())
			}
		})

		AfterEach(func() {
			// Ensure server shutdown completes (if not already consumed by test)
			select {
			case <-serverDone:
				// Already consumed by test
			case <-time.After(4 * time.Second):
				// Timeout - server should have shutdown by now
			}
		})

		It("should start and shutdown gracefully", func() {
			startAndWaitForServer(freePort)
			realShutdown.Stop(context.Background())
			Eventually(serverDone, 2*time.Second).Should(Receive(BeNil()))
		})

		It("should serve HTTP requests before shutdown", func() {
			startAndWaitForServer(freePort)

			// Make actual HTTP request to verify server is functional
			resp, err := http.Get(fmt.Sprintf("http://localhost:%d/v1/clip/", freePort))
			if resp != nil {
				defer resp.Body.Close()
			}
			Expect(err).To(BeNil())

			realShutdown.Stop(context.Background())
			Eventually(serverDone, 2*time.Second).Should(Receive(BeNil()))
		})

		It("should handle startup errors", func() {
			// Don't use helper - test error case directly
			go func() {
				err := server.StartWithShutdownHandler(-1, realShutdown)
				serverDone <- err
			}()

			Eventually(serverDone, 500*time.Millisecond).Should(Receive(HaveOccurred()))
		})
	})

})
