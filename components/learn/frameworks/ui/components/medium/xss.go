package medium

import (
	"github.com/a-h/templ"
	"github.com/amanhigh/go-fun/components/learn/frameworks/ui/components"
	"github.com/templui/templui/components/input"
)

// XSSComponent demonstrates HTML escaping for security
type XSSComponent struct {
	*components.BaseComponent
	name string
}

var _ components.Component = (*XSSComponent)(nil)

// NewXSSComponent creates a new XSS protection demo component
func NewXSSComponent(name string) *XSSComponent {
	c := &XSSComponent{name: name}
	c.BaseComponent = components.NewBaseComponent(
		"xss",
		"Special characters handling with HTML escaping for XSS protection",
		"/medium/xss",
		components.LevelMedium,
		5,
		c.render,
	)
	return c
}

func (c *XSSComponent) render() templ.Component {
	return input.Input(input.Props{
		ID:          "xss-name",
		Name:        "xss-name",
		Placeholder: "Unsafe Name",
		Value:       c.name,
	})
}

// DefaultXSSComponent returns the default XSS component for demo
func DefaultXSSComponent() *XSSComponent {
	return NewXSSComponent("<script>alert('xss')</script>")
}
