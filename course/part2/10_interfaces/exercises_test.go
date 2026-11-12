package lesson10

import (
	"math"
	"testing"
)

// asShape checks that v implements Shape. The tests use it so this file compiles
// before you've added your methods.
func asShape(t *testing.T, v any) Shape {
	t.Helper()
	s, ok := v.(Shape)
	if !ok {
		t.Fatalf("%T does not implement Shape yet. Add Area() float64 and Perimeter() float64 methods", v)
	}
	return s
}

func near(a, b float64) bool { return math.Abs(a-b) < 1e-9 }

func TestEx1Square(t *testing.T) {
	s := asShape(t, Square{Side: 3})
	if !near(s.Area(), 9) || !near(s.Perimeter(), 12) {
		t.Errorf("Ex1: Square{3}: Area %v, Perimeter %v; want 9, 12", s.Area(), s.Perimeter())
	}
}

func TestEx2Circle(t *testing.T) {
	c := asShape(t, Circle{R: 2})
	if !near(c.Area(), 4*math.Pi) || !near(c.Perimeter(), 4*math.Pi) {
		t.Errorf("Ex2: Circle{2}: Area %v, Perimeter %v; want %v, %v", c.Area(), c.Perimeter(), 4*math.Pi, 4*math.Pi)
	}
}

func TestEx3TotalArea(t *testing.T) {
	shapes := []Shape{asShape(t, Square{Side: 2}), asShape(t, Circle{R: 1}), asShape(t, Square{Side: 1})}
	if got, want := TotalArea(shapes), 5+math.Pi; !near(got, want) {
		t.Errorf("Ex3: TotalArea = %v, want %v", got, want)
	}
	if got := TotalArea(nil); got != 0 {
		t.Errorf("Ex3: TotalArea(empty) = %v, want 0", got)
	}
}

func TestEx4Biggest(t *testing.T) {
	small := asShape(t, Square{Side: 1})
	big := asShape(t, Circle{R: 2})
	medium := asShape(t, Square{Side: 3})
	got, ok := Biggest([]Shape{small, big, medium})
	if !ok || got != big {
		t.Errorf("Ex4: Biggest = %v, %v; want Circle{R:2}, true", got, ok)
	}
	if got, ok := Biggest(nil); ok || got != nil {
		t.Errorf("Ex4: Biggest(empty) = %v, %v; want nil, false", got, ok)
	}
}

func TestEx5Chorus(t *testing.T) {
	var d any = Dog{}
	dog, ok := d.(Speaker)
	if !ok {
		t.Fatal("Ex5: Dog does not implement Speaker yet. Add func (d Dog) Speak() string")
	}
	r := &Robot{}
	got := Chorus([]Speaker{dog, r, dog})
	if got != "Woof Beep Woof" {
		t.Errorf("Ex5: Chorus = %q, want %q", got, "Woof Beep Woof")
	}
	if r.spoken != 1 {
		t.Errorf("Ex5: the robot spoke %d times, want 1", r.spoken)
	}
}

func TestEx6Checkout(t *testing.T) {
	tests := []struct {
		p    Payer
		amt  int
		want string
	}{
		{PayPal{Email: "ana@example.com"}, 20, "Order paid: sent $20 from ana@example.com"},
		{CreditCard{Last4: "4242"}, 5, "Order paid: charged $5 to card ending 4242"},
	}
	for _, tt := range tests {
		if got := Checkout(tt.p, tt.amt); got != tt.want {
			t.Errorf("Ex6: Checkout(%T, %d)\n   got:  %q\n   want: %q", tt.p, tt.amt, got, tt.want)
		}
	}
}
