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
			So(encodePerson(person), ShouldEqual, personJson)
		})

		Convey("decode", func() {
			So(decodePerson(personJson), ShouldResemble, person)
		})

	})

}
