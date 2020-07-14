package json

import (
	"fmt"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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

	PContext("Fail", func() {
		It("should throw error on invalid json", func() {
			_, err := decodePerson("abcd")
			Expect(err).To(Not(BeNil()))
		})
	})

})

func BenchmarkEncode(b *testing.B) {
	var per = person{"Zoye", 44, 8983333}

	for n := 0; n < b.N; n++ {
		encodePerson(per)
	}
}
