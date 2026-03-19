package basic

import (
	"context"
	"io"

	"github.com/a-h/templ"
	"github.com/amanhigh/go-fun/components/learn/frameworks/ui/components"
)

// TodoListComponent demonstrates loop rendering with a list of items
type TodoListComponent struct {
	*components.BaseComponent
	todos []string
}

var _ components.Component = (*TodoListComponent)(nil)

// NewTodoListComponent creates a new todo list component
func NewTodoListComponent(todos []string) *TodoListComponent {
	c := &TodoListComponent{todos: todos}
	c.BaseComponent = components.NewBaseComponent(
		"todolist",
		"Todo list that renders items in a loop",
		"/basic/todolist",
		components.LevelBasic,
		3,
		c.render,
	)
	return c
}

func (c *TodoListComponent) render() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, _ = io.WriteString(w, `<div class="todo-list"><h3>Todo Items</h3><ul>`)
		for _, todo := range c.todos {
			_, _ = io.WriteString(w, `<li>`+todo+`</li>`)
		}
		_, _ = io.WriteString(w, `</ul></div>`)
		return nil
	})
}

// DefaultTodoListComponent returns the default todo list component for demo
func DefaultTodoListComponent() *TodoListComponent {
	return NewTodoListComponent([]string{"Learn Templ", "Build UI", "Test Components", "Write Documentation"})
}
