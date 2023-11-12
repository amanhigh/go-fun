// SDK for FunApp with API's for Person Handler using Resty
package clients

import (
	"fmt"
	"strconv"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/models/common"
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

func (self *BaseService) getPaginationParams(offset, limit int) (query string) {
	return "offset=" + strconv.Itoa(offset) + "&limit=" + strconv.Itoa(limit)
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

func (c *PersonService) CreatePerson(person server.PersonRequest) (id string, err common.HttpError) {
	response, err1 := c.client.R().SetBody(person).SetResult(&id).Post(c.VERSION_URL + "/person")
	err = util.ResponseProcessor(response, err1)
	return
}

func (c *PersonService) GetPerson(name string) (person db.Person, err common.HttpError) {
	url := fmt.Sprintf(c.VERSION_URL+"/person/%s", name)
	response, err1 := c.client.R().SetResult(&person).Get(url)
	err = util.ResponseProcessor(response, err1)
	return
}

func (c *PersonService) ListPerson(personQuery server.PersonQuery) (personList server.PersonList, err common.HttpError) {
	url := c.VERSION_URL + "/person?" + c.getPaginationParams(personQuery.Offset, personQuery.Limit)
	response, err1 := c.client.R().SetResult(&personList).Get(url)
	err = util.ResponseProcessor(response, err1)
	return
}

func (c *PersonService) DeletePerson(name string) (err common.HttpError) {
	response, err1 := c.client.R().Delete(fmt.Sprintf(c.VERSION_URL+"/person/%s", name))
	err = util.ResponseProcessor(response, err1)
	return
}
