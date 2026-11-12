package structs

import (
	"errors"
	"slices"
	"testing"
)

func TestPointValueVsPointer(t *testing.T) {
	p := Point{2, 3}
	q := p.ScaledCopy(10)
	if q != (Point{20, 30}) {
		t.Errorf("ScaledCopy = %+v, want {20 30}", q)
	}
	if p != (Point{2, 3}) {
		t.Errorf("ScaledCopy modified the original: %+v", p)
	}

	p.Scale(3) // Go takes &p automatically because Scale has a pointer receiver
	if p != (Point{6, 9}) {
		t.Errorf("after Scale(3): %+v, want {6 9}", p)
	}
}

func TestAccount(t *testing.T) {
	a := NewAccount("ana")
	if a == nil {
		t.Fatal("NewAccount returned nil")
	}
	if a.Owner != "ana" || a.Balance() != 0 {
		t.Fatalf("new account = %+v", a)
	}
	if err := a.Deposit(100); err != nil {
		t.Fatalf("Deposit(100): %v", err)
	}
	if err := a.Deposit(-5); !errors.Is(err, ErrInvalidAmount) {
		t.Errorf("Deposit(-5) err = %v, want ErrInvalidAmount", err)
	}
	if err := a.Withdraw(30); err != nil {
		t.Fatalf("Withdraw(30): %v", err)
	}
	if err := a.Withdraw(500); !errors.Is(err, ErrInsufficientFunds) {
		t.Errorf("Withdraw(500) err = %v, want ErrInsufficientFunds", err)
	}
	if err := a.Withdraw(0); !errors.Is(err, ErrInvalidAmount) {
		t.Errorf("Withdraw(0) err = %v, want ErrInvalidAmount", err)
	}
	if a.Balance() != 70 {
		t.Errorf("Balance = %d, want 70", a.Balance())
	}

	h := a.History()
	if want := []int{100, -30}; !slices.Equal(h, want) {
		t.Fatalf("History = %v, want %v", h, want)
	}
	h[0] = 1_000_000 // try to tamper with it
	if a.History()[0] != 100 {
		t.Errorf("History() leaked the internal slice. Return a copy")
	}
}

func TestDogEmbedding(t *testing.T) {
	d := NewDog("Rex", "Labrador")
	if d.Name != "Rex" || d.Legs != 4 || d.Breed != "Labrador" {
		t.Fatalf("NewDog = %+v", d)
	}
	if got, want := d.Describe(), "Rex has 4 legs"; got != want {
		t.Errorf("Describe = %q, want %q", got, want)
	}
	if got, want := d.Speak(), "Rex the Labrador says woof"; got != want {
		t.Errorf("Speak = %q, want %q", got, want)
	}
}

func TestStackZeroValue(t *testing.T) {
	var s Stack // zero value, no constructor
	if _, ok := s.Pop(); ok {
		t.Error("Pop on empty stack returned ok=true")
	}
	s.Push(1)
	s.Push(2)
	s.Push(3)
	if s.Len() != 3 {
		t.Fatalf("Len = %d, want 3", s.Len())
	}
	if v, ok := s.Peek(); !ok || v != 3 {
		t.Errorf("Peek = (%d, %v), want (3, true)", v, ok)
	}
	for _, want := range []int{3, 2, 1} {
		if v, ok := s.Pop(); !ok || v != want {
			t.Errorf("Pop = (%d, %v), want (%d, true)", v, ok, want)
		}
	}
	if s.Len() != 0 {
		t.Errorf("Len after popping all = %d", s.Len())
	}
}
