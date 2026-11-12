package lesson00

import "testing"

func TestEx1HelloWorld(t *testing.T) {
	got := HelloWorld()
	want := "Hello, World!"
	if got != want {
		t.Errorf("Ex1: HelloWorld() = %q, want %q", got, want)
	}
}

func TestEx2Greet(t *testing.T) {
	for _, name := range []string{"Ana", "Gopher"} {
		got := Greet(name)
		want := "Hello, " + name + "!"
		if got != want {
			t.Errorf("Ex2: Greet(%q) = %q, want %q", name, got, want)
		}
	}
}

func TestEx3FixMe(t *testing.T) {
	got := FixMe()
	want := "Go is fun to learn"
	if got != want {
		t.Errorf("Ex3: FixMe() = %q, want %q", got, want)
	}
}

func TestEx4SecondsInADay(t *testing.T) {
	got := SecondsInADay()
	want := 86400
	if got != want {
		t.Errorf("Ex4: SecondsInADay() = %d, want %d", got, want)
	}
}
