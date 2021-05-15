//go:generate mockgen -package http -destination http_client_mock.go -source http_client.go
package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"

	"github.com/amanhigh/go-fun/apps/models/config"
	log "github.com/sirupsen/logrus"
)

/**
returns:
200 - Decoded Response (if interface provided), status code, no error
Non 200 - Nil Response, status code, error (Non 200 Response)
Deserialization Failed - Nil Response, status code, deserialization error
Http Error - Nil Response, Zero Status Code, error (Http Error that occurred)
*/

const NON_2xx_RESPONSE = "NON_2xx_RESPONSE"

type HttpClientInterface interface {
	DoGet(url string, unmarshalledResponse interface{}) (statusCode int, err error)
	DoPost(url string, body interface{}, unmarshalledResponse interface{}) (statusCode int, err error)
	DoPut(url string, body interface{}, unmarshalledResponse interface{}) (statusCode int, err error)
	DoDelete(url string, body interface{}, unmarshalledResponse interface{}) (statusCode int, err error)

	DoGetWithTimeout(url string, unmarshalledResponse interface{}, timeout time.Duration) (statusCode int, err error)
	DoPostWithTimeout(url string, body interface{}, unmarshalledResponse interface{}, timout time.Duration) (statusCode int, err error)
	DoPutWithTimeout(url string, body interface{}, unmarshalledResponse interface{}, timout time.Duration) (statusCode int, err error)

	DoRequest(request *http.Request, unmarshalledResponse interface{}, timeout time.Duration) (statusCode int, err error)
	SetHeaderMap(headerMap map[string]string)
}

type HttpClient struct {
	Client    *http.Client
	Timeout   time.Duration
	HeaderMap map[string]string
}

/**
dialTimeout: Connect Timeout
requestTimeout: Time allowed for an Http Request

KeepAlive Parameters:
keepAlive: Enable/Disable Keep Alive
idleConnectionsPerHost: Can be -1 if no keep alive or number of Max Idle KeepAlive connections to keep in pool.

enableCompression: Enable/Disable gzip compression
*/
func NewHttpClient(httpClientConfig config.HttpClientConfig) HttpClientInterface {
	log.WithFields(log.Fields{"Config": httpClientConfig}).Debug("CREATE_HTTP_CLIENT")

	defaultHeader := map[string]string{
		"Content-Type": "application/json",
	}
	jar, _ := cookiejar.New(nil)
	return &HttpClient{
		Client: &http.Client{
			Jar: jar,
			Transport: &http.Transport{
				DisableCompression: !httpClientConfig.Compression,
				DisableKeepAlives:  !httpClientConfig.KeepAlive,
				DialContext: (&net.Dialer{
					Timeout: httpClientConfig.DialTimeout, // Connect Timeout
				}).DialContext,
				IdleConnTimeout:     httpClientConfig.IdleConnectionTimeout, //Idle Timeout Before Closing Keepalive Connection
				MaxIdleConnsPerHost: httpClientConfig.IdleConnectionsPerHost,
			},
		},
		Timeout:   httpClientConfig.RequestTimeout, //Request Timeout
		HeaderMap: defaultHeader,
	}
}

func NewHttpClientWithCookies(cookieUrl string, cookies []*http.Cookie, config config.HttpClientConfig) HttpClientInterface {
	client := NewHttpClient(config).(*HttpClient)
	if u, err := url.Parse(cookieUrl); err == nil {
		client.Client.Jar.SetCookies(u, cookies)
	} else {
		log.WithFields(log.Fields{"Error": err}).Error("")
	}
	return client
}

/*
	Makes a Get Request & Unmarshalles Response into unmarshalledResponse if
	provided else returns Status Code.

	Incase of Error returns error that occured
*/
func (self *HttpClient) DoGet(url string, unmarshalledResponse interface{}) (statusCode int, err error) {
	return self.fireRequest("GET", url, nil, unmarshalledResponse, -1)
}

/*
	Makes a Post Request with Given Url & Body under specified timeout.
	Incase of Success you will recieve unmarshalled Response or error otherwise
*/
func (self *HttpClient) DoPost(url string, body interface{}, unmarshalledResponse interface{}) (statusCode int, err error) {
	return self.fireRequest("POST", url, body, unmarshalledResponse, -1)
}

/*
	Makes a Post Request with Given Url & Body under specified timeout.
	Incase of Success you will recieve unmarshalled Response or error otherwise
*/
func (self *HttpClient) DoPut(url string, body interface{}, unmarshalledResponse interface{}) (statusCode int, err error) {
	return self.fireRequest("PUT", url, body, unmarshalledResponse, -1)
}

/*
	Makes a Post Request with Given Url & Body under specified timeout.
	Incase of Success you will recieve unmarshalled Response or error otherwise
*/
func (self *HttpClient) DoDelete(url string, body interface{}, unmarshalledResponse interface{}) (statusCode int, err error) {
	return self.fireRequest("DELETE", url, body, unmarshalledResponse, -1)
}

/**
Ignores Global Timeout of HttpClient and uses provided timeout fo Http call.
*/
func (self *HttpClient) DoGetWithTimeout(url string, unmarshalledResponse interface{}, timeout time.Duration) (statusCode int, err error) {
	return self.fireRequest("GET", url, nil, unmarshalledResponse, timeout)
}

/**
Ignores Global Timeout of HttpClient and uses provided timeout fo Http call.
*/
func (self *HttpClient) DoPostWithTimeout(url string, body interface{}, unmarshalledResponse interface{}, timeout time.Duration) (statusCode int, err error) {
	return self.fireRequest("POST", url, body, unmarshalledResponse, timeout)
}

/**
Ignores Global Timeout of HttpClient and uses provided timeout fo Http call.
*/
func (self *HttpClient) DoPutWithTimeout(url string, body interface{}, unmarshalledResponse interface{}, timeout time.Duration) (statusCode int, err error) {
	return self.fireRequest("PUT", url, body, unmarshalledResponse, timeout)
}

/**
Given a request and unmarshal body, fire Http Client Return Unmarshalled Response
*/
func (self *HttpClient) DoRequest(request *http.Request, unmarshalledResponse interface{}, timeout time.Duration) (statusCode int, err error) {
	var responseBytes []byte
	var response *http.Response

	timeoutContext, cancelFunction := context.WithTimeout(context.Background(), self.getTimeOut(timeout))
	/* Set Header from HeaderMap */
	for key, value := range self.HeaderMap {
		request.Header.Set(key, value)
	}
	/* Execute Request */
	defer cancelFunction()
	if response, err = self.Client.Do(request.WithContext(timeoutContext)); err == nil {
		defer response.Body.Close()

		/* Check If Request was Successful */
		statusCode = response.StatusCode

		/* Decode Body if 2xx else throw error */
		if http.StatusOK <= statusCode && statusCode <= http.StatusMultipleChoices {
			if unmarshalledResponse != nil {
				/* Read Body & Decode if Response came & unmarshal entity is supplied */
				if responseBytes, err = ioutil.ReadAll(response.Body); err == nil {

					/* Return If its string else unmarshal Body */
					if stringInterface, ok := unmarshalledResponse.(*string); ok {
						*stringInterface = string(responseBytes)
					} else {
						err = json.Unmarshal(responseBytes, unmarshalledResponse)
					}
				}
			} else {
				/* Discard body if not read */
				io.Copy(ioutil.Discard, response.Body)
			}
		} else {
			/* Read Body and Print as Part of Error Message */
			if responseBytes, err = ioutil.ReadAll(response.Body); err == nil {
				err = errors.New(fmt.Sprintf("%v STATUS_CODE: %v RESPONSE: %v", NON_2xx_RESPONSE, response.StatusCode, string(responseBytes)))
			}
		}
	}

	return
}

func (self *HttpClient) fireRequest(method string, url string, body interface{}, unmarshalledResponse interface{}, timeout time.Duration) (statusCode int, err error) {
	var requestBody []byte
	var request *http.Request

	/* Return If its string else marshal Body */
	if val, ok := body.(string); ok {
		requestBody = []byte(val)
	} else if body == nil {
		requestBody = nil
	} else {
		requestBody, err = json.Marshal(body)
	}

	/* Check Encode Json Error*/
	if err == nil {
		/* Build Request */
		if request, err = http.NewRequest(method, url, bytes.NewReader(requestBody)); err == nil {
			return self.DoRequest(request, unmarshalledResponse, timeout)
		}
	}
	return
}

/**
Returns timeout if its non Zero or Client Level Timeout otherwise
*/
func (self *HttpClient) getTimeOut(timeout time.Duration) time.Duration {
	if timeout < 0 {
		return self.Timeout
	} else {
		return timeout
	}

}

/**
Set Custom header map for this client, that will be added in all requests by default.
*/
func (self *HttpClient) SetHeaderMap(headerMap map[string]string) {
	self.HeaderMap = headerMap
}
