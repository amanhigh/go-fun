package core

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
)

//go:embed migration/*.sql
var migrationFS embed.FS

// RunMigrations opens a separate mattn/go-sqlite3 connection to the given database path
// and runs all pending golang-migrate migrations. Production only — tests use SetupBarkatDB.
func RunMigrations(dbPath string) error {
	// HACK: Move this to util and make it generic to be used at other places also add unit test using in memory sqlite
	sqlDB, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open sqlite3 for migration: %w", err)
	}
	defer sqlDB.Close()

	return runMigrationsOnDB(sqlDB)
}

func runMigrationsOnDB(sqlDB *sql.DB) error {
	source, err := iofs.New(migrationFS, "migration")
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
	log.Info().Uint("version", version).Bool("dirty", dirty).Msg("Barkat migrations applied")
	return nil
}
