package lesson12

import (
	"errors"
	"slices"
	"strconv"
	"testing"
)

func TestEx1ParseAge(t *testing.T) {
	if n, err := ParseAge("42"); n != 42 || err != nil {
		t.Errorf("Ex1: ParseAge(\"42\") = %d, %v; want 42, nil", n, err)
	}
	n, err := ParseAge("abc")
	if err == nil {
		t.Fatalf("Ex1: ParseAge(\"abc\") returned no error (n = %d)", n)
	}
	if want := `not a number: "abc"`; err.Error() != want || n != 0 {
		t.Errorf("Ex1: ParseAge(\"abc\") = %d, %q; want 0, %q", n, err, want)
	}
	n, err = ParseAge("-3")
	if err == nil || err.Error() != "age cannot be negative" || n != 0 {
		t.Errorf("Ex1: ParseAge(\"-3\") = %d, %v; want 0, \"age cannot be negative\"", n, err)
	}
	if n, err := ParseAge("0"); n != 0 || err != nil {
		t.Errorf("Ex1: ParseAge(\"0\") = %d, %v; want 0, nil", n, err)
	}
}

func TestEx2SafeDivide(t *testing.T) {
	if q, err := SafeDivide(10, 3); q != 3 || err != nil {
		t.Errorf("Ex2: SafeDivide(10, 3) = %d, %v; want 3, nil", q, err)
	}
	q, err := SafeDivide(1, 0)
	if !errors.Is(err, ErrDivideByZero) || q != 0 {
		t.Errorf("Ex2: SafeDivide(1, 0) = %d, %v; want 0, ErrDivideByZero", q, err)
	}
}

func TestEx3SumStrings(t *testing.T) {
	if s, err := SumStrings([]string{"1", "2", "3"}); s != 6 || err != nil {
		t.Errorf("Ex3: SumStrings([1 2 3]) = %d, %v; want 6, nil", s, err)
	}
	if s, err := SumStrings(nil); s != 0 || err != nil {
		t.Errorf("Ex3: SumStrings(empty) = %d, %v; want 0, nil", s, err)
	}
	s, err := SumStrings([]string{"1", "two", "3"})
	if err == nil || s != 0 {
		t.Fatalf("Ex3: SumStrings([1 two 3]) = %d, %v; want 0 and an error", s, err)
	}
	if !errors.Is(err, strconv.ErrSyntax) {
		t.Errorf("Ex3: return the error from strconv.Atoi unchanged; got %v", err)
	}
}

func TestEx4Withdraw(t *testing.T) {
	tests := []struct {
		balance, amount, want int
		wantErr               error
	}{
		{100, 30, 70, nil},
		{100, 100, 0, nil},
		{100, 0, 100, ErrInvalidAmount},
		{100, -5, 100, ErrInvalidAmount},
		{50, 80, 50, ErrInsufficientFunds},
	}
	for _, tt := range tests {
		got, err := Withdraw(tt.balance, tt.amount)
		if got != tt.want || !errors.Is(err, tt.wantErr) || (tt.wantErr == nil && err != nil) {
			t.Errorf("Ex4: Withdraw(%d, %d) = %d, %v; want %d, %v", tt.balance, tt.amount, got, err, tt.want, tt.wantErr)
		}
	}
}

func TestEx5WithdrawMessage(t *testing.T) {
	tests := []struct {
		balance, amount int
		want            string
	}{
		{100, 30, "new balance: 70"},
		{50, 80, "sorry, not enough money"},
		{50, -1, "please enter a positive amount"},
	}
	for _, tt := range tests {
		if got := WithdrawMessage(tt.balance, tt.amount); got != tt.want {
			t.Errorf("Ex5: WithdrawMessage(%d, %d) = %q, want %q", tt.balance, tt.amount, got, tt.want)
		}
	}
}

func TestEx6ParseAll(t *testing.T) {
	nums, bad := ParseAll([]string{"1", "x", "3", "y", "-2"})
	if want := []int{1, 3, -2}; !slices.Equal(nums, want) {
		t.Errorf("Ex6: nums = %v, want %v", nums, want)
	}
	if want := []string{"x", "y"}; !slices.Equal(bad, want) {
		t.Errorf("Ex6: bad = %q, want %q", bad, want)
	}
}

func TestEx7UserName(t *testing.T) {
	users := map[int]string{1: "ana", 2: "bo"}
	if got := UserName(users, 2); got != "bo" {
		t.Errorf("Ex7: UserName(2) = %q, want %q", got, "bo")
	}
	if got := UserName(users, 99); got != "guest" {
		t.Errorf("Ex7: UserName(99) = %q, want %q", got, "guest")
	}
}
