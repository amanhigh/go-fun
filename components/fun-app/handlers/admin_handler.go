package handlers

import (
	"net/http"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/gin-gonic/gin"
)

type AdminHandler interface {
	Stop(c *gin.Context)
}

type AdminHandlerImpl struct {
	Shutdown util.Shutdown
}

func NewAdminHandler(shutdown util.Shutdown) AdminHandler {
	return &AdminHandlerImpl{Shutdown: shutdown}
}

/*
*

	go test fun-app_test.go fun-app.go -coverprofile=coverage.out
	curl http://localhost:8080/admin/stop
	go tool cover -func=coverage.out
*/
func (ah *AdminHandlerImpl) Stop(c *gin.Context) {
	ah.Shutdown.Stop(c.Request.Context())
	c.JSON(http.StatusOK, "Stop Started")
}
