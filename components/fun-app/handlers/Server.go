package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/amanhigh/go-fun/common/telemetry"
	"github.com/amanhigh/go-fun/common/util"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"

	docs "github.com/amanhigh/go-fun/components/fun-app/docs"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

type FunServer struct {
	GinEngine *gin.Engine             `container:"type"`
	Server    *http.Server            `container:"type"`
	Shutdown  *util.GracefullShutdown `container:"type"`
	Tracer    trace.Tracer            `container:"type"`

	/* Handlers */
	PersonHandler *PersonHandler `container:"type"`
	AdminHandler  *AdminHandler  `container:"type"`
}

func (self *FunServer) initRoutes() {
	docs.SwaggerInfo.BasePath = "/v1"
	//Routes

	// Version Group
	v1 := self.GinEngine.Group("/v1")

	personGroup := v1.Group("/person")
	personGroup.GET("/", self.PersonHandler.ListPersons)
	personGroup.GET("/:id", self.PersonHandler.GetPerson)
	personGroup.PUT("/:id", self.PersonHandler.UpdatePerson)
	personGroup.POST("", self.PersonHandler.CreatePerson)
	personGroup.DELETE(":id", self.PersonHandler.DeletePersons)

	adminGroup := self.GinEngine.Group("/admin")
	adminGroup.GET("/stop", self.AdminHandler.Stop)

	//Add Swagger - https://github.com/swaggo/gin-swagger
	//make swag-fun

	self.GinEngine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// http://localhost:8080/debug/statsviz/
	self.GinEngine.GET("/debug/statsviz/*filepath", telemetry.StatvizMetrics)

	//Pprof (Use: http://localhost:8080/debug/pprof/)
	// make profile
	//Load Test:  wrk2 http://localhost:8080/v1/person/all/ -t 2 -c 100 -d 1m -R2000
	//Vegeta: echo "GET http://localhost:9000/v1/person/all" | vegeta attack -max-workers=2 -max-connections=100 -duration=1m -rate=2000/1s | tee results.bin | vegeta report
	//Vegeta Plot: vegeta plot results.bin > ~/Downloads/plot.html
	pprof.Register(self.GinEngine)
}

func (self *FunServer) Start(c context.Context) (err error) {
	self.initRoutes()

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	var errChan = make(chan error, 1)

	go func(errChan chan error) {
		if srvErr := self.Server.ListenAndServe(); srvErr != nil && srvErr != http.ErrServerClosed {
			errChan <- srvErr
		}
	}(errChan)

	//Read Error From GoRoutine or proceed in one second
	select {
	case err = <-errChan:
		log.Trace("Error while Starting Server", errChan)
	case <-time.After(time.Second):
		//No Error Occurred, wait for Graceful Shutdown Signal.
		ctx := self.Shutdown.Wait()

		//Trigger Shutdown Routine
		self.Stop(ctx)
	}

	return
}

func (self *FunServer) Stop(c context.Context) {
	// The context is used to inform the server it has few seconds to finish
	// the request it is currently handling
	ctx, span := self.Tracer.Start(c, "Stop.Server")
	defer span.End()

	ctxTimed, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := self.Server.Shutdown(ctxTimed); err != nil {
		log.WithContext(ctx).Fatal("Forced Shutdown, Graceful Exit Failed: ", err)
	}

	//Stop Tracer
	span.AddEvent("Stopping Tracer")
	telemetry.ShutdownTracerProvider(ctx)

	log.WithContext(ctx).Info("Graceful Exit Successful")
}
