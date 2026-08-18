package common

import (
	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/components/fun-app/dao"
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

// DAO providers return interfaces while delegating to pointer-returning constructors.

func (fi *FunAppInjector) registerDao() {
	container.MustSingleton(fi.di, util.NewBaseDbRepository)
	container.MustSingleton(fi.di, fi.provideStudentDao)
	container.MustSingleton(fi.di, fi.provideEnrollmentDao)
}

func (fi *FunAppInjector) provideStudentDao(base util.BaseDbRepository) dao.StudentDaoInterface {
	return dao.NewStudentDao(base)
}

func (fi *FunAppInjector) provideEnrollmentDao(base util.BaseDbRepository) dao.EnrollmentDaoInterface {
	return dao.NewEnrollmentDao(base)
}
