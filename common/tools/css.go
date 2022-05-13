package tools

import (
	"fmt"
	"github.com/fatih/color"
	"strings"
)

func ClusterSsh(clusterName string) {
	color.Green("Clustering onto " + clusterName)
	ips := ReadClusterFile(clusterName)
	LiveCommand(fmt.Sprintf("cssh %s", strings.Join(ips, " ")))
}
