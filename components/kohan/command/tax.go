package command

import (
	"context"
	"fmt"
	"strconv"

	"github.com/amanhigh/go-fun/components/kohan/core"
	"github.com/amanhigh/go-fun/models/config"
	"github.com/spf13/cobra"
)

var taxCmd = &cobra.Command{
	Use:   "tax",
	Short: "Tax related commands",
}

var computeCmd = &cobra.Command{
	Use:   "compute [YEAR]",
	Short: "Compute and generate tax reports for a given year",
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		year := args[0]
		ctx := context.Background()

		// Convert year to int
		yearInt, err := strconv.Atoi(year)
		if err != nil {
			return fmt.Errorf("invalid year format: %w", err)
		}

		// Initialize config and injector
		kohanConfig, err := config.NewKohanConfig()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		core.SetupKohanInjector(kohanConfig)
		ki := core.GetKohanInterface()
		taxManager, err := ki.GetTaxManager()
		if err != nil {
			return fmt.Errorf("failed to get tax manager: %w", err)
		}

		summary, httpErr := taxManager.GetTaxSummary(ctx, yearInt)
		if httpErr != nil {
			return fmt.Errorf("failed to get tax summary (dir: %s): %w", kohanConfig.Tax.TaxDir, httpErr)
		}

		if err := taxManager.SaveTaxSummaryToExcel(ctx, yearInt, summary); err != nil {
			return fmt.Errorf("failed to save tax summary to excel: %w", err)
		}

		fmt.Printf("Successfully generated tax summary for year %v\n", year)
		return nil
	},
}

func init() {
	appsCmd.AddCommand(taxCmd)
	// TODO: #A Second Year Summary from Vested.
	// TODO: #B Add Interactive Broker Parser.
	taxCmd.AddCommand(vestedCmd)
	// TODO: #C Match Computation for Interactive Broker.
	taxCmd.AddCommand(computeCmd)
}

var vestedCmd = &cobra.Command{
	Use:   "vested",
	Short: "Generate Vested Brokerage Report",
	Long:  `Generate Vested Brokerage Report from DriveWealth Excel file`,
	RunE: func(_ *cobra.Command, _ []string) error {
		ctx := context.Background()

		// Initialize config and injector
		kohanConfig, err := config.NewKohanConfig()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		core.SetupKohanInjector(kohanConfig)
		ki := core.GetKohanInterface()
		driveWealthManager, err := ki.GetDriveWealthManager()
		if err != nil {
			return fmt.Errorf("failed to get drive wealth manager: %w", err)
		}

		// FIXME: #A Match First Year Summary from Vested (ticker download not working).
		info, err := driveWealthManager.Parse()
		if err != nil {
			return fmt.Errorf("failed to parse drive wealth report: %w", err)
		}

		if err := driveWealthManager.GenerateCsv(ctx, info); err != nil {
			return fmt.Errorf("failed to generate csv: %w", err)
		}

		fmt.Println("Successfully generated Vested Brokerage Report")
		return nil
	},
}
