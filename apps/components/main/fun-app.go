package main

import (
	"io/ioutil"

	fun_app "github.com/amanhigh/go-fun/apps/components/fun-app"
	config2 "github.com/amanhigh/go-fun/apps/models/config"
	"github.com/amanhigh/go-fun/apps/models/interfaces"
	"gopkg.in/yaml.v2"
)

func main() {
	var err error
	var config config2.FunAppConfig
	/* Read Config */
	var bytes []byte
	if bytes, err = ioutil.ReadFile("/etc/fun-app/config.yml"); err == nil {
		if err = yaml.Unmarshal(bytes, &config); err == nil {
			//go gometrics.Log(gometrics.DefaultRegistry, 5*time.Second, log.StandardLogger())

			/* Build Injector */
			injector := fun_app.NewFunAppInjector(config)
			var app interface{}

			/* Build App */
			if app, err = injector.BuildApp(); err == nil {
				err = app.(interfaces.ServerInterface).Start()
			}
		}

	}

	if err != nil {
		panic(err)
	}
}
