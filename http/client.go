package http

import (
	"time"
	"net"
	"encoding/json"
	"bytes"
	"io/ioutil"
	"fmt"
	"net/http"
	"context"
	"errors"
)

var NoKeepAliveClient = BuildNonKeepAliveClient()

var KeepAliveClient = BuildKeepAliveClient(60*time.Second, 10)

type HttpClient struct {
	Client *http.Client
}

/* Constructors */
func BuildNonKeepAliveClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			DisableCompression:  true,
			DisableKeepAlives:   true,
			MaxIdleConnsPerHost: -1,
		},
	}
}

func BuildKeepAliveClient(dialTimeout time.Duration, idleConnectionsPerHost int) *HttpClient {
	return &HttpClient{
		Client: &http.Client{
			Transport: &http.Transport{
				Dial: (&net.Dialer{
					Timeout:   dialTimeout,
					KeepAlive: dialTimeout * 60,
				}).Dial,
				MaxIdleConnsPerHost: idleConnectionsPerHost,
			},
		},
	}
}

/* Makes a Post Request with Given Url & Body under specified timeout.

Incase of Success you will recieve unmarshalled Response or error otherwise
 */
func (httpClient *HttpClient) DoPost(url string, body interface{}, unmarshalledResponse interface{}, timeout time.Duration) (err error) {
	var jobData, responseBytes []byte
	var request *http.Request
	var response *http.Response

	/* Encode Json */
	if jobData, err = json.Marshal(body); err == nil {

		/* Build Request */
		if request, err = http.NewRequest("POST", url, bytes.NewReader(jobData)); err == nil {
			timeoutContext, cancelFunction := context.WithTimeout(context.Background(), timeout)

			/* Set Content Type Header */
			request.Header.Set("Content-Type", "application/json")

			/* Execute Request */
			defer cancelFunction()
			if response, err = httpClient.Client.Do(request.WithContext(timeoutContext)); err == nil {
				defer response.Body.Close()

				/* Check If Delegation was Successful */
				if response.StatusCode == http.StatusOK {
					/* Read Body & Decode */
					if responseBytes, err = ioutil.ReadAll(response.Body); err == nil {
						err = json.Unmarshal(responseBytes, unmarshalledResponse)
					}
				} else {
					err = errors.New(fmt.Sprintf("Non 200 Delegation Response. Status Code: %v", response.StatusCode))
				}
			}
		}
	}

	return
}
