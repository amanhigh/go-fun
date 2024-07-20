package handlers

import (
	"net/http"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	Shutdown util.Shutdown
}

func NewAdminHandler(shutdown util.Shutdown) *AdminHandler {
	return &AdminHandler{shutdown}
}

/*
*
go test fun-app_test.go fun-app.go -coverprofile=coverage.out
curl http://localhost:8080/admin/stop
go tool cover -func=coverage.out
*/
func (self *AdminHandler) Stop(c *gin.Context) {
	self.Shutdown.Stop(c.Request.Context())
	c.JSON(http.StatusOK, "Stop Started")
}
