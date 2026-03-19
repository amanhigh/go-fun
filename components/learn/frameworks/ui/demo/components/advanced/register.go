package advanced

import (
	"github.com/amanhigh/go-fun/components/learn/frameworks/ui/demo/components"
)

// RegisterAll registers all advanced components with the given registry
func RegisterAll(r *components.Registry) {
	r.Register(DefaultFullPageComponent())
	r.Register(DefaultDashboardComponent())
}
