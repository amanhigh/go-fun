package common

import (
	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/models/config"
	"github.com/caarlos0/env/v6"
)

func RunFunApp() {
	var err error
	var config config.FunAppConfig
	/* Read Config */
	if err = env.Parse(&config); err == nil {
		// go gometrics.Log(gometrics.DefaultRegistry, 5*time.Second, log.StandardLogger())

		/* Build Injector */
		injector := NewFunAppInjector(config)
		var app any

		/* Build App */
		if app, err = injector.BuildApp(); err == nil {
			err = app.(util.HttpServer).Start() //nolint:errcheck // error is handled via panic
		}
	}

	if err != nil {
		panic(err)
	}
}
