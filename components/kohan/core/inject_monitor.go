package core

import (
	"github.com/amanhigh/go-fun/components/kohan/handler"
	"github.com/amanhigh/go-fun/components/kohan/manager"
	"github.com/golobby/container/v3"
)

// provideOSHandler creates a OSHandler with the given capture path and auto manager.
func provideOSHandler(capturePath string, autoManager manager.AutoManagerInterface) handler.OSHandler {
	return handler.NewOSHandler(capturePath, autoManager)
}

// registerOSDependencies registers all dependencies for the OS feature.
func (ki *KohanInjector) registerOSDependencies(capturePath string, autoManager manager.AutoManagerInterface) {
	container.MustSingleton(ki.di, func() handler.OSHandler {
		return provideOSHandler(capturePath, autoManager)
	})
}
