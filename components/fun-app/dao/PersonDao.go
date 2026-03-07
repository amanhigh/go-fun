package dao

import (
	"context"
	"errors"
	"fmt"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/fun"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// TODO: Rename to Repository package stop using dao in files and package.
type PersonDaoInterface interface {
	util.BaseDbRepositoryInterface
	ListPerson(c context.Context, personQuery fun.PersonQuery) (personList fun.PersonList, err common.HttpError)
	ListPersonAudit(c context.Context, id string) (personAuditList []fun.PersonAudit, err common.HttpError)
}

type PersonDao struct {
	util.BaseDbRepository
}

var _ PersonDaoInterface = (*PersonDao)(nil)

func NewPersonDao(baseRepo util.BaseDbRepository) *PersonDao {
	return &PersonDao{BaseDbRepository: baseRepo}
}

func (pd *PersonDao) ListPerson(c context.Context, personQuery fun.PersonQuery) (personList fun.PersonList, err common.HttpError) {
	var txErr error
	// Add Pagination to Query
	txn := pd.SafeTx(c).Offset(personQuery.Offset).Limit(personQuery.Limit)

	// Add Query Params if Supplied
	if personQuery.Name != "" {
		txn = txn.Where("name like ?", "%"+personQuery.Name+"%")
	}
	if personQuery.Gender != "" {
		txn = txn.Where("gender = ?", personQuery.Gender)
	}

	// Add Sorting to Query
	if personQuery.SortBy != "" {
		txn = txn.Order(fmt.Sprintf("%s %s", personQuery.SortBy, personQuery.Order))
	}

	// Execute Query to Get Records and Count
	if txErr = txn.Find(&personList.Records).Count(&personList.Metadata.Total).Error; txErr != nil && !errors.Is(txErr, gorm.ErrRecordNotFound) {
		zerolog.Ctx(c).Error().Any("Query", personQuery).Err(txErr).Msg("Error Fetching Person List")
		err = util.GormErrorMapper(txErr)
	}

	// Set pagination metadata
	personList.Metadata.Offset = personQuery.Offset
	personList.Metadata.Limit = personQuery.Limit

	return
}

func (pd *PersonDao) ListPersonAudit(c context.Context, id string) (personAuditList []fun.PersonAudit, err common.HttpError) {
	var txErr error
	audit := fun.PersonAudit{Id: id}

	// Fetch Person Audit Records
	if txErr = pd.SafeTx(c).Where(audit).Find(&personAuditList).Error; txErr != nil && !errors.Is(txErr, gorm.ErrRecordNotFound) {
		zerolog.Ctx(c).Error().Str("Id", id).Err(txErr).Msg("Error Fetching Person Audit List")
		err = util.GormErrorMapper(txErr)
	}

	return
}
