package core

import (
	"github.com/amanhigh/go-fun/components/kohan/handler"
	"github.com/amanhigh/go-fun/components/kohan/manager"
	"github.com/golobby/container/v3"
)

// provideMonitorHandler creates a MonitorHandler with the given capture path and auto manager.
func provideMonitorHandler(capturePath string, autoManager manager.AutoManagerInterface) handler.MonitorHandler {
	return handler.NewMonitorHandler(capturePath, autoManager)
}

// registerMonitorDependencies registers all dependencies for the monitor feature.
func (ki *KohanInjector) registerMonitorDependencies(capturePath string, autoManager manager.AutoManagerInterface) {
	container.MustSingleton(ki.di, func() handler.MonitorHandler {
		return provideMonitorHandler(capturePath, autoManager)
	})
}
