package dao

import (
	"context"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/fun"
)

type EnrollmentDaoInterface interface {
	// TODO: Rename to BaseDbRepository Interface & Files and use across Repo in FunApp and Kohan where GORM is used.
	util.BaseDbRepositoryInterface
	FindByStudentID(ctx context.Context, studentID string, enrollment *fun.Enrollment) common.HttpError
}

type EnrollmentDao struct {
	util.BaseDbRepository
}

var _ EnrollmentDaoInterface = (*EnrollmentDao)(nil)

func NewEnrollmentDao(baseRepo util.BaseDbRepository) *EnrollmentDao {
	return &EnrollmentDao{BaseDbRepository: baseRepo}
}

func (ed *EnrollmentDao) FindByStudentID(ctx context.Context, studentID string, enrollment *fun.Enrollment) common.HttpError {
	return util.GormErrorMapper(ed.SafeTx(ctx).Where("student_id = ?", studentID).First(enrollment).Error)
}
