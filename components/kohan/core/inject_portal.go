package core

import (
	"github.com/amanhigh/go-fun/components/kohan/handler"
	"github.com/golobby/container/v3"
)

func provideIndexPortalHandler() handler.IndexPortal {
	return handler.NewIndexPortal()
}

func (ki *KohanInjector) registerPortalDependencies() {
	container.MustSingleton(ki.di, provideIndexPortalHandler)
}
