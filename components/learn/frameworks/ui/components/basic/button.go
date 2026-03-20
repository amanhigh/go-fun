package basic

import (
	"github.com/a-h/templ"
	"github.com/amanhigh/go-fun/components/learn/frameworks/ui/components"
)

// ButtonComponent demonstrates variants, sizes, and states (FR-001 1.1)
type ButtonComponent struct {
	*components.BaseComponent
	text     string
	variant  ButtonVariant
	size     ButtonSize
	disabled bool
}

var _ components.Component = (*ButtonComponent)(nil)

// Text returns the button text
func (c *ButtonComponent) Text() string { return c.text }

// Variant returns the button variant
func (c *ButtonComponent) Variant() ButtonVariant { return c.variant }

// Size returns the button size
func (c *ButtonComponent) Size() ButtonSize { return c.size }

// Disabled returns whether the button is disabled
func (c *ButtonComponent) Disabled() bool { return c.disabled }

// NewButtonComponent creates a new button component
func NewButtonComponent(text string, variant ButtonVariant, size ButtonSize, disabled bool) *ButtonComponent {
	disabledStr := "enabled"
	if disabled {
		disabledStr = "disabled"
	}
	c := &ButtonComponent{text: text, variant: variant, size: size, disabled: disabled}
	c.BaseComponent = components.NewBaseComponent(
		"button",
		"Button with variants, sizes, and states (primary/secondary/tertiary, S/M/L, hover/pressed/disabled/focus)",
		"/basic/button/"+text+"/"+string(variant)+"/"+string(size)+"/"+disabledStr,
		components.LevelBasic,
		4,
		c.render,
	)
	return c
}

func (c *ButtonComponent) render() templ.Component {
	return Button(c.text, c.variant, c.size, c.disabled)
}

// DefaultButtonComponent returns the default button component for demo
func DefaultButtonComponent() *ButtonComponent {
	return NewButtonComponent("ClickMe", ButtonVariantPrimary, ButtonSizeMedium, false)
}
