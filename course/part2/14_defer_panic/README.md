# Lesson 14: `defer`, `panic` and `recover`

**Goal:** clean up reliably with `defer`, understand what a panic does, and know the few situations
where `recover` is the right tool.

**Time:** about 1 hour.

```bash
cd ~/Desktop/learning-go/course/part2/14_defer_panic
```

---

## 1. `defer`: "do this when the function ends"

Python makes sure cleanup happens with `with` or `try/finally`:

```python
with open("data.txt") as f:
    process(f)
# the file is closed here, even if process() raised
```

Go uses **`defer`**. A deferred call runs when the surrounding function **returns**, however it
returns:

```go
func processFile(name string) error {
	f, err := os.Open(name)
	if err != nil {
		return err
	}
	defer f.Close()        // runs when processFile returns, whatever happens below

	// ... read from f, maybe return early in several places ...
	return nil
}
```

The big advantage: the cleanup sits **right next to** the thing it cleans up. You open the file, and
on the very next line you say how it will be closed. You can't forget it at one of the return statements.

Try it:

```go
func main() {
	defer fmt.Println("3. deferred: runs last")
	fmt.Println("1. first")
	fmt.Println("2. second")
}
```

---

## 2. Several defers run in reverse order

Deferred calls are stacked up, and run **last in, first out**, like a stack of plates:

```go
func main() {
	defer fmt.Println("a")
	defer fmt.Println("b")
	defer fmt.Println("c")
	fmt.Println("start")
}
// start
// c
// b
// a
```

That's the right order for cleanup: if you open A, then B (which depends on A), you want to close B
before A.

---

## 3. Arguments are evaluated immediately

The deferred **call** happens later, but its **arguments** are worked out right away, at the
`defer` line:

```go
x := 1
defer fmt.Println("deferred x:", x)   // x is read NOW: 1
x = 2
fmt.Println("x:", x)
// x: 2
// deferred x: 1
```

If you want the value at the end, defer a function that reads it when it runs:

```go
defer func() {
	fmt.Println("x at the end:", x)   // reads x when this runs: 2
}()
```

Note the `()` at the end: you're deferring a **call** to that anonymous function.

### 🤔 Predict the output

```go
func main() {
	for i := range 3 {
		defer fmt.Print(i, " ")
	}
	fmt.Print("go ")
}
```

<details markdown="1">
<summary>Answer</summary>

`go 2 1 0 `. Three defers stack up with i = 0, 1, 2 (evaluated at each `defer`), then run in
reverse when `main` returns.
</details>

---

## 4. Common uses of `defer`

```go
f, _ := os.Open(name)
defer f.Close()               // close files

mu.Lock()
defer mu.Unlock()             // unlock mutexes (Part 3)

start := time.Now()
defer func() { fmt.Println("took", time.Since(start)) }()   // time a function

resp, _ := http.Get(url)
defer resp.Body.Close()       // close network responses (Part 4)
```

⚠️ `defer` runs when the **function** ends, not at the end of a loop iteration. Deferring inside a
loop that opens 10,000 files keeps all 10,000 open until the function returns. Move the loop body
into its own function instead.

---

## 5. Panic: when something is badly wrong

A **panic** stops normal execution. You've already caused some:

```go
var m map[string]int
m["a"] = 1              // panic: assignment to entry in nil map

s := []int{1, 2}
fmt.Println(s[5])       // panic: runtime error: index out of range [5] with length 2
```

You can also panic yourself:

```go
panic("this should never happen")
```

What happens during a panic:
1. The current function stops.
2. Its deferred calls **still run** (that's why `defer` is reliable for cleanup).
3. Then the function's caller stops and runs its defers, and so on up the stack.
4. If nothing stops it, the program crashes and prints the panic message and a stack trace.

It looks like a Python exception travelling up. But in Go it's **not** for normal errors. As
lesson 12 said: errors are for expected problems, **panics are for bugs**.

When it's reasonable to panic yourself:
- A situation that can only happen if the code itself is wrong ("unreachable" branches).
- At program start-up, when a required setting is missing and the program can't do anything useful.
  The standard library has functions named `MustXxx` that panic instead of returning an error, for
  exactly this case, for example `regexp.MustCompile(` + "`" + `[0-9]+` + "`" + `)` for a pattern written in your code.

---

## 6. `recover`: catching a panic

`recover()` stops a panic and returns the value that was passed to `panic`. It only works **inside
a deferred function**:

```go
func safeCall() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("recovered from:", r)
		}
	}()

	panic("oh no")
	fmt.Println("never printed")
}

func main() {
	safeCall()
	fmt.Println("main keeps running")
}
// recovered from: oh no
// main keeps running
```

If there's no panic, `recover()` returns `nil`, so the `if` does nothing.

### Turning a panic into an error

Combine `recover` with a **named result** (lesson 02). A deferred function runs after the `return`
value is set, but **before** the caller gets it, so it can still change a named result:

```go
func safeDivide(a, b int) (result int, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("recovered: %v", r)   // change the returned error
		}
	}()
	return a / b, nil            // integer division by zero panics
}

fmt.Println(safeDivide(10, 2))   // 5 <nil>
fmt.Println(safeDivide(1, 0))    // 0 recovered: runtime error: integer divide by zero
```

### When to use `recover`

Rarely. The legitimate uses are at **boundaries**:
- A web server recovers from a panic in **one request** so the other requests keep working. (Go's
  `net/http` server does this for you.)
- A worker that runs jobs written by other people, so one bad job doesn't kill the whole program.

Don't use `panic`/`recover` as your error-handling system. Return errors.

---

## 7. Common mistakes

| You wrote | What happens | Fix |
|---|---|---|
| `recover()` called directly, not in a deferred function | returns nil, doesn't stop the panic | `defer func() { if r := recover(); ... }()` |
| `defer func() {...}` without `()` at the end | compile error: expression in defer must be function call | add `()` |
| Deferring inside a long loop | resources pile up until the function returns | put the loop body in a function |
| Expecting `defer fmt.Println(x)` to print the final `x` | prints the value at the `defer` line | defer a closure |
| `panic` for invalid user input | the program crashes | return an error |

---

## 8. Exercises

The exercises use a small `Log` type that records lines, so the tests can check what ran and in
which order.

| | Exercise | Practises |
|---|---|---|
| Ex1 | `RunTask`: open, work, close with `defer` | basic `defer` |
| Ex2 | `CloseAll`: defer in a loop | LIFO order |
| Ex3 | `Measure`: "end" is logged even if the work panics | defers run during panics |
| Ex4 | `SafeDivide` | `recover` + named results |
| Ex5 | `SafeRun`: turn any panic into an error | `recover` |
| Ex6 | `MustPositive`: panic on a bug | `panic` |
| Ex7 | `RunJobs`: keep going when a job panics | using `SafeRun` at a boundary |

---

## ✅ Checklist

- [ ] I use `defer` right after acquiring something that needs cleanup
- [ ] I know defers run in reverse order, and that arguments are evaluated at the `defer` line
- [ ] I know what happens step by step during a panic, and that defers still run
- [ ] I can use `recover` in a deferred function to turn a panic into an error
- [ ] I know panics are for bugs, and errors are for everything else

**Next → [Lesson 15: Generics](../15_generics/README.md)**
