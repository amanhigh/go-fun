package util

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

/**
Resolver Policy which sends traffic to slave,
incase of error switches to fallback (master).

At given intervals try to revert config back to slave
based on Pings.
*/
type FallBackPolicy struct {
	//Current Pool which should be Used
	currentPool int
	//DB to do ping checks
	db *gorm.DB
	//Channel to Report Errors
	errChan chan error
}

func NewFallBackPolicy(Db *gorm.DB) *FallBackPolicy {
	return &FallBackPolicy{
		currentPool: 0,
		errChan:     make(chan error, 5),
		db:          Db,
	}
}

/**
Resolve Function Implementation for dbResolver.
*/
func (self *FallBackPolicy) Resolve(connPools []gorm.ConnPool) gorm.ConnPool {
	x := self.GetPool()
	fmt.Println("Pools", len(connPools), x)
	pool := connPools[x]
	return pool
}

/**
Report any DB Errors, can handle nils
*/
func (self *FallBackPolicy) ReportError(err error) {
	self.errChan <- err
}

func (self *FallBackPolicy) GetPool() (poolIndex int) {
	select {
	case err, ok := <-self.errChan:
		//If Channel is Not Closed Update Pool
		fmt.Println("ChanWrite -->", ok, err)
		if ok && err != nil {
			log.WithFields(log.Fields{"Error": err}).Error("Falling Back to Master for Reads")
			//Update New Value in Cache
			self.currentPool = 1
		}
		//Serve Updated Pool
		poolIndex = self.currentPool
	default:
		//Serve Current Pool
		poolIndex = self.currentPool
	}
	return
}

func (self *FallBackPolicy) Ping(resolver dbresolver.Operation) (err error) {
	err = self.db.Clauses(resolver).Raw("show tables").Error
	fmt.Println(resolver, err)
	return
}
