package command

import (
	"errors"
	"fmt"
	"strings"

	"github.com/amanhigh/go-fun/common/tools"
	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/models/config"
	"github.com/rs/zerolog/log"

	"github.com/spf13/cobra"
)

var clusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "Cluster Based Commands",
	Args:  cobra.ExactArgs(1),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		cluster = args[0]
		setLogLevel()
	},
}

var clusterPsshCmd = &cobra.Command{
	Use:   "pssh [Cluster] [Cmd]",
	Short: "Runs Parallel Ssh on Cluster",
	Args:  cobra.ExactArgs(2),
	PreRun: func(cmd *cobra.Command, args []string) {
		command = args[1]
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var selectedPssh tools.Pssh
		if selectedPssh, err = getPsshFromType(tyype); err == nil {
			selectedPssh.RunRange(command, cluster, parallelism, false, index, endIndex)
		}
		return
	},
}

var clusterCssCmd = &cobra.Command{
	Use:   "css [Cluster]",
	Short: "Runs Cluster Ssh",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		tools.ClusterSsh(args[0])
		return
	},
}

var clusterIndexCmd = &cobra.Command{
	Use:   "index [Cluster] [index]",
	Short: "Get Ip for Cluster & Index",
	Args:  cobra.ExactArgs(2),
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		index, err = util.ParseInt(args[1])
		return
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(tools.GetClusterHost(cluster, index))
	},
}

var clusterRemoveCmd = &cobra.Command{
	Use:   "remove [Main Cluster] [Remove Cluster]",
	Short: "Removes Ips in Remove Cluster from Main Cluster",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		count := tools.RemoveCluster(args[0], args[1])
		log.Info().Int("Removed", count).Str("Total", args[0]).Msg("Remove Cluster Ips")
	},
}

var clusterMd5Cmd = &cobra.Command{
	Use:   "md5 [cmd] [cluster(s) Space Separated]",
	Short: "Md5 Verification and Diff",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		for _, cluster := range strings.Fields(args[1]) {
			tools.Md5Checker(args[0], cluster)
		}
	},
}

var clusterSearchCmd = &cobra.Command{
	Use:   "search [keyword] [Opt:Cluster Index] [Opt:Ip Index]",
	Short: "Searches for matching Clusters",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		clusters := tools.SearchCluster(args[0])
		fmt.Println(strings.Join(clusters, "\n"))

		c, i := 1, 1
		switch len(args) {
		case 1:
			return
		case 2:
			c, err = util.ParseInt(args[1])

		case 3:
			if c, err = util.ParseInt(args[1]); err == nil {
				i, err = util.ParseInt(args[2])
			}
		}

		if err == nil {
			/* If Index is Zero do Cluster ssh */
			clusterName := clusters[c-1]
			if i == 0 {
				tools.ClusterSsh(clusterName)
			} else {
				ip := tools.GetClusterHost(clusterName, i)
				tools.LiveCommand("ssh " + ip)
			}
		}
		return
	},
}

func init() {
	clusterPsshCmd.Flags().StringVarP(&tyype, "type", "t", "f", "First alphabet of fast/display/slow")
	clusterPsshCmd.Flags().IntVarP(&parallelism, "parallel", "p", config.DEFAULT_PARALELISM, "Parallelism")
	clusterPsshCmd.Flags().IntVarP(&index, "start", "s", -1, "Starting Index")
	clusterPsshCmd.Flags().IntVarP(&endIndex, "end", "e", -1, "Ending Index")

	RootCmd.AddCommand(clusterCmd)
	clusterCmd.AddCommand(clusterPsshCmd, clusterIndexCmd,
		clusterRemoveCmd, clusterMd5Cmd, clusterSearchCmd, clusterCssCmd)
}

func getPsshFromType(psshType string) (selectedPssh tools.Pssh, err error) {
	switch psshType {
	case "f":
		selectedPssh = tools.FastPssh
	case "s":
		selectedPssh = tools.SlowPssh
	case "d":
		selectedPssh = tools.DisplayPssh
	default:
		err = errors.New("Invalid Pssh Type: " + psshType)
	}
	return
}
