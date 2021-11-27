package main

import (
	"github.com/amanhigh/go-fun/apps/components/fun-app/common"
	config2 "github.com/amanhigh/go-fun/apps/models/config"
	"github.com/amanhigh/go-fun/apps/models/interfaces"
	"github.com/caarlos0/env/v6"
)

func main() {
	var err error
	var config config2.FunAppConfig
	/* Read Config */
	if err = env.Parse(&config); err == nil {
		//go gometrics.Log(gometrics.DefaultRegistry, 5*time.Second, log.StandardLogger())

		/* Build Injector */
		injector := common.NewFunAppInjector(config)
		var app interface{}

		/* Build App */
		if app, err = injector.BuildApp(); err == nil {
			err = app.(interfaces.ServerInterface).Start()
		}
	}

	if err != nil {
		panic(err)
	}
}
