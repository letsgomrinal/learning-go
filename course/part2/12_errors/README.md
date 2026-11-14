# Lesson 12: Errors

**Goal:** return errors from your own functions, handle errors you receive, and use "sentinel"
errors so callers can tell which kind of problem happened.

**Time:** 1.5 hours.

```bash
cd ~/Desktop/learning-go/course/part2/12_errors
```

---

## 1. Errors are values

In Python, a function that fails **raises** an exception, and it flies up through your code until some
`try/except` catches it (or the program crashes):

```python
def parse_age(s):
    n = int(s)              # may raise ValueError
    if n < 0:
        raise ValueError("age cannot be negative")
    return n
```

Go has no exceptions. A function that can fail **returns** an error as its last value, and the caller
checks it. You've done this with `strconv.Atoi` since lesson 04:

```go
n, err := strconv.Atoi("42")
if err != nil {
	// something went wrong
}
```

Why does Go do it this way? Because you can **see** every place that can fail just by reading the
code. In Python, any line might raise, and nothing on the line tells you. The cost is more typing:
you'll write `if err != nil` a lot. Every Go programmer does.

Remember from lesson 11: `error` is just an interface with one method, `Error() string`. `nil`
means "no error".

---

## 2. Returning your own errors

Two ways to create an error:

```go
import (
	"errors"
	"fmt"
)

err1 := errors.New("age cannot be negative")             // a fixed message
err2 := fmt.Errorf("not a number: %q", s)                 // a formatted message, like Sprintf
```

A function that can fail returns `(result, error)`:

```go
func parseAge(s string) (int, error) {
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("not a number: %q", s)
	}
	if n < 0 {
		return 0, errors.New("age cannot be negative")
	}
	return n, nil          // success: the value, and a nil error
}
```

The conventions:
- The error is the **last** return value.
- On failure, return the **zero value** for the result (`0`, `""`, `nil`) along with the error.
- On success, return `nil` for the error.
- Error messages start with a **lowercase** letter and have **no** final punctuation. (They often
  get joined into longer messages, as you'll see in lesson 13.)

---

## 3. Handling errors: the `if err != nil` pattern

```go
age, err := parseAge("abc")
if err != nil {
	fmt.Println("error:", err)     // error: not a number: "abc"
	return
}
fmt.Println("age is", age)
```

The shape you'll write thousands of times:

```go
result, err := doSomething()
if err != nil {
	// handle it, usually by returning early
}
// use result here: we know it's valid
```

Notice the "happy path" stays at the left edge, and error handling goes inside the `if`s. Go code
reads top to bottom: check, return early on problems, continue. Avoid putting the main logic inside
an `else`.

Printing an error with `fmt.Println(err)` or `%v` calls its `Error()` method.

---

## 4. Four ways to handle an error

When you get an error, you have to decide what to do. Ask: **can this function do anything useful
about it?**

**a) Pass it up** to your caller (the most common choice):

```go
func loadAge(s string) (int, error) {
	age, err := parseAge(s)
	if err != nil {
		return 0, err          // not our job to decide; let the caller handle it
	}
	return age, nil
}
```

This is what Python does automatically when you don't catch an exception. In Go you write it.

**b) Use a default value:**

```go
age, err := parseAge(input)
if err != nil {
	age = 18
}
```

**c) Skip the bad item and keep going** (when processing a list):

```go
for _, s := range inputs {
	age, err := parseAge(s)
	if err != nil {
		continue
	}
	// ...
}
```

**d) Stop the program** (only in `main`, for errors you can't recover from):

```go
if err != nil {
	log.Fatal(err)     // prints the error and exits. import "log"
}
```

What you should **never** do is ignore it:

```go
age, _ := parseAge(input)      // ❌ if it failed, age is 0 and you'll never know
```

---

## 5. Sentinel errors: errors you can check for

Sometimes the caller needs to know **which** problem happened, to react differently. In Python you'd
catch different exception classes:

```python
try:
    withdraw(balance, amount)
except InsufficientFunds:
    print("not enough money")
except InvalidAmount:
    print("bad amount")
```

In Go, a package declares its errors as **variables**, so callers can compare against them:

```go
var ErrInsufficientFunds = errors.New("insufficient funds")
var ErrInvalidAmount = errors.New("amount must be positive")

func withdraw(balance, amount int) (int, error) {
	if amount <= 0 {
		return balance, ErrInvalidAmount
	}
	if amount > balance {
		return balance, ErrInsufficientFunds
	}
	return balance - amount, nil
}
```

These are called **sentinel errors**. By convention their names start with `Err`, and they're
exported (capital `E`) so other packages can use them.

The caller checks with **`errors.Is`**:

```go
_, err := withdraw(50, 100)
if errors.Is(err, ErrInsufficientFunds) {
	fmt.Println("not enough money")
} else if errors.Is(err, ErrInvalidAmount) {
	fmt.Println("bad amount")
} else if err != nil {
	fmt.Println("something else went wrong:", err)
}
```

Or with a tagless `switch`, which reads nicely:

```go
switch {
case errors.Is(err, ErrInsufficientFunds):
	fmt.Println("not enough money")
case errors.Is(err, ErrInvalidAmount):
	fmt.Println("bad amount")
case err != nil:
	fmt.Println("error:", err)
}
```

(`err == ErrInsufficientFunds` also works here, but `errors.Is` keeps working when errors get
"wrapped" with extra context, which is lesson 13. Get into the habit of `errors.Is` now.)

**Never** check an error by its message, like `err.Error() == "insufficient funds"`. Messages can
change; the variable doesn't.

You've used sentinels from the standard library already without noticing: `io.EOF` ("end of file")
is the most famous one.

### 🤔 Predict the output

```go
var ErrEmpty = errors.New("empty input")

func first(s string) (string, error) {
	if s == "" {
		return "", ErrEmpty
	}
	return s[:1], nil
}

func main() {
	for _, in := range []string{"go", "", "py"} {
		c, err := first(in)
		if errors.Is(err, ErrEmpty) {
			fmt.Print("? ")
			continue
		}
		fmt.Print(c, " ")
	}
}
```

<details markdown="1">
<summary>Answer</summary>

`g ? p `. The empty string returns `ErrEmpty`, so we print `?` and skip to the next item.
</details>

---

## 6. Errors vs panics

You've seen panics: index out of range, nil map writes, nil pointers. When do you **return an
error**, and when is a **panic** OK?

- **Errors** are for things that can go wrong in normal use: bad user input, a missing file, a
  network timeout. The caller should handle them. **This is almost always the right choice.**
- **Panics** are for **bugs**: situations that should never happen if the code is correct. A
  panic means "the programmer made a mistake", not "the user typed something wrong".

Never use panic as a substitute for exceptions in Go. Lesson 14 covers panic and recover properly.

---

## 7. Common mistakes

| You wrote | What happens | Fix |
|---|---|---|
| `n, _ := strconv.Atoi(s)` | a failure silently becomes 0 | check `err` |
| Using `result` before checking `err` | you may use a garbage zero value | check `err` first, then use the result |
| `return nil, err` in a function returning `(int, error)` | compile error: can't use nil as int | `return 0, err` |
| `errors.New("Invalid amount.")` | works, but against style | lowercase, no final punctuation |
| `if err.Error() == "insufficient funds"` | breaks when the message changes | `errors.Is(err, ErrInsufficientFunds)` |
| `panic(...)` for bad user input | crashes the program | return an error |

---

## 8. Exercises

| | Exercise | Practises |
|---|---|---|
| Ex1 | `ParseAge` with two kinds of error | `errors.New`, `fmt.Errorf`, returning `(T, error)` |
| Ex2 | `SafeDivide` returning a sentinel | sentinel errors |
| Ex3 | `SumStrings`: stop at the first bad value | passing errors up |
| Ex4 | `Withdraw` with two sentinels | choosing which error |
| Ex5 | `WithdrawMessage`: react to each error | `errors.Is`, `switch` |
| Ex6 | `ParseAll`: keep the good, collect the bad | skipping errors in a loop |
| Ex7 | `UserName` with a `"guest"` default | handling by using a default |

---

## ✅ Checklist

- [ ] I can write a function that returns `(value, error)` following the conventions
- [ ] I write `if err != nil` right after every call that can fail, and return early
- [ ] I can name four ways to handle an error, and know never to ignore one
- [ ] I can declare sentinel errors and check them with `errors.Is`
- [ ] I know errors are for expected problems and panics are for bugs

**Next → [Lesson 13: Custom errors and wrapping](../13_error_types/README.md)**
