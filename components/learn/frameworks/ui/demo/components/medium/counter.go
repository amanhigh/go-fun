package medium

import (
	"github.com/a-h/templ"
	"github.com/amanhigh/go-fun/components/learn/frameworks/ui"
	"github.com/amanhigh/go-fun/components/learn/frameworks/ui/demo/components"
)

// CounterComponent demonstrates state-based conditional rendering
type CounterComponent struct {
	*components.BaseComponent
	count int
}

var _ components.Component = (*CounterComponent)(nil)

// NewCounterComponent creates a new counter component
func NewCounterComponent(count int) *CounterComponent {
	c := &CounterComponent{count: count}
	c.BaseComponent = components.NewBaseComponent(
		"counter",
		"Counter with different state rendering (positive, negative, zero)",
		"/medium/counter",
		components.LevelMedium,
		2,
		c.render,
	)
	return c
}

func (c *CounterComponent) render() templ.Component {
	return ui.Counter(c.count)
}

// DefaultCounterComponent returns the default counter component for demo
func DefaultCounterComponent() *CounterComponent {
	return NewCounterComponent(5)
}
