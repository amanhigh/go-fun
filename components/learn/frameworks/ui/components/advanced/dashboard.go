package advanced

import (
	"context"
	"io"

	"github.com/a-h/templ"
	"github.com/amanhigh/go-fun/components/learn/frameworks/ui/components"
)

// DashboardComponent demonstrates a complex dashboard with multiple widgets
type DashboardComponent struct {
	*components.BaseComponent
}

var _ components.Component = (*DashboardComponent)(nil)

// NewDashboardComponent creates a new dashboard component
func NewDashboardComponent() *DashboardComponent {
	c := &DashboardComponent{}
	c.BaseComponent = components.NewBaseComponent(
		"dashboard",
		"Complex dashboard with multiple widgets (table, cards, counters)",
		"/advanced/dashboard",
		components.LevelAdvanced,
		2,
		c.render,
	)
	return c
}

func (c *DashboardComponent) render() templ.Component {
	return Page("Dashboard", dashboardContent())
}

// dashboardContent creates the dashboard layout with multiple widgets
func dashboardContent() templ.Component {
	return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		// Render header
		_, _ = io.WriteString(w, `<div class="greeting"><h1>Hello, Dashboard!</h1><p>Welcome to Templ learning.</p></div>`)

		// Render user cards
		_, _ = io.WriteString(w, `<div class="user-card"><h2>Admin</h2><span class="badge active">Active</span></div>`)
		_, _ = io.WriteString(w, `<div class="user-card"><h2>Guest</h2><span class="badge inactive">Inactive</span></div>`)

		// Render counter
		_, _ = io.WriteString(w, `<div class="counter"><h3>Counter Value</h3><div class="count-display">42</div><div class="count-info"><p>Counter is positive: 42</p></div></div>`)

		// Render data table
		_, _ = io.WriteString(w, `<table class="data-table"><thead><tr><th>ID</th><th>Name</th><th>Age</th></tr></thead><tbody>`)
		_, _ = io.WriteString(w, `<tr><td>1</td><td>Project Alpha</td><td>90</td></tr>`)
		_, _ = io.WriteString(w, `<tr><td>2</td><td>Project Beta</td><td>75</td></tr>`)
		_, _ = io.WriteString(w, `<tr><td>3</td><td>Project Gamma</td><td>60</td></tr>`)
		_, _ = io.WriteString(w, `</tbody></table>`)

		return nil
	})
}

// DefaultDashboardComponent returns the default dashboard component for demo
func DefaultDashboardComponent() *DashboardComponent {
	return NewDashboardComponent()
}
