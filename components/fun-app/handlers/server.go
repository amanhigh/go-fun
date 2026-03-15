package handlers

import (
	"context"

	"github.com/amanhigh/go-fun/common/telemetry"
	"github.com/amanhigh/go-fun/common/util"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"

	docs "github.com/amanhigh/go-fun/components/fun-app/docs"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

// FunAppServerLifecycle implements ServerLifecycle for the FunApp HTTP server.
type FunAppServerLifecycle struct {
	Tracer trace.Tracer `container:"type"`

	/* Handlers */
	PersonHandler     PersonHandler     `container:"type"`
	EnrollmentHandler EnrollmentHandler `container:"type"`
	AdminHandler      AdminHandler      `container:"type"`

	Watermill util.WatermillController `container:"type"`
}

var _ util.ServerLifecycle = (*FunAppServerLifecycle)(nil)

func (fs *FunAppServerLifecycle) RegisterRoutes(engine *gin.Engine) {
	docs.SwaggerInfo.BasePath = "/v1"

	// Version Group
	v1 := engine.Group("/v1")

	personGroup := v1.Group("/person")
	personGroup.GET("/", fs.PersonHandler.ListPersons)
	personGroup.GET("/:id/audit", fs.PersonHandler.ListPersonAudit)
	personGroup.GET("/:id", fs.PersonHandler.GetPerson)
	personGroup.PUT("/:id", fs.PersonHandler.UpdatePerson)
	personGroup.POST("", fs.PersonHandler.CreatePerson)
	personGroup.DELETE(":id", fs.PersonHandler.DeletePersons)

	enrollmentGroup := v1.Group("/enrollments")
	enrollmentGroup.POST("", fs.EnrollmentHandler.CreateEnrollment)
	enrollmentGroup.GET(":personId", fs.EnrollmentHandler.GetEnrollment)

	adminGroup := engine.Group("/admin")
	adminGroup.GET("/stop", fs.AdminHandler.Stop)

	// Pprof (Use: http://localhost:8080/debug/pprof/)
	// make profile
	// Load Test:  wrk2 http://localhost:8080/v1/person/all/ -t 2 -c 100 -d 1m -R2000
	// Vegeta: echo "GET http://localhost:9000/v1/person/all" | vegeta attack -max-workers=2 -max-connections=100 -duration=1m -rate=2000/1s | tee results.bin | vegeta report
	// Vegeta Plot: vegeta plot results.bin > ~/Downloads/plot.html
	pprof.Register(engine)
}

func (fs *FunAppServerLifecycle) RegisterSwagger(engine *gin.Engine) {
	// Add Swagger - https://github.com/swaggo/gin-swagger
	// make swag-fun
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func (fs *FunAppServerLifecycle) BeforeStart(ctx context.Context) {
	fs.Watermill.Start(ctx)
}

func (fs *FunAppServerLifecycle) BeforeShutdown(ctx context.Context) {
	_, span := fs.Tracer.Start(ctx, "Stop.Server")
	defer span.End()

	fs.Watermill.Shutdown(ctx)
}

func (fs *FunAppServerLifecycle) AfterShutdown(ctx context.Context) {
	telemetry.ShutdownTracerProvider(ctx)
}
