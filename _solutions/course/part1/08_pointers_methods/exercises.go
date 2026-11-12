// Reference solution for lesson 08.
package lesson08

func Increment(n *int) {
	*n = *n + 1 // or: *n++
}

func SwapInts(a, b *int) {
	*a, *b = *b, *a
}

type Point struct {
	X, Y int
}

func Move(p *Point, dx, dy int) {
	p.X += dx
	p.Y += dy
}

type Rectangle struct {
	W, H float64
}

func (r Rectangle) Area() float64 {
	return r.W * r.H
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.W + r.H)
}

func (r *Rectangle) Scale(factor float64) {
	r.W *= factor
	r.H *= factor
}

type Counter struct {
	count int
}

func (c *Counter) Inc()       { c.count++ }
func (c *Counter) Add(n int)  { c.count += n }
func (c *Counter) Value() int { return c.count }
func (c *Counter) Reset()     { c.count = 0 }

type Account struct {
	Owner   string
	Balance int
}

func NewAccount(owner string) *Account {
	return &Account{Owner: owner}
}

func (a *Account) Deposit(amount int) {
	if amount <= 0 {
		return
	}
	a.Balance += amount
}

func (a *Account) Withdraw(amount int) bool {
	if amount <= 0 || amount > a.Balance {
		return false
	}
	a.Balance -= amount
	return true
}

type Celsius float64

func (c Celsius) ToFahrenheit() float64 {
	return float64(c)*9/5 + 32
}
