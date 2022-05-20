package tutorial_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

var _ = Describe("Generic", func() {

	var (
		// Initialize a map for the integer values
		ints = map[string]int64{
			"first":  34,
			"second": 12,
		}

		// Initialize a map for the float values
		floats = map[string]float64{
			"first":  35.98,
			"second": 26.99,
		}

		intArray   = []int{1, 2, 3}
		floatArray = []float64{12.3, 2.6, 4.53}
	)

	Context("Custom", func() {
		It("should sum ints", func() {
			Expect(SumNumbers(ints)).To(BeEquivalentTo(46))
		})

		It("should sum floats", func() {
			Expect(SumNumbers(floats)).To(Equal(62.97))
		})
	})

	Context("Slices", func() {
		It("should do contain int", func() {
			Expect(slices.Contains(intArray, 3)).To(BeTrue())
			Expect(slices.Contains(intArray, 7)).To(BeFalse())
		})

		It("should do contain float", func() {
			Expect(slices.Contains(floatArray, 4.53)).To(BeTrue())
			Expect(slices.Contains(floatArray, 7.77)).To(BeFalse())
		})
	})

	Context("Maps", func() {
		It("should get values for int", func() {
			Expect(maps.Values(ints)).To(ContainElements([]int64{34, 12}))
		})

		It("should get values for float", func() {
			Expect(maps.Values(floats)).To(Equal([]float64{35.98, 26.99}))
		})
	})

})

/* Declare Genric Type which obeys few interfaces */
type Number interface {
	int64 | float64
}

// SumNumbers sums the values of map m. It supports both integers
// and floats as map values.
func SumNumbers[K comparable, V Number](m map[K]V) V {
	var s V
	for _, v := range m {
		s += v
	}
	return s
}
