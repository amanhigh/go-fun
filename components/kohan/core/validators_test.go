package core_test

import (
	"github.com/go-playground/validator/v10"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/amanhigh/go-fun/components/kohan/core"
)

var _ = Describe("Validators", func() {
	var (
		v *validator.Validate
	)

	BeforeEach(func() {
		v = validator.New()
		// Register validators with the test validator instance
		_ = v.RegisterValidation("ticker", core.TickerValidator)
		_ = v.RegisterValidation("tag", core.TagValidator)
		_ = v.RegisterValidation("override", core.OverrideValidator)
		_ = v.RegisterValidation("future_date", core.FutureDateValidator)
	})

	Describe("ticker validator", func() {
		Context("valid tickers", func() {
			It("should accept uppercase letters", func() {
				Expect(v.Var("TCS", "ticker")).To(Succeed())
			})

			It("should accept uppercase letters with digits", func() {
				Expect(v.Var("TCS123", "ticker")).To(Succeed())
			})

			It("should accept dots", func() {
				Expect(v.Var("TCS.NS", "ticker")).To(Succeed())
			})

			It("should accept hyphens", func() {
				Expect(v.Var("TCS-TEST", "ticker")).To(Succeed())
			})

			It("should accept ampersands", func() {
				Expect(v.Var("M&M", "ticker")).To(Succeed())
			})

			It("should accept complex combinations", func() {
				Expect(v.Var("TCS.NS-123&TEST", "ticker")).To(Succeed())
			})

			It("should accept empty string", func() {
				Expect(v.Var("", "ticker")).To(Succeed())
			})
		})

		Context("invalid tickers", func() {
			It("should reject lowercase letters", func() {
				Expect(v.Var("tcs", "ticker")).ToNot(Succeed())
			})

			It("should reject spaces", func() {
				Expect(v.Var("TCS NS", "ticker")).ToNot(Succeed())
			})

			It("should reject special characters other than allowed ones", func() {
				Expect(v.Var("TCS@NS", "ticker")).ToNot(Succeed())
				Expect(v.Var("TCS#TEST", "ticker")).ToNot(Succeed())
				Expect(v.Var("TCS+TEST", "ticker")).ToNot(Succeed())
			})

			It("should reject starting with special character", func() {
				Expect(v.Var("-TCS", "ticker")).ToNot(Succeed())
				Expect(v.Var(".TCS", "ticker")).ToNot(Succeed())
				Expect(v.Var("&TCS", "ticker")).ToNot(Succeed())
			})
		})
	})

	Describe("tag validator", func() {
		Context("valid tags", func() {
			It("should accept lowercase letters", func() {
				Expect(v.Var("oe", "tag")).To(Succeed())
			})

			It("should accept uppercase letters", func() {
				Expect(v.Var("TEST", "tag")).To(Succeed())
			})

			It("should accept digits", func() {
				Expect(v.Var("tag123", "tag")).To(Succeed())
			})

			It("should accept hyphens", func() {
				Expect(v.Var("dep-1", "tag")).To(Succeed())
				Expect(v.Var("test-tag-name", "tag")).To(Succeed())
			})

			It("should accept combinations", func() {
				Expect(v.Var("Test-123", "tag")).To(Succeed())
				Expect(v.Var("dep-2-test", "tag")).To(Succeed())
			})

			It("should accept empty string", func() {
				Expect(v.Var("", "tag")).To(Succeed())
			})
		})

		Context("invalid tags", func() {
			It("should reject spaces", func() {
				Expect(v.Var("dep 1", "tag")).ToNot(Succeed())
			})

			It("should reject special characters", func() {
				Expect(v.Var("dep@1", "tag")).ToNot(Succeed())
				Expect(v.Var("dep#1", "tag")).ToNot(Succeed())
				Expect(v.Var("dep.1", "tag")).ToNot(Succeed())
				Expect(v.Var("dep&1", "tag")).ToNot(Succeed())
			})

			It("should accept starting with hyphen", func() {
				Expect(v.Var("-tag", "tag")).To(Succeed())
			})

			It("should accept ending with hyphen", func() {
				Expect(v.Var("tag-", "tag")).To(Succeed())
			})

			It("should accept consecutive hyphens", func() {
				Expect(v.Var("tag--name", "tag")).To(Succeed())
			})
		})
	})

	Describe("override validator", func() {
		Context("valid overrides", func() {
			It("should accept lowercase letters", func() {
				Expect(v.Var("loc", "override")).To(Succeed())
			})

			It("should accept uppercase letters", func() {
				Expect(v.Var("LOC", "override")).To(Succeed())
			})

			It("should accept mixed case letters", func() {
				Expect(v.Var("Location", "override")).To(Succeed())
				Expect(v.Var("abcDef", "override")).To(Succeed())
			})

			It("should accept empty string", func() {
				Expect(v.Var("", "override")).To(Succeed())
			})
		})

		Context("invalid overrides", func() {
			It("should reject spaces", func() {
				Expect(v.Var("test location", "override")).ToNot(Succeed())
			})

			It("should reject special characters", func() {
				Expect(v.Var("test-loc", "override")).ToNot(Succeed())
				Expect(v.Var("test_loc", "override")).ToNot(Succeed())
				Expect(v.Var("test.loc", "override")).ToNot(Succeed())
				Expect(v.Var("test@loc", "override")).ToNot(Succeed())
				Expect(v.Var("test#loc", "override")).ToNot(Succeed())
				Expect(v.Var("test&loc", "override")).ToNot(Succeed())
			})

			It("should reject digits", func() {
				Expect(v.Var("123", "override")).ToNot(Succeed())
				Expect(v.Var("abc1", "override")).ToNot(Succeed())
				Expect(v.Var("ABC123", "override")).ToNot(Succeed())
				Expect(v.Var("test123name", "override")).ToNot(Succeed())
			})
		})
	})

	Describe("reviewed_at validator", func() {
		Context("future date business rule", func() {
			It("should accept empty string", func() {
				Expect(v.Var("", "future_date")).To(Succeed())
			})

			It("should accept past dates", func() {
				Expect(v.Var("2024-01-16", "future_date")).To(Succeed())
				Expect(v.Var("2023-12-31", "future_date")).To(Succeed())
				Expect(v.Var("2020-01-01", "future_date")).To(Succeed())
			})

			It("should accept today's date", func() {
				today := "2024-01-15" // Assuming test runs on 2024-01-15 or later
				Expect(v.Var(today, "future_date")).To(Succeed())
			})

			It("should reject future dates", func() {
				Expect(v.Var("2099-12-31", "future_date")).ToNot(Succeed())
				Expect(v.Var("2100-01-01", "future_date")).ToNot(Succeed())
				Expect(v.Var("2030-07-15", "future_date")).ToNot(Succeed())
			})

			It("should pass through invalid formats to datetime validator", func() {
				// Custom validator should return true for invalid formats, letting datetime validator handle them
				Expect(v.Var("invalid-date", "future_date")).To(Succeed())
				Expect(v.Var("2024-13-32", "future_date")).To(Succeed())
				Expect(v.Var("2024/01/16", "future_date")).To(Succeed())
			})
		})
	})
})
