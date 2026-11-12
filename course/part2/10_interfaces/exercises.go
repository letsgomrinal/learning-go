// Lesson 10 exercises: interfaces. Read README.md first!
//
// Check one exercise at a time:   go test -run Ex1 -v
package lesson10

// Shape is anything with an area and a perimeter.
type Shape interface {
	Area() float64
	Perimeter() float64
}

// Ex1: make Square a Shape by giving it two methods:
//
//	func (s Square) Area() float64       → Side × Side
//	func (s Square) Perimeter() float64  → 4 × Side
//
// Write them below the type. Nothing else is needed: no "implements".
type Square struct {
	Side float64
}

// Ex2: make Circle a Shape too.
// Area is π × R × R and Perimeter is 2 × π × R. Use math.Pi (import "math").
type Circle struct {
	R float64
}

// Ex3: return the sum of the areas of all the shapes.
func TotalArea(shapes []Shape) float64 {
	// TODO
	return 0
}

// Ex4: return the shape with the biggest area, and true.
// For an empty slice return nil and false. (nil is the zero value of an interface.)
// On a tie, return the first one.
func Biggest(shapes []Shape) (Shape, bool) {
	// TODO
	return nil, false
}

// ── Ex5 ──────────────────────────────────────────────────────────────────

// Speaker is anything that can speak.
type Speaker interface {
	Speak() string
}

// Dog: give it a Speak method that returns "Woof". (Use a value receiver.)
type Dog struct{}

// Robot is already a Speaker, but notice its POINTER receiver: it counts how many times it spoke.
// So only *Robot (a pointer) is a Speaker, not Robot. The tests use &Robot{}.
type Robot struct {
	spoken int
}

func (r *Robot) Speak() string {
	r.spoken++
	return "Beep"
}

// Chorus calls Speak on every speaker and joins the results with a single space.
// Chorus([]Speaker{Dog{}, &Robot{}, Dog{}}) == "Woof Beep Woof"
//
// Hint: collect the words in a []string, then strings.Join(words, " ").
func Chorus(speakers []Speaker) string {
	// TODO
	return ""
}

// ── Ex6: design your own interface ───────────────────────────────────────

// CreditCard and PayPal both know how to pay. They were written separately and know nothing
// about each other.
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

// Payer should describe "anything that can Pay(amount int) string".
// TODO: add that method to the interface. (Right now it's empty, which means ANY type fits.)
type Payer interface {
}

// Checkout pays the amount with p and returns "Order paid: " followed by what Pay returned.
// Checkout(PayPal{Email: "ana@example.com"}, 20) == "Order paid: sent $20 from ana@example.com"
//
// This won't compile until Payer has the Pay method. Try it before you add it, and read the error!
func Checkout(p Payer, amount int) string {
	// TODO
	return ""
}

// itoa is a tiny helper (it's strconv.Itoa), kept here so this file needs no extra imports.
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var digits []byte
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	if neg {
		return "-" + string(digits)
	}
	return string(digits)
}
