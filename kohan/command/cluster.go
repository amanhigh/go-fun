package command

import (
	"github.com/amanhigh/go-fun/kohan/commander/components"
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

func init() {
	rootCmd.AddCommand(clusterCmd)

	clusterCmd.AddCommand(clusterSanityCmd)
}
