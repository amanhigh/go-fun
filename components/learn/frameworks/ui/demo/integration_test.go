package main_test

import (
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/amanhigh/go-fun/components/learn/frameworks/ui/components"
	"github.com/amanhigh/go-fun/components/learn/frameworks/ui/components/advanced"
	"github.com/amanhigh/go-fun/components/learn/frameworks/ui/components/basic"
	"github.com/amanhigh/go-fun/components/learn/frameworks/ui/components/medium"
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
		It("should render unified FR-001 showcase", func() {
			comp := basic.DefaultBasicShowcaseComponent()
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", comp.URL(), nil)
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusOK))
			Expect(w.Header().Get("Content-Type")).To(ContainSubstring("text/html"))

			body, _ := io.ReadAll(w.Body)
			html := string(body)
			Expect(html).To(ContainSubstring("FR-001 Showcase"))
			Expect(html).To(ContainSubstring("Username"))
			Expect(html).To(ContainSubstring("Country"))
			Expect(html).To(ContainSubstring("showcase-modal"))
		})

		It("should support button variants (FR-001 1.1)", func() {
			primary := basic.NewButtonComponent("Submit", basic.ButtonVariantPrimary, basic.ButtonSizeMedium, false)
			Expect(primary.Variant()).To(Equal(basic.ButtonVariantPrimary))

			secondary := basic.NewButtonComponent("Cancel", basic.ButtonVariantSecondary, basic.ButtonSizeMedium, false)
			Expect(secondary.Variant()).To(Equal(basic.ButtonVariantSecondary))

			tertiary := basic.NewButtonComponent("Learn", basic.ButtonVariantTertiary, basic.ButtonSizeMedium, false)
			Expect(tertiary.Variant()).To(Equal(basic.ButtonVariantTertiary))
		})

		It("should support button sizes and states (FR-001 1.1)", func() {
			small := basic.NewButtonComponent("S", basic.ButtonVariantPrimary, basic.ButtonSizeSmall, false)
			Expect(small.Size()).To(Equal(basic.ButtonSizeSmall))

			large := basic.NewButtonComponent("L", basic.ButtonVariantPrimary, basic.ButtonSizeLarge, false)
			Expect(large.Size()).To(Equal(basic.ButtonSizeLarge))

			enabled := basic.NewButtonComponent("Enabled", basic.ButtonVariantPrimary, basic.ButtonSizeMedium, false)
			Expect(enabled.Disabled()).To(BeFalse())

			disabled := basic.NewButtonComponent("Disabled", basic.ButtonVariantPrimary, basic.ButtonSizeMedium, true)
			Expect(disabled.Disabled()).To(BeTrue())
		})

		It("should support text input states (FR-001 1.2)", func() {
			defaultInput := basic.DefaultTextInputComponent()
			Expect(defaultInput.State()).To(Equal(basic.InputStateDefault))

			errorInput := basic.ErrorTextInputComponent()
			Expect(errorInput.State()).To(Equal(basic.InputStateError))
			Expect(errorInput.ErrorMessage()).To(Equal("Please enter a valid email address"))

			successInput := basic.SuccessTextInputComponent()
			Expect(successInput.State()).To(Equal(basic.InputStateSuccess))
			Expect(successInput.Value()).To(Equal("12345"))
		})

		It("should implement Component interface for showcase component", func() {
			var _ components.Component = basic.DefaultBasicShowcaseComponent()

			showcase := basic.DefaultBasicShowcaseComponent()
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
