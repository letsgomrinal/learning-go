# Lesson 35: Reflection and struct tags

**Goal:** inspect types and values while the program runs with the `reflect` package, read struct tags
like `encoding/json` does, and know when reflection is the right tool (rarely) and what to use instead.

**Time:** about 1.5 hours.

```bash
cd ~/Desktop/learning-go/course/part5/35_reflection
```

---

## 1. What reflection is

In Python, inspecting objects at runtime is everyday stuff: `type(x)`, `getattr(obj, "name")`,
`obj.__dict__`, `dataclasses.fields(obj)`. That's **reflection**: code that looks at the structure of
values while running.

Go is statically typed, so you rarely need it. But some code genuinely must work with **any** type
without knowing it in advance. `encoding/json` is the famous example: how does `json.Marshal` know your
struct's field names and tags? It uses the `reflect` package. This lesson shows you how, so the
"magic" in JSON, database scanners, validators and config loaders stops being magic.

---

## 2. `reflect.TypeOf` and `reflect.ValueOf`

Two entry points:

```go
type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

u := User{Name: "Ana", Age: 30}

t := reflect.TypeOf(u)      // describes the TYPE: User, its fields, their types and tags
v := reflect.ValueOf(u)     // holds the VALUE: the actual data

fmt.Println(t.Name(), t.Kind())     // User struct
fmt.Println(v.Kind())               // struct
```

**Type vs Kind:** the *type* is the specific named type (`User`, `time.Duration`). The *kind* is the
general category (`struct`, `int64`, `slice`, `map`, `pointer`...). `time.Duration` has type `Duration`
and kind `int64`. Reflection code mostly switches on the kind:

```go
switch v.Kind() {
case reflect.String:
	fmt.Println("a string:", v.String())
case reflect.Int, reflect.Int64:
	fmt.Println("an integer:", v.Int())
case reflect.Struct:
	fmt.Println("a struct with", v.NumField(), "fields")
}
```

---

## 3. Walking a struct's fields

```go
t := reflect.TypeOf(u)
v := reflect.ValueOf(u)

for i := range t.NumField() {
	field := t.Field(i)           // a StructField: Name, Type, Tag, IsExported()...
	value := v.Field(i)           // the value of that field in u
	fmt.Printf("%s (%s) = %v, json tag %q\n",
		field.Name, field.Type, value.Interface(), field.Tag.Get("json"))
}
// Name (string) = Ana, json tag "name"
// Age (int) = 30, json tag "age"
```

- `t.Field(i)` describes the field; `v.Field(i)` is its value.
- `value.Interface()` turns a reflect value back into an ordinary `any`.
- `field.Tag.Get("json")` reads a struct tag. That's all a tag is: a string attached to a field that
  any package can read with reflection.
- `field.IsExported()`: reflection can **see** unexported fields, but can't read their values through
  `Interface()` or change them. That's why `json` ignores lowercase fields.

### Pointers: `Elem`

You'll often receive a **pointer** to a struct. `Elem()` goes from the pointer to the thing it points at:

```go
v := reflect.ValueOf(&u)          // Kind: Ptr
if v.Kind() == reflect.Pointer {
	v = v.Elem()                  // now Kind: Struct
}
```

---

## 4. Your own struct tags

Tags are free-form, so you can invent your own. A tag can hold several keys:

```go
type SignUp struct {
	Email string `json:"email" validate:"required"`
	Name  string `json:"name" validate:"min=2"`
	Age   int    `json:"age" validate:"min=18"`
}

field.Tag.Get("validate")    // "required", "min=2", "min=18"
```

This is how validation libraries like `go-playground/validator` work. You'll build a small one in Ex3.
(Tag syntax rule: `key:"value"` pairs separated by spaces, no space after the colon. `go vet` checks it.)

---

## 5. Changing values with reflection

To **set** a field, you must start from a **pointer** (otherwise you'd be changing a copy, lesson 07),
and the field must be exported:

```go
v := reflect.ValueOf(&u).Elem()       // the struct itself, settable
f := v.FieldByName("Age")
if f.IsValid() && f.CanSet() && f.Kind() == reflect.Int {
	f.SetInt(31)
}
fmt.Println(u.Age)                    // 31
```

- `FieldByName` returns an invalid `Value` if there's no such field: check `IsValid()`.
- `CanSet()` is false for unexported fields, or when you didn't start from a pointer.
- Setting the wrong type **panics**: `reflect.Set: value of type string is not assignable to type int`.
  Check with `reflect.ValueOf(x).Type().AssignableTo(f.Type())` first, or use `f.Set(...)` only after
  checking.

---

## 6. When (not) to use reflection

The Go proverb: **"Clear is better than clever. Reflection is never clear."**

Reflection is:
- **slow** compared with normal code (often 10–100× for the reflective part);
- **unchecked**: mistakes the compiler would catch become runtime panics;
- **hard to read**.

Before reaching for it, try:
1. **Interfaces** (lesson 10): if types can implement a method, ask for the method.
2. **Generics** (lesson 15): if the code is the same for many types.
3. **Type switches** (lesson 11): for a handful of known types.
4. **Code generation**: write a tool that generates ordinary Go code at build time. Mark it with a
   comment like `//go:generate stringer -type=Weekday`, run `go generate ./...`, and commit the result.
   `sqlc` (lesson 28), `stringer` and protobuf all work this way. The generated code is fast and checked
   by the compiler.

Use reflection for **general-purpose libraries** that must handle any user type: encoders, decoders,
validators, dependency injection, test helpers. Application code almost never needs it.

### 🤔 Predict the output

```go
type Temp float64

type Reading struct {
	Value  Temp   `unit:"celsius"`
	Sensor string
	secret int
}

t := reflect.TypeOf(Reading{})
for i := range t.NumField() {
	f := t.Field(i)
	fmt.Println(f.Name, f.Type.Kind(), f.Tag.Get("unit"), f.IsExported())
}
```

<details><summary>Answer</summary>

```
Value float64 celsius true
Sensor string  true
secret int  false
```

`Value`'s type is `Temp`, but its **kind** is `float64`. `Sensor` has no `unit` tag, so `Get` returns
`""` (hence the double space). Reflection sees the unexported `secret` field, but reports it as not exported.
</details>

---

## 7. Common mistakes

| You wrote | What happens | Fix |
|---|---|---|
| `reflect.ValueOf(u).Field(0).SetInt(1)` | panic: unaddressable value | start from `reflect.ValueOf(&u).Elem()` |
| Setting an unexported field | panic | check `CanSet()` |
| `v.NumField()` on a pointer | panic: call of NumField on ptr Value | `v.Elem()` first |
| `v.Interface()` on an unexported field | panic | skip fields where `!field.IsExported()` |
| Comparing `Type()` when you meant `Kind()` | `time.Duration` isn't `int64` | switch on `Kind()` for categories |
| Reflection in hot paths "for flexibility" | slow and fragile | interfaces, generics, or code generation |

---

## 8. Exercises

| | Exercise | Practises |
|---|---|---|
| Ex1 | `FieldNames` of any struct (or pointer to one) | `TypeOf`, `Elem`, `NumField`, `IsExported` |
| Ex2 | `KindOf` any value | `Kind()` vs type |
| Ex3 | `Validate` with `validate:"required"` and `validate:"min=N"` tags | reading custom tags, values by kind |
| Ex4 | `SetField` by name, safely | `Elem`, `FieldByName`, `CanSet`, `AssignableTo` |
| Ex5 | `ToMap` using `json` tags | how `encoding/json` works inside |

---

## ✅ Checklist

- [ ] I know the difference between a value's type and its kind
- [ ] I can walk a struct's fields and read tags with `Tag.Get`
- [ ] I know why setting a field needs a pointer and an exported field
- [ ] I can explain roughly how `encoding/json` uses reflection
- [ ] I reach for interfaces, generics or code generation before reflection

**Next → [Lesson 36: Idioms from the standard library](../36_stdlib_idioms/README.md)**
