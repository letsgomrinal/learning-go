// Package mathxcheck checks Ex4. Run:  go test -v ./mathx_check
//
// If you see an error like "package gomastery/course/part1/09_packages_tests/mathx is not in std"
// or "no required module provides package", the mathx folder doesn't exist yet (or has no .go file).
package mathxcheck

import (
	"testing"

	"gomastery/course/part1/09_packages_tests/mathx"
)

func TestEx4Clamp(t *testing.T) {
	tests := []struct{ v, lo, hi, want int }{
		{5, 0, 10, 5},
		{-3, 0, 10, 0},
		{42, 0, 10, 10},
		{0, 0, 10, 0},
		{10, 0, 10, 10},
		{-7, -10, -5, -7},
	}
	for _, tt := range tests {
		if got := mathx.Clamp(tt.v, tt.lo, tt.hi); got != tt.want {
			t.Errorf("Ex4: Clamp(%d, %d, %d) = %d, want %d", tt.v, tt.lo, tt.hi, got, tt.want)
		}
	}
}
