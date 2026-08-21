package repository

import (
	"context"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/fun"
)

type EnrollmentRepository interface {
	util.BaseDbRepository
	FindByStudentID(ctx context.Context, studentID string, enrollment *fun.Enrollment) common.HttpError
}

type EnrollmentRepositoryImpl struct {
	util.BaseDbRepository
}

var _ EnrollmentRepository = (*EnrollmentRepositoryImpl)(nil)

func NewEnrollmentRepository(baseRepository util.BaseDbRepository) *EnrollmentRepositoryImpl {
	return &EnrollmentRepositoryImpl{BaseDbRepository: baseRepository}
}

func (r *EnrollmentRepositoryImpl) FindByStudentID(ctx context.Context, studentID string, enrollment *fun.Enrollment) common.HttpError {
	return util.GormErrorMapper(r.SafeTx(ctx).Where("student_id = ?", studentID).First(enrollment).Error)
}
