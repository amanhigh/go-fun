package util

import (
	"net"
	"net/http"

	. "github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/config"
	"github.com/dubonzi/otelresty"
	"github.com/go-resty/resty/v2"
)

// ResponseProcessor processes the resty response and error to generate an HttpError.
//
// Parameters:
// - response: a pointer to resty.Response
// - restyErr: an error
// Return type(s):
// - err: an HttpError
func ResponseProcessor(response *resty.Response, restyErr error) (err HttpError) {
	var ok bool
	if restyErr != nil {
		//Rest Client Error hence No Respones
		err = NewServerError(restyErr)
	} else if err, ok = response.Error().(HttpError); ok && err.Code() > 0 {
		//If Error is Http Error & has Data, Use directly.
	} else {
		//Incase we have No Error Honor Status Codes of Http
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
