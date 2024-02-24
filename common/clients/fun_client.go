// SDK for FunApp with API's for Person Handler using Resty
package clients

import (
	"context"
	"fmt"
	"strconv"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/config"
	"github.com/amanhigh/go-fun/models/fun"
	"github.com/go-resty/resty/v2"
)

type FunClient struct {
	PersonService PersonServiceInterface
	AdminService  AdminServiceInterface
}

type PersonServiceInterface interface {
	GetPerson(ctx context.Context, name string) (person fun.Person, err common.HttpError)
	CreatePerson(ctx context.Context, request fun.PersonRequest) (person fun.Person, err common.HttpError)
	UpdatePerson(ctx context.Context, id string, person fun.PersonRequest) (err common.HttpError)
	ListPerson(ctx context.Context, query fun.PersonQuery) (personList fun.PersonList, err common.HttpError)
	DeletePerson(ctx context.Context, name string) (err common.HttpError)
}

type AdminServiceInterface interface {
	Stop(ctx context.Context) (err common.HttpError)
	HealthCheck(ctx context.Context) (err common.HttpError)
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

func (admin *AdminService) Stop(ctx context.Context) (err common.HttpError) {
	response, err1 := admin.client.R().SetContext(ctx).Get("/admin/stop")
	err = util.ResponseProcessor(response, err1)
	return
}

func (admin *AdminService) HealthCheck(ctx context.Context) (err common.HttpError) {
	response, err1 := admin.client.R().SetContext(ctx).Get("/metrics")
	err = util.ResponseProcessor(response, err1)
	return
}

func NewFunAppClient(baseUrl string, httpConfig config.HttpClientConfig) *FunClient {
	client := util.NewRestyClient(baseUrl, httpConfig)

	// Init Base Service
	baseService := BaseService{client: client, VERSION_URL: "/v1"}

	return &FunClient{
		PersonService: &PersonService{BaseService: baseService},
		AdminService:  &AdminService{BaseService: baseService},
	}
}

func (c *PersonService) CreatePerson(ctx context.Context, request fun.PersonRequest) (person fun.Person, err common.HttpError) {
	response, err1 := c.client.R().SetContext(ctx).SetHeader("Content-Type", "application/json").
		SetBody(request).SetResult(&person).Post(c.VERSION_URL + "/person")
	err = util.ResponseProcessor(response, err1)
	return
}

func (c *PersonService) GetPerson(ctx context.Context, name string) (person fun.Person, err common.HttpError) {
	url := fmt.Sprintf(c.VERSION_URL+"/person/%s", name)
	response, err1 := c.client.R().SetContext(ctx).SetResult(&person).Get(url)
	err = util.ResponseProcessor(response, err1)
	return
}

func (c *PersonService) ListPerson(ctx context.Context, personQuery fun.PersonQuery) (personList fun.PersonList, err common.HttpError) {
	response, err1 := c.client.R().SetContext(ctx).SetResult(&personList).Get(c.listPersonUrl(personQuery))
	err = util.ResponseProcessor(response, err1)
	return
}

func (c *PersonService) UpdatePerson(ctx context.Context, id string, person fun.PersonRequest) (err common.HttpError) {
	response, err1 := c.client.R().SetContext(ctx).SetBody(person).Put(fmt.Sprintf(c.VERSION_URL+"/person/%s", id))
	err = util.ResponseProcessor(response, err1)
	return
}

func (c *PersonService) DeletePerson(ctx context.Context, name string) (err common.HttpError) {
	response, err1 := c.request(ctx).Delete(fmt.Sprintf(c.VERSION_URL+"/person/%s", name))
	err = util.ResponseProcessor(response, err1)
	return
}

func (c *PersonService) request(ctx context.Context) *resty.Request {
	return c.client.R().SetContext(ctx).SetError(common.HttpErrorImpl{})
}

// Build Url from personQuery
func (c *PersonService) listPersonUrl(personQuery fun.PersonQuery) (url string) {
	url = c.VERSION_URL + "/person?"

	//Add Pagination Params
	url = url + c.getPaginationParams(personQuery.Offset, personQuery.Limit)

	//Add Sort Params
	if personQuery.SortBy != "" {
		url += "&sort_by=" + personQuery.SortBy
		url += "&order=" + personQuery.Order
	}

	//Add Name and Gender if Provided
	if personQuery.Name != "" {
		url += "&name=" + personQuery.Name
	}
	if personQuery.Gender != "" {
		url += "&gender=" + personQuery.Gender
	}
	return
}
