// Reference solution for lesson 01. Compare with yours AFTER your tests pass.
package basics

import (
	"strconv"
	"strings"
)

func FizzBuzz(n int) []string {
	out := make([]string, 0, n)
	for i := 1; i <= n; i++ {
		switch {
		case i%15 == 0:
			out = append(out, "FizzBuzz")
		case i%3 == 0:
			out = append(out, "Fizz")
		case i%5 == 0:
			out = append(out, "Buzz")
		default:
			out = append(out, strconv.Itoa(i))
		}
	}
	return out
}

func ReverseRunes(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

func CountWords(s string) map[string]int {
	counts := make(map[string]int)
	for _, w := range strings.Fields(s) {
		counts[strings.ToLower(w)]++
	}
	return counts
}

func Stats(nums ...int) (minV, maxV int, avg float64) {
	if len(nums) == 0 {
		return // named results are already 0, 0, 0
	}
	minV, maxV = nums[0], nums[0]
	sum := 0
	for _, n := range nums {
		minV = min(minV, n)
		maxV = max(maxV, n)
		sum += n
	}
	avg = float64(sum) / float64(len(nums)) // convert BEFORE dividing
	return
}

func Divide(a, b int) (int, bool) {
	if b == 0 {
		return 0, false
	}
	return a / b, true
}

func Grade(score int) string {
	switch {
	case score < 0 || score > 100:
		return "invalid"
	case score >= 90:
		return "A"
	case score >= 80:
		return "B"
	case score >= 70:
		return "C"
	case score >= 60:
		return "D"
	default:
		return "F"
	}
}

type Weekday int

const (
	Sunday Weekday = iota
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
)

func (d Weekday) IsWeekend() bool {
	return d == Saturday || d == Sunday
}
