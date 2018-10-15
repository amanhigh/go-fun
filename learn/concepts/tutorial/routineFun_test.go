package main

import (
	"testing"
)

func TestSum(t *testing.T) {
	cases := []struct {
		in, want int
	}{
		{0, 12},
	}
	for _, c := range cases {
		got := sumFun()
		if got != c.want {
			t.Errorf("GoRoutineFun(%v) == %v, want %v", c.in, got, c.want)
		}
	}
}
