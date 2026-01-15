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
	taxCmd.AddCommand(parseCmd)
	taxCmd.AddCommand(computeCmd)
}

var parseCmd = &cobra.Command{
	Use:   "parse YEAR",
	Short: "Parse all broker files and generate CSVs",
	Long:  `Auto-detects and parses DriveWealth, Interactive Brokers files for the specified year. Merges and generates consolidated CSVs.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		ctx := context.Background()

		// Parse year from arguments
		year, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid year format: %w", err)
		}

		kohanConfig, err := config.NewKohanConfig()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		core.SetupKohanInjector(kohanConfig)
		ki := core.GetKohanInterface()

		brokerageManager, err := ki.GetBrokerageManager()
		if err != nil {
			return fmt.Errorf("failed to get brokerage manager: %w", err)
		}

		if err := brokerageManager.ParseAndGenerate(ctx, year); err != nil {
			return fmt.Errorf("failed to parse brokers: %w", err)
		}

		fmt.Printf("Successfully parsed broker files for year %d and generated CSVs\n", year)
		return nil
	},
}
