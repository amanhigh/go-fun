package tools

import (
	"fmt"
	"regexp"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/amanhigh/go-fun/util"
	"github.com/bndr/gojenkins"
)

var compile = regexp.MustCompile("version : (.*)")

type JenkinsClientInterface interface {
	Build(job string, params map[string]string) (err error)
	Status(jobName string) (status string, version string, err error)
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

func (self *JenkinsClient) Build(job string, params map[string]string) (err error) {
	_, err = self.jenkins.BuildJob(job, params)
	return
}

func (self *JenkinsClient) Status(jobName string) (status string, version string, err error) {
	var job *gojenkins.Job
	var build *gojenkins.Build
	if job, err = self.jenkins.GetJob(jobName); err == nil {
		if build, err = job.GetLastBuild(); err == nil {
			if build, err = self.jenkins.GetBuild(jobName, build.GetBuildNumber()); err == nil {
				for i := 1; build.IsRunning(); i++ {
					util.PrintWhite(fmt.Sprintf("Job Running - %v", i))
					time.Sleep(10 * time.Second)
				}
				if match := compile.FindStringSubmatch(build.GetConsoleOutput()); len(match) > 1 {
					version = match[1]
				}
				status = build.GetResult()
			}
		}
	}
	return
}
