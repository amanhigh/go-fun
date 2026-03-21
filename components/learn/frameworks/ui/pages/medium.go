package pages

import (
	"github.com/a-h/templ"
	"github.com/amanhigh/go-fun/components/learn/frameworks/ui/components"
)

// MediumShowcaseComponent renders the medium showcase page.
type MediumShowcaseComponent struct {
	*components.BaseComponent
}

// NewMediumShowcaseComponent creates the medium component showcase.
func NewMediumShowcaseComponent() *MediumShowcaseComponent {
	c := &MediumShowcaseComponent{}
	c.BaseComponent = components.NewBaseComponent(
		"medium-showcase",
		"Single-page showcase for medium components",
		"/medium/showcase",
		components.LevelMedium,
		1,
		c.render,
	)
	return c
}

var _ components.Component = (*MediumShowcaseComponent)(nil)

func (c *MediumShowcaseComponent) render() templ.Component {
	return MediumShowcasePage()
}

// RegisterMedium registers all medium components with the given registry
func RegisterMedium(r *components.Registry) {
	r.Register(NewMediumShowcaseComponent())
}
