package pages

import (
	"github.com/a-h/templ"
	"github.com/amanhigh/go-fun/components/learn/frameworks/ui/components"
)

// BasicShowcaseComponent renders the basic showcase page.
type BasicShowcaseComponent struct {
	*components.BaseComponent
}

// NewBasicShowcaseComponent creates the basic component showcase.
func NewBasicShowcaseComponent() *BasicShowcaseComponent {
	c := &BasicShowcaseComponent{}
	c.BaseComponent = components.NewBaseComponent(
		"basic-showcase",
		"Single-page showcase for basic components",
		"/basic/showcase",
		components.LevelBasic,
		1,
		c.render,
	)
	return c
}

var _ components.Component = (*BasicShowcaseComponent)(nil)

func (c *BasicShowcaseComponent) render() templ.Component {
	return BasicShowcasePage()
}

// RegisterBasic registers all basic components with the given registry.
func RegisterBasic(r *components.Registry) {
	r.Register(NewBasicShowcaseComponent())
}
