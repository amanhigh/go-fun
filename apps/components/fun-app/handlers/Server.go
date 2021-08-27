package handlers

import (
	"context"
	"github.com/amanhigh/go-fun/apps/common/util"
	"net/http"
	"time"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type FunServer struct {
	GinEngine *gin.Engine             `inject:""`
	Server    *http.Server            `inject:""`
	Shutdown  *util.GracefullShutdown `inject:""`

	/* Handlers */
	PersonHandler *PersonHandler `inject:""`
	AdminHandler  *AdminHandler  `inject:""`
}

func (self *FunServer) initRoutes() {
	//Routes
	personGroup := self.GinEngine.Group("/person")
	personGroup.GET("/all", self.PersonHandler.GetAllPerson)
	personGroup.POST("", self.PersonHandler.CreatePerson)
	personGroup.DELETE(":id", self.PersonHandler.DeletePersons)

	adminGroup := self.GinEngine.Group("/admin")
	adminGroup.GET("/stop", self.AdminHandler.Stop)

	//Pprof (Use: http://localhost:8080/debug/pprof/)
	pprof.Register(self.GinEngine)
}

func (self *FunServer) Start() (err error) {
	self.initRoutes()

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if srvErr := self.Server.ListenAndServe(); srvErr != nil && srvErr != http.ErrServerClosed {
			err = srvErr
		}
	}()

	self.Shutdown.Wait()

	// The context is used to inform the server it has few seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := self.Server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Info("Server exiting")
	return
}
