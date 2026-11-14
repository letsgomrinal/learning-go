# Lesson 11: Interfaces in practice

**Goal:** use the interfaces that run through all Go code (`fmt.Stringer`, `io.Writer`, `io.Reader`),
work with `any`, and get the concrete type back out of an interface value.

**Time:** 1.5 hours.

```bash
cd ~/Desktop/learning-go/course/part2/11_interfaces_practice
```

---

## 1. `fmt.Stringer`: Go's `__str__`

In Python, you control how an object prints with `__str__`:

```python
class Money:
    def __init__(self, cents): self.cents = cents
    def __str__(self): return f"${self.cents / 100:.2f}"

print(Money(1250))    # $12.50
```

In Go, the `fmt` package defines this interface:

```go
type Stringer interface {
	String() string
}
```

Give your type a `String() string` method, and every `fmt` print function uses it automatically:

```go
type Money int   // amount in cents

func (m Money) String() string {
	return fmt.Sprintf("$%.2f", float64(m)/100)
}

func main() {
	price := Money(1250)
	fmt.Println(price)                   // $12.50
	fmt.Printf("Total: %v\n", price)     // Total: $12.50
}
```

You never wrote "Money implements Stringer", and `fmt` was written long before your type existed.
This is the implicit interfaces from lesson 10, paying off.

⚠️ Inside `String()`, don't print the value itself with `%v`: `fmt.Sprintf("%v", m)` would call
`String()` again, forever (a stack overflow). Convert to the underlying type first, like
`float64(m)` above.

---

## 2. `any`: a value of any type

An interface with **no methods** is satisfied by every type. It's written `interface{}`, and Go has
a shorter name for it: **`any`**.

```go
var x any
x = 42
x = "hello"
x = []int{1, 2}      // all fine
```

This is as close as Go gets to a Python variable. `fmt.Println` takes `...any`, which is how it
accepts anything.

But an `any` value is almost useless by itself: it has no methods, so you can't do anything with it.

```go
var x any = 5
y := x + 1           // compile error: mismatched types any and int
```

To use it, you need to get the real type back out.

**Use `any` sparingly.** It throws away the compiler's type checking. Most Go code never needs it.

---

## 3. Type assertions: getting the concrete type back

A **type assertion** asks "is the value inside this interface of type T?":

```go
var x any = 42

n, ok := x.(int)
fmt.Println(n, ok)       // 42 true

s, ok := x.(string)
fmt.Println(s, ok)       // "" false   (not a string; s gets the zero value)
```

It's the "comma ok" pattern again (like map lookups). Python equivalent: `isinstance(x, int)`.

There's a short form without `ok`, but it **panics** if the type is wrong:

```go
s := x.(string)          // panic: interface conversion: interface {} is int, not string
```

**Always use the `, ok` form** unless you're absolutely certain.

You can also assert to an **interface**: "does this value have these methods?"

```go
var x any = Money(500)
if s, ok := x.(fmt.Stringer); ok {
	fmt.Println("it can describe itself:", s.String())
}
```

---

## 4. Type switches

When there are several possible types, a **type switch** is cleaner than a chain of assertions:

```go
func describe(v any) string {
	switch x := v.(type) {
	case int:
		return fmt.Sprintf("an int: %d", x)       // x is an int here
	case string:
		return fmt.Sprintf("a string of length %d", len(x))   // x is a string here
	case bool:
		return "a bool"
	case nil:
		return "nothing"
	default:
		return fmt.Sprintf("something else: %T", x)
	}
}
```

`v.(type)` only works inside a `switch`. In each `case`, the variable `x` has that case's type, so you
can use it fully. Cases can also be interfaces, like `case fmt.Stringer:`.

### 🤔 Predict the output

```go
values := []any{3, "go", 2.5, true, 4}
total := 0
for _, v := range values {
	if n, ok := v.(int); ok {
		total += n
	}
}
fmt.Println(total)
```

<details markdown="1">
<summary>Answer</summary>

`7`. Only `3` and `4` are ints. `2.5` is a `float64`, so `v.(int)` gives `ok == false` for it.
</details>

---

## 5. `io.Writer`: anything you can write bytes to

This interface, from package `io`, is everywhere:

```go
type Writer interface {
	Write(p []byte) (n int, err error)
}
```

`[]byte` is a slice of bytes (raw data, which text is made of). Many very different things are Writers:

| Writer | Where the bytes go |
|---|---|
| `os.Stdout` | the terminal |
| a file from `os.Create` | a file on disk |
| `&strings.Builder{}` | a string in memory |
| `&bytes.Buffer{}` | a buffer in memory |
| an HTTP response | the web browser (Part 4) |
| a gzip writer | compressed into another Writer |

`fmt.Fprintf` is `Printf` that writes to **any** Writer you give it (the `F` is for "file"):

```go
func greet(w io.Writer, name string) {
	fmt.Fprintf(w, "Hello, %s!\n", name)
}

func main() {
	greet(os.Stdout, "terminal")             // prints to the screen

	var sb strings.Builder
	greet(&sb, "memory")                     // writes into sb
	fmt.Printf("captured: %q\n", sb.String())
}
```

`greet` doesn't know or care where its text ends up. That's why `io.Writer` is so powerful: write
your function once against `io.Writer`, and it works with the terminal, files, the network, and in
tests with a `strings.Builder` (that's how you test output!).

Python's version of this idea is "file-like objects" with a `.write()` method.

---

## 6. `io.Reader`: anything you can read bytes from

```go
type Reader interface {
	Read(p []byte) (n int, err error)
}
```

Files, network connections, HTTP request bodies, `os.Stdin`, and `strings.NewReader("some text")`
are all Readers.

You rarely call `Read` yourself. Usually you wrap the reader in a `bufio.Scanner`, which reads it line
by line:

```go
func countLines(r io.Reader) int {
	sc := bufio.NewScanner(r)
	n := 0
	for sc.Scan() {          // reads the next line; false at the end
		line := sc.Text()    // the line, without the "\n"
		_ = line
		n++
	}
	return n
}

func main() {
	fmt.Println(countLines(strings.NewReader("a\nb\nc\n")))   // 3
}
```

The same `countLines` works on a file (`os.Open`) or on `os.Stdin`, unchanged. In tests you pass a
`strings.NewReader`, so you never need real files.

Python: `for line in f:`.

---

## 7. Writing your own Writer

You can make your own types into Writers by adding a `Write` method. Here's one that counts how many
bytes pass through it on their way to another Writer:

```go
type CountingWriter struct {
	W     io.Writer
	Count int
}

func (c *CountingWriter) Write(p []byte) (int, error) {
	n, err := c.W.Write(p)   // pass the bytes on
	c.Count += n
	return n, err
}
```

Now `fmt.Fprintf(&CountingWriter{W: os.Stdout}, ...)` works, because `*CountingWriter` is an
`io.Writer`. Wrapping one Writer in another like this is a very common Go pattern (like Python decorators).

---

## 8. `error` is an interface too

You've seen `err` since lesson 04. Here's its secret:

```go
type error interface {
	Error() string
}
```

Any type with an `Error() string` method is an error. Lessons 12 and 13 build on this.

---

## 9. Common mistakes

| You wrote | What happens | Fix |
|---|---|---|
| `x + 1` where `x` is `any` | mismatched types | `n, ok := x.(int)` first |
| `x.(string)` without `ok` on the wrong type | panic | `s, ok := x.(string)` |
| `String()` that uses `%v` on itself | infinite recursion / stack overflow | convert first: `float64(m)` |
| `fmt.Fprintf(sb, ...)` with a `strings.Builder` value | compile error: Write has a pointer receiver | `&sb` |
| `any` everywhere "to be flexible" | you lose compiler checks | use real types or small interfaces |

---

## 10. Exercises

| | Exercise | Practises |
|---|---|---|
| Ex1 | `Money` prints as `$12.50` | `fmt.Stringer` |
| Ex2 | `Point` prints as `(1, 2)` | `fmt.Stringer` |
| Ex3 | `Describe` any value | type switch |
| Ex4 | `SumInts` from a `[]any` | type assertion with `ok` |
| Ex5 | `WriteGreetings` to any Writer | `io.Writer`, `fmt.Fprintf` |
| Ex6 | `LongestLine` from any Reader | `io.Reader`, `bufio.Scanner` |
| Ex7 | `UpperWriter`: your own Writer | implementing `io.Writer` |

---

## ✅ Checklist

- [ ] I can make a type print nicely with a `String()` method
- [ ] I know what `any` is, and why to avoid it when I can
- [ ] I can use a type assertion with `, ok`, and a type switch
- [ ] I can write a function that takes an `io.Writer` or `io.Reader`, and test it with `strings.Builder` / `strings.NewReader`
- [ ] I can make my own type satisfy `io.Writer`

**Next → [Lesson 12: Errors](../12_errors/README.md)**
