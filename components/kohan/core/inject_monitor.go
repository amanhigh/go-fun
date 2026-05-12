package core

import (
	"github.com/amanhigh/go-fun/components/kohan/handler"
	"github.com/amanhigh/go-fun/components/kohan/manager"
	"github.com/golobby/container/v3"
)

// provideOSHandler creates an OSHandler with the configured screenshot path.
func provideOSHandler(osManager manager.OSManagerInterface) handler.OSHandler {
	return handler.NewOSHandler(osManager)
}

// registerOSDependencies registers all dependencies for the OS feature.
func (ki *KohanInjector) registerOSDependencies(osManager manager.OSManagerInterface) {
	container.MustSingleton(ki.di, func() handler.OSHandler {
		return provideOSHandler(osManager)
	})
}
