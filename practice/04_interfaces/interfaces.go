// Lesson 04: Interfaces. Duck typing, checked by the compiler.
//
// Run:  go test -v ./practice/04_interfaces
//
// ─── CONCEPTS ──────────────────────────────────────────────────────────────
//
// 1. IMPLICIT SATISFACTION. There's no `implements` keyword:
//      type Shape interface { Area() float64 }
//      type Square struct{ S float64 }
//      func (s Square) Area() float64 { return s.S * s.S }
//      var sh Shape = Square{2}    // OK: Square has Area(), so it IS a Shape
//    It's Python duck typing, but the compiler checks it. A common trick to assert
//    at compile time:   var _ Shape = Square{}
//
// 2. SMALL INTERFACES WIN. The most important interfaces in Go have ONE method:
//      type Reader interface { Read(p []byte) (n int, err error) }   // io.Reader
//      type Writer interface { Write(p []byte) (n int, err error) }  // io.Writer
//      type Stringer interface { String() string }                   // fmt.Stringer (like __str__)
//    Files, network connections, HTTP bodies, gzip streams and buffers are all
//    Readers/Writers, so code written against io.Reader works with all of them.
//
// 3. "ACCEPT INTERFACES, RETURN STRUCTS"
//    Parameters: take the smallest interface you need (io.Reader, not *os.File).
//    Return values: return concrete types. Define interfaces where they're CONSUMED,
//    not where the types are implemented (the opposite of Java/ABCs).
//
// 4. any (= interface{}) holds any value. You get the concrete type back with:
//      s, ok := v.(string)            // type assertion (ok idiom: no panic)
//      switch x := v.(type) {         // type switch
//      case int:    ...x is an int here
//      case string: ...x is a string here
//      }
//
// 5. METHOD SETS: if the methods have POINTER receivers, only *T satisfies the
//    interface, not T. "T does not implement X (method has pointer receiver)" is
//    an error you'll see often.
//
// 6. An interface value is a pair (type, value). That detail matters a lot for nil
//    checks; lesson 05 shows the trap.
// ───────────────────────────────────────────────────────────────────────────

package interfaces

import (
	"io"
)

// ── Part 1: polymorphism ─────────────────────────────────────────────────

type Shape interface {
	Area() float64
	Perimeter() float64
}

// TODO: make Rect and Circle satisfy Shape by adding Area() and Perimeter() methods.
// Use math.Pi for the circle. Uncomment the compile-time checks once you're done:
//
//	var _ Shape = Rect{}
//	var _ Shape = Circle{}

type Rect struct{ W, H float64 }

type Circle struct{ R float64 }

// TotalArea sums the area of any number of shapes.
// (Once Rect and Circle have methods, change the parameter type to ...Shape.)
func TotalArea(shapes ...any) float64 {
	// TODO: change signature to (shapes ...Shape) and implement
	return 0
}

// ── Part 2: fmt.Stringer (Go's __str__) ──────────────────────────────────

type Temperature float64 // Celsius

// TODO: add a String() method to Temperature so that
// fmt.Sprint(Temperature(21.5)) == "21.5°C"
// Hint: fmt.Sprintf("%.1f°C", float64(t)). Convert to float64 first, otherwise
// Sprintf calls String() again and recurses forever!

// ── Part 3: type switches ────────────────────────────────────────────────

// Describe returns a description of v based on its dynamic type:
//
//	int      → "int: 42"
//	string   → "string of length 5"
//	[]int    → "slice of 3 ints"
//	Shape    → "shape with area 12.00"  (any Shape, use %.2f)
//	nil      → "nil"
//	anything else → "unknown"
func Describe(v any) string {
	// TODO: use switch x := v.(type) { ... }
	return ""
}

// ── Part 4: io.Reader and io.Writer ──────────────────────────────────────

// CountLines counts the lines readable from r. It works with files, strings,
// network connections... anything that's an io.Reader.
//
// Hint: bufio.NewScanner(r); for sc.Scan() { n++ }; return n, sc.Err()
func CountLines(r io.Reader) (int, error) {
	// TODO
	return 0, nil
}

// UpperWriter wraps another io.Writer and upper-cases ASCII text written through it.
// That's the decorator pattern with interfaces, like Python file-like objects.
//
//	var buf bytes.Buffer
//	w := UpperWriter{W: &buf}
//	fmt.Fprintf(w, "hello %s", "go")  // buf now holds "HELLO GO"
type UpperWriter struct {
	W io.Writer
}

// Write implements io.Writer.
// Hint: bytes.ToUpper(p), then write to u.W. Return len(p) on success (the number of
// INPUT bytes consumed), and pass along any error from the inner writer.
func (u UpperWriter) Write(p []byte) (int, error) {
	// TODO
	return 0, nil
}

// CountingReader wraps an io.Reader and counts how many bytes have been read through it.
type CountingReader struct {
	R     io.Reader
	Count int
}

// Read implements io.Reader. It needs a POINTER receiver because it updates Count.
// Hint: n, err := c.R.Read(p); c.Count += n; return n, err
func (c *CountingReader) Read(p []byte) (int, error) {
	// TODO
	return 0, io.EOF
}
