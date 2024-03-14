// mockgen -package gotest -destination json_mock.go -source json.go
package gotest

//Install
//go install go.uber.org/mock/mockgen@latest
//brew install mockery, yay -S mockery-bin

import (
	"encoding/json"
	"errors"
)

//go:generate mockery --name PersonEncoder --inpackage --structname MockEncoder
type PersonEncoder interface {
	EncodePerson(person Person) (jsonString string, err error)
	DecodePerson(encodedPerson string) (person Person, err error)
}

type Person struct {
	Name         string `json:"name"`
	Age          int
	MobileNumber int64
}

type PersonEncoderImpl struct{}

func (p *PersonEncoderImpl) DecodePerson(encodedPerson string) (person Person, err error) {
	err = json.Unmarshal([]byte(encodedPerson), &person)
	return
}

func (p *PersonEncoderImpl) EncodePerson(person Person) (jsonString string, err error) {
	var jsonBytes []byte

	if person.Age < 0 {
		err = errors.New("Invalid Age")
		return
	}
	jsonBytes, err = json.Marshal(person)
	jsonString = string(jsonBytes)
	return
}

func DoSomething(c chan string, shouldClose bool) {
	c <- "Done!"
	if shouldClose {
		close(c)
	}
}
