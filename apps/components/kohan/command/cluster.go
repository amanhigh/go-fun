package command

import (
	"errors"
	"fmt"
	"strings"

	"github.com/amanhigh/go-fun/apps/common/tools"
	"github.com/amanhigh/go-fun/apps/common/util"
	"github.com/amanhigh/go-fun/apps/components/kohan/core"
	util2 "github.com/amanhigh/go-fun/util"
	"github.com/spf13/cobra"
)

/* Vip add frontend port to vip */
var clusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "Cluster Based Commands",
	Args:  cobra.ExactArgs(1),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		cluster = args[0]
	},
}

var clusterSanityCmd = &cobra.Command{
	Use:   "sanity [Cluster] [Package] [Sanity Cmd]",
	Short: "Checks Sanity of Cluster",
	Args:  cobra.ExactArgs(3),
	PreRun: func(cmd *cobra.Command, args []string) {
		pkgName = args[1]
		command = args[2]
	},
	Run: func(cmd *cobra.Command, args []string) {
		core.ClusterSanity(pkgName, command, cluster)
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

var clusterIndexCmd = &cobra.Command{
	Use:   "index [Cluster] [index]",
	Short: "Get Ip for Cluster & Index",
	Args:  cobra.ExactArgs(2),
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		index, err = util.ParseInt(args[1])
		return
	},
	Run: func(cmd *cobra.Command, args []string) {
		tools.IndexedIp(cluster, index)
	},
}

var clusterRemoveCmd = &cobra.Command{
	Use:   "remove [Main Cluster] [Remove Cluster]",
	Short: "Removes Ips in Remove Cluster from Main Cluster",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		count := tools.RemoveCluster(args[0], args[1])
		util2.PrintGreen(fmt.Sprintf("%v items removed from %v", count, args[0]))
	},
}

var clusterMd5Cmd = &cobra.Command{
	Use:   "md5 [cmd] [cluster(s) Space Separated]",
	Short: "Md5 Verification and Diff",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		for _, cluster := range strings.Fields(args[1]) {
			util.Md5Checker(args[0], cluster)
		}
	},
}

func init() {
	clusterPsshCmd.Flags().StringVarP(&tyype, "type", "t", "f", "First alphabet of fast/display/slow")
	clusterPsshCmd.Flags().IntVarP(&parallelism, "parallel", "p", util2.DEFAULT_PARALELISM, "Parallelism")
	clusterPsshCmd.Flags().IntVarP(&index, "start", "s", -1, "Starting Index")
	clusterPsshCmd.Flags().IntVarP(&endIndex, "end", "e", -1, "Ending Index")

	RootCmd.AddCommand(clusterCmd)
	clusterCmd.AddCommand(clusterSanityCmd, clusterPsshCmd, clusterIndexCmd, clusterRemoveCmd, clusterMd5Cmd)
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
