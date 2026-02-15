package util

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/mattn/go-sqlite3" // CGO sqlite3 driver required by golang-migrate
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// RunSqliteMigrations opens a separate mattn/go-sqlite3 connection to the given database path
// and runs all pending golang-migrate migrations from the provided embedded filesystem.
func RunSqliteMigrations(dbPath string, migrationFS embed.FS, migrationDir string) error {
	sqlDB, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open sqlite3 for migration: %w", err)
	}
	defer sqlDB.Close()

	return runMigrationsOnDB(sqlDB, migrationFS, migrationDir)
}

// RunSqliteMigrationsFromDB runs migrations using an existing GORM DB instance.
// This is more generic and allows running migrations in tests via the migration framework.
func RunSqliteMigrationsFromDB(db *gorm.DB, migrationFS embed.FS, migrationDir string) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB from gorm: %w", err)
	}
	return runMigrationsOnDB(sqlDB, migrationFS, migrationDir)
}

func runMigrationsOnDB(sqlDB *sql.DB, migrationFS embed.FS, migrationDir string) error {
	source, err := iofs.New(migrationFS, migrationDir)
	if err != nil {
		return fmt.Errorf("failed to create iofs source: %w", err)
	}

	driver, err := sqlite3.WithInstance(sqlDB, &sqlite3.Config{})
	if err != nil {
		return fmt.Errorf("failed to create sqlite3 driver: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", source, "sqlite3", driver)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migration failed: %w", err)
	}

	version, dirty, _ := m.Version()
	log.Info().Uint("version", version).Bool("dirty", dirty).Msg("Migrations applied")
	return nil
}
