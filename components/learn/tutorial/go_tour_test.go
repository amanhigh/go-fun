package tutorial_test

import (
	"fmt"
	"math"
	"math/rand"
	"regexp"
	"strconv"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var global, second_global = 5, 10

var _ = FDescribe("GoTour", func() {

	Context("Variables", func() {
		// var (
		// 	r, i rune = 8, 9
		// 	g         = 0.867 + 0.5i // complex128
		// )

		It("should have globals", func() {
			Expect(global).To(Not(BeNil()))
			Expect(second_global).To(Not(BeNil()))
		})

		It("should have locals", func() {
			var local string = "localvariable"
			shortHand := "Shorthand Variable"

			Expect(local).To(Not(BeNil()))
			Expect(shortHand).To(Not(BeNil()))
		})

		It("should have constants", func() {
			const constant_string = "Constant"
			Expect(constant_string).To(Not(BeNil()))

		})

		It("should have exported names", func() {
			Expect(math.Pi).To(Not(BeNil()))
		})

		It("should have enums", func() {
			const (
				MONDAY = 1 + iota
				TUESDAY
				WEDNESDAY
				THURSDAY
				FRIDAY
				SATURDAY
				SUNDAY
			)

			Expect(MONDAY).To(Equal(1))
			Expect(THURSDAY).To(Equal(4))
		})

		Context("Type Check", func() {
			var (
				string_var     = "hello"
				genric_var any = string_var
			)

			It("should cast valid string", func() {
				/** Empty Interface */
				casted_var, ok := genric_var.(string) //Type Casting
				Expect(ok).To(BeTrue())
				Expect(casted_var).To(Equal(string_var))
			})

			It("should not cast invalid float", func() {
				_, ok := genric_var.(float64) // Test Statement
				Expect(ok).To(BeFalse())
			})

			It("should fail string convert", func() {
				i, err := strconv.Atoi("XX")
				Expect(err).To(Not(BeNil()))
				Expect(i).To(Equal(0))
			})
		})

		Context("Struct", func() {
			var (
				vertex        Vertex
				pointerVertex *Vertex
			)

			BeforeEach(func() {
				vertex = Vertex{1, 2}
				pointerVertex = &vertex
			})

			It("should only mutate copy", func() {
				vertexCopy := vertex
				vertexCopy.Y = 7
				Expect(vertexCopy.Y).To(Equal(float64(7)))
				Expect(vertex.Y).To(Equal(float64(2)))
			})

			It("should mutate original on pointer", func() {
				pointerVertex.Y = 9
				Expect(vertex.Y).To(Equal(float64(9)))
			})

			Context("Interface", func() {

				var (
					a Abser
				)

				BeforeEach(func() {
					/** Interface */
					// a=vertex /** Gives Error as Abs takes only Pointer */
					a = pointerVertex
				})

				It("should compute absolute", func() {
					// /* While methods with pointer receivers take either a value or a pointer as the receiver when they are called: */
					Expect(vertex.Abs()).To(Equal(pointerVertex.Abs()))
					Expect(vertex.Abs()).To(Equal(a.Abs()))

					/** Null Handling pver.Abs() Would still work but when Abs will try to access X,Y Null Pointer would come. */
					// pver = nil
					// pver.Abs()
					/** Null Handling on Type would be error even if error is called on concrete type */
					// vertex = nil
					// vertex.Abs()
				})
			})
		})

		// fmt.Printf("Variables Type: %T Value: %v\n", r, r)
		// fmt.Printf("Variables Type: %T Value: %v\n", i, i)
		// fmt.Printf("Variables Type: %T Value: %v\n", g, g)

	})

	Context("Math", func() {
		const (
			// Create a huge number by shifting a 1 bit left 100 places.
			// In other words, the binary number that is 1 followed by 100 zeroes.
			Big = 1 << 100
			// Shift it right again 99 places, so we end up with 1<<1, or 2.
			Small = Big >> 99
		)

		It("should work", func() {

			fmt.Println("An untyped constant takes the type needed by its context")
			//Small is 2 and Big is 1^100
			Expect(needInt(Small)).To(Equal(21))
			Expect(needFloat(Small)).To(Equal(float64(0.2)))
			Expect(needFloat(Big)).To(BeNumerically(">", (float64(1.26765))))
		})

		It("should compute sqrt", func() {
			input := 8
			result, err := sqrt(input)
			Expect(err).To(BeNil())
			Expect(result).To(Equal(math.Sqrt((float64(input)))))
		})

		It("should not compute negative sqrt", func() {
			input := -2
			_, err := sqrt(input)
			Expect(err).To(Not(BeNil()))
		})

		It("should generate random", func() {
			rand.Seed(time.Now().UnixNano())
			Expect(rand.Intn(10)).To(Not(BeNil()))
		})

	})

	Context("Regex", func() {
		var (
			mysqlString  = "aman:aman@tcp(mysql:3306)/compute?charset=utf8&parseTime=True&loc=Local"
			mysqlMatcher = regexp.MustCompile(`^(.*)\((.*)\)(.*)$`)
		)

		It("should match", func() {
			matched := mysqlMatcher.FindAllStringSubmatch(mysqlString, 5)
			Expect(matched[0]).To(HaveLen(4))
			Expect(matched[0][2]).To(Equal("mysql:3306"))
			Expect(mysqlMatcher.ReplaceAllString(mysqlString, `$1#$2`)).To(Equal("aman:aman@tcp#mysql:3306"))
		})

	})
})

/* Structs */
type Abser interface {
	Abs() float64
}

type Vertex struct {
	X, Y float64
}

func (v *Vertex) Abs() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

/* Math */
func needInt(x int) int {
	return x*10 + 1
}

func needFloat(x float64) float64 {
	return x * 0.1
}

func sqrt(x int) (float64, error) {
	if x < 0 {
		return math.NaN(), ErrNegativeSqrt(x)
	}
	fX := float64(x)
	z := float64(1)
	z = 1.0
	for i := 0; i < 10; i++ {
		z = z - ((math.Pow(z, 2) - fX) / (2 * z))
	}
	return z, nil
}

/* Error Handling */
type ErrNegativeSqrt float64

func (e ErrNegativeSqrt) Error() string {
	return fmt.Sprintf("cannot Sqrt negativ number: %g", float64(e))
}
