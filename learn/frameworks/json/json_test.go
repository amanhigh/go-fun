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
		name       = "Zoye"
		age        = 44
		number     = int64(88983333)
		personJson = fmt.Sprintf(`{"name":"%s","Age":%d,"MobileNumber":%d}`, name, age, number)
	)
	var per person
	BeforeEach(func() {
		per = person{name, age, number}
	})

	It("should encode Properly", func() {
		Expect(encodePerson(per)).To(Equal(personJson))
	})

	It("should decode Properly", func() {
		Expect(decodePerson(personJson)).To(Equal(per))
	})
})

func BenchmarkEncode(b *testing.B) {
	var per = person{"Zoye", 44, 8983333}

	for n := 0; n < b.N; n++ {
		encodePerson(per)
	}
}
