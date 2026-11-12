// Reference solution for lesson 03.
package structs

import (
	"errors"
	"fmt"
	"slices"
)

type Point struct {
	X, Y int
}

func (p Point) ScaledCopy(k int) Point {
	p.X *= k // p is a copy; modifying it is fine
	p.Y *= k
	return p
}

func (p *Point) Scale(k int) {
	p.X *= k
	p.Y *= k
}

var ErrInsufficientFunds = errors.New("insufficient funds")
var ErrInvalidAmount = errors.New("amount must be positive")

type Account struct {
	Owner   string
	balance int
	history []int
}

func NewAccount(owner string) *Account {
	return &Account{Owner: owner}
}

func (a *Account) Deposit(amount int) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}
	a.balance += amount
	a.history = append(a.history, amount)
	return nil
}

func (a *Account) Withdraw(amount int) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}
	if amount > a.balance {
		return ErrInsufficientFunds
	}
	a.balance -= amount
	a.history = append(a.history, -amount)
	return nil
}

func (a *Account) Balance() int { return a.balance }

func (a *Account) History() []int { return slices.Clone(a.history) }

type Animal struct {
	Name string
	Legs int
}

func (a Animal) Describe() string {
	return fmt.Sprintf("%s has %d legs", a.Name, a.Legs)
}

type Dog struct {
	Animal
	Breed string
}

func NewDog(name, breed string) Dog {
	return Dog{Animal: Animal{Name: name, Legs: 4}, Breed: breed}
}

func (d Dog) Speak() string {
	return fmt.Sprintf("%s the %s says woof", d.Name, d.Breed)
}

type Stack struct {
	items []int
}

func (s *Stack) Push(v int) { s.items = append(s.items, v) }

func (s *Stack) Pop() (v int, ok bool) {
	if len(s.items) == 0 {
		return 0, false
	}
	v = s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return v, true
}

func (s *Stack) Peek() (v int, ok bool) {
	if len(s.items) == 0 {
		return 0, false
	}
	return s.items[len(s.items)-1], true
}

func (s *Stack) Len() int { return len(s.items) }
