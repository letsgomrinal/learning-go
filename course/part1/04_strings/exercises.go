// Lesson 04 exercises: strings. Read README.md first!
//
// Check one exercise at a time:   go test -run Ex1 -v
//
// You'll need to add imports at the top as you go, for example:
//
//	import (
//		"strconv"
//		"strings"
//	)
//
// Remember: an import you don't use yet is a compile error, so add each one when you need it.
// (If your editor has the Go extension, it adds imports for you when you save.)
package lesson04

// Ex1: remove spaces around the text, make it uppercase, and add "!" at the end.
// Shout("  hello  ") == "HELLO!"
func Shout(s string) string {
	// TODO
	return ""
}

// Ex2: return true if the text, after trimming spaces, ends with "?".
// IsQuestion("how are you?  ") == true, IsQuestion("fine.") == false
func IsQuestion(s string) bool {
	// TODO
	return false
}

// Ex3: return the uppercase first letter of each word.
// Initials("ada lovelace") == "AL", Initials("  grace   brewster hopper ") == "GBH"
// You can assume the words are plain English letters (one byte each).
//
// Hints:
//   - strings.Fields(name) splits into words and ignores extra spaces
//   - word[:1] is a string holding the first byte of word
//   - strings.ToUpper works on that piece
func Initials(name string) string {
	// TODO
	return ""
}

// Ex4: count the vowels (a, e, i, o, u, upper or lower case).
// CountVowels("Hello World") == 3
//
// Hint: loop with  for _, r := range strings.ToLower(s)  and use a switch:
//
//	switch r {
//	case 'a', 'e', 'i', 'o', 'u':   ← single quotes: these are runes
func CountVowels(s string) int {
	// TODO
	return 0
}

// Ex5: convert both strings to numbers and add them.
// Return the sum and true. If either string isn't a number, return 0 and false.
// AddStrings("2", "40") == 42, true
// AddStrings("2", "forty") == 0, false
//
// Hint: strconv.Atoi returns (number, err). Check err != nil after EACH conversion.
func AddStrings(a, b string) (int, bool) {
	// TODO
	return 0, false
}

// Ex6: return the number of CHARACTERS (not bytes).
// CharCount("héllo") == 5, CharCount("日本語") == 3
func CharCount(s string) int {
	// TODO
	return 0
}

// Ex7: reverse the string by characters.
// Reverse("héllo") == "olléh", Reverse("ab🙂") == "🙂ba"
//
// Hint:
//  1. runes := []rune(s)
//  2. build a new []rune going backwards, OR swap from both ends:
//     for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 { runes[i], runes[j] = runes[j], runes[i] }
//  3. return string(runes)
func Reverse(s string) string {
	// TODO
	return ""
}

// Ex8: return true if the text reads the same backwards, ignoring upper/lower case and spaces.
// IsPalindrome("Never odd or even") == true, IsPalindrome("Gopher") == false
//
// Hint: lowercase it, remove spaces (strings.ReplaceAll(s, " ", "")), then compare it with
// Reverse(...) of itself. You can call your own function from Ex7!
func IsPalindrome(s string) bool {
	// TODO
	return false
}
