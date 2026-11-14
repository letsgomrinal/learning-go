# Lesson 13: Custom errors and wrapping

**Goal:** add context to errors as they travel up through your code, create error types that carry
extra information, and inspect errors with `errors.Is` and `errors.As`.

**Time:** 1.5 hours.

```bash
cd ~/Desktop/learning-go/course/part2/13_error_types
```

---

## 1. The problem: errors without context

Imagine a program that starts a server, which loads a config, which opens a file. The file is missing.
If each layer just passes the error up with `return err`, the message that reaches you is:

```
open app.yaml: no such file or directory
```

Which part of the program tried to open it? Why? In Python you'd get a full traceback. In Go, you build
the story yourself by **adding context** at each layer.

---

## 2. Wrapping with `%w`

`fmt.Errorf` has a special verb, **`%w`** ("wrap"). It puts the original error inside a new one with
more context:

```go
func loadConfig(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("load config %s: %w", path, err)
	}
	return string(data), nil
}

func startServer() error {
	_, err := loadConfig("app.yaml")
	if err != nil {
		return fmt.Errorf("start server: %w", err)
	}
	return nil
}
```

Now the message reads like a breadcrumb trail, from the outermost layer to the root cause:

```
start server: load config app.yaml: open app.yaml: no such file or directory
```

This is why error messages are lowercase with no final punctuation: they get joined with `: `.

Python comparison: `raise ConfigError("load config") from original_error`.

### `%w` vs `%v`

Both produce the same **text**. The difference: `%w` **keeps the original error inside**, so code
higher up can still check for it. `%v` only copies its message, and the original is lost.

**Use `%w` when you're adding context to an error you received.**

---

## 3. `errors.Is` looks through the wrapping

This is why lesson 12 told you to use `errors.Is` instead of `==`:

```go
var ErrNotFound = errors.New("not found")

err := fmt.Errorf("get user 42: %w", ErrNotFound)

fmt.Println(err == ErrNotFound)            // false: err is a NEW error that contains ErrNotFound
fmt.Println(errors.Is(err, ErrNotFound))   // true:  errors.Is unwraps and finds it
```

Wrapped errors form a **chain**. `errors.Is` walks the whole chain looking for a match, however many
layers deep.

### 🤔 Predict the output

```go
base := errors.New("disk full")
a := fmt.Errorf("save file: %w", base)
b := fmt.Errorf("export report: %v", a)    // note: %v
c := fmt.Errorf("export report: %w", a)

fmt.Println(errors.Is(a, base), errors.Is(b, base), errors.Is(c, base))
fmt.Println(b)
```

<details markdown="1">
<summary>Answer</summary>

```
true false true
export report: save file: disk full
```

`b` used `%v`, so it has the same message but the chain was cut: `base` is no longer inside it.
</details>

---

## 4. Custom error types: errors that carry data

Sometimes the caller needs more than "which error". A validation error should say **which field**
was wrong. A sentinel can't carry that, but a struct can:

```go
type ValidationError struct {
	Field   string
	Problem string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("invalid %s: %s", e.Field, e.Problem)
}
```

Because it has an `Error() string` method, `*ValidationError` **is an error** (lesson 11, section 8).
Return it like any other error:

```go
func validateEmail(email string) error {
	if !strings.Contains(email, "@") {
		return &ValidationError{Field: "email", Problem: "missing @"}
	}
	return nil
}
```

Conventions: the type name ends in `Error`, and you usually return a **pointer** (`&ValidationError{...}`)
with a pointer receiver on `Error()`.

Python comparison: defining `class ValidationError(Exception)` with extra attributes.

---

## 5. `errors.As`: getting your error type back out

`errors.Is` asks "is this particular error value somewhere in the chain?". **`errors.As`** asks
"is there an error of **this type** somewhere in the chain? If so, give it to me":

```go
err := fmt.Errorf("register ana: %w", validateEmail("ana.example.com"))

var ve *ValidationError
if errors.As(err, &ve) {
	fmt.Println("bad field:", ve.Field)        // bad field: email
	fmt.Println("problem:", ve.Problem)        // problem: missing @
}
```

How it works, step by step:
1. Declare an empty variable of the type you're looking for: `var ve *ValidationError`.
2. Pass its **address**: `errors.As(err, &ve)`. (It needs the address so it can fill `ve` in, just
   like the pointer exercises in lesson 08.)
3. If it finds a match anywhere in the chain, it fills in `ve` and returns `true`.

Python comparison: `except ValidationError as ve: print(ve.field)`.

**Which one when?**

| You want to know... | Use |
|---|---|
| "Is it this specific error?" (a sentinel) | `errors.Is(err, ErrNotFound)` |
| "Is it this kind of error, and what's inside?" (a type) | `errors.As(err, &target)` |

---

## 6. ⚠️ The nil interface trap

This is one of Go's most famous gotchas, and a classic interview question. It catches experienced
Go developers too.

An interface value holds two things: a **type** and a **value**. It only equals `nil` when **both**
are empty. Watch:

```go
func findUser(id int) error {
	var err *NotFoundError          // a nil pointer of type *NotFoundError
	if id < 0 {
		err = &NotFoundError{ID: id}
	}
	return err                      // ⚠️ BUG
}

func main() {
	err := findUser(5)              // a valid id, so this should succeed...
	fmt.Println(err == nil)         // false!!
}
```

When `findUser` returns `err`, Go puts it into an `error` interface. The interface now holds
**type = `*NotFoundError`, value = nil**. The type part isn't empty, so the interface isn't `nil`,
and the caller thinks something failed.

Picture an envelope: an empty envelope labelled "*NotFoundError" is still an envelope. `nil` means
"no envelope at all".

**The fix:** never return a typed pointer variable as an `error`. Return a literal `nil`:

```go
func findUser(id int) error {
	if id < 0 {
		return &NotFoundError{ID: id}
	}
	return nil
}
```

Rule: **functions that return `error` should `return nil` explicitly on success**, never a variable
of a concrete error type.

---

## 7. A complete example

```go
var ErrNotFound = errors.New("not found")

type ValidationError struct{ Field, Problem string }

func (e *ValidationError) Error() string { return "invalid " + e.Field + ": " + e.Problem }

func statusCode(err error) int {
	var ve *ValidationError
	switch {
	case err == nil:
		return 200
	case errors.Is(err, ErrNotFound):
		return 404
	case errors.As(err, &ve):
		return 400
	default:
		return 500
	}
}
```

This is how a real web server decides which HTTP status to send back, whatever layers of wrapping the
error went through. You'll write this in the exercises.

---

## 8. Common mistakes

| You wrote | What happens | Fix |
|---|---|---|
| `fmt.Errorf("load: %v", err)` when callers need to check it | the chain is cut | `%w` |
| `err == ErrNotFound` on a wrapped error | false | `errors.Is` |
| `errors.As(err, ve)` without `&` | `go vet` flags it; at runtime it panics: `errors: target must be a non-nil pointer` | `errors.As(err, &ve)` |
| `var ve ValidationError` when `Error()` has a pointer receiver | `errors.As` never matches | `var ve *ValidationError` |
| Returning a nil `*MyError` variable as `error` | caller sees a non-nil error | `return nil` |
| Wrapping at every single line | very long messages | add context where it helps explain *what you were doing* |

---

## 9. Exercises

| | Exercise | Practises |
|---|---|---|
| Ex1 | `LoadConfig`: wrap an error with context | `%w` |
| Ex2 | `ValidationError.Error()` | a custom error type |
| Ex3 | `ValidateUser` returns a `*ValidationError` | returning custom errors |
| Ex4 | `RegisterUser` wraps the validation error | wrapping custom errors |
| Ex5 | `BadField`: dig the field name out of any error | `errors.As` |
| Ex6 | Fix `FindUser`'s nil interface bug | the nil trap |
| Ex7 | `StatusCode` for a web server | `errors.Is` + `errors.As` together |

---

## ✅ Checklist

- [ ] I add context to errors with `fmt.Errorf("doing X: %w", err)`
- [ ] I know the difference between `%w` and `%v`
- [ ] I can write a custom error type and return it
- [ ] I know when to use `errors.Is` and when to use `errors.As`
- [ ] I can explain the nil interface trap, and how to avoid it

**Next → [Lesson 14: defer, panic and recover](../14_defer_panic/README.md)**
