package pages

import (
	"github.com/a-h/templ"
	"github.com/amanhigh/go-fun/components/learn/frameworks/ui/components"
)

// DataShowcaseComponent renders the data presentation showcase page.
type DataShowcaseComponent struct {
	components.BaseComponent
}

// NewDataShowcaseComponent creates the data presentation showcase.
func NewDataShowcaseComponent() *DataShowcaseComponent {
	return &DataShowcaseComponent{
		BaseComponent: components.NewBaseComponent(
			"data-showcase",
			"📊 Data Presentation - Display structured data with tables, cards, and status indicators",
			"/data",
			components.LevelMedium,
			1,
		),
	}
}

var _ components.Component = (*DataShowcaseComponent)(nil)

func (c *DataShowcaseComponent) Render() templ.Component {
	return DataShowcasePage()
}

// RegisterMedium registers all data presentation components with the given registry
func RegisterMedium(r *components.Registry) {
	r.Register(NewDataShowcaseComponent())
}
