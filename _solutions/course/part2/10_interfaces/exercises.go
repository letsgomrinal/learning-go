// Reference solution for lesson 10. Compare with yours AFTER your tests pass.
package lesson10

import (
	"math"
	"strconv"
	"strings"
)

type Shape interface {
	Area() float64
	Perimeter() float64
}

type Square struct {
	Side float64
}

func (s Square) Area() float64      { return s.Side * s.Side }
func (s Square) Perimeter() float64 { return 4 * s.Side }

type Circle struct {
	R float64
}

func (c Circle) Area() float64      { return math.Pi * c.R * c.R }
func (c Circle) Perimeter() float64 { return 2 * math.Pi * c.R }

// Compile-time checks.
var (
	_ Shape = Square{}
	_ Shape = Circle{}
)

func TotalArea(shapes []Shape) float64 {
	total := 0.0
	for _, s := range shapes {
		total += s.Area()
	}
	return total
}

func Biggest(shapes []Shape) (Shape, bool) {
	if len(shapes) == 0 {
		return nil, false
	}
	best := shapes[0]
	for _, s := range shapes {
		if s.Area() > best.Area() {
			best = s
		}
	}
	return best, true
}

type Speaker interface {
	Speak() string
}

type Dog struct{}

func (d Dog) Speak() string { return "Woof" }

type Robot struct {
	spoken int
}

func (r *Robot) Speak() string {
	r.spoken++
	return "Beep"
}

func Chorus(speakers []Speaker) string {
	var words []string
	for _, s := range speakers {
		words = append(words, s.Speak())
	}
	return strings.Join(words, " ")
}

type CreditCard struct {
	Last4 string
}

func (c CreditCard) Pay(amount int) string {
	return "charged $" + itoa(amount) + " to card ending " + c.Last4
}

type PayPal struct {
	Email string
}

func (p PayPal) Pay(amount int) string {
	return "sent $" + itoa(amount) + " from " + p.Email
}

type Payer interface {
	Pay(amount int) string
}

func Checkout(p Payer, amount int) string {
	return "Order paid: " + p.Pay(amount)
}

func itoa(n int) string { return strconv.Itoa(n) }
