package config

import (
	"os"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestConfig(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Config Suite")
}

var _ = Describe("KohanConfig", func() {
	Describe("NewKohanConfig", func() {
		It("should expand ${HOME} environment variables in file paths", func() {
			config, err := NewKohanConfig()
			Expect(err).ToNot(HaveOccurred())

			// Get actual HOME value for exact comparison
			homeDir := os.Getenv("HOME")

			// Verify ${HOME} was expanded to exact paths
			Expect(config.Tax.TaxDir).To(Equal(homeDir + "/Downloads/FACompute"))
			Expect(config.Tax.TickerCacheDir).To(Equal(homeDir + "/Downloads/FACompute/Tickers"))
			Expect(config.Tax.TradesPath).To(Equal(homeDir + "/Downloads/FACompute/trades.csv"))

			// Ensure no literal ${HOME} remains
			Expect(config.Tax.TaxDir).ToNot(ContainSubstring("${HOME}"))
		})

		It("should parse configuration successfully", func() {
			config, err := NewKohanConfig()
			Expect(err).ToNot(HaveOccurred())

			// Verify basic configuration is loaded
			Expect(config.Tax.AlphaBaseURL).To(Equal("https://www.alphavantage.co/query"))
			Expect(config.Tax.SBIBaseURL).To(ContainSubstring("raw.githubusercontent.com"))
			Expect(config.Tax.TaxDir).ToNot(BeEmpty())
		})
	})
})
