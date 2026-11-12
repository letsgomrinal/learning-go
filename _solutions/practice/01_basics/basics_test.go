package basics

// Read the tests too: this is what idiomatic TABLE-DRIVEN tests look like in Go.
// There's no pytest; `testing` is in the standard library.

import (
	"maps"
	"math"
	"slices"
	"testing"
)

func TestFizzBuzz(t *testing.T) {
	got := FizzBuzz(15)
	want := []string{"1", "2", "Fizz", "4", "Buzz", "Fizz", "7", "8", "Fizz", "Buzz", "11", "Fizz", "13", "14", "FizzBuzz"}
	// Slices can't be compared with ==. Use slices.Equal (or reflect.DeepEqual).
	if !slices.Equal(got, want) {
		t.Errorf("FizzBuzz(15)\n got: %q\nwant: %q", got, want)
	}
	if got := FizzBuzz(0); len(got) != 0 {
		t.Errorf("FizzBuzz(0) = %q, want empty", got)
	}
}

func TestReverseRunes(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{"", ""},
		{"abc", "cba"},
		{"héllo, 世界", "界世 ,olléh"},
		{"🙂ok", "ko🙂"},
	}
	for _, tt := range tests {
		if got := ReverseRunes(tt.in); got != tt.want {
			t.Errorf("ReverseRunes(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestCountWords(t *testing.T) {
	got := CountWords("Go go  GO\nstop\tStop now")
	want := map[string]int{"go": 3, "stop": 2, "now": 1}
	if !maps.Equal(got, want) {
		t.Errorf("CountWords = %v, want %v", got, want)
	}
	if got := CountWords("   "); len(got) != 0 {
		t.Errorf("CountWords(blank) = %v, want empty", got)
	}
}

func TestStats(t *testing.T) {
	tests := []struct {
		name     string
		in       []int
		min, max int
		avg      float64
	}{
		{"empty", nil, 0, 0, 0},
		{"one", []int{5}, 5, 5, 5},
		{"mixed", []int{3, -1, 4, 1, 5}, -1, 5, 2.4},
		{"int division trap", []int{1, 2}, 1, 2, 1.5},
	}
	for _, tt := range tests {
		// t.Run creates a named subtest: go test -run TestStats/mixed
		t.Run(tt.name, func(t *testing.T) {
			mn, mx, avg := Stats(tt.in...)
			if mn != tt.min || mx != tt.max || math.Abs(avg-tt.avg) > 1e-9 {
				t.Errorf("Stats(%v) = (%d, %d, %v), want (%d, %d, %v)", tt.in, mn, mx, avg, tt.min, tt.max, tt.avg)
			}
		})
	}
}

func TestDivide(t *testing.T) {
	if q, ok := Divide(7, 2); !ok || q != 3 {
		t.Errorf("Divide(7, 2) = (%d, %v), want (3, true)", q, ok)
	}
	if q, ok := Divide(-7, 2); !ok || q != -3 {
		// Note: Go truncates toward zero. Python's -7 // 2 == -4!
		t.Errorf("Divide(-7, 2) = (%d, %v), want (-3, true)", q, ok)
	}
	if q, ok := Divide(1, 0); ok || q != 0 {
		t.Errorf("Divide(1, 0) = (%d, %v), want (0, false)", q, ok)
	}
}

func TestGrade(t *testing.T) {
	tests := map[int]string{100: "A", 90: "A", 89: "B", 80: "B", 75: "C", 60: "D", 59: "F", 0: "F", -1: "invalid", 101: "invalid"}
	for score, want := range tests {
		if got := Grade(score); got != want {
			t.Errorf("Grade(%d) = %q, want %q", score, got, want)
		}
	}
}

func TestIsWeekend(t *testing.T) {
	for d := Sunday; d <= Saturday; d++ {
		want := d == Saturday || d == Sunday
		if got := d.IsWeekend(); got != want {
			t.Errorf("Weekday(%d).IsWeekend() = %v, want %v", d, got, want)
		}
	}
}
