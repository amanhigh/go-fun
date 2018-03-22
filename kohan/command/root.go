package command

import (
	"github.com/spf13/cobra"
	. "github.fkinternal.com/Flipkart/elb/elb/util/os"
	"os"
)

var (
	rootCmd = &cobra.Command{}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		PrintRed(err.Error())
		os.Exit(1)
	} else {
		PrintGreen("Command Successful")
	}
}