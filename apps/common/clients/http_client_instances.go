package clients

import (
	"time"

	"github.com/amanhigh/go-fun/apps/models/config"
)

const (
	DIAL_TIMEOUT     = 2 * time.Second
	REQUEST_TIMEOUT  = 10 * time.Second
	IDLE_TIMEOUT     = time.Minute
	IDLE_CONNECTIONS = 64
)

var (
	DefaultHttpClientConfig = config.HttpClientConfig{
		DialTimeout:            DIAL_TIMEOUT,
		RequestTimeout:         REQUEST_TIMEOUT,
		IdleConnectionTimeout:  IDLE_TIMEOUT,
		IdleConnectionsPerHost: IDLE_CONNECTIONS,
		KeepAlive:              true,
		Compression:            false,
	}
)

var TestHttpClient = NewHttpClient(config.HttpClientConfig{
	DialTimeout:            time.Millisecond * 200,
	RequestTimeout:         time.Second,
	IdleConnectionTimeout:  time.Second * 5,
	KeepAlive:              true,
	Compression:            false,
	IdleConnectionsPerHost: 5,
})
