package metrics

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rcrowley/go-metrics"
)

/**
Sets path matched for this request
*/
func MatchedPath(c *gin.Context) {
	url := c.Request.URL.String()
	for _, p := range c.Params {
		url = strings.Replace(url, p.Value, ":"+p.Key, 1)
	}
	c.Set("matched_path", url)
}

/**
Sets Access Metrics for this request
*/
func AccessMetrics(context *gin.Context) {
	//TODO:Replace prometheus
	path, _ := context.Get("matched_path")
	matchedPath := path.(string)

	/* Time Taken */
	timer := metrics.GetOrRegister(matchedPath, metrics.NewTimer()).(metrics.Timer)
	timer.Time(context.Next)
	/* Status Counter */
	status := context.Writer.Status()

	statusCounter := metrics.GetOrRegister(fmt.Sprintf("%v.%v", matchedPath, status), metrics.NewCounter()).(metrics.Counter)
	statusCounter.Inc(1)

	/* Error Counter */
	if status != http.StatusOK {
		errorCounter := metrics.GetOrRegister("error."+matchedPath, metrics.NewCounter()).(metrics.Counter)
		errorCounter.Inc(1)
	}
}
