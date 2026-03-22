package main_test

import (
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/components/learn/frameworks/frontend/ui/components"
	"github.com/amanhigh/go-fun/components/learn/frameworks/frontend/ui/pages"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// Test registration is in demo_suite_test.go

var _ = Describe("UI Component Handler Tests", func() {
	var (
		router   *gin.Engine
		registry *components.Registry
	)

	BeforeEach(func() {
		gin.SetMode(gin.TestMode)
		router = gin.New()
		registry = components.NewRegistry()

		// Serve static files (JS, CSS, images)
		router.Static("/assets", "../assets")

		// Create registry and register all components
		pages.RegisterBasic(registry)

		// Register component routes
		for _, comp := range registry.All() {
			url := comp.URL()
			c := comp // capture for closure
			router.GET(url, func(ctx *gin.Context) {
				ctx.Header("Content-Type", "text/html")
				c.Render().Render(ctx.Request.Context(), ctx.Writer)
			})
		}
	})

	Context("Form Essentials Showcase", func() {
		var (
			comp *pages.FormShowcaseComponent
			w    *httptest.ResponseRecorder
			html string
		)

		BeforeEach(func() {
			comp = pages.NewFormShowcaseComponent()
			var req *http.Request
			req, w = util.CreateHTMLTestRequest("GET", comp.URL())
			router.ServeHTTP(w, req)
			html = w.Body.String()
		})

		It("should render unified showcase with all core components", func() {
			Expect(w.Code).To(Equal(http.StatusOK))
			Expect(w.Header().Get("Content-Type")).To(ContainSubstring("text/html"))
			Expect(html).To(ContainSubstring("Basic Components Showcase"))
		})

		It("should showcase Button components", func() {
			// Verify button actions are present
			Expect(html).To(ContainSubstring("Save Draft"))
			Expect(html).To(ContainSubstring("Submit"))
			Expect(html).To(ContainSubstring("Preview Modal"))
		})

		It("should showcase Text Input components", func() {
			// Verify inputs with labels
			Expect(html).To(ContainSubstring("Username"))
			Expect(html).To(ContainSubstring("Email"))
			Expect(html).To(ContainSubstring("id=\"username\""))
			Expect(html).To(ContainSubstring("id=\"email\""))
			Expect(html).To(ContainSubstring("type=\"email\""))
		})

		It("should showcase Text Area component", func() {
			// Verify textarea
			Expect(html).To(ContainSubstring("Notes"))
			Expect(html).To(ContainSubstring("<textarea"))
			Expect(html).To(ContainSubstring("id=\"notes\""))
		})

		It("should showcase Select Box component", func() {
			// Verify select box for country
			Expect(html).To(ContainSubstring("Country"))
			Expect(html).To(ContainSubstring("id=\"country\""))
			Expect(html).To(ContainSubstring("United States"))
			Expect(html).To(ContainSubstring("India"))
		})

		It("should showcase Badge components", func() {
			// Verify status badges
			Expect(html).To(ContainSubstring("Ready"))
			Expect(html).To(ContainSubstring("Review"))
			Expect(html).To(ContainSubstring("Info"))
		})

		It("should showcase Radio Button component", func() {
			// Verify radio buttons for plan selection
			Expect(html).To(ContainSubstring("Subscription Plan"))
			Expect(html).To(ContainSubstring("Starter"))
			Expect(html).To(ContainSubstring("Pro"))
			Expect(html).To(ContainSubstring("id=\"plan-starter\""))
			Expect(html).To(ContainSubstring("id=\"plan-pro\""))
		})

		It("should showcase Checkbox component", func() {
			// Verify checkbox for terms
			Expect(html).To(ContainSubstring("terms and conditions"))
			Expect(html).To(ContainSubstring("id=\"terms\""))
		})

		It("should showcase Modal/Dialog component", func() {
			// Verify dialog structure
			Expect(html).To(ContainSubstring("id=\"showcase-dialog\""))
			Expect(html).To(ContainSubstring("Submission Preview"))
			Expect(html).To(ContainSubstring("Confirm"))
		})

		It("should implement Component interface for form showcase component", func() {
			var _ components.Component = pages.NewFormShowcaseComponent()

			showcase := pages.NewFormShowcaseComponent()
			Expect(showcase.Name()).To(Equal("form-showcase"))
			Expect(showcase.Level()).To(Equal(components.LevelBasic))
		})
	})

	Context("Hello Page", func() {
		var (
			comp *pages.HelloComponent
			w    *httptest.ResponseRecorder
			html string
		)

		BeforeEach(func() {
			comp = pages.NewHelloComponent()
			w = httptest.NewRecorder()
			req, _ := http.NewRequest("GET", comp.URL(), nil)
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusOK))
			body, _ := io.ReadAll(w.Body)
			html = string(body)
		})

		It("should render hello page with selectbox", func() {
			Expect(html).To(ContainSubstring("Hello World Showcase"))
			Expect(html).To(ContainSubstring("Demonstrating both TemplUI components and native HTML templates"))
		})

		It("should showcase Select Box component", func() {
			// Verify select box structure
			Expect(html).To(ContainSubstring("Country"))
			Expect(html).To(ContainSubstring("id=\"country\""))
			Expect(html).To(ContainSubstring("Choose your country"))
			Expect(html).To(ContainSubstring("United States"))
			Expect(html).To(ContainSubstring("India"))
			Expect(html).To(ContainSubstring("United Kingdom"))
		})

		It("should include selectbox script", func() {
			// Verify selectbox JavaScript is loaded
			Expect(html).To(ContainSubstring("selectbox.min.js"))
		})

		It("should include Tailwind CSS", func() {
			// Verify Tailwind CSS is loaded
			Expect(html).To(ContainSubstring("/assets/css/app.css"))
		})

		It("should implement Component interface for hello component", func() {
			var _ components.Component = pages.NewHelloComponent()

			hello := pages.NewHelloComponent()
			Expect(hello.Name()).To(Equal("hello"))
			Expect(hello.Level()).To(Equal(components.LevelBasic))
			Expect(hello.URL()).To(Equal("/hello"))
		})
	})

	Context("Component Registry", func() {
		It("should have correct number of basic components", func() {
			Expect(registry.Basic()).To(HaveLen(2))
		})

		It("should find component by URL", func() {
			comp := registry.FindByURL("/form")
			Expect(comp).ToNot(BeNil())
			Expect(comp.Name()).To(Equal("form-showcase"))
		})

		It("should return nil for unknown URL", func() {
			comp := registry.FindByURL("/unknown/path")
			Expect(comp).To(BeNil())
		})

		It("should serve static JS files", func() {
			req, w := util.CreateHTMLTestRequest("GET", "/assets/js/app.js")
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusOK))
			Expect(w.Header().Get("Content-Type")).To(ContainSubstring("javascript"))
		})

		It("should serve CSS files", func() {
			req, w := util.CreateHTMLTestRequest("GET", "/assets/css/showcase.css")
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusOK))
			Expect(w.Header().Get("Content-Type")).To(ContainSubstring("text/css"))
			Expect(w.Body.String()).To(ContainSubstring("CSS_LOADED_SUCCESSFULLY"))
		})

		It("should serve image files", func() {
			req, w := util.CreateHTMLTestRequest("GET", "/assets/images/sample-logo.png")
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusOK))
			Expect(w.Header().Get("Content-Type")).To(ContainSubstring("image"))
		})

		It("should return 404 for non-existent static files", func() {
			req, w := util.CreateHTMLTestRequest("GET", "/assets/css/nonexistent.css")
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusNotFound))
		})
	})
})
