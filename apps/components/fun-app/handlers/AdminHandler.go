package handlers

import (
	"github.com/amanhigh/go-fun/apps/common/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AdminHandler struct {
	Shutdown *util.GracefullShutdown `inject:""`
}

func (self *AdminHandler) Stop(c *gin.Context) {
	self.Shutdown.Stop()
	//TODO:Add check to prevent accidental stop
	c.JSON(http.StatusOK, "Stop Started")
}
