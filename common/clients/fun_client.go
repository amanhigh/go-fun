// SDK for FunApp with API's for Person Handler using Resty
package clients

import (
	"fmt"

	"github.com/amanhigh/go-fun/models/fun-app/db"
	"github.com/go-resty/resty/v2"
)

type FunClient struct {
	client *resty.Client
}

func NewFunAppClient(BASE_URL string) *FunClient {
	//TODO: Configuration of Http Timeouts
	client := resty.New().SetBaseURL(BASE_URL)
	client.SetHeader("Content-Type", "application/json")

	return &FunClient{
		client: client,
	}
}

func (c *FunClient) GetPerson(name string) (person db.Person, err error) {
	url := fmt.Sprintf("/person/%s", name)
	_, err = c.client.R().SetResult(&person).Get(url)
	return
}
