// SDK for FunApp with API's for Student Handler using Resty
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
	StudentService    StudentServiceInterface
	AdminService      AdminServiceInterface
	EnrollmentService EnrollmentServiceInterface
}

type StudentServiceInterface interface {
	GetStudent(ctx context.Context, name string) (student fun.Student, err common.HttpError)
	CreateStudent(ctx context.Context, request fun.StudentRequest) (student fun.Student, err common.HttpError)
	UpdateStudent(ctx context.Context, id string, student fun.StudentRequest) (err common.HttpError)
	ListStudent(ctx context.Context, query fun.StudentQuery) (studentList fun.StudentList, err common.HttpError)
	ListStudentAudit(ctx context.Context, id string) (studentAuditList []fun.StudentAudit, err common.HttpError)
	DeleteStudent(ctx context.Context, name string) (err common.HttpError)
}

type EnrollmentServiceInterface interface {
	CreateEnrollment(ctx context.Context, request fun.EnrollmentRequest) (enrollment fun.Enrollment, err common.HttpError)
	GetEnrollment(ctx context.Context, studentID string) (enrollment fun.Enrollment, err common.HttpError)
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

type StudentService struct {
	BaseService
}

type EnrollmentService struct {
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
		StudentService:    &StudentService{BaseService: baseService},
		AdminService:      &AdminService{BaseService: baseService},
		EnrollmentService: &EnrollmentService{BaseService: baseService},
	}
}

func (c *StudentService) CreateStudent(ctx context.Context, request fun.StudentRequest) (student fun.Student, err common.HttpError) {
	response, err1 := c.request(ctx).SetHeader("Content-Type", "application/json").
		SetBody(request).SetResult(&student).Post(c.VersionUrl + "/student")
	err = util.ResponseProcessor(response, err1)
	return
}

func (c *StudentService) GetStudent(ctx context.Context, name string) (student fun.Student, err common.HttpError) {
	url := fmt.Sprintf(c.VersionUrl+"/student/%s", name)
	response, err1 := c.request(ctx).SetResult(&student).Get(url)
	err = util.ResponseProcessor(response, err1)
	return
}

func (c *StudentService) ListStudent(ctx context.Context, studentQuery fun.StudentQuery) (studentList fun.StudentList, err common.HttpError) {
	response, err1 := c.request(ctx).SetResult(&studentList).Get(c.listStudentUrl(studentQuery))
	err = util.ResponseProcessor(response, err1)
	return
}

func (c *StudentService) ListStudentAudit(ctx context.Context, id string) (studentAuditList []fun.StudentAudit, err common.HttpError) {
	response, err1 := c.request(ctx).SetResult(&studentAuditList).Get(fmt.Sprintf(c.VersionUrl+"/student/%s/audit", id))
	err = util.ResponseProcessor(response, err1)
	return
}

func (c *StudentService) UpdateStudent(ctx context.Context, id string, student fun.StudentRequest) (err common.HttpError) {
	response, err1 := c.request(ctx).SetBody(student).Put(fmt.Sprintf(c.VersionUrl+"/student/%s", id))
	err = util.ResponseProcessor(response, err1)
	return
}

func (c *StudentService) DeleteStudent(ctx context.Context, name string) (err common.HttpError) {
	response, err1 := c.request(ctx).Delete(fmt.Sprintf(c.VersionUrl+"/student/%s", name))
	err = util.ResponseProcessor(response, err1)
	return
}

func (e *EnrollmentService) CreateEnrollment(ctx context.Context, request fun.EnrollmentRequest) (enrollment fun.Enrollment, err common.HttpError) {
	res, err1 := e.request(ctx).SetHeader("Content-Type", "application/json").
		SetBody(request).SetResult(&enrollment).Post(e.VersionUrl + "/enrollments")
	err = util.ResponseProcessor(res, err1)
	return
}

func (e *EnrollmentService) GetEnrollment(ctx context.Context, studentID string) (enrollment fun.Enrollment, err common.HttpError) {
	res, err1 := e.request(ctx).SetResult(&enrollment).Get(fmt.Sprintf(e.VersionUrl+"/enrollments/%s", studentID))
	err = util.ResponseProcessor(res, err1)
	return
}

// Build Url from studentQuery
func (c *StudentService) listStudentUrl(studentQuery fun.StudentQuery) (url string) {
	url = c.VersionUrl + "/student?"

	// Add Pagination Params
	url += c.getPaginationParams(studentQuery.Offset, studentQuery.Limit)

	// Add Sort Params
	if studentQuery.SortBy != "" {
		url += "&sort_by=" + studentQuery.SortBy
		url += "&sort-order=" + string(studentQuery.SortOrder)
	}

	// Add Name and Gender if Provided
	if studentQuery.Name != "" {
		url += "&name=" + studentQuery.Name
	}
	if studentQuery.Gender != "" {
		url += "&gender=" + studentQuery.Gender
	}
	return
}
