package errs

import (
	"errors"
	"slices"
	"strconv"
	"strings"
	"testing"
)

func TestWithdraw(t *testing.T) {
	if got, err := Withdraw(100, 30); err != nil || got != 70 {
		t.Errorf("Withdraw(100, 30) = (%d, %v), want (70, nil)", got, err)
	}
	_, err := Withdraw(100, 150)
	if !errors.Is(err, ErrInsufficientFunds) {
		t.Fatalf("errors.Is(err, ErrInsufficientFunds) = false; err = %v (did you use %%w?)", err)
	}
	if want := "withdraw 150 from 100: insufficient funds"; err.Error() != want {
		t.Errorf("message = %q, want %q", err.Error(), want)
	}
}

func TestParseAge(t *testing.T) {
	if n, err := ParseAge("42"); err != nil || n != 42 {
		t.Errorf("ParseAge(42) = (%d, %v)", n, err)
	}

	_, err := ParseAge("abc")
	if !errors.Is(err, strconv.ErrSyntax) {
		t.Errorf("ParseAge(abc): errors.Is(err, strconv.ErrSyntax) = false; err = %v", err)
	}
	if err != nil && !strings.HasPrefix(err.Error(), `parse age "abc": `) {
		t.Errorf("ParseAge(abc) message = %q", err.Error())
	}

	for _, s := range []string{"-1", "151"} {
		_, err := ParseAge(s)
		var ve *ValidationError
		if !errors.As(err, &ve) {
			t.Fatalf("ParseAge(%s): want *ValidationError, got %v", s, err)
		}
		if ve.Field != "age" || ve.Reason != "out of range" {
			t.Errorf("ParseAge(%s): %+v", s, ve)
		}
		if ve.Error() != "invalid age: out of range" {
			t.Errorf("Error() = %q", ve.Error())
		}
	}
}

func TestRegisterUser(t *testing.T) {
	if err := RegisterUser("ana", "30"); err != nil {
		t.Errorf("valid user: %v", err)
	}

	var ve *ValidationError
	if err := RegisterUser("", "30"); !errors.As(err, &ve) || ve.Field != "name" {
		t.Errorf("empty name: got %v", err)
	}

	err := RegisterUser("bob", "200")
	if !errors.As(err, &ve) || ve.Field != "age" {
		t.Fatalf("bad age: errors.As failed through wrapping; err = %v", err)
	}
	if want := "register bob: invalid age: out of range"; err.Error() != want {
		t.Errorf("message = %q, want %q", err.Error(), want)
	}
}

func TestFindUserNilTrap(t *testing.T) {
	if err := FindUser(1); err != nil {
		t.Errorf("FindUser(1) != nil even though the pointer inside is nil: (%T, %v). Fix the nil-interface bug", err, err)
	}
	var nf *NotFoundError
	if err := FindUser(-5); !errors.As(err, &nf) || nf.ID != -5 {
		t.Errorf("FindUser(-5) = %v", err)
	}
}

func TestSafeDivide(t *testing.T) {
	if r, err := SafeDivide(10, 3); err != nil || r != 3 {
		t.Errorf("SafeDivide(10, 3) = (%d, %v)", r, err)
	}
	r, err := SafeDivide(1, 0)
	if err == nil || !strings.HasPrefix(err.Error(), "recovered: ") {
		t.Errorf("SafeDivide(1, 0) = (%d, %v), want a 'recovered: ...' error", r, err)
	}
}

func TestDeferOrder(t *testing.T) {
	got := DeferOrder()
	want := []string{"start", "end", "deferred 3", "deferred 2", "deferred 1"}
	if !slices.Equal(got, want) {
		t.Errorf("DeferOrder = %q, want %q", got, want)
	}
}
