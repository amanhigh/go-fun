package gotest

import "testing"

var per = person{"Zoye", 44, 8983333}

func TestGoEncode(t *testing.T) {
	_, err := encodePerson(per)
	if err != nil {
		t.Errorf("Encode Failied for Person %v", per)
	}
}

func BenchmarkEncode(b *testing.B) {

	for n := 0; n < b.N; n++ {
		encodePerson(per)
	}
}
