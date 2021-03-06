package command

import (
	"fmt"
	"strings"

	"github.com/amanhigh/go-fun/apps/common/tools"
	"github.com/amanhigh/go-fun/apps/components/kohan/core"
	"github.com/amanhigh/go-fun/util"
	"github.com/spf13/cobra"
)

var allCmd = &cobra.Command{
	Use:   "all",
	Short: "Uncategorized Commands",
	Args:  cobra.ExactArgs(1),
}

var getVersionCmd = &cobra.Command{
	Use:   "getVersion [Package Name] [Dpkg/Latest] [Host] [Comment]",
	Short: "Get Version for Package Latest or Dpkg",
	Args:  cobra.ExactArgs(4),
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		err = util.ValidateEnumArg(args[1], []string{"dpkg", "latest"})
		return
	},
	Run: func(cmd *cobra.Command, args []string) {
		core.GetVersion(args[0], args[2], args[1], args[3])
	},
}

var printfCmd = &cobra.Command{
	Use:   "printf [Template File] [Param File]",
	Short: "Substitution Helper",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		core.Printf(args[0], args[1], marker)
	},
}

var syncCmd = &cobra.Command{
	Use:   "sync [srcHost] [srcDir] [targetHost(s) Space Separated]",
	Short: "Syncs Remote Host Directory with target hosts",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		tools.Sync(args[0], args[1], args[1], strings.Fields(args[2]))
	},
}

var investingCmd = &cobra.Command{
	Use:   "investing fileDirectory",
	Short: "Converts Historical Data Format",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		for _, file := range util.ListFiles(args[0]) {
			fmt.Println(fmt.Sprintf("Processing: %v", file))
			if err = core.ReformatInvestingFile(file); err != nil {
				break
			}
		}
		return
	},
}

func init() {
	printfCmd.Flags().StringVarP(&marker, "marker", "m", "#", "Marker in Template File")

	RootCmd.AddCommand(allCmd)
	allCmd.AddCommand(getVersionCmd, printfCmd, syncCmd, investingCmd)
}
