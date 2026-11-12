package textutil

import "testing"

func TestCapitalize(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"gopher", "Gopher"},
		{"Gopher", "Gopher"},
		{"", ""},
		{"go is fun", "Go is fun"},
		{"x", "X"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := Capitalize(tt.input); got != tt.want {
				t.Errorf("Capitalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
