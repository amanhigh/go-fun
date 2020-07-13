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
