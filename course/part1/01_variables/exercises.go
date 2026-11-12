// Lesson 01 exercises: variables and types. Read README.md first!
//
// Check one exercise at a time:   go test -run Ex1 -v
package lesson01

import "fmt"

// Ex1: show the zero values.
//
// Inside the function:
//  1. declare three variables WITHOUT giving them a value:
//     var n int
//     var s string
//     var b bool
//  2. return fmt.Sprintf("%d %q %t", n, s, b)
//
// Before you run the test, guess what the string will be!
func ZeroValues() string {
	// TODO
	return ""
}

// Ex2: convert a temperature from Celsius to Fahrenheit.
// The formula is:  F = C × 9 / 5 + 32
// CelsiusToFahrenheit(100) == 212, CelsiusToFahrenheit(37.5) == 99.5
//
// Hint: c is already a float64, so the math is decimal math. Just write the formula.
func CelsiusToFahrenheit(c float64) float64 {
	// TODO
	return 0
}

// Ex3: return the average of three whole numbers, as a decimal.
// Average(1, 2, 2) == 1.6666666666666667, Average(2, 4, 6) == 4
//
// Hint: a+b+c is an int. Convert it with float64(...) BEFORE dividing by 3,
// or you'll get integer division (lesson section 8, "The average trap").
func Average(a, b, c int) float64 {
	// TODO
	return 0
}

// Ex4: a bill of totalCents is split between some people. How many whole cents does each pay?
// SharePerPerson(1000, 3) == 333   (1000 cents = $10.00, three people)
//
// Hint: with two ints, / already does whole-number division. That's what we want here!
func SharePerPerson(totalCents, people int) int {
	// TODO
	return 0
}

// Ex5: how many cents are left over after splitting evenly?
// LeftoverCents(1000, 3) == 1   (3 × 333 = 999, so 1 cent is left)
//
// Hint: the remainder operator.
func LeftoverCents(totalCents, people int) int {
	// TODO
	return 0
}

// Ex6: build a sentence with fmt.Sprintf.
// Describe("Ana", 30, 1.68) == "Ana is 30 years old and 1.68 m tall"
// Describe("Bo", 7, 1.2)    == "Bo is 7 years old and 1.20 m tall"
//
// Hint: %s for the string, %d for the int, %.2f for the float with 2 decimals.
func Describe(name string, age int, height float64) string {
	// TODO
	return fmt.Sprintf("")
}

// Ex7: return a sentence like "30 years is about 10950 days".
//
// Steps:
//  1. declare a constant inside the function:  const daysPerYear = 365
//  2. compute the days with a new variable:     days := years * daysPerYear
//  3. return fmt.Sprintf(...) with two %d verbs
func AgeInDays(years int) string {
	// TODO
	return ""
}
