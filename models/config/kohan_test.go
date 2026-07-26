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
			Expect(config.Tax.TickerCacheDir).To(Equal(homeDir + "/Downloads/FACompute/Data/Tickers"))
			Expect(config.Tax.TradesPath).To(Equal(homeDir + "/Downloads/FACompute/Input/Parsed/trades.csv"))
			Expect(config.Barkat.ScreenshotPath).To(Equal(homeDir + "/Downloads/Screenshots"))

			// Ensure no literal ${HOME} remains
			Expect(config.Tax.TaxDir).ToNot(ContainSubstring("${HOME}"))
		})

		Context("Explicitly configured paths with ${HOME} expansion", func() {
			var originalTaxDir string

			BeforeEach(func() {
				originalTaxDir = os.Getenv("TAX_DIR")
				os.Setenv("TAX_DIR", "/tmp/${HOME}/custom-tax")
			})

			AfterEach(func() {
				if originalTaxDir == "" {
					os.Unsetenv("TAX_DIR")
				} else {
					os.Setenv("TAX_DIR", originalTaxDir)
				}
			})

			It("should expand ${HOME} in explicitly configured TAX_DIR path", func() {
				config, err := NewKohanConfig()
				Expect(err).ToNot(HaveOccurred())

				homeDir := os.Getenv("HOME")
				Expect(config.Tax.TaxDir).To(Equal("/tmp/" + homeDir + "/custom-tax"))
			})
		})

		It("should parse configuration successfully", func() {
			config, err := NewKohanConfig()
			Expect(err).ToNot(HaveOccurred())

			// Verify basic configuration is loaded
			Expect(config.Tax.YahooBaseURL).To(Equal("https://query1.finance.yahoo.com"))
			Expect(config.Tax.SBIBaseURL).To(ContainSubstring("raw.githubusercontent.com"))
			Expect(config.Tax.TaxDir).ToNot(BeEmpty())
			Expect(config.Barkat.ScreenshotPath).ToNot(BeEmpty())
		})
	})
})
