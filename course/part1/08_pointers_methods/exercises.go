// Lesson 08 exercises: pointers and methods. Read README.md first!
//
// Check one exercise at a time:   go test -run Ex1 -v
package lesson08

// Ex1: add 1 to the int that n points at.
//
//	x := 5
//	Increment(&x)
//	// x is now 6
func Increment(n *int) {
	// TODO
}

// Ex2: swap the values that a and b point at.
//
//	x, y := 1, 2
//	SwapInts(&x, &y)
//	// x is now 2, y is now 1
func SwapInts(a, b *int) {
	// TODO
}

type Point struct {
	X, Y int
}

// Ex3: move the point by dx and dy. The caller's Point must change.
//
//	p := Point{X: 1, Y: 1}
//	Move(&p, 2, 3)
//	// p is now {3 4}
//
// Remember: with a pointer to a struct you can write p.X directly.
func Move(p *Point, dx, dy int) {
	// TODO
}

type Rectangle struct {
	W, H float64
}

// Ex4: the methods Area and Perimeter only READ the rectangle, so they use value receivers.

// Area returns W × H.
func (r Rectangle) Area() float64 {
	// TODO
	return 0
}

// Perimeter returns 2 × (W + H).
func (r Rectangle) Perimeter() float64 {
	// TODO
	return 0
}

// Ex5: multiply both sides by factor, changing the rectangle itself.
// Look at the receiver: (r *Rectangle). That's what makes the change visible to the caller.
func (r *Rectangle) Scale(factor float64) {
	// TODO
}

// Ex6: a counter. Its zero value must work:
//
//	var c Counter
//	c.Inc()
//	c.Add(5)
//	c.Value() == 6
//	c.Reset()
//	c.Value() == 0
//
// The field is lowercase (count), which keeps it private to this package. Lesson 09 explains.
type Counter struct {
	count int
}

// Inc adds 1.
func (c *Counter) Inc() {
	// TODO
}

// Add adds n.
func (c *Counter) Add(n int) {
	// TODO
}

// Value returns the current count.
func (c *Counter) Value() int {
	// TODO
	return 0
}

// Reset sets the count back to 0.
func (c *Counter) Reset() {
	// TODO
}

// Ex7: a bank account.
type Account struct {
	Owner   string
	Balance int
}

// NewAccount returns a POINTER to a new Account for owner, with a balance of 0.
func NewAccount(owner string) *Account {
	// TODO
	return nil
}

// Deposit adds amount to the balance. Ignore amounts of 0 or less (do nothing).
func (a *Account) Deposit(amount int) {
	// TODO
}

// Withdraw removes amount from the balance and returns true.
// If amount is 0 or less, or bigger than the balance, change nothing and return false.
func (a *Account) Withdraw(amount int) bool {
	// TODO
	return false
}

// Ex8: a method on a non-struct type.
type Celsius float64

// ToFahrenheit converts: F = C × 9 / 5 + 32
// Hint: convert c to a float64 first: float64(c)
func (c Celsius) ToFahrenheit() float64 {
	// TODO
	return 0
}
