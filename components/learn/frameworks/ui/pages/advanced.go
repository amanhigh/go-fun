package pages

import (
	"github.com/a-h/templ"
	"github.com/amanhigh/go-fun/components/learn/frameworks/ui/components"
)

// InteractiveShowcaseComponent renders the interactive behaviors showcase page.
type InteractiveShowcaseComponent struct {
	*components.BaseComponent
}

// NewInteractiveShowcaseComponent creates the interactive behaviors showcase.
func NewInteractiveShowcaseComponent() *InteractiveShowcaseComponent {
	c := &InteractiveShowcaseComponent{}
	c.BaseComponent = components.NewBaseComponent(
		"interactive-showcase",
		"⚡ Interactive Behaviors - Dynamic client-side interactions with Alpine.js",
		"/interactive",
		components.LevelAdvanced,
		1,
		c.render,
	)
	return c
}

func (c *InteractiveShowcaseComponent) render() templ.Component {
	return InteractiveShowcasePage()
}

// RegisterAdvanced registers the interactive behaviors showcase component with the given registry.
func RegisterAdvanced(r *components.Registry) {
	r.Register(NewInteractiveShowcaseComponent())
}
