package core

import (
	"embed"
	"fmt"

	"github.com/amanhigh/go-fun/common/util"
)

//go:embed migration/*.sql
var migrationFS embed.FS

// RunMigrations opens a separate mattn/go-sqlite3 connection to the given database path
// and runs all pending golang-migrate migrations. Production only — tests use SetupBarkatDB.
func RunMigrations(dbPath string) error {
	if err := util.RunSqliteMigrations(dbPath, migrationFS, "migration"); err != nil {
		return fmt.Errorf("failed to run barkat migrations: %w", err)
	}
	return nil
}
