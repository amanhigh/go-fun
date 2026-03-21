package pages

import (
	"github.com/a-h/templ"
	"github.com/amanhigh/go-fun/components/learn/frameworks/ui/components"
)

// AdvancedShowcaseComponent renders the advanced showcase page.
type AdvancedShowcaseComponent struct {
	*components.BaseComponent
}

// NewAdvancedShowcaseComponent creates the advanced component showcase.
func NewAdvancedShowcaseComponent() *AdvancedShowcaseComponent {
	c := &AdvancedShowcaseComponent{}
	c.BaseComponent = components.NewBaseComponent(
		"showcase",
		"Single-page showcase for advanced components",
		"/advanced/showcase",
		components.LevelAdvanced,
		1,
		c.render,
	)
	return c
}

func (c *AdvancedShowcaseComponent) render() templ.Component {
	return AdvancedShowcasePage()
}

// RegisterAdvanced registers the advanced showcase component with the given registry.
func RegisterAdvanced(r *components.Registry) {
	r.Register(NewAdvancedShowcaseComponent())
}
