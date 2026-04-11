package core

import (
	"github.com/amanhigh/go-fun/components/kohan/handler"
	"github.com/golobby/container/v3"
)

func provideIndexPortalHandler() handler.IndexPortal {
	return handler.NewIndexPortal()
}

func provideJournalPortalHandler() handler.JournalPortal {
	return handler.NewJournalPortal()
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
