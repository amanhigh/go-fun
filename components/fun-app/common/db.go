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
		&fun.Person{},
		&fun.PersonAudit{},
		&fun.Enrollment{},
	)
	return
}

// DAO providers return interfaces while delegating to pointer-returning constructors.

func (fi *FunAppInjector) registerDao() {
	container.MustSingleton(fi.di, util.NewBaseDao)
	container.MustSingleton(fi.di, fi.providePersonDao)
	container.MustSingleton(fi.di, fi.provideEnrollmentDao)
}

func (fi *FunAppInjector) providePersonDao(base util.BaseDao) dao.PersonDaoInterface {
	return dao.NewPersonDao(base)
}

func (fi *FunAppInjector) provideEnrollmentDao(base util.BaseDao) dao.EnrollmentDaoInterface {
	return dao.NewEnrollmentDao(base)
}
