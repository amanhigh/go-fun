package pages

import (
	"github.com/a-h/templ"
	"github.com/amanhigh/go-fun/components/learn/frameworks/frontend/ui/components"
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
