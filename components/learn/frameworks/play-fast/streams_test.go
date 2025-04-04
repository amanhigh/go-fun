package play_fast

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/samber/lo"
	"github.com/thoas/go-funk"
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

	Context("Funk", func() {

		It("should do unique", func() {
			uniqueNames := funk.Uniq(funk.Map(people, func(person Person) (name string) { return person.Name }))
			Expect(uniqueNames).To(HaveLen(6))
			Expect(uniqueNames).To(ContainElements("John Smith", "Jane Doe"))
		})

		It("should map", func() {
			personMap := funk.ToMap(people, "Name")
			ageList := funk.Map(personMap, func(_ string, person Person) (age int) { return person.Age })

			Expect(ageList).To(HaveLen(6))
			Expect(ageList).To(ContainElements(13, 15, 32, 17, 20, 35))
		})
	})

	Context("Lodash", func() {
		It("should filter", func() {
			females := lo.Filter(people, func(person Person, _ int) bool {
				return !person.Male
			})

			Expect(females).To(HaveLen(4))
			lo.ForEach(females, func(person Person, _ int) {
				Expect(person.Male).To(BeFalse())
			})
		})

		It("should do unique", func() {
			uniqueNames := lo.Uniq(lo.Map(people, func(person Person, _ int) string { return person.Name }))

			Expect(uniqueNames).To(HaveLen(6))
			Expect(uniqueNames).To(ContainElements("John Smith", "Jane Doe"))
		})

		It("should sum age (reduce)", func() {
			totalAge := lo.Reduce(people, func(total int, person Person, _ int) int {
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

			Expect(groupedByGender["male"]).To(HaveLen(3))

			Expect(groupedByGender["female"]).To(HaveLen(4))
		})
		Context("AgeMap (Associate)", func() {
			var ageMap map[int]Person
			BeforeEach(func() {
				ageMap = lo.Associate(people, func(person Person) (key int, value Person) {
					return person.Age, person
				})
			})

			It("should build AgeMap (Associate)", func() {
				Expect(ageMap).To(HaveLen(6)) // Unique ages
			})
			It("should extract keys (ages)", func() {
				ages := lo.Keys(ageMap)
				Expect(ages).To(HaveLen(6))
				Expect(ages[0]).ShouldNot(Equal(0))
			})

			It("should extract values (persons)", func() {
				persons := lo.Values(ageMap)
				Expect(persons).To(HaveLen(6))
				Expect(persons[0].Name).ShouldNot(BeEmpty())
			})

			It("should invert", func() {
				invertedMap := lo.Invert(ageMap)
				Expect(invertedMap).To(HaveLen(6))
			})

			It("should map to slice (Age_Name)", func() {
				ageNames := lo.MapToSlice(ageMap, func(age int, person Person) string { return fmt.Sprintf("%d_%s", age, person.Name) })
				Expect(ageNames).To(HaveLen(6))
				Expect(ageNames[0]).Should(ContainSubstring("_"))
			})

			It("should pick by keys (age > 30)", func() {
				// Define the keys you want to pick
				keys := []int{32, 35}

				olderPeople := lo.PickByKeys(ageMap, keys)
				Expect(olderPeople).To(HaveLen(2))
			})
		})

		It("should drop and drop right", func() {
			droppedPeople := lo.Drop(people, 2)
			Expect(droppedPeople).To(HaveLen(5))

			droppedRightPeople := lo.DropRight(people, 2)
			Expect(droppedRightPeople).To(HaveLen(5))
		})

		It("should reject (odd age)", func() {
			evenAgePeople := lo.Reject(people, func(person Person, _ int) bool {
				return person.Age%2 != 0
			})

			Expect(evenAgePeople).To(HaveLen(2)) // People with even ages
		})

		Context("Map Age", func() {
			var ages []int
			BeforeEach(func() {
				ages = lo.Map(people, func(person Person, _ int) int { return person.Age })
			})

			It("should work", func() {
				Expect(ages).To(HaveLen(7))
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
				Expect(reversedAges).To(HaveLen(7))
				Expect(reversedAges[0]).To(Equal(15))
				Expect(reversedAges[len(reversedAges)-1]).To(Equal(32))
			})

			It("should chunk", func() {
				chunks := lo.Chunk(ages, 3)
				Expect(chunks).To(HaveLen(3))
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

				Expect(agesBack).To(HaveLen(7))
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

			It("should check Contains", func() {
				Expect(lo.Contains(ages, 35)).To(BeTrue())
			})

			It("should check ContainsBy (Age 35)", func() {
				Expect(lo.ContainsBy(people, func(person Person) bool { return person.Age == 35 })).To(BeTrue())
			})

			It("should check Every", func() {
				Expect(lo.Every(ages, []int{32, 17})).To(BeTrue())
			})

			It("should check EveryBy (Age > 0)", func() {
				Expect(lo.EveryBy(ages, func(age int) bool { return age > 0 })).To(BeTrue())
			})

			It("should check Some", func() {
				Expect(lo.Some(ages, []int{13})).To(BeTrue())
			})

			It("should check SomeBy (Age between 20 and 30)", func() {
				Expect(lo.SomeBy(ages, func(age int) bool { return age >= 20 && age <= 30 })).To(BeTrue())
			})

			It("should check None (Age > 100)", func() {
				Expect(lo.None(ages, []int{101})).To(BeTrue())
			})

			It("should check IndexOf", func() {
				Expect(lo.IndexOf(ages, 35)).To(Equal(3)) // Assuming the ages are in the same order as defined
			})

			It("should check FindOrElse", func() {
				foundPerson := lo.FindOrElse(people, Person{Name: "Default Person"}, func(person Person) bool { return person.Name == "Nonexistent Person" })
				Expect(foundPerson.Name).To(Equal("Default Person"))
			})

		})

		Context("Conditions", func() {
			ages := []int{32, 17, 20, 35, 35, 13, 15}

			It("should check Without", func() {
				Expect(lo.Without(ages, 32, 17)).To(Equal([]int{20, 35, 35, 13, 15}))
			})

			It("should check WithoutEmpty", func() {
				Expect(lo.WithoutEmpty([]string{"John", "", "Jane"})).To(Equal([]string{"John", "Jane"}))
			})

			It("should check Nth", func() {
				Expect(lo.Nth(ages, 2)).To(Equal(20)) // Zero-indexed
			})

			It("should check Ternary", func() {
				Expect(lo.Ternary(true, "yes", "no")).To(Equal("yes"))
				Expect(lo.Ternary(false, "yes", "no")).To(Equal("no"))
			})

			It("should check lo.If and lo.ElseIf", func() {
				result := lo.If(len(ages) > 0, "PositiveAge").ElseIf(len(ages) == 0, "ZeroAge").Else("unknown")
				Expect(result).To(Equal("PositiveAge"))
			})

			It("should check lo.IsNil", func() {
				Expect(lo.IsNil(nil)).To(BeTrue())
				Expect(lo.IsNil(ages)).To(BeFalse())
			})

			It("should check lo.Empty", func() {
				Expect(lo.Empty[int]()).To(Equal(0))
				Expect(lo.Empty[string]()).To(Equal(""))
			})

			It("should validate", func() {
				slice := []string{"a"}
				err := lo.Validate(len(slice) == 0, "Slice should be empty but contains %v", slice)
				Expect(err).Should(HaveOccurred())
			})
		})

		It("should Find Person", func() {
			foundPerson, ok := lo.Find(people, func(person Person) bool { return person.Name == "Jane Doe" })
			Expect(ok).To(BeTrue())
			Expect(foundPerson.Name).To(Equal("Jane Doe"))
		})

		It("should shuffle", func() {
			shuffledPeople := lo.Shuffle(people)

			Expect(shuffledPeople).To(HaveLen(7))
		})
	})

})
