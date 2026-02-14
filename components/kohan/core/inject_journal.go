package core

import (
	"fmt"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/components/kohan/handler"
	"github.com/amanhigh/go-fun/components/kohan/manager"
	"github.com/amanhigh/go-fun/components/kohan/repository"
	barkatmodels "github.com/amanhigh/go-fun/models/barkat"
	"github.com/golobby/container/v3"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ---- Journal Helpers ----

// SetupBarkatDB runs AutoMigrate for all barkat tables on the given database.
// Used by both DI providers and tests to ensure schema is ready.
func SetupBarkatDB(db *gorm.DB) error {
	if err := db.AutoMigrate(&barkatmodels.Entry{}, &barkatmodels.Image{}); err != nil {
		return fmt.Errorf("failed to migrate barkat tables: %w", err)
	}
	return nil
}

// ---- Journal Providers ----

func (ki *KohanInjector) provideBarkatDB() (*gorm.DB, error) {
	db, err := util.CreateSqliteDb(ki.config.Barkat.DbPath, logger.Warn)
	if err != nil {
		return nil, err
	}
	if err := SetupBarkatDB(db); err != nil {
		return nil, err
	}
	return db, nil
}

func provideJournalRepository(db *gorm.DB) repository.JournalRepository {
	return repository.NewJournalRepository(db)
}

func provideJournalManager(repo repository.JournalRepository) manager.JournalManager {
	return manager.NewJournalManager(repo)
}

func provideJournalHandler(mgr manager.JournalManager) handler.JournalHandler {
	return handler.NewJournalHandler(mgr)
}

// registerJournalDependencies registers all dependencies for the journal feature.
func (ki *KohanInjector) registerJournalDependencies() error {
	container.MustSingleton(ki.di, ki.provideBarkatDB)
	container.MustSingleton(ki.di, provideJournalRepository)
	container.MustSingleton(ki.di, provideJournalManager)
	container.MustSingleton(ki.di, provideJournalHandler)
	return nil
}
