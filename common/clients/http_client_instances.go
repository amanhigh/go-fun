package clients

import (
	config2 "github.com/amanhigh/go-fun/models/config"
	"time"
)

const (
	DIAL_TIMEOUT     = 2 * time.Second
	REQUEST_TIMEOUT  = 10 * time.Second
	IDLE_TIMEOUT     = time.Minute
	IDLE_CONNECTIONS = 64
)

var (
	DefaultHttpClientConfig = config2.HttpClientConfig{
		DialTimeout:            DIAL_TIMEOUT,
		RequestTimeout:         REQUEST_TIMEOUT,
		IdleConnectionTimeout:  IDLE_TIMEOUT,
		IdleConnectionsPerHost: IDLE_CONNECTIONS,
		KeepAlive:              true,
		Compression:            false,
	}
)

var TestHttpClient = NewHttpClient(config2.HttpClientConfig{
	DialTimeout:            time.Millisecond * 200,
	RequestTimeout:         time.Second,
	IdleConnectionTimeout:  time.Second * 5,
	KeepAlive:              true,
	Compression:            false,
	IdleConnectionsPerHost: 5,
})
