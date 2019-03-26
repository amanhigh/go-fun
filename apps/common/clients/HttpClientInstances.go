package clients

import (
	"time"

	"github.com/amanhigh/go-fun/apps/models/config"
)

var TestHttpClient = NewHttpClient(config.HttpClientConfig{
	DialTimeout:            time.Millisecond * 200,
	RequestTimeout:         time.Second,
	IdleConnectionTimeout:  time.Second * 5,
	KeepAlive:              true,
	Compression:            false,
	IdleConnectionsPerHost: 5,
})
