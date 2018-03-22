package command

import (
	"fmt"
	"github.com/amanhigh/go-fun/kohan/commander/components"
	"github.com/amanhigh/go-fun/kohan/commander/tools"
	"github.com/amanhigh/go-fun/util"
	"github.com/spf13/cobra"
)

/* Vip add frontend port to vip */
var clusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "Cluster Based Commands",
	Long:  `Cluster Based Commands`,
	Args:  cobra.ExactArgs(1),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		cluster = args[0]
	},
}

var clusterSanityCmd = &cobra.Command{
	Use:   "sanity [Cluster] [Package] [Sanity Cmd]",
	Short: "Checks Sanity of Cluster",
	Long:  `Checks Sanity of Cluster`,
	Args:  cobra.ExactArgs(3),
	PreRun: func(cmd *cobra.Command, args []string) {
		pkgName = args[1]
		command = args[2]
	},
	Run: func(cmd *cobra.Command, args []string) {
		components.ClusterSanity(pkgName, command, cluster)
	},
}

var clusterPsshCmd = &cobra.Command{
	Use:   "pssh [Cluster] [Cmd]",
	Short: "Runs Parallel Ssh on Cluster",
	Long:  `Runs Parallel Ssh on Cluster`,
	Args:  cobra.ExactArgs(2),
	PreRun: func(cmd *cobra.Command, args []string) {
		command = args[2]
	},
	Run: func(cmd *cobra.Command, args []string) {
		selectedPssh := getPsshFromType(tyype)
		selectedPssh.RunRange(command,cluster,parallelism,false,startIndex,endIndex)
	},
}

func init() {
	rootCmd.AddCommand(clusterCmd)
	clusterPsshCmd.Flags().StringVarP(&tyype, "type", "t", "fast", "fast/display/slow")
	clusterPsshCmd.Flags().IntVarP(&parallelism, "parallel", "p", util.DEFAULT_PARALELISM, "Parallelism")
	clusterPsshCmd.Flags().IntVarP(&startIndex, "start", "s", -1, "Starting Index")
	clusterPsshCmd.Flags().IntVarP(&endIndex, "end", "e", -1, "Ending Index")
	clusterCmd.AddCommand(clusterSanityCmd)
}

func getPsshFromType(psshType string) tools.Pssh {
	var selectedPssh tools.Pssh
	switch psshType {
	case "fast":
		selectedPssh = tools.FastPssh
		break
	case "slow":
		selectedPssh = tools.SlowPssh
	case "display":
		selectedPssh = tools.DisplayPssh

	}
	util.PrintYellow(fmt.Sprintf("Using %v PSSH", psshType))
	return selectedPssh
}
