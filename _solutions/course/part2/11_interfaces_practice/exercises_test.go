package lesson11

import (
	"fmt"
	"io"
	"strings"
	"testing"
)

func TestEx1Money(t *testing.T) {
	var v any = Money(0)
	if _, ok := v.(fmt.Stringer); !ok {
		t.Fatal("Ex1: Money has no String() method yet")
	}
	tests := map[Money]string{1250: "$12.50", 5: "$0.05", 100: "$1.00", 0: "$0.00"}
	for m, want := range tests {
		if got := fmt.Sprint(m); got != want {
			t.Errorf("Ex1: fmt.Sprint(Money(%d)) = %q, want %q", int(m), got, want)
		}
	}
}

func TestEx2Point(t *testing.T) {
	var v any = Point{}
	if _, ok := v.(fmt.Stringer); !ok {
		t.Fatal("Ex2: Point has no String() method yet")
	}
	if got := fmt.Sprint(Point{X: 1, Y: 2}); got != "(1, 2)" {
		t.Errorf("Ex2: fmt.Sprint(Point{1, 2}) = %q, want %q", got, "(1, 2)")
	}
	if got := fmt.Sprintf("from %v to %v", Point{0, 0}, Point{-3, 4}); got != "from (0, 0) to (-3, 4)" {
		t.Errorf("Ex2: got %q", got)
	}
}

func TestEx3Describe(t *testing.T) {
	tests := []struct {
		in   any
		want string
	}{
		{42, "int 42"},
		{"hi", `string "hi"`},
		{true, "bool true"},
		{nil, "nil"},
		{3.5, "other"},
		{[]int{1}, "other"},
	}
	for _, tt := range tests {
		if got := Describe(tt.in); got != tt.want {
			t.Errorf("Ex3: Describe(%#v) = %q, want %q", tt.in, got, tt.want)
		}
	}
	var m any = Money(150)
	if _, ok := m.(fmt.Stringer); ok {
		if got := Describe(Money(150)); got != "stringer $1.50" {
			t.Errorf("Ex3: Describe(Money(150)) = %q, want %q", got, "stringer $1.50")
		}
	} else {
		t.Log("Ex3: (the Stringer case is checked once Ex1 is done)")
	}
}

func TestEx4SumInts(t *testing.T) {
	if got := SumInts([]any{1, "two", 3, 4.5, 5}); got != 9 {
		t.Errorf("Ex4: SumInts([1 two 3 4.5 5]) = %d, want 9", got)
	}
	if got := SumInts(nil); got != 0 {
		t.Errorf("Ex4: SumInts(empty) = %d, want 0", got)
	}
}

func TestEx5WriteGreetings(t *testing.T) {
	var sb strings.Builder
	WriteGreetings(&sb, []string{"Ana", "Bo"})
	if got, want := sb.String(), "Hello, Ana!\nHello, Bo!\n"; got != want {
		t.Errorf("Ex5: wrote %q, want %q", got, want)
	}
}

func TestEx6LongestLine(t *testing.T) {
	tests := map[string]string{
		"go\ngopher\nhi":    "gopher",
		"same\nsize\n":      "same",
		"":                  "",
		"one line no end":   "one line no end",
		"a\n\nlongest\nb\n": "longest",
	}
	for in, want := range tests {
		if got := LongestLine(strings.NewReader(in)); got != want {
			t.Errorf("Ex6: LongestLine(%q) = %q, want %q", in, got, want)
		}
	}
}

type failWriter struct{}

func (failWriter) Write(p []byte) (int, error) { return 0, io.ErrClosedPipe }

func TestEx7UpperWriter(t *testing.T) {
	var sb strings.Builder
	u := UpperWriter{W: &sb}
	fmt.Fprintf(u, "hello %s! ", "go")
	io.WriteString(u, "bye")
	if got := sb.String(); got != "HELLO GO! BYE" {
		t.Errorf("Ex7: wrote %q, want %q", got, "HELLO GO! BYE")
	}
	n, err := UpperWriter{W: &sb}.Write([]byte("abc"))
	if n != 3 || err != nil {
		t.Errorf("Ex7: Write(\"abc\") returned %d, %v; want 3, nil", n, err)
	}
	if _, err := (UpperWriter{W: failWriter{}}).Write([]byte("x")); err == nil {
		t.Errorf("Ex7: when the inner writer fails, Write must return its error")
	}
}
