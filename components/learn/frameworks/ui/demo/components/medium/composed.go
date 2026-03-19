package medium

import (
	"context"
	"io"

	"github.com/a-h/templ"
	"github.com/amanhigh/go-fun/components/learn/frameworks/ui"
	"github.com/amanhigh/go-fun/components/learn/frameworks/ui/demo/components"
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
	greeting := ui.Greeting(c.name)
	todoList := ui.TodoList(c.todos)
	return ui.Page("Composed View", composedContent(greeting, todoList))
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
