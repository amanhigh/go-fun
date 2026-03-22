package main_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/amanhigh/go-fun/components/learn/frameworks/frontend/ui/components"
	"github.com/amanhigh/go-fun/components/learn/frameworks/frontend/ui/pages"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const testPort = 18081

var (
	server    *http.Server
	serverURL string
)

var _ = BeforeSuite(func() {
	By("Starting the actual HTTP server for integration testing")

	serverURL = fmt.Sprintf("http://localhost:%d", testPort)

	// Create gin router using same pattern as server
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	// Create components once
	components := []components.Component{
		pages.NewFormShowcaseComponent(),
		pages.NewHelloComponent(),
	}

	// Serve static files (JS, CSS, images)
	router.Static("/assets", "../assets")

	// Index page
	router.GET("/", func(c *gin.Context) {
		c.Header("Content-Type", "text/html")
		pages.IndexPage(components).Render(c.Request.Context(), c.Writer)
	})

	// Register component routes
	for _, comp := range components {
		comp := comp // capture for closure
		router.GET(comp.URL(), func(ctx *gin.Context) {
			ctx.Header("Content-Type", "text/html")
			comp.Render().Render(ctx.Request.Context(), ctx.Writer)
		})
	}

	server = &http.Server{
		Addr:              fmt.Sprintf(":%d", testPort),
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		server.ListenAndServe()
	}()

	// Wait for server to be ready
	Eventually(func() error {
		resp, err := http.Get(serverURL) //nolint:gosec // Test URL is constant
		if err != nil {
			return err
		}
		resp.Body.Close()
		return nil
	}, "10s", "100ms").Should(Succeed())

	By(fmt.Sprintf("Server started successfully on %s", serverURL))
})

var _ = AfterSuite(func() {
	By("Stopping the HTTP server")
	if server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		server.Shutdown(ctx)
	}
})

// Server Smoke Tests - Tests one component from each level on real HTTP server
// This validates the full HTTP stack works correctly for each complexity level
var _ = Describe("Server Smoke Tests", func() {
	Context("Index and Component Pages", func() {
		It("should serve index page with component links", func() {
			resp, err := http.Get(serverURL + "/")
			Expect(err).ToNot(HaveOccurred())
			defer resp.Body.Close()

			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			body, _ := io.ReadAll(resp.Body)
			html := string(body)

			Expect(html).To(ContainSubstring("Templ UI Component Demo"))
			Expect(html).To(ContainSubstring("/form"))
			Expect(html).To(ContainSubstring("/hello"))
		})

		It("should serve individual component pages", func() {
			testCases := []string{"/form", "/hello"}
			for _, url := range testCases {
				resp, err := http.Get(serverURL + url)
				Expect(err).ToNot(HaveOccurred())
				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
			}
		})
	})

	Context("Component Content Validation", func() {
		It("should render form showcase with expected content", func() {
			resp, err := http.Get(serverURL + "/form")
			Expect(err).ToNot(HaveOccurred())
			defer resp.Body.Close()

			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			body, _ := io.ReadAll(resp.Body)
			html := string(body)
			Expect(html).To(ContainSubstring("Basic Components Showcase"))
			Expect(html).To(ContainSubstring("Username"))
		})

		It("should render hello page with expected content", func() {
			resp, err := http.Get(serverURL + "/hello")
			Expect(err).ToNot(HaveOccurred())
			defer resp.Body.Close()

			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			body, _ := io.ReadAll(resp.Body)
			html := string(body)
			Expect(html).To(ContainSubstring("Hello World Showcase"))
			Expect(html).To(ContainSubstring("Country"))
		})
	})
})
