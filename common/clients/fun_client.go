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
	ListPersonAudit(ctx context.Context, id string) (personAuditList []fun.PersonAudit, err common.HttpError)
	DeletePerson(ctx context.Context, name string) (err common.HttpError)
}

type AdminServiceInterface interface {
	Stop(ctx context.Context) (err common.HttpError)
	HealthCheck(ctx context.Context) (err common.HttpError)
}

type BaseService struct {
	client     *resty.Client
	VersionUrl string
}

// Takes offset and limit as parameters and returns a query string.
func (bs *BaseService) getPaginationParams(offset, limit int) (query string) {
	return "offset=" + strconv.Itoa(offset) + "&limit=" + strconv.Itoa(limit)
}

/*
Builds Base Request for REST Interaction.

@param ctx - Context

Return type(s):
- *resty.Request
*/
func (bs *BaseService) request(ctx context.Context) *resty.Request {
	return bs.client.R().SetContext(ctx).SetError(common.HttpErrorImpl{})
}

type PersonService struct {
	BaseService
}

type AdminService struct {
	BaseService
}

func (as *AdminService) Stop(ctx context.Context) (err common.HttpError) {
	response, err1 := as.client.R().SetContext(ctx).Get("/admin/stop")
	err = util.ResponseProcessor(response, err1)
	return
}

func (as *AdminService) HealthCheck(ctx context.Context) (err common.HttpError) {
	response, err1 := as.client.R().SetContext(ctx).Get("/metrics")
	err = util.ResponseProcessor(response, err1)
	return
}

func NewFunAppClient(baseUrl string, httpConfig config.HttpClientConfig) *FunClient {
	client := NewRestyClient(baseUrl, httpConfig)

	// Init Base Service
	baseService := BaseService{client: client, VersionUrl: "/v1"}

	return &FunClient{
		PersonService: &PersonService{BaseService: baseService},
		AdminService:  &AdminService{BaseService: baseService},
	}
}

func (c *PersonService) CreatePerson(ctx context.Context, request fun.PersonRequest) (person fun.Person, err common.HttpError) {
	response, err1 := c.request(ctx).SetHeader("Content-Type", "application/json").
		SetBody(request).SetResult(&person).Post(c.VersionUrl + "/person")
	err = util.ResponseProcessor(response, err1)
	return
}

func (c *PersonService) GetPerson(ctx context.Context, name string) (person fun.Person, err common.HttpError) {
	url := fmt.Sprintf(c.VersionUrl+"/person/%s", name)
	response, err1 := c.request(ctx).SetResult(&person).Get(url)
	err = util.ResponseProcessor(response, err1)
	return
}

func (c *PersonService) ListPerson(ctx context.Context, personQuery fun.PersonQuery) (personList fun.PersonList, err common.HttpError) {
	response, err1 := c.request(ctx).SetResult(&personList).Get(c.listPersonUrl(personQuery))
	err = util.ResponseProcessor(response, err1)
	return
}

func (c *PersonService) ListPersonAudit(ctx context.Context, id string) (personAuditList []fun.PersonAudit, err common.HttpError) {
	response, err1 := c.request(ctx).SetResult(&personAuditList).Get(fmt.Sprintf(c.VersionUrl+"/person/%s/audit", id))
	err = util.ResponseProcessor(response, err1)
	return
}

func (c *PersonService) UpdatePerson(ctx context.Context, id string, person fun.PersonRequest) (err common.HttpError) {
	response, err1 := c.request(ctx).SetBody(person).Put(fmt.Sprintf(c.VersionUrl+"/person/%s", id))
	err = util.ResponseProcessor(response, err1)
	return
}

func (c *PersonService) DeletePerson(ctx context.Context, name string) (err common.HttpError) {
	response, err1 := c.request(ctx).Delete(fmt.Sprintf(c.VersionUrl+"/person/%s", name))
	err = util.ResponseProcessor(response, err1)
	return
}

// Build Url from personQuery
func (c *PersonService) listPersonUrl(personQuery fun.PersonQuery) (url string) {
	url = c.VersionUrl + "/person?"

	// Add Pagination Params
	url += c.getPaginationParams(personQuery.Offset, personQuery.Limit)

	// Add Sort Params
	if personQuery.SortBy != "" {
		url += "&sort_by=" + personQuery.SortBy
		url += "&order=" + personQuery.Order
	}

	// Add Name and Gender if Provided
	if personQuery.Name != "" {
		url += "&name=" + personQuery.Name
	}
	if personQuery.Gender != "" {
		url += "&gender=" + personQuery.Gender
	}
	return
}
