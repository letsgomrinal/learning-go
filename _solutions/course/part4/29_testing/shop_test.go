// Reference solution for lesson 29: the tests YOU were asked to write.
package lesson29

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"
)

// Ex1
func TestSlugify(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"simple", "Hello World", "hello-world"},
		{"punctuation", "Hello, World!", "hello-world"},
		{"repeated separators", "  Go -- is   fun  ", "go-is-fun"}, // this case found the bug
		{"unicode letters", "Café au lait", "café-au-lait"},
		{"digits", "Top 10 tips", "top-10-tips"},
		{"empty", "", ""},
		{"only punctuation", "!!!", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Slugify(tt.in); got != tt.want {
				t.Errorf("Slugify(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

// Ex2
type fakeClock struct {
	t time.Time
}

func (f fakeClock) Now() time.Time { return f.t }

func at(hour, min int) fakeClock {
	return fakeClock{t: time.Date(2026, 1, 1, hour, min, 0, 0, time.UTC)}
}

func TestGreeting(t *testing.T) {
	tests := []struct {
		clock fakeClock
		want  string
	}{
		{at(0, 0), "Good morning"},
		{at(11, 59), "Good morning"},
		{at(12, 0), "Good afternoon"},
		{at(17, 59), "Good afternoon"},
		{at(18, 0), "Good evening"},
		{at(23, 30), "Good evening"},
	}
	for _, tt := range tests {
		t.Run(tt.clock.t.Format("15:04"), func(t *testing.T) {
			if got := Greeting(tt.clock); got != tt.want {
				t.Errorf("at %s: got %q, want %q", tt.clock.t.Format("15:04"), got, tt.want)
			}
		})
	}
}

// Ex3
type fakeMailer struct {
	sent   []string
	failOn string // Send fails for this address
}

func (f *fakeMailer) Send(to, body string) error {
	if to == f.failOn {
		return errors.New("smtp down")
	}
	f.sent = append(f.sent, to)
	return nil
}

func TestNotifyOverdue(t *testing.T) {
	now := time.Date(2026, 3, 10, 12, 0, 0, 0, time.UTC)
	invoices := []Invoice{
		{Customer: "late@x.com", Amount: 1250, Due: now.Add(-48 * time.Hour)},
		{Customer: "paid@x.com", Amount: 500, Due: now.Add(-48 * time.Hour), Paid: true},
		{Customer: "future@x.com", Amount: 500, Due: now.Add(24 * time.Hour)},
		{Customer: "exact@x.com", Amount: 500, Due: now}, // due exactly now: not overdue yet
		{Customer: "late2@x.com", Amount: 99, Due: now.Add(-time.Hour)},
	}

	t.Run("sends to overdue customers only", func(t *testing.T) {
		m := &fakeMailer{}
		n, err := NotifyOverdue(m, invoices, now)
		if err != nil || n != 2 {
			t.Fatalf("got %d, %v; want 2, nil", n, err)
		}
		if strings.Join(m.sent, ",") != "late@x.com,late2@x.com" {
			t.Errorf("sent to %v", m.sent)
		}
	})

	t.Run("stops when sending fails", func(t *testing.T) {
		m := &fakeMailer{failOn: "late2@x.com"}
		n, err := NotifyOverdue(m, invoices, now)
		if n != 1 || err == nil || !strings.Contains(err.Error(), "notify late2@x.com") {
			t.Errorf("got %d, %v; want 1 and an error mentioning late2@x.com", n, err)
		}
	})
}

// Ex4
func BenchmarkSlugify(b *testing.B) {
	for b.Loop() {
		Slugify("Hello, World! This is a fairly long title for Go")
	}
}

// Ex5
func FuzzSlugify(f *testing.F) {
	f.Add("Hello World")
	f.Add("  Go -- is fun!  ")
	f.Add("Café")
	f.Fuzz(func(t *testing.T, s string) {
		got := Slugify(s)
		if strings.Contains(got, " ") {
			t.Errorf("Slugify(%q) = %q contains a space", s, got)
		}
		if strings.Contains(got, "--") {
			t.Errorf("Slugify(%q) = %q contains a double dash", s, got)
		}
		if strings.HasPrefix(got, "-") || strings.HasSuffix(got, "-") {
			t.Errorf("Slugify(%q) = %q starts or ends with a dash", s, got)
		}
	})
}

// Ex6
func ExampleSlugify() {
	fmt.Println(Slugify("Hello, World!"))
	// Output: hello-world
}
