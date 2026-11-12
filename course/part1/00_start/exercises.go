// Lesson 00 exercises. Read README.md first!
//
// Check one exercise at a time:   go test -run Ex1 -v
package lesson00

// Ex1: return the text "Hello, World!" (exactly, with the comma and the exclamation mark).
func HelloWorld() string {
	// TODO: replace the empty string
	return "Hello, World!"
}

// Ex2: return "Hello, " + name + "!"
// Greet("Ana") should return "Hello, Ana!"
//
// `name string` means: this function receives one value, called name, and it's a string.
// You can use name like a variable.
func Greet(name string) string {
	// TODO
	return "Hello, " + name + "!"
}

// Ex3: this one is already written, but the text has a typo.
// Run `go test -run Ex3 -v`, compare what you got with what was wanted, and fix it.
func FixMe() string {
	return "Go is fun to learn"
}

// Ex4: return the number of seconds in one day.
// Write it as a calculation (hours × minutes × seconds), not as a number you worked out.
// `int` means the function returns a whole number.
func SecondsInADay() int {
	// TODO
	return 24 * 60 * 60
}
