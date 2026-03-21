package pages

import (
	"github.com/a-h/templ"
	"github.com/amanhigh/go-fun/components/learn/frameworks/ui/components"
)

// LayoutShowcaseComponent renders the layout & composition showcase page.
type LayoutShowcaseComponent struct {
	components.BaseComponent
}

// NewLayoutShowcaseComponent creates the layout & composition showcase.
func NewLayoutShowcaseComponent() *LayoutShowcaseComponent {
	return &LayoutShowcaseComponent{
		BaseComponent: components.NewBaseComponent(
			"layout-showcase",
			"🎨 Layout & Composition - Complex page layouts and responsive design patterns",
			"/layout",
			components.LevelAdvanced,
			2,
		),
	}
}

func (c *LayoutShowcaseComponent) Render() templ.Component {
	return LayoutShowcasePage()
}

// RegisterLayout registers the layout & composition showcase component with the given registry.
func RegisterLayout(r *components.Registry) {
	r.Register(NewLayoutShowcaseComponent())
}
