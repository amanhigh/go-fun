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
	if assert.NoError(t, err) {
		assert.Equal(t, personJson, result)
	}
}

func TestDecode(t *testing.T) {
	decodedPerson, err := personEncoder.DecodePerson(personJson)
	if assert.NoError(t, err) {
		assert.Equal(t, person, decodedPerson)
	}
}

func TestEncodeTable(t *testing.T) {
	tests := []struct {
		name    string
		person  Person
		want    string
		wantErr bool
	}{
		{
			name:    "Valid Person",
			person:  Person{"Zoye", 44, 8983333},
			want:    `{"name":"Zoye","Age":44,"MobileNumber":8983333}`,
			wantErr: false,
		},
		// Add more negative test cases here
		{
			name:    "Invalid Person - Negative Age",
			person:  Person{"Zoye", -44, 8983333},
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := personEncoder.EncodePerson(tt.person)
			if (err != nil) != tt.wantErr {
				t.Errorf("EncodePerson() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("EncodePerson() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkEncode(b *testing.B) {
	for n := 0; n < b.N; n++ {
		personEncoder.EncodePerson(person)
	}
}
