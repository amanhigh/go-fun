package config

import (
	"errors"
	"fmt"
	"time"

	"github.com/rs/zerolog"
	"gorm.io/gorm/logger"
)

// https://github.com/caarlos0/env

type Server struct {
	Host     string        `env:"HOST"`
	Port     int           `env:"PORT" envDefault:"8080"`
	LogLevel zerolog.Level `env:"LOG_LEVEL" envDefault:"info"`
}

type RateLimit struct {
	// Skip Redis Host to Disable Rate Limiting
	RedisHost      string `env:"REDIS_RATE_LIMIT"`
	PerMinuteLimit int64  `env:"PER_MIN_LIMIT" envDefault:"-1"`
}

func (self *Server) GetUrl(uri string) string {
	return fmt.Sprintf("http://%v:%v%v", self.Host, self.Port, uri)
}

type Vault struct {
	Server `yaml:",inline"`
	Token  string `yaml:"token"`
}

type Db struct {
	DbType string `env:"DB_TYPE" envDefault:"sqlite"` //mysql,postgres,sqlite
	Url    string `env:"DB_URL" envDefault:"aman:aman@tcp(mysql:3306)/compute?charset=utf8&parseTime=True&loc=Local"`
	//BUG: Add Migration Scripts Proper
	//migration_source: /Users/amanpreet.singh/IdeaProjects/Go/go-fun/learn/frameworks/orm/db/go-migrate/migration
	MigrationSource string `env:"DB_MIGRATION_SOURCE"`
	MaxIdle         int    `env:"DB_MAX_IDLE"  envDefault:"2"`
	MaxOpen         int    `env:"DB_MAX_OPEN"  envDefault:"10"`
	AutoMigrate     bool   `env:"DB_AUTO_MIGRATE"  envDefault:"true"`
	//Log level: 4 (Info), 3 (Warn), 2 (Error)
	LogLevel logger.LogLevel `env:"DB_LOG_LEVEL"  envDefault:"2"`
}

type Tracing struct {
	Type     string `env:"TRACING_TYPE" envDefault:"noop"` // noop,console,otlp
	Endpoint string `env:"TRACING_URL" envDefault:"docker:4317"`
	Publish  string `env:"TRACING_PUBLISH" envDefault:"sync"` //sync, batch (production)
}

var DefaultHttpConfig = HttpClientConfig{
	DialTimeout:            200 * time.Millisecond,
	RequestTimeout:         2 * time.Second,
	IdleConnectionTimeout:  30 * time.Second,
	KeepAlive:              true,
	Compression:            false,
	IdleConnectionsPerHost: 20,
}

type HttpClientConfig struct {
	/* Timeouts */
	DialTimeout           time.Duration `env:"HTTP_DIAL_TIMEOUT" envDefault:"200ms"`
	RequestTimeout        time.Duration `env:"HTTP_REQUEST_TIMEOUT" envDefault:"2s"`
	IdleConnectionTimeout time.Duration `env:"HTTP_IDLE_CONNECTION_TIMEOUT" envDefault:"30s"`

	/* Flags */
	KeepAlive   bool `env:"HTTP_KEEP_ALIVE" envDefault:"true"`
	Compression bool `env:"HTTP_COMPRESSION" envDefault:"false"`

	IdleConnectionsPerHost int `env:"HTTP_IDLE_CONN_PER_HOST" envDefault:"20"`
}

type ZoneMap map[string]Server

func (self ZoneMap) GetUrl(zone, uri string) (url string, err error) {
	if server, ok := self[zone]; ok {
		return server.GetUrl(uri), nil
	} else {
		return "", errors.New("INVALID_ZONE: " + zone)
	}
}
