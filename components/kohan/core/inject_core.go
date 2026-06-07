package core

import (
	"fmt"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/models/config"
	"github.com/golobby/container/v3"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// registerConfig registers shared configuration in the container.
func (ki *KohanInjector) registerConfig() {
	container.MustSingleton(ki.di, func() config.BarkatConfig { return ki.config.Barkat })
}

// createDB opens the database, runs migrations, and registers *gorm.DB in the container.
func (ki *KohanInjector) createDB() error {
	log.Info().Str("db_path", ki.config.Barkat.DbPath).Msg("Opening Barkat database")
	db, err := util.CreateSqliteDb(ki.config.Barkat.DbPath, logger.Warn)
	if err != nil {
		return fmt.Errorf("failed to create barkat db: %w", err)
	}
	if err := util.RunMigrations(db, migrationFS, "migration"); err != nil {
		return fmt.Errorf("failed to run barkat migrations: %w", err)
	}

	container.MustSingleton(ki.di, func() *gorm.DB { return db })
	return nil
}

// registerCoreDependencies registers shared dependencies used by all features.
func (ki *KohanInjector) registerCoreDependencies() error {
	ki.registerConfig()
	return ki.createDB()
}
