// Lesson 35 exercises: reflection and struct tags. Read README.md first!
//
// Check one exercise at a time:   go test -run Ex1 -v
package lesson35

import (
	"errors"
	"reflect"
)

var _ = reflect.TypeOf // the reflect package is what you'll use in every exercise

// Ex1: return the names of the EXPORTED fields of a struct, in order.
// v may be a struct or a pointer to a struct. For anything else (including nil), return nil.
//
// FieldNames(User{}) == ["Name", "Email"]   (User also has an unexported field "password")
//
// Hints: t := reflect.TypeOf(v); if t == nil → nil; if t.Kind() == reflect.Pointer { t = t.Elem() };
// check t.Kind() == reflect.Struct; loop over t.NumField() and use t.Field(i).IsExported()
func FieldNames(v any) []string {
	// TODO
	return nil
}

// Ex2: return the kind of v as a string: "int", "string", "struct", "slice", "map", "ptr" ...
// (that's what reflect.Kind's String method returns). For nil, return "nil".
//
// KindOf(time.Second) == "int64"   (its TYPE is time.Duration, but its KIND is int64)
func KindOf(v any) string {
	// TODO
	return ""
}

// Ex3: validate a struct (or a pointer to one) using `validate` tags. Check fields in order and
// return the FIRST problem as an error, or nil if everything is valid. Supported tags:
//
//	validate:"required"  string: must not be ""            → error "<Field>: is required"
//	                     int:    must not be 0             → error "<Field>: is required"
//	validate:"min=N"     string: at least N characters     → error "<Field>: must be at least N characters"
//	                     int:    at least N                → error "<Field>: must be at least N"
//
// Fields without a validate tag are ignored. If v isn't a struct, return an error "not a struct".
//
// Hints:
//   - tag := t.Field(i).Tag.Get("validate")
//   - strings.HasPrefix(tag, "min="), then strconv.Atoi(strings.TrimPrefix(tag, "min="))
//   - value := rv.Field(i); value.Kind() tells you string vs int; value.String(), value.Int()
//   - count characters with utf8.RuneCountInString (lesson 04!)
func Validate(v any) error {
	// TODO
	return errors.New("not implemented")
}

// Ex4: set the field called name in the struct that ptr points to, to value.
// Return an error (don't panic!) if:
//   - ptr isn't a non-nil pointer to a struct            → "need a pointer to a struct"
//   - there's no field with that name                    → "no field <name>"
//   - the field can't be set (unexported)                → "cannot set <name>"
//   - value's type can't be assigned to the field's type → "cannot assign <value type> to <name>"
//
// SetField(&u, "Name", "Bo") sets u.Name to "Bo".
// SetField(&u, "Name", 42) returns "cannot assign int to Name".
//
// Hints: rv := reflect.ValueOf(ptr); rv.Kind() == reflect.Pointer && !rv.IsNil(); s := rv.Elem();
// f := s.FieldByName(name); f.IsValid(); f.CanSet(); val := reflect.ValueOf(value);
// val.Type().AssignableTo(f.Type()); f.Set(val)
func SetField(ptr any, name string, value any) error {
	// TODO
	return errors.New("not implemented")
}

// Ex5: convert a struct (or pointer to one) into a map, like encoding/json does internally:
//
//   - only exported fields
//
//   - the key is the json tag's NAME (the part before any comma), or the field name if there's no tag
//
//   - skip fields tagged json:"-"
//
//   - the value is the field's value (value.Interface())
//
//     type Product struct {
//     ID    int     `json:"id"`
//     Name  string  `json:"name,omitempty"`
//     Price float64
//     cost  float64
//     Notes string  `json:"-"`
//     }
//     ToMap(Product{ID: 1, Name: "pen", Price: 2.5}) == map[id:1 name:pen Price:2.5]
//
// For anything that isn't a struct, return nil. Hint: strings.Cut(tag, ",")
func ToMap(v any) map[string]any {
	// TODO
	return nil
}
