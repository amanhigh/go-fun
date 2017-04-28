package main

import (
	"testing"
)

func Test_encodePerson(t *testing.T) {
	type args struct {
		p1 person
	}
	per := person{"Zoye", 44, 8983333}
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
