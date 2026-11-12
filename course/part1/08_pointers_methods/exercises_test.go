package lesson08

import "testing"

func TestEx1Increment(t *testing.T) {
	x := 5
	Increment(&x)
	Increment(&x)
	if x != 7 {
		t.Errorf("Ex1: after Increment twice, x = %d, want 7", x)
	}
}

func TestEx2SwapInts(t *testing.T) {
	x, y := 1, 2
	SwapInts(&x, &y)
	if x != 2 || y != 1 {
		t.Errorf("Ex2: after SwapInts, x=%d y=%d; want x=2 y=1", x, y)
	}
}

func TestEx3Move(t *testing.T) {
	p := Point{X: 1, Y: 1}
	Move(&p, 2, 3)
	if p != (Point{X: 3, Y: 4}) {
		t.Errorf("Ex3: after Move(&p, 2, 3), p = %+v, want {X:3 Y:4}", p)
	}
}

func TestEx4AreaPerimeter(t *testing.T) {
	r := Rectangle{W: 3, H: 4}
	if got := r.Area(); got != 12 {
		t.Errorf("Ex4: Area() = %v, want 12", got)
	}
	if got := r.Perimeter(); got != 14 {
		t.Errorf("Ex4: Perimeter() = %v, want 14", got)
	}
}

func TestEx5Scale(t *testing.T) {
	r := Rectangle{W: 3, H: 4}
	r.Scale(2)
	if r != (Rectangle{W: 6, H: 8}) {
		t.Errorf("Ex5: after Scale(2), r = %+v, want {W:6 H:8}", r)
	}
}

func TestEx6Counter(t *testing.T) {
	var c Counter
	c.Inc()
	c.Inc()
	c.Add(5)
	if got := c.Value(); got != 7 {
		t.Errorf("Ex6: after Inc, Inc, Add(5): Value() = %d, want 7", got)
	}
	c.Reset()
	if got := c.Value(); got != 0 {
		t.Errorf("Ex6: after Reset: Value() = %d, want 0", got)
	}
}

func TestEx7Account(t *testing.T) {
	a := NewAccount("Ana")
	if a == nil {
		t.Fatal("Ex7: NewAccount returned nil. Return &Account{...}")
	}
	if a.Owner != "Ana" || a.Balance != 0 {
		t.Fatalf("Ex7: NewAccount(\"Ana\") = %+v, want Owner Ana and Balance 0", *a)
	}
	a.Deposit(100)
	a.Deposit(-50) // ignored
	if a.Balance != 100 {
		t.Errorf("Ex7: after Deposit(100) and Deposit(-50): Balance = %d, want 100", a.Balance)
	}
	if !a.Withdraw(30) {
		t.Errorf("Ex7: Withdraw(30) = false, want true")
	}
	if a.Withdraw(500) {
		t.Errorf("Ex7: Withdraw(500) = true, want false (not enough money)")
	}
	if a.Withdraw(0) {
		t.Errorf("Ex7: Withdraw(0) = true, want false")
	}
	if a.Balance != 70 {
		t.Errorf("Ex7: final Balance = %d, want 70", a.Balance)
	}
}

func TestEx8Celsius(t *testing.T) {
	tests := map[Celsius]float64{0: 32, 100: 212, -40: -40, 37: 98.6}
	for c, want := range tests {
		got := c.ToFahrenheit()
		if got < want-1e-9 || got > want+1e-9 {
			t.Errorf("Ex8: Celsius(%v).ToFahrenheit() = %v, want %v", float64(c), got, want)
		}
	}
}
