package basic

import (
	"github.com/a-h/templ"
	"github.com/amanhigh/go-fun/components/learn/frameworks/ui/components"
)

// ButtonComponent demonstrates attribute handling with disabled state
type ButtonComponent struct {
	*components.BaseComponent
	text     string
	disabled bool
}

var _ components.Component = (*ButtonComponent)(nil)

// NewButtonComponent creates a new button component
func NewButtonComponent(text string, disabled bool) *ButtonComponent {
	disabledStr := "enabled"
	if disabled {
		disabledStr = "disabled"
	}
	c := &ButtonComponent{text: text, disabled: disabled}
	c.BaseComponent = components.NewBaseComponent(
		"button",
		"Button with disabled state attribute handling",
		"/basic/button/"+text+"/"+disabledStr,
		components.LevelBasic,
		4,
		c.render,
	)
	return c
}

func (c *ButtonComponent) render() templ.Component {
	return Button(c.text, c.disabled)
}

// DefaultButtonComponent returns the default button component for demo
func DefaultButtonComponent() *ButtonComponent {
	return NewButtonComponent("ClickMe", false)
}
