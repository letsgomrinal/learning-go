package lesson04

import "testing"

func TestEx1Shout(t *testing.T) {
	tests := map[string]string{"  hello  ": "HELLO!", "go": "GO!", "Hi there": "HI THERE!"}
	for in, want := range tests {
		if got := Shout(in); got != want {
			t.Errorf("Ex1: Shout(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestEx2IsQuestion(t *testing.T) {
	tests := map[string]bool{"how are you?  ": true, "fine.": false, "?": true, "what? no": false, "": false}
	for in, want := range tests {
		if got := IsQuestion(in); got != want {
			t.Errorf("Ex2: IsQuestion(%q) = %v, want %v", in, got, want)
		}
	}
}

func TestEx3Initials(t *testing.T) {
	tests := map[string]string{
		"ada lovelace":               "AL",
		"  grace   brewster hopper ": "GBH",
		"Linus":                      "L",
		"":                           "",
	}
	for in, want := range tests {
		if got := Initials(in); got != want {
			t.Errorf("Ex3: Initials(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestEx4CountVowels(t *testing.T) {
	tests := map[string]int{"Hello World": 3, "AEIOU aeiou": 10, "rhythm": 0, "": 0}
	for in, want := range tests {
		if got := CountVowels(in); got != want {
			t.Errorf("Ex4: CountVowels(%q) = %d, want %d", in, got, want)
		}
	}
}

func TestEx5AddStrings(t *testing.T) {
	tests := []struct {
		a, b   string
		want   int
		wantOK bool
	}{
		{"2", "40", 42, true},
		{"-5", "5", 0, true},
		{"2", "forty", 0, false},
		{"two", "2", 0, false},
		{"", "1", 0, false},
	}
	for _, tt := range tests {
		got, ok := AddStrings(tt.a, tt.b)
		if got != tt.want || ok != tt.wantOK {
			t.Errorf("Ex5: AddStrings(%q, %q) = %d, %v; want %d, %v", tt.a, tt.b, got, ok, tt.want, tt.wantOK)
		}
	}
}

func TestEx6CharCount(t *testing.T) {
	tests := map[string]int{"hello": 5, "héllo": 5, "日本語": 3, "Go🙂": 3, "": 0}
	for in, want := range tests {
		if got := CharCount(in); got != want {
			t.Errorf("Ex6: CharCount(%q) = %d, want %d", in, got, want)
			if got == len(in) {
				t.Logf("   Hint: you returned len(s), which counts bytes")
			}
		}
	}
}

func TestEx7Reverse(t *testing.T) {
	tests := map[string]string{"abc": "cba", "héllo": "olléh", "ab🙂": "🙂ba", "": "", "x": "x"}
	for in, want := range tests {
		if got := Reverse(in); got != want {
			t.Errorf("Ex7: Reverse(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestEx8IsPalindrome(t *testing.T) {
	tests := map[string]bool{
		"Never odd or even": true,
		"racecar":           true,
		"Gopher":            false,
		"A Santa at NASA":   true,
		"ab":                false,
	}
	for in, want := range tests {
		if got := IsPalindrome(in); got != want {
			t.Errorf("Ex8: IsPalindrome(%q) = %v, want %v", in, got, want)
		}
	}
}
