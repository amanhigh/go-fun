package tools

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/amanhigh/go-fun/util"
	"github.com/bndr/gojenkins"
)

type JenkinsClientInterface interface {
	Build(job string, params map[string]string) (buildNumber int64, err error)
}

type JenkinsClient struct {
	jenkins *gojenkins.Jenkins
}

func NewJenkinsClient(ip string, userName string, apiKey string) JenkinsClientInterface {
	jenkins := gojenkins.CreateJenkins(util.KeepAliveInsecureClient.(*util.HttpClient).Client, fmt.Sprintf("https://%v/", ip), userName, apiKey)
	if _, err := jenkins.Init(); err != nil {
		log.WithFields(log.Fields{"Error": err}).Fatal("Error Building Jenkins Client")
	}
	return &JenkinsClient{jenkins: jenkins}
}

func (self *JenkinsClient) Build(job string, params map[string]string) (buildNumber int64, err error) {
	return self.jenkins.BuildJob(job, params)
}