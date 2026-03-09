package util_test

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/models/config"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const (
	// testWaitTime is the standard timeout for test operations
	testWaitTime = 500 * time.Millisecond
	// testErrorWaitTime is shorter timeout for error paths
	testErrorWaitTime = 200 * time.Millisecond
	// testStartupTime is timeout for server startup
	testStartupTime = 500 * time.Millisecond
)

// testLifecycle implements ServerLifecycle with optional func overrides
type testLifecycle struct {
	registerRoutes func(*gin.Engine)
	beforeStart    func(context.Context)
	beforeShutdown func(context.Context)
	afterShutdown  func(context.Context)
}

func (t *testLifecycle) RegisterRoutes(e *gin.Engine) {
	if t.registerRoutes != nil {
		t.registerRoutes(e)
	}
}
func (t *testLifecycle) BeforeStart(ctx context.Context) {
	if t.beforeStart != nil {
		t.beforeStart(ctx)
	}
}
func (t *testLifecycle) BeforeShutdown(ctx context.Context) {
	if t.beforeShutdown != nil {
		t.beforeShutdown(ctx)
	}
}
func (t *testLifecycle) AfterShutdown(ctx context.Context) {
	if t.afterShutdown != nil {
		t.afterShutdown(ctx)
	}
}

var _ = Describe("HttpServer", func() {
	var (
		server   *util.HttpServerImpl
		shutdown util.Shutdown
	)

	BeforeEach(func() {
		gin.SetMode(gin.TestMode)
		shutdown = util.NewGracefulShutdown()
	})

	Context("Constructor", func() {
		It("should create server with correct configuration", func() {
			server = util.NewHttpServer(config.HttpServerConfig{Name: "test", Port: 0}, gin.Default(), shutdown)
			Expect(server).ToNot(BeNil())
			Expect(server.Name).To(Equal("test"))
			Expect(server.Engine).ToNot(BeNil())
			Expect(server.Server).ToNot(BeNil())
		})
	})

	Context("Server operations", func() {
		var (
			serverDone chan error
			freePort   int
		)

		BeforeEach(func() {
			serverDone = make(chan error, 1)
			freePort = getFreePort()
			server = util.NewHttpServer(config.HttpServerConfig{Name: "test", Port: freePort}, gin.Default(), shutdown)
		})

		AfterEach(func() {
			select {
			case <-serverDone:
			default:
			}
		})

		It("should start and shutdown gracefully", func() {
			startServer(server, serverDone, freePort)
			shutdown.Stop(context.Background())
			Eventually(serverDone, testWaitTime).Should(Receive(BeNil()))
		})

		It("should serve default endpoints", func() {
			startServer(server, serverDone, freePort)

			// Test health endpoint
			resp, err := http.Get(fmt.Sprintf("http://localhost:%d/health", freePort))
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			resp.Body.Close()

			// Test statsviz endpoint
			resp, err = http.Get(fmt.Sprintf("http://localhost:%d/debug/statsviz/", freePort))
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			resp.Body.Close()

			shutdown.Stop(context.Background())
			Eventually(serverDone, testWaitTime).Should(Receive(BeNil()))
		})

		It("should handle startup errors", func() {
			badServer := util.NewHttpServer(config.HttpServerConfig{Name: "bad", Port: -1}, gin.Default(), shutdown)
			go func() {
				serverDone <- badServer.Start()
			}()
			Eventually(serverDone, testWaitTime).Should(Receive(HaveOccurred()))
		})

		It("should handle impossible port", func() {
			badServer := util.NewHttpServer(config.HttpServerConfig{Name: "bad", Port: 65536}, gin.Default(), shutdown)
			err := badServer.Start()
			Expect(err).To(HaveOccurred())
		})
	})

	Context("Lifecycle hooks", func() {
		var (
			serverDone chan error
			freePort   int
		)

		BeforeEach(func() {
			serverDone = make(chan error, 1)
			freePort = getFreePort()
			server = util.NewHttpServer(config.HttpServerConfig{Name: "test", Port: freePort}, gin.Default(), shutdown)
		})

		AfterEach(func() {
			select {
			case <-serverDone:
			default:
			}
		})

		It("should call all lifecycle hooks in correct order", func() {
			var hookOrder []string
			server.SetLifecycle(&testLifecycle{
				registerRoutes: func(*gin.Engine) {
					hookOrder = append(hookOrder, "registerRoutes")
				},
				beforeStart: func(context.Context) {
					hookOrder = append(hookOrder, "beforeStart")
				},
				beforeShutdown: func(context.Context) {
					hookOrder = append(hookOrder, "beforeShutdown")
				},
				afterShutdown: func(context.Context) {
					hookOrder = append(hookOrder, "afterShutdown")
				},
			})

			startServer(server, serverDone, freePort)
			shutdown.Stop(context.Background())
			Eventually(serverDone, testWaitTime).Should(Receive(BeNil()))

			Expect(hookOrder).To(Equal([]string{
				"registerRoutes", "beforeStart", "beforeShutdown", "afterShutdown",
			}))
		})

		It("should register custom routes via lifecycle", func() {
			server.SetLifecycle(&testLifecycle{
				registerRoutes: func(engine *gin.Engine) {
					engine.GET("/custom", func(c *gin.Context) {
						c.JSON(http.StatusOK, gin.H{"custom": true})
					})
				},
			})

			startServer(server, serverDone, freePort)

			resp, err := http.Get(fmt.Sprintf("http://localhost:%d/custom", freePort))
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			resp.Body.Close()

			shutdown.Stop(context.Background())
			Eventually(serverDone, testWaitTime).Should(Receive(BeNil()))
		})

		It("should call BeforeStart hook", func() {
			called := false
			server.SetLifecycle(&testLifecycle{
				beforeStart: func(_ context.Context) {
					called = true
				},
			})

			startServer(server, serverDone, freePort)
			Expect(called).To(BeTrue())

			shutdown.Stop(context.Background())
			Eventually(serverDone, testWaitTime).Should(Receive(BeNil()))
		})

		It("should call BeforeShutdown hook", func() {
			called := false
			server.SetLifecycle(&testLifecycle{
				beforeShutdown: func(_ context.Context) {
					called = true
				},
			})

			startServer(server, serverDone, freePort)
			shutdown.Stop(context.Background())
			Eventually(serverDone, testWaitTime).Should(Receive(BeNil()))
			Expect(called).To(BeTrue())
		})

		It("should call AfterShutdown hook", func() {
			called := false
			server.SetLifecycle(&testLifecycle{
				afterShutdown: func(_ context.Context) {
					called = true
				},
			})

			startServer(server, serverDone, freePort)
			shutdown.Stop(context.Background())
			Eventually(serverDone, testWaitTime).Should(Receive(BeNil()))
			Expect(called).To(BeTrue())
		})
	})

	Context("Configuration", func() {
		It("should set correct server address and timeouts", func() {
			testServer := util.NewHttpServer(config.HttpServerConfig{Name: "test", Port: 8080}, gin.Default(), shutdown)
			Expect(testServer.Server.Addr).To(Equal(":8080"))
			Expect(testServer.Server.ReadHeaderTimeout).To(Equal(5 * time.Second))
			Expect(testServer.Server.ReadTimeout).To(Equal(5 * time.Second))
			Expect(testServer.Server.WriteTimeout).To(Equal(5 * time.Second))
		})

		It("should use provided gin engine", func() {
			customEngine := gin.New()
			testServer := util.NewHttpServer(config.HttpServerConfig{Name: "test", Port: 8080}, customEngine, shutdown)
			Expect(testServer.Engine).To(Equal(customEngine))
		})
	})

	Context("API methods", func() {
		It("should allow setting lifecycle multiple times", func() {
			server = util.NewHttpServer(config.HttpServerConfig{Name: "test", Port: 0}, gin.Default(), shutdown)

			// Test that lifecycle can be set and actually used
			called := false
			lifecycle := &testLifecycle{
				beforeStart: func(context.Context) {
					called = true
				},
			}

			server.SetLifecycle(lifecycle)

			// Verify the lifecycle is actually used by starting the server
			serverDone := make(chan error, 1)
			freePort := getFreePort()
			server = util.NewHttpServer(config.HttpServerConfig{Name: "test", Port: freePort}, gin.Default(), shutdown)
			server.SetLifecycle(lifecycle)

			startServer(server, serverDone, freePort)
			Expect(called).To(BeTrue()) // This proves the lifecycle was actually set and used

			shutdown.Stop(context.Background())
			Eventually(serverDone, testWaitTime).Should(Receive(BeNil()))
		})

		It("should stop server directly", func() {
			serverDone := make(chan error, 1)
			freePort := getFreePort()
			server = util.NewHttpServer(config.HttpServerConfig{Name: "test", Port: freePort}, gin.Default(), shutdown)

			startServer(server, serverDone, freePort)
			server.Stop(context.Background())
			Eventually(serverDone, testWaitTime).Should(Receive(BeNil()))
		})

		It("should handle multiple Stop calls gracefully", func() {
			serverDone := make(chan error, 1)
			freePort := getFreePort()
			server = util.NewHttpServer(config.HttpServerConfig{Name: "test", Port: freePort}, gin.Default(), shutdown)

			startServer(server, serverDone, freePort)
			server.Stop(context.Background())
			server.Stop(context.Background()) // Second call should not panic
			Eventually(serverDone, testWaitTime).Should(Receive(BeNil()))
		})

		It("should handle context cancellation during shutdown", func() {
			serverDone := make(chan error, 1)
			freePort := getFreePort()
			server = util.NewHttpServer(config.HttpServerConfig{Name: "test", Port: freePort}, gin.Default(), shutdown)

			startServer(server, serverDone, freePort)

			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			server.Stop(ctx)
			Eventually(serverDone, testWaitTime).Should(Receive(BeNil()))
		})

		It("should handle graceful shutdown failure", func() {
			// Test the error path in waitForShutdown (lines 140-142)
			serverDone := make(chan error, 1)
			freePort := getFreePort()
			server = util.NewHttpServer(config.HttpServerConfig{Name: "test", Port: freePort}, gin.Default(), shutdown)

			startServer(server, serverDone, freePort)

			// Force shutdown by closing the connection immediately
			// This should trigger an error in Server.Shutdown()
			go func() {
				// Give server a moment to start, then force close
				time.Sleep(10 * time.Millisecond)
				server.Server.Close()
			}()

			// Stop should handle the error gracefully
			server.Stop(context.Background())

			// Should complete quickly since we're testing error handling
			Eventually(serverDone, testErrorWaitTime).Should(Receive())
		})
	})
})

// getFreePort returns an available port for testing
func getFreePort() int {
	listener, err := net.Listen("tcp", ":0") //nolint:gosec // Binding to all interfaces intentionally for testing
	Expect(err).ToNot(HaveOccurred())
	defer listener.Close()
	return listener.Addr().(*net.TCPAddr).Port
}

// startServer starts the server and waits for it to be ready
func startServer(server *util.HttpServerImpl, serverDone chan error, port int) {
	go func() {
		serverDone <- server.Start()
	}()

	Eventually(func() error {
		conn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", port))
		if conn != nil {
			conn.Close()
		}
		return err
	}, testStartupTime, 10*time.Millisecond).Should(Succeed())
}
