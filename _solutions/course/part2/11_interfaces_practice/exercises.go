// Reference solution for lesson 11.
package lesson11

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
)

type Money int

func (m Money) String() string {
	return fmt.Sprintf("$%.2f", float64(m)/100)
}

type Point struct {
	X, Y int
}

func (p Point) String() string {
	return fmt.Sprintf("(%d, %d)", p.X, p.Y)
}

func Describe(v any) string {
	switch x := v.(type) {
	case nil:
		return "nil"
	case int:
		return fmt.Sprintf("int %d", x)
	case string:
		return fmt.Sprintf("string %q", x)
	case bool:
		return fmt.Sprintf("bool %t", x)
	case fmt.Stringer:
		return "stringer " + x.String()
	default:
		return "other"
	}
}

func SumInts(values []any) int {
	total := 0
	for _, v := range values {
		if n, ok := v.(int); ok {
			total += n
		}
	}
	return total
}

func WriteGreetings(w io.Writer, names []string) {
	for _, name := range names {
		fmt.Fprintf(w, "Hello, %s!\n", name)
	}
}

func LongestLine(r io.Reader) string {
	longest := ""
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		if line := sc.Text(); len(line) > len(longest) {
			longest = line
		}
	}
	return longest
}

type UpperWriter struct {
	W io.Writer
}

func (u UpperWriter) Write(p []byte) (int, error) {
	if _, err := u.W.Write(bytes.ToUpper(p)); err != nil {
		return 0, err
	}
	return len(p), nil
}
