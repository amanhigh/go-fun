//go:generate mockgen -package util -destination http_mock.go -source http.go
package util

import (
	"bytes"
	"context"
	"crypto/tls"
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

	log "github.com/Sirupsen/logrus"
)

const (
	DIAL_TIMEOUT     = 2 * time.Second
	REQUEST_TIMEOUT  = 10 * time.Second
	IDLE_CONNECTIONS = 64
)

var NoKeepAliveClient = NewHttpClient(DIAL_TIMEOUT, REQUEST_TIMEOUT, false, -1, false, false)
var KeepAliveClient = NewHttpClient(DIAL_TIMEOUT, REQUEST_TIMEOUT, true, 64, false, false)
var KeepAliveInsecureClient = NewHttpClient(DIAL_TIMEOUT, REQUEST_TIMEOUT, true, 64, false, true)

/**
returns:
200 - Decoded Response (if interface provided), status code, no error
Non 200 - Nil Response, status code, error (Non 200 Response)
Deserialization Failed - Nil Response, status code, deserialization error
Http Error - Nil Response, Zero Status Code, error (Http Error that occurred)
*/
type HttpClientInterface interface {
	DoGet(url string, unmarshalledResponse interface{}) (statusCode int, err error)
	DoPost(url string, body interface{}, unmarshalledResponse interface{}) (statusCode int, err error)
	DoPut(url string, body interface{}, unmarshalledResponse interface{}) (statusCode int, err error)
	DoDelete(url string, body interface{}, unmarshalledResponse interface{}) (statusCode int, err error)

	DoGetWithTimeout(url string, unmarshalledResponse interface{}, timeout time.Duration) (statusCode int, err error)
	DoPostWithTimeout(url string, body interface{}, unmarshalledResponse interface{}, timout time.Duration) (statusCode int, err error)

	DoRequest(request *http.Request, unmarshalledResponse interface{}, timeout time.Duration) (statusCode int, err error)
	AddHeader(key string, value string)
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
func NewHttpClient(dialTimeout time.Duration, requestTimeout time.Duration, enableKeepAlive bool, idleConnectionsPerHost int, enableCompression bool, allowInsecure bool) HttpClientInterface {
	jar, _ := cookiejar.New(nil)
	return &HttpClient{
		Client: &http.Client{
			Jar: jar,
			Transport: &http.Transport{
				DisableCompression: !enableCompression,
				DisableKeepAlives:  !enableKeepAlive,
				DialContext: (&net.Dialer{
					Timeout: dialTimeout, // Connect Timeout
				}).DialContext,
				MaxIdleConnsPerHost: idleConnectionsPerHost,
				IdleConnTimeout:     (dialTimeout + requestTimeout) * 120, //Idle Timeout Before Closing Keepalive Connection
				TLSClientConfig:     &tls.Config{InsecureSkipVerify: allowInsecure},
			},
			CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
				return errors.New("It is a Redirect!")
			},
		},
		Timeout: requestTimeout, //Request Timeout
	}
}

func NewHttpClientWithCookies(cookieUrl string, cookies []*http.Cookie, keepAlive bool, compression bool) HttpClientInterface {
	client := NewHttpClient(DIAL_TIMEOUT, REQUEST_TIMEOUT, keepAlive, IDLE_CONNECTIONS, compression, false).(*HttpClient)
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
Given a request and unmarshal body, fire Http Client Return Unmarshalled Response
*/
func (self *HttpClient) DoRequest(request *http.Request, unmarshalledResponse interface{}, timeout time.Duration) (statusCode int, err error) {
	var responseBytes []byte
	var response *http.Response

	timeoutContext, cancelFunction := context.WithTimeout(context.Background(), self.getTimeOut(timeout))
	/* Set Content Type Header */
	request.Header.Set("Content-Type", "application/json")
	/* Execute Request */
	defer cancelFunction()
	if response, err = self.Client.Do(request.WithContext(timeoutContext)); err == nil {
		defer response.Body.Close()

		/* Check If Request was Successful */
		statusCode = response.StatusCode

		/* Decode Body if 200 else throw error */
		if response.StatusCode == http.StatusOK {
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
			/* Discard body if not read */
			io.Copy(ioutil.Discard, response.Body)
			err = errors.New(fmt.Sprintf("Non 200 Response. Status Code: %v", response.StatusCode))
		}
	}

	return
}

func (self *HttpClient) fireRequest(method string, url string, body interface{}, unmarshalledResponse interface{}, timeout time.Duration) (statusCode int, err error) {
	var requestBody []byte
	var request *http.Request

	/* Encode Json */
	if requestBody, err = json.Marshal(body); err == nil {
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
Add header required for the request
*/
func (self *HttpClient) AddHeader(key string, value string) {
	self.HeaderMap[key] = value
}
