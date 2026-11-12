// Lesson 05: Errors are values. Goodbye try/except.
//
// Run:  go test -v ./practice/05_errors
//
// ─── CONCEPTS ──────────────────────────────────────────────────────────────
//
// 1. error IS JUST AN INTERFACE
//      type error interface { Error() string }
//    Functions that can fail return it as their LAST value:
//      n, err := strconv.Atoi(s)
//      if err != nil {
//          return 0, err          // handle it or pass it up. Never ignore it silently.
//      }
//    It's verbose, and that's on purpose: every failure path is visible in the code.
//
// 2. THREE KINDS OF ERRORS
//    a) Sentinel errors: package-level values you compare against
//         var ErrNotFound = errors.New("not found")
//         if errors.Is(err, ErrNotFound) { ... }       // like `except NotFound:`
//    b) Custom error types: carry structured data
//         type ValidationError struct { Field string }
//         func (e *ValidationError) Error() string { ... }
//         var ve *ValidationError
//         if errors.As(err, &ve) { fmt.Println(ve.Field) }   // like `except X as e:`
//    c) Ad-hoc: fmt.Errorf("bad input %q", s)
//
// 3. WRAPPING adds context while keeping the original (like `raise X from e`):
//      return fmt.Errorf("load config %s: %w", path, err)   // %w wraps; %v doesn't
//    errors.Is / errors.As look through the whole chain of wrapped errors.
//    The message reads like a breadcrumb trail:
//      "start server: load config app.yaml: open app.yaml: no such file or directory"
//    Convention: lowercase messages, no trailing punctuation.
//
// 4. defer runs when the surrounding function RETURNS (LIFO order). It's Go's `finally`/`with`:
//      f, err := os.Open(p); if err != nil { return err }
//      defer f.Close()
//
// 5. panic/recover are for programmer bugs (index out of range, nil map write) and truly
//    unrecoverable states, NOT for normal control flow. recover() only works inside a
//    deferred function. Libraries should almost never let a panic escape.
//
// 6. THE NIL INTERFACE TRAP (a classic Go interview question)
//    An interface value is (type, value). It's nil only when BOTH are nil.
//      var p *MyErr = nil
//      var err error = p     // err is (*MyErr, nil), which is NOT nil!
//      err != nil            // true!!
//    So never return a typed nil pointer as an error. Return a literal `nil`.
// ───────────────────────────────────────────────────────────────────────────

package errs

import (
	"errors"
)

// ── Part 1: sentinel errors and wrapping ─────────────────────────────────

var ErrInsufficientFunds = errors.New("insufficient funds")

// Withdraw returns balance-amount. If amount > balance it returns an error that
// WRAPS ErrInsufficientFunds and has the exact message:
//
//	"withdraw 150 from 100: insufficient funds"
//
// Hint: fmt.Errorf with %d and %w.
func Withdraw(balance, amount int) (int, error) {
	// TODO
	return 0, nil
}

// ── Part 2: custom error types ───────────────────────────────────────────

// ValidationError reports which field was invalid and why.
type ValidationError struct {
	Field  string
	Reason string
}

// Error implements the error interface. Format: "invalid <field>: <reason>"
func (e *ValidationError) Error() string {
	// TODO
	return ""
}

// ParseAge parses an age string.
//   - Not a number → return an error that WRAPS the strconv error, with the message
//     prefix `parse age "abc": ` (hint: %q quotes a string). errors.Is(err, strconv.ErrSyntax)
//     must be true.
//   - Negative or > 150 → return &ValidationError{Field: "age", Reason: "out of range"}
func ParseAge(s string) (int, error) {
	// TODO
	return 0, nil
}

// RegisterUser validates name and age. It returns:
//   - &ValidationError{Field: "name", Reason: "empty"} if name is ""
//   - the error from ParseAge wrapped as "register <name>: <err>" otherwise, if any
//
// The caller should still be able to errors.As(err, &ve) through the wrapping.
func RegisterUser(name, age string) error {
	// TODO
	return nil
}

// ── Part 3: the nil interface trap ───────────────────────────────────────

type NotFoundError struct{ ID int }

func (e *NotFoundError) Error() string { return "not found" }

// FindUser should return nil for valid ids (>= 0) and a *NotFoundError otherwise.
// BUG: this code has the nil-interface bug. The test fails. Fix it.
func FindUser(id int) error {
	var e *NotFoundError
	if id < 0 {
		e = &NotFoundError{ID: id}
	}
	return e
}

// ── Part 4: defer, panic, recover ────────────────────────────────────────

// SafeDivide divides a by b. Integer division by zero PANICS in Go. Instead of
// checking b == 0, practise recover: use a deferred function that recovers from the
// panic and sets the named result err to fmt.Errorf("recovered: %v", r).
//
//	func SafeDivide(a, b int) (res int, err error) {
//	    defer func() {
//	        if r := recover(); r != nil { ... }
//	    }()
//	    ...
//	}
//
// The deferred func can modify `err` because it's a NAMED result.
func SafeDivide(a, b int) (res int, err error) {
	// TODO
	return 0, nil
}

// DeferOrder returns the order in which things ran. Implement it so it
// appends "start", then defers appending "deferred 1", "deferred 2", "deferred 3"
// (in that source order), then appends "end".
// Predict the final slice before you run the test! (Hint: LIFO.)
//
// Note: a deferred func runs AFTER the return value is evaluated, so return the
// log through a pointer or a named result. Here we use a named result.
func DeferOrder() (log []string) {
	// TODO
	return nil
}
