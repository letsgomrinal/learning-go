# Lesson 06: Maps (Go's dicts)

**Goal:** store and look up key/value pairs, check whether a key exists, count things, and make sets.

**Time:** about 1 hour.

```bash
cd ~/Desktop/learning-go/course/part1/06_maps
```

---

## 1. Creating a map

Python:

```python
ages = {"ana": 30, "bo": 25}
```

Go:

```go
ages := map[string]int{
	"ana": 30,
	"bo":  25,
}
fmt.Println(ages)          // map[ana:30 bo:25]
```

`map[string]int` reads as "a map from string keys to int values". Like slices, one type for all keys
and one type for all values.

An empty map:

```go
scores := map[string]int{}       // or: make(map[string]int)
```

---

## 2. Setting, getting, deleting

```go
ages := map[string]int{}
ages["ana"] = 30             // add or update, like Python
ages["bo"] = 25
fmt.Println(ages["ana"])     // 30
fmt.Println(len(ages))       // 2
delete(ages, "bo")           // Python: del ages["bo"]
```

### Missing keys don't crash

In Python, `ages["zed"]` raises `KeyError`. In Go, a missing key gives you the **zero value**:

```go
fmt.Println(ages["zed"])     // 0   (no error!)
```

That's handy for counting (see section 4), but it means you can't tell "missing" from "stored a 0".

### The "comma ok" check

To know whether a key is really there, ask for a second value:

```go
age, ok := ages["zed"]
if ok {
	fmt.Println("found:", age)
} else {
	fmt.Println("no such person")
}
```

`ok` is `true` if the key exists. Python equivalent: `if "zed" in ages:`. You can combine it with the
short `if` from lesson 03:

```go
if age, ok := ages["ana"]; ok {
	fmt.Println("ana is", age)
}
```

This is called the "comma ok" idiom. You'll see it a lot in Go.

---

## 3. Looping over a map

```go
for name, age := range ages {
	fmt.Println(name, age)
}
```

Python: `for name, age in ages.items():`

**⚠️ The order is random.** Python dicts keep insertion order; Go maps don't, and Go deliberately
shuffles the order each time to stop you relying on it. Run this a few times:

```go
m := map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}
for k := range m {
	fmt.Print(k, " ")
}
```

If you need a fixed order, collect the keys and sort them:

```go
keys := make([]string, 0, len(m))
for k := range m {
	keys = append(keys, k)
}
slices.Sort(keys)
for _, k := range keys {
	fmt.Println(k, m[k])
}
```

(A shortcut, once you know iterators in Part 2: `slices.Sorted(maps.Keys(m))`.)

`fmt.Println(m)` always prints maps with sorted keys, to make output readable. Only loops are random.

---

## 4. The counting pattern

Because missing keys give 0, counting needs no special setup. Python needs `defaultdict` or `Counter`:

```go
words := []string{"go", "is", "go", "fun", "go"}
counts := map[string]int{}
for _, w := range words {
	counts[w]++                   // missing key starts at 0, then becomes 1
}
fmt.Println(counts)               // map[fun:1 go:3 is:1]
```

### Grouping into slices

The same trick works with slices as values, because appending to a missing (nil) slice just works:

```go
byLetter := map[string][]string{}
for _, name := range []string{"ana", "amy", "bo"} {
	first := name[:1]
	byLetter[first] = append(byLetter[first], name)
}
fmt.Println(byLetter)             // map[a:[ana amy] b:[bo]]
```

### 🤔 Predict the output

```go
m := map[string]int{"x": 1}
m["y"] += 5
v, ok := m["z"]
delete(m, "x")
fmt.Println(len(m), m["y"], v, ok)
```

<details><summary>Answer</summary>

`1 5 0 false`. `m["y"] += 5` starts from the zero value 0. Looking up `"z"` gives 0 and `false`
(and doesn't add it). After deleting `"x"`, only `"y"` is left.
</details>

---

## 5. Sets

Go has no `set` type. Use a map whose values you ignore. The most common version uses `bool`:

```go
seen := map[string]bool{}
seen["go"] = true
seen["py"] = true

if seen["go"] {                   // missing keys give false, so this reads nicely
	fmt.Println("seen go")
}
```

You'll also see `map[string]struct{}` in real code. `struct{}` is an empty type that takes no
memory, so it's a slightly more efficient set. `map[string]bool` is easier to read while learning.

---

## 6. The nil map trap

A map declared with `var` and no value is `nil`. You can **read** from it, but **writing panics**:

```go
var m map[string]int
fmt.Println(m["a"])     // 0, fine
m["a"] = 1              // panic: assignment to entry in nil map
```

Always create maps with `map[string]int{}` or `make(...)` before writing to them.

---

## 7. Maps are shared, not copied

Assigning a map to another variable, or passing it to a function, does **not** copy it. Both names
refer to the same map (same as Python dicts):

```go
a := map[string]int{"x": 1}
b := a
b["x"] = 99
fmt.Println(a["x"])     // 99
```

Use `maps.Clone(a)` (package `"maps"`) for a real copy.

---

## 8. Common mistakes

| You wrote | What happens | Fix |
|---|---|---|
| `var m map[string]int` then `m["a"] = 1` | panic: assignment to entry in nil map | `m := map[string]int{}` |
| Expecting the loop order to match insertion order | random order | sort the keys |
| `if m["k"] != 0` to check existence | wrong when the stored value is 0 | `if _, ok := m["k"]; ok` |
| `"k" in m` | syntax error | comma ok |

---

## 9. Exercises

| | Exercise | Practises |
|---|---|---|
| Ex1 | `Lookup`: a price, or "not found" | comma ok |
| Ex2 | `WordCount` | the counting pattern |
| Ex3 | `HasDuplicates` | a set |
| Ex4 | `Unique`: remove duplicates, keep order | set + slice together |
| Ex5 | `Invert`: swap keys and values | building a map in a loop |
| Ex6 | `SortedKeys` | collecting keys, `slices.Sort` |
| Ex7 | `GroupByFirstLetter` | map of slices |
| Ex8 | `MostCommon` word | counting, then finding the max |

---

## ✅ Checklist

- [ ] I can create a map, set/get/delete values, and get its length
- [ ] I know a missing key gives the zero value, and how to use `value, ok := m[key]`
- [ ] I know map iteration order is random, and how to get sorted output
- [ ] I can count with `m[k]++` and build a set with `map[T]bool`
- [ ] I know writing to a nil map panics

**Next → [Lesson 07: Structs](../07_structs/README.md)**
