# Lesson 07: Structs

**Goal:** group related data into your own types, create and use them, put them in slices and maps,
and understand that structs are **copied** when you pass them around.

**Time:** about 1 hour.

```bash
cd ~/Desktop/learning-go/course/part1/07_structs
```

---

## 1. Why structs?

Say you're tracking books. With what you know so far, you'd need separate variables or parallel slices:

```go
titles := []string{"Dune", "Emma"}
pages := []int{412, 474}
```

That gets messy fast. In Python you'd write a class or a dataclass:

```python
@dataclass
class Book:
    title: str
    author: str
    pages: int
```

Go has no classes. It has **structs**: a named group of fields.

```go
type Book struct {
	Title  string
	Author string
	Pages  int
}
```

`type Book struct { ... }` creates a new type called `Book`. It goes **outside** any function,
usually near the top of the file.

---

## 2. Creating a struct value

```go
b := Book{Title: "Dune", Author: "Frank Herbert", Pages: 412}
fmt.Println(b)            // {Dune Frank Herbert 412}
fmt.Printf("%+v\n", b)    // {Title:Dune Author:Frank Herbert Pages:412}   ← %+v shows field names
```

Always use the `Field: value` form. (You can also write `Book{"Dune", "Frank Herbert", 412}` in field
order, but it breaks as soon as someone adds a field, so avoid it.)

### Fields you leave out get their zero value

```go
b := Book{Title: "Draft"}
fmt.Printf("%+v\n", b)    // {Title:Draft Author: Pages:0}

var empty Book            // every field is a zero value
```

---

## 3. Reading and changing fields

Use a dot, like Python attributes:

```go
b := Book{Title: "Dune", Pages: 412}
fmt.Println(b.Title)      // Dune
b.Pages = 500
fmt.Println(b.Pages)      // 500
```

Unlike Python objects, you can't add new fields at runtime. `b.Rating = 5` is a compile error,
because `Book` has no `Rating` field. The struct definition is the complete list.

---

## 4. Structs in slices and maps

```go
library := []Book{
	{Title: "Dune", Author: "Frank Herbert", Pages: 412},
	{Title: "Emma", Author: "Jane Austen", Pages: 474},
}

total := 0
for _, book := range library {
	total += book.Pages
}
fmt.Println(total)        // 886
```

Inside `[]Book{...}` you can leave out the type name on each element.

A map with struct values works the same way:

```go
byTitle := map[string]Book{
	"Dune": {Title: "Dune", Pages: 412},
}
fmt.Println(byTitle["Dune"].Pages)   // 412
```

---

## 5. ⚠️ Structs are copied

This is the most important thing in this lesson, and it's different from Python.

In Python, assigning an object or passing it to a function shares the **same** object:

```python
b = Book("Dune", "Herbert", 412)
c = b
c.pages = 1
print(b.pages)    # 1: b and c are the same object
```

In Go, assigning a struct **copies all its fields**:

```go
b := Book{Title: "Dune", Pages: 412}
c := b            // c is a separate copy
c.Pages = 1
fmt.Println(b.Pages)    // 412: b didn't change
```

The same happens when you pass a struct to a function:

```go
func addPages(book Book, n int) {
	book.Pages += n       // changes the function's own copy
}

b := Book{Pages: 100}
addPages(b, 50)
fmt.Println(b.Pages)      // 100: unchanged!
```

And when you loop over a slice of structs (remember lesson 05: `range` gives a copy):

```go
for _, book := range library {
	book.Pages = 0        // changes the copy only
}
// library is unchanged
```

**So how do I change a struct inside a function?** Two options:

1. **Return the changed copy** (simple, and often the cleanest):
   ```go
   func withMorePages(book Book, n int) Book {
   	book.Pages += n
   	return book
   }
   b = withMorePages(b, 50)
   ```
2. **Use a pointer**, which is the next lesson.

To change structs inside a slice, use the index: `library[i].Pages = 0`.

### 🤔 Predict the output

```go
type Point struct{ X, Y int }

func move(p Point) Point {
	p.X += 10
	return p
}

func main() {
	a := Point{X: 1, Y: 2}
	b := move(a)
	a.Y = 5
	fmt.Println(a, b)
}
```

<details><summary>Answer</summary>

`{1 5} {11 2}`. `move` gets a copy of `a` and returns a changed copy. Changing `a.Y` afterwards
doesn't affect `b`. (Also notice `X, Y int`: fields of the same type can share a line.)
</details>

---

## 6. Comparing structs

You can compare two structs with `==` if all their fields can be compared:

```go
p1 := Point{1, 2}
p2 := Point{1, 2}
fmt.Println(p1 == p2)     // true: same field values
```

This doesn't work if the struct contains a slice or a map (those can't be compared with `==`).

---

## 7. Structs inside structs

```go
type Address struct {
	City    string
	Country string
}

type Person struct {
	Name    string
	Address Address
}

p := Person{
	Name:    "Ana",
	Address: Address{City: "Lisbon", Country: "Portugal"},
}
fmt.Println(p.Address.City)   // Lisbon
```

---

## 8. "Constructor" functions

Go has no `__init__`. When you want to set defaults or check values, write a normal function. By
convention its name starts with `New`:

```go
func NewBook(title, author string) Book {
	return Book{Title: title, Author: author, Pages: 100}
}
```

Many Go types don't need one at all: their zero value is already useful.

---

## 9. Common mistakes

| You wrote | What happens | Fix |
|---|---|---|
| `type Book struct` inside `main` | it works, but the type only exists inside `main` | define types at the top level |
| Changed a struct inside a function and expected the caller to see it | the caller has a different copy | return the struct, or use a pointer (lesson 08) |
| `for _, b := range books { b.Pages = 0 }` | the slice doesn't change | `books[i].Pages = 0` |
| `b.Rating = 5` when there's no such field | `b.Rating undefined` | add the field to the struct |

---

## 10. Exercises

`exercises.go` already defines `Book` and `Rectangle` for you.

| | Exercise | Practises |
|---|---|---|
| Ex1 | `NewBook`: create a Book | struct literals |
| Ex2 | `Summary`: `"Dune by Frank Herbert (412 pages)"` | reading fields |
| Ex3 | `TotalPages` of a library | looping over structs |
| Ex4 | `LongestBook` | finding the max by a field |
| Ex5 | `BooksBy` an author | filtering structs |
| Ex6 | `Area` of a `Rectangle` | functions that take structs |
| Ex7 | `Scaled`: return a bigger copy | copies, returning structs |
| Ex8 | `ResetPages`: set every book's pages to 0 in place | changing structs in a slice |
| Ex9 | `CountByAuthor` | structs + maps together |

---

## ✅ Checklist

- [ ] I can define a struct type and create values with `Type{Field: value}`
- [ ] I know missing fields get their zero value
- [ ] I understand that assigning or passing a struct **copies** it
- [ ] I can change a struct by returning a new copy, or through a slice index
- [ ] I can put structs in slices and maps

**Next → [Lesson 08: Pointers and methods](../08_pointers_methods/README.md)**
