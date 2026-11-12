# Lesson 29: Testing in depth

**Goal:** write the tests yourself this time: table-driven tests with subtests, fakes for things like
clocks and email, benchmarks, fuzz tests and examples. You'll also use your tests to find a real bug.

**Time:** about 2 hours.

```bash
cd ~/Desktop/learning-go/course/part4/29_testing
```

Until now, the tests were written for you. In this lesson, **you write the tests** for the code in
`shop.go`, in a file you create called `shop_test.go`. A checker (`checker_test.go`) tells you whether
each exercise is done.

---

## 1. Quick recap

From lesson 09: test files end in `_test.go`, test functions are `func TestXxx(t *testing.T)`, and
you report problems with `t.Errorf` (keep going) or `t.Fatalf` (stop this test). The **got / want**
style and **table-driven tests** are the Go conventions.

---

## 2. Subtests with `t.Run`

```go
func TestSlugify(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"simple", "Hello World", "hello-world"},
		{"extra spaces", "  Go   is fun ", "go-is-fun"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Slugify(tt.in); got != tt.want {
				t.Errorf("Slugify(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}
```

Each `t.Run` is a separate subtest with its own name:

```bash
go test -v -run TestSlugify               # all cases
go test -v -run 'TestSlugify/extra'       # just one case
```

A `t.Fatalf` inside a subtest only stops that subtest, and the other cases still run.

---

## 3. Helpers and cleanup

```go
func newTestDB(t *testing.T) *Store {
	t.Helper()                       // errors are reported at the CALLER's line, not here
	s, err := Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	t.Cleanup(func() { s.Close() })  // runs when the test finishes, like defer for the whole test
	return s
}
```

- **`t.Helper()`**: put it first in any test helper function, so failures point at the test that
  called it.
- **`t.Cleanup(f)`**: register cleanup from inside a helper. It runs after the test.
- **`t.TempDir()`**: a fresh folder, deleted afterwards (lesson 23).
- **`t.Parallel()`**: let this test run at the same time as other parallel tests.

---

## 4. Fakes: testing code that depends on the outside world

How do you test code that uses the current time, sends emails, or calls a payment API? You don't want
your tests to depend on the clock, spam real inboxes, or charge real cards.

The Go way: **the code depends on a small interface** (lesson 10), and the test passes in a **fake**.

```go
// in the real code
type Clock interface {
	Now() time.Time
}

func Greeting(c Clock) string {
	h := c.Now().Hour()
	...
}

// the real implementation, used by main
type RealClock struct{}
func (RealClock) Now() time.Time { return time.Now() }
```

```go
// in the test
type fakeClock struct {
	t time.Time
}

func (f fakeClock) Now() time.Time { return f.t }

func TestGreetingMorning(t *testing.T) {
	c := fakeClock{t: time.Date(2026, 1, 1, 9, 0, 0, 0, time.UTC)}   // always 9:00
	if got := Greeting(c); got != "Good morning" {
		t.Errorf("got %q", got)
	}
}
```

The fake is a few lines of normal Go: no mocking library needed. Python developers often reach for
`unittest.mock.patch` to replace things at runtime; in Go you design the code so dependencies are
**passed in**, and tests pass in fakes. This is called **dependency injection**, and it's one of the
main reasons Go code uses small interfaces.

A fake can also **record** what happened, so you can check it afterwards:

```go
type fakeMailer struct {
	sent []string            // who we sent to
	err  error               // set this to make Send fail
}

func (f *fakeMailer) Send(to, body string) error {
	if f.err != nil {
		return f.err
	}
	f.sent = append(f.sent, to)
	return nil
}
```

Notice the pointer receiver: `Send` changes the fake, so the test must pass `&fakeMailer{}`
(lesson 10, pointer receivers and interfaces).

---

## 5. Tests that find bugs

Writing tests isn't just checking code that already works. Thinking up cases like empty strings,
punctuation, repeated separators and unicode is how you **find** bugs. `shop.go` contains a real bug
in `Slugify`. Your table test in Ex1 should expose it, and then you fix it.

When a test fails, read the got / want carefully before touching the code. Sometimes the test is wrong!

---

## 6. Coverage

```bash
go test -cover                          # coverage: 87.5% of statements
go test -coverprofile=cover.out
go tool cover -html=cover.out           # opens a page: green = tested, red = never ran
```

Coverage shows which lines your tests **never run**. 100% isn't the goal, but red lines in important
logic are worth a test.

---

## 7. Benchmarks

A **benchmark** measures how fast code is:

```go
func BenchmarkSlugify(b *testing.B) {
	for b.Loop() {
		Slugify("Hello, World! This is Go")
	}
}
```

```bash
go test -bench=Slugify -benchmem
# BenchmarkSlugify-10    2459380    487.3 ns/op    96 B/op    3 allocs/op
```

It runs the loop as many times as needed for a stable measurement, then reports the time per operation
and the memory allocated. `b.Loop()` is the modern style (Go 1.24+). You'll see `for i := 0; i < b.N; i++`
in older code. Benchmarks don't run with plain `go test`, only with `-bench`.

Python comparison: `timeit`, built into the test tool.

---

## 8. Fuzzing

A **fuzz test** throws thousands of generated inputs at your function, looking for crashes or broken
rules. You describe a **property** that should always hold, instead of exact outputs:

```go
func FuzzSlugify(f *testing.F) {
	f.Add("Hello World")              // starting examples ("seed corpus")
	f.Add("  Go!  ")
	f.Fuzz(func(t *testing.T, s string) {
		got := Slugify(s)
		if strings.Contains(got, " ") {
			t.Errorf("Slugify(%q) = %q contains a space", s, got)
		}
		if got != strings.ToLower(got) {
			t.Errorf("Slugify(%q) = %q has upper case", s, got)
		}
	})
}
```

```bash
go test -run FuzzSlugify              # plain go test: only runs the seed examples
go test -fuzz=FuzzSlugify -fuzztime=10s   # generate new inputs for 10 seconds
```

If fuzzing finds a failing input, Go saves it in `testdata/fuzz/`, and it becomes a permanent test
case. Python comparison: the `hypothesis` library.

---

## 9. Examples: tests that are also documentation

```go
func ExampleSlugify() {
	fmt.Println(Slugify("Hello World"))
	// Output: hello-world
}
```

An example function prints something, and the `// Output:` comment says what it must print. `go test`
runs it and compares. It also shows up in `go doc` and on pkg.go.dev as documentation, which is how
standard library docs get their examples.

### 🤔 Predict the output

```go
func TestCases(t *testing.T) {
	for _, n := range []int{1, 2, 3} {
		t.Run(fmt.Sprint(n), func(t *testing.T) {
			if n == 2 {
				t.Fatalf("two is bad")
			}
			t.Log("ran", n)
		})
	}
}
```

With `go test -v`, which subtests pass?

<details><summary>Answer</summary>

`TestCases/1` and `TestCases/3` pass; `TestCases/2` fails. `t.Fatalf` only stops the subtest it's in,
so case 3 still runs. (The parent `TestCases` is reported as failed because a subtest failed.)
</details>

---

## 10. Common mistakes

| You wrote | What happens | Fix |
|---|---|---|
| Test depends on `time.Now()` | passes in the morning, fails at night | inject a `Clock`, pass a fake |
| A helper without `t.Helper()` | failures point at the helper, not the test | add `t.Helper()` |
| Real network/email in unit tests | slow, flaky, side effects | fakes behind small interfaces |
| Only testing the happy path | bugs hide in edge cases | empty, zero, unicode, errors, boundaries |
| Fake with value receiver that records calls | the recording is lost (it's on a copy) | pointer receiver, pass `&fake` |
| `t.Fatal` from a goroutine you started | not allowed: the test may hang or misreport | send results back, check in the test goroutine |

---

## 11. Exercises

Create **`shop_test.go`** in this folder (`package lesson29`) and write your tests there. Check your
progress with:

```bash
go test -v -run Ex          # the checker's verdict for each exercise
go test -v                  # everything, including your own tests
```

| | Exercise | What to write |
|---|---|---|
| Ex1 | `TestSlugify`: table-driven, with `t.Run`, at least 5 cases | Your cases should expose the bug in `Slugify`. Then **fix the bug** in `shop.go` |
| Ex2 | `fakeClock` + `TestGreeting` | test morning (before 12), afternoon (12–17), evening (18+), including the boundaries 12:00 and 18:00 |
| Ex3 | `fakeMailer` + `TestNotifyOverdue` | check who gets emailed, and what happens when sending fails |
| Ex4 | `BenchmarkSlugify` | use `b.Loop()`, run it with `go test -bench=. -benchmem` |
| Ex5 | `FuzzSlugify` | at least 2 seed inputs, check 2 properties; fuzz it for 10 seconds |
| Ex6 | `ExampleSlugify` | with an `// Output:` comment |

---

## ✅ Checklist

- [ ] I write table-driven tests with `t.Run` subtests, and use `t.Helper` and `t.Cleanup`
- [ ] I design code so dependencies (time, email, APIs) come in through small interfaces
- [ ] I can write a fake that returns fixed values or records calls
- [ ] I can measure coverage, write a benchmark, a fuzz test and an example
- [ ] I've used a test to find and fix a real bug

**Next → [Lesson 30: Shipping](../30_shipping/README.md)**
