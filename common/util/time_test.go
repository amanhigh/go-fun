package util_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/amanhigh/go-fun/common/util"
)

var _ = Describe("Time", func() {
	var _ = Describe("Util", func() {

		Context("FormatTime", func() {
			It("formats time correctly", func() {
				t := time.Date(2022, 4, 1, 14, 30, 0, 0, time.UTC)
				formatted := util.FormatTime(t, util.PRINT_LAYOUT)
				Expect(formatted).To(Equal("Apr 1, 2022 at 2:30pm (UTC)"))
			})
		})

		Context("TimeAfter", func() {
			It("calculates time after correctly", func() {
				duration, _ := time.ParseDuration("48h")
				expected := time.Now().Add(duration)
				result := util.TimeAfter(duration)
				Expect(result).To(BeTemporally("~", expected, time.Second))
			})
		})

		Context("DaysHour", func() {
			It("computes date and time correctly", func() {
				computedTime := util.DaysHour("24 4:05AM")
				Expect(computedTime.Day()).To(Equal(24))
				Expect(computedTime.Hour()).To(Equal(4))
				Expect(computedTime.Minute()).To(Equal(5))
				Expect(computedTime.Second()).To(Equal(0))
				Expect(computedTime.Nanosecond()).To(Equal(0))
				Expect(computedTime.Location()).To(Equal(time.UTC))
			})

			It("should return an same string for invalid date", func() {
				t, err := time.Parse(time.RFC3339, "2022-03-30T14:30:00Z")
				Expect(err).NotTo(HaveOccurred())

				layout := "invalid layout"
				result := util.FormatTime(t, layout)
				Expect(result).To(Equal(layout))
			})
		})

	})
})
