package util

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
	"time"
)

/**
Resolver Policy which sends traffic to slave,
incase of error switches to fallback (master).

At given intervals try to revert config back to slave
based on Pings. It needs Ping Table and a Ping dbresolver for it
to work.
*/
type FallBackPolicy struct {
	//Current Pool which should be Used
	currentPool int
	//DB to do ping checks
	db *gorm.DB
	//Channel to Report Errors
	errChan chan error
	//Interval retry to revert to original Config
	ticker *time.Ticker
	//Model for Table to get count as a Ping
	pingTable string
}

const (
	PING          = "ping"
	POOL_PRIMARY  = 0
	POOL_FALLBACK = 1
)

/**
DB with read,write and ping resolver configured.
RetryInterval at which restore to primary would be tried.
PingTable Name used to ping db with count query to check connectivity.
*/
func NewFallBackPolicy(Db *gorm.DB, retryInterval time.Duration, pingTable string) *FallBackPolicy {
	return &FallBackPolicy{
		currentPool: 0,
		errChan:     make(chan error, 5),
		db:          Db,
		ticker:      time.NewTicker(retryInterval),
		pingTable:   pingTable,
	}
}

/**
Resolve Function Implementation for dbResolver.
*/
func (self *FallBackPolicy) Resolve(connPools []gorm.ConnPool) gorm.ConnPool {
	x := self.GetPool()
	log.WithFields(log.Fields{"Pool Count": len(connPools), "Current Pool": x}).Trace("Pool Info")
	pool := connPools[x]
	return pool
}

/**
Report any DB Errors, can handle nils
*/
func (self *FallBackPolicy) ReportError(err error) {
	self.errChan <- err
}

/**
Returns Current Pool that should be getting Traffic.
On Error switches to fallback and reverts post
successful Ping
*/
func (self *FallBackPolicy) GetPool() (poolIndex int) {
	select {
	case err, ok := <-self.errChan:
		//If Channel is Not Closed Update Pool
		//TODO:Handle only disconnection error
		if ok && err != nil {
			log.WithFields(log.Fields{"Error": err}).Error("Falling Back to Master for Reads")
			//Update New Value in Cache
			self.currentPool = POOL_FALLBACK
		}
		//Serve Updated Pool
		poolIndex = self.currentPool
	case <-self.ticker.C:
		//If we have switched to Fallback ping primary
		if self.currentPool == POOL_FALLBACK {
			//Try to Ping Slave if it is up
			if err := self.Ping(); err == nil {
				self.currentPool = POOL_PRIMARY
				log.Info("Slave Up reverting config for Reads.")
			} else {
				log.WithFields(log.Fields{"Error": err}).Warning("Pinged Slave still not up. Reads continue on master")
			}
			poolIndex = self.currentPool
		}
	default:
		//Serve Current Pool
		poolIndex = self.currentPool
	}
	return
}

func (self *FallBackPolicy) Ping() (err error) {
	c := int64(0)
	err = self.db.Clauses(dbresolver.Use(PING)).Table(self.pingTable).Count(&c).Error
	return
}
