package handlers

import (
	"context"
	"net/http"
	"time"

	util2 "github.com/amanhigh/go-fun/common/util"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	docs "github.com/amanhigh/go-fun/components/fun-app/docs"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

type FunServer struct {
	GinEngine *gin.Engine              `inject:""`
	Server    *http.Server             `inject:""`
	Shutdown  *util2.GracefullShutdown `inject:""`

	/* Handlers */
	PersonHandler *PersonHandler `inject:""`
	AdminHandler  *AdminHandler  `inject:""`
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
	//Init/Update: swag i --parseDependency true  (in main.go dir)
	//URL: http://localhost:8080/swagger/index.html
	self.GinEngine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//Pprof (Use: http://localhost:8080/debug/pprof/)
	//go tool pprof -http=:8000 --seconds=30 http://localhost:8080/debug/pprof/profile
	//go tool pprof -http=:8001 http://localhost:8080/debug/pprof/heap
	//Load Test:  wrk2 http://localhost:8080/v1/person/all/ -t 2 -c 100 -d 1m -R2000
	//Vegeta: echo "GET http://localhost:9000/v1/person/all" | vegeta attack -max-workers=2 -max-connections=100 -duration=1m -rate=2000/1s | tee results.bin | vegeta report
	//Vegeta Plot: vegeta plot results.bin > ~/Downloads/plot.html
	pprof.Register(self.GinEngine)
}

func (self *FunServer) Start() (err error) {
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
		self.Shutdown.Wait()

		// The context is used to inform the server it has few seconds to finish
		// the request it is currently handling
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := self.Server.Shutdown(ctx); err != nil {
			log.Fatal("Server forced to shutdown: ", err)
		}

		log.Info("Server exiting")
	}

	return
}
