package json

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
	"time"
)

func TestJson(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Books Suite")
}

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

	Context("Interesting Assertions", func() {
		Context("Channel", func() {
			var (
				c chan string
			)
			BeforeEach(func() {
				c = make(chan string, 0)

			})

			It("should receive", func() {
				go DoSomething(c, true)
				Eventually(c).Should(BeClosed())
			})

			It("should receive", func() {
				go DoSomething(c, true)
				Eventually(c).Should(Receive(Equal("Done!")))
				Eventually(c, time.Nanosecond).ShouldNot(BeClosed())
			})

			It("Channel Check Content", func() {
				go DoSomething(c, false)
				Expect(<-c).To(ContainSubstring("Done!"))
				Eventually(c).ShouldNot(BeClosed())
			})
		})
	})
})

func BenchmarkEncode(b *testing.B) {
	var per = person{"Zoye", 44, 8983333}

	for n := 0; n < b.N; n++ {
		encodePerson(per)
	}
}
