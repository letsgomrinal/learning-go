package lesson16

import (
	"iter"
	"slices"
	"testing"
)

// collect gathers up to max values, so a broken infinite iterator can't hang the tests.
func collect[T any](seq iter.Seq[T], max int) []T {
	var out []T
	for v := range seq {
		out = append(out, v)
		if len(out) >= max {
			break
		}
	}
	return out
}

// stopsEarly checks that the iterator respects break (if it doesn't, Go panics).
func stopsEarly[T any](t *testing.T, name string, seq iter.Seq[T]) {
	t.Helper()
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("%s kept going after the loop broke: %v. Use: if !yield(v) { return }", name, r)
		}
	}()
	for range seq {
		break
	}
}

func TestEx1Countdown(t *testing.T) {
	if got := collect(Countdown(3), 100); !slices.Equal(got, []int{3, 2, 1}) {
		t.Errorf("Ex1: Countdown(3) produced %v, want [3 2 1]", got)
	}
	if got := collect(Countdown(0), 100); len(got) != 0 {
		t.Errorf("Ex1: Countdown(0) produced %v, want nothing", got)
	}
	stopsEarly(t, "Ex1: Countdown", Countdown(5))
}

func TestEx2EvensBelow(t *testing.T) {
	if got := collect(EvensBelow(7), 100); !slices.Equal(got, []int{0, 2, 4, 6}) {
		t.Errorf("Ex2: EvensBelow(7) produced %v, want [0 2 4 6]", got)
	}
	if got := collect(EvensBelow(6), 100); !slices.Equal(got, []int{0, 2, 4}) {
		t.Errorf("Ex2: EvensBelow(6) produced %v, want [0 2 4]", got)
	}
	stopsEarly(t, "Ex2: EvensBelow", EvensBelow(100))
}

func TestEx3Fibonacci(t *testing.T) {
	got := collect(Fibonacci(), 10)
	if want := []int{0, 1, 1, 2, 3, 5, 8, 13, 21, 34}; !slices.Equal(got, want) {
		t.Errorf("Ex3: first 10 Fibonacci numbers = %v, want %v", got, want)
	}
	stopsEarly(t, "Ex3: Fibonacci", Fibonacci())
}

func TestEx4Take(t *testing.T) {
	letters := slices.Values([]string{"a", "b", "c"})
	if got := Take(letters, 2); !slices.Equal(got, []string{"a", "b"}) {
		t.Errorf("Ex4: Take([a b c], 2) = %q, want [a b]", got)
	}
	if got := Take(letters, 10); !slices.Equal(got, []string{"a", "b", "c"}) {
		t.Errorf("Ex4: Take([a b c], 10) = %q, want [a b c]", got)
	}
	if got := Take(letters, 0); len(got) != 0 {
		t.Errorf("Ex4: Take(seq, 0) = %q, want empty", got)
	}
	// An infinite sequence that doesn't depend on your Ex3:
	naturals := func(yield func(int) bool) {
		for i := 0; ; i++ {
			if !yield(i) {
				return
			}
		}
	}
	if got := Take(naturals, 4); !slices.Equal(got, []int{0, 1, 2, 3}) {
		t.Errorf("Ex4: Take(naturals, 4) = %v, want [0 1 2 3] (did you break when you had enough?)", got)
	}
}

func TestEx5FilterSeq(t *testing.T) {
	odd := func(n int) bool { return n%2 == 1 }
	got := collect(FilterSeq(slices.Values([]int{1, 2, 3, 4, 5}), odd), 100)
	if !slices.Equal(got, []int{1, 3, 5}) {
		t.Errorf("Ex5: FilterSeq([1 2 3 4 5], odd) produced %v, want [1 3 5]", got)
	}
	stopsEarly(t, "Ex5: FilterSeq", FilterSeq(slices.Values([]int{1, 3, 5}), odd))
}

func TestEx6NumberedLines(t *testing.T) {
	var nums []int
	var lines []string
	for n, line := range NumberedLines("first\n\nthird\nfourth") {
		nums = append(nums, n)
		lines = append(lines, line)
	}
	if !slices.Equal(nums, []int{1, 3, 4}) || !slices.Equal(lines, []string{"first", "third", "fourth"}) {
		t.Errorf("Ex6: got numbers %v and lines %q; want [1 3 4] and [first third fourth]", nums, lines)
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Ex6: NumberedLines kept going after the loop broke: %v", r)
			}
		}()
		for range NumberedLines("a\nb\nc") {
			break
		}
	}()
}

func TestEx7SortedWords(t *testing.T) {
	got := SortedWords(map[string]int{"pear": 1, "apple": 5, "fig": 2})
	if want := []string{"apple", "fig", "pear"}; !slices.Equal(got, want) {
		t.Errorf("Ex7: SortedWords = %q, want %q", got, want)
	}
}
