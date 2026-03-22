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

	// Create gin router
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	// Serve static files (JS, CSS, images)
	router.Static("/assets", "../assets")

	// Index page
	router.GET("/", func(c *gin.Context) {
		levels := []pages.LevelInfo{
			{
				Name:        "Form Essentials",
				Path:        "/form",
				Description: "Master form inputs, validation, and user data collection with professional UI components.",
				BadgeClass:  "badge-basic",
			},
		}
		c.Header("Content-Type", "text/html")
		pages.IndexPage(levels).Render(c.Request.Context(), c.Writer)
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
			Expect(html).To(ContainSubstring("/form"))
		})

		It("should serve individual component pages", func() {
			resp, err := http.Get(serverURL + "/form")
			Expect(err).ToNot(HaveOccurred())
			defer resp.Body.Close()

			Expect(resp.StatusCode).To(Equal(http.StatusOK))
		})
	})

	Context("Form Essentials Showcase", func() {
		var (
			cshowcase *pages.FormShowcaseComponent
		)

		BeforeEach(func() {
			cshowcase = pages.NewFormShowcaseComponent()
		})

		It("should render basic showcase via HTTP", func() {
			resp, err := http.Get(serverURL + cshowcase.URL())
			Expect(err).ToNot(HaveOccurred())
			defer resp.Body.Close()

			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(resp.Header.Get("Content-Type")).To(ContainSubstring("text/html"))

			body, _ := io.ReadAll(resp.Body)
			html := string(body)
			Expect(html).To(ContainSubstring("Basic Components Showcase"))
			Expect(html).To(ContainSubstring("Username"))
		})

		It("should match direct rendering with HTTP response", func() {
			// Direct rendering
			var buf testBuffer
			err := cshowcase.Render().Render(context.Background(), &buf)
			Expect(err).ToNot(HaveOccurred())

			// HTTP rendering
			resp, err := http.Get(serverURL + cshowcase.URL())
			Expect(err).ToNot(HaveOccurred())
			defer resp.Body.Close()

			httpBody, _ := io.ReadAll(resp.Body)
			direct := buf.String()
			httpHTML := string(httpBody)
			Expect(httpHTML).To(ContainSubstring("Basic Components Showcase"))
			Expect(direct).To(ContainSubstring("Basic Components Showcase"))
			Expect(httpHTML).To(ContainSubstring("id=\"showcase-dialog\""))
			Expect(direct).To(ContainSubstring("id=\"showcase-dialog\""))
		})
	})

	Context("Hello Page", func() {
		var (
			hshowcase *pages.HelloComponent
		)

		BeforeEach(func() {
			hshowcase = pages.NewHelloComponent()
		})

		It("should render hello page via HTTP", func() {
			resp, err := http.Get(serverURL + hshowcase.URL())
			Expect(err).ToNot(HaveOccurred())
			defer resp.Body.Close()

			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(resp.Header.Get("Content-Type")).To(ContainSubstring("text/html"))

			body, _ := io.ReadAll(resp.Body)
			html := string(body)
			Expect(html).To(ContainSubstring("Hello World Showcase"))
			Expect(html).To(ContainSubstring("Country"))
			Expect(html).To(ContainSubstring("selectbox.min.js"))
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

	Context("Static File Serving", func() {
		It("should serve JavaScript files", func() {
			resp, err := http.Get(serverURL + "/assets/js/app.js")
			Expect(err).ToNot(HaveOccurred())
			defer resp.Body.Close()

			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(resp.Header.Get("Content-Type")).To(ContainSubstring("javascript"))

			body, _ := io.ReadAll(resp.Body)
			Expect(string(body)).To(ContainSubstring("alpine:init"))
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

// Ensure pages package form showcase constructor is available
var _ = pages.NewFormShowcaseComponent()
