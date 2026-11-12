package interfaces

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"strings"
	"testing"
)

// These tests use `any` and type assertions so the file compiles before you add
// your methods. Once you're done, all of them should pass.

func TestShapes(t *testing.T) {
	var r any = Rect{3, 4}
	var c any = Circle{1}

	rs, ok := r.(Shape)
	if !ok {
		t.Fatal("Rect does not implement Shape yet (add Area and Perimeter)")
	}
	cs, ok := c.(Shape)
	if !ok {
		t.Fatal("Circle does not implement Shape yet")
	}
	if rs.Area() != 12 || rs.Perimeter() != 14 {
		t.Errorf("Rect{3,4}: area %v perim %v, want 12 14", rs.Area(), rs.Perimeter())
	}
	if math.Abs(cs.Area()-math.Pi) > 1e-9 || math.Abs(cs.Perimeter()-2*math.Pi) > 1e-9 {
		t.Errorf("Circle{1}: area %v perim %v", cs.Area(), cs.Perimeter())
	}

	got := TotalArea(Rect{3, 4}, Circle{1}, Rect{1, 1})
	if want := 13 + math.Pi; math.Abs(got-want) > 1e-9 {
		t.Errorf("TotalArea = %v, want %v", got, want)
	}
}

func TestTemperatureStringer(t *testing.T) {
	var v any = Temperature(21.5)
	if _, ok := v.(fmt.Stringer); !ok {
		t.Fatal("Temperature does not implement fmt.Stringer yet")
	}
	if got := fmt.Sprint(Temperature(21.5)); got != "21.5°C" {
		t.Errorf("Sprint = %q, want %q", got, "21.5°C")
	}
}

func TestDescribe(t *testing.T) {
	tests := []struct {
		in   any
		want string
	}{
		{42, "int: 42"},
		{"hello", "string of length 5"},
		{[]int{1, 2, 3}, "slice of 3 ints"},
		{Rect{3, 4}, "shape with area 12.00"},
		{nil, "nil"},
		{3.14, "unknown"},
	}
	for _, tt := range tests {
		if got := Describe(tt.in); got != tt.want {
			t.Errorf("Describe(%#v) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestCountLines(t *testing.T) {
	n, err := CountLines(strings.NewReader("a\nb\nc\n"))
	if err != nil || n != 3 {
		t.Errorf("CountLines = (%d, %v), want (3, nil)", n, err)
	}
	n, _ = CountLines(strings.NewReader("no newline at end"))
	if n != 1 {
		t.Errorf("CountLines(single) = %d, want 1", n)
	}
}

func TestUpperWriter(t *testing.T) {
	var buf bytes.Buffer
	w := UpperWriter{W: &buf}
	var _ io.Writer = w // compile-time check: UpperWriter is an io.Writer

	fmt.Fprintf(w, "hello %s! ", "go")
	io.WriteString(w, "bye")
	if got := buf.String(); got != "HELLO GO! BYE" {
		t.Errorf("got %q, want %q", got, "HELLO GO! BYE")
	}
}

func TestCountingReader(t *testing.T) {
	cr := &CountingReader{R: strings.NewReader("hello, world")}
	var _ io.Reader = cr

	// io.ReadAll calls Read repeatedly until io.EOF.
	data, err := io.ReadAll(cr)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "hello, world" || cr.Count != 12 {
		t.Errorf("read %q, count %d; want %q, 12", data, cr.Count, "hello, world")
	}
}
