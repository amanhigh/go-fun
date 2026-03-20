package medium

import (
	"github.com/a-h/templ"
	"github.com/amanhigh/go-fun/components/learn/frameworks/ui/components"
	"github.com/amanhigh/go-fun/components/learn/frameworks/ui/components/advanced"
	"github.com/templui/templui/components/input"
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
	content := input.Input(input.Props{
		ID:          "nested-name",
		Name:        "nested-name",
		Placeholder: "Hello, " + c.name + "!",
		Value:       c.name,
		Attributes: templ.Attributes{
			"title": "Nested composition using templUI input",
		},
	})
	return advanced.Page(c.title, content)
}

// DefaultNestedComponent returns the default nested component for demo
func DefaultNestedComponent() *NestedComponent {
	return NewNestedComponent("Welcome Page", "Bob")
}
