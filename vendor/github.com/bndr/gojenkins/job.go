// Copyright 2015 Vadim Kravcenko
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package gojenkins

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/url"
	"strconv"
)

type Job struct {
	Raw     *jobResponse
	Jenkins *Jenkins
	Base    string
}

type jobBuild struct {
	Number int64
	URL    string
}

type job struct {
	Name  string `json:"name"`
	Url   string `json:"url"`
	Color string `json:"color"`
}

type parameterDefinition struct {
	DefaultParameterValue struct {
		Name  string `json:"name"`
		Value bool   `json:"value"`
	} `json:"defaultParameterValue"`
	Description string `json:"description"`
	Name        string `json:"name"`
	Type        string `json:"type"`
}

type jobResponse struct {
	Actions            []generalObj
	Buildable          bool `json:"buildable"`
	Builds             []jobBuild
	Color              string      `json:"color"`
	ConcurrentBuild    bool        `json:"concurrentBuild"`
	Description        string      `json:"description"`
	DisplayName        string      `json:"displayName"`
	DisplayNameOrNull  interface{} `json:"displayNameOrNull"`
	DownstreamProjects []job       `json:"downstreamProjects"`
	FirstBuild         jobBuild
	HealthReport       []struct {
		Description   string `json:"description"`
		IconClassName string `json:"iconClassName"`
		IconUrl       string `json:"iconUrl"`
		Score         int64  `json:"score"`
	} `json:"healthReport"`
	InQueue               bool     `json:"inQueue"`
	KeepDependencies      bool     `json:"keepDependencies"`
	LastBuild             jobBuild `json:"lastBuild"`
	LastCompletedBuild    jobBuild `json:"lastCompletedBuild"`
	LastFailedBuild       jobBuild `json:"lastFailedBuild"`
	LastStableBuild       jobBuild `json:"lastStableBuild"`
	LastSuccessfulBuild   jobBuild `json:"lastSuccessfulBuild"`
	LastUnstableBuild     jobBuild `json:"lastUnstableBuild"`
	LastUnsuccessfulBuild jobBuild `json:"lastUnsuccessfulBuild"`
	Name                  string   `json:"name"`
	NextBuildNumber       int64    `json:"nextBuildNumber"`
	Property              []struct {
		ParameterDefinitions []parameterDefinition `json:"parameterDefinitions"`
	} `json:"property"`
	QueueItem        interface{} `json:"queueItem"`
	Scm              struct{}    `json:"scm"`
	UpstreamProjects []job       `json:"upstreamProjects"`
	URL              string      `json:"url"`
}

func (j *Job) GetName() string {
	return j.Raw.Name
}

func (j *Job) GetDescription() string {
	return j.Raw.Description
}

func (j *Job) GetDetails() *jobResponse {
	return j.Raw
}

func (j *Job) GetBuild(id int64) (*Build, error) {
	build := Build{Jenkins: j.Jenkins, Job: j, Raw: new(buildResponse), Depth: 1, Base: "/job/" + j.GetName() + "/" + strconv.FormatInt(id, 10)}
	status, err := build.Poll()
	if err != nil {
		return nil, err
	}
	if status == 200 {
		return &build, nil
	}
	return nil, errors.New(strconv.Itoa(status))
}

func (j *Job) getBuildByType(buildType string) (*Build, error) {
	allowed := map[string]jobBuild{
		"lastStableBuild":     j.Raw.LastStableBuild,
		"lastSuccessfulBuild": j.Raw.LastSuccessfulBuild,
		"lastBuild":           j.Raw.LastBuild,
		"lastCompletedBuild":  j.Raw.LastCompletedBuild,
		"firstBuild":          j.Raw.FirstBuild,
		"lastFailedBuild":     j.Raw.LastFailedBuild,
	}
	number := ""
	if val, ok := allowed[buildType]; ok {
		number = strconv.FormatInt(val.Number, 10)
	} else {
		panic("No Such Build")
	}
	build := Build{
		Jenkins: j.Jenkins,
		Depth:   1,
		Job:     j,
		Raw:     new(buildResponse),
		Base:    "/job/" + j.GetName() + "/" + number}
	status, err := build.Poll()
	if err != nil {
		return nil, err
	}
	if status == 200 {
		return &build, nil
	}
	return nil, errors.New(strconv.Itoa(status))
}

func (j *Job) GetLastSuccessfulBuild() (*Build, error) {
	return j.getBuildByType("lastSuccessfulBuild")
}

func (j *Job) GetFirstBuild() (*Build, error) {
	return j.getBuildByType("firstBuild")
}

func (j *Job) GetLastBuild() (*Build, error) {
	return j.getBuildByType("lastBuild")
}

func (j *Job) GetLastStableBuild() (*Build, error) {
	return j.getBuildByType("lastStableBuild")
}

func (j *Job) GetLastFailedBuild() (*Build, error) {
	return j.getBuildByType("lastFailedBuild")
}

func (j *Job) GetLastCompletedBuild() (*Build, error) {
	return j.getBuildByType("lastCompletedBuild")
}

// Returns All Builds with Number and URL
func (j *Job) GetAllBuildIds() ([]jobBuild, error) {
	var buildsResp struct {
		Builds []jobBuild `json:"allBuilds"`
	}
	_, err := j.Jenkins.Requester.GetJSON(j.Base, &buildsResp, map[string]string{"tree": "allBuilds[number,url]"})
	if err != nil {
		return nil, err
	}
	return buildsResp.Builds, nil
}

func (j *Job) GetUpstreamJobsMetadata() []job {
	return j.Raw.UpstreamProjects
}

func (j *Job) GetDownstreamJobsMetadata() []job {
	return j.Raw.DownstreamProjects
}

func (j *Job) GetUpstreamJobs() ([]*Job, error) {
	jobs := make([]*Job, len(j.Raw.UpstreamProjects))
	for i, job := range j.Raw.UpstreamProjects {
		ji, err := j.Jenkins.GetJob(job.Name)
		if err != nil {
			return nil, err
		}
		jobs[i] = ji
	}
	return jobs, nil
}

func (j *Job) GetDownstreamJobs() ([]*Job, error) {
	jobs := make([]*Job, len(j.Raw.DownstreamProjects))
	for i, job := range j.Raw.DownstreamProjects {
		ji, err := j.Jenkins.GetJob(job.Name)
		if err != nil {
			return nil, err
		}
		jobs[i] = ji
	}
	return jobs, nil
}

func (j *Job) Enable() (bool, error) {
	resp, err := j.Jenkins.Requester.Post(j.Base+"/enable", nil, nil, nil)
	if err != nil {
		return false, err
	}
	if resp.StatusCode != 200 {
		return false, errors.New(strconv.Itoa(resp.StatusCode))
	}
	return true, nil
}

func (j *Job) Disable() (bool, error) {
	resp, err := j.Jenkins.Requester.Post(j.Base+"/disable", nil, nil, nil)
	if err != nil {
		return false, err
	}
	if resp.StatusCode != 200 {
		return false, errors.New(strconv.Itoa(resp.StatusCode))
	}
	return true, nil
}

func (j *Job) Delete() (bool, error) {
	resp, err := j.Jenkins.Requester.Post(j.Base+"/doDelete", nil, nil, nil)
	if err != nil {
		return false, err
	}
	if resp.StatusCode != 200 {
		return false, errors.New(strconv.Itoa(resp.StatusCode))
	}
	return true, nil
}

func (j *Job) Rename(name string) (bool, error) {
	data := url.Values{}
	data.Set("newName", name)
	_, err := j.Jenkins.Requester.Post(j.Base+"/doRename", bytes.NewBufferString(data.Encode()), nil, nil)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (j *Job) Create(config string, qr ...interface{}) (*Job, error) {
	var querystring map[string]string
	if len(qr) > 0 {
		querystring = qr[0].(map[string]string)
	}
	resp, err := j.Jenkins.Requester.PostXML("/createItem", config, j.Raw, querystring)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == 200 {
		j.Poll()
		return j, nil
	}
	return nil, errors.New(strconv.Itoa(resp.StatusCode))
}

func (j *Job) Copy(from string, newName string) (*Job, error) {
	qr := map[string]string{"name": newName, "from": from, "mode": "copy"}
	resp, err := j.Jenkins.Requester.Post("/createItem", nil, nil, qr)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == 200 {
		return j, nil
	}
	return nil, errors.New(strconv.Itoa(resp.StatusCode))
}

func (j *Job) GetConfig() (string, error) {
	var data string
	_, err := j.Jenkins.Requester.GetXML(j.Base+"/config.xml", &data, nil)
	if err != nil {
		return "", err
	}
	return data, nil
}

func (j *Job) GetParameters() ([]parameterDefinition, error) {
	_, err := j.Poll()
	if err != nil {
		return nil, err
	}
	var parameters []parameterDefinition
	for _, property := range j.Raw.Property {
		for _, param := range property.ParameterDefinitions {
			parameters = append(parameters, param)
		}
	}
	return parameters, nil
}

func (j *Job) IsQueued() (bool, error) {
	if _, err := j.Poll(); err != nil {
		return false, err
	}
	return j.Raw.InQueue, nil
}

func (j *Job) IsRunning() (bool, error) {
	if _, err := j.Poll(); err != nil {
		return false, err
	}
	lastBuild, err := j.GetLastBuild()
	if err != nil {
		return false, err
	}
	return lastBuild.IsRunning(), nil
}

func (j *Job) IsEnabled() (bool, error) {
	if _, err := j.Poll(); err != nil {
		return false, err
	}
	return j.Raw.Color != "disabled", nil
}

func (j *Job) HasQueuedBuild() {

}

func (j *Job) InvokeSimple(params map[string]string) (bool, error) {
	isQueued, err := j.IsQueued()
	if err != nil {
		return false, err
	}
	if isQueued {
		Error.Printf("%s is already running", j.GetName())
		return false, nil
	}

	endpoint := "/build"
	parameters, err := j.GetParameters()
	if err != nil {
		return false, err
	}
	if len(parameters) > 0 {
		endpoint = "/buildWithParameters"
	}
	data := url.Values{}
	for k, v := range params {
		data.Set(k, v)
	}
	resp, err := j.Jenkins.Requester.Post(j.Base+endpoint, bytes.NewBufferString(data.Encode()), nil, nil)
	if err != nil {
		return false, err
	}
	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		Error.Println("Could not invoke job %s", j.GetName())
		return false, errors.New("Could not invoke job " + j.GetName())
	}
	return true, nil
}

func (j *Job) Invoke(files []string, skipIfRunning bool, params map[string]string, cause string, securityToken string) (bool, error) {
	isQueued, err := j.IsQueued()
	if err != nil {
		return false, err
	}
	if isQueued {
		Error.Printf("%s is already running", j.GetName())
		return false, nil
	}
	isRunning, err := j.IsRunning()
	if err != nil {
		return false, err
	}
	if isRunning && skipIfRunning {
		Warning.Printf("%s Will not request new build because %s is already running", j.GetName())
	}

	base := "/build"

	// If parameters are specified - url is /builWithParameters
	if params != nil {
		base = "/buildWithParameters"
	} else {
		params = make(map[string]string)
	}

	// If files are specified - url is /build
	if files != nil {
		base = "/build"
	}
	reqParams := map[string]string{}
	buildParams := map[string]string{}
	if securityToken != "" {
		reqParams["token"] = securityToken
	}

	buildParams["json"] = string(makeJson(params))
	b, _ := json.Marshal(buildParams)
	resp, err := j.Jenkins.Requester.PostFiles(j.Base+base, bytes.NewBuffer(b), nil, reqParams, files)
	if err != nil {
		return false, err
	}
	if resp.StatusCode == 200 || resp.StatusCode == 201 {
		return true, nil
	}
	return false, errors.New(strconv.Itoa(resp.StatusCode))
}

func (j *Job) Poll() (int, error) {
	_, err := j.Jenkins.Requester.GetJSON(j.Base, j.Raw, nil)
	if err != nil {
		return 0, err
	}
	return j.Jenkins.Requester.LastResponse.StatusCode, nil
}
