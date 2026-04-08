package main

import (
	"net/http"
	"path/filepath"
	"runtime"

	"github.com/amanhigh/go-fun/components/learn/frameworks/frontend/ui/components"
	"github.com/amanhigh/go-fun/components/learn/frameworks/frontend/ui/pages"
	"github.com/gin-gonic/gin"
	"github.com/templui/templui/utils"
)

// CreateComponents creates the standard set of UI components
func CreateComponents() []components.Component {
	return []components.Component{
		pages.NewHelloComponent(),
		pages.NewFormShowcaseComponent(),
		pages.NewStudentListComponent(),
	}
}

// SetupRoutes configures all routes on the given gin engine with provided components
func SetupRoutes(r *gin.Engine, components []components.Component) {
	r.Static("/assets", assetsDir())

	// Serve templui JavaScript files using embedded assets
	mux := http.NewServeMux()
	utils.SetupScriptRoutes(mux, true) // true for development
	r.Any("/templui/*filepath", gin.WrapH(mux))

	// Initialize and register student API routes
	studentHandler := NewStudentHandler()
	studentHandler.RegisterRoutes(r)

	// Index page - shows all available components
	r.GET("/", func(c *gin.Context) {
		c.Header("Content-Type", "text/html")
		pages.IndexPage(components).Render(c.Request.Context(), c.Writer)
	})

	// Register component routes
	for _, comp := range components {
		r.GET(comp.URL(), func(c *gin.Context) {
			c.Header("Content-Type", "text/html")
			comp.Render().Render(c.Request.Context(), c.Writer)
		})
	}
}

func assetsDir() string {
	_, currentFile, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(currentFile), "..", "assets")
}

// UIServer holds the HTTP server configuration and components
type UIServer struct {
	port       string
	components []components.Component
}

// NewUIServer creates a new UI server instance
func NewUIServer(port string) *UIServer {
	return &UIServer{
		port:       port,
		components: CreateComponents(),
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
	SetupRoutes(r, s.components)
}

// Components returns the server's components for testing
func (s *UIServer) Components() []components.Component {
	return s.components
}
