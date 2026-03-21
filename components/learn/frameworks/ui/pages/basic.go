package pages

import (
	"github.com/a-h/templ"
	"github.com/amanhigh/go-fun/components/learn/frameworks/ui/components"
)

// BasicShowcase defines the behavior for the basic showcase component.
type BasicShowcase interface {
	components.Component
}

// BasicShowcaseComponent renders the FR-001 basic showcase page.
type BasicShowcaseComponent struct {
	*components.BaseComponent
}

// NewBasicShowcaseComponent creates the unified basic component showcase.
func NewBasicShowcaseComponent() *BasicShowcaseComponent {
	c := &BasicShowcaseComponent{}
	c.BaseComponent = components.NewBaseComponent(
		"basic-showcase",
		"Single-page showcase for FR-001 basic components",
		"/basic/showcase",
		components.LevelBasic,
		1,
		c.render,
	)
	return c
}

var _ BasicShowcase = (*BasicShowcaseComponent)(nil)

func (c *BasicShowcaseComponent) render() templ.Component {
	return BasicShowcasePage()
}

// DefaultBasicShowcaseComponent returns the default FR-001 basic showcase.
func DefaultBasicShowcaseComponent() *BasicShowcaseComponent {
	return NewBasicShowcaseComponent()
}

// RegisterBasic registers all basic components with the given registry.
func RegisterBasic(r *components.Registry) {
	r.Register(DefaultBasicShowcaseComponent())
}
