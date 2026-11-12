package lesson15

import (
	"slices"
	"strings"
	"testing"
)

func TestEx1Contains(t *testing.T) {
	if !Contains([]string{"a", "b"}, "b") {
		t.Error("Ex1: Contains([a b], b) = false, want true")
	}
	if Contains([]int{1, 2}, 5) {
		t.Error("Ex1: Contains([1 2], 5) = true, want false")
	}
	if Contains([]float64{}, 1.5) {
		t.Error("Ex1: Contains(empty, 1.5) = true, want false")
	}
	type point struct{ X, Y int }
	if !Contains([]point{{1, 2}, {3, 4}}, point{3, 4}) {
		t.Error("Ex1: Contains should work with comparable structs too")
	}
}

func TestEx2Max(t *testing.T) {
	if v, ok := Max([]int{3, 9, 2}); v != 9 || !ok {
		t.Errorf("Ex2: Max([3 9 2]) = %v, %v; want 9, true", v, ok)
	}
	if v, ok := Max([]int{-5, -2, -9}); v != -2 || !ok {
		t.Errorf("Ex2: Max([-5 -2 -9]) = %v, %v; want -2, true", v, ok)
	}
	if v, ok := Max([]string{"go", "py", "c"}); v != "py" || !ok {
		t.Errorf("Ex2: Max([go py c]) = %q, %v; want \"py\", true", v, ok)
	}
	if v, ok := Max([]float64{}); v != 0 || ok {
		t.Errorf("Ex2: Max(empty) = %v, %v; want 0, false", v, ok)
	}
}

func TestEx3Map(t *testing.T) {
	got := Map([]string{"go", "rust", "c"}, func(s string) int { return len(s) })
	if want := []int{2, 4, 1}; !slices.Equal(got, want) {
		t.Errorf("Ex3: Map(lengths) = %v, want %v", got, want)
	}
	upper := Map([]string{"a", "b"}, strings.ToUpper)
	if want := []string{"A", "B"}; !slices.Equal(upper, want) {
		t.Errorf("Ex3: Map(ToUpper) = %q, want %q", upper, want)
	}
}

func TestEx4Filter(t *testing.T) {
	got := Filter([]int{1, 2, 3, 4, 5, 6}, func(n int) bool { return n%2 == 0 })
	if want := []int{2, 4, 6}; !slices.Equal(got, want) {
		t.Errorf("Ex4: Filter(even) = %v, want %v", got, want)
	}
	long := Filter([]string{"go", "gopher", "gc"}, func(s string) bool { return len(s) > 2 })
	if want := []string{"gopher"}; !slices.Equal(long, want) {
		t.Errorf("Ex4: Filter(long words) = %q, want %q", long, want)
	}
}

func TestEx5Sum(t *testing.T) {
	if got := Sum([]int{1, 2, 3}); got != 6 {
		t.Errorf("Ex5: Sum([1 2 3]) = %v, want 6", got)
	}
	if got := Sum([]float64{0.5, 0.25}); got != 0.75 {
		t.Errorf("Ex5: Sum([0.5 0.25]) = %v, want 0.75", got)
	}
	if got := Sum([]int{}); got != 0 {
		t.Errorf("Ex5: Sum(empty) = %v, want 0", got)
	}
	if got := Sum([]Celsius{20.5, 22}); got != 42.5 {
		t.Errorf("Ex5: Sum([]Celsius{20.5, 22}) = %v, want 42.5", got)
	}
}

func TestEx6Stack(t *testing.T) {
	var s Stack[string]
	if _, ok := s.Pop(); ok {
		t.Error("Ex6: Pop on an empty stack returned ok = true")
	}
	s.Push("a")
	s.Push("b")
	s.Push("c")
	if s.Len() != 3 {
		t.Fatalf("Ex6: Len = %d, want 3", s.Len())
	}
	if v, ok := s.Peek(); v != "c" || !ok {
		t.Errorf("Ex6: Peek = %q, %v; want \"c\", true", v, ok)
	}
	for _, want := range []string{"c", "b", "a"} {
		if v, ok := s.Pop(); v != want || !ok {
			t.Errorf("Ex6: Pop = %q, %v; want %q, true", v, ok, want)
		}
	}
	nums := &Stack[int]{}
	nums.Push(42)
	if v, _ := nums.Pop(); v != 42 {
		t.Errorf("Ex6: Stack[int] Pop = %d, want 42", v)
	}
}

func TestEx7SortedKeys(t *testing.T) {
	got := SortedKeys(map[string]int{"pear": 1, "apple": 2, "fig": 3})
	if want := []string{"apple", "fig", "pear"}; !slices.Equal(got, want) {
		t.Errorf("Ex7: SortedKeys(strings) = %q, want %q", got, want)
	}
	ints := SortedKeys(map[int]bool{3: true, 1: false, 2: true})
	if want := []int{1, 2, 3}; !slices.Equal(ints, want) {
		t.Errorf("Ex7: SortedKeys(ints) = %v, want %v", ints, want)
	}
}
