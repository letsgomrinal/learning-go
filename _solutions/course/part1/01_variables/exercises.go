// Reference solution for lesson 01.
package lesson01

import "fmt"

func ZeroValues() string {
	var n int
	var s string
	var b bool
	return fmt.Sprintf("%d %q %t", n, s, b)
}

func CelsiusToFahrenheit(c float64) float64 {
	return c*9/5 + 32
}

func Average(a, b, c int) float64 {
	return float64(a+b+c) / 3
}

func SharePerPerson(totalCents, people int) int {
	return totalCents / people
}

func LeftoverCents(totalCents, people int) int {
	return totalCents % people
}

func Describe(name string, age int, height float64) string {
	return fmt.Sprintf("%s is %d years old and %.2f m tall", name, age, height)
}

func AgeInDays(years int) string {
	const daysPerYear = 365
	days := years * daysPerYear
	return fmt.Sprintf("%d years is about %d days", years, days)
}
