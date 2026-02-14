package core

import (
	"fmt"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/components/kohan/handler"
	"github.com/amanhigh/go-fun/components/kohan/manager"
	"github.com/amanhigh/go-fun/components/kohan/repository"
	barkatmodels "github.com/amanhigh/go-fun/models/barkat"
	"github.com/golobby/container/v3"
	"gorm.io/gorm/logger"
)

// ---- Journal Providers ----

func (ki *KohanInjector) provideJournalRepository() (repository.JournalRepository, error) {
	db, err := util.CreateSqliteDb(ki.config.Barkat.DbPath, logger.Warn)
	if err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(&barkatmodels.Entry{}, &barkatmodels.Image{}); err != nil {
		return nil, fmt.Errorf("failed to migrate barkat tables: %w", err)
	}
	return repository.NewJournalRepository(db), nil
}

func provideJournalManager(repo repository.JournalRepository) manager.JournalManager {
	return manager.NewJournalManager(repo)
}

func provideJournalHandler(mgr manager.JournalManager) *handler.JournalHandler {
	return handler.NewJournalHandler(mgr)
}

// registerJournalDependencies registers all dependencies for the journal feature.
func (ki *KohanInjector) registerJournalDependencies() error {
	container.MustSingleton(ki.di, ki.provideJournalRepository)
	container.MustSingleton(ki.di, provideJournalManager)
	container.MustSingleton(ki.di, provideJournalHandler)
	return nil
}
