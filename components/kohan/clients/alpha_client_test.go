package clients_test

import (
	"github.com/amanhigh/go-fun/components/kohan/clients"
	"github.com/go-resty/resty/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("AlphaClient", func() {
	var (
		restyClient *resty.Client
	)

	BeforeEach(func() {
		restyClient = resty.New()
	})

	Context("ValidateAPIKey", func() {
		It("should return error when API key is empty", func() {
			client := clients.NewAlphaClient(restyClient, "https://test.com", "")

			err := client.ValidateAPIKey()

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("alpha Vantage API key is required for ticker download"))
		})

		It("should return error when API key is whitespace only", func() {
			client := clients.NewAlphaClient(restyClient, "https://test.com", "   ")

			err := client.ValidateAPIKey()

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("alpha Vantage API key is required for ticker download"))
		})

		It("should return nil when API key is valid", func() {
			client := clients.NewAlphaClient(restyClient, "https://test.com", "valid-api-key")

			err := client.ValidateAPIKey()

			Expect(err).ToNot(HaveOccurred())
		})
	})
})
