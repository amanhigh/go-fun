package play_fast_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
	"github.com/thoas/go-funk"
	"github.com/wesovilabs/koazee"
	"github.com/wesovilabs/koazee/stream"
	"strings"
)

type Person struct {
	Name string
	Male bool
	Age  int
}

var _ = Describe("Streams", func() {
	var (
		people []Person
	)
	BeforeEach(func() {
		people = []Person{
			{"John Smith", true, 32},
			{"Peter Pan", true, 17},
			{"Jane Doe", false, 20},
			{"Anna Wallace", false, 35},
			{"Anna Wallace", false, 35},
			{"Tim O'Brian", true, 13},
			{"Celia Hills", false, 15},
		}
	})

	Context("Koazee", func() {
		/* Map Functions not supported yet, but has multiple operations */
		var (
			peopleStream stream.Stream
		)

		BeforeEach(func() {
			peopleStream = koazee.StreamOf(people)
		})

		It("should filter and sort", func() {
			stream := peopleStream.Filter(func(person Person) bool {
				return !person.Male
			}).Sort(func(person, otherPerson Person) int {
				return strings.Compare(person.Name, otherPerson.Name)
			}).ForEach(func(person Person) {
				logrus.Debugf("%s is %d years old", person.Name, person.Age)
				Expect(person.Male).To(BeFalse())
			})

			logrus.Debug("Operations are not evaluated until we perform stream.Do, Count etc")
			//stream.Do()
			Expect(stream.Count()).To(Equal(4))

		})

		It("should do unique", func() {
			uniqueNames := peopleStream.Map(func(person Person) (name string) { return person.Name }).RemoveDuplicates().Out().Val()
			Expect(uniqueNames).To(HaveLen(6))
			Expect(uniqueNames).To(ContainElements("John Smith", "Jane Doe"))
		})
	})

	Context("Funk", func() {

		It("should do unique", func() {
			uniqueNames := funk.Uniq(funk.Map(people, func(person Person) (name string) { return person.Name }))
			Expect(uniqueNames).To(HaveLen(6))
			Expect(uniqueNames).To(ContainElements("John Smith", "Jane Doe"))
		})

		It("should map", func() {
			personMap := funk.ToMap(people, "Name")
			ageList := funk.Map(personMap, func(name string, person Person) (age int) { return person.Age })

			Expect(ageList).To(HaveLen(6))
			Expect(ageList).To(ContainElements(13, 15, 32, 17, 20, 35))
		})

	})

})
