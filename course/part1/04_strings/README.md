# Lesson 04: Strings

**Goal:** work with text using the `strings` and `strconv` packages, turn text into numbers (with
your first error check), and understand why Go strings are bytes, not characters.

**Time:** 1–1.5 hours.

```bash
cd ~/Desktop/learning-go/course/part1/04_strings
```

---

## 1. String basics

```go
s := "Hello"
t := s + ", Go"           // joining with + works like Python
fmt.Println(t)            // Hello, Go
fmt.Println(len(t))       // 9
fmt.Println(t == "Hello, Go")  // true: compare with ==
```

Strings are **immutable**, like in Python: you can't change a string, only make a new one.

Two kinds of string literal:

```go
a := "line one\nline two"          // double quotes: \n means newline
b := `C:\no\escapes\here
and it can span lines`              // backticks: raw text, exactly as written
```

Backticks are like Python's `r"..."` and `"""..."""` combined.

---

## 2. The `strings` package

Python puts these as methods on the string (`s.upper()`). Go puts them as functions in the
`strings` package, which you pass the string to. Add `"strings"` to your imports.

| Python | Go |
|---|---|
| `s.upper()` | `strings.ToUpper(s)` |
| `s.lower()` | `strings.ToLower(s)` |
| `s.strip()` | `strings.TrimSpace(s)` |
| `"go" in s` | `strings.Contains(s, "go")` |
| `s.startswith("x")` | `strings.HasPrefix(s, "x")` |
| `s.endswith("x")` | `strings.HasSuffix(s, "x")` |
| `s.find("x")` | `strings.Index(s, "x")` (−1 if not found) |
| `s.replace("a", "b")` | `strings.ReplaceAll(s, "a", "b")` |
| `s.split(",")` | `strings.Split(s, ",")` |
| `s.split()` (on any whitespace) | `strings.Fields(s)` |
| `",".join(parts)` | `strings.Join(parts, ",")` |
| `s * 3` | `strings.Repeat(s, 3)` |
| `s.count("a")` | `strings.Count(s, "a")` |

Try a few:

```go
s := "  Hello, Gopher  "
fmt.Println(strings.TrimSpace(s))              // "Hello, Gopher"
fmt.Println(strings.ToUpper(s))                // "  HELLO, GOPHER  "
fmt.Println(strings.Contains(s, "Go"))         // true
fmt.Println(strings.Fields("a  b   c"))        // [a b c]
fmt.Println(strings.Join([]string{"x", "y"}, "-"))   // x-y
```

`strings.Split` and `strings.Fields` return a **slice** of strings, written `[]string`. Slices are
Go's lists, and lesson 05 is all about them. For now, you only need to know you can loop over one:

```go
for _, word := range strings.Fields("go is fun") {
	fmt.Println(word)
}
```

`for _, word := range` gives you each element. (The `_` skips the position number; lesson 05 explains.)

Explore the whole package from your terminal:

```bash
go doc strings          # list everything
go doc strings.Fields   # details of one function
```

Getting used to `go doc` now will save you a lot of time later.

---

## 3. Building strings efficiently

Joining with `+` in a big loop creates a new string every time, which is slow (same as in Python,
where you'd use `"".join(...)`). Use `strings.Builder`:

```go
var sb strings.Builder
for i := range 3 {
	sb.WriteString("go")
	sb.WriteString(strconv.Itoa(i))
}
fmt.Println(sb.String())    // go0go1go2
```

For a few pieces, `+` or `fmt.Sprintf` is perfectly fine.

---

## 4. Text ↔ numbers with `strconv`, and your first error check

Number to text is easy and can't fail:

```go
s := strconv.Itoa(42)            // "42"
```

Text to number **can fail** (what number is `"hello"`?). In Python you'd catch an exception:

```python
try:
    n = int(text)
except ValueError:
    print("not a number")
```

Go has no exceptions. Instead, the function returns **two values**: the result and an error:

```go
n, err := strconv.Atoi("42")
if err != nil {
	fmt.Println("not a number:", err)
} else {
	fmt.Println(n + 1)           // 43
}

n, err = strconv.Atoi("hello")
fmt.Println(n, err)
// 0 strconv.Atoi: parsing "hello": invalid syntax
```

How to read this:
- `err` is `nil` when everything worked. `nil` means "nothing", like Python's `None`.
- `err != nil` means something went wrong, and `err` describes what.
- When there's an error, the number is just the zero value, `0`. Don't use it.

This `value, err := ...` / `if err != nil` pattern is the single most common thing in Go code.
Part 2 covers errors in depth. For now: **always check `err` before using the value.**

Decimals work the same way:

```go
f, err := strconv.ParseFloat("3.14", 64)   // 64 means "give me a float64"
```

---

## 5. Strings are bytes, not characters

This is the one part of Go strings that's really different from Python.

In Python, a string is a sequence of characters. In Go, a string is a sequence of **bytes**,
holding text encoded as UTF-8. English letters are 1 byte each, but many characters take more:

```go
fmt.Println(len("hello"))   // 5
fmt.Println(len("héllo"))   // 6   ← é takes 2 bytes
fmt.Println(len("日本"))     // 6   ← each of these takes 3 bytes
```

So `len` counts **bytes**. And indexing gives you a byte, not a character:

```go
s := "héllo"
fmt.Println(s[0])           // 104      (the byte value of 'h')
fmt.Println(s[1])           // 195      (the first half of 'é'!)
```

### Runes: Go's word for a character

A single character (a Unicode code point) is called a **rune** in Go. Rune literals use **single quotes**:

```go
r := 'é'
fmt.Println(r)              // 233  (a rune is just a number: the character's code)
fmt.Println(string(r))      // é
```

That's why `'hello'` with single quotes is an error: single quotes hold exactly one character.

### Looping over characters

`for ... range` over a string gives you **runes**, which is almost always what you want:

```go
for i, r := range "héllo" {
	fmt.Println(i, string(r))
}
// 0 h
// 1 é
// 3 l      ← position jumped from 1 to 3, because é is 2 bytes
// 4 l
// 5 o
```

`i` is the byte position, `r` is the character.

### Counting and reversing characters

Convert to a slice of runes, `[]rune(s)`, to work character by character:

```go
s := "héllo"
runes := []rune(s)
fmt.Println(len(runes))         // 5  (characters!)
fmt.Println(string(runes[1]))   // é
```

`utf8.RuneCountInString(s)` (package `"unicode/utf8"`) also counts characters.

**The rule of thumb:** for plain ASCII text, bytes and characters are the same and you can use
`len` and `s[i]`. When the text may contain accents, emoji or other languages, use `range` or `[]rune`.

### 🤔 Predict the output

```go
s := "Go🙂"
count := 0
for range s {
	count++
}
fmt.Println(len(s), count, len([]rune(s)))
```

<details markdown="1">
<summary>Answer</summary>

`6 3 3`. `G` and `o` are 1 byte each, and the emoji is 4 bytes, so `len` is 6. But there are only
3 characters, which both `range` and `[]rune` see.
</details>

---

## 6. Useful character checks

The `unicode` package works on runes:

```go
unicode.IsUpper('A')     // true
unicode.IsDigit('7')     // true
unicode.IsLetter('é')    // true
unicode.IsSpace(' ')     // true
unicode.ToUpper('a')     // 'A'
```

---

## 7. Common mistakes

| You wrote | What happens | Fix |
|---|---|---|
| `s.upper()` | `s.upper undefined` | `strings.ToUpper(s)` |
| `'hello'` | error: more than one character in rune literal | `"hello"` |
| `n := strconv.Atoi(s)` | `assignment mismatch: 1 variable but strconv.Atoi returns 2 values` | `n, err := strconv.Atoi(s)` |
| `s[0] = 'H'` | cannot assign to `s[0]` (strings are immutable) | build a new string |
| `len(s)` to count letters in `"héllo"` | 6, not 5 | `len([]rune(s))` |
| `string(65)` to get `"65"` | `"A"` | `strconv.Itoa(65)` |

---

## 8. Exercises

| | Exercise | Practises |
|---|---|---|
| Ex1 | `Shout`: trim, uppercase, add `!` | `strings` functions |
| Ex2 | `IsQuestion`: does the trimmed text end with `?` | `HasSuffix` |
| Ex3 | `Initials`: `"ada lovelace"` → `"AL"` | `Fields`, looping, bytes |
| Ex4 | `CountVowels` | `range` over a string, `switch` |
| Ex5 | `AddStrings`: `"2"` + `"40"` → `42`, or an error flag | `strconv.Atoi`, checking `err` |
| Ex6 | `CharCount`: count characters, not bytes | runes |
| Ex7 | `Reverse` a string, accents and emoji included | `[]rune` |
| Ex8 | `IsPalindrome`, ignoring case and spaces | putting it together |

---

## ✅ Checklist

- [ ] I can find the Go version of common Python string methods (and use `go doc strings`)
- [ ] I can convert text to a number and check `err != nil`
- [ ] I can explain why `len("héllo")` is 6
- [ ] I know what a rune is, and when to use `range` or `[]rune` instead of indexing

**Next → [Lesson 05: Slices](../05_slices/README.md)**
