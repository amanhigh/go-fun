package gotest

import (
	"encoding/json"
)

type person struct {
	Name         string `json:"name"`
	Age          int
	MobileNumber int64
}

func decodePerson(encodedPerson string) (p person, err error) {
	err = json.Unmarshal([]byte(encodedPerson), &p)
	return
}

func encodePerson(p1 person) (jsonString string, err error) {
	var jsonBytes []byte
	jsonBytes, err = json.Marshal(p1)
	jsonString = string(jsonBytes)
	return
}

func DoSomething(c chan string, shouldClose bool) {
	c <- "Done!"
	if shouldClose {
		close(c)
	}
}
