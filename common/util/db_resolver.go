package util

import (
	"errors"
	"net"
	"time"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

/*
*
Resolver Policy which sends traffic to slave,
incase of error switches to fallback (master).

At given intervals try to revert config back to slave
based on Pings. It needs Ping Table and a Ping dbresolver for it
to work.

Usage:
policy := util.NewFallBackPolicy(db, time.Second * 2,"verticals")

	db.Use(dbresolver.Register(dbresolver.Config{
		Replicas: []gorm.Dialector{
			slave,
			master,
		},
		Policy: policy,
	}).Register(dbresolver.Config{
		Replicas: []gorm.Dialector{
			slave,
		},
	},"ping"))
*/
type FallBackPolicy struct {
	// Current Pool which should be Used
	currentPool int
	// DB to do ping checks with ping resolver configured
	db *gorm.DB
	// Channel to Report Errors
	errChan chan error
	// Interval retry to revert to original Config
	ticker *time.Ticker
	// TableName to get count as a PingCheck on DB
	pingTable string
}

const (
	PING          = "ping" // DBResolver Which Represents Primary for Reads
	POOL_PRIMARY  = 0      // PoolIndex Which Represents Primary
	POOL_FALLBACK = 1      // PoolIndex Which Represents Fallback
)

/*
*
DB with read,write and ping (primary for reads) resolver configured.
RetryInterval at which restore to primary would be tried.
PingTable Name used to ping db with count query to check connectivity.
*/
func NewFallBackPolicy(db *gorm.DB, retryInterval time.Duration, pingTable string) *FallBackPolicy {
	return &FallBackPolicy{
		currentPool: POOL_PRIMARY,
		errChan:     make(chan error, 5),
		db:          db,
		ticker:      time.NewTicker(retryInterval),
		pingTable:   pingTable,
	}
}

/*
*
Resolve Function Implementation for dbResolver.
*/
func (fb *FallBackPolicy) Resolve(connPools []gorm.ConnPool) gorm.ConnPool {
	x := fb.GetPool()
	log.Trace().Int("Pool Count", len(connPools)).Int("Current Pool", x).Msg("Pool Info")
	return connPools[x]
}

/*
*
Report any DB Errors which will be used to fallback if applicable
This ignores any nil or non relevant errors.
*/
func (fb *FallBackPolicy) ReportError(err error) {
	var netErr net.Error
	if errors.As(err, &netErr) {
		fb.errChan <- err
	}
}

/*
*
Returns Current Pool that should be getting Traffic.
On Error switches to fallback and reverts post
successful Ping
*/
func (fb *FallBackPolicy) GetPool() (poolIndex int) {
	select {
	case err, ok := <-fb.errChan:
		// Process if errorChannel Open we are on Primary Pool
		if ok && fb.currentPool == POOL_PRIMARY {
			log.Error().Err(err).Msg("Falling Back to Master for Reads")
			fb.currentPool = POOL_FALLBACK
		}
		// Serve Updated Pool
		poolIndex = fb.currentPool

	case <-fb.ticker.C:
		// At each Retry Interval, If we have switched to Fallback
		if fb.currentPool == POOL_FALLBACK {
			// Try to Ping Slave if it is up
			if err := fb.Ping(); err == nil {
				fb.currentPool = POOL_PRIMARY
				log.Info().Int("Pool", fb.currentPool).Msg("Slave Up reverting config for Reads.")
			} else {
				log.Warn().Int("Pool", fb.currentPool).Err(err).Msg("Pinged Slave still not up. Reads continue on master")
			}
			poolIndex = fb.currentPool
		}
	default:
		//Serve Current Pool
		poolIndex = fb.currentPool
	}
	return
}

/*
*
Triest to do a Count * on configured Ping Table to check connectivity.
Uses Ping Resolver as Read Resolver is switched by our Fallback Policy and is point
to master.
*/
func (fb *FallBackPolicy) Ping() (err error) {
	c := int64(0)
	err = fb.db.Clauses(dbresolver.Use(PING)).Table(fb.pingTable).Count(&c).Error
	return
}
