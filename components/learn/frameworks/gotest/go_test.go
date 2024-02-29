package gotest

import "testing"

var per = Person{"Zoye", 44, 8983333}
var personEncoder PersonEncoder = &PersonEncoderImpl{}

func TestGoEncode(t *testing.T) {
	_, err := personEncoder.EncodePerson(per)
	if err != nil {
		t.Errorf("Encode Failied for Person %v", per)
	}
}

func BenchmarkEncode(b *testing.B) {
	for n := 0; n < b.N; n++ {
		personEncoder.EncodePerson(per)
	}
}
