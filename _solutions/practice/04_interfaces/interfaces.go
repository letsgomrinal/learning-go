// Reference solution for lesson 04.
package interfaces

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math"
)

type Shape interface {
	Area() float64
	Perimeter() float64
}

var (
	_ Shape = Rect{}
	_ Shape = Circle{}
)

type Rect struct{ W, H float64 }

func (r Rect) Area() float64      { return r.W * r.H }
func (r Rect) Perimeter() float64 { return 2 * (r.W + r.H) }

type Circle struct{ R float64 }

func (c Circle) Area() float64      { return math.Pi * c.R * c.R }
func (c Circle) Perimeter() float64 { return 2 * math.Pi * c.R }

func TotalArea(shapes ...Shape) float64 {
	total := 0.0
	for _, s := range shapes {
		total += s.Area()
	}
	return total
}

type Temperature float64

func (t Temperature) String() string {
	return fmt.Sprintf("%.1f°C", float64(t))
}

func Describe(v any) string {
	switch x := v.(type) {
	case nil:
		return "nil"
	case int:
		return fmt.Sprintf("int: %d", x)
	case string:
		return fmt.Sprintf("string of length %d", len(x))
	case []int:
		return fmt.Sprintf("slice of %d ints", len(x))
	case Shape: // an interface can be a case too
		return fmt.Sprintf("shape with area %.2f", x.Area())
	default:
		return "unknown"
	}
}

func CountLines(r io.Reader) (int, error) {
	sc := bufio.NewScanner(r)
	n := 0
	for sc.Scan() {
		n++
	}
	return n, sc.Err()
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

type CountingReader struct {
	R     io.Reader
	Count int
}

func (c *CountingReader) Read(p []byte) (int, error) {
	n, err := c.R.Read(p)
	c.Count += n
	return n, err
}
