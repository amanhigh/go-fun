package util

import (
	"context"
	"errors"
	"time"

	"github.com/amanhigh/go-fun/models"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

const TX_TIMEOUT = 30 * time.Second

type BaseDaoInterface interface {
	FindById(c context.Context, id any, entity any) (err common.HttpError)
	FindPaginated(c context.Context, pageParams common.Pagination, result any) (count int64, err common.HttpError)
	Create(c context.Context, entity any, omit ...string) (err common.HttpError)
	Update(c context.Context, entity any, omit ...string) (err common.HttpError)
	DeleteById(c context.Context, id any, entity any) (err common.HttpError)
	GetCount(c context.Context, entity any) (count int64, err common.HttpError)
	UseOrCreateTx(c context.Context, run DbRun, readOnly ...bool) (err common.HttpError)
}

type BaseDao struct {
	Db *gorm.DB
}

func NewBaseDao(db *gorm.DB) BaseDao {
	return BaseDao{Db: db}
}

type DbRun func(c context.Context) (err common.HttpError)

func (self *BaseDao) FindById(c context.Context, id any, entity any) (err common.HttpError) {
	var txErr error
	if txErr = Tx(c).First(entity, "id=?", id).Error; txErr != nil && !errors.Is(txErr, gorm.ErrRecordNotFound) {
		log.Ctx(c).Error().Any("Id", id).Any("Entity", entity).Err(txErr).Msg("Error Fetching Entity")
	}
	err = GormErrorMapper(txErr)
	return
}

func (self *BaseDao) FindPaginated(c context.Context, pageParams common.Pagination, result any) (count int64, err common.HttpError) {
	var txErr error
	if txErr = Tx(c).Offset(pageParams.Offset).Limit(pageParams.Limit).Find(result).Error; txErr != nil && !errors.Is(txErr, gorm.ErrRecordNotFound) {
		log.Ctx(c).Error().Any("paginationParams", pageParams).Int64("TotalCount", count).
			Err(txErr).Msg("Error Fetching Paginated Entity")
		err = GormErrorMapper(txErr)
	} else {
		//Add count to Paginated Result
		count, err = self.GetCount(c, result)
	}
	return
}

func (self *BaseDao) Create(c context.Context, entity any, omit ...string) (err common.HttpError) {
	var txErr error
	if txErr = Tx(c).Omit(omit...).Create(entity).Error; txErr != nil {
		log.Ctx(c).Error().Any("Entity", entity).Err(txErr).Msg("Entity Create Failed")
	}
	//Error Conversion
	err = GormErrorMapper(txErr)
	return
}

func (self *BaseDao) Update(c context.Context, entity any, omit ...string) (err common.HttpError) {
	if txErr := Tx(c).Omit(omit...).Save(entity).Error; txErr != nil {
		log.Ctx(c).Error().
			Any("Entity", entity).Err(txErr).Msg("Entity Update Failed")
		err = GormErrorMapper(txErr)
	}
	return
}

func (self *BaseDao) DeleteById(c context.Context, id any, entity any) (err common.HttpError) {
	if txErr := Tx(c).Delete(entity, "id=?", id).Error; txErr != nil {
		log.Ctx(c).Error().
			Any("Id", id).Any("Entity", entity).Err(txErr).Msg("Entity Delete Failed")
		err = GormErrorMapper(txErr)
	}
	return
}

func (self *BaseDao) GetCount(c context.Context, entity any) (count int64, err common.HttpError) {
	var txErr error
	if txErr = Tx(c).Model(entity).Count(&count).Error; txErr != nil && !errors.Is(txErr, gorm.ErrRecordNotFound) {
		log.Ctx(c).Error().Any("Entity", entity).Err(txErr).Msg("Error Getting Entity Count")
	}
	err = GormErrorMapper(txErr)
	return
}

func (self *BaseDao) SetPagination(query *gorm.DB, offset, limit int) {
	query.Offset(offset)
	if limit > 0 {
		query.Limit(limit)
	}
}

/*
Transaction Handling to use already created transaction or Init New.
Needs State, hence placed in BaseDao (Not Util)
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
