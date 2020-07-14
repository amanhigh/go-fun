package gotest

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"time"
)

var _ = Describe("Json Encode/Decode", func() {
	var (
		name           = "Zoye"
		age            = 44
		number         = int64(88983333)
		originalPerson person
		personJson     string
	)
	BeforeEach(func() {
		originalPerson = person{name, age, number}
		personJson = fmt.Sprintf(`{"name":"%s","Age":%d,"MobileNumber":%d}`, name, age, number)
	})
	Context("Success", func() {
		It("should encode Properly", func() {
			jsonString, err := encodePerson(originalPerson)
			Expect(err).To(BeNil())
			Expect(jsonString).To(Equal(personJson))
		})

		It("should decode Properly", func() {
			decodedPerson, err := decodePerson(personJson)
			Expect(err).To(BeNil())
			Expect(decodedPerson).To(Equal(originalPerson))
		})
	})

	Context("Fail", func() {
		BeforeEach(func() {
			originalPerson.Name = "Bob"
		})

		It("should throw error on invalid json", func() {
			_, err := decodePerson("abcd")
			Expect(err).To(Not(BeNil()))
		})

		It("should not match original person", func() {
			jsonString, err := encodePerson(originalPerson)
			Expect(err).To(BeNil())
			Expect(jsonString).To(Not(Equal(personJson)))
		})
	})

	Context("Encode", func() {
		var (
			jsonString string
			err        error
		)

		AfterEach(func() {
			Expect(err).To(BeNil())

		})

		JustAfterEach(func() {
			//Creation
			jsonString, err = encodePerson(originalPerson)
		})

		Context("Success", func() {
			AfterEach(func() {
				Expect(jsonString).To(Equal(personJson))
			})

			It("should encode Properly", func() {
			})
		})

		Context("Fail", func() {
			//Assertions
			AfterEach(func() {
				Expect(jsonString).To(Not(Equal(personJson)))
			})

			//Configuration
			It("with changed age", func() {
				originalPerson.Age = 88
			})

			It("with changed name", func() {
				originalPerson.Name = "Bob"
			})
		})
	})

	Measure("it should do something hard efficiently", func(b Benchmarker) {
		runtime := b.Time("Encode", func() {
			output, _ := encodePerson(originalPerson)
			Expect(output).To(Equal(personJson))
		})

		Ω(runtime.Seconds()).Should(BeNumerically("<", 0.2), "SomethingHard() shouldn't take too long.")

		b.RecordValue("disk usage (in MB)", 1)
	}, 10)

	Context("Interesting Assertions", func() {
		var (
			err error
		)

		It("should match", func() {
			//Symbol Equivalent to Expect
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("should async", func() {
			Eventually(func() []int {
				time.Sleep(time.Millisecond * 10)
				return []int{2, 3}
			}, time.Millisecond*100, time.Millisecond*2).Should(HaveLen(2))
		})

		It("should deep equal", func() {
			pizza := "Cheeseboard Pizza"
			type FoodSrce string

			Ω(FoodSrce(pizza)).ShouldNot(Equal(pizza))       //will fail
			Ω(FoodSrce(pizza)).Should(BeEquivalentTo(pizza)) //will pass
		})

		It("should match collection", func() {
			//Array
			Ω([]string{"Foo", "FooBar"}).Should(ConsistOf("FooBar", "Foo"))
			Ω([]string{"Foo", "FooBar"}).Should(ConsistOf(ContainSubstring("Bar"), "Foo"))
			Ω([]string{"Foo", "FooBar"}).Should(ConsistOf(ContainSubstring("Foo"), ContainSubstring("Foo")))
			Ω([]string{"Foo", "FooBar"}).Should(ConsistOf([]string{"FooBar", "Foo"}))

			//Map
			Ω(map[string]string{"Foo": "Bar", "BazFoo": "Duck"}).Should(HaveKey(MatchRegexp(`.+Foo$`)))
			Ω(map[string]int{"Foo": 3, "BazFoo": 4}).Should(HaveKeyWithValue(MatchRegexp(`.+Foo$`), BeNumerically(">", 3)))
		})

		It("should panic", func() {
			Ω(func() { panic("FooBarBaz") }).Should(Panic())
		})

		Context("Channel", func() {
			var (
				c chan string
			)

			BeforeEach(func() {
				c = make(chan string, 0)
			})

			Context("With Close", func() {
				AfterEach(func() {
					Eventually(c).Should(BeClosed())
				})

				It("should receive", func() {
					go DoSomething(c, true)
				})

				It("should receive", func() {
					go DoSomething(c, true)
					Eventually(c).Should(Receive(Equal("Done!")))
				})
			})

			Context("No Close", func() {
				AfterEach(func() {
					Eventually(c).ShouldNot(BeClosed())
				})

				It("Channel Check Content", func() {
					go DoSomething(c, false)
					Expect(<-c).To(ContainSubstring("Done!"))
				})

				It("should not receive", func() {
					Consistently(c, time.Millisecond*100).ShouldNot(Receive())
				})
			})
		})
	})
})
