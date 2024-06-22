package tools

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
)

func ClusterSsh(clusterName string) {
	log.Info().Str("Cluster", clusterName).Msg("Clustering")
	ips := ReadClusterFile(clusterName)
	LiveCommand(fmt.Sprintf("cssh %s", strings.Join(ips, " ")))
}
