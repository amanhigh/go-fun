package advanced

import (
	"github.com/a-h/templ"
	"github.com/amanhigh/go-fun/components/learn/frameworks/ui/components"
	"github.com/templui/templui/components/textarea"
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
	content := textarea.Textarea(textarea.Props{
		ID:          "advanced-notes",
		Name:        "advanced-notes",
		Placeholder: "Capture advanced page notes",
		Rows:        4,
	})
	return Page("Advanced Full Page Demo", content)
}

// DefaultFullPageComponent returns the default full page component for demo
func DefaultFullPageComponent() *FullPageComponent {
	return NewFullPageComponent()
}
