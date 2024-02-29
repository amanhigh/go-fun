package gotest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var person = Person{"Zoye", 44, 8983333}
var personJson = `{"name":"Zoye","Age":44,"MobileNumber":8983333}`
var personEncoder PersonEncoder = &PersonEncoderImpl{}

func TestEncode(t *testing.T) {
	result, err := personEncoder.EncodePerson(person)
	assert.NoError(t, err)
	assert.Equal(t, personJson, result)
}

func TestDecode(t *testing.T) {
	decodedPerson, err := personEncoder.DecodePerson(personJson)
	assert.NoError(t, err)
	assert.Equal(t, person, decodedPerson)
}

func BenchmarkEncode(b *testing.B) {
	for n := 0; n < b.N; n++ {
		personEncoder.EncodePerson(person)
	}
}
