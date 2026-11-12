// Reference solution for lesson 35.
package lesson35

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unicode/utf8"
)

// structValue returns v (following one pointer) as a struct Value, or ok=false.
func structValue(v any) (reflect.Value, bool) {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Pointer {
		if rv.IsNil() {
			return reflect.Value{}, false
		}
		rv = rv.Elem()
	}
	return rv, rv.Kind() == reflect.Struct
}

func FieldNames(v any) []string {
	t := reflect.TypeOf(v)
	if t == nil {
		return nil
	}
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil
	}
	var names []string
	for i := range t.NumField() {
		if f := t.Field(i); f.IsExported() {
			names = append(names, f.Name)
		}
	}
	return names
}

func KindOf(v any) string {
	if v == nil {
		return "nil"
	}
	return reflect.TypeOf(v).Kind().String()
}

func Validate(v any) error {
	rv, ok := structValue(v)
	if !ok {
		return errors.New("not a struct")
	}
	t := rv.Type()
	for i := range t.NumField() {
		field := t.Field(i)
		tag := field.Tag.Get("validate")
		if tag == "" || !field.IsExported() {
			continue
		}
		value := rv.Field(i)
		switch {
		case tag == "required":
			if value.IsZero() {
				return fmt.Errorf("%s: is required", field.Name)
			}
		case strings.HasPrefix(tag, "min="):
			min, err := strconv.Atoi(strings.TrimPrefix(tag, "min="))
			if err != nil {
				return fmt.Errorf("%s: bad tag %q", field.Name, tag)
			}
			switch value.Kind() {
			case reflect.String:
				if utf8.RuneCountInString(value.String()) < min {
					return fmt.Errorf("%s: must be at least %d characters", field.Name, min)
				}
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				if value.Int() < int64(min) {
					return fmt.Errorf("%s: must be at least %d", field.Name, min)
				}
			}
		}
	}
	return nil
}

func SetField(ptr any, name string, value any) error {
	rv := reflect.ValueOf(ptr)
	if rv.Kind() != reflect.Pointer || rv.IsNil() || rv.Elem().Kind() != reflect.Struct {
		return errors.New("need a pointer to a struct")
	}
	f := rv.Elem().FieldByName(name)
	if !f.IsValid() {
		return fmt.Errorf("no field %s", name)
	}
	if !f.CanSet() {
		return fmt.Errorf("cannot set %s", name)
	}
	val := reflect.ValueOf(value)
	if !val.IsValid() || !val.Type().AssignableTo(f.Type()) {
		return fmt.Errorf("cannot assign %T to %s", value, name)
	}
	f.Set(val)
	return nil
}

func ToMap(v any) map[string]any {
	rv, ok := structValue(v)
	if !ok {
		return nil
	}
	t := rv.Type()
	out := map[string]any{}
	for i := range t.NumField() {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}
		key := field.Name
		if tag := field.Tag.Get("json"); tag != "" {
			name, _, _ := strings.Cut(tag, ",")
			if name == "-" {
				continue
			}
			if name != "" {
				key = name
			}
		}
		out[key] = rv.Field(i).Interface()
	}
	return out
}
