package util

import (
	"net"
	"net/http"

	"github.com/amanhigh/go-fun/models/config"
	"github.com/dubonzi/otelresty"
	"github.com/go-resty/resty/v2"
)

func NewRestyClient(baseUrl string, httpConfig config.HttpClientConfig) (client *resty.Client) {
	//Init Client
	client = resty.New().SetBaseURL(baseUrl)

	//Default Header
	// client.SetHeader("Content-Type", "application/json")

	//Tracing
	otelresty.TraceClient(client, otelresty.WithTracerName("resty-sdk"))

	//Configure Http Config
	client.SetTimeout(httpConfig.RequestTimeout)

	transport := http.Transport{
		DisableCompression: !httpConfig.Compression,
		DisableKeepAlives:  !httpConfig.KeepAlive,
		DialContext: (&net.Dialer{
			Timeout: httpConfig.DialTimeout, // Connect Timeout
		}).DialContext,
		IdleConnTimeout:     httpConfig.IdleConnectionTimeout, // Idle Timeout Before Closing Keepalive Connection
		MaxIdleConnsPerHost: httpConfig.IdleConnectionsPerHost,
	}
	client.SetTransport(&transport)
	return
}
