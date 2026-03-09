package play_fast

import (
	"time"

	"github.com/jinzhu/copier"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gmeasure"
	deepcopy "github.com/tiendc/go-deepcopy"
)

// Types for Copy tests (must be at package level for method definitions)
type Coffee struct {
	Name          string
	Origin        string
	RoastLevel    string
	Price         float64
	Age           int32          // For method-to-field tests
	ProcessedRole string         // Target for field-to-method copying
	Location      Location       // For nested struct copying tests
	Flavors       []string       // For slice copying tests - coffee flavor notes
	Ratings       map[string]int // For map copying tests - coffee ratings by source
	TastingNotes  *string        // For pointer copying tests - detailed tasting notes
	Quality       CoffeeQuality  // For interface copying tests - quality assessment
}

// CoffeeQuality interface for quality assessment
type CoffeeQuality interface {
	GetScore() int
	GetDescription() string
}

// QualityRating implements CoffeeQuality interface
type QualityRating struct {
	Score       int
	Description string
}

func (qr QualityRating) GetScore() int {
	return qr.Score
}

func (qr QualityRating) GetDescription() string {
	return qr.Description
}

// Location nested struct for coffee origin information
type Location struct {
	City     string
	Country  string
	Region   string
	Altitude int // meters above sea level
}

// DoubleAge is a method that copier can use for method-to-field copying
func (c *Coffee) DoubleAge() int32 {
	return 2 * c.Age
}

// ProcessRole is a setter method that copier can use for field-to-method copying
func (c *Coffee) ProcessRole(role string) {
	c.ProcessedRole = "Super " + role
}

type CoffeeDTO struct {
	Name        string
	Origin      string `copier:"-"` // For field exclusion tests
	RoastLevel  string
	Price       float64
	Age         int32
	DoubleAge   int32          // Target for method-to-field copying
	ProcessRole string         // Source for field-to-method copying
	Location    Location       // For nested struct copying tests
	Flavors     []string       // For slice copying tests - coffee flavor notes
	Ratings     map[string]int // For map copying tests - coffee ratings by source
}

var _ = Describe("Copy", func() {

	Context("Copier", func() {

		It("should build", func() {
			src := Coffee{}
			dst := Coffee{}
			err := copier.Copy(&dst, &src)
			Expect(err).ToNot(HaveOccurred())
		})

		Context("Basic", func() {
			var (
				src Coffee
				dst CoffeeDTO
			)

			BeforeEach(func() {
				src = Coffee{
					Name:       "Ethiopian Yirgacheffe",
					Origin:     "Ethiopia",
					RoastLevel: "Light",
					Price:      18.99,
					Age:        0, // Initialize all fields
					Location: Location{
						City:     "Yirgacheffe",
						Country:  "Ethiopia",
						Region:   "Sidamo",
						Altitude: 1800,
					},
				}
				err := copier.Copy(&dst, &src)
				Expect(err).ToNot(HaveOccurred())
			})

			It("1.1 should copy simple struct fields", func() {
				Expect(dst.Name).To(Equal(src.Name))
				Expect(dst.Origin).To(BeEmpty()) // Ignored via copier:"-" tag
				Expect(dst.RoastLevel).To(Equal(src.RoastLevel))
				Expect(dst.Price).To(Equal(src.Price))
				Expect(dst.Age).To(Equal(src.Age))
			})

			It("1.2 should copy nested struct fields", func() {
				src := Coffee{
					Name:     "Alice",
					Age:      30,
					Location: Location{City: "Seattle", Country: "USA", Region: "Washington", Altitude: 50},
				}
				var dst CoffeeDTO
				err := copier.Copy(&dst, &src)
				Expect(err).ToNot(HaveOccurred())
				Expect(dst.Location.City).To(Equal("Seattle"))
				Expect(dst.Location.Country).To(Equal("USA"))
			})
		})

		Context("Medium", func() {
			Context("Field Mapping and Transformation", func() {
				It("2.1 should copy from getter methods to fields (method-to-field)", func() {
					By("Copying Coffee with DoubleAge method to CoffeeDTO with DoubleAge field")
					coffee := Coffee{Name: "John", Age: 25}
					var dto CoffeeDTO
					err := copier.Copy(&dto, &coffee)
					Expect(err).ToNot(HaveOccurred())

					Expect(dto.Name).To(Equal("John"))
					Expect(dto.Age).To(Equal(int32(25)))
					Expect(dto.DoubleAge).To(Equal(int32(50)))
				})

				It("2.2 should copy from fields to setter methods (field-to-method)", func() {
					By("Copying CoffeeDTO ProcessRole field to Coffee ProcessRole method")

					// Source: DTO with ProcessRole field
					source := CoffeeDTO{
						Name:        "John",
						Age:         25,
						ProcessRole: "Admin", // This field should trigger ProcessRole method
					}

					// Destination: Coffee with ProcessRole method and ProcessedRole field
					var destination Coffee
					err := copier.Copy(&destination, &source)
					Expect(err).ToNot(HaveOccurred())

					// Verify field-to-method copying worked
					Expect(destination.ProcessedRole).To(Equal("Super Admin")) // ProcessRole method added "Super " prefix
					Expect(destination.Name).To(Equal("John"))
					Expect(destination.Age).To(Equal(int32(25)))
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
					base := Coffee{Name: "House Blend", Origin: "Colombia", RoastLevel: "Medium", Price: 12.99}
					override := Coffee{Name: "House Blend Special", Price: 14.99}

					var result Coffee
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
					base := Coffee{Name: "Base", Origin: "Kenya", RoastLevel: "Dark", Price: 15.0}
					override := Coffee{Name: "Override", Price: 20.0}

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
					src := Coffee{Name: "Espresso", Origin: "Italy", RoastLevel: "Dark", Price: 22.0, Age: 30}
					var dst CoffeeDTO
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
					var src *Coffee
					dst := &Coffee{}
					err := copier.Copy(dst, src)
					Expect(err).To(HaveOccurred())
				})

				It("1.4 should handle empty structs", func() {
					src := Coffee{}
					var dst Coffee
					err := copier.Copy(&dst, &src)
					Expect(err).ToNot(HaveOccurred())
					Expect(dst.Name).To(BeEmpty())
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

			It("2.7 should handle type conversion failures gracefully", func() {
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

		It("should build", func() {
			src := Coffee{Name: "test"}
			var dst Coffee
			err := deepcopy.Copy(&dst, &src)
			Expect(err).ToNot(HaveOccurred())
			Expect(dst.Name).To(Equal("test"))
		})

		Context("Deep Copy Basics", func() {
			It("should clone simple structs without shared references", func() {
				// Using Coffee struct with realistic coffee data
				src := Coffee{Name: "Ethiopian Yirgacheffe", Price: 28.99}
				var dst Coffee
				err := deepcopy.Copy(&dst, &src)
				Expect(err).ToNot(HaveOccurred())

				By("Verifying no shared references")
				src.Name = "modified"
				Expect(dst.Name).To(Equal("Ethiopian Yirgacheffe"))
			})
		})

		Context("Nested Structure Cloning", func() {
			var (
				src Coffee
				dst Coffee
			)

			BeforeEach(func() {
				notes := "Complex floral aroma with bright citrus acidity and smooth finish"
				src = Coffee{
					Name:         "Ethiopian Yirgacheffe",
					Origin:       "Ethiopia",
					RoastLevel:   "Light",
					Price:        28.99,
					Age:          1,
					TastingNotes: &notes,
					Location: Location{
						City:     "Yirgacheffe",
						Country:  "Ethiopia",
						Region:   "Sidamo",
						Altitude: 1800,
					},
					Flavors: []string{"floral", "citrus", "bergamot", "tea-like"},
					Ratings: map[string]int{
						"aroma":      10,
						"acidity":    9,
						"body":       7,
						"flavor":     10,
						"aftertaste": 9,
						"sweetness":  8,
					},
				}
				err := deepcopy.Copy(&dst, &src)
				Expect(err).ToNot(HaveOccurred())
			})

			It("should deep copy all fields", func() {
				Expect(dst.Name).To(Equal("Ethiopian Yirgacheffe"))
				Expect(dst.Location.City).To(Equal("Yirgacheffe"))
				Expect(dst.TastingNotes).ToNot(BeNil())
				Expect(*dst.TastingNotes).To(Equal("Complex floral aroma with bright citrus acidity and smooth finish"))
				Expect(dst.Flavors).To(HaveLen(4))
				Expect(dst.Flavors).To(ContainElements("floral", "citrus", "bergamot", "tea-like"))
				Expect(dst.Ratings).To(HaveKeyWithValue("aroma", 10))
			})

			It("should have no shared slice references", func() {
				src.Flavors[0] = "modified"
				Expect(dst.Flavors[0]).To(Equal("floral"))
			})

			It("should have no shared map references", func() {
				src.Ratings["aroma"] = 0
				Expect(dst.Ratings["aroma"]).To(Equal(10))
			})
		})

		Context("Pointer Handling", func() {
			It("3.2 should clone structs with pointer fields independently", func() {
				notes := "Complex floral aroma with bright citrus acidity and smooth finish"
				src := Coffee{
					Name:         "Ethiopian Yirgacheffe",
					Price:        28.99,
					TastingNotes: &notes,
				}
				var dst Coffee
				err := deepcopy.Copy(&dst, &src)
				Expect(err).ToNot(HaveOccurred())

				By("Verifying values are copied")
				Expect(dst.Name).To(Equal("Ethiopian Yirgacheffe"))
				Expect(dst.Price).To(Equal(28.99))
				Expect(dst.TastingNotes).ToNot(BeNil())
				Expect(*dst.TastingNotes).To(Equal("Complex floral aroma with bright citrus acidity and smooth finish"))

				By("Verifying pointed-to data is independent")
				*src.TastingNotes = "Modified tasting notes"
				Expect(*dst.TastingNotes).To(Equal("Complex floral aroma with bright citrus acidity and smooth finish"))
			})

			It("1.3 should handle nil pointers", func() {
				src := Coffee{Name: "Colombian Supremo", TastingNotes: nil}
				var dst Coffee
				err := deepcopy.Copy(&dst, &src)
				Expect(err).ToNot(HaveOccurred())
				Expect(dst.Name).To(Equal("Colombian Supremo"))
				Expect(dst.TastingNotes).To(BeNil())
			})
		})

		Context("Slice and Array Cloning", func() {
			It("3.3 should deep copy coffee flavors", func() {
				src := []string{"floral", "citrus", "chocolate", "caramel"}
				var dst []string
				err := deepcopy.Copy(&dst, &src)
				Expect(err).ToNot(HaveOccurred())

				Expect(dst).To(Equal([]string{"floral", "citrus", "chocolate", "caramel"}))

				By("Verifying modifications to clone don't affect original")
				dst[0] = "modified"
				Expect(src[0]).To(Equal("floral"))
			})
		})

		Context("Map Cloning", func() {
			It("3.4 should deep copy coffee ratings map", func() {
				src := map[string]int{
					"aroma":      9,
					"acidity":    8,
					"body":       7,
					"flavor":     9,
					"aftertaste": 8,
				}
				var dst map[string]int
				err := deepcopy.Copy(&dst, &src)
				Expect(err).ToNot(HaveOccurred())

				Expect(dst).To(HaveKeyWithValue("aroma", 9))

				By("Verifying clone has independent map instance")
				src["aroma"] = 0
				Expect(dst["aroma"]).To(Equal(9))
			})
		})

		Context("Interface Fields", func() {
			It("3.6 should deep copy structs with interface fields", func() {
				quality := QualityRating{
					Score:       9,
					Description: "Exceptional quality with complex flavor profile",
				}
				src := Coffee{
					Name:    "Ethiopian Yirgacheffe",
					Price:   28.99,
					Quality: quality,
				}
				var dst Coffee
				err := deepcopy.Copy(&dst, &src)
				Expect(err).ToNot(HaveOccurred())

				Expect(dst.Name).To(Equal("Ethiopian Yirgacheffe"))
				Expect(dst.Price).To(Equal(28.99))
				Expect(dst.Quality).ToNot(BeNil())
				Expect(dst.Quality.GetScore()).To(Equal(9))
				Expect(dst.Quality.GetDescription()).To(Equal("Exceptional quality with complex flavor profile"))

				By("Verifying interface field independence")
				// Modify the original quality rating
				originalQuality, ok := src.Quality.(QualityRating)
				Expect(ok).To(BeTrue())
				originalQuality.Score = 5
				Expect(dst.Quality.GetScore()).To(Equal(9)) // Should remain unchanged
			})
		})

		Context("Performance Benchmarks", FlakeAttempts(3), func() {
			It("3.8 should benchmark deep copy operations", func() {
				experiment := gmeasure.NewExperiment("DeepCopy Operations")
				AddReportEntry(experiment.Name, experiment)

				experiment.SampleDuration("simple-struct", func(_ int) {
					src := Coffee{Name: "benchmark"}
					var dst Coffee
					deepcopy.Copy(&dst, &src)
				}, gmeasure.SamplingConfig{N: 10000})

				experiment.SampleDuration("complex-struct", func(_ int) {
					src := Coffee{
						Name:       "Ethiopian Yirgacheffe",
						Origin:     "Ethiopia",
						RoastLevel: "Light",
						Price:      28.99,
						Age:        1,
						Location: Location{
							City:     "Yirgacheffe",
							Country:  "Ethiopia",
							Region:   "Sidamo",
							Altitude: 1800,
						},
						Flavors: []string{"floral", "citrus", "bergamot", "tea-like"},
						Ratings: map[string]int{
							"aroma":      10,
							"acidity":    9,
							"body":       7,
							"flavor":     10,
							"aftertaste": 9,
							"sweetness":  8,
						},
					}
					var dst Coffee
					deepcopy.Copy(&dst, &src)
				}, gmeasure.SamplingConfig{N: 10000})

				AddReportEntry("Simple Struct Stats", experiment.GetStats("simple-struct"))
				AddReportEntry("Complex Struct Stats", experiment.GetStats("complex-struct"))

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
		It("6.1 should demonstrate difference between Copier (shallow) and go-deepcopy (deep)", func() {
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
