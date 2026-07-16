package tax_test

import (
	"math"
	"net/http"
	"time"

	. "github.com/amanhigh/go-fun/models/tax"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("YahooSplit", func() {
	const ticker = "TEST"

	type httpErr interface {
		Code() int
		Error() string
	}

	assertInvalidSplit := func(err httpErr, category string) {
		Expect(err).To(HaveOccurred())
		Expect(err.Code()).To(Equal(http.StatusBadRequest))
		Expect(err.Error()).To(ContainSubstring(ticker))
		Expect(err.Error()).To(ContainSubstring(category))
	}

	Describe("Validate", func() {
		Context("valid split", func() {
			It("should accept valid timestamp and positive finite numerator/denominator", func() {
				split := YahooSplit{Date: 1609459200, Numerator: 5, Denominator: 1}
				Expect(split.Validate(ticker)).To(Succeed())
			})
		})

		Context("invalid split", func() {
			It("should reject non-positive timestamp with HTTP 400 and ticker/timestamp context", func() {
				By("zero timestamp")
				err := YahooSplit{Date: 0, Numerator: 5, Denominator: 1}.Validate(ticker)
				assertInvalidSplit(err, "non-positive timestamp")

				By("negative timestamp")
				err = YahooSplit{Date: -1, Numerator: 5, Denominator: 1}.Validate(ticker)
				assertInvalidSplit(err, "non-positive timestamp")
			})

			It("should reject zero or negative numerator with HTTP 400 and ticker context", func() {
				By("zero numerator")
				split := YahooSplit{Date: 1609459200, Numerator: 0, Denominator: 1}
				err := split.Validate(ticker)
				assertInvalidSplit(err, "numerator")

				By("negative numerator")
				err = YahooSplit{Date: 1609459200, Numerator: -1, Denominator: 1}.Validate(ticker)
				assertInvalidSplit(err, "numerator")
			})

			It("should reject NaN or +Inf numerator with HTTP 400 and ticker context", func() {
				By("NaN numerator")
				split := YahooSplit{Date: 1609459200, Numerator: math.NaN(), Denominator: 1}
				err := split.Validate(ticker)
				assertInvalidSplit(err, "numerator")

				By("+Inf numerator")
				err = YahooSplit{Date: 1609459200, Numerator: math.Inf(1), Denominator: 1}.Validate(ticker)
				assertInvalidSplit(err, "numerator")
			})

			It("should reject zero or negative denominator with HTTP 400 and ticker context", func() {
				By("zero denominator")
				split := YahooSplit{Date: 1609459200, Numerator: 5, Denominator: 0}
				err := split.Validate(ticker)
				assertInvalidSplit(err, "denominator")

				By("negative denominator")
				err = YahooSplit{Date: 1609459200, Numerator: 5, Denominator: -1}.Validate(ticker)
				assertInvalidSplit(err, "denominator")
			})

			It("should reject NaN or +Inf denominator with HTTP 400 and ticker context", func() {
				By("NaN denominator")
				split := YahooSplit{Date: 1609459200, Numerator: 5, Denominator: math.NaN()}
				err := split.Validate(ticker)
				assertInvalidSplit(err, "denominator")

				By("+Inf denominator")
				err = YahooSplit{Date: 1609459200, Numerator: 5, Denominator: math.Inf(1)}.Validate(ticker)
				assertInvalidSplit(err, "denominator")
			})
		})
	})

	Describe("Ratio", func() {
		It("should return expected forward ratio for valid splits", func() {
			Expect(YahooSplit{Numerator: 5, Denominator: 1}.Ratio()).To(Equal(5.0))
		})

		It("should return expected reverse ratio for valid splits", func() {
			Expect(YahooSplit{Numerator: 1, Denominator: 5}.Ratio()).To(Equal(0.2))
		})
	})

	Describe("EffectiveDate", func() {
		It("should convert non-midnight Unix timestamp to UTC calendar midnight", func() {
			// Jan 2, 2021 14:30:00 UTC → should truncate to Jan 2, 2021 00:00:00 UTC
			split := YahooSplit{Date: 1609597800}
			expected := time.Date(2021, time.January, 2, 0, 0, 0, 0, time.UTC)
			Expect(split.EffectiveDate()).To(Equal(expected))
		})
	})
})
