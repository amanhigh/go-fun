package clients

import (
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"

	"github.com/amanhigh/go-fun/models/config"
	"github.com/dubonzi/otelresty"
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
)

const (
	DIAL_TIMEOUT     = 2 * time.Second
	REQUEST_TIMEOUT  = 10 * time.Second
	IDLE_TIMEOUT     = time.Minute
	IDLE_CONNECTIONS = 64
)

var DefaultHttpClient = resty.New().SetTimeout(REQUEST_TIMEOUT).SetTransport(&http.Transport{
	IdleConnTimeout:    IDLE_TIMEOUT,
	MaxIdleConns:       IDLE_CONNECTIONS,
	DisableKeepAlives:  false,
	DisableCompression: true,
})

func NewHttpClientWithCookies(cookieUrl string, cookies []*http.Cookie, client *resty.Client) *resty.Client {
	cookieJar, _ := cookiejar.New(nil)
	if u, err := url.Parse(cookieUrl); err == nil {
		cookieJar.SetCookies(u, cookies)
	} else {
		log.Error().Err(err).Msg("Error Setting Cookies in HttpClient")
	}
	return client.SetCookieJar(cookieJar)
}

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
