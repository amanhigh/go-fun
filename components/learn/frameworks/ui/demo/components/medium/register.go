package medium

import (
	"github.com/amanhigh/go-fun/components/learn/frameworks/ui/demo/components"
)

// RegisterAll registers all medium components with the given registry
func RegisterAll(r *components.Registry) {
	r.Register(DefaultNestedComponent())
	r.Register(DefaultCounterComponent())
	r.Register(DefaultDataTableComponent())
	r.Register(DefaultComposedComponent())
	r.Register(DefaultXSSComponent())
	r.Register(DefaultEmptyTableComponent())
}
