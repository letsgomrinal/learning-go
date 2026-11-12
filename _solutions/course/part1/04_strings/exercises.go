// Reference solution for lesson 04.
package lesson04

import (
	"strconv"
	"strings"
)

func Shout(s string) string {
	return strings.ToUpper(strings.TrimSpace(s)) + "!"
}

func IsQuestion(s string) bool {
	return strings.HasSuffix(strings.TrimSpace(s), "?")
}

func Initials(name string) string {
	result := ""
	for _, word := range strings.Fields(name) {
		result += strings.ToUpper(word[:1])
	}
	return result
}

func CountVowels(s string) int {
	count := 0
	for _, r := range strings.ToLower(s) {
		switch r {
		case 'a', 'e', 'i', 'o', 'u':
			count++
		}
	}
	return count
}

func AddStrings(a, b string) (int, bool) {
	x, err := strconv.Atoi(a)
	if err != nil {
		return 0, false
	}
	y, err := strconv.Atoi(b)
	if err != nil {
		return 0, false
	}
	return x + y, true
}

func CharCount(s string) int {
	return len([]rune(s))
}

func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func IsPalindrome(s string) bool {
	cleaned := strings.ReplaceAll(strings.ToLower(s), " ", "")
	return cleaned == Reverse(cleaned)
}
