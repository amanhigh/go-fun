package main_test

import (
	"net/http"

	"github.com/amanhigh/go-fun/common/util"
	demo "github.com/amanhigh/go-fun/components/learn/frameworks/frontend/demo"
	"github.com/amanhigh/go-fun/components/learn/frameworks/frontend/ui/components"
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
		// Create gin router
		router = gin.New()

		// Use server functions for consistent setup
		comps = demo.CreateComponents()
		demo.SetupRoutes(router, comps)
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

		Context("TypeScript Source Files", func() {
			It("should fetch the custom TypeScript source file", func() {
				req, w := util.CreateHTMLTestRequest("GET", "/assets/js/custom.ts")
				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusOK))
				Expect(w.Body.String()).To(ContainSubstring("AlpineCounterElement"))
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
