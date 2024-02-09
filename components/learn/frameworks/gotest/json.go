// mockgen -package gotest -destination json_mock.go -source json.go
package gotest

//Install
//go install go.uber.org/mock/mockgen@latest
//brew install mockery, yay -S mockery-bin

import (
	"encoding/json"
)

//go:generate mockery --name PersonEncoder --inpackage --structname MockEncoder
type PersonEncoder interface {
	EncodePerson(p Person) (jsonString string, err error)
	DecodePerson(encodedPerson string) (p Person, err error)
}

type Person struct {
	Name         string `json:"name"`
	Age          int
	MobileNumber int64
}

func decodePerson(encodedPerson string) (p Person, err error) {
	err = json.Unmarshal([]byte(encodedPerson), &p)
	return
}

func encodePerson(p Person) (jsonString string, err error) {
	var jsonBytes []byte
	jsonBytes, err = json.Marshal(p)
	jsonString = string(jsonBytes)
	return
}

func DoSomething(c chan string, shouldClose bool) {
	c <- "Done!"
	if shouldClose {
		close(c)
	}
}
