package lesson36

import (
	"errors"
	"io"
	"io/fs"
	"strings"
	"testing"
	"testing/iotest"
	"time"
)

func TestEx1LimitedReader(t *testing.T) {
	r := &LimitedReader{R: strings.NewReader("hello, gopher"), N: 5}
	data, err := io.ReadAll(r)
	if err != nil || string(data) != "hello" {
		t.Errorf("Ex1: read %q, %v; want \"hello\", nil", data, err)
	}
	if r.N != 0 {
		t.Errorf("Ex1: N = %d after reading, want 0", r.N)
	}
	r = &LimitedReader{R: strings.NewReader("hi"), N: 10}
	if data, _ := io.ReadAll(r); string(data) != "hi" {
		t.Errorf("Ex1: limit bigger than the data: read %q, want \"hi\"", data)
	}
	// iotest.TestReader checks a Reader follows all the io.Reader rules.
	if err := iotest.TestReader(&LimitedReader{R: strings.NewReader("abcdefghij"), N: 6}, []byte("abcdef")); err != nil {
		t.Errorf("Ex1: iotest.TestReader: %v", err)
	}
}

type failingWriter struct{}

func (failingWriter) Write(p []byte) (int, error) { return 0, errors.New("disk full") }

type shortWriter struct{}

func (shortWriter) Write(p []byte) (int, error) { return len(p) / 2, nil }

func TestEx2MultiWriter(t *testing.T) {
	var a, b, c strings.Builder
	w := MultiWriter(&a, &b, &c)
	n, err := io.WriteString(w, "copy me")
	if n != 7 || err != nil {
		t.Errorf("Ex2: Write returned %d, %v; want 7, nil", n, err)
	}
	if a.String() != "copy me" || b.String() != "copy me" || c.String() != "copy me" {
		t.Errorf("Ex2: writers got %q, %q, %q; want all \"copy me\"", a.String(), b.String(), c.String())
	}
	var d strings.Builder
	if _, err := io.WriteString(MultiWriter(&d, failingWriter{}), "x"); err == nil || err.Error() != "disk full" {
		t.Errorf("Ex2: with a failing writer, err = %v, want \"disk full\"", err)
	}
	if _, err := io.WriteString(MultiWriter(shortWriter{}), "abcd"); !errors.Is(err, io.ErrShortWrite) {
		t.Errorf("Ex2: with a short writer, err = %v, want io.ErrShortWrite", err)
	}
}

func TestEx3ProcessorFunc(t *testing.T) {
	var p any = ProcessorFunc(strings.ToUpper)
	upper, ok := p.(Processor)
	if !ok {
		t.Fatal("Ex3: ProcessorFunc doesn't implement Processor yet. Add the Process method")
	}
	trim := any(ProcessorFunc(strings.TrimSpace)).(Processor)
	exclaim := any(ProcessorFunc(func(s string) string { return s + "!" })).(Processor)
	if got := Pipeline("  hello ", trim, upper, exclaim); got != "HELLO!" {
		t.Errorf("Ex3: Pipeline = %q, want \"HELLO!\"", got)
	}
}

func TestEx4QueryError(t *testing.T) {
	var e any = &QueryError{}
	if _, ok := e.(error); !ok {
		t.Fatal("Ex4: *QueryError isn't an error yet. Add Error() string")
	}
	if _, ok := e.(interface{ Unwrap() error }); !ok {
		t.Fatal("Ex4: *QueryError has no Unwrap() error method yet")
	}
	err := any(&QueryError{Query: "SELECT * FROM users", Err: fs.ErrPermission}).(error)
	if got, want := err.Error(), `query "SELECT * FROM users": permission denied`; got != want {
		t.Errorf("Ex4: Error() = %q, want %q", got, want)
	}
	wrapped := errors.Join(errors.New("handler failed"), err)
	if !errors.Is(wrapped, fs.ErrPermission) {
		t.Error("Ex4: errors.Is can't see through QueryError. Unwrap must return Err")
	}
	// (Written with an interface target so this file compiles before your methods exist.)
	var target interface {
		Error() string
		Unwrap() error
	}
	qe, _ := any(target).(*QueryError)
	if errors.As(wrapped, &target) {
		qe, _ = any(target).(*QueryError)
	}
	if qe == nil || qe.Query != "SELECT * FROM users" {
		t.Error("Ex4: errors.As should find the *QueryError")
	}
}

func TestEx5Options(t *testing.T) {
	c := NewClient("https://api.example.com")
	want := Client{BaseURL: "https://api.example.com", Timeout: 10 * time.Second, Retries: 3, UserAgent: "gomastery/1.0"}
	if *c != want {
		t.Errorf("Ex5: defaults = %+v, want %+v", *c, want)
	}
	c = NewClient("u", WithTimeout(2*time.Second), WithRetries(0), WithUserAgent("bot/2"))
	want = Client{BaseURL: "u", Timeout: 2 * time.Second, Retries: 0, UserAgent: "bot/2"}
	if *c != want {
		t.Errorf("Ex5: with options = %+v, want %+v", *c, want)
	}
	if c := NewClient("u", WithRetries(-5)); c.Retries != 0 {
		t.Errorf("Ex5: WithRetries(-5) gave %d retries, want 0", c.Retries)
	}
}

func TestEx6SortPeople(t *testing.T) {
	people := []Person{
		{"Zoe", "Adams", 30},
		{"Ana", "Silva", 25},
		{"Ana", "Adams", 40},
		{"Ana", "Adams", 60},
		{"Bo", "Silva", 33},
	}
	SortPeople(people)
	want := []Person{
		{"Ana", "Adams", 60},
		{"Ana", "Adams", 40},
		{"Zoe", "Adams", 30},
		{"Ana", "Silva", 25},
		{"Bo", "Silva", 33},
	}
	for i := range want {
		if people[i] != want[i] {
			t.Fatalf("Ex6: sorted = %v\n   want %v", people, want)
		}
	}
}
