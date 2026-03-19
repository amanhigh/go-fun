package advanced

import (
	"context"
	"io"

	"github.com/a-h/templ"
	"github.com/amanhigh/go-fun/components/learn/frameworks/ui"
	"github.com/amanhigh/go-fun/components/learn/frameworks/ui/demo/components"
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
	// Create dashboard content with multiple components
	rows := []ui.TableRow{
		{ID: 1, Name: "Project Alpha", Age: 90},
		{ID: 2, Name: "Project Beta", Age: 75},
		{ID: 3, Name: "Project Gamma", Age: 60},
	}

	return ui.Page("Dashboard", dashboardContent(rows))
}

// dashboardContent creates the dashboard layout with multiple widgets
func dashboardContent(rows []ui.TableRow) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		// Render header
		if err := ui.Greeting("Dashboard").Render(ctx, w); err != nil {
			return err
		}

		// Render user cards
		if err := ui.UserCard("Admin", true).Render(ctx, w); err != nil {
			return err
		}
		if err := ui.UserCard("Guest", false).Render(ctx, w); err != nil {
			return err
		}

		// Render counters
		if err := ui.Counter(42).Render(ctx, w); err != nil {
			return err
		}

		// Render data table
		return ui.DataTable(rows).Render(ctx, w)
	})
}

// DefaultDashboardComponent returns the default dashboard component for demo
func DefaultDashboardComponent() *DashboardComponent {
	return NewDashboardComponent()
}
