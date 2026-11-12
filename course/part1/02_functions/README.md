# Lesson 02: Functions

**Goal:** write your own functions, return one or several values, and treat functions as values
you can store and pass around.

**Time:** about 1 hour.

```bash
cd ~/Desktop/learning-go/course/part1/02_functions
```

**Important:** in this lesson you'll write functions in `try/main.go`. They go **outside** `main`,
either above or below it. Then you call them from inside `main`:

```go
package main

import "fmt"

func double(n int) int {     // ← your function, outside main
	return n * 2
}

func main() {
	fmt.Println(double(21))  // ← call it from main
}
```

---

## 1. Your first function

Python:

```python
def add(a: int, b: int) -> int:
    return a + b
```

Go:

```go
func add(a int, b int) int {
	return a + b
}
```

Reading it piece by piece:

```
func   add   (a int, b int)   int   {
 │      │          │           │
 │      │          │           └── the type it returns
 │      │          └── parameters: name first, THEN type
 │      └── the function's name
 └── keyword, like Python's def
```

In Go, **the name comes before the type**: `a int`, not `int a`. In Python type hints are optional
and nothing checks them at runtime. In Go types are required and the compiler checks every call:

```go
add(1, 2)        // 3
add(1, "two")    // compile error: cannot use "two" (untyped string constant) as int value
add(1)           // compile error: not enough arguments in call to add
```

When several parameters in a row have the same type, you can write the type once:

```go
func add(a, b int) int {     // same as (a int, b int)
	return a + b
}
```

---

## 2. Functions that return nothing

Leave out the return type:

```go
func greet(name string) {
	fmt.Println("Hello,", name)
}
```

Unlike Python, there is no hidden `None` result. A function either declares what it returns, or it
returns nothing at all, and `x := greet("Ana")` is a compile error.

---

## 3. Returning several values

This is one of Go's most-used features. In Python you'd return a tuple:

```python
def min_max(a, b):
    if a < b:
        return a, b
    return b, a

lo, hi = min_max(7, 3)
```

In Go, list the return types in parentheses:

```go
func minMax(a, b int) (int, int) {
	if a < b {
		return a, b
	}
	return b, a
}

func main() {
	lo, hi := minMax(7, 3)
	fmt.Println(lo, hi)   // 3 7
}
```

(Don't worry about the `if` details yet. That's the next lesson. It works like Python's `if`, just
with braces and no colon.)

You'll see this pattern **everywhere** in Go, because it's how Go handles errors:

```go
value, err := someFunctionThatCanFail()
```

The function returns the result **and** an error. Part 2 covers errors properly, but you'll start
noticing this pattern now.

### Ignoring a value with `_`

If you only need one of the results, put `_` (the blank identifier) in place of the other:

```go
_, hi := minMax(7, 3)    // we don't care about lo
fmt.Println(hi)
```

You have to do this, because an unused variable is a compile error (remember lesson 00).

---

## 4. Named return values

You can give the return values names. They start at their zero value, and a bare `return` returns them:

```go
func divide(a, b int) (quotient, remainder int) {
	quotient = a / b
	remainder = a % b
	return     // returns quotient and remainder
}
```

The names also document what each value means. Use this in short functions; in long ones a bare
`return` makes it hard to see what's being returned.

---

## 5. Functions are values

Like in Python, a function is a value. You can put it in a variable and pass it to another function:

```go
func apply(n int, f func(int) int) int {
	return f(n)
}

func double(n int) int { return n * 2 }

func main() {
	fmt.Println(apply(5, double))   // 10
}
```

The type `func(int) int` means "a function that takes an int and returns an int".

### Anonymous functions (Python's `lambda`, but any length)

```go
square := func(n int) int {
	return n * n
}
fmt.Println(square(4))                          // 16
fmt.Println(apply(3, func(n int) int { return n + 100 }))   // 103
```

Python's `lambda` can only hold one expression. Go's anonymous functions can hold any code.

### Closures

An anonymous function can use (and change) variables from around it, just like a Python closure:

```go
func makeCounter() func() int {
	count := 0
	return func() int {
		count++
		return count
	}
}

func main() {
	next := makeCounter()
	fmt.Println(next())   // 1
	fmt.Println(next())   // 2
	fmt.Println(next())   // 3
}
```

Each call to `makeCounter()` creates a new, separate `count`. (In Python you'd need `nonlocal count`;
Go doesn't need anything special.)

### 🤔 Predict the output

```go
func calc(a, b int) (int, int) {
	return a + b, a * b
}

func main() {
	s, p := calc(3, 4)
	_, p2 := calc(s, p)
	fmt.Println(s, p, p2)
}
```

<details><summary>Answer</summary>

`7 12 84`. First call: 3+4=7 and 3×4=12. Second call: calc(7, 12), and we keep only the product, 84.
</details>

---

## 6. What Go functions don't have

Coming from Python, you might look for these. Go leaves them out on purpose, to keep code simple and explicit:

| Python | Go |
|---|---|
| Default parameter values `def f(x=10)` | Not available. Write two functions, or pass a value explicitly. |
| Keyword arguments `f(name="Ana")` | Not available. Arguments are always by position. |
| `*args` | Yes, as "variadic" functions (`func sum(nums ...int)`). Covered in lesson 05, after slices. |
| `**kwargs` | Not available. Use a struct (lesson 07). |
| Overloading (same name, different params) | Not available. Use different names. |

---

## 7. Common mistakes

| You wrote | What happens | Fix |
|---|---|---|
| `func add(int a, int b)` | syntax error | Name first: `a int` |
| Defined a function inside `main` with `func name() {...}` | syntax error | Define functions outside `main`, or use `name := func() {...}` inside |
| Forgot `return` in a function with a return type | `missing return` | Every path must return a value |
| `x := minMax(1, 2)` | `assignment mismatch: 1 variable but minMax returns 2 values` | `lo, hi := minMax(1, 2)` |

---

## 8. Exercises

From now on, the exercises ask you to write the **whole function body**. The function signatures
(first lines) are already there.

| | Exercise | Practises |
|---|---|---|
| Ex1 | `RectangleArea` | parameters, return |
| Ex2 | `Perimeter` with shared parameter types | `(w, h float64)` style |
| Ex3 | `Swap` two strings | returning 2 values |
| Ex4 | `Divmod`: quotient and remainder | named results |
| Ex5 | `SumAndProduct` of three numbers | returning 2 values |
| Ex6 | `ApplyTwice`: call a function twice | functions as parameters |
| Ex7 | `MakeMultiplier`: return a function | closures |

```bash
go test -run Ex1 -v
```

---

## ✅ Checklist

- [ ] I can write a function with parameters and a return type, name before type
- [ ] I can return two values and receive them with `a, b := f()`
- [ ] I know how to ignore a return value with `_`
- [ ] I can pass a function to another function, and write an anonymous function
- [ ] I can explain how `makeCounter` remembers `count` between calls

**Next → [Lesson 03: Control flow](../03_control_flow/README.md)**
