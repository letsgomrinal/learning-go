// Lesson 11 exercises: interfaces in practice. Read README.md first!
//
// Check one exercise at a time:   go test -run Ex1 -v
//
// Add imports as you need them ("fmt", "io", "bufio", "strings", "bytes").
package lesson11

import "io"

// Ex1: Money is an amount in cents. Give it a String() method so that
// fmt.Sprint(Money(1250)) == "$12.50" and fmt.Sprint(Money(5)) == "$0.05".
//
// Hint: fmt.Sprintf("$%.2f", float64(m)/100). Convert first (README section 1, the ⚠️ part).
type Money int

// Ex2: Point should print as "(1, 2)". Add a String() method.
type Point struct {
	X, Y int
}

// Ex3: describe a value using a type switch (README section 4):
//
//	an int        → "int 42"
//	a string      → "string \"hi\"" (use %q for the quotes)
//	a bool        → "bool true"
//	a fmt.Stringer → "stringer " + its String()   (so Money(150) → "stringer $1.50", after Ex1)
//	nil           → "nil"
//	anything else → "other"
//
// Careful with order: put the concrete types (int, string, bool) before fmt.Stringer.
func Describe(v any) string {
	// TODO
	return ""
}

// Ex4: add up only the values that are ints. Skip everything else.
// SumInts([]any{1, "two", 3, 4.5, 5}) == 9
//
// Hint: n, ok := v.(int)
func SumInts(values []any) int {
	// TODO
	return 0
}

// Ex5: write one line per name to w, in the form "Hello, Ana!\n".
// WriteGreetings(&sb, []string{"Ana", "Bo"}) leaves "Hello, Ana!\nHello, Bo!\n" in sb.
//
// Hint: fmt.Fprintf(w, ...). This one function can then write to the terminal (os.Stdout),
// a file, or a strings.Builder in a test.
func WriteGreetings(w io.Writer, names []string) {
	// TODO
}

// Ex6: return the longest line readable from r (without its "\n").
// If there are no lines, return "". On a tie, return the first.
// LongestLine(strings.NewReader("go\ngopher\nhi")) == "gopher"
//
// Hint: bufio.NewScanner(r), then for sc.Scan() { line := sc.Text() ... }
// Compare lengths with len(line).
func LongestLine(r io.Reader) string {
	// TODO
	return ""
}

// Ex7: UpperWriter wraps another Writer, and makes all ASCII text written through it uppercase.
//
//	var sb strings.Builder
//	u := UpperWriter{W: &sb}
//	fmt.Fprintf(u, "hello %s", "go")    // sb now holds "HELLO GO"
//
// Write must:
//  1. uppercase the bytes:  upper := bytes.ToUpper(p)   (import "bytes")
//  2. write them to u.W
//  3. return len(p) and nil on success, or 0 and the error if u.W.Write failed
type UpperWriter struct {
	W io.Writer
}

func (u UpperWriter) Write(p []byte) (int, error) {
	// TODO
	return 0, nil
}
