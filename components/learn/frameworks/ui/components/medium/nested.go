package medium

import (
	"github.com/a-h/templ"
	"github.com/amanhigh/go-fun/components/learn/frameworks/ui"
	"github.com/amanhigh/go-fun/components/learn/frameworks/ui/components"
)

// NestedComponent demonstrates nested component composition
type NestedComponent struct {
	*components.BaseComponent
	title string
	name  string
}

var _ components.Component = (*NestedComponent)(nil)

// NewNestedComponent creates a new nested component
func NewNestedComponent(title, name string) *NestedComponent {
	c := &NestedComponent{title: title, name: name}
	c.BaseComponent = components.NewBaseComponent(
		"nested",
		"Nested page component with greeting inside a page wrapper",
		"/medium/nested",
		components.LevelMedium,
		1,
		c.render,
	)
	return c
}

func (c *NestedComponent) render() templ.Component {
	greeting := ui.Greeting(c.name)
	return ui.Page(c.title, greeting)
}

// DefaultNestedComponent returns the default nested component for demo
func DefaultNestedComponent() *NestedComponent {
	return NewNestedComponent("Welcome Page", "Bob")
}
