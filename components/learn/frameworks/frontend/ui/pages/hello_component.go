package pages

import (
	"github.com/a-h/templ"
	"github.com/amanhigh/go-fun/components/learn/frameworks/frontend/ui/components"
)

type HelloComponent struct {
	components.BaseComponent
}

func NewHelloComponent() *HelloComponent {
	return &HelloComponent{
		components.NewBaseComponent(
			"hello",
			"Hello - Single component test",
			"/hello",
			components.LevelBasic,
			0,
		),
	}
}

var _ components.Component = (*HelloComponent)(nil)

func (c *HelloComponent) Render() templ.Component {
	return HelloPage()
}
