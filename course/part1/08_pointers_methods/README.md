# Lesson 08: Pointers and methods

**Goal:** understand what a pointer is (from scratch), use pointers to let a function change your
data, and attach methods to your own types.

**Time:** 1.5 hours. Take it slowly; pointers are new if you've only written Python.

```bash
cd ~/Desktop/learning-go/course/part1/08_pointers_methods
```

---

## 1. The problem pointers solve

From lesson 07: Go copies values when you pass them to a function. So this doesn't work:

```go
func birthday(age int) {
	age = age + 1          // changes the function's own copy
}

func main() {
	myAge := 30
	birthday(myAge)
	fmt.Println(myAge)     // 30, not 31
}
```

Python never has this problem because it always passes references to objects. Go makes you choose:
pass a copy (the default), or pass a **pointer**.

---

## 2. What a pointer is

Every variable lives somewhere in the computer's memory, and that place has an **address**.
A pointer is simply a value that holds an address: it says *where* a variable lives.

Think of it this way: a variable is a house, its value is what's inside, and a pointer is a piece
of paper with the house's street address written on it. If you give someone the address, they can
go to the house and rearrange the furniture. If you give them a photo of the inside (a copy), they can't.

Two operators:

| Operator | Name | Meaning | Example |
|---|---|---|---|
| `&x` | "address of" | give me a pointer to `x` | `p := &age` |
| `*p` | "dereference" | the value at the address `p` holds | `*p = 31` |

```go
age := 30
p := &age              // p holds the address of age

fmt.Println(p)         // 0x14000012078 (some address; yours will differ)
fmt.Println(*p)        // 30  (the value at that address)

*p = 31                // change the value AT that address
fmt.Println(age)       // 31  ← age changed, through the pointer
```

### Pointer types

The type of `&age` is `*int`, read "pointer to int". A `*` before a **type** means "pointer to".
A `*` before a **variable** means "the value it points at". Same symbol, two jobs:

```go
var p *int     // p is a pointer to an int  (a * in a type)
*p = 5         // set the int that p points at  (a * on a variable)
```

---

## 3. Fixing `birthday` with a pointer

```go
func birthday(age *int) {      // takes a pointer to an int
	*age = *age + 1            // change the value at that address
}

func main() {
	myAge := 30
	birthday(&myAge)           // pass the address of myAge
	fmt.Println(myAge)         // 31 ✅
}
```

Read it step by step:
1. `&myAge` gives the function the address of `myAge`.
2. Inside, `age` is a pointer (`*int`).
3. `*age = *age + 1` means "read the value at the address, add 1, store it back at the address".

### 🤔 Predict the output

```go
func swap(a, b *int) {
	*a, *b = *b, *a
}

func main() {
	x, y := 1, 2
	swap(&x, &y)
	fmt.Println(x, y)
}
```

<details><summary>Answer</summary>

`2 1`. `swap` receives the addresses of `x` and `y`, and swaps the values stored there.
</details>

---

## 4. Pointers to structs

Pointers are used most with structs, because you often want a function to change a struct:

```go
type Player struct {
	Name  string
	Score int
}

func addPoints(p *Player, points int) {
	p.Score += points         // Go lets you write p.Score instead of (*p).Score
}

func main() {
	ana := Player{Name: "Ana"}
	addPoints(&ana, 10)
	fmt.Println(ana.Score)    // 10 ✅
}
```

Notice `p.Score`: with a pointer to a struct, Go automatically follows the pointer for you.
You don't have to write `(*p).Score`.

You can also create a struct and get a pointer to it in one step:

```go
p := &Player{Name: "Bo"}      // p is a *Player
p.Score = 5
```

### `nil`: a pointer to nothing

The zero value of a pointer is `nil`: it points nowhere. Using it crashes:

```go
var p *Player                 // nil
fmt.Println(p == nil)         // true
p.Score = 1                   // panic: runtime error: invalid memory address or nil pointer dereference
```

That error message is famous. When you see it, some pointer was `nil` when you used it. Check
with `if p != nil` when a pointer might be `nil`.

---

## 5. Methods

A **method** is a function attached to a type. In Python, methods live inside the class:

```python
class Rectangle:
    def __init__(self, w, h):
        self.w = w
        self.h = h

    def area(self):
        return self.w * self.h
```

In Go, you define the struct, and then write methods separately, with a **receiver** in front of
the name:

```go
type Rectangle struct {
	W, H float64
}

func (r Rectangle) Area() float64 {
	return r.W * r.H
}

func main() {
	rect := Rectangle{W: 3, H: 4}
	fmt.Println(rect.Area())     // 12
}
```

`(r Rectangle)` is the receiver. It plays the role of Python's `self`, but you choose the name.
Go convention: a short name, often the first letter of the type (`r`, `p`, `b`). Never `self` or `this`.

A method is really just a function with a special first parameter. `rect.Area()` is roughly the same as
calling `Area(rect)`.

---

## 6. Value receivers vs pointer receivers

This is where pointers and methods meet.

A **value receiver** `(r Rectangle)` gets a **copy**. It can read the fields but changes are lost:

```go
func (r Rectangle) DoubleBroken() {
	r.W *= 2                 // changes the copy only
}
```

A **pointer receiver** `(r *Rectangle)` gets the address, so it can change the original:

```go
func (r *Rectangle) Double() {
	r.W *= 2
	r.H *= 2
}

func main() {
	rect := Rectangle{W: 3, H: 4}
	rect.DoubleBroken()
	fmt.Println(rect)        // {3 4}  (unchanged)
	rect.Double()
	fmt.Println(rect)        // {6 8}  ✅
}
```

Notice you call `rect.Double()`, not `(&rect).Double()`. Go takes the address for you automatically.

**How to choose:**
- Method needs to **change** the struct → pointer receiver `(r *Rectangle)`.
- Method only **reads** → a value receiver works.
- If **any** method of a type uses a pointer receiver, the Go convention is to make **all** of
  them pointer receivers, for consistency. Many Go programmers just use pointer receivers for
  structs by default, and that's a fine habit while learning.

### 🤔 Predict the output

```go
type Counter struct{ n int }

func (c Counter) IncValue()    { c.n++ }
func (c *Counter) IncPointer() { c.n++ }

func main() {
	var c Counter
	c.IncValue()
	c.IncPointer()
	c.IncPointer()
	c.IncValue()
	fmt.Println(c.n)
}
```

<details><summary>Answer</summary>

`2`. Only the pointer-receiver calls change `c`. (Also notice `var c Counter` works without
any constructor: the zero value `{n: 0}` is ready to use.)
</details>

---

## 7. Constructors return pointers

When a type is meant to be changed by its methods, its `New...` function usually returns a pointer:

```go
type Account struct {
	Owner   string
	Balance int
}

func NewAccount(owner string) *Account {
	return &Account{Owner: owner}
}

acc := NewAccount("Ana")     // acc is a *Account
acc.Balance += 100
```

Returning the address of a value created inside a function is completely safe in Go (unlike C).
Go notices and keeps the value alive as long as something points to it.

---

## 8. Methods on any type you define

Methods aren't only for structs. You can add them to any type you create:

```go
type Celsius float64

func (c Celsius) ToFahrenheit() float64 {
	return float64(c)*9/5 + 32
}

temp := Celsius(100)
fmt.Println(temp.ToFahrenheit())   // 212
```

`type Celsius float64` creates a new type based on `float64`. This is a common way to give plain
numbers a meaning (and stop you mixing up Celsius and Fahrenheit by accident).

---

## 9. When to use pointers: summary

| Use a pointer when... | Example |
|---|---|
| a function or method must **change** the caller's value | `func (a *Account) Deposit(n int)` |
| the struct is large and copying it every call is wasteful | big config structs |
| "no value" is a meaningful state (`nil`) | an optional result |

Don't use pointers for small values that you only read (ints, small structs). Copies are cheap and
safer. And you rarely need pointers to slices or maps: they already share their data.

---

## 10. Common mistakes

| You wrote | What happens | Fix |
|---|---|---|
| Method that changes fields, with a value receiver | changes are lost silently | `(r *Type)` |
| `var p *Player` then `p.Score = 1` | panic: nil pointer dereference | `p := &Player{}` |
| `birthday(myAge)` when it takes `*int` | `cannot use myAge (variable of type int) as *int value` | `birthday(&myAge)` |
| `func (self Rectangle)` | works, but not Go style | `func (r Rectangle)` |

---

## 11. Exercises

| | Exercise | Practises |
|---|---|---|
| Ex1 | `Increment(n *int)` | `*` to change a value |
| Ex2 | `SwapInts(a, b *int)` | two pointers |
| Ex3 | `Point.Move` | a function taking `*Point` |
| Ex4 | `Rectangle.Area` and `Perimeter` methods | value receivers |
| Ex5 | `Rectangle.Scale` | a pointer receiver |
| Ex6 | `Counter` with `Inc`, `Add`, `Value`, `Reset` | a small type with methods |
| Ex7 | `Account` with `Deposit` and `Withdraw` | methods that validate |
| Ex8 | `Celsius.ToFahrenheit` | methods on a non-struct type |

---

## ✅ Checklist

- [ ] I can explain what `&x` and `*p` do, in my own words
- [ ] I can write a function that changes the caller's variable using a pointer
- [ ] I know the difference between `(r Rectangle)` and `(r *Rectangle)` receivers
- [ ] I know what "nil pointer dereference" means
- [ ] I know when a pointer is the right choice, and when a copy is better

**Next → [Lesson 09: Packages and tests](../09_packages_tests/README.md)**
