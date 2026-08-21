package common

import (
	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/components/fun-app/repository"
	"github.com/amanhigh/go-fun/models/config"
	"github.com/amanhigh/go-fun/models/fun"
	"github.com/golobby/container/v3"

	// Blank import for mysql driver registration
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/gorm"
)

func newDb(config config.FunAppConfig) (db *gorm.DB, err error) {
	db = util.MustCreateDb(config.Db)

	/** Gorm AutoMigrate Schema */
	err = db.AutoMigrate(
		&fun.Student{},
		&fun.StudentAudit{},
		&fun.Enrollment{},
	)
	return
}

// Repository providers return interfaces while delegating to pointer-returning constructors.

func (fi *FunAppInjector) registerRepositories() {
	container.MustSingleton(fi.di, provideBaseDbRepository)
	container.MustSingleton(fi.di, fi.provideStudentRepository)
	container.MustSingleton(fi.di, fi.provideEnrollmentRepository)
}

func provideBaseDbRepository(db *gorm.DB) util.BaseDbRepository {
	return util.NewBaseDbRepository(db)
}

func (fi *FunAppInjector) provideStudentRepository(baseRepository util.BaseDbRepository) repository.StudentRepository {
	return repository.NewStudentRepository(baseRepository)
}

func (fi *FunAppInjector) provideEnrollmentRepository(baseRepository util.BaseDbRepository) repository.EnrollmentRepository {
	return repository.NewEnrollmentRepository(baseRepository)
}
