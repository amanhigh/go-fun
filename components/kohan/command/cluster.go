package command

import (
	"errors"
	"fmt"
	tools2 "github.com/amanhigh/go-fun/common/tools"
	util2 "github.com/amanhigh/go-fun/common/util"
	core2 "github.com/amanhigh/go-fun/components/kohan/core"
	config2 "github.com/amanhigh/go-fun/models/config"
	"github.com/fatih/color"
	"strings"

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
	Use:   "sanity [Cluster] [Package(Direct/Regex)] [Sanity Cmd]",
	Short: "Checks Sanity of Cluster",
	Args:  cobra.ExactArgs(3),
	PreRun: func(cmd *cobra.Command, args []string) {
		pkgName = args[1]
		command = args[2]
	},
	Run: func(cmd *cobra.Command, args []string) {
		core2.ClusterSanity(pkgName, command, cluster)
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
		var selectedPssh tools2.Pssh
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
		tools2.ClusterSsh(args[0])
		return
	},
}

var clusterIndexCmd = &cobra.Command{
	Use:   "index [Cluster] [index]",
	Short: "Get Ip for Cluster & Index",
	Args:  cobra.ExactArgs(2),
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		index, err = util2.ParseInt(args[1])
		return
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(tools2.GetClusterHost(cluster, index))
	},
}

var clusterRemoveCmd = &cobra.Command{
	Use:   "remove [Main Cluster] [Remove Cluster]",
	Short: "Removes Ips in Remove Cluster from Main Cluster",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		count := tools2.RemoveCluster(args[0], args[1])
		color.Green("%v items removed from %v", count, args[0])
	},
}

var clusterMd5Cmd = &cobra.Command{
	Use:   "md5 [cmd] [cluster(s) Space Separated]",
	Short: "Md5 Verification and Diff",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		for _, cluster := range strings.Fields(args[1]) {
			tools2.Md5Checker(args[0], cluster)
		}
	},
}

var clusterSearchCmd = &cobra.Command{
	Use:   "search [keyword] [Opt:Cluster Index] [Opt:Ip Index]",
	Short: "Searches for matching Clusters",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		clusters := tools2.SearchCluster(args[0])
		fmt.Println(strings.Join(clusters, "\n"))

		c, i := 1, 1
		switch len(args) {
		case 1:
			return
		case 2:
			c, err = util2.ParseInt(args[1])

		case 3:
			if c, err = util2.ParseInt(args[1]); err == nil {
				i, err = util2.ParseInt(args[2])
			}
		}

		if err == nil {
			/* If Index is Zero do Cluster ssh */
			clusterName := clusters[c-1]
			if i == 0 {
				tools2.ClusterSsh(clusterName)
			} else {
				ip := tools2.GetClusterHost(clusterName, i)
				tools2.LiveCommand("ssh " + ip)
			}
		}
		return
	},
}

func init() {
	clusterPsshCmd.Flags().StringVarP(&tyype, "type", "t", "f", "First alphabet of fast/display/slow")
	clusterPsshCmd.Flags().IntVarP(&parallelism, "parallel", "p", config2.DEFAULT_PARALELISM, "Parallelism")
	clusterPsshCmd.Flags().IntVarP(&index, "start", "s", -1, "Starting Index")
	clusterPsshCmd.Flags().IntVarP(&endIndex, "end", "e", -1, "Ending Index")

	RootCmd.AddCommand(clusterCmd)
	clusterCmd.AddCommand(clusterSanityCmd, clusterPsshCmd, clusterIndexCmd,
		clusterRemoveCmd, clusterMd5Cmd, clusterSearchCmd, clusterCssCmd)
}

func getPsshFromType(psshType string) (selectedPssh tools2.Pssh, err error) {
	switch psshType {
	case "f":
		selectedPssh = tools2.FastPssh
	case "s":
		selectedPssh = tools2.SlowPssh
	case "d":
		selectedPssh = tools2.DisplayPssh
	default:
		err = errors.New("Invalid Pssh Type: " + psshType)
	}
	return
}
