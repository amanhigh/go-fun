package core

import (
	"github.com/amanhigh/go-fun/components/kohan/handler"
	"github.com/amanhigh/go-fun/components/kohan/manager"
	"github.com/golobby/container/v3"
)

func provideIndexPortalHandler() handler.IndexPortal {
	return handler.NewIndexPortal()
}

func provideJournalPortalHandler(mgr manager.JournalManager) handler.JournalPortal {
	return handler.NewJournalPortal(mgr)
}

func providePortalHandlers(indexPortal handler.IndexPortal, journalPortal handler.JournalPortal) PortalHandlers {
	return PortalHandlers{
		IndexPortal:   indexPortal,
		JournalPortal: journalPortal,
	}
}

func (ki *KohanInjector) registerPortalDependencies() {
	container.MustSingleton(ki.di, provideIndexPortalHandler)
	container.MustSingleton(ki.di, provideJournalPortalHandler)
	container.MustSingleton(ki.di, providePortalHandlers)
}
