package config

import (
	"errors"
	"fmt"
	"time"

	"github.com/rs/zerolog"
	"gorm.io/gorm/logger"
)

// https://github.com/caarlos0/env
const (
	LOG_FORMATTER_JSON   = "json"
	LOG_FORMATTER_PRETTY = "pretty"

	// HTTP Configuration Constants
	DefaultDialTimeout            = 200 * time.Millisecond
	DefaultRequestTimeout         = 2 * time.Second
	DefaultIdleConnectionTimeout  = 30 * time.Second
	DefaultIdleTimeout            = 30 * time.Second
	DefaultReadHeaderTimeout      = 2 * time.Second
	DefaultIdleConnectionsPerHost = 20
	DefaultMaxRetries             = 3
	DefaultJitterFactor           = 0.1
	DefaultFailureThreshold       = 5
	DefaultBreakerDelay           = 10 * time.Second
	DefaultSuccessThreshold       = 3
)

var DefaultLogConfig = Log{
	LogLevel:  zerolog.InfoLevel,
	Formatter: LOG_FORMATTER_PRETTY,
}

type Server struct {
	Host string `env:"HOST"`
	Port int    `env:"PORT" envDefault:"8080"`
}

type Log struct {
	LogLevel  zerolog.Level `env:"LOG_LEVEL" envDefault:"info"`
	Formatter string        `env:"LOG_FORMATTER" envDefault:"pretty"` // json,pretty
}

type RateLimit struct {
	// Skip Redis Host to Disable Rate Limiting
	RedisHost      string `env:"REDIS_RATE_LIMIT"`
	PerMinuteLimit int64  `env:"PER_MIN_LIMIT" envDefault:"-1"`
}

func (s *Server) GetUrl(uri string) string {
	return fmt.Sprintf("http://%v:%v%v", s.Host, s.Port, uri)
}

type Vault struct {
	Server `yaml:",inline"`
	Token  string `yaml:"token"`
}

type Db struct {
	DbType          string `env:"DB_TYPE" envDefault:"sqlite"` // mysql,postgres,sqlite
	Url             string `env:"DB_URL" envDefault:"aman:aman@tcp(mysql:3306)/compute?charset=utf8&parseTime=True&loc=Local"`
	MigrationSource string `env:"DB_MIGRATION_SOURCE"`
	MaxIdle         int    `env:"DB_MAX_IDLE"  envDefault:"2"`
	MaxOpen         int    `env:"DB_MAX_OPEN"  envDefault:"10"`
	AutoMigrate     bool   `env:"DB_AUTO_MIGRATE"  envDefault:"true"`
	// Log level: 4 (Info), 3 (Warn), 2 (Error)
	LogLevel logger.LogLevel `env:"DB_LOG_LEVEL"  envDefault:"2"`
}

type Tracing struct {
	Type     string `env:"TRACING_TYPE" envDefault:"noop"` // noop,console,otlp
	Endpoint string `env:"TRACING_URL" envDefault:"docker:4317"`
	Publish  string `env:"TRACING_PUBLISH" envDefault:"batch"` // sync, batch (production)
}

type HttpClientConfig struct {
	/* Timeouts */
	DialTimeout           time.Duration `env:"HTTP_DIAL_TIMEOUT" envDefault:"200ms"`
	RequestTimeout        time.Duration `env:"HTTP_REQUEST_TIMEOUT" envDefault:"2s"`
	IdleConnectionTimeout time.Duration `env:"HTTP_IDLE_CONNECTION_TIMEOUT" envDefault:"30s"`
	ReadTimeout           time.Duration `env:"HTTP_READ_TIMEOUT" envDefault:"1s"`
	WriteTimeout          time.Duration `env:"HTTP_WRITE_TIMEOUT" envDefault:"1s"`
	IdleTimeout           time.Duration `env:"HTTP_IDLE_TIMEOUT" envDefault:"30s"`
	ReadHeaderTimeout     time.Duration `env:"HTTP_READ_HEADER_TIMEOUT" envDefault:"2s"`

	/* Flags */
	KeepAlive   bool `env:"HTTP_KEEP_ALIVE" envDefault:"true"`
	Compression bool `env:"HTTP_COMPRESSION" envDefault:"false"`

	IdleConnectionsPerHost int `env:"HTTP_IDLE_CONN_PER_HOST" envDefault:"20"`

	Failsafe FailsafeConfig
}

type FailsafeConfig struct {
	Retry   RetryConfig
	Breaker BreakerConfig
}

type RetryConfig struct {
	MaxRetries   int           `env:"HTTP_RETRY_MAX" envDefault:"3"`
	Delay        time.Duration `env:"HTTP_RETRY_DELAY" envDefault:"1s"`
	JitterFactor float32       `env:"HTTP_RETRY_JITTER" envDefault:"0.1"`
}

type BreakerConfig struct {
	Delay            time.Duration `env:"HTTP_BREAKER_DELAY" envDefault:"10s"`
	FailureThreshold uint          `env:"HTTP_BREAKER_FAILURE_THRESHOLD" envDefault:"5"`
	SuccessThreshold uint          `env:"HTTP_BREAKER_SUCCESS_THRESHOLD" envDefault:"3"`
}

var DefaultHttpConfig = HttpClientConfig{
	DialTimeout:            DefaultDialTimeout,
	RequestTimeout:         DefaultRequestTimeout,
	IdleConnectionTimeout:  DefaultIdleConnectionTimeout,
	ReadTimeout:            1 * time.Second,
	WriteTimeout:           1 * time.Second,
	IdleTimeout:            DefaultIdleTimeout,
	ReadHeaderTimeout:      DefaultReadHeaderTimeout,
	KeepAlive:              true,
	Compression:            false,
	IdleConnectionsPerHost: DefaultIdleConnectionsPerHost,
	Failsafe: FailsafeConfig{
		Retry: RetryConfig{
			MaxRetries:   DefaultMaxRetries,
			Delay:        time.Second,
			JitterFactor: DefaultJitterFactor,
		},
		Breaker: BreakerConfig{
			FailureThreshold: DefaultFailureThreshold,
			Delay:            DefaultBreakerDelay,
			SuccessThreshold: DefaultSuccessThreshold,
		},
	},
}

type ZoneMap map[string]Server

func (s ZoneMap) GetUrl(zone, uri string) (url string, err error) {
	if server, ok := s[zone]; ok {
		return server.GetUrl(uri), nil
	}
	return "", errors.New("INVALID_ZONE: " + zone)
}
