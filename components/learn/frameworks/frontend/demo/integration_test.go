package main_test

import (
	"net/http"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/components/learn/frameworks/frontend/ui/components"
	"github.com/amanhigh/go-fun/components/learn/frameworks/frontend/ui/pages"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// Test registration is in demo_suite_test.go

var _ = Describe("UI Server Integration Tests", func() {
	var (
		router *gin.Engine
		comps  []components.Component
	)

	BeforeEach(func() {
		gin.SetMode(gin.TestMode)
		router = gin.New()

		// Serve static files (JS, CSS, images)
		router.Static("/assets", "../assets")

		// Create components and register routes like the server does
		comps = []components.Component{
			pages.NewFormShowcaseComponent(),
			pages.NewHelloComponent(),
		}

		// Register component routes
		for _, comp := range comps {
			comp := comp // capture for closure
			router.GET(comp.URL(), func(c *gin.Context) {
				c.Header("Content-Type", "text/html")
				comp.Render().Render(c.Request.Context(), c.Writer)
			})
		}

		// Add index page
		router.GET("/", func(c *gin.Context) {
			c.Header("Content-Type", "text/html")
			pages.IndexPage(comps).Render(c.Request.Context(), c.Writer)
		})
	})

	Context("Server Routes", func() {
		It("should register all component routes", func() {
			// Test that all component routes are registered and working
			for _, comp := range comps {
				req, w := util.CreateHTMLTestRequest("GET", comp.URL())
				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusOK))
				Expect(w.Header().Get("Content-Type")).To(ContainSubstring("text/html"))
			}
		})

		It("should serve index page with component links", func() {
			req, w := util.CreateHTMLTestRequest("GET", "/")
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusOK))
			html := w.Body.String()

			// Verify index page structure
			Expect(html).To(ContainSubstring("Templ UI Component Demo"))
			Expect(html).To(ContainSubstring("/form"))
			Expect(html).To(ContainSubstring("/hello"))
		})

		It("should return 404 for unknown routes", func() {
			req, w := util.CreateHTMLTestRequest("GET", "/unknown")
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusNotFound))
		})
	})

	Context("Static Assets", func() {
		Context("CSS Files", func() {
			It("should serve CSS files with correct content type", func() {
				req, w := util.CreateHTMLTestRequest("GET", "/assets/css/showcase.css")
				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusOK))
				Expect(w.Header().Get("Content-Type")).To(ContainSubstring("text/css"))
				Expect(w.Body.String()).To(ContainSubstring("CSS_LOADED_SUCCESSFULLY"))
			})
		})

		Context("JavaScript Files", func() {
			It("should serve JavaScript files with correct content type", func() {
				req, w := util.CreateHTMLTestRequest("GET", "/assets/js/app.js")
				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusOK))
				Expect(w.Header().Get("Content-Type")).To(ContainSubstring("javascript"))
			})
		})

		Context("Image Files", func() {
			It("should serve image files with correct content type", func() {
				req, w := util.CreateHTMLTestRequest("GET", "/assets/images/sample-logo.png")
				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusOK))
				Expect(w.Header().Get("Content-Type")).To(ContainSubstring("image"))
			})
		})

		It("should return 404 for non-existent static files", func() {
			req, w := util.CreateHTMLTestRequest("GET", "/assets/css/nonexistent.css")
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusNotFound))
		})
	})
})
