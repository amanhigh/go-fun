package libraries

import (
	"fmt"
	"github.com/thoas/go-funk"
	"github.com/wesovilabs/koazee"
	"strings"
)

type Person struct {
	Name string
	Male bool
	Age  int
}

var people = []Person{
	{"John Smith", true, 32},
	{"Peter Pan", true, 17},
	{"Jane Doe", false, 20},
	{"Anna Wallace", false, 35},
	{"Anna Wallace", false, 35},
	{"Tim O'Brian", true, 13},
	{"Celia Hills", false, 15},
}

func StreamFun() {
	/* Map Functions not supported yet, but has multiple operations */
	fmt.Println("** Koazee **")
	peopleStream := koazee.StreamOf(people)
	stream := peopleStream.
		Filter(func(person Person) bool {
			return !person.Male
		}).
		Sort(func(person, otherPerson Person) int {
			return strings.Compare(person.Name, otherPerson.Name)
		}).
		ForEach(func(person Person) {
			fmt.Printf("%s is %d years old\n", person.Name, person.Age)
		})

	fmt.Println("Operations are not evaluated until we perform stream.Do()\n")
	stream.Do()

	fmt.Println("Uniq Names", peopleStream.Map(func(person Person) (name string) { return person.Name }).RemoveDuplicates().Out().Val())

	fmt.Println("** FUNK **")
	fmt.Println("Uniq Names", funk.Uniq(funk.Map(people, func(person Person) (name string) { return person.Name })))

	personMap := funk.ToMap(people, "Name")
	fmt.Println("PeopleMap", personMap)
	fmt.Println("Age", funk.Map(personMap, func(name string, person Person) (age int) { return person.Age }))
}
