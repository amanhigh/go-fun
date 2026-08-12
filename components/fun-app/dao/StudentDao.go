package dao

import (
	"context"
	"errors"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/fun"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// TODO: Rename to Repository package stop using dao in files and package.
type StudentDaoInterface interface {
	util.BaseDbRepositoryInterface
	ListStudent(c context.Context, studentQuery fun.StudentQuery) (studentList fun.StudentList, err common.HttpError)
	ListStudentAudit(c context.Context, id string) (studentAuditList []fun.StudentAudit, err common.HttpError)
}

type StudentDao struct {
	util.BaseDbRepository
}

var _ StudentDaoInterface = (*StudentDao)(nil)

func NewStudentDao(baseRepo util.BaseDbRepository) *StudentDao {
	return &StudentDao{BaseDbRepository: baseRepo}
}

func (pd *StudentDao) ListStudent(c context.Context, studentQuery fun.StudentQuery) (studentList fun.StudentList, err common.HttpError) {
	var txErr error
	// Add Pagination to Query
	txn := pd.SafeTx(c).Offset(studentQuery.Offset).Limit(studentQuery.Limit)

	// Add Query Params if Supplied
	if studentQuery.Name != "" {
		txn = txn.Where("name like ?", "%"+studentQuery.Name+"%")
	}
	if studentQuery.Gender != "" {
		txn = txn.Where("gender = ?", studentQuery.Gender)
	}

	// Add Sorting to Query
	txn = util.ApplySort(txn, util.SortOptions{
		SortBy:    studentQuery.SortBy,
		SortOrder: studentQuery.SortOrder,
	})

	// Execute Query to Get Records and Count
	if txErr = txn.Find(&studentList.Records).Count(&studentList.Metadata.Total).Error; txErr != nil && !errors.Is(txErr, gorm.ErrRecordNotFound) {
		zerolog.Ctx(c).Error().Any("Query", studentQuery).Err(txErr).Msg("Error Fetching Student List")
		err = util.GormErrorMapper(txErr)
	}

	// Set pagination metadata
	studentList.Metadata.Offset = studentQuery.Offset
	studentList.Metadata.Limit = studentQuery.Limit

	return
}

func (pd *StudentDao) ListStudentAudit(c context.Context, id string) (studentAuditList []fun.StudentAudit, err common.HttpError) {
	var txErr error
	audit := fun.StudentAudit{Id: id}

	// Fetch Student Audit Records
	if txErr = pd.SafeTx(c).Where(audit).Find(&studentAuditList).Error; txErr != nil && !errors.Is(txErr, gorm.ErrRecordNotFound) {
		zerolog.Ctx(c).Error().Str("Id", id).Err(txErr).Msg("Error Fetching Student Audit List")
		err = util.GormErrorMapper(txErr)
	}

	return
}
