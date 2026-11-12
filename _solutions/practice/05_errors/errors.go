// Reference solution for lesson 05.
package errs

import (
	"errors"
	"fmt"
	"strconv"
)

var ErrInsufficientFunds = errors.New("insufficient funds")

func Withdraw(balance, amount int) (int, error) {
	if amount > balance {
		return 0, fmt.Errorf("withdraw %d from %d: %w", amount, balance, ErrInsufficientFunds)
	}
	return balance - amount, nil
}

type ValidationError struct {
	Field  string
	Reason string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("invalid %s: %s", e.Field, e.Reason)
}

func ParseAge(s string) (int, error) {
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("parse age %q: %w", s, err)
	}
	if n < 0 || n > 150 {
		return 0, &ValidationError{Field: "age", Reason: "out of range"}
	}
	return n, nil
}

func RegisterUser(name, age string) error {
	if name == "" {
		return &ValidationError{Field: "name", Reason: "empty"}
	}
	if _, err := ParseAge(age); err != nil {
		return fmt.Errorf("register %s: %w", name, err)
	}
	return nil
}

type NotFoundError struct{ ID int }

func (e *NotFoundError) Error() string { return "not found" }

// Fixed: return an untyped nil on the success path.
func FindUser(id int) error {
	if id < 0 {
		return &NotFoundError{ID: id}
	}
	return nil
}

func SafeDivide(a, b int) (res int, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("recovered: %v", r)
		}
	}()
	return a / b, nil
}

func DeferOrder() (log []string) {
	log = append(log, "start")
	for i := 1; i <= 3; i++ {
		defer func() { log = append(log, fmt.Sprintf("deferred %d", i)) }()
	}
	log = append(log, "end")
	return log
}
