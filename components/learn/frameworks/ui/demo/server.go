package main

import (
	"net/http"

	"github.com/amanhigh/go-fun/components/learn/frameworks/ui"
	"github.com/gin-gonic/gin"
)

// UIServer provides HTTP endpoints to view Templ components in browser
type UIServer struct {
	port string
}

// NewUIServer creates a new UI server instance
func NewUIServer(port string) *UIServer {
	return &UIServer{
		port: port,
	}
}

// Start starts the HTTP server and serves the UI demo pages
func (s *UIServer) Start() error {
	r := gin.Default()

	// Main demo page
	r.GET("/", s.indexHandler)

	// Individual component demos
	r.GET("/greeting/:name", s.greetingHandler)
	r.GET("/usercard/:username/:active", s.userCardHandler)
	r.GET("/todolist", s.todoListHandler)
	r.GET("/button/:text/:disabled", s.buttonHandler)
	r.GET("/counter/:count", s.counterHandler)

	return r.Run(":" + s.port)
}

// indexHandler serves the main demo page with links to all examples
func (s *UIServer) indexHandler(c *gin.Context) {
	html := `<!DOCTYPE html>
<html>
<head>
    <title>Templ UI Demo</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        .demo-section { margin: 30px 0; padding: 20px; border: 1px solid #ddd; border-radius: 8px; }
        .demo-link { color: #007bff; text-decoration: none; font-weight: bold; }
        .demo-link:hover { text-decoration: underline; }
        .back-link { margin-top: 20px; display: block; }
    </style>
</head>
<body>
    <h1>Templ UI Component Demo</h1>
    <p>This demo showcases the first 5 basic Templ components. Click on any component to see it rendered:</p>
    
    <div class="demo-section">
        <h2>1. Greeting Component</h2>
        <p>A simple greeting component with a name parameter.</p>
        <a href="/greeting/Alice" class="demo-link">View Greeting: "Alice"</a> | 
        <a href="/greeting/Bob" class="demo-link">View Greeting: "Bob"</a>
    </div>

    <div class="demo-section">
        <h2>2. UserCard Component</h2>
        <p>A user card with conditional active/inactive badge rendering.</p>
        <a href="/usercard/John/true" class="demo-link">View UserCard: John (Active)</a> | 
        <a href="/usercard/Jane/false" class="demo-link">View UserCard: Jane (Inactive)</a>
    </div>

    <div class="demo-section">
        <h2>3. TodoList Component</h2>
        <p>A todo list that renders items in a loop.</p>
        <a href="/todolist" class="demo-link">View TodoList</a>
    </div>

    <div class="demo-section">
        <h2>4. Button Component</h2>
        <p>A button with disabled state attribute handling.</p>
        <a href="/button/ClickMe/false" class="demo-link">View Button: Enabled</a> | 
        <a href="/button/Disabled/true" class="demo-link">View Button: Disabled</a>
    </div>

    <div class="demo-section">
        <h2>5. Counter Component</h2>
        <p>A counter with different state rendering (positive, negative, zero).</p>
        <a href="/counter/0" class="demo-link">View Counter: 0</a> | 
        <a href="/counter/5" class="demo-link">View Counter: 5</a> | 
        <a href="/counter/-3" class="demo-link">View Counter: -3</a>
    </div>
</body>
</html>`

	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, html)
}

// greetingHandler renders the Greeting component
func (s *UIServer) greetingHandler(c *gin.Context) {
	name := c.Param("name")
	component := ui.Greeting(name)

	c.Header("Content-Type", "text/html")
	component.Render(c.Request.Context(), c.Writer)
}

// userCardHandler renders the UserCard component
func (s *UIServer) userCardHandler(c *gin.Context) {
	username := c.Param("username")
	activeStr := c.Param("active")
	isActive := activeStr == "true"

	component := ui.UserCard(username, isActive)

	c.Header("Content-Type", "text/html")
	component.Render(c.Request.Context(), c.Writer)
}

// todoListHandler renders the TodoList component
func (s *UIServer) todoListHandler(c *gin.Context) {
	todos := []string{"Learn Templ", "Build UI", "Test Components", "Write Documentation"}
	component := ui.TodoList(todos)

	c.Header("Content-Type", "text/html")
	component.Render(c.Request.Context(), c.Writer)
}

// buttonHandler renders the Button component
func (s *UIServer) buttonHandler(c *gin.Context) {
	text := c.Param("text")
	disabledStr := c.Param("disabled")
	isDisabled := disabledStr == "true"

	component := ui.Button(text, isDisabled)

	c.Header("Content-Type", "text/html")
	component.Render(c.Request.Context(), c.Writer)
}

// counterHandler renders the Counter component
func (s *UIServer) counterHandler(c *gin.Context) {
	countStr := c.Param("count")
	// For simplicity, we'll use a fixed count. In a real app, you'd parse this properly
	var count int
	switch countStr {
	case "0":
		count = 0
	case "5":
		count = 5
	case "-3":
		count = -3
	default:
		count = 0
	}

	component := ui.Counter(count)

	c.Header("Content-Type", "text/html")
	component.Render(c.Request.Context(), c.Writer)
}
