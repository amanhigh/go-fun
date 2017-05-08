package main

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestJson(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Books Suite")
}

var _ = Describe("Json Encode/Decode", func() {
	var per person
	personJson := `{"name":"Zoye","Age":44,"MobileNumber":8983333}`

	BeforeEach(func() {
		per = person{"Zoye", 44, 8983333}
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
