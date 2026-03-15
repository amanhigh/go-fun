package core_test

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-sql/civil"
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
		_ = v.RegisterValidation("not_future", core.NotFutureValidator)
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

	Describe("not_future validator", func() {
		Context("future date business rule", func() {
			type TestStruct struct {
				ReviewedAt *civil.Date `validate:"not_future"`
			}

			It("should accept past dates", func() {
				pastDate1 := civil.Date{Year: 2024, Month: 1, Day: 16}
				pastDate2 := civil.Date{Year: 2023, Month: 12, Day: 31}
				pastDate3 := civil.Date{Year: 2020, Month: 1, Day: 1}

				test1 := TestStruct{ReviewedAt: &pastDate1}
				test2 := TestStruct{ReviewedAt: &pastDate2}
				test3 := TestStruct{ReviewedAt: &pastDate3}

				Expect(v.Struct(test1)).To(Succeed())
				Expect(v.Struct(test2)).To(Succeed())
				Expect(v.Struct(test3)).To(Succeed())
			})

			It("should accept today's date", func() {
				today := civil.DateOf(time.Now())
				test := TestStruct{ReviewedAt: &today}
				Expect(v.Struct(test)).To(Succeed())
			})

			It("should reject future dates", func() {
				futureDate1 := civil.Date{Year: 2099, Month: 12, Day: 31}
				futureDate2 := civil.Date{Year: 2100, Month: 1, Day: 1}
				futureDate3 := civil.Date{Year: 2030, Month: 7, Day: 15}

				test1 := TestStruct{ReviewedAt: &futureDate1}
				test2 := TestStruct{ReviewedAt: &futureDate2}
				test3 := TestStruct{ReviewedAt: &futureDate3}

				Expect(v.Struct(test1)).ToNot(Succeed())
				Expect(v.Struct(test2)).ToNot(Succeed())
				Expect(v.Struct(test3)).ToNot(Succeed())
			})
		})
	})
})
