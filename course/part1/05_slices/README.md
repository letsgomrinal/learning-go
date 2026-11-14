# Lesson 05: Slices (Go's lists)

**Goal:** create and grow slices, loop over them, cut pieces out of them, and write functions that
take any number of arguments.

**Time:** 1.5 hours. This is one of the most important lessons.

```bash
cd ~/Desktop/learning-go/course/part1/05_slices
```

---

## 1. Creating a slice

Python:

```python
nums = [10, 20, 30]
```

Go:

```go
nums := []int{10, 20, 30}
fmt.Println(nums)          // [10 20 30]
```

`[]int` reads as "a slice of ints". Unlike a Python list, a slice holds **one type only**:

```go
mixed := []int{1, "two", 3}   // compile error
```

Other examples:

```go
names := []string{"Ana", "Bo", "Cy"}
prices := []float64{9.99, 5.50}
var empty []int               // an empty slice (its zero value is nil, and that's fine)
```

---

## 2. Reading and changing elements

Indexes start at 0, like Python:

```go
names := []string{"Ana", "Bo", "Cy"}
fmt.Println(names[0])         // Ana
fmt.Println(len(names))       // 3
names[1] = "Bea"              // change an element
fmt.Println(names)            // [Ana Bea Cy]
```

Two differences from Python:
- **No negative indexes.** `names[-1]` is a compile error. The last element is `names[len(names)-1]`.
- **Going out of range crashes the program** (a "panic"):

```go
fmt.Println(names[5])
// panic: runtime error: index out of range [5] with length 3
```

That's Go's version of Python's `IndexError`. Read the message: it tells you the index and the length.

---

## 3. Adding elements with `append`

Python: `nums.append(40)`. Go:

```go
nums := []int{10, 20, 30}
nums = append(nums, 40)
nums = append(nums, 50, 60)       // several at once
fmt.Println(nums)                 // [10 20 30 40 50 60]
```

**You must write `nums = append(nums, ...)`.** `append` doesn't change `nums` in place; it *returns*
the new slice, and you store it back. Forgetting the `nums =` is a compile error
(`append(nums, 40) (value of type []int) is not used`), which helps.

Appending a whole slice to another uses `...` (like Python's `*` unpacking):

```go
a := []int{1, 2}
b := []int{3, 4}
a = append(a, b...)               // [1 2 3 4]   Python: a.extend(b)
```

### Building a slice in a loop

A very common pattern: start empty, append as you go.

```go
var squares []int
for i := range 5 {
	squares = append(squares, i*i)
}
fmt.Println(squares)              // [0 1 4 9 16]
```

Python: `squares = [i*i for i in range(5)]`. Go has no list comprehensions; you write the loop.

---

## 4. Looping over a slice

```go
fruits := []string{"apple", "banana", "cherry"}

for i, fruit := range fruits {
	fmt.Println(i, fruit)
}
// 0 apple
// 1 banana
// 2 cherry
```

Python: `for i, fruit in enumerate(fruits):`

`range` always gives you **two** things: the index and the element. If you only want the element,
put `_` for the index. If you only want the index, leave out the second variable:

```go
for _, fruit := range fruits { ... }   // Python: for fruit in fruits
for i := range fruits { ... }          // just the indexes
```

Important: `fruit` is a **copy** of the element. Changing it doesn't change the slice:

```go
nums := []int{1, 2, 3}
for _, n := range nums {
	n = n * 10                  // changes the copy only
}
fmt.Println(nums)               // [1 2 3]   (unchanged!)

for i := range nums {
	nums[i] = nums[i] * 10      // change through the index instead
}
fmt.Println(nums)               // [10 20 30]
```

### 🤔 Predict the output

```go
words := []string{"go", "is", "fun"}
total := 0
for i, w := range words {
	if i == 1 {
		continue
	}
	total += len(w)
}
fmt.Println(total)
```

<details markdown="1">
<summary>Answer</summary>

`5`. We skip index 1 ("is"), so it adds len("go") + len("fun") = 2 + 3.
</details>

---

## 5. Slicing: taking a piece

Same syntax as Python, `s[start:end]`, with the end **not included**:

```go
nums := []int{0, 10, 20, 30, 40, 50}
fmt.Println(nums[1:4])      // [10 20 30]
fmt.Println(nums[:2])       // [0 10]
fmt.Println(nums[3:])       // [30 40 50]
```

No negative numbers and no step (`[::2]` doesn't exist).

### ⚠️ The big difference: a slice of a slice shares memory

In Python, `lst[1:4]` makes a **copy**. In Go, it does **not**. It gives you a *window* onto the
same underlying data:

```go
nums := []int{0, 10, 20, 30, 40}
part := nums[1:3]           // [10 20]
part[0] = 999
fmt.Println(nums)           // [0 999 20 30 40]   ← the original changed!
```

This is a deliberate design choice (it makes slicing very fast), but it surprises every Python developer.
When you need an independent copy, say so:

```go
part := slices.Clone(nums[1:3])    // import "slices". A real copy.
```

Lesson 09 (and practice set 02) goes deeper into how slices work under the hood, and why `append`
can sometimes cause the same kind of surprise. For now, remember: **slicing shares, `slices.Clone` copies.**

---

## 6. `make`: a slice of a given size

```go
zeros := make([]int, 5)
fmt.Println(zeros)          // [0 0 0 0 0]
zeros[2] = 7
```

`make([]int, 5)` creates a slice with 5 elements, all set to the zero value. Python: `[0] * 5`.

When you know roughly how many elements you'll append, you can reserve space up front. This is
an optional speed-up:

```go
results := make([]int, 0, 100)   // length 0, but room for 100 before it needs to grow
```

---

## 7. The `slices` package

The standard `slices` package (add `"slices"` to imports) has many helpers:

| Python | Go |
|---|---|
| `x in lst` | `slices.Contains(lst, x)` |
| `lst.index(x)` | `slices.Index(lst, x)` (−1 if missing) |
| `lst.sort()` | `slices.Sort(lst)` (sorts in place) |
| `max(lst)` / `min(lst)` | `slices.Max(lst)` / `slices.Min(lst)` (panic if empty!) |
| `lst.reverse()` | `slices.Reverse(lst)` |
| `lst == other` | `slices.Equal(lst, other)`. Note: `==` doesn't work on slices! |
| `lst.copy()` | `slices.Clone(lst)` |

`go doc slices` lists them all. In the exercises you'll write some of these yourself first, because
writing them teaches you loops and slices. In real code, use the package.

---

## 8. Functions with any number of arguments (variadic)

Python's `*args`:

```python
def total(*nums):
    return sum(nums)

total(1, 2, 3)
```

Go:

```go
func total(nums ...int) int {
	sum := 0
	for _, n := range nums {
		sum += n
	}
	return sum
}

total(1, 2, 3)            // 6
total()                   // 0
values := []int{4, 5}
total(values...)          // 9   ← pass a slice with ... (Python: total(*values))
```

Inside the function, `nums` is a normal `[]int`. `fmt.Println` itself is variadic: that's why it
accepts any number of arguments.

---

## 9. Slices of slices (a grid)

```go
grid := [][]int{
	{1, 2, 3},
	{4, 5, 6},
}
fmt.Println(grid[1][2])       // 6
for _, row := range grid {
	fmt.Println(row)
}
```

`[][]int` is "a slice of slices of ints". (Note the trailing comma after the last row: Go requires
it when the closing `}` is on its own line.)

---

## 10. Arrays (briefly)

You'll occasionally see `[3]int` (with a number inside the brackets). That's an **array**: a fixed
size that can never grow. Arrays are rarely used directly; slices are what you'll use 99% of the time.
If you see a number in the brackets, it's an array. If the brackets are empty, it's a slice.

---

## 11. Common mistakes

| You wrote | What happens | Fix |
|---|---|---|
| `append(nums, 4)` on its own | `value is not used` | `nums = append(nums, 4)` |
| `nums[-1]` | compile error: index must be non-negative | `nums[len(nums)-1]` |
| `for x := range nums` expecting elements | `x` is the **index** | `for _, x := range nums` |
| `if a == b` on two slices | compile error: slice can only be compared to nil | `slices.Equal(a, b)` |
| `n = n * 2` inside `for _, n := range` | the slice doesn't change | `nums[i] = ...` using the index |
| `slices.Max(empty)` | panic | check `len(s) == 0` first |

---

## 12. Exercises

| | Exercise | Practises |
|---|---|---|
| Ex1 | `Sum` of a slice | `range` |
| Ex2 | `Largest`, and whether there was one | handling an empty slice |
| Ex3 | `Evens`: keep the even numbers | building a slice with `append` |
| Ex4 | `Doubled`: a new slice, original untouched | `make` or `append`, copies |
| Ex5 | `DoubleInPlace`: change the original | writing through the index |
| Ex6 | `IndexOf` | loops with early `return` |
| Ex7 | `Average` of any number of floats | variadic functions |
| Ex8 | `Reversed`: a reversed copy | indexes, not touching the input |
| Ex9 | `FirstAndLast` | `len` and safe indexing |

---

## ✅ Checklist

- [ ] I can create a slice, read and change elements, and get its length
- [ ] I always write `s = append(s, x)`
- [ ] I know `range` gives (index, copy of element)
- [ ] I know slicing shares memory with the original, and `slices.Clone` copies
- [ ] I can write a variadic function and pass a slice to it with `...`

**Next → [Lesson 06: Maps](../06_maps/README.md)**
