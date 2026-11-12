package lesson35

import (
	"maps"
	"slices"
	"testing"
	"time"
)

type User struct {
	Name     string
	Email    string
	password string
}

func TestEx1FieldNames(t *testing.T) {
	if got := FieldNames(User{}); !slices.Equal(got, []string{"Name", "Email"}) {
		t.Errorf("Ex1: FieldNames(User{}) = %q, want [Name Email]", got)
	}
	if got := FieldNames(&User{}); !slices.Equal(got, []string{"Name", "Email"}) {
		t.Errorf("Ex1: FieldNames(&User{}) = %q, want [Name Email] (follow the pointer with Elem)", got)
	}
	for _, v := range []any{42, "x", nil, []int{1}} {
		if got := FieldNames(v); got != nil {
			t.Errorf("Ex1: FieldNames(%#v) = %q, want nil", v, got)
		}
	}
}

func TestEx2KindOf(t *testing.T) {
	tests := []struct {
		v    any
		want string
	}{
		{42, "int"},
		{"hi", "string"},
		{time.Second, "int64"},
		{User{}, "struct"},
		{&User{}, "ptr"},
		{[]int{}, "slice"},
		{map[string]int{}, "map"},
		{nil, "nil"},
	}
	for _, tt := range tests {
		if got := KindOf(tt.v); got != tt.want {
			t.Errorf("Ex2: KindOf(%#v) = %q, want %q", tt.v, got, tt.want)
		}
	}
}

type SignUp struct {
	Email string `validate:"required"`
	Name  string `validate:"min=2"`
	Age   int    `validate:"min=18"`
	Code  int    `validate:"required"`
	Notes string
}

func TestEx3Validate(t *testing.T) {
	valid := SignUp{Email: "a@b.c", Name: "Jo", Age: 18, Code: 7}
	if err := Validate(valid); err != nil {
		t.Errorf("Ex3: valid sign-up returned %v", err)
	}
	if err := Validate(&valid); err != nil {
		t.Errorf("Ex3: a pointer to a valid sign-up returned %v", err)
	}
	tests := []struct {
		in   SignUp
		want string
	}{
		{SignUp{Name: "Jo", Age: 18, Code: 7}, "Email: is required"},
		{SignUp{Email: "a@b.c", Name: "J", Age: 18, Code: 7}, "Name: must be at least 2 characters"},
		{SignUp{Email: "a@b.c", Name: "Zoë", Age: 17, Code: 7}, "Age: must be at least 18"},
		{SignUp{Email: "a@b.c", Name: "Jo", Age: 30}, "Code: is required"},
		{SignUp{}, "Email: is required"}, // the first problem, in field order
	}
	for _, tt := range tests {
		err := Validate(tt.in)
		if err == nil || err.Error() != tt.want {
			t.Errorf("Ex3: Validate(%+v) = %v, want %q", tt.in, err, tt.want)
		}
	}
	if err := Validate(42); err == nil || err.Error() != "not a struct" {
		t.Errorf("Ex3: Validate(42) = %v, want \"not a struct\"", err)
	}
}

func TestEx4SetField(t *testing.T) {
	u := User{Name: "Ana"}
	if err := SetField(&u, "Name", "Bo"); err != nil || u.Name != "Bo" {
		t.Errorf("Ex4: SetField(&u, Name, Bo): err %v, u.Name %q", err, u.Name)
	}
	tests := []struct {
		ptr   any
		name  string
		value any
		want  string
	}{
		{u, "Name", "x", "need a pointer to a struct"},
		{(*User)(nil), "Name", "x", "need a pointer to a struct"},
		{&u, "Phone", "x", "no field Phone"},
		{&u, "password", "x", "cannot set password"},
		{&u, "Name", 42, "cannot assign int to Name"},
	}
	for _, tt := range tests {
		var err error
		func() {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("Ex4: SetField(%T, %q, %v) panicked: %v. Check before setting", tt.ptr, tt.name, tt.value, r)
				}
			}()
			err = SetField(tt.ptr, tt.name, tt.value)
		}()
		if err == nil || err.Error() != tt.want {
			t.Errorf("Ex4: SetField(%T, %q, %v) = %v, want %q", tt.ptr, tt.name, tt.value, err, tt.want)
		}
	}
}

type Product struct {
	ID    int    `json:"id"`
	Name  string `json:"name,omitempty"`
	Price float64
	cost  float64
	Notes string `json:"-"`
}

func TestEx5ToMap(t *testing.T) {
	p := Product{ID: 1, Name: "pen", Price: 2.5, cost: 1, Notes: "secret"}
	want := map[string]any{"id": 1, "name": "pen", "Price": 2.5}
	if got := ToMap(p); !maps.Equal(got, want) {
		t.Errorf("Ex5: ToMap = %v, want %v", got, want)
	}
	if got := ToMap(&p); !maps.Equal(got, want) {
		t.Errorf("Ex5: ToMap(pointer) = %v, want %v", got, want)
	}
	if got := ToMap("nope"); got != nil {
		t.Errorf("Ex5: ToMap(string) = %v, want nil", got)
	}
}
