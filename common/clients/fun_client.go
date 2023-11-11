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
	PersonService *PersonService
}

type BaseService struct {
	client      *resty.Client
	VERSION_URL string
}

type PersonService struct {
	BaseService
}

func NewFunAppClient(BASE_URL string) *FunClient {
	//TODO: Configuration of Http Timeouts
	client := resty.New().SetBaseURL(BASE_URL)
	client.SetHeader("Content-Type", "application/json")

	// Init Base Service
	baseService := BaseService{client: client, VERSION_URL: "/v1"}

	return &FunClient{
		PersonService: &PersonService{BaseService: baseService},
	}
}

func (c *PersonService) CreatePerson(person server.PersonRequest) (err error) {
	response, err := c.client.R().SetBody(person).Post(c.VERSION_URL + "/person")
	err = helper.ResponseProcessor(response, err)
	return
}

func (c *PersonService) GetPerson(name string) (person db.Person, err error) {
	url := fmt.Sprintf(c.VERSION_URL+"/person/%s", name)
	response, err := c.client.R().SetResult(&person).Get(url)
	err = helper.ResponseProcessor(response, err)
	return
}

func (c *PersonService) GetAllPersons() (persons []db.Person, err error) {
	response, err := c.client.R().SetResult(&persons).Get(c.VERSION_URL + "/person/all")
	err = helper.ResponseProcessor(response, err)
	return
}

func (c *PersonService) DeletePerson(name string) (err error) {
	response, err := c.client.R().Delete(fmt.Sprintf(c.VERSION_URL+"/person/%s", name))
	err = helper.ResponseProcessor(response, err)
	return
}
