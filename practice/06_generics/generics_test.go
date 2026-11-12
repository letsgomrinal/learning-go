package generics

import (
	"slices"
	"strconv"
	"testing"
)

func TestMapFilterReduce(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5}

	strs := Map(nums, func(n int) string { return strconv.Itoa(n * n) })
	if want := []string{"1", "4", "9", "16", "25"}; !slices.Equal(strs, want) {
		t.Errorf("Map = %q, want %q", strs, want)
	}

	evens := Filter(nums, func(n int) bool { return n%2 == 0 })
	if want := []int{2, 4}; !slices.Equal(evens, want) {
		t.Errorf("Filter = %v, want %v", evens, want)
	}

	joined := Reduce(nums, "", func(acc string, n int) string { return acc + strconv.Itoa(n) })
	if joined != "12345" {
		t.Errorf("Reduce = %q, want %q", joined, "12345")
	}
}

func TestSum(t *testing.T) {
	if got := Sum([]int{1, 2, 3}); got != 6 {
		t.Errorf("Sum(ints) = %v", got)
	}
	if got := Sum([]float64{0.5, 0.25}); got != 0.75 {
		t.Errorf("Sum(floats) = %v", got)
	}
	if got := Sum([]Celsius{20, 1.5}); got != 21.5 {
		t.Errorf("Sum(Celsius) = %v", got)
	}
	if got := Sum[int](nil); got != 0 {
		t.Errorf("Sum(nil) = %v", got)
	}
}

func TestSet(t *testing.T) {
	s := NewSet("a", "b", "a")
	if s.Len() != 2 || !s.Has("a") || !s.Has("b") || s.Has("z") {
		t.Fatalf("NewSet(a, b, a): len %d", s.Len())
	}

	var zero Set[int] // zero value must not panic on Add
	zero.Add(1)
	if !zero.Has(1) {
		t.Error("Add on zero-value Set didn't work")
	}

	u := s.Union(NewSet("c", "a"))
	if u.Len() != 3 {
		t.Errorf("Union len = %d, want 3", u.Len())
	}
	if s.Len() != 2 {
		t.Errorf("Union modified the receiver")
	}

	var got []string
	for v := range u.All() {
		got = append(got, v)
	}
	slices.Sort(got)
	if want := []string{"a", "b", "c"}; !slices.Equal(got, want) {
		t.Errorf("All() yielded %v, want %v", got, want)
	}

	// Breaking early must not panic (you must respect yield's return value).
	for range u.All() {
		break
	}
}

func TestFibonacciAndTake(t *testing.T) {
	got := Take(Fibonacci(), 10)
	if want := []int{0, 1, 1, 2, 3, 5, 8, 13, 21, 34}; !slices.Equal(got, want) {
		t.Errorf("Take(Fibonacci, 10) = %v, want %v", got, want)
	}
	if got := Take(slices.Values([]string{"x", "y"}), 5); !slices.Equal(got, []string{"x", "y"}) {
		t.Errorf("Take(short seq) = %v", got)
	}
}
