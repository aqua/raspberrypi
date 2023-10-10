package main

import (
	"testing"
)

func TestIsPrime(t *testing.T) {
	cases := []struct {
		val  uint
		want bool
	}{
		{0, false},
		{1, false},
		{2, true},
		{3, true},
		{4, false},
		{5, true},
		{6, false},
		{7, true},
		{8, false},
		{9, false},
		{10, false},
		{11, true},
	}
	for _, c := range cases {
		if got := is_prime(c.val); got != c.want {
			t.Errorf("got %d prime=%v, want %v", c.val, got, c.want)
		}
	}
}
