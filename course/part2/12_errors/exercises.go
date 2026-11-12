// Lesson 12 exercises: errors. Read README.md first!
//
// Check one exercise at a time:   go test -run Ex1 -v
package lesson12

import (
	"errors"
)

// Ex1: parse an age.
//   - If s isn't a whole number, return 0 and an error with the message: not a number: "abc"
//     (that's fmt.Errorf("not a number: %q", s))
//   - If the number is negative, return 0 and the error message: age cannot be negative
//   - Otherwise return the number and nil.
//
// You'll need strconv.Atoi and the "fmt" and "strconv" imports.
func ParseAge(s string) (int, error) {
	// TODO
	return 0, nil
}

// ErrDivideByZero is a sentinel error.
var ErrDivideByZero = errors.New("division by zero")

// Ex2: return a / b. If b is 0, return 0 and ErrDivideByZero (the variable above, not a new error).
func SafeDivide(a, b int) (int, error) {
	// TODO
	return 0, nil
}

// Ex3: convert every string to a number and return the sum.
// If any string isn't a number, stop and return 0 and the error that strconv.Atoi gave you
// (pass it up unchanged).
// SumStrings([]string{"1", "2", "3"}) == 6, nil
func SumStrings(values []string) (int, error) {
	// TODO
	return 0, nil
}

var (
	ErrInvalidAmount     = errors.New("amount must be positive")
	ErrInsufficientFunds = errors.New("insufficient funds")
)

// Ex4: take amount out of balance and return the new balance.
//   - amount <= 0       → return balance (unchanged) and ErrInvalidAmount
//   - amount > balance  → return balance (unchanged) and ErrInsufficientFunds
//   - otherwise         → return balance - amount and nil
func Withdraw(balance, amount int) (int, error) {
	// TODO
	return 0, nil
}

// Ex5: call Withdraw and turn the result into a message for the user:
//   - success                    → "new balance: 70"   (with the real number)
//   - ErrInsufficientFunds       → "sorry, not enough money"
//   - ErrInvalidAmount           → "please enter a positive amount"
//
// Use errors.Is (README section 5). A switch reads nicely here.
func WithdrawMessage(balance, amount int) string {
	// TODO
	return ""
}

// Ex6: parse every string. Return the numbers that parsed, AND the strings that didn't,
// both in their original order. Don't stop at the first bad one.
// ParseAll([]string{"1", "x", "3", "y"}) == [1 3], ["x" "y"]
func ParseAll(inputs []string) (nums []int, bad []string) {
	// TODO
	return nil, nil
}

// ErrNotFound means no user has that id.
var ErrNotFound = errors.New("user not found")

// FindUser is already written: it returns the user's name, or ErrNotFound.
func FindUser(users map[int]string, id int) (string, error) {
	name, ok := users[id]
	if !ok {
		return "", ErrNotFound
	}
	return name, nil
}

// Ex7: return the user's name, or "guest" if the user isn't found.
// Use FindUser and handle the error by using a default (README section 4b).
func UserName(users map[int]string, id int) string {
	// TODO
	return ""
}
