package config

import (
	"errors"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm/logger"
)

type Server struct {
	Host     string       `env:"HOST"`
	Port     int          `env:"PORT" envDefault:"8080"`
	LogLevel logrus.Level `env:"LOG_LEVEL" envDefault:"info"`
}

type RateLimit struct {
	// Skip Redis Host to Disable Rate Limiting
	RedisHost      string `env:"REDIS_RATE_LIMIT"`
	PerMinuteLimit int64  `env:"PORT" envDefault:"50"`
}

func (self *Server) GetUrl(uri string) string {
	return fmt.Sprintf("http://%v:%v%v", self.Host, self.Port, uri)
}

type Vault struct {
	Server `yaml:",inline"`
	Token  string `yaml:"token"`
}

type Db struct {
	//aman:aman@tcp(mysql:3306)/compute?charset=utf8&parseTime=True&loc=Local
	Url string `env:"DB_URL"`
	//TODO: Add Migration Scripts Proper
	//migration_source: /Users/amanpreet.singh/IdeaProjects/Go/go-fun/learn/frameworks/orm/db/go-migrate/migration
	MigrationSource string `env:"DB_MIGRATION_SOURCE"`
	MaxIdle         int    `env:"DB_MAX_IDLE"  envDefault:"2"`
	MaxOpen         int    `env:"DB_MAX_OPEN"  envDefault:"10"`
	AutoMigrate     bool   `env:"DB_AUTO_MIGRATE"  envDefault:"true"`
	//Log level: 4 (Info), 3 (Warn)
	LogLevel logger.LogLevel `env:"DB_LOG_LEVEL"  envDefault:"4"`
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
