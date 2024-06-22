package telemetry

import (
	"net/http"
	"strings"

	"github.com/arl/statsviz"
	"github.com/gin-gonic/gin"
)

var ws http.HandlerFunc
var index http.HandlerFunc

func init() {
	// Create statsviz server.
	srv, _ := statsviz.NewServer()

	ws = srv.Ws()
	index = srv.Index()
}

/*
*
Access Metrics with Path Params resolved
*/
func AccessMetrics(c *gin.Context) (matchedPath string) {
	matchedPath = c.Request.URL.String()
	for _, p := range c.Params {
		matchedPath = strings.Replace(matchedPath, p.Value, ":"+p.Key, 1)
	}
	return
}

/*
*
Sets Up Heap, GC and Goroutine Metric Graphs
*/
func StatvizMetrics(context *gin.Context) {
	if context.Param("filepath") == "/ws" {
		ws(context.Writer, context.Request)
		return
	}
	index(context.Writer, context.Request)
}
