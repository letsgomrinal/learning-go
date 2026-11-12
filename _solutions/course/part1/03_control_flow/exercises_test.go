package lesson03

import "testing"

func TestEx1Sign(t *testing.T) {
	tests := map[int]string{5: "positive", -3: "negative", 0: "zero", 1: "positive", -1: "negative"}
	for n, want := range tests {
		if got := Sign(n); got != want {
			t.Errorf("Ex1: Sign(%d) = %q, want %q", n, got, want)
		}
	}
}

func TestEx2CanVote(t *testing.T) {
	tests := []struct {
		age     int
		citizen bool
		want    bool
	}{
		{18, true, true},
		{40, true, true},
		{17, true, false},
		{30, false, false},
		{10, false, false},
	}
	for _, tt := range tests {
		if got := CanVote(tt.age, tt.citizen); got != tt.want {
			t.Errorf("Ex2: CanVote(%d, %v) = %v, want %v", tt.age, tt.citizen, got, tt.want)
		}
	}
}

func TestEx3SumTo(t *testing.T) {
	tests := map[int]int{1: 1, 4: 10, 10: 55, 100: 5050, 0: 0, -5: 0}
	for n, want := range tests {
		if got := SumTo(n); got != want {
			t.Errorf("Ex3: SumTo(%d) = %d, want %d", n, got, want)
		}
	}
}

func TestEx4Factorial(t *testing.T) {
	tests := map[int]int{0: 1, 1: 1, 3: 6, 5: 120, 10: 3628800}
	for n, want := range tests {
		if got := Factorial(n); got != want {
			t.Errorf("Ex4: Factorial(%d) = %d, want %d", n, got, want)
		}
	}
}

func TestEx5FizzBuzzWord(t *testing.T) {
	tests := map[int]string{1: "", 3: "Fizz", 5: "Buzz", 9: "Fizz", 10: "Buzz", 15: "FizzBuzz", 30: "FizzBuzz", 7: ""}
	for n, want := range tests {
		if got := FizzBuzzWord(n); got != want {
			t.Errorf("Ex5: FizzBuzzWord(%d) = %q, want %q", n, got, want)
		}
	}
}

func TestEx6DayType(t *testing.T) {
	tests := map[string]string{
		"sat": "weekend", "sun": "weekend",
		"mon": "weekday", "tue": "weekday", "wed": "weekday", "thu": "weekday", "fri": "weekday",
		"funday": "unknown", "": "unknown",
	}
	for day, want := range tests {
		if got := DayType(day); got != want {
			t.Errorf("Ex6: DayType(%q) = %q, want %q", day, got, want)
		}
	}
}

func TestEx7CountDigits(t *testing.T) {
	tests := map[int]int{0: 1, 7: 1, 10: 2, 99: 2, 4096: 4, 1000000: 7}
	for n, want := range tests {
		if got := CountDigits(n); got != want {
			t.Errorf("Ex7: CountDigits(%d) = %d, want %d", n, got, want)
		}
	}
}

func TestEx8IsPrime(t *testing.T) {
	primes := []int{2, 3, 5, 7, 11, 13, 97}
	notPrimes := []int{-7, 0, 1, 4, 9, 15, 25, 100}
	for _, n := range primes {
		if !IsPrime(n) {
			t.Errorf("Ex8: IsPrime(%d) = false, want true", n)
		}
	}
	for _, n := range notPrimes {
		if IsPrime(n) {
			t.Errorf("Ex8: IsPrime(%d) = true, want false", n)
		}
	}
}

func TestEx9CollatzSteps(t *testing.T) {
	tests := map[int]int{1: 0, 2: 1, 6: 8, 7: 16, 27: 111}
	for n, want := range tests {
		if got := CollatzSteps(n); got != want {
			t.Errorf("Ex9: CollatzSteps(%d) = %d, want %d", n, got, want)
		}
	}
}
