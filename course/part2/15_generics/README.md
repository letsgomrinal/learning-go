# Lesson 15: Generics

**Goal:** write one function or type that works with many element types, while keeping full type
checking, and know when generics are the right tool and when an interface is.

**Time:** 1.5 hours.

```bash
cd ~/Desktop/learning-go/course/part2/15_generics
```

---

## 1. The problem: the same code, many times

Without generics, a "find the biggest number" function only works for one type:

```go
func MaxInt(values []int) int { ... }
func MaxFloat(values []float64) float64 { ... }   // identical code, different type
func MaxString(values []string) string { ... }    // and again...
```

You could use `any` (lesson 11), but then you lose type checking, need type assertions, and get
`any` back instead of the real type. Not great.

Python doesn't have this problem at runtime (everything is dynamic), but if you use type hints you've
seen the same idea there: `def max_of(values: list[T]) -> T` with a `TypeVar`. Go's generics are
that, except the compiler enforces them.

---

## 2. Your first generic function

```go
func Max[T cmp.Ordered](values []T) T {
	biggest := values[0]
	for _, v := range values {
		if v > biggest {
			biggest = v
		}
	}
	return biggest
}
```

The new part is `[T cmp.Ordered]`, right after the function name. It declares a **type parameter**:

```
func Max [T cmp.Ordered] (values []T) T
         │  │
         │  └── the constraint: which types T is allowed to be
         └── the type parameter's name (T is the convention, like in Python)
```

Read it as: "`Max` works for any type `T` that can be ordered (compared with `<` and `>`). It takes
a slice of `T` and returns a `T`."

`cmp.Ordered` comes from the standard `cmp` package (import `"cmp"`). It means "all the types that
support `<`": ints, floats and strings.

Using it:

```go
fmt.Println(Max([]int{3, 9, 2}))             // 9
fmt.Println(Max([]float64{1.5, 0.5}))        // 1.5
fmt.Println(Max([]string{"go", "py", "c"}))  // py
```

You don't write the type when calling. Go **infers** `T` from the argument. (You *can* write it:
`Max[int](nums)`. Sometimes it's needed when there's nothing to infer from.)

And the type checking is real:

```go
Max([]bool{true, false})   // compile error: bool does not satisfy cmp.Ordered
```

---

## 3. Constraints: what can T be?

A constraint says what `T` must support. The common ones:

| Constraint | Allows | Lets you use |
|---|---|---|
| `any` | every type | nothing special: store, pass around, return |
| `comparable` | types that work with `==` | `==`, `!=`, map keys |
| `cmp.Ordered` | numbers and strings | `<`, `>`, `<=`, `>=` (and `==`) |

The rule: **you can only do with `T` what the constraint allows.**

```go
func Contains[T comparable](items []T, target T) bool {
	for _, item := range items {
		if item == target {        // allowed: T is comparable
			return true
		}
	}
	return false
}
```

With `[T any]`, `item == target` would be a compile error, because not every type supports `==`
(slices and maps don't).

### Making your own constraint

A constraint is an interface. It can list types with `|` ("or"):

```go
type Number interface {
	~int | ~int64 | ~float64
}

func Sum[T Number](values []T) T {
	var total T              // the zero value of T: 0 for any number type
	for _, v := range values {
		total += v           // allowed: every type in Number supports +
	}
	return total
}
```

The `~` means "this type, **or any type built on it**". Remember `type Celsius float64` from lesson 08?
`~float64` includes `Celsius`; plain `float64` would not.

```go
type Celsius float64
temps := []Celsius{20.5, 22}
fmt.Println(Sum(temps))      // 42.5, and the result is a Celsius
```

---

## 4. Several type parameters

```go
func Map[T, U any](items []T, f func(T) U) []U {
	result := make([]U, 0, len(items))
	for _, item := range items {
		result = append(result, f(item))
	}
	return result
}

words := []string{"go", "rust", "c"}
lengths := Map(words, func(w string) int { return len(w) })
fmt.Println(lengths)         // [2 4 1]
```

`T` is the input element type, `U` the output element type. Here `T` = `string`, `U` = `int`.
Python equivalent: `list(map(len, words))`.

### 🤔 Predict the output

```go
func Count[T comparable](items []T, target T) int {
	n := 0
	for _, it := range items {
		if it == target {
			n++
		}
	}
	return n
}

func main() {
	fmt.Println(Count([]string{"a", "b", "a"}, "a"), Count([]int{1, 2, 3}, 4))
}
```

<details><summary>Answer</summary>

`2 0`. The same function works for strings and ints. `T` is inferred separately at each call.
</details>

---

## 5. Generic types

Types can have type parameters too. A stack that works for any element type:

```go
type Stack[T any] struct {
	items []T
}

func (s *Stack[T]) Push(v T) {
	s.items = append(s.items, v)
}

func (s *Stack[T]) Pop() (T, bool) {
	var zero T
	if len(s.items) == 0 {
		return zero, false
	}
	last := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return last, true
}
```

Using it, you say which type it holds:

```go
var ints Stack[int]
ints.Push(1)
ints.Push(2)
v, _ := ints.Pop()           // v is an int: 2

names := &Stack[string]{}
names.Push("ana")
// names.Push(42)            // compile error: cannot use 42 as string value
```

Notes:
- In methods, the receiver repeats the type parameter: `(s *Stack[T])`.
- `var zero T` is how you get "the zero value of whatever T is", for the "nothing to return" case.

---

## 6. You've already used generics

The `slices` and `maps` packages from lessons 05 and 06 are generic. That's how `slices.Contains`,
`slices.Max` and `slices.Sort` work for any element type. Look:

```bash
go doc slices.Contains
# func Contains[S ~[]E, E comparable](s S, v E) bool
```

Now you can read that signature: it works on any slice type `S` whose elements `E` are comparable.

---

## 7. Generics or interfaces?

Both let one piece of code work with different types. They solve different problems:

| Use **generics** when... | Use **interfaces** when... |
|---|---|
| the code is the **same** for every type, only the element type changes | each type **behaves differently** behind the same methods |
| containers: stacks, sets, caches | `Shape.Area()`, `io.Writer`, `error` |
| helpers on slices and maps: `Map`, `Filter`, `Contains`, `Max` | plug-in behaviour: different payment providers, storage backends |

A good rule: **write normal code first.** Reach for generics when you notice you're writing the same
function for the third type. Go worked without generics for its first ten years, and most Go code
still doesn't need them.

---

## 8. Common mistakes

| You wrote | What happens | Fix |
|---|---|---|
| `a == b` with `[T any]` | `incomparable types in type set` | `[T comparable]` |
| `a < b` with `[T comparable]` | `type parameter T cannot use operator <` | `[T cmp.Ordered]` |
| `return nil` from `func F[T any]() T` | `cannot use nil as T value` | `var zero T; return zero` |
| `float64` in a constraint, used with `Celsius` | `Celsius does not satisfy Number (possibly missing ~ for float64 in Number)` | `~float64` |

---

## 9. Exercises

| | Exercise | Practises |
|---|---|---|
| Ex1 | `Contains` for any comparable type | `comparable` |
| Ex2 | `Max` that handles an empty slice | `cmp.Ordered`, `var zero T` |
| Ex3 | `Map` | two type parameters |
| Ex4 | `Filter` | functions as parameters + generics |
| Ex5 | `Sum` that works with `Celsius` | your own constraint, `~` |
| Ex6 | A generic `Stack[T]` | generic types and methods |
| Ex7 | `SortedKeys` of any map | constraints on map keys |

---

## ✅ Checklist

- [ ] I can read and write a function signature with a type parameter, like `func F[T any](x T) T`
- [ ] I know what `any`, `comparable` and `cmp.Ordered` allow
- [ ] I can write a constraint with `|` and `~`
- [ ] I can write a generic type with methods
- [ ] I can explain when to use generics and when to use an interface

**Next → [Lesson 16: Iterators](../16_iterators/README.md)**
