package tools

import (
	"fmt"
	"github.com/amanhigh/go-fun/util"
)

func ClusterSsh(clusterName string) {
	util.PrintGreen("Clustering onto " + clusterName)
	LiveCommand(fmt.Sprintf("i2cssh -b -f /tmp/clusters/%v.txt", clusterName))
}
