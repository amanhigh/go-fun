package medium

import (
	"github.com/a-h/templ"
	"github.com/amanhigh/go-fun/components/learn/frameworks/ui"
	"github.com/amanhigh/go-fun/components/learn/frameworks/ui/demo/components"
)

// EmptyTableComponent demonstrates graceful handling of empty table data
type EmptyTableComponent struct {
	*components.BaseComponent
}

var _ components.Component = (*EmptyTableComponent)(nil)

// NewEmptyTableComponent creates a new empty table component
func NewEmptyTableComponent() *EmptyTableComponent {
	c := &EmptyTableComponent{}
	c.BaseComponent = components.NewBaseComponent(
		"emptytable",
		"Empty table demonstrating graceful handling of no data rows",
		"/medium/emptytable",
		components.LevelMedium,
		6,
		c.render,
	)
	return c
}

func (c *EmptyTableComponent) render() templ.Component {
	return ui.DataTable([]ui.TableRow{})
}

// DefaultEmptyTableComponent returns the default empty table component for demo
func DefaultEmptyTableComponent() *EmptyTableComponent {
	return NewEmptyTableComponent()
}
