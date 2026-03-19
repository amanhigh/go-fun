package basic

import (
	"context"
	"io"

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
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		disabledAttr := ""
		if c.disabled {
			disabledAttr = " disabled"
		}
		_, _ = io.WriteString(w, `<button type="button" class="btn"`+disabledAttr+`>`+c.text+`</button>`)
		return nil
	})
}

// DefaultButtonComponent returns the default button component for demo
func DefaultButtonComponent() *ButtonComponent {
	return NewButtonComponent("ClickMe", false)
}
