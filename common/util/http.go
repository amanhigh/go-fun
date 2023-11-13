package util

import (
	"net"
	"net/http"

	. "github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/config"
	"github.com/go-resty/resty/v2"
)

// Error Proccessor Mapping Http Code to Http Error
func ResponseProcessor(response *resty.Response, restyErr error) (err HttpError) {
	if restyErr != nil {
		err = NewServerError(restyErr)
	} else {
		switch response.StatusCode() {
		case http.StatusBadRequest:
			err = ErrBadRequest
		case http.StatusNotFound:
			err = ErrNotFound
		case http.StatusUnauthorized:
			err = ErrNotAuthorized
		case http.StatusForbidden:
			err = ErrNotAuthenticated
		default:
			err = nil
		}
	}
	return
}

func NewRestyClient(baseUrl string, httpConfig config.HttpClientConfig) (client *resty.Client) {
	//Init Client
	client = resty.New().SetBaseURL(baseUrl)

	//Default Header
	client.SetHeader("Content-Type", "application/json")

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
