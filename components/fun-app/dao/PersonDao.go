package dao

import (
	"context"
	"errors"
	"fmt"

	"github.com/amanhigh/go-fun/common/util"
	. "github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/fun"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type PersonDaoInterface interface {
	util.BaseDaoInterface
	ListPerson(c context.Context, personQuery fun.PersonQuery) (personList fun.PersonList, err common.HttpError)
	ListPersonAudit(c context.Context, id string) (personAuditList []fun.PersonAudit, err common.HttpError)
}

type PersonDao struct {
	BaseDao
}

func NewPersonDao(baseDao util.BaseDao) PersonDaoInterface {
	return &PersonDao{BaseDao: baseDao}
}

func (self *PersonDao) ListPerson(c context.Context, personQuery fun.PersonQuery) (personList fun.PersonList, err common.HttpError) {
	var txErr error
	//Add Pagination to Query
	txn := Tx(c).Offset(personQuery.Offset).Limit(personQuery.Limit)

	//Add Query Params if Supplied
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

	//Execute Query to Get Records and Count
	if txErr = txn.Find(&personList.Records).Count(&personList.Metadata.Total).Error; txErr != nil && !errors.Is(txErr, gorm.ErrRecordNotFound) {
		zerolog.Ctx(c).Error().Any("Query", personQuery).Err(txErr).Msg("Error Fetching Person List")
		err = GormErrorMapper(txErr)
	}

	return
}

func (self *PersonDao) ListPersonAudit(c context.Context, id string) (personAuditList []fun.PersonAudit, err common.HttpError) {
	var txErr error
	var audit = fun.PersonAudit{}
	audit.Id = id

	//Fetch Person Audit Records
	if txErr = Tx(c).Where(audit).Find(&personAuditList).Error; txErr != nil && !errors.Is(txErr, gorm.ErrRecordNotFound) {
		zerolog.Ctx(c).Error().Str("Id", id).Err(txErr).Msg("Error Fetching Person Audit List")
		err = GormErrorMapper(txErr)
	}

	return
}
