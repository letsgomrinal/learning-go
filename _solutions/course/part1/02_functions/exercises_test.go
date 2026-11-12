package lesson02

import "testing"

func TestEx1RectangleArea(t *testing.T) {
	if got := RectangleArea(3, 4); got != 12 {
		t.Errorf("Ex1: RectangleArea(3, 4) = %v, want 12", got)
	}
	if got := RectangleArea(2.5, 2); got != 5 {
		t.Errorf("Ex1: RectangleArea(2.5, 2) = %v, want 5", got)
	}
}

func TestEx2Perimeter(t *testing.T) {
	if got := Perimeter(3, 4); got != 14 {
		t.Errorf("Ex2: Perimeter(3, 4) = %v, want 14", got)
	}
	if got := Perimeter(0.5, 1); got != 3 {
		t.Errorf("Ex2: Perimeter(0.5, 1) = %v, want 3", got)
	}
}

func TestEx3Swap(t *testing.T) {
	a, b := Swap("hello", "world")
	if a != "world" || b != "hello" {
		t.Errorf("Ex3: Swap(\"hello\", \"world\") = %q, %q; want \"world\", \"hello\"", a, b)
	}
}

func TestEx4Divmod(t *testing.T) {
	tests := []struct{ a, b, q, r int }{
		{17, 5, 3, 2},
		{10, 2, 5, 0},
		{3, 7, 0, 3},
	}
	for _, tt := range tests {
		q, r := Divmod(tt.a, tt.b)
		if q != tt.q || r != tt.r {
			t.Errorf("Ex4: Divmod(%d, %d) = %d, %d; want %d, %d", tt.a, tt.b, q, r, tt.q, tt.r)
		}
	}
}

func TestEx5SumAndProduct(t *testing.T) {
	s, p := SumAndProduct(2, 3, 4)
	if s != 9 || p != 24 {
		t.Errorf("Ex5: SumAndProduct(2, 3, 4) = %d, %d; want 9, 24", s, p)
	}
	s, p = SumAndProduct(-1, 0, 5)
	if s != 4 || p != 0 {
		t.Errorf("Ex5: SumAndProduct(-1, 0, 5) = %d, %d; want 4, 0", s, p)
	}
}

func TestEx6ApplyTwice(t *testing.T) {
	add3 := func(x int) int { return x + 3 }
	if got := ApplyTwice(add3, 10); got != 16 {
		t.Errorf("Ex6: ApplyTwice(add3, 10) = %d, want 16", got)
	}
	square := func(x int) int { return x * x }
	if got := ApplyTwice(square, 3); got != 81 {
		t.Errorf("Ex6: ApplyTwice(square, 3) = %d, want 81", got)
	}
}

func TestEx7MakeMultiplier(t *testing.T) {
	triple := MakeMultiplier(3)
	if triple == nil {
		t.Fatal("Ex7: MakeMultiplier(3) returned nil. Return a func(n int) int { ... }")
	}
	if got := triple(5); got != 15 {
		t.Errorf("Ex7: MakeMultiplier(3)(5) = %d, want 15", got)
	}
	double := MakeMultiplier(2)
	if got := double(21); got != 42 {
		t.Errorf("Ex7: MakeMultiplier(2)(21) = %d, want 42", got)
	}
	if got := triple(10); got != 30 {
		t.Errorf("Ex7: triple still works after making double: triple(10) = %d, want 30", got)
	}
}
