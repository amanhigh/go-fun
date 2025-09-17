package clients_test

import (
	"context"
	"net/http"
	"net/http/httptest"

	"github.com/amanhigh/go-fun/components/kohan/clients"
	"github.com/go-resty/resty/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("SBIClient", func() {
	var (
		server *httptest.Server
		client *clients.SBIClientImpl
		ctx    context.Context
	)

	BeforeEach(func() {
		ctx = context.Background()
	})

	AfterEach(func() {
		server.Close()
	})

	Context("FetchExchangeRates", func() {
		BeforeEach(func() {
			server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "text/csv")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("USD,1.0,1.1"))
			}))
			client = clients.NewSBIClient(resty.New(), server.URL)
		})

		It("should fetch rates successfully", func() {
			result, err := client.FetchExchangeRates(ctx)
			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(Equal("USD,1.0,1.1"))
		})
	})
})
