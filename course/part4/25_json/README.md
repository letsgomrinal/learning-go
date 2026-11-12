# Lesson 25: JSON

**Goal:** turn Go structs into JSON and back, control field names with struct tags, handle nested and
unknown JSON, and stream JSON from a reader.

**Time:** about 1.5 hours.

```bash
cd ~/Desktop/learning-go/course/part4/25_json
```

---

## 1. Python vs Go

In Python, JSON becomes dicts and lists:

```python
import json
data = json.loads('{"name": "Ana", "age": 30}')
print(data["name"])            # a dict: anything goes
```

In Go, you usually decode JSON into a **struct** you define. The compiler then knows exactly what
fields exist and their types, and you get `user.Name` instead of `data["name"]`. (The closest Python
habit is a dataclass or a pydantic model.)

Everything is in the standard package `encoding/json`.

---

## 2. Struct → JSON: `json.Marshal`

```go
type User struct {
	Name  string
	Age   int
	Admin bool
}

u := User{Name: "Ana", Age: 30}
data, err := json.Marshal(u)       // data is []byte
if err != nil {
	return err
}
fmt.Println(string(data))          // {"Name":"Ana","Age":30,"Admin":false}
```

"Marshal" means "convert to JSON". It returns `[]byte`; convert with `string(data)` to print it.

⚠️ **Only exported (capitalized) fields are included.** A field called `password` (lowercase) is
silently left out. That's lesson 09's rule again, and the `json` package can't see private fields.

---

## 3. Struct tags: controlling the JSON names

JSON APIs usually use `lowercase` or `snake_case` keys, but Go fields must be capitalized to be
exported. **Struct tags** connect the two:

```go
type User struct {
	Name     string   `json:"name"`
	Age      int      `json:"age"`
	Email    string   `json:"email,omitempty"`   // left out when empty
	Password string   `json:"-"`                 // never included
	Tags     []string `json:"tags"`
}

u := User{Name: "Ana", Age: 30, Password: "secret"}
data, _ := json.Marshal(u)
fmt.Println(string(data))     // {"name":"Ana","age":30,"tags":null}
```

A tag is the text in backticks after the field type. The `json` package reads it:

| Tag | Meaning |
|---|---|
| `` `json:"name"` `` | use `name` as the JSON key |
| `` `json:"email,omitempty"` `` | skip the field if it has its zero value (`""`, `0`, `false`, nil) |
| `` `json:"-"` `` | never put this field in JSON (or read it from JSON) |

Notice `"tags":null`: a nil slice becomes `null`. If you want `[]`, use an empty slice:
`Tags: []string{}`. (Frontends often prefer `[]`, as you saw in practice set 09.)

---

## 4. JSON → struct: `json.Unmarshal`

```go
text := `{"name": "Bo", "age": 25, "extra": "ignored"}`

var u User
err := json.Unmarshal([]byte(text), &u)      // note the & : it fills u in
if err != nil {
	return err
}
fmt.Println(u.Name, u.Age)                   // Bo 25
```

Things to know:
- Pass a **pointer** (`&u`), so `Unmarshal` can fill in your variable (lesson 08).
- Keys in the JSON that don't match any field (`extra`) are **ignored**.
- Fields missing from the JSON keep their **zero value**.
- Key matching is case-insensitive (`"NAME"` would fill `Name` too), but use the exact tag names anyway.
- A type mismatch is an error: `{"age": "old"}` gives
  `json: cannot unmarshal string into Go struct field User.age of type int`.

---

## 5. Nested JSON

Nested objects become nested structs, and arrays become slices:

```go
type Address struct {
	City string `json:"city"`
	Zip  string `json:"zip"`
}

type Customer struct {
	Name    string    `json:"name"`
	Address Address   `json:"address"`
	Orders  []Order   `json:"orders"`
}

type Order struct {
	ID     int     `json:"id"`
	Amount float64 `json:"amount"`
}
```

That matches JSON like:

```json
{
  "name": "Ana",
  "address": {"city": "Lisbon", "zip": "1000"},
  "orders": [{"id": 1, "amount": 9.99}, {"id": 2, "amount": 20}]
}
```

You only need fields for the parts you care about. Everything else is ignored.

---

## 6. Pretty printing

```go
data, err := json.MarshalIndent(u, "", "  ")    // prefix "", indent with 2 spaces
fmt.Println(string(data))
```

```json
{
  "name": "Ana",
  "age": 30
}
```

---

## 7. JSON with an unknown shape: `map[string]any`

When you really don't know the structure in advance, decode into `map[string]any` (or `any`), like
Python's `json.loads`:

```go
var data map[string]any
json.Unmarshal([]byte(`{"name": "Ana", "age": 30, "langs": ["go", "py"]}`), &data)

name := data["name"].(string)          // type assertions, lesson 11
age := data["age"].(float64)           // ⚠️ float64, not int!
langs := data["langs"].([]any)         // []any, not []string
```

What JSON types become:

| JSON | Go (inside `any`) |
|---|---|
| object | `map[string]any` |
| array | `[]any` |
| string | `string` |
| number | **`float64`** (always, even `30`) |
| true / false | `bool` |
| null | `nil` |

It works, but you lose the compiler's help and need type assertions everywhere. **Prefer structs.**
Use maps only for truly free-form data.

---

## 8. Streaming: `json.Encoder` and `json.Decoder`

`Marshal`/`Unmarshal` work with `[]byte` in memory. To write JSON straight to an `io.Writer`, or read
it from an `io.Reader` (a file, an HTTP body, stdin), use an encoder or decoder:

```go
// write JSON to any Writer (a file, os.Stdout, an HTTP response...)
enc := json.NewEncoder(os.Stdout)
enc.Encode(u)                       // writes the JSON plus a newline

// read JSON from any Reader
dec := json.NewDecoder(file)
var u User
err := dec.Decode(&u)
```

A decoder can also read **many JSON values in a row**, for example a log file with one JSON object per
line ("JSON Lines"):

```go
dec := json.NewDecoder(r)
for {
	var e Event
	err := dec.Decode(&e)
	if errors.Is(err, io.EOF) {
		break                      // no more values: normal end
	}
	if err != nil {
		return err                 // broken JSON
	}
	// use e
}
```

`io.EOF` is the sentinel error meaning "end of input" (lesson 12).

### Rejecting unknown fields

By default, unknown keys are ignored. For config files you may want typos to be errors:

```go
dec := json.NewDecoder(strings.NewReader(text))
dec.DisallowUnknownFields()
err := dec.Decode(&cfg)     // {"prot": 8080} → error: json: unknown field "prot"
```

### 🤔 Predict the output

```go
type Item struct {
	Name  string  `json:"name"`
	Price float64 `json:"price,omitempty"`
	note  string
}

data, _ := json.Marshal(Item{Name: "pen", note: "blue"})
fmt.Println(string(data))

var it Item
json.Unmarshal([]byte(`{"name":"cup","price":2.5,"color":"red"}`), &it)
fmt.Println(it.Name, it.Price)
```

<details><summary>Answer</summary>

```
{"name":"pen"}
cup 2.5
```

`Price` is 0, so `omitempty` drops it. `note` is unexported, so it's invisible to `json`. When
decoding, the unknown key `color` is ignored.
</details>

---

## 9. Common mistakes

| You wrote | What happens | Fix |
|---|---|---|
| Lowercase field names | missing from the JSON, silently | capitalize, and use a tag for the JSON name |
| `json.Unmarshal(data, u)` without `&` | error: `json: Unmarshal(non-pointer lesson25.User)` | `&u` |
| `` `json: "name"` `` (space after the colon) | the tag is ignored, `go vet` warns | `` `json:"name"` `` |
| `data["age"].(int)` on a `map[string]any` | panic: it's a `float64` | `.(float64)`, or use a struct |
| Ignoring the error from `Unmarshal` | half-filled structs | always check it |
| Nil slice when the client expects `[]` | `null` in the output | `[]T{}` |

---

## 10. Exercises

| | Exercise | Practises |
|---|---|---|
| Ex1 | Add tags to `Book`, then `BookJSON` | struct tags, `omitempty`, `-`, `json.Marshal` |
| Ex2 | `ParseBook` | `json.Unmarshal` |
| Ex3 | `TotalSpent` from a JSON array of orders | slices of structs |
| Ex4 | Add tags so `Profile` matches an API's nested JSON | nested structs, snake_case keys |
| Ex5 | `Pretty` | `json.MarshalIndent` |
| Ex6 | `Lookup` a key in unknown JSON | `map[string]any`, `float64` numbers |
| Ex7 | `ReadEvents` from JSON Lines | `json.Decoder`, `io.EOF` |
| Ex8 | `LoadConfig` that rejects unknown keys | `DisallowUnknownFields` |

---

## ✅ Checklist

- [ ] I can marshal a struct to JSON and unmarshal JSON into a struct (with `&`)
- [ ] I can use tags for key names, `omitempty` and `-`
- [ ] I know unexported fields are invisible to `encoding/json`
- [ ] I can map nested JSON to nested structs
- [ ] I know numbers in `map[string]any` are `float64`, and why structs are better
- [ ] I can stream JSON with `json.NewEncoder` / `json.NewDecoder`

**Next → [Lesson 26: HTTP servers](../26_http_server/README.md)**
