package medium

import (
	"context"
	"io"

	"github.com/a-h/templ"
	"github.com/amanhigh/go-fun/components/learn/frameworks/ui/components"
	"github.com/amanhigh/go-fun/components/learn/frameworks/ui/components/advanced"
	"github.com/amanhigh/go-fun/components/learn/frameworks/ui/components/basic"
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
	greeting := basic.Greeting(c.name)
	todoList := basic.TodoList(c.todos)
	return advanced.Page("Composed View", composedContent(greeting, todoList))
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
