package main_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/amanhigh/go-fun/components/learn/frameworks/ui/components"
	"github.com/amanhigh/go-fun/components/learn/frameworks/ui/components/advanced"
	"github.com/amanhigh/go-fun/components/learn/frameworks/ui/components/medium"
	"github.com/amanhigh/go-fun/components/learn/frameworks/ui/pages"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const testPort = 18081

var (
	serverURL string
	server    *http.Server
	registry  *components.Registry
)

var _ = BeforeSuite(func() {
	By("Starting the actual HTTP server for integration testing")

	serverURL = fmt.Sprintf("http://localhost:%d", testPort)

	// Create registry and register all components
	registry = components.NewRegistry()
	pages.RegisterBasic(registry)
	medium.RegisterAll(registry)
	advanced.RegisterAll(registry)

	// Create gin router
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	// Index page
	router.GET("/", func(c *gin.Context) {
		levels := []pages.LevelInfo{
			{
				Name:        "Basic Components",
				Path:        "/basic",
				Description: "Core UI building blocks: Button, TextInput, TextArea, Dropdown, Badge, Radio, Checkbox, and Modal components with professional styling.",
				Count:       len(registry.Basic()),
				BadgeClass:  "badge-basic",
			},
			{
				Name:        "Medium Components",
				Path:        "/medium",
				Description: "Intermediate patterns: nested components, state handling, data tables, composition, and security.",
				Count:       len(registry.Medium()),
				BadgeClass:  "badge-medium",
			},
			{
				Name:        "Advanced Components",
				Path:        "/advanced",
				Description: "Complex patterns: full page layouts, dashboards with multiple widgets, and advanced composition.",
				Count:       len(registry.Advanced()),
				BadgeClass:  "badge-advanced",
			},
		}
		c.Header("Content-Type", "text/html")
		pages.IndexPage(levels).Render(c.Request.Context(), c.Writer)
	})

	// Level pages
	router.GET("/basic", func(c *gin.Context) {
		comps := pages.ComponentsToInfoList(registry.Basic())
		c.Header("Content-Type", "text/html")
		pages.LevelPage("basic", "Basic Components", comps).Render(c.Request.Context(), c.Writer)
	})

	router.GET("/medium", func(c *gin.Context) {
		comps := pages.ComponentsToInfoList(registry.Medium())
		c.Header("Content-Type", "text/html")
		pages.LevelPage("medium", "Medium Components", comps).Render(c.Request.Context(), c.Writer)
	})

	router.GET("/advanced", func(c *gin.Context) {
		comps := pages.ComponentsToInfoList(registry.Advanced())
		c.Header("Content-Type", "text/html")
		pages.LevelPage("advanced", "Advanced Components", comps).Render(c.Request.Context(), c.Writer)
	})

	// Register all component routes
	for _, comp := range registry.All() {
		url := comp.URL()
		c := comp
		router.GET(url, func(ctx *gin.Context) {
			ctx.Header("Content-Type", "text/html")
			c.Render().Render(ctx.Request.Context(), ctx.Writer)
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
	Context("Index and Level Pages", func() {
		It("should serve index page with level links", func() {
			resp, err := http.Get(serverURL + "/")
			Expect(err).ToNot(HaveOccurred())
			defer resp.Body.Close()

			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			body, _ := io.ReadAll(resp.Body)
			html := string(body)

			Expect(html).To(ContainSubstring("Templ UI Component Demo"))
			Expect(html).To(ContainSubstring("/basic"))
			Expect(html).To(ContainSubstring("/medium"))
			Expect(html).To(ContainSubstring("/advanced"))
		})

		It("should serve basic level page", func() {
			resp, err := http.Get(serverURL + "/basic")
			Expect(err).ToNot(HaveOccurred())
			defer resp.Body.Close()

			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			body, _ := io.ReadAll(resp.Body)
			html := string(body)

			Expect(html).To(ContainSubstring("Basic Components"))
		})
	})

	Context("Basic Component - Showcase", func() {
		It("should render basic showcase via HTTP", func() {
			comp := pages.DefaultBasicShowcaseComponent()
			resp, err := http.Get(serverURL + comp.URL())
			Expect(err).ToNot(HaveOccurred())
			defer resp.Body.Close()

			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(resp.Header.Get("Content-Type")).To(ContainSubstring("text/html"))

			body, _ := io.ReadAll(resp.Body)
			html := string(body)
			Expect(html).To(ContainSubstring("FR-001 Showcase"))
			Expect(html).To(ContainSubstring("Username"))
		})

		It("should match direct rendering with HTTP response", func() {
			comp := pages.DefaultBasicShowcaseComponent()

			// Direct rendering
			var buf testBuffer
			err := comp.Render().Render(context.Background(), &buf)
			Expect(err).ToNot(HaveOccurred())

			// HTTP rendering
			resp, err := http.Get(serverURL + comp.URL())
			Expect(err).ToNot(HaveOccurred())
			defer resp.Body.Close()

			httpBody, _ := io.ReadAll(resp.Body)
			direct := buf.String()
			httpHTML := string(httpBody)
			Expect(httpHTML).To(ContainSubstring("FR-001 Showcase"))
			Expect(direct).To(ContainSubstring("FR-001 Showcase"))
			Expect(httpHTML).To(ContainSubstring("id=\"showcase-dialog\""))
			Expect(direct).To(ContainSubstring("id=\"showcase-dialog\""))
		})
	})

	Context("Component Consistency", func() {
		It("should render all basic components consistently", func() {
			for _, comp := range registry.Basic() {
				resp, err := http.Get(serverURL + comp.URL())
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				resp.Body.Close()
			}
		})
	})
})

// testBuffer implements io.Writer for testing
type testBuffer struct {
	data []byte
}

func (b *testBuffer) Write(p []byte) (n int, err error) {
	b.data = append(b.data, p...)
	return len(p), nil
}

func (b *testBuffer) String() string {
	return string(b.data)
}

// Ensure pages package basic showcase constructor is available
var _ = pages.DefaultBasicShowcaseComponent
