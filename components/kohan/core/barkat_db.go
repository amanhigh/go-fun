package core

import (
	"fmt"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/models/barkat"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// 3.1 FIXME: Remove barkat_db.go and replace direct DB setup with common util helpers registered via DI config.
// SetupBarkatDB opens (or creates) the SQLite database and auto-migrates barkat tables.
func SetupBarkatDB(dbPath string) (*gorm.DB, error) {
	db, err := util.CreateSqliteDb(dbPath, logger.Warn)
	if err != nil {
		return nil, err
	}

	// Auto-migrate barkat tables
	if err := db.AutoMigrate(&barkat.Entry{}, &barkat.Image{}); err != nil {
		return nil, fmt.Errorf("failed to migrate barkat tables: %w", err)
	}

	return db, nil
}
