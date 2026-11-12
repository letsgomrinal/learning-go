package lesson06

import (
	"maps"
	"slices"
	"testing"
)

func TestEx1Lookup(t *testing.T) {
	prices := map[string]float64{"apple": 0.5, "free sample": 0}
	if p, ok := Lookup(prices, "apple"); p != 0.5 || !ok {
		t.Errorf("Ex1: Lookup(apple) = %v, %v; want 0.5, true", p, ok)
	}
	if p, ok := Lookup(prices, "free sample"); p != 0 || !ok {
		t.Errorf("Ex1: Lookup(free sample) = %v, %v; want 0, true (it exists, with price 0)", p, ok)
	}
	if p, ok := Lookup(prices, "caviar"); p != 0 || ok {
		t.Errorf("Ex1: Lookup(caviar) = %v, %v; want 0, false", p, ok)
	}
}

func TestEx2WordCount(t *testing.T) {
	got := WordCount([]string{"go", "is", "go", "fun", "go"})
	want := map[string]int{"go": 3, "is": 1, "fun": 1}
	if !maps.Equal(got, want) {
		t.Errorf("Ex2: WordCount = %v, want %v", got, want)
	}
	if got := WordCount(nil); len(got) != 0 {
		t.Errorf("Ex2: WordCount(empty) = %v, want empty", got)
	}
}

func TestEx3HasDuplicates(t *testing.T) {
	tests := []struct {
		in   []int
		want bool
	}{
		{[]int{1, 2, 3, 2}, true},
		{[]int{1, 2, 3}, false},
		{[]int{5, 5}, true},
		{nil, false},
	}
	for _, tt := range tests {
		if got := HasDuplicates(tt.in); got != tt.want {
			t.Errorf("Ex3: HasDuplicates(%v) = %v, want %v", tt.in, got, tt.want)
		}
	}
}

func TestEx4Unique(t *testing.T) {
	got := Unique([]string{"b", "a", "b", "c", "a"})
	if want := []string{"b", "a", "c"}; !slices.Equal(got, want) {
		t.Errorf("Ex4: Unique = %q, want %q", got, want)
	}
}

func TestEx5Invert(t *testing.T) {
	got := Invert(map[string]int{"one": 1, "two": 2, "three": 3})
	want := map[int]string{1: "one", 2: "two", 3: "three"}
	if !maps.Equal(got, want) {
		t.Errorf("Ex5: Invert = %v, want %v", got, want)
	}
}

func TestEx6SortedKeys(t *testing.T) {
	got := SortedKeys(map[string]int{"pear": 1, "apple": 2, "fig": 3})
	if want := []string{"apple", "fig", "pear"}; !slices.Equal(got, want) {
		t.Errorf("Ex6: SortedKeys = %q, want %q", got, want)
	}
}

func TestEx7GroupByFirstLetter(t *testing.T) {
	got := GroupByFirstLetter([]string{"ana", "bo", "amy", "", "cy"})
	want := map[string][]string{"a": {"ana", "amy"}, "b": {"bo"}, "c": {"cy"}}
	if !maps.EqualFunc(got, want, slices.Equal[[]string]) {
		t.Errorf("Ex7: GroupByFirstLetter = %q, want %q", got, want)
	}
}

func TestEx8MostCommon(t *testing.T) {
	tests := []struct {
		in   []string
		want string
	}{
		{[]string{"a", "b", "b", "c"}, "b"},
		{[]string{"go", "py", "go", "py", "rs"}, "go"}, // tie: alphabetical
		{[]string{"zed"}, "zed"},
		{nil, ""},
	}
	for _, tt := range tests {
		if got := MostCommon(tt.in); got != tt.want {
			t.Errorf("Ex8: MostCommon(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}
