package clients

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"

	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
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

var TestHttpClient = resty.New()

func NewHttpClientWithCookies(cookieUrl string, cookies []*http.Cookie, client *resty.Client) *resty.Client {
	cookieJar, _ := cookiejar.New(nil)
	if u, err := url.Parse(cookieUrl); err == nil {
		cookieJar.SetCookies(u, cookies)
	} else {
		log.WithFields(log.Fields{"Error": err}).Error("Error Setting Cookies in HttpClient")
	}
	return client.SetCookieJar(cookieJar)
}
