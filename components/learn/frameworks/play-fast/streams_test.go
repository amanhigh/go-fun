package play_fast_test

import (
	"fmt"
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

		Context("AgeMap (Associate)", func() {
			var ageMap map[int]Person
			BeforeEach(func() {
				ageMap = lo.Associate(people, func(person Person) (key int, value Person) {
					return person.Age, person
				})
			})

			It("should build AgeMap (Associate)", func() {
				Expect(len(ageMap)).To(Equal(6)) // Unique ages
			})

			It("should extract keys (ages)", func() {
				ages := lo.Keys(ageMap)
				Expect(len(ages)).To(Equal(6))
				Expect(ages[0]).ShouldNot(Equal(0))
			})

			It("should extract values (persons)", func() {
				persons := lo.Values(ageMap)
				Expect(len(persons)).To(Equal(6))
				Expect(persons[0].Name).ShouldNot(BeEmpty())
			})

			It("should invert", func() {
				invertedMap := lo.Invert(ageMap)
				Expect(len(invertedMap)).To(Equal(6))
			})

			It("should map to slice (Age_Name)", func() {
				ageNames := lo.MapToSlice(ageMap, func(age int, person Person) string { return fmt.Sprintf("%d_%s", age, person.Name) })
				Expect(len(ageNames)).To(Equal(6))
				Expect(ageNames[0]).Should(ContainSubstring("_"))
			})

			It("should pick by keys (age > 30)", func() {
				// Define the keys you want to pick
				keys := []int{32, 35}

				olderPeople := lo.PickByKeys(ageMap, keys)
				Expect(len(olderPeople)).To(Equal(2))
			})
		})

		It("should drop and drop right", func() {
			droppedPeople := lo.Drop(people, 2)
			Expect(len(droppedPeople)).To(Equal(5))

			droppedRightPeople := lo.DropRight(people, 2)
			Expect(len(droppedRightPeople)).To(Equal(5))
		})

		It("should reject (odd age)", func() {
			evenAgePeople := lo.Reject(people, func(person Person, index int) bool {
				return person.Age%2 != 0
			})

			Expect(len(evenAgePeople)).To(Equal(2)) // People with even ages
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

			It("should chunk", func() {
				chunks := lo.Chunk(ages, 3)
				Expect(len(chunks)).To(Equal(3))
			})

			It("should count by (age < 30)", func() {
				count := lo.CountBy(ages, func(age int) bool { return age < 30 })
				Expect(count).To(Equal(4)) // Number of people with age < 30
			})

			It("should sum", func() {
				sum := lo.Sum(ages)
				Expect(sum).To(Equal(167)) // Sum of ages
			})

			It("should convert slice to channel and back", func() {
				ageChan := lo.SliceToChannel(2, ages)
				agesBack := lo.ChannelToSlice(ageChan)

				Expect(len(agesBack)).To(Equal(7))
				Expect(agesBack).To(ContainElements(32, 17, 20, 35, 35, 13, 15))
			})

			It("should buffer (read from channel)", func() {
				ageChan := lo.SliceToChannel(1, ages)
				bufferedAges, length, _, ok := lo.Buffer(ageChan, 3)

				Expect(ok).To(BeTrue())
				Expect(length).To(Equal(3))
				for _, age := range bufferedAges {
					Expect(age).To(BeElementOf(ages))
				}
			})

			It("should contain", func() {
				Expect(lo.Contains(ages, 32)).To(BeTrue())
				Expect(lo.Contains(ages, 100)).To(BeFalse())
			})
		})

		It("should shuffle", func() {
			shuffledPeople := lo.Shuffle(people)

			Expect(len(shuffledPeople)).To(Equal(7))
		})
	})

})
