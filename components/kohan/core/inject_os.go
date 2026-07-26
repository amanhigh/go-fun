package core

import (
	"fmt"

	"github.com/amanhigh/go-fun/components/kohan/handler"
	"github.com/amanhigh/go-fun/components/kohan/manager"
	"github.com/go-co-op/gocron/v2"
	"github.com/golobby/container/v3"
)

// provideScheduler creates a gocron scheduler, panicking on construction failure.
func provideScheduler() gocron.Scheduler {
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		panic(fmt.Errorf("failed to create scheduler: %w", err))
	}
	return scheduler
}

// provideOSHandler creates an OSHandler from the DI-resolved OSManagerInterface.
func provideOSHandler(osManager manager.OSManagerInterface) handler.OSHandler {
	return handler.NewOSHandler(osManager)
}

// provideOSManager creates an OSManager from the injector's config and a gocron scheduler.
func (ki *KohanInjector) provideOSManager(scheduler gocron.Scheduler) manager.OSManagerInterface {
	return manager.NewOSManager(ki.config.OSWaitInterval, ki.config.Barkat.ScreenshotPath, scheduler)
}

// registerOSDependencies registers all dependencies for the OS feature.
func (ki *KohanInjector) registerOSDependencies() {
	container.MustSingleton(ki.di, provideScheduler)
	container.MustSingleton(ki.di, ki.provideOSManager)
	container.MustSingleton(ki.di, provideOSHandler)
}
