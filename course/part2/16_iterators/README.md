# Lesson 16: Iterators

**Goal:** write your own functions that work with `for ... range`, which is Go's version of Python
generators, and use the iterators the standard library provides.

**Time:** about 1 hour. This lesson uses generics from lesson 15.

```bash
cd ~/Desktop/learning-go/course/part2/16_iterators
```

---

## 1. Python generators, and Go's version

In Python, a generator produces values one at a time with `yield`:

```python
def count_up(n):
    for i in range(n):
        yield i

for x in count_up(3):
    print(x)          # 0 1 2
```

Values are made **on demand**, so a generator can even be infinite.

Go (since version 1.23) lets you `range` over a **function**. The function receives a `yield`
function, and calls it once for each value:

```go
func CountUp(n int) iter.Seq[int] {
	return func(yield func(int) bool) {
		for i := range n {
			if !yield(i) {
				return
			}
		}
	}
}

func main() {
	for x := range CountUp(3) {
		fmt.Println(x)      // 0 1 2
	}
}
```

It looks more complicated than Python's, so let's take it apart.

---

## 2. Reading an iterator, piece by piece

```go
func CountUp(n int) iter.Seq[int] {
```

`iter.Seq[int]` (from package `"iter"`) is the type "a sequence of ints". It's a generic type
(lesson 15), defined in the standard library as:

```go
type Seq[V any] func(yield func(V) bool)
```

So an `iter.Seq[int]` is **a function that takes a `yield` function**. `CountUp` returns one:

```go
	return func(yield func(int) bool) {
```

Inside it, you call `yield(value)` for each value you want to produce, like Python's `yield value`:

```go
		for i := range n {
			if !yield(i) {
				return
			}
		}
```

**Why `if !yield(i) { return }`?** `yield` returns a `bool`:
- `true` means "the loop wants more, keep going".
- `false` means "the loop body did `break` (or `return`), **stop now**".

You **must** stop when `yield` returns false. If you keep calling it, the program panics with
`range function continued iteration after function for loop body returned false`.

What happens when you write `for x := range CountUp(3) { ... }`: Go calls the function, and passes it
a `yield` whose job is "run the loop body with this value". Each `yield(i)` runs the loop body once.

---

## 3. Infinite sequences

Like Python generators, a sequence can go on forever, because values are only made when the loop
asks for them:

```go
func Naturals() iter.Seq[int] {
	return func(yield func(int) bool) {
		for i := 1; ; i++ {          // no end condition: loops forever
			if !yield(i) {
				return               // ...until the caller breaks
			}
		}
	}
}

for n := range Naturals() {
	if n > 5 {
		break                         // yield returns false, the iterator returns
	}
	fmt.Print(n, " ")                 // 1 2 3 4 5
}
```

---

## 4. Pairs: `iter.Seq2`

`range` over a slice gives two values (index and element). Your iterators can do that too with
**`iter.Seq2[K, V]`**, where `yield` takes two arguments:

```go
func Enumerate(words []string) iter.Seq2[int, string] {
	return func(yield func(int, string) bool) {
		for i, w := range words {
			if !yield(i+1, w) {      // numbering from 1
				return
			}
		}
	}
}

for n, w := range Enumerate([]string{"go", "py"}) {
	fmt.Println(n, w)                 // 1 go
}                                     // 2 py
```

### 🤔 Predict the output

```go
func Squares() iter.Seq[int] {
	return func(yield func(int) bool) {
		for i := 1; ; i++ {
			if !yield(i * i) {
				return
			}
		}
	}
}

func main() {
	sum := 0
	for sq := range Squares() {
		if sq > 20 {
			break
		}
		sum += sq
	}
	fmt.Println(sum)
}
```

<details markdown="1">
<summary>Answer</summary>

`30`. The squares are 1, 4, 9, 16, 25... The loop adds 1 + 4 + 9 + 16 = 30, then sees 25 > 20 and breaks.
</details>

---

## 5. Iterators in the standard library

The `slices` and `maps` packages produce and consume iterators:

| Function | What it does |
|---|---|
| `slices.Values(s)` | an `iter.Seq` of the slice's elements |
| `slices.All(s)` | an `iter.Seq2` of (index, element) |
| `maps.Keys(m)` | an `iter.Seq` of the map's keys |
| `maps.Values(m)` | an `iter.Seq` of the map's values |
| `maps.All(m)` | an `iter.Seq2` of (key, value) |
| `slices.Collect(seq)` | collect any `iter.Seq` into a slice. Python: `list(gen)` |
| `slices.Sorted(seq)` | collect and sort |

That's why `slices.Sorted(maps.Keys(m))` gives you a map's keys in order, in one line. Compare that
with the loop you wrote in lesson 06.

---

## 6. Iterators that transform other iterators

Because a sequence is just a value, you can write functions that take a sequence and return a new
one, building a chain (a "pipeline"):

```go
func Filter[T any](seq iter.Seq[T], keep func(T) bool) iter.Seq[T] {
	return func(yield func(T) bool) {
		for v := range seq {
			if keep(v) {
				if !yield(v) {
					return
				}
			}
		}
	}
}

evens := Filter(Naturals(), func(n int) bool { return n%2 == 0 })
for n := range evens {
	if n > 10 {
		break
	}
	fmt.Print(n, " ")                  // 2 4 6 8 10
}
```

Nothing is computed until the final loop asks for values, and it stops as soon as you `break`, even
though `Naturals()` is infinite. Python: `(n for n in naturals() if n % 2 == 0)`.

---

## 7. When to write an iterator

Most of the time a plain slice is simpler: build it, return it. Write an iterator when:
- the sequence is **infinite** or very long, so you can't build it all up front;
- values are **expensive** to produce, and the caller might stop early;
- you're exposing the contents of your own **collection type** (like a `Set`) without giving away its
  internal slice or map.

---

## 8. Common mistakes

| You wrote | What happens | Fix |
|---|---|---|
| Ignored `yield`'s result | panic when the loop breaks | `if !yield(v) { return }` |
| `yield(v)` after the loop body broke | same panic | return as soon as `yield` gives false |
| Forgot `return func(yield ...) {...}` and wrote the loop directly in `CountUp` | compile errors | the function must **return** a function |
| `iter.Seq[int]` without importing `"iter"` | `undefined: iter` | `import "iter"` |
| An iterator for a small fixed list | more complex than needed | return a slice |

---

## 9. Exercises

| | Exercise | Practises |
|---|---|---|
| Ex1 | `Countdown`: n, n−1, ..., 1 | your first iterator |
| Ex2 | `EvensBelow` | yield and stopping |
| Ex3 | `Fibonacci`, an infinite sequence | infinite iterators |
| Ex4 | `Take`: the first n values of any sequence | consuming an iterator with `range` and `break` |
| Ex5 | `FilterSeq` | a generic iterator that wraps another one |
| Ex6 | `NumberedLines` | `iter.Seq2` |
| Ex7 | `WordsByLength` | `maps.Keys`, `slices.Sorted`, `slices.Collect` |

The tests also check that your iterators stop correctly when the loop breaks early.

---

## ✅ Checklist: end of Part 2 🎉

- [ ] I can write a function that returns `iter.Seq[T]` and use it in a `for ... range`
- [ ] I know why I must check `yield`'s return value
- [ ] I can write an infinite sequence, and stop it with `break`
- [ ] I can use `slices.Collect`, `slices.Sorted` and `maps.Keys`

---

## What's next

You've finished Part 2. Before Part 3:

1. **Practise.** Do these practice sets in order. They're harder and mix everything:
   - `practice/04_interfaces`
   - `practice/05_errors`
   - `practice/06_generics`
2. **Build something.** The Part 2 milestone in [roadmap.html](../../../roadmap.html): a plugin-style
   file processor, built around a small interface and proper error wrapping.

Then Part 3: concurrency, Go's superpower.

**Next → [Part 3, Lesson 17: Goroutines and WaitGroup](../../part3/17_goroutines/README.md)**
