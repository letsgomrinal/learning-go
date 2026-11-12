// Lesson 16 exercises: iterators. Read README.md first!
//
// Check one exercise at a time:   go test -run Ex1 -v
//
// Every iterator must stop when yield returns false:   if !yield(v) { return }
package lesson16

import "iter"

// Ex1: produce n, n-1, ..., 1. For n <= 0, produce nothing.
//
//	for x := range Countdown(3) { ... }   // 3, 2, 1
func Countdown(n int) iter.Seq[int] {
	return func(yield func(int) bool) {
		// TODO
	}
}

// Ex2: produce the even numbers from 0 up to (but not including) limit.
// EvensBelow(7) produces 0, 2, 4, 6
func EvensBelow(limit int) iter.Seq[int] {
	return func(yield func(int) bool) {
		// TODO
	}
}

// Ex3: produce the Fibonacci numbers forever: 0, 1, 1, 2, 3, 5, 8, 13...
// (Each number is the sum of the previous two.) It only works because the caller breaks.
//
// Hint: a, b := 0, 1, then loop forever: yield a, then a, b = b, a+b
func Fibonacci() iter.Seq[int] {
	return func(yield func(int) bool) {
		// TODO
	}
}

// Ex4: collect the first n values of any sequence into a slice. If the sequence has fewer
// than n values, return all of them. For n <= 0 return an empty slice.
// Take(Fibonacci(), 5) == [0 1 1 2 3]
//
// Hint: for v := range seq { ... append ...; if len(result) == n { break } }
// (Careful: with n <= 0, don't range at all.)
func Take[T any](seq iter.Seq[T], n int) []T {
	// TODO
	return nil
}

// Ex5: return a new sequence with only the values of seq for which keep returns true.
// It's lazy: nothing happens until someone ranges over the result. (README section 6.)
func FilterSeq[T any](seq iter.Seq[T], keep func(T) bool) iter.Seq[T] {
	return func(yield func(T) bool) {
		// TODO
	}
}

// Ex6: produce (line number, line) pairs for the lines in text, numbered from 1.
// Lines are separated by "\n". Skip empty lines, but keep counting them.
//
//	for n, line := range NumberedLines("a\n\nb") { ... }   // (1, "a"), (3, "b")
//
// Hint: strings.Split(text, "\n") (import "strings"), then loop with the index.
func NumberedLines(text string) iter.Seq2[int, string] {
	return func(yield func(int, string) bool) {
		// TODO
	}
}

// Ex7: given word counts, return the words sorted alphabetically.
// Do it in ONE line with the standard library: slices.Sorted and maps.Keys
// (import "maps" and "slices"). README section 5.
func SortedWords(counts map[string]int) []string {
	// TODO
	return nil
}
