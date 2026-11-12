// Lesson 03 exercises: if, for, switch. Read README.md first!
//
// Check one exercise at a time:   go test -run Ex1 -v
package lesson03

// Ex1: return "positive", "negative" or "zero".
func Sign(n int) string {
	// TODO
	return ""
}

// Ex2: a person can vote if they are 18 or older AND a citizen.
func CanVote(age int, isCitizen bool) bool {
	// TODO
	return false
}

// Ex3: return 1 + 2 + 3 + ... + n. For n <= 0, return 0.
// SumTo(4) == 10
//
// Hint: start with total := 0 and add in a loop. (for i := 1; i <= n; i++)
func SumTo(n int) int {
	// TODO
	return 0
}

// Ex4: return n! = 1 × 2 × 3 × ... × n. Factorial(0) is 1.
// Factorial(5) == 120
//
// Hint: start with result := 1 (why not 0?).
func Factorial(n int) int {
	// TODO
	return 0
}

// Ex5: return the FizzBuzz word for ONE number:
// "FizzBuzz" if n is divisible by both 3 and 5, "Fizz" if by 3, "Buzz" if by 5,
// otherwise "" (an empty string).
//
// Use a switch with no value (README section 3). Order matters: which case must come first?
// Hint: "n is divisible by 3" is n%3 == 0
func FizzBuzzWord(n int) string {
	// TODO
	return ""
}

// Ex6: return "weekend" for "sat" and "sun", "weekday" for "mon".."fri",
// and "unknown" for anything else.
//
// Use a switch on the value, with several values in one case: case "sat", "sun":
func DayType(day string) string {
	// TODO
	return ""
}

// Ex7: return how many digits n has. CountDigits(0) == 1, CountDigits(4096) == 4.
// You can assume n >= 0.
//
// Hint: keep dividing by 10 until n is 0, counting the rounds ("while"-style loop: for n > 0 { ... }).
// Careful with 0!
func CountDigits(n int) int {
	// TODO
	return 0
}

// Ex8: return true if n is a prime number (only divisible by 1 and itself).
// Numbers below 2 are not prime. IsPrime(2) == true, IsPrime(9) == false, IsPrime(13) == true
//
// Hint: try every divisor d from 2 while d*d <= n. If n%d == 0, return false straight away.
// If the loop finishes, return true.
func IsPrime(n int) bool {
	// TODO
	return false
}

// Ex9 (a little harder): the Collatz sequence.
// Start from n. If it's even, halve it. If it's odd, make it 3n + 1. Stop when you reach 1.
// Return how many steps it took. CollatzSteps(1) == 0, CollatzSteps(6) == 8
//
//	6 → 3 → 10 → 5 → 16 → 8 → 4 → 2 → 1   (8 steps)
func CollatzSteps(n int) int {
	// TODO
	return 0
}
