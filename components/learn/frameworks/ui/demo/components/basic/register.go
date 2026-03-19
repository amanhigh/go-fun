package basic

import (
	"github.com/amanhigh/go-fun/components/learn/frameworks/ui/demo/components"
)

// RegisterAll registers all basic components with the given registry
func RegisterAll(r *components.Registry) {
	r.Register(DefaultGreetingComponent())
	r.Register(DefaultUserCardComponent())
	r.Register(DefaultTodoListComponent())
	r.Register(DefaultButtonComponent())
	r.Register(DefaultEmptyListComponent())
}
