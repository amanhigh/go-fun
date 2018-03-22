package command

import (
	"github.com/amanhigh/go-fun/kohan/commander/components"
	"github.com/spf13/cobra"
	"github.fkinternal.com/Flipkart/elb/elb/cli/util"
)

var allCmd = &cobra.Command{
	Use:   "all",
	Short: "Uncategorized Commands",
	Long:  `Uncategorized Commands`,
	Args:  cobra.ExactArgs(1),
}

var getVersionCmd = &cobra.Command{
	Use:   "getVersion [Package Name] [Host] [Dpkg/Latest] [Comment]",
	Short: "Get Version for Package Latest or Dpkg",
	Long:  `Get Version for Package Latest or Dpkg`,
	Args:  cobra.ExactArgs(4),
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		err = util.ValidateEnumArg(args[2], []string{"dpkg", "latest"})
		return
	},
	Run: func(cmd *cobra.Command, args []string) {
		components.GetVersion(args[0], args[1], args[2], args[3])
	},
}

var printfCmd = &cobra.Command{
	Use:   "printf [Template File] [Param File]",
	Short: "Substitution Helper",
	Long:  `Substitution Helper`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		components.Printf(args[0], args[1], marker)
	},
}

func init() {
	printfCmd.Flags().StringVarP(&marker, "marker", "m", "#", "Marker in Template File")

	rootCmd.AddCommand(allCmd)
	allCmd.AddCommand(getVersionCmd,printfCmd)
}
