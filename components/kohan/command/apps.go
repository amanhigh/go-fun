package command

import (
	"github.com/spf13/cobra"
)

var appsCmd = &cobra.Command{
	Use:   "apps",
	Short: "Manage various Kohan applications",
	Long:  `Collection of commands for managing different Kohan applications`,
}

func init() {
	RootCmd.AddCommand(appsCmd)
}
