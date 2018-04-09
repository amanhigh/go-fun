package tools

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/amanhigh/go-fun/util"
	"github.com/bndr/gojenkins"
	"regexp"
	"time"
)

var compile = regexp.MustCompile("FINAL_DEB=.*_(.*)_all.deb")

type JenkinsClientInterface interface {
	Build(job string, params map[string]string) (buildNumber int64, err error)
	Status(jobName string, buildNumber int64) (status string, version string, err error)
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

func (self *JenkinsClient) Status(jobName string, buildNumber int64) (status string, version string, err error) {
	var build *gojenkins.Build
	if build, err = self.jenkins.GetBuild(jobName, buildNumber); err == nil {
		for ; build.IsRunning(); {
			util.PrintWhite(fmt.Sprintf("Job Running: %v", build.GetDuration()))
			time.Sleep(5 * time.Second)
		}
		if match := compile.FindStringSubmatch(build.GetConsoleOutput()); len(match) > 1 {
			version = match[1]
		}
		status = build.GetResult()
	}
	return
}
