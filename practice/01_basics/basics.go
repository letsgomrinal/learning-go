// Lesson 01: Basics. Speaking Go instead of Python-with-braces.
//
// Run:  go test -v ./practice/01_basics
//
// ─── CONCEPTS ──────────────────────────────────────────────────────────────
//
// 1. DECLARATIONS AND ZERO VALUES
//      var count int        // 0      Python has no equivalent: there is no "unset"
//      var name string      // ""
//      var ok bool          // false
//      var nums []int       // nil (but len(nums) == 0 and append works!)
//      x := 42              // short declaration, type inferred. Only inside functions.
//    Every variable ALWAYS has a value. Unused local variables and unused imports are
//    COMPILE ERRORS. That's on purpose.
//
// 2. MULTIPLE RETURN VALUES. Python returns a tuple; Go returns real multiple values:
//      func divmod(a, b int) (int, int) { return a / b, a % b }
//      q, r := divmod(7, 2)
//      _, r = divmod(7, 2)  // _ discards a value
//    The "comma ok" idiom is everywhere:
//      v, ok := m["key"]    // map lookup
//      n, err := strconv.Atoi("42")
//
// 3. ONE LOOP: for
//      for i := 0; i < 10; i++ {}     // C-style
//      for cond {}                     // while
//      for {}                          // while True
//      for i, v := range slice {}      // enumerate(slice)
//      for i := range 10 {}            // range(10) (Go 1.22+)
//      for k, v := range m {}          // m.items(), but in RANDOM order
//
// 4. SWITCH has no fallthrough by default, and cases can be expressions:
//      switch {
//      case n < 0: ...
//      case n == 0: ...
//      default: ...
//      }
//
// 5. STRINGS ARE BYTES. This is the big one for Python devs.
//      s := "héllo"
//      len(s)            // 6 BYTES, not 5 characters ('é' is 2 bytes in UTF-8)
//      s[1]              // a byte (uint8), NOT a character
//      []rune(s)         // convert to Unicode code points (like Python str indexing)
//      for i, r := range s {}  // range over a string yields RUNES (i = byte offset)
//    Strings are immutable. Build them with strings.Builder (like "".join()).
//
// 6. NO IMPLICIT CONVERSIONS
//      var i int = 3
//      var f float64 = i          // COMPILE ERROR
//      var f float64 = float64(i) // explicit
//    int + string doesn't work either. Use strconv.Itoa / fmt.Sprintf.
//
// 7. VARIADIC FUNCTIONS (like *args, but typed)
//      func sum(nums ...int) int  // nums is a []int
//      sum(1, 2, 3)
//      sum(slice...)              // like sum(*lst) in Python
//
// Useful docs:  go doc strings   go doc strconv   go doc unicode
// ───────────────────────────────────────────────────────────────────────────

package basics

// FizzBuzz returns a slice of n strings for 1..n:
// "Fizz" for multiples of 3, "Buzz" for multiples of 5, "FizzBuzz" for both,
// otherwise the number itself as a string.
//
// Hints: make([]string, 0, n) preallocates capacity. strconv.Itoa converts int → string.
// A tagless `switch { case ...: }` reads nicer than if/else chains.
func FizzBuzz(n int) []string {
	// TODO
	return nil
}

// ReverseRunes reverses a string by CHARACTERS (runes), not bytes.
// ReverseRunes("héllo, 世界") == "界世 ,olléh"
//
// Hint: reversing the bytes would corrupt multi-byte characters. Convert to []rune first,
// swap with the two-index trick:  r[i], r[j] = r[j], r[i]  (works just like Python).
func ReverseRunes(s string) string {
	// TODO
	return ""
}

// CountWords counts case-insensitive word frequencies.
// Words are separated by any whitespace. CountWords("Go go  GO\nstop") == {"go": 3, "stop": 1}
//
// Hints: strings.Fields splits on any whitespace (like Python's str.split() with no args).
// strings.ToLower. A map's missing keys return the zero value, so m[w]++ just works,
// no defaultdict needed!
func CountWords(s string) map[string]int {
	// TODO
	return nil
}

// Stats returns the min, max and average of the given numbers.
// If no numbers are given, return 0, 0, 0.
//
// Concepts: variadic params, NAMED RESULTS (minV, maxV, avg are already declared
// and zero-valued; a bare `return` returns them), and int → float64 conversion.
// Go 1.21+ also has builtin min() and max() functions.
func Stats(nums ...int) (minV, maxV int, avg float64) {
	// TODO
	return
}

// Divide returns a/b and true, or 0 and false if b is zero (the "comma ok" idiom).
func Divide(a, b int) (int, bool) {
	// TODO
	return 0, false
}

// Grade turns a score (0..100) into a letter:
// >=90 "A", >=80 "B", >=70 "C", >=60 "D", otherwise "F".
// Scores outside 0..100 return "invalid".
// Use a tagless switch.
func Grade(score int) string {
	// TODO
	return ""
}

// Weekday shows how Go does enums: a named type plus constants built with iota.
//
//	const (
//	    Sunday Weekday = iota // 0
//	    Monday                // 1  (iota auto-increments, the type repeats)
//	    ...
//	)
type Weekday int

const (
	Sunday Weekday = iota
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
)

// IsWeekend reports whether d is Saturday or Sunday.
// (This is a METHOD: a function with a receiver `d`. Lesson 03 covers methods properly.)
func (d Weekday) IsWeekend() bool {
	// TODO
	return false
}
