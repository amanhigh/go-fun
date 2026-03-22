package main

import (
	"net/http"

	"github.com/amanhigh/go-fun/components/learn/frameworks/frontend/ui/components"
	"github.com/amanhigh/go-fun/components/learn/frameworks/frontend/ui/pages"
	"github.com/gin-gonic/gin"
	"github.com/templui/templui/utils"
)

// UIServer holds the HTTP server configuration and components
type UIServer struct {
	port       string
	components []components.Component
}

// NewUIServer creates a new UI server instance
func NewUIServer(port string) *UIServer {
	// Create components once
	components := []components.Component{
		pages.NewFormShowcaseComponent(),
		pages.NewHelloComponent(),
	}

	return &UIServer{
		port:       port,
		components: components,
	}
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

	// Register component routes using server components
	for _, comp := range s.components {
		comp := comp // capture for closure
		r.GET(comp.URL(), func(c *gin.Context) {
			c.Header("Content-Type", "text/html")
			comp.Render().Render(c.Request.Context(), c.Writer)
		})
	}
}

// indexHandler serves the main index page with links to all showcases
func (s *UIServer) indexHandler(c *gin.Context) {
	c.Header("Content-Type", "text/html")

	// Use server components directly
	pages.IndexPage(s.components).Render(c.Request.Context(), c.Writer)
}

// Ensure UIServer is not used directly
var _ = http.StatusOK
