package lesson13

import (
	"errors"
	"fmt"
	"testing"
)

func TestEx1LoadConfig(t *testing.T) {
	if c, err := LoadConfig("app.yaml"); c != "port: 8080" || err != nil {
		t.Errorf("Ex1: LoadConfig(\"app.yaml\") = %q, %v; want \"port: 8080\", nil", c, err)
	}
	_, err := LoadConfig("missing.yaml")
	if err == nil {
		t.Fatal("Ex1: LoadConfig(\"missing.yaml\") returned no error")
	}
	if want := "load config missing.yaml: open missing.yaml: file not found"; err.Error() != want {
		t.Errorf("Ex1: message\n   got:  %q\n   want: %q", err.Error(), want)
	}
	if !errors.Is(err, ErrFileNotFound) {
		t.Errorf("Ex1: errors.Is(err, ErrFileNotFound) = false. Did you wrap with %%w (not %%v)?")
	}
}

func TestEx2ValidationErrorMessage(t *testing.T) {
	e := &ValidationError{Field: "email", Problem: "missing @"}
	if got := e.Error(); got != "invalid email: missing @" {
		t.Errorf("Ex2: Error() = %q, want %q", got, "invalid email: missing @")
	}
}

func TestEx3ValidateUser(t *testing.T) {
	if err := ValidateUser("ana", "ana@example.com"); err != nil {
		t.Errorf("Ex3: valid user returned %v, want nil", err)
	}
	tests := []struct {
		name, email, field, problem string
	}{
		{"", "ana@example.com", "name", "is required"},
		{"ana", "ana.example.com", "email", "missing @"},
		{"", "nope", "name", "is required"}, // name is checked first
	}
	for _, tt := range tests {
		err := ValidateUser(tt.name, tt.email)
		ve, ok := err.(*ValidationError)
		if !ok {
			t.Errorf("Ex3: ValidateUser(%q, %q) = %v (%T), want a *ValidationError", tt.name, tt.email, err, err)
			continue
		}
		if ve.Field != tt.field || ve.Problem != tt.problem {
			t.Errorf("Ex3: ValidateUser(%q, %q) = %+v, want Field %q, Problem %q", tt.name, tt.email, *ve, tt.field, tt.problem)
		}
	}
}

func TestEx4RegisterUser(t *testing.T) {
	if err := RegisterUser("ana", "ana@example.com"); err != nil {
		t.Errorf("Ex4: valid user returned %v", err)
	}
	err := RegisterUser("ana", "nope")
	if err == nil {
		t.Fatal("Ex4: RegisterUser(\"ana\", \"nope\") returned no error")
	}
	if want := "register ana: invalid email: missing @"; err.Error() != want {
		t.Errorf("Ex4: message\n   got:  %q\n   want: %q", err.Error(), want)
	}
	var ve *ValidationError
	if !errors.As(err, &ve) {
		t.Errorf("Ex4: errors.As can't find the *ValidationError. Wrap with %%w")
	}
}

func TestEx5BadField(t *testing.T) {
	wrapped := fmt.Errorf("outer: %w", fmt.Errorf("inner: %w", &ValidationError{Field: "age", Problem: "too low"}))
	if got := BadField(wrapped); got != "age" {
		t.Errorf("Ex5: BadField(two layers of wrapping) = %q, want %q", got, "age")
	}
	if got := BadField(errors.New("boom")); got != "" {
		t.Errorf("Ex5: BadField(plain error) = %q, want \"\"", got)
	}
	if got := BadField(nil); got != "" {
		t.Errorf("Ex5: BadField(nil) = %q, want \"\"", got)
	}
}

func TestEx6FindUserNilTrap(t *testing.T) {
	if err := FindUser(5); err != nil {
		t.Errorf("Ex6: FindUser(5) != nil. It holds (%T, %v): a typed nil inside an interface. Return a plain nil", err, err)
	}
	var nf *NotFoundError
	if err := FindUser(-1); !errors.As(err, &nf) || nf.ID != -1 {
		t.Errorf("Ex6: FindUser(-1) = %v, want a *NotFoundError with ID -1", err)
	}
}

func TestEx7StatusCode(t *testing.T) {
	_, notFound := ReadFile("nope.txt") // a wrapped ErrFileNotFound
	tests := []struct {
		name string
		err  error
		want int
	}{
		{"nil", nil, 200},
		{"wrapped not found", fmt.Errorf("handler: %w", notFound), 404},
		{"validation", &ValidationError{Field: "x", Problem: "y"}, 400},
		{"wrapped validation", fmt.Errorf("register: %w", &ValidationError{Field: "x", Problem: "y"}), 400},
		{"other", errors.New("database down"), 500},
	}
	for _, tt := range tests {
		if got := StatusCode(tt.err); got != tt.want {
			t.Errorf("Ex7: StatusCode(%s) = %d, want %d", tt.name, got, tt.want)
		}
	}
}
