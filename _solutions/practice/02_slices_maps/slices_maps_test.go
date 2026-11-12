package slicesmaps

import (
	"maps"
	"slices"
	"testing"
)

func TestRemoveAt(t *testing.T) {
	orig := []int{10, 20, 30, 40}
	got := RemoveAt(orig, 1)
	if want := []int{10, 30, 40}; !slices.Equal(got, want) {
		t.Errorf("RemoveAt = %v, want %v", got, want)
	}
	if want := []int{10, 20, 30, 40}; !slices.Equal(orig, want) {
		t.Errorf("input was modified: %v, want %v (did you use append(s[:i], s[i+1:]...)?)", orig, want)
	}
	if got := RemoveAt(orig, 3); !slices.Equal(got, []int{10, 20, 30}) {
		t.Errorf("RemoveAt(last) = %v", got)
	}
	if got := RemoveAt(orig, 9); !slices.Equal(got, orig) {
		t.Errorf("RemoveAt(out of range) = %v, want %v", got, orig)
	}
}

func TestSafeAppend(t *testing.T) {
	base := make([]int, 3, 10) // len 3, cap 10: spare capacity!

	// This is what goes wrong with a plain append. Run it and see:
	x := append(base, 1)
	y := append(base, 2) // writes into the SAME slot x[3] uses
	if x[3] != 2 {
		t.Fatalf("expected the aliasing bug to show up: x[3]=%d", x[3])
	}
	_ = y

	a := SafeAppend(base, 1)
	b := SafeAppend(base, 2)
	if len(a) != 4 || len(b) != 4 {
		t.Fatalf("len(a)=%d len(b)=%d, want 4 and 4", len(a), len(b))
	}
	if a[3] != 1 || b[3] != 2 {
		t.Errorf("a[3]=%d b[3]=%d, want 1 and 2. Results share memory", a[3], b[3])
	}
	a[0] = 42
	if base[0] == 42 {
		t.Errorf("modifying the result changed base. They share memory")
	}
}

func TestChunk(t *testing.T) {
	got := Chunk([]int{1, 2, 3, 4, 5}, 2)
	want := [][]int{{1, 2}, {3, 4}, {5}}
	if !slices.EqualFunc(got, want, slices.Equal[[]int]) {
		t.Errorf("Chunk = %v, want %v", got, want)
	}
	if got := Chunk([]int{1, 2}, 5); len(got) != 1 || !slices.Equal(got[0], []int{1, 2}) {
		t.Errorf("Chunk(size>len) = %v", got)
	}
	if got := Chunk([]int{1}, 0); got != nil {
		t.Errorf("Chunk(size 0) = %v, want nil", got)
	}
}

func TestDedupe(t *testing.T) {
	got := Dedupe([]string{"b", "a", "b", "c", "a"})
	if want := []string{"b", "a", "c"}; !slices.Equal(got, want) {
		t.Errorf("Dedupe = %v, want %v", got, want)
	}
}

func TestGroupByLen(t *testing.T) {
	got := GroupByLen([]string{"go", "is", "fun", "yes", "a"})
	want := map[int][]string{1: {"a"}, 2: {"go", "is"}, 3: {"fun", "yes"}}
	if !maps.EqualFunc(got, want, slices.Equal[[]string]) {
		t.Errorf("GroupByLen = %v, want %v", got, want)
	}
}

func TestSortedKeys(t *testing.T) {
	got := SortedKeys(map[string]int{"pear": 1, "apple": 2, "fig": 3})
	if want := []string{"apple", "fig", "pear"}; !slices.Equal(got, want) {
		t.Errorf("SortedKeys = %v, want %v", got, want)
	}
}

func TestMatrix(t *testing.T) {
	m := Matrix(3, 2)
	if len(m) != 3 {
		t.Fatalf("rows = %d, want 3", len(m))
	}
	for i, row := range m {
		if len(row) != 2 {
			t.Fatalf("row %d len = %d, want 2", i, len(row))
		}
	}
	m[0][0] = 7
	if m[1][0] == 7 || m[2][0] == 7 {
		t.Errorf("rows share memory")
	}
}
