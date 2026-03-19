package basic

import (
	"github.com/a-h/templ"
	"github.com/amanhigh/go-fun/components/learn/frameworks/ui"
	"github.com/amanhigh/go-fun/components/learn/frameworks/ui/demo/components"
)

// GreetingComponent demonstrates a simple greeting with name parameter
type GreetingComponent struct {
	*components.BaseComponent
	name string
}

var _ components.Component = (*GreetingComponent)(nil)

// NewGreetingComponent creates a new greeting component with the given name
func NewGreetingComponent(name string) *GreetingComponent {
	c := &GreetingComponent{name: name}
	c.BaseComponent = components.NewBaseComponent(
		"greeting",
		"Simple greeting component with a name parameter",
		"/basic/greeting/"+name,
		components.LevelBasic,
		1,
		c.render,
	)
	return c
}

func (c *GreetingComponent) render() templ.Component {
	return ui.Greeting(c.name)
}

// DefaultGreetingComponent returns the default greeting component for demo
func DefaultGreetingComponent() *GreetingComponent {
	return NewGreetingComponent("Alice")
}
