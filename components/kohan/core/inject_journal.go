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

// SetupBarkatDB runs migrations for all barkat tables on the given database.
func SetupBarkatDB(db *gorm.DB) error {
	// HACK: Migrate using Migration go-migrate framework and remove AutoMigrate (Use it only for tests)
	if err := db.AutoMigrate(&barkatmodels.Entry{}, &barkatmodels.Image{}, &barkatmodels.Tag{}, &barkatmodels.Note{}); err != nil {
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

func provideBaseHTTPServer(port int, shutdown util.Shutdown) *util.BaseHTTPServer {
	return util.NewBaseHTTPServer("kohan", port, shutdown)
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
