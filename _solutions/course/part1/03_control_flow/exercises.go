// Reference solution for lesson 03.
package lesson03

func Sign(n int) string {
	if n > 0 {
		return "positive"
	} else if n < 0 {
		return "negative"
	}
	return "zero"
}

func CanVote(age int, isCitizen bool) bool {
	return age >= 18 && isCitizen
}

func SumTo(n int) int {
	total := 0
	for i := 1; i <= n; i++ {
		total += i
	}
	return total
}

func Factorial(n int) int {
	result := 1 // starting at 0 would make every product 0
	for i := 2; i <= n; i++ {
		result *= i
	}
	return result
}

func FizzBuzzWord(n int) string {
	switch {
	case n%3 == 0 && n%5 == 0: // must come first
		return "FizzBuzz"
	case n%3 == 0:
		return "Fizz"
	case n%5 == 0:
		return "Buzz"
	default:
		return ""
	}
}

func DayType(day string) string {
	switch day {
	case "sat", "sun":
		return "weekend"
	case "mon", "tue", "wed", "thu", "fri":
		return "weekday"
	default:
		return "unknown"
	}
}

func CountDigits(n int) int {
	if n == 0 {
		return 1
	}
	count := 0
	for n > 0 {
		n /= 10
		count++
	}
	return count
}

func IsPrime(n int) bool {
	if n < 2 {
		return false
	}
	for d := 2; d*d <= n; d++ {
		if n%d == 0 {
			return false
		}
	}
	return true
}

func CollatzSteps(n int) int {
	steps := 0
	for n != 1 {
		if n%2 == 0 {
			n /= 2
		} else {
			n = 3*n + 1
		}
		steps++
	}
	return steps
}
