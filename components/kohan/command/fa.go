package command

import (
	"context"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"

	"github.com/amanhigh/go-fun/components/kohan/core"
	"github.com/amanhigh/go-fun/models/fa"
)

var (
	year  int
	faCmd = &cobra.Command{
		Use:   "fa",
		Short: "Foreign Asset (FA) related commands",
	}

	faAnalyzeCmd = &cobra.Command{
		Use:   "analyze [tickers]",
		Short: "Analyze Foreign Assets for Schedule FA",
		Long: `Analyzes one or more tickers for Schedule FA reporting.
               Provides both USD prices and INR conversions using SBI TT rates.
               Example: kohan fa analyze AMZN,SIVR`,
		Args: cobra.ExactArgs(1),
		RunE: runFAAnalyze,
	}
)

func init() {
	faAnalyzeCmd.Flags().IntVarP(&year, "year", "y", time.Now().Year(), "Tax year for analysis")
	faCmd.AddCommand(faAnalyzeCmd)
	RootCmd.AddCommand(faCmd)
}

func runFAAnalyze(cmd *cobra.Command, args []string) error {
	tickers := parseTickers(args[0])

	faManager := core.GetKohanInterface().GetFAManager()
	analysis, err := faManager.ProcessTickers(context.Background(), tickers, year)
	if err != nil {
		return err
	}

	printFAAnalysis(analysis)
	return nil
}

func parseTickers(tickerString string) []string {
	return strings.Split(tickerString, ",")
}

func printFAAnalysis(analysis []fa.TickerInfo) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer w.Flush()

	// Print Headers
	fmt.Fprintln(w, "Ticker\tPeak Date\tPeak Price (USD)\tYear-End Date\tYear-End Price (USD)\tTTBR (Peak)\tTTBR (Year-End)\tPeak Price (INR)\tYear-End Price (INR)")

	// Print Data
	for _, a := range analysis {
		fmt.Fprintf(w, "%s\t%s\t$%.2f\t%s\t$%.2f\t₹%.2f\t₹%.2f\t₹%.2f\t₹%.2f\n",
			a.TickerAnalysis.Ticker,
			a.TickerAnalysis.PeakDate,
			a.PeakPrice,
			a.TickerAnalysis.YearEndDate,
			a.YearEndPrice,
			a.PeakTTRate,
			a.YearEndTTRate,
			a.PeakPriceINR,
			a.YearEndPriceINR)
	}
}
