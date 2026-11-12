package lesson09

import (
	"os"
	"strings"
	"testing"

	"gomastery/course/part1/09_packages_tests/textutil"
)

// If this file doesn't compile ("undefined: textutil.WordCount"), do Ex1 first.

func TestEx1Exported(t *testing.T) {
	if got := textutil.WordCount("  go is   fun "); got != 3 {
		t.Errorf("Ex1: textutil.WordCount = %d, want 3", got)
	}
	if got := textutil.Capitalize("gopher"); got != "Gopher" {
		t.Errorf("Ex1: textutil.Capitalize(\"gopher\") = %q, want \"Gopher\"", got)
	}
}

func TestEx2Truncate(t *testing.T) {
	tests := []struct {
		in   string
		max  int
		want string
	}{
		{"Hello, Gopher", 5, "Hello..."},
		{"Go", 5, "Go"},
		{"exact", 5, "exact"},
		{"héllo wörld", 5, "héllo..."},
		{"", 3, ""},
	}
	for _, tt := range tests {
		if got := textutil.Truncate(tt.in, tt.max); got != tt.want {
			t.Errorf("Ex2: Truncate(%q, %d) = %q, want %q", tt.in, tt.max, got, tt.want)
		}
	}
}

func TestEx3YourOwnTest(t *testing.T) {
	data, err := os.ReadFile("textutil/textutil_test.go")
	if err != nil {
		t.Fatal("Ex3: textutil/textutil_test.go doesn't exist yet. Create it (README section 6, Ex3)")
	}
	src := string(data)
	if !strings.Contains(src, "func Test") {
		t.Error("Ex3: textutil_test.go has no test function. Tests are named func TestXxx(t *testing.T)")
	}
	if !strings.Contains(src, "Capitalize") {
		t.Error("Ex3: your test should call Capitalize")
	}
	if !strings.Contains(src, "range") {
		t.Error("Ex3: write it table-driven: a slice of test cases and a for ... range loop")
	}
	t.Log("Now run your test itself: go test -v ./textutil")
}
