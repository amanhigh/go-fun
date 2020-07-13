package json

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestJsonEncode(t *testing.T) {
	var (
		name       = "Zoye"
		age        = 44
		number     = int64(88983333)
		personJson = fmt.Sprintf(`{"name":"%s","Age":%d,"MobileNumber":%d}`, name, age, number)
		person     = person{name, age, number}
	)

	Convey("Json", t, func() {
		Convey("encode", func() {
			age = 55
			So(encodePerson(person), ShouldEqual, personJson)
		})

		Convey("decode", func() {
			So(decodePerson(personJson), ShouldResemble, person)
		})

	})

	//Assertions - https://github.com/smartystreets/goconvey/blob/master/examples/assertion_examples_test.go
	Convey("Interesting Assertions", t, func() {
		So(1, ShouldAlmostEqual, 1.000000000000001)
		So([]int{1, 2, 3}, ShouldContain, 2)
		So(map[int]int{1: 1, 2: 2, 3: 3}, ShouldContainKey, 2)
		So(1, ShouldBeIn, []int{1, 2, 3})
		So(func() {}, ShouldNotPanic)
		So(1, ShouldNotHaveSameTypeAs, "1")
		So(bytes.NewBufferString(""), ShouldImplement, (*io.Reader)(nil))
	})
}
