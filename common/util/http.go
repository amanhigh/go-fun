package util

import (
	"net/http"

	"github.com/amanhigh/go-fun/models/config"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/rs/zerolog/log"
)

func UnwrapRequest(c *gin.Context, requestParam, pathParam, queryParam any, callback func(c *gin.Context)) {
	var err error
	if requestParam != nil {
		err = c.ShouldBindBodyWith(requestParam, binding.JSON)
	}
	if pathParam != nil && err == nil {
		err = c.ShouldBindUri(pathParam)
	}
	if queryParam != nil && err == nil {
		err = c.ShouldBindQuery(queryParam)
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		log.Error().Err(err).Msg("Request Decoding Failied")
	}

	callback(c)
}

// Helper to create server with timeouts from DefaultHttpConfig
func NewTestServer(addr string) *http.Server {
	return &http.Server{
		Addr:              addr,
		ReadTimeout:       config.DefaultHttpConfig.ReadTimeout,
		WriteTimeout:      config.DefaultHttpConfig.WriteTimeout,
		IdleTimeout:       config.DefaultHttpConfig.IdleTimeout,
		ReadHeaderTimeout: config.DefaultHttpConfig.ReadHeaderTimeout,
	}
}
