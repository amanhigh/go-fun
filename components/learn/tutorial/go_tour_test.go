package tutorial_test

import (
	"strconv"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var global, second_global = 5, 10

type Vertex struct {
	X, Y float64
}

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
				vertex = Vertex{1, 2}
			)

			It("should only mutate copy", func() {
				vertexCopy := vertex
				vertexCopy.Y = 7
				Expect(vertexCopy.Y).To(Equal(float64(7)))
				Expect(vertex.Y).To(Equal(float64(2)))
			})

			It("should mutate original on pointer", func() {
				pointerVertex := &vertex
				pointerVertex.Y = 9

				Expect(vertex.Y).To(Equal(float64(9)))
			})

			/** Interface */
			// var a Abser
			// //a=vertex /** Gives Error as Abs takes only Pointer */
			// a = pointerVertex

			// /* While methods with pointer receivers take either a value or a pointer as the receiver when they are called: */
			// fmt.Println("Vertex Method: ", vertex.Abs(), pointerVertex.Abs(), a.Abs()) //Methods linked to Struct

			/** Null Handling pver.Abs() Would still work but when Abs will try to access X,Y Null Pointer would come. */
			//pver=nil;pver.Abs();
			/** Null Handling on Type would be error even if error is called on concrete type */
			//vertex=nil;vertex.Abs()

		})

		// fmt.Printf("Variables Type: %T Value: %v\n", r, r)
		// fmt.Printf("Variables Type: %T Value: %v\n", i, i)
		// fmt.Printf("Variables Type: %T Value: %v\n", g, g)

	})

})
