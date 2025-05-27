package command

import (
	"fmt"

	"github.com/amanhigh/go-fun/components/kohan/core"
	"github.com/amanhigh/go-fun/models/config"
	"github.com/spf13/cobra"
)

var dariusCmd = &cobra.Command{
	Use:   "darius",
	Short: "Kohan Commander TUI",
	RunE: func(_ *cobra.Command, _ []string) (err error) {
		config := config.DariusConfig{
			MakeDir:             makeFileDir,
			SelectedServiceFile: tmpServiceFile,
		}
		darius, err := core.GetKohanInterface().GetDariusApp(config)
		if err != nil {
			return fmt.Errorf("failed to get darius app: %w", err)
		}
		return darius.Run()
	},
}

func init() {
	RootCmd.AddCommand(dariusCmd)

	// Flags
	dariusCmd.Flags().StringVarP(&makeFileDir, "makedir", "", makeFileDir, "Makefile Directory")
	dariusCmd.Flags().StringVarP(&tmpServiceFile, "tmpsvc", "", tmpServiceFile, "Temp Service File")
}
