package json

import (
	"encoding/json"
	"fmt"
)

type person struct {
	Name         string `json:"name"`
	Age          int
	MobileNumber int64
}

func main() {
	p1 := person{"Aman", 29, 9844415553}
	//fmt.Println(p1)
	encodedPerson := encodePerson(p1)
	fmt.Println("Encoded:", string(encodedPerson))

	decodedPerson := decodePerson(encodedPerson)
	fmt.Printf("Decoded: %+v\n", decodedPerson)
}
func decodePerson(encodedPerson string) person {
	var pDecoded person
	json.Unmarshal([]byte(encodedPerson), &pDecoded)
	return pDecoded
}

func encodePerson(p1 person) string {
	if decodedPerson, e := json.Marshal(p1); e == nil {
		return string(decodedPerson)
	} else {
		fmt.Println("Error:", e)
		return ""
	}
}
