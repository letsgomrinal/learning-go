package lesson01

import (
	"math"
	"testing"
)

func TestEx1ZeroValues(t *testing.T) {
	got := ZeroValues()
	want := `0 "" false`
	if got != want {
		t.Errorf("Ex1: ZeroValues() = %s, want %s", got, want)
	}
}

func TestEx2CelsiusToFahrenheit(t *testing.T) {
	tests := []struct{ c, want float64 }{
		{0, 32},
		{100, 212},
		{37.5, 99.5},
		{-40, -40},
	}
	for _, tt := range tests {
		got := CelsiusToFahrenheit(tt.c)
		if math.Abs(got-tt.want) > 1e-9 {
			t.Errorf("Ex2: CelsiusToFahrenheit(%v) = %v, want %v", tt.c, got, tt.want)
		}
	}
}

func TestEx3Average(t *testing.T) {
	tests := []struct {
		a, b, c int
		want    float64
	}{
		{2, 4, 6, 4},
		{1, 2, 2, 5.0 / 3.0},
		{1, 1, 2, 4.0 / 3.0},
		{0, 0, 1, 1.0 / 3.0},
	}
	for _, tt := range tests {
		got := Average(tt.a, tt.b, tt.c)
		if math.Abs(got-tt.want) > 1e-9 {
			t.Errorf("Ex3: Average(%d, %d, %d) = %v, want %v", tt.a, tt.b, tt.c, got, tt.want)
			if got == math.Trunc(got) && tt.want != math.Trunc(tt.want) {
				t.Log("   Hint: your answer has no decimals. Did integer division happen before the conversion?")
			}
		}
	}
}

func TestEx4SharePerPerson(t *testing.T) {
	tests := []struct{ total, people, want int }{
		{1000, 3, 333},
		{1000, 4, 250},
		{99, 100, 0},
	}
	for _, tt := range tests {
		if got := SharePerPerson(tt.total, tt.people); got != tt.want {
			t.Errorf("Ex4: SharePerPerson(%d, %d) = %d, want %d", tt.total, tt.people, got, tt.want)
		}
	}
}

func TestEx5LeftoverCents(t *testing.T) {
	tests := []struct{ total, people, want int }{
		{1000, 3, 1},
		{1000, 4, 0},
		{1001, 7, 0},
		{99, 100, 99},
	}
	for _, tt := range tests {
		if got := LeftoverCents(tt.total, tt.people); got != tt.want {
			t.Errorf("Ex5: LeftoverCents(%d, %d) = %d, want %d", tt.total, tt.people, got, tt.want)
		}
	}
}

func TestEx6Describe(t *testing.T) {
	tests := []struct {
		name   string
		age    int
		height float64
		want   string
	}{
		{"Ana", 30, 1.68, "Ana is 30 years old and 1.68 m tall"},
		{"Bo", 7, 1.2, "Bo is 7 years old and 1.20 m tall"},
	}
	for _, tt := range tests {
		if got := Describe(tt.name, tt.age, tt.height); got != tt.want {
			t.Errorf("Ex6: Describe(%q, %d, %v)\n   got:  %q\n   want: %q", tt.name, tt.age, tt.height, got, tt.want)
		}
	}
}

func TestEx7AgeInDays(t *testing.T) {
	tests := []struct {
		years int
		want  string
	}{
		{30, "30 years is about 10950 days"},
		{1, "1 years is about 365 days"},
	}
	for _, tt := range tests {
		if got := AgeInDays(tt.years); got != tt.want {
			t.Errorf("Ex7: AgeInDays(%d)\n   got:  %q\n   want: %q", tt.years, got, tt.want)
		}
	}
}
