# Lesson 03: Control flow (`if`, `for`, `switch`)

**Goal:** make decisions with `if` and `switch`, and repeat things with Go's only loop, `for`.

**Time:** 1–1.5 hours.

```bash
cd ~/Desktop/learning-go/course/part1/03_control_flow
```

---

## 1. `if`

Python:

```python
if temperature > 30:
    print("hot")
elif temperature > 15:
    print("nice")
else:
    print("cold")
```

Go:

```go
if temperature > 30 {
	fmt.Println("hot")
} else if temperature > 15 {
	fmt.Println("nice")
} else {
	fmt.Println("cold")
}
```

Differences to notice:
- **No parentheses** around the condition (you *can* add them, but `gofmt` removes them).
- **Braces are required**, even for a one-line body.
- `elif` is written `else if`.
- `else` must be on the same line as the closing `}`: `} else {`.

### Logic operators

| Python | Go |
|---|---|
| `and` | `&&` |
| `or` | `\|\|` |
| `not x` | `!x` |
| `==`, `!=`, `<`, `>=` ... | the same |

```go
age := 25
hasTicket := true
if age >= 18 && hasTicket {
	fmt.Println("come in")
}
if !hasTicket || age < 18 {
	fmt.Println("sorry")
}
```

### Only real booleans

In Python, `if 0:`, `if "":` and `if []:` are all false ("falsy"). **Go has no truthiness.**
The condition must be a `bool`:

```go
n := 5
if n {          // compile error: non-boolean condition in if statement
}
if n != 0 {     // correct
}
```

### `if` with a short statement

You can run a short statement before the condition, separated by `;`. A variable created there
only exists inside the `if`/`else`:

```go
if length := len("gopher"); length > 5 {
	fmt.Println("long word:", length)
}
// length doesn't exist here
```

This looks odd at first, but you'll see it constantly with errors in Part 2:

```go
if err := doSomething(); err != nil {
	// handle the error
}
```

---

## 2. `for`: Go's only loop

Go has **no `while` loop**. `for` does everything, in four shapes.

### Shape 1: count with `range`

```go
for i := range 5 {
	fmt.Println(i)      // 0 1 2 3 4
}
```

Python: `for i in range(5):`

### Shape 2: the classic three-part loop

```go
for i := 1; i <= 10; i += 2 {
	fmt.Println(i)      // 1 3 5 7 9
}
```

The three parts, separated by `;`: **start** (`i := 1`), **keep going while** (`i <= 10`), **after each round** (`i += 2`).

Python: `for i in range(1, 11, 2):`

### Shape 3: "while"

Give `for` just a condition:

```go
n := 1
for n < 100 {
	n *= 2
}
fmt.Println(n)          // 128
```

Python: `while n < 100:`

### Shape 4: forever

```go
for {
	// runs until something calls break or return
}
```

Python: `while True:`

### `break` and `continue`

They work exactly as in Python:

```go
for i := range 10 {
	if i%2 == 0 {
		continue    // skip even numbers
	}
	if i > 7 {
		break       // stop the loop completely
	}
	fmt.Println(i)  // 1 3 5 7
}
```

Looping over lists, strings and dicts (`for i, x := range ...`) comes in lessons 04–06.

### 🤔 Predict the output

```go
total := 0
for i := 1; i <= 4; i++ {
	total += i
}
count := 0
for total > 1 {
	total /= 2
	count++
}
fmt.Println(total, count)
```

<details><summary>Answer</summary>

`1 3`. First loop: total = 1+2+3+4 = 10. Second loop halves it with integer division: 10 → 5 → 2 → 1,
which takes 3 rounds.
</details>

---

## 3. `switch`

A cleaner way to write long `if / else if` chains. Python 3.10 has `match`; Go's `switch` is simpler.

### Switch on a value

```go
day := "sat"
switch day {
case "sat", "sun":            // several values in one case
	fmt.Println("weekend")
case "fri":
	fmt.Println("almost weekend")
default:                      // like a final else
	fmt.Println("weekday")
}
```

**No `break` needed.** Unlike C or Java, Go runs only the matching case and then stops.

### Switch with no value (a tidy if/else chain)

Leave out the value, and write a condition on each case. The first true one wins:

```go
score := 72
switch {
case score >= 90:
	fmt.Println("A")
case score >= 80:
	fmt.Println("B")
case score >= 70:
	fmt.Println("C")
default:
	fmt.Println("F")
}
```

This prints `C`. Go programmers prefer this to long `if / else if` chains.

### 🤔 Predict the output

```go
for i := range 4 {
	switch {
	case i == 0:
		fmt.Print("zero ")
	case i%2 == 1:
		fmt.Print("odd ")
	default:
		fmt.Print("even ")
	}
}
```

<details><summary>Answer</summary>

`zero odd even odd`. (`fmt.Print` doesn't add a newline.) For i = 0 the first case matches and the
others are skipped, even though 0 is also even.
</details>

---

## 4. Common mistakes

| You wrote | What happens | Fix |
|---|---|---|
| `if x > 5:` | syntax error | `if x > 5 {` |
| `}` then `else {` on the next line | syntax error | `} else {` on one line |
| `while x < 5 {` | `while` is undefined | `for x < 5 {` |
| `if n {` with an int | `non-boolean condition` | `if n != 0 {` |
| `x and y` | syntax error | `x && y` |
| `for i := range 5` and `i` unused | `declared and not used` | `for range 5 {` (no variable) |

---

## 5. Exercises

| | Exercise | Practises |
|---|---|---|
| Ex1 | `Sign`: "positive", "negative" or "zero" | `if / else if / else` |
| Ex2 | `CanVote`: age and citizenship | `&&` |
| Ex3 | `SumTo`: 1 + 2 + ... + n | `for` with a counter |
| Ex4 | `Factorial` | loop multiplication |
| Ex5 | `FizzBuzzWord`: the FizzBuzz word for one number | `switch` with no value |
| Ex6 | `DayType`: weekday or weekend | `switch` on a value, several values per case |
| Ex7 | `CountDigits` | "while"-style loop |
| Ex8 | `IsPrime` | loops with an early `return` |
| Ex9 | `CollatzSteps` (a little harder) | everything together |

---

## ✅ Checklist

- [ ] I can write `if / else if / else` without parentheses, with braces
- [ ] I know Go has no truthiness: conditions must be real `bool`s
- [ ] I can write all four shapes of `for`: `range n`, three-part, "while" and forever
- [ ] I can use both kinds of `switch`, and I know cases don't fall through

**Next → [Lesson 04: Strings](../04_strings/README.md)**
