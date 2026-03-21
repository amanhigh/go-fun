package main_test

import (
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/components/learn/frameworks/ui/components"
	"github.com/amanhigh/go-fun/components/learn/frameworks/ui/components/advanced"
	"github.com/amanhigh/go-fun/components/learn/frameworks/ui/pages"
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
		router.Static("/static", "../static")

		// Create registry and register all components
		pages.RegisterBasic(registry)
		pages.RegisterMedium(registry)
		advanced.RegisterAll(registry)

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

	Context("Basic Components - Core Showcase", func() {
		var (
			comp *pages.BasicShowcaseComponent
			w    *httptest.ResponseRecorder
			html string
		)

		BeforeEach(func() {
			comp = pages.NewBasicShowcaseComponent()
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

		It("should implement Component interface for showcase component", func() {
			var _ components.Component = pages.NewBasicShowcaseComponent()

			showcase := pages.NewBasicShowcaseComponent()
			Expect(showcase.Name()).To(Equal("basic-showcase"))
			Expect(showcase.Level()).To(Equal(components.LevelBasic))
		})
	})

	Context("Medium Components", func() {
		var (
			comp *pages.MediumShowcaseComponent
			w    *httptest.ResponseRecorder
			html string
		)

		BeforeEach(func() {
			comp = pages.NewMediumShowcaseComponent()
			w = httptest.NewRecorder()
			req, _ := http.NewRequest("GET", comp.URL(), nil)
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusOK))
			body, _ := io.ReadAll(w.Body)
			html = string(body)
		})

		It("should render medium showcase page with all components", func() {
			Expect(html).To(ContainSubstring("Medium Components Showcase"))
			Expect(html).To(ContainSubstring("Layout & Content Blocks"))
		})

		It("should render Card components section", func() {
			Expect(html).To(ContainSubstring("Card Components"))
			Expect(html).To(ContainSubstring("Cards provide flexible content containers"))
			Expect(html).To(ContainSubstring("Basic Card"))
			Expect(html).To(ContainSubstring("Feature Card"))
			Expect(html).To(ContainSubstring("Product Showcase"))
			Expect(html).To(ContainSubstring("Analytics Dashboard"))
			Expect(html).To(ContainSubstring("card component structure"))
			Expect(html).To(ContainSubstring("key performance indicators"))
		})

		It("should render Data Tables section", func() {
			Expect(html).To(ContainSubstring("Data Tables"))
			Expect(html).To(ContainSubstring("Structured data presentation"))
			Expect(html).To(ContainSubstring("Sample user data"))
			Expect(html).To(ContainSubstring("Alice Johnson"))
			Expect(html).To(ContainSubstring("Bob Smith"))
			Expect(html).To(ContainSubstring("Carol Davis"))
			Expect(html).To(ContainSubstring("Developer"))
			Expect(html).To(ContainSubstring("Designer"))
			Expect(html).To(ContainSubstring("Manager"))
		})

		It("should render Status Indicators section", func() {
			Expect(html).To(ContainSubstring("Status Indicators"))
			Expect(html).To(ContainSubstring("Visual state representations"))
			Expect(html).To(ContainSubstring("Project Status"))
			Expect(html).To(ContainSubstring("✅ Ready"))
			Expect(html).To(ContainSubstring("🕒 In Progress"))
			Expect(html).To(ContainSubstring("⚠️ Review"))
			Expect(html).To(ContainSubstring("❌ Blocked"))
			Expect(html).To(ContainSubstring("ℹ️ Info"))
			Expect(html).To(ContainSubstring("Priority Levels"))
			Expect(html).To(ContainSubstring("🔴 High"))
			Expect(html).To(ContainSubstring("🟡 Medium"))
			Expect(html).To(ContainSubstring("🟢 Low"))
		})

		It("should render Content Organization section", func() {
			Expect(html).To(ContainSubstring("Content Organization"))
			Expect(html).To(ContainSubstring("Content hierarchy and organization"))
			Expect(html).To(ContainSubstring("Project Overview"))
			Expect(html).To(ContainSubstring("Key Objectives"))
			Expect(html).To(ContainSubstring("Demonstrate card component usage"))
			Expect(html).To(ContainSubstring("Show table data presentation"))
			Expect(html).To(ContainSubstring("Display status indicators"))
			Expect(html).To(ContainSubstring("Organize content hierarchically"))
			Expect(html).To(ContainSubstring("Last updated: 2024-03-21"))
		})

		It("should implement Component interface for medium showcase component", func() {
			var _ components.Component = pages.NewMediumShowcaseComponent()

			showcase := pages.NewMediumShowcaseComponent()
			Expect(showcase.Name()).To(Equal("medium-showcase"))
			Expect(showcase.Level()).To(Equal(components.LevelMedium))
		})
	})

	Context("Advanced Components", func() {
		var (
			w    *httptest.ResponseRecorder
			html string
		)

		BeforeEach(func() {
			w = httptest.NewRecorder()
		})

		It("should render fullpage component", func() {
			comp := advanced.DefaultFullPageComponent()
			req, _ := http.NewRequest("GET", comp.URL(), nil)
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusOK))
			body, _ := io.ReadAll(w.Body)
			html = string(body)
			Expect(html).To(ContainSubstring("<!doctype html>"))
			Expect(html).To(ContainSubstring("Advanced Full Page Demo"))
		})

		It("should render dashboard component", func() {
			comp := advanced.DefaultDashboardComponent()
			req, _ := http.NewRequest("GET", comp.URL(), nil)
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusOK))
			body, _ := io.ReadAll(w.Body)
			html = string(body)
			Expect(html).To(ContainSubstring("Dashboard"))
			Expect(html).To(ContainSubstring("Admin"))
			Expect(html).To(ContainSubstring("<table"))
		})
	})

	Context("Component Registry", func() {
		It("should have correct number of basic components", func() {
			Expect(registry.Basic()).To(HaveLen(1))
		})

		It("should have correct number of medium components", func() {
			Expect(registry.Medium()).To(HaveLen(1))
		})

		It("should have correct number of advanced components", func() {
			Expect(registry.Advanced()).To(HaveLen(2))
		})

		It("should find component by URL", func() {
			comp := registry.FindByURL("/basic/showcase")
			Expect(comp).ToNot(BeNil())
			Expect(comp.Name()).To(Equal("basic-showcase"))
		})

		It("should return nil for unknown URL", func() {
			comp := registry.FindByURL("/unknown/path")
			Expect(comp).To(BeNil())
		})

		It("should serve static JS files", func() {
			req, w := util.CreateHTMLTestRequest("GET", "/static/js/basic.js")
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusOK))
			Expect(w.Header().Get("Content-Type")).To(ContainSubstring("javascript"))
		})

		It("should serve CSS files", func() {
			req, w := util.CreateHTMLTestRequest("GET", "/static/css/showcase.css")
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusOK))
			Expect(w.Header().Get("Content-Type")).To(ContainSubstring("text/css"))
			Expect(w.Body.String()).To(ContainSubstring("CSS_LOADED_SUCCESSFULLY"))
		})

		It("should serve image files", func() {
			req, w := util.CreateHTMLTestRequest("GET", "/static/images/sample-logo.png")
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusOK))
			Expect(w.Header().Get("Content-Type")).To(ContainSubstring("image"))
		})

		It("should return 404 for non-existent static files", func() {
			req, w := util.CreateHTMLTestRequest("GET", "/static/css/nonexistent.css")
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusNotFound))
		})
	})
})
