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

func (self *BaseDao) FindFirst(c context.Context, entity interface{}) (err common.HttpError) {
	var txErr error
	if txErr = Tx(c).Where(entity).First(entity).Error; txErr != nil && !errors.Is(txErr, gorm.ErrRecordNotFound) {
		log.WithContext(c).WithFields(log.Fields{"Entity": entity, "Error": txErr}).Error("Error Fetching Entity")
	}
	err = util.GormErrorMapper(err)
	return
}

/*
*

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
