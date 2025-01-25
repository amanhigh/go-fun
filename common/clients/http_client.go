package clients

import (
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"

	"github.com/amanhigh/go-fun/models/config"
	"github.com/dubonzi/otelresty"
	"github.com/failsafe-go/failsafe-go/circuitbreaker"
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
	// Init Client
	client = resty.New().SetBaseURL(baseUrl)

	// Default Header
	// client.SetHeader("Content-Type", "application/json")
	// Tracing
	otelresty.TraceClient(client, otelresty.WithTracerName("resty-sdk"))

	// Configure Http Config
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

	retryPolicy := buildRetryPolicy(httpConfig.Failsafe.Retry)
	circuitBreaker := buildCircuitBreaker(httpConfig.Failsafe.Breaker)

	return failsafehttp.NewRoundTripper(transport, retryPolicy, circuitBreaker)
}

func buildRetryPolicy(config config.RetryConfig) retrypolicy.RetryPolicy[*http.Response] {
	// TASK: Implement Metrics
	return failsafehttp.RetryPolicyBuilder().
		WithDelay(config.Delay).
		WithJitterFactor(config.JitterFactor).
		WithMaxRetries(config.MaxRetries).
		Build()
}

/*
*
When the number of recent execution failures exceed a configured threshold, the breaker is opened
and further executions will fail with circuitbreaker.ErrOpen.
After a delay, the breaker is half-opened and trial executions are allowed which determine
whether the breaker should be closed or opened again. If the trial executions meet a
success threshold, the breaker is closed again and executions will proceed as normal, otherwise it’s re-opened.
*/
func buildCircuitBreaker(config config.BreakerConfig) circuitbreaker.CircuitBreaker[*http.Response] {
	return circuitbreaker.Builder[*http.Response]().
		WithFailureThreshold(config.FailureThreshold).
		WithDelay(config.Delay).
		WithSuccessThreshold(config.SuccessThreshold).
		Build()
}
