package main

import (
	"net/http"

	"github.com/amanhigh/go-fun/components/learn/frameworks/ui/demo/components"
	"github.com/amanhigh/go-fun/components/learn/frameworks/ui/demo/components/advanced"
	"github.com/amanhigh/go-fun/components/learn/frameworks/ui/demo/components/basic"
	"github.com/amanhigh/go-fun/components/learn/frameworks/ui/demo/components/medium"
	"github.com/amanhigh/go-fun/components/learn/frameworks/ui/demo/pages"
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
	basic.RegisterAll(registry)
	medium.RegisterAll(registry)
	advanced.RegisterAll(registry)

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
	// Index page
	r.GET("/", s.indexHandler)

	// Level pages
	r.GET("/basic", s.basicPageHandler)
	r.GET("/medium", s.mediumPageHandler)
	r.GET("/advanced", s.advancedPageHandler)

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

// indexHandler serves the main index page
func (s *UIServer) indexHandler(c *gin.Context) {
	c.Header("Content-Type", "text/html")
	pages.IndexPage().Render(c.Request.Context(), c.Writer)
}

// basicPageHandler serves the basic components page
func (s *UIServer) basicPageHandler(c *gin.Context) {
	comps := pages.ComponentsToInfoList(s.registry.Basic())
	c.Header("Content-Type", "text/html")
	pages.LevelPage("basic", "Basic Components", comps).Render(c.Request.Context(), c.Writer)
}

// mediumPageHandler serves the medium components page
func (s *UIServer) mediumPageHandler(c *gin.Context) {
	comps := pages.ComponentsToInfoList(s.registry.Medium())
	c.Header("Content-Type", "text/html")
	pages.LevelPage("medium", "Medium Components", comps).Render(c.Request.Context(), c.Writer)
}

// advancedPageHandler serves the advanced components page
func (s *UIServer) advancedPageHandler(c *gin.Context) {
	comps := pages.ComponentsToInfoList(s.registry.Advanced())
	c.Header("Content-Type", "text/html")
	pages.LevelPage("advanced", "Advanced Components", comps).Render(c.Request.Context(), c.Writer)
}

// GetComponent returns a component by URL for testing
func (s *UIServer) GetComponent(url string) components.Component {
	return s.registry.FindByURL(url)
}

// Ensure UIServer is not used directly
var _ = http.StatusOK
