package dao

import (
	"context"
	"errors"
	"time"

	"github.com/amanhigh/go-fun/common/util"
	. "github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/models"
	"github.com/amanhigh/go-fun/models/common"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const TX_TIMEOUT = 30 * time.Second

type BaseDao struct {
	Db *gorm.DB `inject:""`
}

type DbRun func(c context.Context) (err common.HttpError)

func (self *BaseDao) FindById(c context.Context, id any, entity any) (err common.HttpError) {
	var txErr error
	if txErr = Tx(c).First(entity, "id=?", id).Error; txErr != nil && !errors.Is(txErr, gorm.ErrRecordNotFound) {
		log.WithContext(c).WithFields(log.Fields{"Id": id, "Entity": entity, "Error": txErr}).Error("Error Fetching Entity")
	}
	err = util.GormErrorMapper(txErr)
	return
}

func (self *BaseDao) Find(c context.Context, query any, result any) (err common.HttpError) {
	var txErr error
	if txErr = Tx(c).Where(query).Find(result).Error; txErr != nil && !errors.Is(txErr, gorm.ErrRecordNotFound) {
		log.WithContext(c).WithFields(log.Fields{"query": query, "Error": txErr}).Error("Error Fetching Entity")
	}
	err = util.GormErrorMapper(txErr)
	return
}

func (self *BaseDao) Create(c context.Context, entity any, omit ...string) (err common.HttpError) {
	var txErr error
	if txErr = Tx(c).Omit(omit...).Create(entity).Error; txErr != nil {
		log.WithContext(c).WithFields(log.Fields{"Entity": entity}).WithField("Err", txErr).Error("Entity Create Failed")
	}
	//Error Conversion
	err = util.GormErrorMapper(txErr)
	return
}

func (self *BaseDao) Update(c context.Context, entity any, omit ...string) (err common.HttpError) {
	if txErr := Tx(c).Omit(omit...).Save(entity).Error; txErr != nil {
		log.WithContext(c).WithFields(log.Fields{"Entity": entity, "Error": txErr}).Error("Entity Update Failed")
		err = util.GormErrorMapper(txErr)
	}
	return
}

func (self *BaseDao) DeleteById(c context.Context, id any, entity any) (err common.HttpError) {
	if txErr := Tx(c).Delete(entity, id).Error; txErr != nil {
		log.WithContext(c).WithFields(log.Fields{"Id": id, "Entity": entity, "Error": txErr}).
			Error("Entity Delete Failed")
		err = util.GormErrorMapper(txErr)
	}
	return
}

func (self *BaseDao) GetCount(c context.Context, entity any) (count int64, err common.HttpError) {
	var txErr error
	tx := Tx(c).Model(entity).Where(entity)
	if txErr = tx.Count(&count).Error; txErr != nil && !errors.Is(txErr, gorm.ErrRecordNotFound) {
		log.WithContext(c).WithFields(log.Fields{"Entity": entity, "Error": txErr}).Error("Error Fetching Entity")
	}
	err = util.GormErrorMapper(txErr)
	return
}

/*
Transaction Handling
*/
func (self *BaseDao) UseOrCreateTx(c context.Context, run DbRun, readOnly ...bool) (err common.HttpError) {
	//Check if Context has Tx
	if Tx(c) != nil {
		//First Preference to use existing tx if supplied
		err = run(c)
	} else if len(readOnly) > 0 && readOnly[0] {
		// Set Timeout on DB
		ctx, cancel := context.WithTimeout(c, TX_TIMEOUT)
		if cancel != nil {
			defer cancel()
		}
		// Avoid Creating New Transaction for Readonly (Use DB)
		err = run(context.WithValue(c, models.CONTEXT_TX, self.Db.WithContext(ctx)))
	} else {
		// Create Transaction With Timeout
		ctx, cancel := context.WithTimeout(c, TX_TIMEOUT)
		if cancel != nil {
			defer cancel()
		}

		// Error Returned after running completes in Transaction.
		var txErr error

		// Inject Transaction in Context
		txErr = self.Db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			return run(context.WithValue(c, models.CONTEXT_TX, tx))
		})

		/* Morph Transaction Error to Http Error */
		if ok := errors.As(txErr, &err); txErr != nil && !ok {
			//This Should Not Happen.
			err = common.NewServerError(txErr)
		}
	}

	return
}
