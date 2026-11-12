// Reference solution for lesson 16.
package lesson16

import (
	"iter"
	"maps"
	"slices"
	"strings"
)

func Countdown(n int) iter.Seq[int] {
	return func(yield func(int) bool) {
		for i := n; i >= 1; i-- {
			if !yield(i) {
				return
			}
		}
	}
}

func EvensBelow(limit int) iter.Seq[int] {
	return func(yield func(int) bool) {
		for i := 0; i < limit; i += 2 {
			if !yield(i) {
				return
			}
		}
	}
}

func Fibonacci() iter.Seq[int] {
	return func(yield func(int) bool) {
		a, b := 0, 1
		for {
			if !yield(a) {
				return
			}
			a, b = b, a+b
		}
	}
}

func Take[T any](seq iter.Seq[T], n int) []T {
	result := []T{}
	if n <= 0 {
		return result
	}
	for v := range seq {
		result = append(result, v)
		if len(result) == n {
			break
		}
	}
	return result
}

func FilterSeq[T any](seq iter.Seq[T], keep func(T) bool) iter.Seq[T] {
	return func(yield func(T) bool) {
		for v := range seq {
			if keep(v) && !yield(v) {
				return
			}
		}
	}
}

func NumberedLines(text string) iter.Seq2[int, string] {
	return func(yield func(int, string) bool) {
		for i, line := range strings.Split(text, "\n") {
			if line == "" {
				continue
			}
			if !yield(i+1, line) {
				return
			}
		}
	}
}

func SortedWords(counts map[string]int) []string {
	return slices.Sorted(maps.Keys(counts))
}
