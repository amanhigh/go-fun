package pages

import (
	"github.com/a-h/templ"
	"github.com/amanhigh/go-fun/components/learn/frameworks/ui/components"
)

// FormShowcaseComponent renders the form essentials showcase page.
type FormShowcaseComponent struct {
	components.BaseComponent
}

// NewFormShowcaseComponent creates the form essentials showcase.
func NewFormShowcaseComponent() *FormShowcaseComponent {
	return &FormShowcaseComponent{
		BaseComponent: components.NewBaseComponent(
			"form-showcase",
			"📝 Form Essentials - Master form inputs, validation, and user data collection patterns",
			"/form",
			components.LevelBasic,
			1,
		),
	}
}

var _ components.Component = (*FormShowcaseComponent)(nil)

func (c *FormShowcaseComponent) Render() templ.Component {
	return FormShowcasePage()
}

// RegisterBasic registers all form components with the given registry.
func RegisterBasic(r *components.Registry) {
	r.Register(NewHelloComponent())
	r.Register(NewFormShowcaseComponent())
}
