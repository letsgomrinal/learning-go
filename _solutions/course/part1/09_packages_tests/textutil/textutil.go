// Package textutil has small helpers for working with text. (Reference solution.)
package textutil

import "strings"

// WordCount returns the number of words in s, split on any whitespace.
func WordCount(s string) int {
	return len(strings.Fields(s))
}

// Capitalize returns s with its first letter in upper case.
// It assumes plain ASCII text.
func Capitalize(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func isVowel(r rune) bool {
	return strings.ContainsRune("aeiouAEIOU", r)
}

// CountVowels returns the number of vowels in s.
func CountVowels(s string) int {
	n := 0
	for _, r := range s {
		if isVowel(r) {
			n++
		}
	}
	return n
}

// Truncate shortens s to at most max characters, adding "..." if it cut anything.
func Truncate(s string, max int) string {
	runes := []rune(s)
	if len(runes) <= max {
		return s
	}
	return string(runes[:max]) + "..."
}
