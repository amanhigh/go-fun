// SDK for FunApp with API's for Person Handler using Resty
package clients

import (
	"fmt"

	"github.com/amanhigh/go-fun/common/helper"
	"github.com/amanhigh/go-fun/models/fun-app/db"
	"github.com/amanhigh/go-fun/models/fun-app/server"
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
func (c *FunClient) CreatePerson(person server.PersonRequest) (err error) {
	response, err := c.client.R().SetBody(person).Post("/person")
	err = helper.ResponseProcessor(response, err)
	return
}
func (c *FunClient) GetPerson(name string) (person db.Person, err error) {
	url := fmt.Sprintf("/person/%s", name)
	response, err := c.client.R().SetResult(&person).Get(url)
	err = helper.ResponseProcessor(response, err)
	return
}

func (c *FunClient) GetAllPersons() (persons []db.Person, err error) {
	response, err := c.client.R().SetResult(&persons).Get("/person/all")
	err = helper.ResponseProcessor(response, err)
	return
}

func (c *FunClient) DeletePerson(name string) (err error) {
	response, err := c.client.R().Delete(fmt.Sprintf("/person/%s", name))
	err = helper.ResponseProcessor(response, err)
	return
}