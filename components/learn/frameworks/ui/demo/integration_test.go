package main_test

import (
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/components/learn/frameworks/ui/components"
	"github.com/amanhigh/go-fun/components/learn/frameworks/ui/components/advanced"
	"github.com/amanhigh/go-fun/components/learn/frameworks/ui/components/medium"
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

		// Register all components
		pages.RegisterBasic(registry)
		medium.RegisterAll(registry)
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
		It("should render unified showcase with all core components", func() {
			comp := pages.DefaultBasicShowcaseComponent()
			req, w := util.CreateHTMLTestRequest("GET", comp.URL())
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusOK))
			Expect(w.Header().Get("Content-Type")).To(ContainSubstring("text/html"))

			html := w.Body.String()
			Expect(html).To(ContainSubstring("Basic Components Showcase"))
		})

		It("should showcase Button components", func() {
			comp := pages.DefaultBasicShowcaseComponent()
			req, w := util.CreateHTMLTestRequest("GET", comp.URL())
			router.ServeHTTP(w, req)

			html := w.Body.String()
			// Verify button actions are present
			Expect(html).To(ContainSubstring("Save Draft"))
			Expect(html).To(ContainSubstring("Submit"))
			Expect(html).To(ContainSubstring("Preview Modal"))
		})

		It("should showcase Text Input components", func() {
			comp := pages.DefaultBasicShowcaseComponent()
			req, w := util.CreateHTMLTestRequest("GET", comp.URL())
			router.ServeHTTP(w, req)

			html := w.Body.String()
			// Verify inputs with labels
			Expect(html).To(ContainSubstring("Username"))
			Expect(html).To(ContainSubstring("Email"))
			Expect(html).To(ContainSubstring("id=\"username\""))
			Expect(html).To(ContainSubstring("id=\"email\""))
			Expect(html).To(ContainSubstring("type=\"email\""))
		})

		It("should showcase Text Area component", func() {
			comp := pages.DefaultBasicShowcaseComponent()
			req, w := util.CreateHTMLTestRequest("GET", comp.URL())
			router.ServeHTTP(w, req)

			html := w.Body.String()
			// Verify textarea
			Expect(html).To(ContainSubstring("Notes"))
			Expect(html).To(ContainSubstring("<textarea"))
			Expect(html).To(ContainSubstring("id=\"notes\""))
		})

		It("should showcase Select/Dropdown component", func() {
			comp := pages.DefaultBasicShowcaseComponent()
			req, w := util.CreateHTMLTestRequest("GET", comp.URL())
			router.ServeHTTP(w, req)

			html := w.Body.String()
			// Verify select box for country
			Expect(html).To(ContainSubstring("Country"))
			Expect(html).To(ContainSubstring("id=\"country\""))
			Expect(html).To(ContainSubstring("United States"))
			Expect(html).To(ContainSubstring("India"))
		})

		It("should showcase Badge components", func() {
			comp := pages.DefaultBasicShowcaseComponent()
			req, w := util.CreateHTMLTestRequest("GET", comp.URL())
			router.ServeHTTP(w, req)

			html := w.Body.String()
			// Verify status badges
			Expect(html).To(ContainSubstring("Ready"))
			Expect(html).To(ContainSubstring("Review"))
			Expect(html).To(ContainSubstring("Info"))
		})

		It("should showcase Radio Button components", func() {
			comp := pages.DefaultBasicShowcaseComponent()
			req, w := util.CreateHTMLTestRequest("GET", comp.URL())
			router.ServeHTTP(w, req)

			html := w.Body.String()
			// Verify radio buttons for plan selection
			Expect(html).To(ContainSubstring("Subscription Plan"))
			Expect(html).To(ContainSubstring("Starter"))
			Expect(html).To(ContainSubstring("Pro"))
			Expect(html).To(ContainSubstring("id=\"plan-starter\""))
			Expect(html).To(ContainSubstring("id=\"plan-pro\""))
		})

		It("should showcase Checkbox component", func() {
			comp := pages.DefaultBasicShowcaseComponent()
			req, w := util.CreateHTMLTestRequest("GET", comp.URL())
			router.ServeHTTP(w, req)

			html := w.Body.String()
			// Verify checkbox for terms
			Expect(html).To(ContainSubstring("terms and conditions"))
			Expect(html).To(ContainSubstring("id=\"terms\""))
		})

		It("should showcase Modal/Dialog component", func() {
			comp := pages.DefaultBasicShowcaseComponent()
			req, w := util.CreateHTMLTestRequest("GET", comp.URL())
			router.ServeHTTP(w, req)

			html := w.Body.String()
			// Verify dialog structure
			Expect(html).To(ContainSubstring("id=\"showcase-dialog\""))
			Expect(html).To(ContainSubstring("Submission Preview"))
			Expect(html).To(ContainSubstring("Confirm"))
		})

		It("should implement Component interface for showcase component", func() {
			var _ components.Component = pages.DefaultBasicShowcaseComponent()

			showcase := pages.DefaultBasicShowcaseComponent()
			Expect(showcase.Name()).To(Equal("basic-showcase"))
			Expect(showcase.Level()).To(Equal(components.LevelBasic))
		})
	})

	Context("Medium Components", func() {
		It("should render nested component", func() {
			comp := medium.DefaultNestedComponent()
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", comp.URL(), nil)
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusOK))
			body, _ := io.ReadAll(w.Body)
			html := string(body)
			Expect(html).To(ContainSubstring("Welcome Page"))
			Expect(html).To(ContainSubstring("Hello, Bob!"))
		})

		It("should render counter component", func() {
			comp := medium.DefaultCounterComponent()
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", comp.URL(), nil)
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusOK))
			body, _ := io.ReadAll(w.Body)
			html := string(body)
			Expect(html).To(ContainSubstring("Counter Value"))
			Expect(html).To(ContainSubstring("Counter is positive"))
		})

		It("should render datatable component", func() {
			comp := medium.DefaultDataTableComponent()
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", comp.URL(), nil)
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusOK))
			body, _ := io.ReadAll(w.Body)
			html := string(body)
			Expect(html).To(ContainSubstring("<table"))
			Expect(html).To(ContainSubstring("Alice"))
			Expect(html).To(ContainSubstring("Bob"))
		})

		It("should render composed component", func() {
			comp := medium.DefaultComposedComponent()
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", comp.URL(), nil)
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusOK))
			body, _ := io.ReadAll(w.Body)
			html := string(body)
			Expect(html).To(ContainSubstring("Team"))
			Expect(html).To(ContainSubstring("Review code"))
		})

		It("should render xss component with escaped content", func() {
			comp := medium.DefaultXSSComponent()
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", comp.URL(), nil)
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusOK))
			body, _ := io.ReadAll(w.Body)
			html := string(body)
			Expect(html).ToNot(ContainSubstring("<script>"))
			Expect(html).To(ContainSubstring("&lt;script&gt;"))
		})

		It("should render emptytable component", func() {
			comp := medium.DefaultEmptyTableComponent()
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", comp.URL(), nil)
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusOK))
			body, _ := io.ReadAll(w.Body)
			html := string(body)
			Expect(html).To(ContainSubstring("<table"))
			Expect(html).To(ContainSubstring("<thead>"))
		})
	})

	Context("Advanced Components", func() {
		It("should render fullpage component", func() {
			comp := advanced.DefaultFullPageComponent()
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", comp.URL(), nil)
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusOK))
			body, _ := io.ReadAll(w.Body)
			html := string(body)
			Expect(html).To(ContainSubstring("<!doctype html>"))
			Expect(html).To(ContainSubstring("Advanced Full Page Demo"))
		})

		It("should render dashboard component", func() {
			comp := advanced.DefaultDashboardComponent()
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", comp.URL(), nil)
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusOK))
			body, _ := io.ReadAll(w.Body)
			html := string(body)
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
			Expect(registry.Medium()).To(HaveLen(6))
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
	})
})
