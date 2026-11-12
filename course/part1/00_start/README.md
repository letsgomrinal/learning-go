# Lesson 00: Your first program

**Goal:** run Go code, understand every line of a tiny program, read compiler errors without
panicking, and learn how the exercises in this course work.

**Time:** about 45 minutes.

---

## 1. Open a terminal in this folder

Every lesson is done from inside its own folder:

```bash
cd ~/Desktop/learning-go/course/part1/00_start
```

Check that Go works:

```bash
go version
```

You should see something like `go version go1.27.4 darwin/arm64`.

---

## 2. Run your first program

Open `try/main.go`. It contains:

```go
package main

import "fmt"

func main() {
	fmt.Println("Hello, Go!")
}
```

Run it:

```bash
go run ./try
```

You should see `Hello, Go!`. 🎉

`./try` means "the program in the `try` folder". Go runs programs by **folder**, not by file.

---

## 3. What each line means

Here is the same program in Python, for comparison:

```python
print("Hello, Go!")
```

Go needs a few more lines. Here's what each one does:

```go
package main
```
Every Go file starts by saying which **package** it belongs to. A package is a folder of Go files
that belong together. The name `main` is special: it means "this package is a program you can run".

```go
import "fmt"
```
Like Python's `import`. `fmt` (short for "format") is the standard package for printing.

```go
func main() {
```
`func` declares a function, like Python's `def`. A program starts running at the function named
`main` inside package `main`. Python runs a file from top to bottom; Go always starts at `main()`.

The `{` **must** be on the same line as `func main()`. Go uses braces `{ }` instead of indentation
to mark blocks.

```go
	fmt.Println("Hello, Go!")
```
Call the `Println` function from the `fmt` package. It prints its arguments and then a newline,
like Python's `print`. Notice:
- Strings use **double quotes** `"..."`. Single quotes mean something else in Go (a single character, covered in lesson 04).
- No semicolon is needed at the end of the line.
- `Println` starts with a capital letter. That's not a style choice; lesson 09 explains why.

```go
}
```
The end of the function.

---

## 4. Try it: change the program

Edit `try/main.go` so it prints three lines:

```go
func main() {
	fmt.Println("Hello, Go!")
	fmt.Println("My name is", "Ana")   // Println puts a space between its arguments
	fmt.Println(2 + 3)                  // you can print numbers and math too
}
```

Use your own name. Run `go run ./try` again.

`//` starts a comment, like `#` in Python.

---

## 5. Go is compiled

When you run a Python file, the Python interpreter reads it and runs it straight away.

Go works differently. A **compiler** first checks the whole program and turns it into a machine-code
program, and only then runs it. `go run` does both steps for you. You can also do the steps separately:

```bash
go build -o hello ./try     # compile into a program file called "hello"
./hello                     # run it. No Go needed anymore!
```

The `hello` file is a complete program. You could copy it to another Mac and run it there without
installing Go. That's one big reason people like Go for tools and servers.

Delete it when you're done: `rm hello`

---

## 6. Compiler errors are your friend

In Python, many mistakes only show up when that line runs. The Go compiler checks everything
**before** the program runs. That means you'll see a lot of compiler errors at first, and that's
a good thing: each one is a bug caught early.

Let's make some on purpose. Make each change below in `try/main.go`, run `go run ./try`, read
the error, then undo the change.

**Error 1: a variable you never use.** Add this line inside `main`, under the `Println`:

```go
	x := 5
```

You'll see:

```
./try/main.go:7:2: declared and not used: x
```

How to read it: **file** `try/main.go`, **line** 7, **column** 2, then the message. In Go, an unused
variable is an error, not a warning. It keeps code clean. (`x := 5` creates a variable; lesson 01
explains it.)

**Error 2: an import you never use.** Add `"strings"` to the imports:

```go
import (
	"fmt"
	"strings"
)
```

You'll see `"strings" imported and not used`. Unused imports are errors too.

(This also shows how to import several packages: inside `( )`, one per line.)

**Error 3: a typo.** Change `fmt.Println` to `fmt.println` (lowercase p):

```
./try/main.go:6:6: undefined: fmt.println (but have Println)
```

Go is case-sensitive, and `Println` must be spelled exactly. The compiler even suggests the right name.

**Error 4: wrong brace position.** Move the `{` onto its own line:

```go
func main()
{
```

You'll get `syntax error: unexpected semicolon or newline before {`. This is why Go code always
puts `{` at the end of the line.

---

## 7. Let the tools format your code

Go has one official code style, and a tool applies it for you:

```bash
gofmt -w try/main.go
```

In VS Code with the Go extension, this happens automatically every time you save. Nobody argues
about formatting in Go. Tabs for indentation, and everything else is decided by `gofmt`.

---

## 8. How the exercises work

Open `exercises.go`. It has small functions with a `TODO` inside. Your job is to fill them in.

Open `exercises_test.go`. These are **tests**: code that calls your functions and checks the answers.
Run the test for exercise 1:

```bash
go test -run Ex1 -v
```

It fails, because you haven't done it yet. You'll see something like:

```
=== RUN   TestEx1HelloWorld
    exercises_test.go:9: Ex1: HelloWorld() = "", want "Hello, World!"
--- FAIL: TestEx1HelloWorld (0.00s)
FAIL
```

How to read it:
- `TestEx1HelloWorld` is the test that ran.
- `exercises_test.go:9` is the line of the test that complained.
- `HelloWorld() = ""` is what your function **returned** (an empty string, because of the `TODO`).
- `want "Hello, World!"` is what it **should** return.

The quotes around values come from a printing option called `%q`. It makes empty strings and extra
spaces easy to see.

When you fix the function and run the test again, you'll see:

```
--- PASS: TestEx1HelloWorld (0.00s)
PASS
```

---

## 9. Exercises

You haven't learned functions properly yet (that's lesson 02). For now, all you need to know:

```go
func HelloWorld() string {
	return ""
}
```

means "a function called `HelloWorld` that takes nothing and **returns a string**". Whatever comes
after `return` is the result. To finish an exercise, change what it returns.

| Exercise | Run with |
|---|---|
| **Ex1** `HelloWorld`: return the text `Hello, World!` | `go test -run Ex1 -v` |
| **Ex2** `Greet`: return `Hello, ` + the name + `!`. In Go, `+` joins strings, just like Python. | `go test -run Ex2 -v` |
| **Ex3** `FixMe`: this one is already written, but the text has a typo. Run the test, compare what you got with what was wanted, and fix it. | `go test -run Ex3 -v` |
| **Ex4** `SecondsInADay`: return how many seconds are in a day. Write it as a calculation like `24 * 60 * 60`, not a number you worked out yourself. | `go test -run Ex4 -v` |

When all four pass, run everything once:

```bash
go test -v
```

---

## ✅ Checklist

You're ready for lesson 01 when you can say yes to all of these:

- [ ] I can run a program with `go run ./try`
- [ ] I know why every Go program has `package main` and `func main()`
- [ ] I can read a compiler error and find the file and line it points at
- [ ] I know that unused variables and imports are errors in Go
- [ ] I can run a single exercise's test and read a FAIL message (what I got vs. what was wanted)

**Next → [Lesson 01: Variables and types](../01_variables/README.md)**
