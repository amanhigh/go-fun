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

type BaseDbRepositoryInterface interface {
	FindById(c context.Context, id any, entity any) (err common.HttpError)
	FindPaginated(c context.Context, pageParams common.Pagination, result any) (count int64, err common.HttpError)
	Create(c context.Context, entity any, omit ...string) (err common.HttpError)
	Update(c context.Context, entity any, omit ...string) (err common.HttpError)
	DeleteById(c context.Context, id any, entity any) (err common.HttpError)
	GetCount(c context.Context, entity any) (count int64, err common.HttpError)
	UseOrCreateTx(c context.Context, run DbRun, readOnly ...bool) (err common.HttpError)
	GetByExternalId(c context.Context, externalId string, entity any) (err common.HttpError)
}

type BaseDbRepository struct {
	Db *gorm.DB
}

func NewBaseDbRepository(db *gorm.DB) BaseDbRepository {
	return BaseDbRepository{Db: db}
}

type DbRun func(c context.Context) (err common.HttpError)

func (b *BaseDbRepository) FindById(c context.Context, id, entity any) (err common.HttpError) {
	var txErr error
	query := b.SafeTx(c)
	if txErr = query.First(entity, "id=?", id).Error; txErr != nil && !errors.Is(txErr, gorm.ErrRecordNotFound) {
		log.Ctx(c).Error().Any("Id", id).Any("Entity", entity).Err(txErr).Msg("Error Fetching Entity")
	}
	err = GormErrorMapper(txErr)
	return
}

func (b *BaseDbRepository) FindPaginated(c context.Context, pageParams common.Pagination, result any) (count int64, err common.HttpError) {
	var txErr error
	query := b.SafeTx(c)
	if txErr = query.Offset(pageParams.Offset).Limit(pageParams.Limit).Find(result).Error; txErr != nil && !errors.Is(txErr, gorm.ErrRecordNotFound) {
		log.Ctx(c).Error().Any("paginationParams", pageParams).Int64("TotalCount", count).
			Err(txErr).Msg("Error Fetching Paginated Entity")
		err = GormErrorMapper(txErr)
	} else {
		// Add count to Paginated Result
		count, err = b.GetCount(c, result)
	}
	return
}

func (b *BaseDbRepository) Create(c context.Context, entity any, omit ...string) (err common.HttpError) {
	var txErr error
	query := b.SafeTx(c)
	if txErr = query.Omit(omit...).Create(entity).Error; txErr != nil {
		log.Ctx(c).Error().Any("Entity", entity).Err(txErr).Msg("Entity Create Failed")
	}
	// Error Conversion
	err = GormErrorMapper(txErr)
	return
}

func (b *BaseDbRepository) Update(c context.Context, entity any, omit ...string) (err common.HttpError) {
	var txErr error
	query := b.SafeTx(c)
	if txErr = query.Omit(omit...).Save(entity).Error; txErr != nil {
		log.Ctx(c).Error().Any("Entity", entity).Err(txErr).Msg("Entity Update Failed")
	}
	err = GormErrorMapper(txErr)
	return
}

func (b *BaseDbRepository) DeleteById(c context.Context, id, entity any) (err common.HttpError) {
	query := b.SafeTx(c)
	result := query.Delete(entity, "id=?", id)
	if result.Error != nil {
		return GormErrorMapper(result.Error)
	}
	if result.RowsAffected == 0 {
		return common.ErrNotFound
	}
	return nil
}

func (b *BaseDbRepository) GetCount(c context.Context, entity any) (count int64, err common.HttpError) {
	var txErr error
	query := b.SafeTx(c)
	if txErr = query.Model(entity).Count(&count).Error; txErr != nil && !errors.Is(txErr, gorm.ErrRecordNotFound) {
		return 0, GormErrorMapper(txErr)
	}
	return count, nil
}

func (b *BaseDbRepository) SetPagination(query *gorm.DB, offset, limit int) {
	query.Offset(offset)
	if limit > 0 {
		query.Limit(limit)
	}
}

/*
Transaction Handling to use already created transaction or Init New.
Needs State, hence placed in BaseDbRepository (Not Util)
*/
func (b *BaseDbRepository) UseOrCreateTx(c context.Context, run DbRun, readOnly ...bool) (err common.HttpError) {
	// Check if Context has Tx
	switch {
	case Tx(c) != nil:
		// First Preference to use existing tx if supplied
		err = run(c)
	case len(readOnly) > 0 && readOnly[0]:
		// Set Timeout on DB
		ctx, cancel := context.WithTimeout(c, TX_TIMEOUT)
		if cancel != nil {
			defer cancel()
		}
		// Avoid Creating New Transaction for Readonly (Use DB)
		err = run(context.WithValue(c, models.ContextTx, b.Db.WithContext(ctx)))
	default:
		// Create Transaction With Timeout
		ctx, cancel := context.WithTimeout(c, TX_TIMEOUT)
		if cancel != nil {
			defer cancel()
		}

		// Error Returned after running completes in Transaction.
		// Inject Transaction in Context
		txErr := b.Db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			return run(context.WithValue(c, models.ContextTx, tx))
		})

		// Morph Transaction Error to Http Error
		if ok := errors.As(txErr, &err); txErr != nil && !ok {
			// This Should Not Happen.
			err = common.NewServerError(txErr)
		}
	}

	return
}

// SafeTx returns a database query with automatic fallback to the repository's database
// when no transaction is present in the context. This eliminates the need for manual nil checks.
func (b *BaseDbRepository) SafeTx(c context.Context) *gorm.DB {
	query := Tx(c)
	if query == nil {
		query = b.Db.WithContext(c)
	}
	return query
}

// GetByExternalId finds an entity by its external_id field
func (b *BaseDbRepository) GetByExternalId(c context.Context, externalId string, entity any) (err common.HttpError) {
	var txErr error
	query := b.SafeTx(c)
	if txErr = query.First(entity, "external_id=?", externalId).Error; txErr != nil && !errors.Is(txErr, gorm.ErrRecordNotFound) {
		log.Ctx(c).Error().Str("ExternalId", externalId).Any("Entity", entity).Err(txErr).Msg("Error Fetching Entity by External ID")
	}
	err = GormErrorMapper(txErr)
	return
}
