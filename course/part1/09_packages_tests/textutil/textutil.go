// Package textutil has small helpers for working with text.
//
// Lesson 09 exercises happen in this file. Read ../README.md first!
package textutil

import "strings"

// Ex1: these two functions are lowercase, so other packages can't use them.
// Rename them so they're exported: wordCount → WordCount, capitalize → Capitalize.
// (Doc comments start with the function's name, so update those too.)

// wordCount returns the number of words in s, split on any whitespace.
func wordCount(s string) int {
	return len(strings.Fields(s))
}

// capitalize returns s with its first letter in upper case.
// It assumes plain ASCII text.
func capitalize(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// isVowel is an internal helper. Leave it lowercase: nothing outside this package needs it.
func isVowel(r rune) bool {
	return strings.ContainsRune("aeiouAEIOU", r)
}

// CountVowels is exported and uses the private helper above.
// Code inside the package can use private names freely.
func CountVowels(s string) int {
	n := 0
	for _, r := range s {
		if isVowel(r) {
			n++
		}
	}
	return n
}

// Ex2: Truncate shortens s to at most max characters.
// If s already has max characters or fewer, return it unchanged.
// Otherwise return the first max characters followed by "...".
//
//	Truncate("Hello, Gopher", 5) == "Hello..."
//	Truncate("Go", 5)            == "Go"
//	Truncate("héllo wörld", 5)   == "héllo..."   ← count characters, not bytes (lesson 04!)
//
// Hint: runes := []rune(s); if len(runes) <= max { ... }; string(runes[:max]) + "..."
func Truncate(s string, max int) string {
	// TODO
	return ""
}
