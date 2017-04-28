package main

import (
	"testing"
)

var per = person{"Zoye", 44, 8983333}

func Test_encodePerson(t *testing.T) {
	type args struct {
		p1 person
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"Encode Test", args{per}, "{\"name\":\"Zoye\",\"Age\":44,\"MobileNumber\":8983333}"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := encodePerson(tt.args.p1); got != tt.want {
				t.Errorf("encodePerson() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkEncode(b *testing.B) {
	for n := 0; n < b.N; n++ {
		encodePerson(per)
	}
}
