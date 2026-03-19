package medium

import (
	"github.com/a-h/templ"
	"github.com/amanhigh/go-fun/components/learn/frameworks/ui/components"
)

// DataTableComponent demonstrates table rendering with structured data
type DataTableComponent struct {
	*components.BaseComponent
	rows []TableRow
}

var _ components.Component = (*DataTableComponent)(nil)

// NewDataTableComponent creates a new data table component
func NewDataTableComponent(rows []TableRow) *DataTableComponent {
	c := &DataTableComponent{rows: rows}
	c.BaseComponent = components.NewBaseComponent(
		"datatable",
		"Data table with multiple rows and columns",
		"/medium/datatable",
		components.LevelMedium,
		3,
		c.render,
	)
	return c
}

func (c *DataTableComponent) render() templ.Component {
	return DataTable(c.rows)
}

// DefaultDataTableComponent returns the default data table component for demo
func DefaultDataTableComponent() *DataTableComponent {
	rows := []TableRow{
		{ID: 1, Name: "Alice", Age: 25},
		{ID: 2, Name: "Bob", Age: 30},
		{ID: 3, Name: "Charlie", Age: 35},
	}
	return NewDataTableComponent(rows)
}
