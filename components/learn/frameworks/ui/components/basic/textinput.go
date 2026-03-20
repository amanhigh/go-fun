package basic

import (
	"github.com/a-h/templ"
	"github.com/amanhigh/go-fun/components/learn/frameworks/ui/components"
)

// TextInputConfig holds configuration for text input components
type TextInputConfig struct {
	ID           string
	Label        string
	Placeholder  string
	Value        string
	State        InputState
	HelperText   string
	ErrorMessage string
}

// TextInputComponent demonstrates label, helper, and error treatments (FR-001 1.2)
type TextInputComponent struct {
	*components.BaseComponent
	config TextInputConfig
}

var _ components.Component = (*TextInputComponent)(nil)

// NewTextInputComponent creates a new text input component
func NewTextInputComponent(config TextInputConfig) *TextInputComponent {
	c := &TextInputComponent{config: config}
	c.BaseComponent = components.NewBaseComponent(
		"textinput",
		"Text input with label, helper, and error treatments",
		"/basic/textinput/"+config.ID,
		components.LevelBasic,
		5,
		c.render,
	)
	return c
}

func (c *TextInputComponent) render() templ.Component {
	return TextInput(c.config.ID, c.config.Label, c.config.Placeholder, c.config.Value, c.config.State, c.config.HelperText, c.config.ErrorMessage)
}

// DefaultTextInputComponent returns the default text input component for demo
func DefaultTextInputComponent() *TextInputComponent {
	return NewTextInputComponent(TextInputConfig{
		ID:           "username",
		Label:        "Username",
		Placeholder:  "Enter your username",
		Value:        "",
		State:        InputStateDefault,
		HelperText:   "Username should be 3-20 characters",
		ErrorMessage: "",
	})
}

// ErrorTextInputComponent returns a text input component with error state for demo
func ErrorTextInputComponent() *TextInputComponent {
	return NewTextInputComponent(TextInputConfig{
		ID:           "email",
		Label:        "Email",
		Placeholder:  "Enter your email",
		Value:        "invalid-email",
		State:        InputStateError,
		HelperText:   "",
		ErrorMessage: "Please enter a valid email address",
	})
}

// SuccessTextInputComponent returns a text input component with success state for demo
func SuccessTextInputComponent() *TextInputComponent {
	return NewTextInputComponent(TextInputConfig{
		ID:           "zipcode",
		Label:        "Zip Code",
		Placeholder:  "Enter zip code",
		Value:        "12345",
		State:        InputStateSuccess,
		HelperText:   "Valid zip code format",
		ErrorMessage: "",
	})
}
