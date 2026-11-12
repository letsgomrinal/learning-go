// Lesson 02 exercises: functions. Read README.md first!
//
// Check one exercise at a time:   go test -run Ex1 -v
package lesson02

// Ex1: return the area of a rectangle (width × height).
// RectangleArea(3, 4) == 12
func RectangleArea(width float64, height float64) float64 {
	// TODO
	return 0
}

// Ex2: return the perimeter of a rectangle: 2 × (w + h).
// Perimeter(3, 4) == 14
//
// Notice the signature uses the short form: (w, h float64) means both are float64.
func Perimeter(w, h float64) float64 {
	// TODO
	return 0
}

// Ex3: return the two strings in the opposite order.
// Swap("hello", "world") returns "world", "hello"
func Swap(a, b string) (string, string) {
	// TODO
	return "", ""
}

// Ex4: return the quotient and the remainder of a divided by b.
// Divmod(17, 5) returns 3, 2   (17 = 5×3 + 2)
//
// The results are NAMED (q and r). Set them, then use a bare `return`.
// You can assume b is never 0.
func Divmod(a, b int) (q, r int) {
	// TODO
	return
}

// Ex5: return the sum and the product of three numbers.
// SumAndProduct(2, 3, 4) returns 9, 24
func SumAndProduct(a, b, c int) (int, int) {
	// TODO
	return 0, 0
}

// Ex6: call f on n, then call f again on the result, and return that.
// ApplyTwice(func(x int) int { return x + 3 }, 10) == 16   (10 → 13 → 16)
//
// `f func(int) int` means f is a function that takes an int and returns an int.
// You call it like any function: f(n)
func ApplyTwice(f func(int) int, n int) int {
	// TODO
	return 0
}

// Ex7: return a NEW function that multiplies its input by factor.
//
//	triple := MakeMultiplier(3)
//	triple(5) == 15
//	triple(10) == 30
//
// Hint: return func(n int) int { ... }. The inner function can use `factor`
// (that's a closure, README section 5).
func MakeMultiplier(factor int) func(int) int {
	// TODO: replace nil with a function
	return nil
}
