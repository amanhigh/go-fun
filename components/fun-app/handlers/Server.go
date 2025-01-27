package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/amanhigh/go-fun/common/telemetry"
	"github.com/amanhigh/go-fun/common/util"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"

	docs "github.com/amanhigh/go-fun/components/fun-app/docs"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

type FunServer struct {
	GinEngine *gin.Engine   `container:"type"`
	Server    *http.Server  `container:"type"`
	Shutdown  util.Shutdown `container:"type"`
	Tracer    trace.Tracer  `container:"type"`

	/* Handlers */
	PersonHandler *PersonHandler `container:"type"`
	AdminHandler  *AdminHandler  `container:"type"`
}

func (fs *FunServer) initRoutes() {
	docs.SwaggerInfo.BasePath = "/v1"
	// Routes

	// Version Group
	v1 := fs.GinEngine.Group("/v1")

	personGroup := v1.Group("/person")
	personGroup.GET("/", fs.PersonHandler.ListPersons)
	personGroup.GET("/:id/audit", fs.PersonHandler.ListPersonAudit)
	personGroup.GET("/:id", fs.PersonHandler.GetPerson)
	personGroup.PUT("/:id", fs.PersonHandler.UpdatePerson)
	personGroup.POST("", fs.PersonHandler.CreatePerson)
	personGroup.DELETE(":id", fs.PersonHandler.DeletePersons)

	adminGroup := fs.GinEngine.Group("/admin")
	adminGroup.GET("/stop", fs.AdminHandler.Stop)

	// Add Swagger - https://github.com/swaggo/gin-swagger
	// make swag-fun

	fs.GinEngine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// http://localhost:8080/debug/statsviz/
	fs.GinEngine.GET("/debug/statsviz/*filepath", telemetry.StatvizMetrics)

	// Pprof (Use: http://localhost:8080/debug/pprof/)
	// make profile
	// Load Test:  wrk2 http://localhost:8080/v1/person/all/ -t 2 -c 100 -d 1m -R2000
	// Vegeta: echo "GET http://localhost:9000/v1/person/all" | vegeta attack -max-workers=2 -max-connections=100 -duration=1m -rate=2000/1s | tee results.bin | vegeta report
	// Vegeta Plot: vegeta plot results.bin > ~/Downloads/plot.html
	pprof.Register(fs.GinEngine)
}

func (fs *FunServer) Start(c context.Context) (err error) {
	fs.initRoutes()

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	var errChan = make(chan error, 1)

	go func(errChan chan error) {
		if srvErr := fs.Server.ListenAndServe(); srvErr != nil && srvErr != http.ErrServerClosed {
			errChan <- srvErr
		}
	}(errChan)

	// Read Error From GoRoutine or proceed in one second
	select {
	case err = <-errChan:
		zerolog.Ctx(c).Trace().Err(err).Msg("Failed To Start Server")
	case <-time.After(time.Second):
		// No Error Occurred, wait for Graceful Shutdown Signal.
		ctx := fs.Shutdown.Wait()

		// Trigger Shutdown Routine
		fs.Stop(ctx)
	}

	return
}

func (fs *FunServer) Stop(c context.Context) {
	// The context is used to inform the server it has few seconds to finish
	// the request it is currently handling
	ctx, span := fs.Tracer.Start(c, "Stop.Server")
	defer span.End()

	ctxTimed, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := fs.Server.Shutdown(ctxTimed); err != nil {
		zerolog.Ctx(c).Fatal().Ctx(ctx).Err(err).Msg("Forced Shutdown, Graceful Exit Failed: ")
	}

	// Stop Tracer
	span.AddEvent("Stopping Tracer")
	telemetry.ShutdownTracerProvider(ctx)

	zerolog.Ctx(c).Info().Ctx(ctx).Msg("Bye...")
}
