package play_fast

import (
	"time"

	"github.com/jinzhu/copier"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gmeasure"
	deepcopy "github.com/tiendc/go-deepcopy"
)

// FIXME: Remove Copier in all Names below.
// Types for Copier tests (must be at package level for method definitions)
type copierCoffee struct {
	Name       string
	Origin     string
	RoastLevel string
	Price      float64
}

type copierCoffeeDTO struct {
	Name       string
	Origin     string
	RoastLevel string
	Price      float64
}

type copierCoffeeWithTag struct {
	Name       string
	Origin     string `copier:"-"`
	RoastLevel string
	Price      float64
}

// HACK: Reuse struct with diffierent fields for different Tags.
type copierCoffeeRequired struct {
	Name       string `copier:"must"`
	Origin     string
	RoastLevel string
	Price      float64
}

type copierCoffeeRequiredNoPanic struct {
	Name       string `copier:"must,nopanic"`
	Origin     string
	RoastLevel string
	Price      float64
}

type copierEmployee struct {
	Name string
	Age  int32
	Role string
}

// DoubleAge is a method that copier can use for method-to-field copying
func (e *copierEmployee) DoubleAge() int32 {
	return 2 * e.Age
}

type copierManager struct {
	Name      string
	Age       int32
	DoubleAge int32
	SuperRole string
}

// Role is a setter method that copier can use for field-to-method copying
func (m *copierManager) Role(role string) {
	m.SuperRole = "Super " + role
}

type copierAddress struct {
	City    string
	Country string
}

type copierPerson struct {
	Name    string
	Age     int
	Address copierAddress
	Tags    []string
	Scores  map[string]int
}

type copierPersonDTO struct {
	Name    string
	Age     int
	Address copierAddress
	Tags    []string
	Scores  map[string]int
}

var _ = Describe("Copy", func() {

	Context("Copier", func() {

		It("should build", func() {
			src := copierCoffee{}
			dst := copierCoffeeDTO{}
			err := copier.Copy(&dst, &src)
			Expect(err).ToNot(HaveOccurred())
		})

		Context("Basic", func() {
			var (
				src copierCoffee
				dst copierCoffeeDTO
			)

			BeforeEach(func() {
				src = copierCoffee{
					Name:       "Ethiopian Yirgacheffe",
					Origin:     "Ethiopia",
					RoastLevel: "Light",
					Price:      18.99,
				}
				err := copier.Copy(&dst, &src)
				Expect(err).ToNot(HaveOccurred())
			})

			It("1.1 should copy simple struct fields", func() {
				Expect(dst.Name).To(Equal(src.Name))
				Expect(dst.Origin).To(Equal(src.Origin))
				Expect(dst.RoastLevel).To(Equal(src.RoastLevel))
				Expect(dst.Price).To(Equal(src.Price))
			})

			It("1.2 should copy nested struct fields", func() {
				src := copierPerson{
					Name:    "Alice",
					Age:     30,
					Address: copierAddress{City: "Seattle", Country: "USA"},
				}
				var dst copierPersonDTO
				err := copier.Copy(&dst, &src)
				Expect(err).ToNot(HaveOccurred())
				Expect(dst.Address.City).To(Equal("Seattle"))
				Expect(dst.Address.Country).To(Equal("USA"))
			})
		})

		Context("Medium", func() {
			Context("Field Mapping and Transformation", func() {
				It("2.1 should copy from getter methods to fields (method-to-field)", func() {
					By("Copying Employee with DoubleAge method to Manager with DoubleAge field")
					emp := copierEmployee{Name: "John", Age: 25, Role: "Admin"}
					var mgr copierManager
					err := copier.Copy(&mgr, &emp)
					Expect(err).ToNot(HaveOccurred())

					Expect(mgr.Name).To(Equal("John"))
					Expect(mgr.Age).To(Equal(int32(25)))
					Expect(mgr.DoubleAge).To(Equal(int32(50)))
				})

				It("2.2 should copy from fields to setter methods (field-to-method)", func() {
					By("Copying Employee Role field via Manager Role setter method")
					emp := copierEmployee{Name: "John", Age: 25, Role: "Admin"}
					var mgr copierManager
					err := copier.Copy(&mgr, &emp)
					Expect(err).ToNot(HaveOccurred())

					Expect(mgr.SuperRole).To(Equal("Super Admin"))
				})

				It("2.3 should map source fields to different destination field names", func() {
					type Src struct {
						FullName string
						Years    int
					}
					type Dst struct {
						Name string
						Age  int
					}

					src := Src{FullName: "Alice", Years: 30}
					var dst Dst
					err := copier.CopyWithOption(&dst, &src, copier.Option{
						FieldNameMapping: []copier.FieldNameMapping{
							{SrcType: Src{}, DstType: Dst{}, Mapping: map[string]string{
								"FullName": "Name",
								"Years":    "Age",
							}},
						},
					})
					Expect(err).ToNot(HaveOccurred())
					Expect(dst.Name).To(Equal("Alice"))
					Expect(dst.Age).To(Equal(30))
				})
			})

			Context("Configuration Merging and Overriding", func() {
				It("2.4 should merge configurations from multiple sources", func() {
					base := copierCoffee{Name: "House Blend", Origin: "Colombia", RoastLevel: "Medium", Price: 12.99}
					override := copierCoffee{Name: "House Blend Special", Price: 14.99}

					var result copierCoffee
					err := copier.Copy(&result, &base)
					Expect(err).ToNot(HaveOccurred())

					By("Applying override with IgnoreEmpty to preserve non-overridden fields")
					err = copier.CopyWithOption(&result, &override, copier.Option{IgnoreEmpty: true})
					Expect(err).ToNot(HaveOccurred())

					Expect(result.Name).To(Equal("House Blend Special"))
					Expect(result.Origin).To(Equal("Colombia"))
					Expect(result.RoastLevel).To(Equal("Medium"))
					Expect(result.Price).To(Equal(14.99))
				})

				It("2.5 should override all fields without IgnoreEmpty", func() {
					base := copierCoffee{Name: "Base", Origin: "Kenya", RoastLevel: "Dark", Price: 15.0}
					override := copierCoffee{Name: "Override", Price: 20.0}

					err := copier.Copy(&base, &override)
					Expect(err).ToNot(HaveOccurred())

					Expect(base.Name).To(Equal("Override"))
					Expect(base.Origin).To(BeEmpty()) // Overwritten with zero value
					Expect(base.Price).To(Equal(20.0))
				})
			})

			Context("Map Merging", func() {
				It("2.6 should copy maps", func() {
					src := map[string]int{"a": 1, "b": 2}
					dst := map[string]int{"b": 3, "c": 4}
					err := copier.Copy(&dst, &src)
					Expect(err).ToNot(HaveOccurred())

					Expect(dst["a"]).To(Equal(1))
					Expect(dst["b"]).To(Equal(2)) // Overridden by source
					Expect(dst["c"]).To(Equal(4)) // Preserved from destination
				})
			})

			Context("Field Filtering", func() {
				It("2.7 should ignore fields with copier:\"-\" tag", func() {
					src := copierCoffee{Name: "Espresso", Origin: "Italy", RoastLevel: "Dark", Price: 22.0}
					var dst copierCoffeeWithTag
					err := copier.Copy(&dst, &src)
					Expect(err).ToNot(HaveOccurred())

					Expect(dst.Name).To(Equal("Espresso"))
					Expect(dst.Origin).To(BeEmpty()) // Ignored via tag
					Expect(dst.RoastLevel).To(Equal("Dark"))
				})
			})

			Context("Type Conversion", func() {
				It("2.8 should copy between compatible types", func() {
					type Src struct {
						Name string
						Age  int
					}
					type Dst struct {
						Name string
						Age  int64
					}

					src := Src{Name: "Alice", Age: 30}
					var dst Dst
					err := copier.Copy(&dst, &src)
					Expect(err).ToNot(HaveOccurred())
					Expect(dst.Name).To(Equal("Alice"))
					Expect(dst.Age).To(Equal(int64(30)))
				})

				It("2.9 should handle type conversion failures gracefully", func() {
					type Src struct {
						Value string
					}
					type Dst struct {
						Value int
					}

					src := Src{Value: "not-a-number"}
					var dst Dst
					// Copier silently skips incompatible types (no error, field stays zero)
					err := copier.Copy(&dst, &src)
					Expect(err).ToNot(HaveOccurred())
					Expect(dst.Value).To(Equal(0))
				})
			})

			Context("Error Handling", func() {
				It("1.3 should handle nil pointers gracefully", func() {
					var src *copierCoffee
					dst := &copierCoffeeDTO{}
					err := copier.Copy(dst, src)
					Expect(err).To(HaveOccurred())
				})

				It("1.4 should handle empty structs", func() {
					src := copierCoffee{}
					var dst copierCoffeeDTO
					err := copier.Copy(&dst, &src)
					Expect(err).ToNot(HaveOccurred())
					Expect(dst.Name).To(BeEmpty())
				})

				It("1.2 should copy nested struct fields", func() {
					src := copierPerson{
						Name:    "Alice",
						Age:     30,
						Address: copierAddress{City: "Seattle", Country: "USA"},
					}
					var dst copierPersonDTO
					err := copier.Copy(&dst, &src)
					Expect(err).ToNot(HaveOccurred())
					Expect(src.Address.City).To(Equal(dst.Address.City))
					Expect(src.Address.Country).To(Equal(dst.Address.Country))
				})
				type Src struct {
					Name string
					Age  int
				}
				type Dst struct {
					Name string
					Age  int64
				}

				src := Src{Name: "Alice", Age: 30}
				var dst Dst
				err := copier.Copy(&dst, &src)
				Expect(err).ToNot(HaveOccurred())
				Expect(dst.Name).To(Equal("Alice"))
				Expect(dst.Age).To(Equal(int64(30)))
			})

			It("should handle type conversion failures gracefully", func() {
				type Src struct {
					Value string
				}
				type Dst struct {
					Value int
				}

				src := Src{Value: "not-a-number"}
				var dst Dst
				// Copier silently skips incompatible types (no error, field stays zero)
				err := copier.Copy(&dst, &src)
				Expect(err).ToNot(HaveOccurred())
				Expect(dst.Value).To(Equal(0))
			})
		})
	})

	Context("DeepCopy", func() {

		// FIXME: Have common types for both Copier and DeepCopy
		type Inner struct {
			Value string
		}

		type Nested struct {
			Name   string
			Inner  *Inner
			Tags   []string
			Scores map[string]int
		}

		It("should build", func() {
			src := Inner{Value: "test"}
			var dst Inner
			err := deepcopy.Copy(&dst, &src)
			Expect(err).ToNot(HaveOccurred())
			Expect(dst.Value).To(Equal("test"))
		})

		Context("Deep Copy Basics", func() {
			It("should clone simple structs without shared references", func() {
				// FIXME: make example more meaningful with real world example
				src := Inner{Value: "original"}
				var dst Inner
				err := deepcopy.Copy(&dst, &src)
				Expect(err).ToNot(HaveOccurred())

				Expect(dst.Value).To(Equal("original"))

				By("Verifying no shared references")
				src.Value = "modified"
				Expect(dst.Value).To(Equal("original"))
			})
		})

		Context("Nested Structure Cloning", func() {
			var (
				src Nested
				dst Nested
			)

			BeforeEach(func() {
				src = Nested{
					Name:   "root",
					Inner:  &Inner{Value: "nested"},
					Tags:   []string{"a", "b", "c"},
					Scores: map[string]int{"math": 95, "english": 88},
				}
				err := deepcopy.Copy(&dst, &src)
				Expect(err).ToNot(HaveOccurred())
			})

			It("should deep copy all fields", func() {
				Expect(dst.Name).To(Equal("root"))
				Expect(dst.Inner.Value).To(Equal("nested"))
				Expect(dst.Tags).To(Equal([]string{"a", "b", "c"}))
				Expect(dst.Scores).To(HaveKeyWithValue("math", 95))
			})

			It("should have no shared pointer references", func() {
				src.Inner.Value = "changed"
				Expect(dst.Inner.Value).To(Equal("nested"))
			})

			It("should have no shared slice references", func() {
				src.Tags[0] = "modified"
				Expect(dst.Tags[0]).To(Equal("a"))
			})

			It("should have no shared map references", func() {
				src.Scores["math"] = 0
				Expect(dst.Scores["math"]).To(Equal(95))
			})
		})

		Context("Pointer Handling", func() {
			// FIXME: Add Operation Ids matching with PRD.
			It("should clone structs with pointer fields independently", func() {
				type WithPointers struct {
					IntPtr    *int
					StringPtr *string
				}

				intVal := 42
				strVal := "hello"
				src := WithPointers{IntPtr: &intVal, StringPtr: &strVal}
				var dst WithPointers
				err := deepcopy.Copy(&dst, &src)
				Expect(err).ToNot(HaveOccurred())

				By("Verifying values are copied")
				Expect(*dst.IntPtr).To(Equal(42))
				Expect(*dst.StringPtr).To(Equal("hello"))

				By("Verifying pointed-to data is independent")
				*src.IntPtr = 100
				*src.StringPtr = "world"
				Expect(*dst.IntPtr).To(Equal(42))
				Expect(*dst.StringPtr).To(Equal("hello"))
			})

			It("should handle nil pointers", func() {
				src := Nested{Name: "no-inner", Inner: nil}
				var dst Nested
				err := deepcopy.Copy(&dst, &src)
				Expect(err).ToNot(HaveOccurred())
				Expect(dst.Name).To(Equal("no-inner"))
				Expect(dst.Inner).To(BeNil())
			})
		})

		Context("Slice and Array Cloning", func() {
			It("should deep copy slices", func() {
				src := []int{1, 2, 3, 4, 5}
				var dst []int
				err := deepcopy.Copy(&dst, &src)
				Expect(err).ToNot(HaveOccurred())

				Expect(dst).To(Equal([]int{1, 2, 3, 4, 5}))

				By("Verifying modifications to clone don't affect original")
				dst[0] = 100
				Expect(src[0]).To(Equal(1))
			})

			It("should deep copy slice of pointers", func() {
				a, b := 1, 2
				src := []*int{&a, &b}
				var dst []*int
				err := deepcopy.Copy(&dst, &src)
				Expect(err).ToNot(HaveOccurred())

				Expect(*dst[0]).To(Equal(1))
				*src[0] = 100
				Expect(*dst[0]).To(Equal(1))
			})
		})

		Context("Map Cloning", func() {
			It("should deep copy map fields", func() {
				src := map[string][]int{
					"a": {1, 2, 3},
					"b": {4, 5, 6},
				}
				var dst map[string][]int
				err := deepcopy.Copy(&dst, &src)
				Expect(err).ToNot(HaveOccurred())

				Expect(dst["a"]).To(Equal([]int{1, 2, 3}))

				By("Verifying clone has independent map instance")
				src["a"][0] = 100
				Expect(dst["a"][0]).To(Equal(1))
			})
		})

		Context("Interface Fields", func() {
			It("should deep copy structs with interface fields", func() {
				// HACK: Don't use With Prefix have more meaningful names.
				type WithInterface struct {
					Name  string
					Value interface{}
				}

				src := WithInterface{Name: "test", Value: map[string]int{"a": 1}}
				var dst WithInterface
				err := deepcopy.Copy(&dst, &src)
				Expect(err).ToNot(HaveOccurred())

				Expect(dst.Name).To(Equal("test"))
				Expect(dst.Value).To(Equal(map[string]int{"a": 1}))
			})
		})

		Context("Performance Benchmarks", FlakeAttempts(3), func() {
			// FIXME: Time the benchmark and simplify if it takes more time.
			// BUG: Change to Gingk Benchmark not it block for all benchmarks in this package.
			It("should benchmark deep copy operations", func() {
				experiment := gmeasure.NewExperiment("DeepCopy Operations")
				AddReportEntry(experiment.Name, experiment)

				experiment.SampleDuration("simple-struct", func(_ int) {
					src := Inner{Value: "benchmark"}
					var dst Inner
					deepcopy.Copy(&dst, &src)
				}, gmeasure.SamplingConfig{N: 10000})

				experiment.SampleDuration("nested-struct", func(_ int) {
					src := Nested{
						Name:   "root",
						Inner:  &Inner{Value: "nested"},
						Tags:   []string{"a", "b", "c"},
						Scores: map[string]int{"x": 1, "y": 2},
					}
					var dst Nested
					deepcopy.Copy(&dst, &src)
				}, gmeasure.SamplingConfig{N: 10000})

				AddReportEntry("Simple Struct Stats", experiment.GetStats("simple-struct"))
				AddReportEntry("Nested Struct Stats", experiment.GetStats("nested-struct"))

				Expect(experiment.GetStats("simple-struct").DurationFor(gmeasure.StatMedian)).To(
					BeNumerically("<", 10*time.Microsecond), "Median simple deep copy should be fast")
			})
		})

		// NOT SUPPORTED:
		// - Unexported field copying (limited support)
		// - Custom deep copy methods
		// - Channels and function copying
		// - Circular reference handling (will cause infinite recursion)
	})

	Context("Shallow vs Deep Copy Comparison", func() {
		It("should demonstrate difference between Copier (shallow) and go-deepcopy (deep)", func() {
			type Data struct {
				Name   string
				Values []int
			}

			original := Data{Name: "original", Values: []int{1, 2, 3}}

			By("Shallow copy with Copier (default) - shares slice reference")
			var shallowDst Data
			err := copier.Copy(&shallowDst, &original)
			Expect(err).ToNot(HaveOccurred())
			shallowDst.Values[0] = 999
			// Shallow copy: modifying dst affects original's slice
			Expect(original.Values[0]).To(Equal(999))

			By("Resetting original for deep copy test")
			original.Values[0] = 1

			By("Deep copy with go-deepcopy - independent slice")
			var deepDst Data
			err = deepcopy.Copy(&deepDst, &original)
			Expect(err).ToNot(HaveOccurred())
			deepDst.Values[0] = 888
			// Deep copy: modifying dst does NOT affect original
			Expect(original.Values[0]).To(Equal(1))
		})
	})
})
