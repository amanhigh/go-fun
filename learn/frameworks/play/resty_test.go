package play_test

import (
	"github.com/amanhigh/go-fun/apps/models/fun-app/db"
	"github.com/go-resty/resty/v2"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
)

type BinAnyResponse struct {
	Headers map[string]string `json:"headers"`
	//Data db.Person `json:"data"`
	Method string `json:"method"`
}

var _ = Describe("Resty", func() {

	var (
		client = resty.New()
		err    error
		resp   *resty.Response
	)

	BeforeEach(func() {
		err = nil
		resp = nil

	})

	It("should build", func() {
		Expect(client).To(Not(BeNil()))
	})

	It("should do get", func() {
		resp, err = client.R().Get("https://httpbin.org/status/200")
		Expect(err).To(BeNil())
		Expect(resp.StatusCode()).To(Equal(http.StatusOK))
	})

	Context("Custom Request", func() {
		var (
			person      db.Person
			binResponse BinAnyResponse
			headerKey   = "Myheader"
			headerValue = "MyHeaderValue"
		)

		BeforeEach(func() {
			binResponse = BinAnyResponse{}
			person = db.Person{
				Name:   "Aman",
				Age:    18,
				Gender: "Male",
			}
		})

		It("should build custom Request", func() {
			resp, err = client.R().SetHeader(headerKey, headerValue).SetBody(person).SetResult(&binResponse).Post("https://httpbin.org/anything")
			Expect(err).To(BeNil())
			Expect(resp.StatusCode()).To(Equal(http.StatusOK))
			Expect(binResponse.Method).To(Equal(http.MethodPost))
			Expect(len(binResponse.Headers)).To(BeNumerically(">", 2))
			Expect(binResponse.Headers).To(HaveKeyWithValue(headerKey, headerValue))
		})
	})

})
