package play_fast_test

import (
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"github.com/thoas/go-funk"
	"github.com/wesovilabs/koazee"
	"github.com/wesovilabs/koazee/stream"
)

type Person struct {
	Name string
	Male bool
	Age  int
}

var _ = FDescribe("Streams", func() {
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

	Context("Lodash", func() {
		It("should filter", func() {
			females := lo.Filter(people, func(person Person, index int) bool {
				return !person.Male
			})

			Expect(len(females)).To(Equal(4))
			lo.ForEach(females, func(person Person, index int) {
				Expect(person.Male).To(BeFalse())
			})
		})

		It("should do unique", func() {
			uniqueNames := lo.Uniq(lo.Map(people, func(person Person, index int) string { return person.Name }))

			Expect(len(uniqueNames)).To(Equal(6))
			Expect(uniqueNames).To(ContainElements("John Smith", "Jane Doe"))
		})

		It("should sum age (reduce)", func() {
			totalAge := lo.Reduce(people, func(total int, person Person, index int) int {
				return total + person.Age
			}, 0)

			Expect(totalAge).To(Equal(167)) // Sum of ages
		})

		It("should group by gender", func() {
			groupedByGender := lo.GroupBy(people, func(person Person) string {
				if person.Male {
					return "male"
				}
				return "female"
			})

			Expect(len(groupedByGender["male"])).To(Equal(3))
			Expect(len(groupedByGender["female"])).To(Equal(4))
		})

		Context("Map Age", func() {
			var ages []int
			BeforeEach(func() {
				ages = lo.Map(people, func(person Person, index int) int { return person.Age })
			})

			It("should work", func() {
				Expect(len(ages)).To(Equal(7))
				Expect(ages).To(ContainElements(32, 17, 20, 35, 35, 13, 15))
			})

			It("should find Max", func() {
				maxAge := lo.Max(ages)
				Expect(maxAge).To(Equal(35))
			})

			It("should find Min", func() {
				minAge := lo.Min(ages)
				Expect(minAge).To(Equal(13))
			})

			It("should reverse", func() {
				reversedAges := lo.Reverse(ages)
				Expect(len(reversedAges)).To(Equal(7))
				Expect(reversedAges[0]).To(Equal(15))
				Expect(reversedAges[len(reversedAges)-1]).To(Equal(32))
			})
		})

		It("should shuffle", func() {
			shuffledPeople := lo.Shuffle(people)

			Expect(len(shuffledPeople)).To(Equal(7))
		})
	})

})
