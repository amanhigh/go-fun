package clients

import (
	"net/http"
	"time"

	config2 "github.com/amanhigh/go-fun/models/config"
	"github.com/go-resty/resty/v2"
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

var DefaultHttpClient = resty.New().SetTimeout(REQUEST_TIMEOUT).SetTransport(&http.Transport{
	IdleConnTimeout:    IDLE_TIMEOUT,
	MaxIdleConns:       IDLE_CONNECTIONS,
	DisableKeepAlives:  false,
	DisableCompression: true,
})

var TestHttpClient = resty.New().SetHeader("Content-Type", "application/json")
