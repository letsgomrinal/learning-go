# Lesson 10: Interfaces

**Goal:** understand what an interface is, why Go uses them instead of inheritance, and write
functions that work with any type that has the right methods.

**Time:** 1.5 hours. Interfaces are the most important idea in Part 2.

```bash
cd ~/Desktop/learning-go/course/part2/10_interfaces
```

---

## 1. The problem

Say you have two shapes, each with an `Area` method (lesson 08):

```go
type Rect struct{ W, H float64 }
type Circle struct{ R float64 }

func (r Rect) Area() float64   { return r.W * r.H }
func (c Circle) Area() float64 { return 3.14159 * c.R * c.R }
```

Now you want a function that adds up the area of a list of shapes, rects and circles mixed.
What type does the list have? `[]Rect` can only hold rects. `[]Circle` can only hold circles.

In Python you wouldn't even think about it. You'd write:

```python
def total_area(shapes):
    return sum(s.area() for s in shapes)
```

and it works for anything with an `.area()` method. That's **duck typing**: "if it walks like a duck
and quacks like a duck, it's a duck". The catch: if you pass something without `.area()`, you only
find out when the program crashes.

Go gives you duck typing **checked by the compiler**. The tool for it is an **interface**.

---

## 2. Defining an interface

An interface is a type that lists **methods**, without any code:

```go
type Shape interface {
	Area() float64
}
```

Read it as: "a `Shape` is anything that has an `Area() float64` method".

Now `Rect` and `Circle` are both `Shape`s, **automatically**. You don't write `implements Shape`
anywhere. Go sees that they have an `Area() float64` method, and that's enough.

```go
func totalArea(shapes []Shape) float64 {
	total := 0.0
	for _, s := range shapes {
		total += s.Area()
	}
	return total
}

func main() {
	shapes := []Shape{
		Rect{W: 2, H: 3},
		Circle{R: 1},
	}
	fmt.Println(totalArea(shapes))   // 9.14159
}
```

`[]Shape` holds rects and circles side by side, because both are Shapes. Inside `totalArea`, all
you know about `s` is that it has an `Area()` method, and that's all you need.

### 🤔 Try it

Type the code above into `try/main.go` and run it. Then add a third type:

```go
type Square struct{ Side float64 }
```

and put `Square{Side: 2}` in the slice **without** giving it an `Area` method. Read the error:

```
cannot use Square{…} (value of struct type Square) as Shape value in array or slice literal:
Square does not implement Shape (missing method Area)
```

That's the compiler doing duck typing for you, **before** the program runs. Now add the method:

```go
func (s Square) Area() float64 { return s.Side * s.Side }
```

and it compiles.

---

## 3. Implicit is the point

In Java or C#, a class must say `implements Shape`. Python's abstract base classes are similar:
`class Rect(Shape):`. Go doesn't need any of that. This sounds like a small detail, but it changes
how you design code:

- **You can write an interface for types you didn't write.** If some library has a type with a
  `Save()` method, you can define `type Saver interface { Save() error }` in your own code, and the
  library's type satisfies it, without the library knowing.
- **Interfaces are defined where they're used**, by the code that needs the behaviour, not by the
  types that provide it.

---

## 4. Interfaces with several methods

```go
type Shape interface {
	Area() float64
	Perimeter() float64
}
```

Now a type must have **both** methods to be a Shape. If it only has `Area`, the error says
`missing method Perimeter`.

---

## 5. What you can and can't do with an interface value

When a variable's type is an interface, you can **only call the interface's methods**:

```go
var s Shape = Rect{W: 2, H: 3}
fmt.Println(s.Area())    // fine
fmt.Println(s.W)         // compile error: s.W undefined (type Shape has no field or method W)
```

Even though `s` holds a `Rect`, the compiler only knows it's "some Shape". (Lesson 11 shows how
to get the `Rect` back out when you really need it.)

### 🤔 Predict the output

```go
type Greeter interface{ Greet() string }

type English struct{}
type Spanish struct{ Name string }

func (English) Greet() string   { return "Hello" }
func (s Spanish) Greet() string { return "Hola " + s.Name }

func main() {
	gs := []Greeter{English{}, Spanish{Name: "Ana"}, English{}}
	for _, g := range gs {
		fmt.Print(g.Greet(), " | ")
	}
}
```

<details markdown="1">
<summary>Answer</summary>

`Hello | Hola Ana | Hello | `. Each element calls its own type's `Greet`. (Also notice
`func (English) Greet()`: when a method doesn't use its receiver, you can leave out the name.
And `English{}` is a struct with no fields at all, which is perfectly fine.)
</details>

---

## 6. ⚠️ Pointer receivers and interfaces

This trips up everyone once. If a method has a **pointer receiver**, only the **pointer** satisfies
the interface:

```go
type Counter struct{ n int }

func (c *Counter) Area() float64 { return float64(c.n) }   // pointer receiver

var s Shape = Counter{}    // compile error!
// Counter does not implement Shape (method Area has pointer receiver)

var s Shape = &Counter{}   // fine: *Counter has the method
```

**Rule:** if the methods use `(c *Counter)`, put `&Counter{...}` in the interface. If they use
`(c Counter)`, both `Counter{}` and `&Counter{}` work.

When you see "method has pointer receiver" in an error, add a `&`.

---

## 7. Checking at compile time

A common trick to make the compiler confirm a type satisfies an interface, even if nothing uses it yet:

```go
var _ Shape = Rect{}
```

It creates a variable named `_` (thrown away) of type `Shape`, from a `Rect`. If `Rect` is missing a
method, this line fails to compile. You'll see this in real Go code near type definitions.

---

## 8. Small interfaces are best

Go programmers keep interfaces **small**, often just **one method**. The standard library is full of them:

```go
type Stringer interface { String() string }                       // package fmt
type Reader interface { Read(p []byte) (n int, err error) }       // package io
type Writer interface { Write(p []byte) (n int, err error) }      // package io
```

One-method interfaces are usually named after the method plus "er": `Reader`, `Writer`, `Stringer`,
`Closer`. A small interface is easy to satisfy, so many types can be used with your function.

A famous Go proverb: **"The bigger the interface, the weaker the abstraction."**

And a design rule you'll hear: **"Accept interfaces, return structs."** Functions should ask for the
smallest interface they need as a parameter, and return concrete types.

Lesson 11 puts `Stringer`, `Reader` and `Writer` to work.

---

## 9. Common mistakes

| You wrote | What happens | Fix |
|---|---|---|
| `type Rect struct {...} implements Shape` | syntax error | Nothing to write. Just add the methods |
| A method missing or with a different signature | `does not implement Shape (missing method X)` or `(wrong type for method X)` | Match the method name, parameters and return type exactly |
| Pointer-receiver methods, then `var s Shape = T{}` | `method has pointer receiver` | `&T{}` |
| `s.W` where `s` is a `Shape` | `s.W undefined` | Only interface methods are callable (lesson 11 shows type assertions) |
| Huge interfaces with 10 methods | hard to satisfy and to fake in tests | Split into small interfaces |

---

## 10. Exercises

The tests are written so the file compiles **before** your methods exist. Until you add a method,
the test tells you which type doesn't implement the interface yet.

| | Exercise | Practises |
|---|---|---|
| Ex1 | Make `Square` a `Shape` | adding methods to satisfy an interface |
| Ex2 | Make `Circle` a `Shape` | same, with `math.Pi` |
| Ex3 | `TotalArea` of any shapes | an interface as a parameter |
| Ex4 | `Biggest` shape | an interface as a return value |
| Ex5 | Make `Dog` a `Speaker`, then write `Chorus` | pointer receivers in interfaces (`Robot`) |
| Ex6 | Finish the `Payer` interface and `Checkout` | designing your own interface |

---

## ✅ Checklist

- [ ] I can define an interface and explain what it means
- [ ] I know a type satisfies an interface just by having the methods
- [ ] I can write a function that accepts an interface, and pass different types to it
- [ ] I know why `var s Shape = Counter{}` fails when the method has a pointer receiver
- [ ] I know why small interfaces are preferred

**Next → [Lesson 11: Interfaces in practice](../11_interfaces_practice/README.md)**
