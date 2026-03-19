package main_test

import (
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/amanhigh/go-fun/components/learn/frameworks/ui/demo/components"
	"github.com/amanhigh/go-fun/components/learn/frameworks/ui/demo/components/advanced"
	"github.com/amanhigh/go-fun/components/learn/frameworks/ui/demo/components/basic"
	"github.com/amanhigh/go-fun/components/learn/frameworks/ui/demo/components/medium"
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
		basic.RegisterAll(registry)
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

	Context("Basic Components", func() {
		It("should render greeting component", func() {
			comp := basic.DefaultGreetingComponent()
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", comp.URL(), nil)
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusOK))
			Expect(w.Header().Get("Content-Type")).To(ContainSubstring("text/html"))

			body, _ := io.ReadAll(w.Body)
			html := string(body)
			Expect(html).To(ContainSubstring("Hello, Alice!"))
			Expect(html).To(ContainSubstring("Welcome to Templ learning"))
		})

		It("should render usercard component", func() {
			comp := basic.DefaultUserCardComponent()
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", comp.URL(), nil)
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusOK))
			body, _ := io.ReadAll(w.Body)
			html := string(body)
			Expect(html).To(ContainSubstring("John"))
			Expect(html).To(ContainSubstring("Active"))
		})

		It("should render todolist component", func() {
			comp := basic.DefaultTodoListComponent()
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", comp.URL(), nil)
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusOK))
			body, _ := io.ReadAll(w.Body)
			html := string(body)
			Expect(html).To(ContainSubstring("Todo Items"))
			Expect(html).To(ContainSubstring("Learn Templ"))
		})

		It("should render button component", func() {
			comp := basic.DefaultButtonComponent()
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", comp.URL(), nil)
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusOK))
			body, _ := io.ReadAll(w.Body)
			html := string(body)
			Expect(html).To(ContainSubstring("ClickMe"))
			Expect(html).ToNot(ContainSubstring("disabled"))
		})

		It("should render emptylist component", func() {
			comp := basic.DefaultEmptyListComponent()
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", comp.URL(), nil)
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusOK))
			body, _ := io.ReadAll(w.Body)
			html := string(body)
			Expect(html).To(ContainSubstring("Todo Items"))
			Expect(html).To(ContainSubstring("<ul>"))
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
			Expect(registry.Basic()).To(HaveLen(5))
		})

		It("should have correct number of medium components", func() {
			Expect(registry.Medium()).To(HaveLen(6))
		})

		It("should have correct number of advanced components", func() {
			Expect(registry.Advanced()).To(HaveLen(2))
		})

		It("should find component by URL", func() {
			comp := registry.FindByURL("/basic/greeting/Alice")
			Expect(comp).ToNot(BeNil())
			Expect(comp.Name()).To(Equal("greeting"))
		})

		It("should return nil for unknown URL", func() {
			comp := registry.FindByURL("/unknown/path")
			Expect(comp).To(BeNil())
		})
	})
})
