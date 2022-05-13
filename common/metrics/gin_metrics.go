package metrics

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

/**
Access Metrics with Path Params resolved
*/
func AccessMetrics(c *gin.Context) (matchedPath string) {
	matchedPath = c.Request.URL.String()
	for _, p := range c.Params {
		matchedPath = strings.Replace(matchedPath, p.Value, ":"+p.Key, 1)
	}
	fmt.Println(matchedPath)
	return
}
