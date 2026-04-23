package core

import (
	"github.com/amanhigh/go-fun/components/kohan/handler"
	"github.com/amanhigh/go-fun/models/config"
	"github.com/golobby/container/v3"
)

func provideIndexPortalHandler() handler.IndexPortal {
	return handler.NewIndexPortal()
}

func provideJournalPortalHandler(cfg config.BarkatConfig) handler.JournalPortal {
	return handler.NewJournalPortal(cfg.ScreenshotPath)
}

func providePortalHandlers(indexPortal handler.IndexPortal, journalPortal handler.JournalPortal) PortalHandlers {
	return PortalHandlers{
		IndexPortal:   indexPortal,
		JournalPortal: journalPortal,
	}
}

func (ki *KohanInjector) registerPortalDependencies() {
	container.MustSingleton(ki.di, func() config.BarkatConfig {
		return ki.config.Barkat
	})
	container.MustSingleton(ki.di, provideIndexPortalHandler)
	container.MustSingleton(ki.di, provideJournalPortalHandler)
	container.MustSingleton(ki.di, providePortalHandlers)
}
