package dao

import (
	"context"
	"errors"
	"fmt"

	"github.com/amanhigh/go-fun/common/util"
	. "github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/fun"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type PersonDaoInterface interface {
	util.BaseDaoInterface
	ListPerson(c context.Context, personQuery fun.PersonQuery) (personList fun.PersonList, err common.HttpError)
}

type PersonDao struct {
	BaseDao `inject:"inline"`
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
		log.WithContext(c).WithFields(log.Fields{"Query": personQuery, "Error": txErr}).Error("Error Fetching Person List")
		err = GormErrorMapper(txErr)
	}

	return
}
