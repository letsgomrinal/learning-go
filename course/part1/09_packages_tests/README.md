# Lesson 09: Packages and tests

**Goal:** split code into packages, understand public vs private names, and write your own tests,
which is how you've been checking every exercise so far.

**Time:** 1.5 hours.

```bash
cd ~/Desktop/learning-go/course/part1/09_packages_tests
```

---

## 1. Modules and packages

Two words you'll see everywhere:

- A **package** is a folder of `.go` files that belong together. Every file in the folder starts
  with the same `package name` line. (Python: a module or package.)
- A **module** is a whole project: a tree of packages with a `go.mod` file at the top.
  (Python: roughly a project with a `pyproject.toml`.)

Look at the `go.mod` file at the top of this course (`~/Desktop/learning-go/go.mod`):

```
module gomastery

go 1.24
```

`gomastery` is this module's name. Every package inside it is imported by **module name + folder path**.
This lesson's `textutil` folder is the package:

```
gomastery/course/part1/09_packages_tests/textutil
```

For your own new projects, you create a module with:

```bash
mkdir myproject && cd myproject
go mod init myproject        # creates go.mod
```

(Projects that will live on GitHub use the repo path as the name, like `github.com/yourname/myproject`.)

---

## 2. Public and private: capital letters

This is the simplest rule in Go, and you've been seeing it since lesson 00:

> **A name that starts with a capital letter is exported (public). Lowercase is private to its package.**

```go
package textutil

func WordCount(s string) int { ... }   // exported: other packages can call textutil.WordCount
func helper() { ... }                  // private: only code in package textutil can use it
```

That's why it's `fmt.Println` and `strings.ToUpper`: those functions are exported from their packages.
It works the same for types, struct fields, methods, constants and variables:

```go
type Book struct {
	Title string    // exported field
	price int       // private field: other packages can't read or set it
}
```

Python has a convention (`_private`) that anyone can ignore. Go's rule is enforced by the compiler.

Private fields are how you protect data. In lesson 08's `Counter`, the `count` field was lowercase,
so code in other packages can only change it through `Inc`, `Add` and `Reset`.

---

## 3. Using a package

Open `textutil/textutil.go`. It's a small package with text helpers. To use it from another package,
you import its full path and call its exported names with the package name in front:

```go
package main

import (
	"fmt"

	"gomastery/course/part1/09_packages_tests/textutil"
)

func main() {
	fmt.Println(textutil.WordCount("go is fun"))
}
```

Standard library imports (`"fmt"`) come first, then your own, separated by a blank line. `gofmt` keeps
them sorted.

Right now this **won't compile**, because the functions in `textutil` are lowercase. Fixing that is
your first exercise.

### Multiple files, one package

A package can be split over many files in its folder. They all share everything, private names
included, as if they were one big file. There's no need to import between files of the same package.

### Documentation comments

A comment directly above an exported name is its documentation. Start it with the name:

```go
// WordCount returns the number of words in s, split on whitespace.
func WordCount(s string) int {
```

`go doc` shows it:

```bash
go doc ./textutil
go doc ./textutil.WordCount
```

---

## 4. Third-party packages (for later)

To use someone else's package, you `go get` it, then import it:

```bash
go get github.com/google/uuid
```

```go
import "github.com/google/uuid"

id := uuid.New()
```

`go get` records the dependency in `go.mod` (like `pip install` + `requirements.txt`), and `go.sum`
stores checksums so every build uses exactly the same code. You don't need any third-party packages
in this course yet. Go's standard library covers most things.

---

## 5. How Go tests work

You've been running tests since lesson 00. Now let's write them.

The rules:
1. Tests live in files ending in **`_test.go`**, next to the code they test. `go build` ignores
   these files; only `go test` uses them.
2. A test is a function named **`TestSomething`** that takes **`t *testing.T`**.
3. When something is wrong, report it with **`t.Errorf`** (keeps going) or **`t.Fatalf`** (stops this test).
4. A test that reports nothing passes.

```go
package textutil

import "testing"

func TestWordCount(t *testing.T) {
	got := WordCount("go is fun")
	want := 3
	if got != want {
		t.Errorf("WordCount(\"go is fun\") = %d, want %d", got, want)
	}
}
```

There's no `assert` like in pytest. You write a normal `if` and describe what went wrong. The
**got / want** wording is the Go convention.

Notice the test file uses `package textutil`, the same package as the code. That means tests can
also check private functions.

### Table-driven tests

To check many inputs, put them in a slice of test cases and loop. This is **the** standard Go test
style, and every test in this course uses it. Open any `exercises_test.go` and you'll recognise it now:

```go
func TestWordCount(t *testing.T) {
	tests := []struct {
		input string
		want  int
	}{
		{"go is fun", 3},
		{"", 0},
		{"   spaces   everywhere  ", 2},
	}
	for _, tt := range tests {
		got := WordCount(tt.input)
		if got != tt.want {
			t.Errorf("WordCount(%q) = %d, want %d", tt.input, got, tt.want)
		}
	}
}
```

`[]struct{ input string; want int }` is a slice of an **anonymous struct**: a struct type with no
name, defined right where it's used. (Python equivalent: `@pytest.mark.parametrize`.)

### Subtests with `t.Run`

`t.Run` gives each case its own name, so you can see and run them separately:

```go
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := WordCount(tt.input); got != tt.want {
				t.Errorf("got %d, want %d", got, tt.want)
			}
		})
	}
```

### Running tests

```bash
go test ./textutil              # run the tests in one package
go test -v ./textutil           # verbose: show every test
go test -run WordCount ./...    # only tests whose name matches "WordCount", in all packages
go test -cover ./textutil       # how much of the code the tests run
```

`./...` means "this folder and everything below it".

---

## 6. Exercises

This lesson's exercises are a bit different: you'll edit files in **several packages**.

### Ex1: export the functions

Run the lesson's tests:

```bash
go test -v
```

They **don't even compile**:

```
./exercises_test.go:...: undefined: textutil.WordCount
```

Open `textutil/textutil.go`. The functions exist, but they're lowercase, so `exercises_test.go`
(in another package) can't see them. Rename `wordCount` → `WordCount` and `capitalize` → `Capitalize`.
Leave `isVowel` lowercase: it's an internal helper. Then `go test -run Ex1 -v` should pass.

### Ex2: write `Truncate`

In `textutil/textutil.go`, fill in `Truncate`. Instructions are in the comment there.
Check with `go test -run Ex2 -v`.

### Ex3: write your own test

Create a new file `textutil/textutil_test.go` and write a **table-driven** test for `Capitalize`
with at least 4 cases (include an empty string and a string that's already capitalized).
Run it with:

```bash
go test -v ./textutil
```

Then make your test fail on purpose (change a `want` value) to see your error message, and change it back.
`go test -run Ex3 -v` checks that your file exists and contains a test.

### Ex4: create a package from scratch

1. Create a folder `mathx` inside this lesson's folder.
2. In it, create `mathx.go` with `package mathx` and an exported function
   `Clamp(v, lo, hi int) int` that returns `v` limited to the range `lo`..`hi`:
   `Clamp(5, 0, 10)` is 5, `Clamp(-3, 0, 10)` is 0, `Clamp(42, 0, 10)` is 10.
3. Add a doc comment, and check it shows up with `go doc ./mathx.Clamp`.
4. Use it from `try/main.go` and run `go run ./try`.
5. Check it with the tests in the `mathx_check` folder:

```bash
go test -v ./mathx_check
```

Until your `mathx` package exists, that command fails with an error about the package not being found.
That's expected: reading that error is part of the exercise.

---

## ✅ Checklist: end of Part 1 🎉

- [ ] I know the difference between a module and a package, and what `go.mod` is for
- [ ] I know capital letters export names, and use lowercase to keep things private
- [ ] I can import and use a package from my own module
- [ ] I can write a table-driven test with `t.Errorf` and run it with `go test`
- [ ] I can create a new package from an empty folder

---

## What's next

You've finished the foundations. Before Part 2:

1. **Practise.** Do the practice sets, in order. They mix everything and are harder on purpose:
   - `practice/01_basics`
   - `practice/02_slices_maps` (goes deeper into how slices share memory)
   - `practice/03_structs_methods` (adds embedding)
2. **Build something.** The Part 1 milestone project in [roadmap.html](../../../roadmap.html):
   a small `wc`/`grep` clone. You'll need `os.Args` and `bufio.Scanner`. Look them up with `go doc`,
   which is a skill in itself.

Then Part 2: interfaces, errors, generics.

**Next → [Part 2, Lesson 10: Interfaces](../../part2/10_interfaces/README.md)**
