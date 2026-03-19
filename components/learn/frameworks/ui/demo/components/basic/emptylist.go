package basic

import (
	"github.com/a-h/templ"
	"github.com/amanhigh/go-fun/components/learn/frameworks/ui"
	"github.com/amanhigh/go-fun/components/learn/frameworks/ui/demo/components"
)

// EmptyListComponent demonstrates graceful handling of empty lists
type EmptyListComponent struct {
	*components.BaseComponent
}

var _ components.Component = (*EmptyListComponent)(nil)

// NewEmptyListComponent creates a new empty list component
func NewEmptyListComponent() *EmptyListComponent {
	c := &EmptyListComponent{}
	c.BaseComponent = components.NewBaseComponent(
		"emptylist",
		"Empty list demonstrating graceful handling of no items",
		"/basic/emptylist",
		components.LevelBasic,
		5,
		c.render,
	)
	return c
}

func (c *EmptyListComponent) render() templ.Component {
	return ui.TodoList([]string{})
}

// DefaultEmptyListComponent returns the default empty list component for demo
func DefaultEmptyListComponent() *EmptyListComponent {
	return NewEmptyListComponent()
}
