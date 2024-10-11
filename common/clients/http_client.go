package clients

import (
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"

	"github.com/amanhigh/go-fun/models/config"
	"github.com/dubonzi/otelresty"
	"github.com/failsafe-go/failsafe-go/failsafehttp"
	"github.com/failsafe-go/failsafe-go/retrypolicy"
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

	// Set the transport using the new transport builder
	client.SetTransport(buildNewTransport(httpConfig))

	return
}

func buildNewTransport(httpConfig config.HttpClientConfig) http.RoundTripper {
	transport := &http.Transport{
		DisableCompression: !httpConfig.Compression,
		DisableKeepAlives:  !httpConfig.KeepAlive,
		DialContext: (&net.Dialer{
			Timeout: httpConfig.DialTimeout, // Connect Timeout
		}).DialContext,
		IdleConnTimeout:     httpConfig.IdleConnectionTimeout, // Idle Timeout Before Closing Keepalive Connection
		MaxIdleConnsPerHost: httpConfig.IdleConnectionsPerHost,
	}

	if httpConfig.Retries > 0 {
		retryPolicy := buildRetryPolicy(httpConfig.Retries)
		return failsafehttp.NewRoundTripper(transport, retryPolicy)
	}

	return transport
}

func buildRetryPolicy(retries int) retrypolicy.RetryPolicy[*http.Response] {
	return failsafehttp.RetryPolicyBuilder().
		WithDelay(time.Second).
		WithJitterFactor(0.1). //Avoid Thundering Hurd
		WithMaxRetries(retries).
		Build()
}
