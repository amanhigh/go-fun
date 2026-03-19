package advanced

import (
	"github.com/a-h/templ"
	"github.com/amanhigh/go-fun/components/learn/frameworks/ui"
	"github.com/amanhigh/go-fun/components/learn/frameworks/ui/demo/components"
)

// FullPageComponent demonstrates a complete page with all component types
type FullPageComponent struct {
	*components.BaseComponent
}

var _ components.Component = (*FullPageComponent)(nil)

// NewFullPageComponent creates a new full page component
func NewFullPageComponent() *FullPageComponent {
	c := &FullPageComponent{}
	c.BaseComponent = components.NewBaseComponent(
		"fullpage",
		"Complete page demonstrating all component types together",
		"/advanced/fullpage",
		components.LevelAdvanced,
		1,
		c.render,
	)
	return c
}

func (c *FullPageComponent) render() templ.Component {
	// Create a complex page with multiple nested components
	greeting := ui.Greeting("Advanced User")
	return ui.Page("Advanced Full Page Demo", greeting)
}

// DefaultFullPageComponent returns the default full page component for demo
func DefaultFullPageComponent() *FullPageComponent {
	return NewFullPageComponent()
}
