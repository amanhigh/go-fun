// SDK for FunApp with API's for Person Handler using Resty
package clients

import (
	"fmt"
	"strconv"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/fun"
	"github.com/go-resty/resty/v2"
)

type FunClient struct {
	PersonService *PersonService
	AdminService  *AdminService
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

type AdminService struct {
	BaseService
}

func (admin *AdminService) Stop() (err common.HttpError) {
	response, err1 := admin.client.R().Get("/admin/stop")
	err = util.ResponseProcessor(response, err1)
	return
}

func (admin *AdminService) HealthCheck() (err common.HttpError) {
	response, err1 := admin.client.R().Get("/metrics")
	err = util.ResponseProcessor(response, err1)
	return
}

func NewFunAppClient(BASE_URL string) *FunClient {
	//TODO: Configuration of Http Timeouts
	client := resty.New().SetBaseURL(BASE_URL)
	client.SetHeader("Content-Type", "application/json")

	// Init Base Service
	baseService := BaseService{client: client, VERSION_URL: "/v1"}

	return &FunClient{
		PersonService: &PersonService{BaseService: baseService},
		AdminService:  &AdminService{BaseService: baseService},
	}
}

func (c *PersonService) CreatePerson(person fun.PersonRequest) (id string, err common.HttpError) {
	response, err1 := c.client.R().SetBody(person).SetResult(&id).Post(c.VERSION_URL + "/person")
	err = util.ResponseProcessor(response, err1)
	return
}

func (c *PersonService) GetPerson(name string) (person fun.Person, err common.HttpError) {
	url := fmt.Sprintf(c.VERSION_URL+"/person/%s", name)
	response, err1 := c.client.R().SetResult(&person).Get(url)
	err = util.ResponseProcessor(response, err1)
	return
}

func (c *PersonService) ListPerson(personQuery fun.PersonQuery) (personList fun.PersonList, err common.HttpError) {
	response, err1 := c.client.R().SetResult(&personList).Get(c.listPersonUrl(personQuery))
	err = util.ResponseProcessor(response, err1)
	return
}

func (c *PersonService) UpdatePerson(id string, person fun.PersonRequest) (err common.HttpError) {
	response, err1 := c.client.R().SetBody(person).Put(fmt.Sprintf(c.VERSION_URL+"/person/%s", id))
	err = util.ResponseProcessor(response, err1)
	return
}

func (c *PersonService) DeletePerson(name string) (err common.HttpError) {
	response, err1 := c.client.R().Delete(fmt.Sprintf(c.VERSION_URL+"/person/%s", name))
	err = util.ResponseProcessor(response, err1)
	return
}

// Build Url from personQuery
func (c *PersonService) listPersonUrl(personQuery fun.PersonQuery) (url string) {
	url = c.VERSION_URL + "/person?"

	//Add Pagination Params
	url = url + c.getPaginationParams(personQuery.Offset, personQuery.Limit)

	//Add Name and Gender if Provided
	if personQuery.Name != "" {
		url += "&name=" + personQuery.Name
	}
	if personQuery.Gender != "" {
		url += "&gender=" + personQuery.Gender
	}
	return
}
