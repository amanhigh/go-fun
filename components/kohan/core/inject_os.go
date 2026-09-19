package core

import (
	"github.com/amanhigh/go-fun/components/kohan/handler"
	"github.com/amanhigh/go-fun/components/kohan/manager"
	"github.com/golobby/container/v3"
)

// provideOSHandler creates an OSHandler from the DI-resolved OSManagerInterface.
func provideOSHandler(osManager manager.OSManagerInterface) handler.OSHandler {
	return handler.NewOSHandler(osManager)
}

// provideOSManager creates an OSManager from the injector's screenshot path.
func (ki *KohanInjector) provideOSManager() manager.OSManagerInterface {
	return manager.NewOSManager(ki.config.Barkat.ScreenshotPath)
}

// registerOSDependencies registers all dependencies for the OS feature.
func (ki *KohanInjector) registerOSDependencies() {
	container.MustSingleton(ki.di, ki.provideOSManager)
	container.MustSingleton(ki.di, provideOSHandler)
}
