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

// Registry holds all registered components organized by level
type Registry struct {
	basic    []Component
	medium   []Component
	advanced []Component
}

// NewRegistry creates a new component registry
func NewRegistry() *Registry {
	return &Registry{
		basic:    make([]Component, 0),
		medium:   make([]Component, 0),
		advanced: make([]Component, 0),
	}
}

// Register adds a component to the registry
func (r *Registry) Register(c Component) {
	switch c.Level() {
	case LevelBasic:
		r.basic = append(r.basic, c)
	case LevelMedium:
		r.medium = append(r.medium, c)
	case LevelAdvanced:
		r.advanced = append(r.advanced, c)
	}
}

// Basic returns all basic level components sorted by order
func (r *Registry) Basic() []Component {
	return sortByOrder(r.basic)
}

// Medium returns all medium level components sorted by order
func (r *Registry) Medium() []Component {
	return sortByOrder(r.medium)
}

// Advanced returns all advanced level components sorted by order
func (r *Registry) Advanced() []Component {
	return sortByOrder(r.advanced)
}

// All returns all components across all levels
func (r *Registry) All() []Component {
	all := make([]Component, 0, len(r.basic)+len(r.medium)+len(r.advanced))
	all = append(all, r.Basic()...)
	all = append(all, r.Medium()...)
	all = append(all, r.Advanced()...)
	return all
}

// FindByURL finds a component by its URL path
func (r *Registry) FindByURL(url string) Component {
	for _, c := range r.All() {
		if c.URL() == url {
			return c
		}
	}
	return nil
}

// sortByOrder sorts components by their order value
func sortByOrder(components []Component) []Component {
	sorted := make([]Component, len(components))
	copy(sorted, components)

	for i := 0; i < len(sorted)-1; i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[i].Order() > sorted[j].Order() {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}
	return sorted
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
