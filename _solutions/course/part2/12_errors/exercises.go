// Reference solution for lesson 12.
package lesson12

import (
	"errors"
	"fmt"
	"strconv"
)

func ParseAge(s string) (int, error) {
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("not a number: %q", s)
	}
	if n < 0 {
		return 0, errors.New("age cannot be negative")
	}
	return n, nil
}

var ErrDivideByZero = errors.New("division by zero")

func SafeDivide(a, b int) (int, error) {
	if b == 0 {
		return 0, ErrDivideByZero
	}
	return a / b, nil
}

func SumStrings(values []string) (int, error) {
	total := 0
	for _, v := range values {
		n, err := strconv.Atoi(v)
		if err != nil {
			return 0, err
		}
		total += n
	}
	return total, nil
}

var (
	ErrInvalidAmount     = errors.New("amount must be positive")
	ErrInsufficientFunds = errors.New("insufficient funds")
)

func Withdraw(balance, amount int) (int, error) {
	if amount <= 0 {
		return balance, ErrInvalidAmount
	}
	if amount > balance {
		return balance, ErrInsufficientFunds
	}
	return balance - amount, nil
}

func WithdrawMessage(balance, amount int) string {
	newBalance, err := Withdraw(balance, amount)
	switch {
	case errors.Is(err, ErrInsufficientFunds):
		return "sorry, not enough money"
	case errors.Is(err, ErrInvalidAmount):
		return "please enter a positive amount"
	case err != nil:
		return "unexpected error: " + err.Error()
	}
	return fmt.Sprintf("new balance: %d", newBalance)
}

func ParseAll(inputs []string) (nums []int, bad []string) {
	for _, s := range inputs {
		n, err := strconv.Atoi(s)
		if err != nil {
			bad = append(bad, s)
			continue
		}
		nums = append(nums, n)
	}
	return nums, bad
}

var ErrNotFound = errors.New("user not found")

func FindUser(users map[int]string, id int) (string, error) {
	name, ok := users[id]
	if !ok {
		return "", ErrNotFound
	}
	return name, nil
}

func UserName(users map[int]string, id int) string {
	name, err := FindUser(users, id)
	if err != nil {
		return "guest"
	}
	return name
}
