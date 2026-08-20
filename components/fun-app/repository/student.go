package repository

import (
	"context"
	"errors"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/fun"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type StudentRepository interface {
	util.BaseDbRepository
	ListStudent(c context.Context, studentQuery fun.StudentQuery) (studentList fun.StudentList, err common.HttpError)
	ListStudentAudit(c context.Context, id string) (studentAuditList []fun.StudentAudit, err common.HttpError)
}

type StudentRepositoryImpl struct {
	util.BaseDbRepository
}

var _ StudentRepository = (*StudentRepositoryImpl)(nil)

func NewStudentRepository(baseRepository util.BaseDbRepository) *StudentRepositoryImpl {
	return &StudentRepositoryImpl{BaseDbRepository: baseRepository}
}

func (r *StudentRepositoryImpl) ListStudent(c context.Context, studentQuery fun.StudentQuery) (studentList fun.StudentList, err common.HttpError) {
	var txErr error

	// Build base filtered query (no pagination)
	txn := r.SafeTx(c)
	if studentQuery.Name != "" {
		txn = txn.Where("name like ?", "%"+studentQuery.Name+"%")
	}
	if studentQuery.Gender != "" {
		txn = txn.Where("gender = ?", studentQuery.Gender)
	}

	// Count total filtered records before applying offset/limit
	if txErr = txn.Model(&fun.Student{}).Count(&studentList.Metadata.Total).Error; txErr != nil && !errors.Is(txErr, gorm.ErrRecordNotFound) {
		zerolog.Ctx(c).Error().Any("Query", studentQuery).Err(txErr).Msg("Error Counting Student List")
		err = util.GormErrorMapper(txErr)
		return
	}

	// Apply pagination and sorting, then fetch records
	txn = util.ApplySort(txn.Offset(studentQuery.Offset).Limit(studentQuery.Limit), util.SortOptions{
		SortBy:           studentQuery.SortBy,
		SortOrder:        studentQuery.SortOrder,
		DefaultSortBy:    "name",
		DefaultSortOrder: common.SortOrderAsc,
	})

	if txErr = txn.Find(&studentList.Records).Error; txErr != nil && !errors.Is(txErr, gorm.ErrRecordNotFound) {
		zerolog.Ctx(c).Error().Any("Query", studentQuery).Err(txErr).Msg("Error Fetching Student List")
		err = util.GormErrorMapper(txErr)
	}

	// Set pagination metadata
	studentList.Metadata.Offset = studentQuery.Offset
	studentList.Metadata.Limit = studentQuery.Limit

	return
}

func (r *StudentRepositoryImpl) ListStudentAudit(c context.Context, id string) (studentAuditList []fun.StudentAudit, err common.HttpError) {
	var txErr error
	audit := fun.StudentAudit{Id: id}

	// Fetch Student Audit Records ordered deterministically by audit_id
	query := r.SafeTx(c).Where(audit)
	query = util.ApplySort(query, util.SortOptions{
		DefaultSortBy:    "audit_id",
		DefaultSortOrder: common.SortOrderAsc,
	})
	if txErr = query.Find(&studentAuditList).Error; txErr != nil && !errors.Is(txErr, gorm.ErrRecordNotFound) {
		zerolog.Ctx(c).Error().Str("Id", id).Err(txErr).Msg("Error Fetching Student Audit List")
		err = util.GormErrorMapper(txErr)
	}

	return
}
