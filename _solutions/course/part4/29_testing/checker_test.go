package lesson29

// This file checks your work. You don't need to edit it.
// Write YOUR tests in shop_test.go, then run:  go test -v -run Ex

import (
	"os"
	"regexp"
	"testing"
)

func yourTests(t *testing.T) string {
	t.Helper()
	data, err := os.ReadFile("shop_test.go")
	if err != nil {
		t.Fatal("shop_test.go doesn't exist yet. Create it in this folder, starting with: package lesson29")
	}
	return string(data)
}

func mustContain(t *testing.T, src, pattern, hint string) {
	t.Helper()
	if !regexp.MustCompile(pattern).MatchString(src) {
		t.Error(hint)
	}
}

func TestEx1Check(t *testing.T) {
	src := yourTests(t)
	mustContain(t, src, `func TestSlugify\(t \*testing\.T\)`, "Ex1: write func TestSlugify(t *testing.T) in shop_test.go")
	mustContain(t, src, `t\.Run\(`, "Ex1: use subtests: t.Run(tt.name, func(t *testing.T) { ... })")
	mustContain(t, src, `range`, "Ex1: make it table-driven (a slice of cases and a for ... range loop)")

	// Is the bug fixed?
	cases := map[string]string{
		"Hello, World!":      "hello-world",
		"  Go -- is   fun  ": "go-is-fun",
		"Café au lait":       "café-au-lait",
		"a!!b":               "a-b",
		"":                   "",
	}
	for in, want := range cases {
		if got := Slugify(in); got != want {
			t.Errorf("Ex1: the Slugify bug isn't fixed yet: Slugify(%q) = %q, want %q. Add a case like this to your table, then fix shop.go", in, got, want)
			break
		}
	}
}

func TestEx2Check(t *testing.T) {
	src := yourTests(t)
	mustContain(t, src, `type fakeClock struct`, "Ex2: define a fake clock: type fakeClock struct { ... }")
	mustContain(t, src, `func \(\w+ \*?fakeClock\) Now\(\) time\.Time`, "Ex2: give fakeClock a method: func (f fakeClock) Now() time.Time")
	mustContain(t, src, `func TestGreeting`, "Ex2: write func TestGreeting(t *testing.T)")
	mustContain(t, src, `Good afternoon`, "Ex2: test all three greetings, including \"Good afternoon\"")
	mustContain(t, src, `Good evening`, "Ex2: test all three greetings, including \"Good evening\"")
}

func TestEx3Check(t *testing.T) {
	src := yourTests(t)
	mustContain(t, src, `type fakeMailer struct`, "Ex3: define type fakeMailer struct { ... } that records who it sent to")
	mustContain(t, src, `func \(\w+ \*fakeMailer\) Send\(`, "Ex3: give fakeMailer a POINTER-receiver Send method (so it can record calls)")
	mustContain(t, src, `func TestNotifyOverdue`, "Ex3: write func TestNotifyOverdue(t *testing.T)")
	mustContain(t, src, `errors\.(New|Is)|fmt\.Errorf`, "Ex3: also test what happens when Send fails (make the fake return an error)")
}

func TestEx4Check(t *testing.T) {
	src := yourTests(t)
	mustContain(t, src, `func BenchmarkSlugify\(b \*testing\.B\)`, "Ex4: write func BenchmarkSlugify(b *testing.B)")
	mustContain(t, src, `b\.Loop\(\)`, "Ex4: use for b.Loop() { ... }")
}

func TestEx5Check(t *testing.T) {
	src := yourTests(t)
	mustContain(t, src, `func FuzzSlugify\(f \*testing\.F\)`, "Ex5: write func FuzzSlugify(f *testing.F)")
	if n := len(regexp.MustCompile(`f\.Add\(`).FindAllString(src, -1)); n < 2 {
		t.Error("Ex5: add at least 2 seed inputs with f.Add(...)")
	}
	mustContain(t, src, `f\.Fuzz\(`, "Ex5: call f.Fuzz(func(t *testing.T, s string) { ... })")
}

func TestEx6Check(t *testing.T) {
	src := yourTests(t)
	mustContain(t, src, `func ExampleSlugify\(\)`, "Ex6: write func ExampleSlugify()")
	mustContain(t, src, `// Output:`, "Ex6: end the example with a // Output: comment")
}
