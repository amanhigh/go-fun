package main

import (
	"net/http"

	"github.com/amanhigh/go-fun/components/learn/frameworks/frontend/ui/components"
	"github.com/amanhigh/go-fun/components/learn/frameworks/frontend/ui/pages"
	"github.com/gin-gonic/gin"
	"github.com/templui/templui/utils"
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
	r.Static("/assets", "../assets")

	// Serve templui JavaScript files using embedded assets
	mux := http.NewServeMux()
	utils.SetupScriptRoutes(mux, true) // true for development
	r.Any("/templui/*filepath", gin.WrapH(mux))

	// Index page - shows all available components
	r.GET("/", s.indexHandler)

	// Register all component routes dynamically using the registry
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

	// Build levels dynamically from registry
	levels := make([]pages.LevelInfo, 0, len(s.registry.All()))
	for _, comp := range s.registry.All() {
		badgeClass := "badge-basic"
		switch comp.Level() {
		case components.LevelBasic:
			badgeClass = "badge-basic"
		case components.LevelMedium:
			badgeClass = "badge-medium"
		case components.LevelAdvanced:
			badgeClass = "badge-advanced"
		}

		levels = append(levels, pages.LevelInfo{
			Name:        comp.Name(),
			Path:        comp.URL(),
			Description: comp.Description(),
			Count:       1,
			BadgeClass:  badgeClass,
		})
	}

	pages.IndexPage(levels).Render(c.Request.Context(), c.Writer)
}

// GetComponent returns a component by URL for testing
func (s *UIServer) GetComponent(url string) components.Component {
	return s.registry.FindByURL(url)
}

// Ensure UIServer is not used directly
var _ = http.StatusOK
