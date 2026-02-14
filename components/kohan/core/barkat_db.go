package core

import (
	"fmt"

	"github.com/amanhigh/go-fun/models/barkat"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// SetupBarkatDB opens (or creates) the SQLite database and auto-migrates barkat tables.
func SetupBarkatDB(dbPath string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open barkat db at %s: %w", dbPath, err)
	}

	// Enable WAL mode for better concurrent read performance
	if err := db.Exec("PRAGMA journal_mode=WAL").Error; err != nil {
		return nil, fmt.Errorf("failed to enable WAL mode: %w", err)
	}

	// Auto-migrate barkat tables
	if err := db.AutoMigrate(&barkat.Entry{}, &barkat.Image{}); err != nil {
		return nil, fmt.Errorf("failed to migrate barkat tables: %w", err)
	}

	return db, nil
}
