// Lesson 03: Structs, methods, pointers and composition. Go has no classes.
//
// Run:  go test -v ./practice/03_structs_methods
//
// ─── CONCEPTS ──────────────────────────────────────────────────────────────
//
// 1. STRUCTS hold data; METHODS are functions with a RECEIVER:
//      type Point struct { X, Y int }
//      func (p Point) String() string { ... }   // value receiver: p is a COPY
//      func (p *Point) Move(dx int) { p.X += dx } // pointer receiver: mutates the original
//    Python's `self` is always a reference. In Go YOU choose.
//    Rule of thumb: if ANY method needs a pointer receiver (it mutates, or the struct is
//    big, or holds a mutex), make ALL methods pointer receivers for consistency.
//
// 2. EVERYTHING IS PASSED BY VALUE
//      func reset(p Point)  { p.X = 0 }   // changes a copy. The caller sees nothing.
//      func reset(p *Point) { p.X = 0 }   // changes the caller's Point
//    &x takes an address; *p dereferences. Go auto-dereferences for field access and
//    method calls: p.X works whether p is a Point or a *Point.
//
// 3. NO CONSTRUCTORS. By convention, write NewXxx functions:
//      func NewAccount(owner string) *Account { return &Account{owner: owner} }
//    Returning a pointer to a local variable is SAFE in Go (escape analysis moves it to
//    the heap). Also: design types so their ZERO VALUE is useful (sync.Mutex, bytes.Buffer,
//    strings.Builder all work without initialization).
//
// 4. VISIBILITY = CAPITALIZATION (per PACKAGE, not per type)
//      type Account struct {
//          Owner   string   // exported: visible outside the package
//          balance int      // unexported: only this package can touch it
//      }
//    There is no `_private` convention; the compiler enforces it.
//
// 5. EMBEDDING = COMPOSITION, NOT INHERITANCE
//      type Animal struct { Name string }
//      func (a Animal) Describe() string { ... }
//      type Dog struct {
//          Animal          // embedded: Dog "gets" Name and Describe() (promoted)
//          Breed string
//      }
//      d := Dog{Animal: Animal{Name: "Rex"}, Breed: "Lab"}
//      d.Name, d.Describe()      // both work
//    But a Dog IS NOT an Animal: you can't pass a Dog where an Animal is expected.
//    There's no super() and no virtual dispatch. Interfaces (lesson 04) give you
//    polymorphism instead.
//
// 6. STRUCTS ARE COMPARABLE with == if all their fields are comparable (no slices/maps).
// ───────────────────────────────────────────────────────────────────────────

package structs

import "errors"

// ── Part 1: value vs pointer ─────────────────────────────────────────────

type Point struct {
	X, Y int
}

// ScaledCopy returns a NEW point with both coordinates multiplied by k.
// The original must be unchanged (value receiver).
func (p Point) ScaledCopy(k int) Point {
	// TODO
	return Point{}
}

// Scale multiplies the point's coordinates by k IN PLACE (pointer receiver).
func (p *Point) Scale(k int) {
	// TODO
}

// ── Part 2: encapsulation and constructors ───────────────────────────────

var ErrInsufficientFunds = errors.New("insufficient funds")
var ErrInvalidAmount = errors.New("amount must be positive")

// Account is a bank account. balance is unexported so other packages can't
// set it directly. They must go through the methods.
type Account struct {
	Owner   string
	balance int
	history []int // every successful deposit (+) and withdrawal (-)
}

// NewAccount creates an account for owner with a zero balance.
func NewAccount(owner string) *Account {
	// TODO
	return nil
}

// Deposit adds amount. It returns ErrInvalidAmount if amount <= 0.
func (a *Account) Deposit(amount int) error {
	// TODO
	return nil
}

// Withdraw subtracts amount. It returns ErrInvalidAmount if amount <= 0, and
// ErrInsufficientFunds if the balance is too low (the balance stays unchanged).
func (a *Account) Withdraw(amount int) error {
	// TODO
	return nil
}

// Balance returns the current balance.
func (a *Account) Balance() int {
	// TODO
	return 0
}

// History returns a COPY of the transaction history, so callers can't mutate our
// internal slice (remember lesson 02: returning a.history would share memory!).
func (a *Account) History() []int {
	// TODO
	return nil
}

// ── Part 3: embedding ────────────────────────────────────────────────────

type Animal struct {
	Name string
	Legs int
}

// Describe returns e.g. "Rex has 4 legs".
// Hint: fmt.Sprintf("%s has %d legs", ...)
func (a Animal) Describe() string {
	// TODO
	return ""
}

// Dog embeds Animal, so it automatically has Name, Legs and Describe().
type Dog struct {
	Animal
	Breed string
}

// NewDog builds a Dog with 4 legs.
func NewDog(name, breed string) Dog {
	// TODO
	return Dog{}
}

// Speak returns e.g. "Rex the Labrador says woof".
// Notice you can use d.Name directly (promoted from Animal).
func (d Dog) Speak() string {
	// TODO
	return ""
}

// ── Part 4: a useful zero value ──────────────────────────────────────────

// Stack is a LIFO stack of ints. Its zero value must be ready to use:
//
//	var s Stack   // no constructor needed
//	s.Push(1)
type Stack struct {
	items []int
}

func (s *Stack) Push(v int) {
	// TODO
}

// Pop removes and returns the top item. ok is false if the stack is empty.
func (s *Stack) Pop() (v int, ok bool) {
	// TODO
	return 0, false
}

// Peek returns the top item without removing it.
func (s *Stack) Peek() (v int, ok bool) {
	// TODO
	return 0, false
}

func (s *Stack) Len() int {
	// TODO
	return 0
}
