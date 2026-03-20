package medium

import (
	"context"
	"io"

	"github.com/a-h/templ"
	"github.com/amanhigh/go-fun/components/learn/frameworks/ui/components"
	"github.com/amanhigh/go-fun/components/learn/frameworks/ui/components/advanced"
	"github.com/templui/templui/components/input"
	"github.com/templui/templui/components/textarea"
)

// ComposedComponent demonstrates multiple components composed together
type ComposedComponent struct {
	*components.BaseComponent
	name  string
	todos []string
}

var _ components.Component = (*ComposedComponent)(nil)

// NewComposedComponent creates a new composed component
func NewComposedComponent(name string, todos []string) *ComposedComponent {
	c := &ComposedComponent{name: name, todos: todos}
	c.BaseComponent = components.NewBaseComponent(
		"composed",
		"Multiple components composed together (greeting + todo list)",
		"/medium/composed",
		components.LevelMedium,
		4,
		c.render,
	)
	return c
}

func (c *ComposedComponent) render() templ.Component {
	profile := input.Input(input.Props{
		ID:          "composed-team",
		Name:        "composed-team",
		Placeholder: c.name,
		Value:       c.name,
	})

	notes := textarea.Textarea(textarea.Props{
		ID:          "composed-notes",
		Name:        "composed-notes",
		Placeholder: "Tasks",
		Value:       "Review code\nDeploy app\nWrite docs",
		Rows:        4,
	})

	return advanced.Page("Composed View", composedContent(profile, notes))
}

// composedContent wraps multiple components in a single templ component
func composedContent(greeting, todoList templ.Component) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		if err := greeting.Render(ctx, w); err != nil {
			return err
		}
		return todoList.Render(ctx, w)
	})
}

// DefaultComposedComponent returns the default composed component for demo
func DefaultComposedComponent() *ComposedComponent {
	return NewComposedComponent("Team", []string{"Review code", "Deploy app", "Write docs"})
}
