package dao

import (
	"context"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/fun"
)

type EnrollmentDaoInterface interface {
	util.BaseDaoInterface
	FindByPersonID(ctx context.Context, personID string, enrollment *fun.Enrollment) common.HttpError
}

type EnrollmentDao struct {
	util.BaseDao
}

func NewEnrollmentDao(baseDao util.BaseDao) EnrollmentDaoInterface {
	return &EnrollmentDao{BaseDao: baseDao}
}

func (ed *EnrollmentDao) FindByPersonID(ctx context.Context, personID string, enrollment *fun.Enrollment) common.HttpError {
	query := ed.Db.WithContext(ctx)
	if tx := util.Tx(ctx); tx != nil {
		query = tx
	}
	return util.GormErrorMapper(query.Where("person_id = ?", personID).First(enrollment).Error)
}
