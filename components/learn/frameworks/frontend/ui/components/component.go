package components

import (
	"github.com/a-h/templ"
)

// Level represents the complexity level of a component
type Level string

const (
	LevelBasic    Level = "basic"
	LevelMedium   Level = "medium"
	LevelAdvanced Level = "advanced"
)

// Component defines the interface for all UI demo components
type Component interface {
	// Name returns the unique identifier for the component
	Name() string

	// Description returns a human-readable description of what the component demonstrates
	Description() string

	// URL returns the relative URL path for viewing this component
	URL() string

	// Level returns the complexity level (basic, medium, advanced)
	Level() Level

	// Order returns the display order within its level (lower numbers first)
	Order() int

	// Render returns the templ.Component for rendering
	Render() templ.Component
}

// BaseComponent provides base implementation with common fields that can be embedded
type BaseComponent struct {
	name        string
	description string
	url         string
	level       Level
	order       int
}

// NewBaseComponent creates a new base component with common fields
func NewBaseComponent(name, description, url string, level Level, order int) BaseComponent {
	return BaseComponent{
		name:        name,
		description: description,
		url:         url,
		level:       level,
		order:       order,
	}
}

func (b *BaseComponent) Name() string        { return b.name }
func (b *BaseComponent) Description() string { return b.description }
func (b *BaseComponent) URL() string         { return b.url }
func (b *BaseComponent) Level() Level        { return b.level }
func (b *BaseComponent) Order() int          { return b.order }
