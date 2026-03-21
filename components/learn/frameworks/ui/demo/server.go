package main

import (
	"net/http"

	"github.com/amanhigh/go-fun/components/learn/frameworks/ui/components"
	"github.com/amanhigh/go-fun/components/learn/frameworks/ui/pages"
	"github.com/gin-gonic/gin"
)

// UIServer provides HTTP endpoints to view Templ components in browser
type UIServer struct {
	port     string
	registry *components.Registry
}

// NewUIServer creates a new UI server instance with all components registered
func NewUIServer(port string) *UIServer {
	registry := components.NewRegistry()

	// Register all components
	pages.RegisterBasic(registry)
	pages.RegisterMedium(registry)
	pages.RegisterAdvanced(registry)
	pages.RegisterLayout(registry)

	return &UIServer{
		port:     port,
		registry: registry,
	}
}

// Registry returns the component registry for testing
func (s *UIServer) Registry() *components.Registry {
	return s.registry
}

// Start starts the HTTP server and serves the UI demo pages
func (s *UIServer) Start() error {
	r := gin.Default()
	s.SetupRoutes(r)
	return r.Run(":" + s.port)
}

// SetupRoutes configures all routes on the given gin engine
func (s *UIServer) SetupRoutes(r *gin.Engine) {
	// Serve static files (JS, CSS, images) - path relative to demo directory
	r.Static("/static", "../static")

	// Index page
	r.GET("/", s.indexHandler)

	// Showcase pages (new feature-based names)
	r.GET("/form", s.formShowcaseHandler)
	r.GET("/data", s.dataShowcaseHandler)
	r.GET("/interactive", s.interactiveShowcaseHandler)
	r.GET("/layout", s.layoutShowcaseHandler)

	// Register all component routes dynamically
	for _, comp := range s.registry.All() {
		s.registerComponentRoute(r, comp)
	}
}

// registerComponentRoute registers a route for a single component
func (s *UIServer) registerComponentRoute(r *gin.Engine, comp components.Component) {
	url := comp.URL()
	r.GET(url, func(c *gin.Context) {
		c.Header("Content-Type", "text/html")
		comp.Render().Render(c.Request.Context(), c.Writer)
	})
}

// indexHandler serves the main index page with links to all showcases
func (s *UIServer) indexHandler(c *gin.Context) {
	c.Header("Content-Type", "text/html")

	levels := []pages.LevelInfo{
		{
			Name:        "📝 Form Essentials",
			Path:        "/form",
			Description: "Master form inputs, validation, and user data collection patterns. Text fields, dropdowns, checkboxes, and radio buttons with proper validation.",
			Count:       1,
			BadgeClass:  "badge-basic",
		},
		{
			Name:        "📊 Data Presentation",
			Path:        "/data",
			Description: "Display structured data with tables, cards, status indicators, and content organization patterns for dashboards and reports.",
			Count:       1,
			BadgeClass:  "badge-medium",
		},
		{
			Name:        "⚡ Interactive Behaviors",
			Path:        "/interactive",
			Description: "Dynamic client-side interactions with Alpine.js. Modals, character counters, real-time updates, and state management patterns.",
			Count:       1,
			BadgeClass:  "badge-advanced",
		},
		{
			Name:        "🎨 Layout & Composition",
			Path:        "/layout",
			Description: "Complex page layouts, grid systems, responsive design, and component composition for production-ready applications.",
			Count:       1,
			BadgeClass:  "badge-advanced",
		},
	}

	pages.IndexPage(levels).Render(c.Request.Context(), c.Writer)
}

// formShowcaseHandler serves the form essentials showcase
func (s *UIServer) formShowcaseHandler(c *gin.Context) {
	c.Header("Content-Type", "text/html")
	pages.FormShowcasePage().Render(c.Request.Context(), c.Writer)
}

// dataShowcaseHandler serves the data presentation showcase
func (s *UIServer) dataShowcaseHandler(c *gin.Context) {
	c.Header("Content-Type", "text/html")
	pages.DataShowcasePage().Render(c.Request.Context(), c.Writer)
}

// interactiveShowcaseHandler serves the interactive behaviors showcase
func (s *UIServer) interactiveShowcaseHandler(c *gin.Context) {
	c.Header("Content-Type", "text/html")
	pages.InteractiveShowcasePage().Render(c.Request.Context(), c.Writer)
}

// layoutShowcaseHandler serves the layout & composition showcase
func (s *UIServer) layoutShowcaseHandler(c *gin.Context) {
	c.Header("Content-Type", "text/html")
	pages.LayoutShowcasePage().Render(c.Request.Context(), c.Writer)
}

// GetComponent returns a component by URL for testing
func (s *UIServer) GetComponent(url string) components.Component {
	return s.registry.FindByURL(url)
}

// Ensure UIServer is not used directly
var _ = http.StatusOK
