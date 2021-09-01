package tools

import (
	"fmt"
	"github.com/amanhigh/go-fun/util"
	"strings"
)

func ClusterSsh(clusterName string) {
	util.PrintGreen("Clustering onto " + clusterName)
	ips := ReadClusterFile(clusterName)
	LiveCommand(fmt.Sprintf("cssh %s", strings.Join(ips, " ")))
}
