package util_test

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("BaseHTTPServer", func() {
	var (
		server   *util.BaseHTTPServer
		shutdown util.Shutdown
	)

	BeforeEach(func() {
		gin.SetMode(gin.TestMode)
		shutdown = util.NewGracefulShutdown()
		server = util.NewBaseHTTPServer("test", 0, shutdown)
	})

	Context("Constructor", func() {
		It("should create server with correct configuration", func() {
			Expect(server).ToNot(BeNil())
			Expect(server.Name).To(Equal("test"))
			Expect(server.Engine).ToNot(BeNil())
			Expect(server.Server).ToNot(BeNil())
		})
	})

	Context("Lifecycle", func() {
		var (
			serverDone            chan error
			freePort              int
			startAndWaitForServer func()
		)

		BeforeEach(func() {
			serverDone = make(chan error, 1)

			// Get free port
			listener, err := net.Listen("tcp", ":0") //nolint:gosec // Binding to all interfaces intentionally for testing
			Expect(err).ToNot(HaveOccurred())
			tcpAddr, ok := listener.Addr().(*net.TCPAddr)
			Expect(ok).To(BeTrue(), "Expected TCP address")
			freePort = tcpAddr.Port
			err = listener.Close()
			Expect(err).ToNot(HaveOccurred())

			// Recreate server with the free port
			shutdown = util.NewGracefulShutdown()
			server = util.NewBaseHTTPServer("test", freePort, shutdown)

			startAndWaitForServer = func() {
				go func() {
					err := server.Start()
					serverDone <- err
				}()

				Eventually(func() error {
					conn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", freePort))
					if conn != nil {
						conn.Close()
					}
					return err
				}, 1*time.Second, 50*time.Millisecond).Should(Succeed())
			}
		})

		AfterEach(func() {
			select {
			case <-serverDone:
			default:
			}
		})

		It("should start and shutdown gracefully", func() {
			startAndWaitForServer()
			time.Sleep(100 * time.Millisecond)
			shutdown.Stop(context.Background())
			Eventually(serverDone, 2*time.Second).Should(Receive(BeNil()))
		})

		It("should serve statsviz endpoint", func() {
			startAndWaitForServer()
			time.Sleep(100 * time.Millisecond)

			resp, err := http.Get(fmt.Sprintf("http://localhost:%d/debug/statsviz/", freePort))
			if resp != nil {
				defer resp.Body.Close()
			}
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))

			shutdown.Stop(context.Background())
			Eventually(serverDone, 2*time.Second).Should(Receive(BeNil()))
		})

		It("should serve health endpoint", func() {
			startAndWaitForServer()
			time.Sleep(100 * time.Millisecond)

			resp, err := http.Get(fmt.Sprintf("http://localhost:%d/health", freePort))
			if resp != nil {
				defer resp.Body.Close()
			}
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))

			shutdown.Stop(context.Background())
			Eventually(serverDone, 2*time.Second).Should(Receive(BeNil()))
		})

		It("should call RegisterRoutes hook", func() {
			routesCalled := false
			server.RegisterRoutes = func(engine *gin.Engine) {
				routesCalled = true
				engine.GET("/custom", func(c *gin.Context) {
					c.JSON(http.StatusOK, gin.H{"custom": true})
				})
			}

			startAndWaitForServer()
			time.Sleep(100 * time.Millisecond)
			Expect(routesCalled).To(BeTrue())

			resp, err := http.Get(fmt.Sprintf("http://localhost:%d/custom", freePort))
			if resp != nil {
				defer resp.Body.Close()
			}
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))

			shutdown.Stop(context.Background())
			Eventually(serverDone, 2*time.Second).Should(Receive(BeNil()))
		})

		It("should call BeforeStart hook", func() {
			beforeStartCalled := false
			server.BeforeStart = func(_ context.Context) {
				beforeStartCalled = true
			}

			startAndWaitForServer()
			Expect(beforeStartCalled).To(BeTrue())

			shutdown.Stop(context.Background())
			Eventually(serverDone, 2*time.Second).Should(Receive(BeNil()))
		})

		It("should handle startup errors", func() {
			errShutdown := util.NewGracefulShutdown()
			errServer := util.NewBaseHTTPServer("test", -1, errShutdown)
			go func() {
				err := errServer.Start()
				serverDone <- err
			}()

			Eventually(serverDone, 500*time.Millisecond).Should(Receive(HaveOccurred()))
		})
	})
})
