package lesson05

import (
	"slices"
	"testing"
)

func TestEx1Sum(t *testing.T) {
	tests := []struct {
		in   []int
		want int
	}{
		{[]int{1, 2, 3}, 6},
		{[]int{10, -10, 5}, 5},
		{[]int{}, 0},
		{nil, 0},
	}
	for _, tt := range tests {
		if got := Sum(tt.in); got != tt.want {
			t.Errorf("Ex1: Sum(%v) = %d, want %d", tt.in, got, tt.want)
		}
	}
}

func TestEx2Largest(t *testing.T) {
	tests := []struct {
		in     []int
		want   int
		wantOK bool
	}{
		{[]int{3, 9, 2}, 9, true},
		{[]int{-5, -2, -9}, -2, true},
		{[]int{7}, 7, true},
		{[]int{}, 0, false},
	}
	for _, tt := range tests {
		got, ok := Largest(tt.in)
		if got != tt.want || ok != tt.wantOK {
			t.Errorf("Ex2: Largest(%v) = %d, %v; want %d, %v", tt.in, got, ok, tt.want, tt.wantOK)
		}
	}
}

func TestEx3Evens(t *testing.T) {
	got := Evens([]int{1, 2, 3, 4, 6, 7})
	want := []int{2, 4, 6}
	if !slices.Equal(got, want) {
		t.Errorf("Ex3: Evens([1 2 3 4 6 7]) = %v, want %v", got, want)
	}
	if got := Evens([]int{1, 3}); len(got) != 0 {
		t.Errorf("Ex3: Evens([1 3]) = %v, want empty", got)
	}
}

func TestEx4Doubled(t *testing.T) {
	in := []int{1, 2, 3}
	got := Doubled(in)
	if want := []int{2, 4, 6}; !slices.Equal(got, want) {
		t.Errorf("Ex4: Doubled([1 2 3]) = %v, want %v", got, want)
	}
	if !slices.Equal(in, []int{1, 2, 3}) {
		t.Errorf("Ex4: Doubled changed its input to %v. It should return a NEW slice", in)
	}
}

func TestEx5DoubleInPlace(t *testing.T) {
	nums := []int{1, 2, 3}
	DoubleInPlace(nums)
	if want := []int{2, 4, 6}; !slices.Equal(nums, want) {
		t.Errorf("Ex5: after DoubleInPlace, nums = %v, want %v", nums, want)
	}
}

func TestEx6IndexOf(t *testing.T) {
	items := []string{"a", "b", "c", "b"}
	tests := map[string]int{"a": 0, "b": 1, "c": 2, "z": -1}
	for target, want := range tests {
		if got := IndexOf(items, target); got != want {
			t.Errorf("Ex6: IndexOf(%q, %q) = %d, want %d", items, target, got, want)
		}
	}
	if got := IndexOf(nil, "a"); got != -1 {
		t.Errorf("Ex6: IndexOf(empty, \"a\") = %d, want -1", got)
	}
}

func TestEx7Average(t *testing.T) {
	if got := Average(1, 2, 3, 4); got != 2.5 {
		t.Errorf("Ex7: Average(1, 2, 3, 4) = %v, want 2.5", got)
	}
	if got := Average(); got != 0 {
		t.Errorf("Ex7: Average() = %v, want 0", got)
	}
	vals := []float64{10, 20}
	if got := Average(vals...); got != 15 {
		t.Errorf("Ex7: Average(10, 20) = %v, want 15", got)
	}
}

func TestEx8Reversed(t *testing.T) {
	in := []int{1, 2, 3, 4}
	got := Reversed(in)
	if want := []int{4, 3, 2, 1}; !slices.Equal(got, want) {
		t.Errorf("Ex8: Reversed([1 2 3 4]) = %v, want %v", got, want)
	}
	if !slices.Equal(in, []int{1, 2, 3, 4}) {
		t.Errorf("Ex8: Reversed changed its input to %v", in)
	}
	if got := Reversed(nil); len(got) != 0 {
		t.Errorf("Ex8: Reversed(empty) = %v, want empty", got)
	}
}

func TestEx9FirstAndLast(t *testing.T) {
	tests := []struct {
		in          []string
		first, last string
		ok          bool
	}{
		{[]string{"x", "y", "z"}, "x", "z", true},
		{[]string{"solo"}, "solo", "solo", true},
		{[]string{}, "", "", false},
	}
	for _, tt := range tests {
		f, l, ok := FirstAndLast(tt.in)
		if f != tt.first || l != tt.last || ok != tt.ok {
			t.Errorf("Ex9: FirstAndLast(%q) = %q, %q, %v; want %q, %q, %v", tt.in, f, l, ok, tt.first, tt.last, tt.ok)
		}
	}
}
