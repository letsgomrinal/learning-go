# Lesson 01: Variables and types

**Goal:** create variables, understand Go's basic types, know what a "zero value" is, convert
between types, and print values nicely.

**Time:** 1–1.5 hours.

```bash
cd ~/Desktop/learning-go/course/part1/01_variables
```

Type every example below into `try/main.go` (inside `func main() { ... }`) and run `go run ./try`.

---

## 1. The big difference: every variable has a type

In Python, a variable can hold anything, and can change what it holds:

```python
x = 5
x = "hello"   # fine in Python
```

In Go, every variable has **one type, fixed forever**. The compiler checks it:

```go
x := 5
x = "hello"   // compile error: cannot use "hello" (untyped string constant) as int value
```

This feels strict at first, but it means a whole category of bugs (like adding a number to a string
by accident) can never reach production.

---

## 2. Creating variables: `:=`

The most common way to create a variable:

```go
name := "Ana"
age := 30
height := 1.68
isStudent := false

fmt.Println(name, age, height, isStudent)
```

`:=` means "**create** a new variable and figure out its type from the value". Go sees `30` and decides
`age` is an `int`.

To **change** an existing variable, use plain `=`:

```go
age := 30      // create
age = 31       // change (no colon!)
age := 32      // compile error: no new variables on left side of :=
```

Remember it like this: **`:=` creates, `=` changes.**

---

## 3. The basic types

| Go type | What it holds | Python equivalent | Example |
|---|---|---|---|
| `int` | whole numbers | `int` | `42`, `-7` |
| `float64` | decimal numbers | `float` | `3.14`, `-0.5` |
| `string` | text | `str` | `"hello"` |
| `bool` | true or false | `bool` | `true`, `false` (lowercase!) |

There are more number types (`int8`, `int64`, `uint`, `float32`...) for special cases. You can ignore
them for now: use `int` for whole numbers and `float64` for decimals.

Note: Python's `int` can grow without limit. Go's `int` has a maximum (about 9.2 quintillion on your Mac).
You'll almost never hit it, but it exists.

You can ask Go for a variable's type with `%T` (more on `Printf` below):

```go
height := 1.68
fmt.Printf("%T\n", height)   // float64
```

---

## 4. The other way: `var`, and zero values

You can also create a variable with `var`, giving the type yourself:

```go
var city string = "Lisbon"
var count int = 10
```

And here's something Python doesn't have. If you don't give a value, the variable gets its type's
**zero value**:

```go
var n int       // 0
var s string    // "" (empty string)
var b bool      // false
var f float64   // 0

fmt.Println(n, s, b, f)
```

There is no "undefined" or uninitialized variable in Go. Every variable always holds a valid value.
(There is a `nil`, similar to Python's `None`, but only for certain types. You'll meet it in lesson 05.)

**Which one should I use?**
- Use `:=` almost always, inside functions.
- Use `var x int` when you want to start from the zero value, for example a counter: `var total int`.

### 🤔 Predict the output

```go
var a int
b := 3
a = a + b
var s string
fmt.Println(a, s == "", b)
```

<details><summary>Answer</summary>

`3 true 3`. `a` starts at 0, then becomes 0 + 3. `s` is the empty string, so `s == ""` is `true`.
</details>

---

## 5. Several variables at once

```go
x, y := 1, 2
fmt.Println(x, y)    // 1 2

x, y = y, x          // swap! Exactly like Python
fmt.Println(x, y)    // 2 1
```

---

## 6. Constants

A value that never changes is a **constant**:

```go
const pi = 3.14159
const appName = "Gopher Notes"

pi = 3      // compile error: cannot assign to pi
```

Python has no real constants (only the `UPPER_CASE` convention). In Go, the compiler enforces it.
Go constants use normal `camelCase` names, not `UPPER_CASE`.

---

## 7. Math, and the integer-division trap

```go
fmt.Println(7 + 2)    // 9
fmt.Println(7 - 2)    // 5
fmt.Println(7 * 2)    // 14
fmt.Println(7 / 2)    // 3     ← !!
fmt.Println(7 % 2)    // 1     (remainder)
fmt.Println(7.0 / 2)  // 3.5
```

**`7 / 2` is `3` in Go.** When both sides are `int`, `/` does whole-number division and drops the
decimal part. (Like Python's `//`, except Go rounds toward zero, so `-7 / 2` is `-3`, where Python's `-7 // 2` is `-4`.)

There's no `**` operator for powers. Use `math.Pow(2, 10)` (import `"math"`), or just multiply.

`x += 1` works like in Python. Go also has `x++` and `x--` to add or subtract 1.

---

## 8. Converting between types

Go **never** converts types automatically. This doesn't compile:

```go
count := 3
price := 2.5
total := count * price   // compile error: mismatched types int and float64
```

You must convert explicitly, by using the type name like a function:

```go
total := float64(count) * price
fmt.Println(total)    // 7.5
```

Converting a float to an int **cuts off** the decimals (it doesn't round):

```go
fmt.Println(int(2.99))   // compile error! a constant 2.99 can't become an int
f := 2.99
fmt.Println(int(f))      // 2
```

(The first line fails because Go checks constant values at compile time, and refuses to silently
lose the `.99`. With a variable, it does what you asked.)

### The average trap

```go
a, b := 1, 2
avg := (a + b) / 2
fmt.Println(avg)          // 1   ← integer division!

avg2 := float64(a+b) / 2
fmt.Println(avg2)         // 1.5
```

Convert **before** dividing. You'll use this in the exercises.

### Numbers ↔ strings

This is a common surprise:

```go
n := 65
s := string(n)
fmt.Println(s)            // "A" !!  not "65"
```

`string(65)` makes the character with code 65, which is `A`. To turn a number into its digits, use
the `strconv` package (import `"strconv"`):

```go
s := strconv.Itoa(65)     // "65"   (Itoa = "Integer to ASCII")
```

Going the other way (`"65"` → `65`) can fail, so it needs error handling. That's in lesson 04.

---

## 9. Printing nicely with `Printf`

`fmt.Println` is like Python's `print`. `fmt.Printf` is like Python's f-strings or `%` formatting.
You write the text with **verbs** (placeholders starting with `%`), then the values:

```go
name := "Ana"
age := 30
height := 1.68

fmt.Printf("%s is %d years old and %.2f m tall\n", name, age, height)
// Ana is 30 years old and 1.68 m tall
```

| Verb | Meaning |
|---|---|
| `%v` | any value, default format (use it when unsure) |
| `%d` | an integer |
| `%f`, `%.2f` | a float; `.2` means 2 decimal places |
| `%s` | a string |
| `%q` | a string with quotes around it, `"like this"` |
| `%t` | a bool |
| `%T` | the **type** of the value |

`Printf` does **not** add a newline. End the text with `\n`.

The same in Python: `print(f"{name} is {age} years old and {height:.2f} m tall")`

`fmt.Sprintf` works exactly like `Printf`, but **returns** the string instead of printing it. You'll
need it in the exercises:

```go
msg := fmt.Sprintf("%s is %d", name, age)   // msg is "Ana is 30"
```

### 🤔 Predict the output

```go
x := 10
y := 4
fmt.Printf("%d %v %.1f %T\n", x/y, x%y, float64(x)/float64(y), x/y)
```

<details><summary>Answer</summary>

`2 2 2.5 int`. `10 / 4` is 2 (integer division), `10 % 4` is 2, the float division is 2.5,
and `x/y` is an `int`.
</details>

---

## 10. Common mistakes

| You wrote | What happens | Fix |
|---|---|---|
| `x := 1` then later `x := 2` | `no new variables on left side of :=` | Use `x = 2` |
| `x := 5` and never use it | `declared and not used` | Use it or delete it |
| `true`/`True` | `True` is undefined | Go uses lowercase `true` / `false` |
| `'hello'` | error: more than one character in rune literal | Strings use double quotes |
| `count * price` with int and float64 | `mismatched types` | `float64(count) * price` |
| `string(42)` | gives `"*"`, not `"42"` | `strconv.Itoa(42)` |

---

## 11. Exercises

Open `exercises.go`. Do them in order. Each one has hints.

| | Exercise | Practises |
|---|---|---|
| Ex1 | `ZeroValues`: declare variables without values and show what they hold | zero values, `var` |
| Ex2 | `CelsiusToFahrenheit` | float math |
| Ex3 | `Average` of three ints, as a float64 | converting before dividing |
| Ex4 | `SharePerPerson`: split a bill in cents | `/` with ints |
| Ex5 | `LeftoverCents`: the cents that can't be split evenly | `%` with ints |
| Ex6 | `Describe`: build a sentence with `fmt.Sprintf` | `Sprintf` and verbs |
| Ex7 | `AgeInDays` returning a string like `"30 years is about 10950 days"` | constants, `Sprintf` |

```bash
go test -run Ex1 -v
```

---

## ✅ Checklist

- [ ] I know when to use `:=` and when to use `=`
- [ ] I can name the zero value of `int`, `float64`, `string` and `bool`
- [ ] I know why `7 / 2` is `3`, and how to get `3.5`
- [ ] I can convert between `int` and `float64`, and turn an int into a string with `strconv.Itoa`
- [ ] I can use `Printf` / `Sprintf` with `%d`, `%s`, `%.2f` and `%v`

**Next → [Lesson 02: Functions](../02_functions/README.md)**
