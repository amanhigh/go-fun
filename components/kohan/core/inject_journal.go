package core

import (
	"embed"
	"fmt"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/components/kohan/handler"
	"github.com/amanhigh/go-fun/components/kohan/manager"
	"github.com/amanhigh/go-fun/components/kohan/repository"
	"github.com/amanhigh/go-fun/models/barkat"
	"github.com/amanhigh/go-fun/models/config"
	"github.com/gin-gonic/gin"
	"github.com/golobby/container/v3"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

//go:embed migration/*.sql
var migrationFS embed.FS

// ---- Journal Helpers ----

// CreateTestBarkatDB creates a test database using util package with proper migrations.
// This is the recommended approach for integration tests as it uses the same migration system as production.
func CreateTestBarkatDB() (*gorm.DB, error) {
	// Create in-memory SQLite database for testing
	db, err := util.CreateTestDb(logger.Warn)
	if err != nil {
		return nil, fmt.Errorf("failed to create test database: %w", err)
	}

	// Use AutoMigrate for test database (faster and more reliable for in-memory DB)
	if err := db.AutoMigrate(&barkat.Journal{}, &barkat.Image{}, &barkat.Tag{}, &barkat.Note{}); err != nil {
		return nil, fmt.Errorf("failed to auto-migrate barkat tables: %w", err)
	}

	return db, nil
}

// ---- Journal Providers ----

func (ki *KohanInjector) provideBarkatDB() (*gorm.DB, error) {
	// Create database first
	db, err := util.CreateSqliteDb(ki.config.Barkat.DbPath, logger.Warn)
	if err != nil {
		return nil, fmt.Errorf("failed to create barkat db: %w", err)
	}

	// Run migrations using the created GORM DB
	if err := util.RunMigrations(db, migrationFS, "migration"); err != nil {
		return nil, fmt.Errorf("failed to run barkat migrations: %w", err)
	}

	return db, nil
}

func provideHttpServer(cfg config.HttpServerConfig, shutdown util.Shutdown) util.HttpServer {
	return util.NewHttpServer(cfg, gin.Default(), shutdown)
}

func provideKohanServer(
	httpServer util.HttpServer,
	lifecycle util.ServerLifecycle,
) util.HttpServer {
	httpServer.SetLifecycle(lifecycle)
	return httpServer
}

// ---- Entry ----

func provideJournalRepository(db *gorm.DB) repository.JournalRepository {
	return repository.NewJournalRepository(db)
}

func provideJournalManager(repo repository.JournalRepository) manager.JournalManager {
	return manager.NewJournalManager(repo)
}

func provideJournalHandler(mgr manager.JournalManager) handler.JournalHandler {
	return handler.NewJournalHandler(mgr)
}

// ---- Image ----

func provideImageRepository(db *gorm.DB) repository.ImageRepository {
	return repository.NewImageRepository(db)
}

func provideImageManager(entryMgr manager.JournalManager, repo repository.ImageRepository) manager.ImageManager {
	return manager.NewImageManager(entryMgr, repo)
}

func provideImageHandler(mgr manager.ImageManager) handler.ImageHandler {
	return handler.NewImageHandler(mgr)
}

// ---- Note ----

func provideNoteRepository(db *gorm.DB) repository.NoteRepository {
	return repository.NewNoteRepository(db)
}

func provideNoteManager(entryMgr manager.JournalManager, repo repository.NoteRepository) manager.NoteManager {
	return manager.NewNoteManager(entryMgr, repo)
}

func provideNoteHandler(mgr manager.NoteManager) handler.NoteHandler {
	return handler.NewNoteHandler(mgr)
}

// ---- Tag ----

func provideTagRepository(db *gorm.DB) repository.TagRepository {
	return repository.NewTagRepository(db)
}

func provideTagManager(entryMgr manager.JournalManager, repo repository.TagRepository) manager.TagManager {
	return manager.NewTagManager(entryMgr, repo)
}

func provideTagHandler(mgr manager.TagManager) handler.TagHandler {
	return handler.NewTagHandler(mgr)
}

// registerJournalDependencies registers all dependencies for the journal feature.
func (ki *KohanInjector) registerJournalDependencies() error {
	container.MustSingleton(ki.di, ki.provideBarkatDB)

	// Entry
	container.MustSingleton(ki.di, provideJournalRepository)
	container.MustSingleton(ki.di, provideJournalManager)
	container.MustSingleton(ki.di, provideJournalHandler)

	// Image
	container.MustSingleton(ki.di, provideImageRepository)
	container.MustSingleton(ki.di, provideImageManager)
	container.MustSingleton(ki.di, provideImageHandler)

	// Note
	container.MustSingleton(ki.di, provideNoteRepository)
	container.MustSingleton(ki.di, provideNoteManager)
	container.MustSingleton(ki.di, provideNoteHandler)

	// Tag
	container.MustSingleton(ki.di, provideTagRepository)
	container.MustSingleton(ki.di, provideTagManager)
	container.MustSingleton(ki.di, provideTagHandler)

	return nil
}
