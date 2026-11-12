// Reference solution for lesson 05.
package lesson05

func Sum(nums []int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

func Largest(nums []int) (int, bool) {
	if len(nums) == 0 {
		return 0, false
	}
	biggest := nums[0]
	for _, n := range nums {
		if n > biggest {
			biggest = n
		}
	}
	return biggest, true
}

func Evens(nums []int) []int {
	var result []int
	for _, n := range nums {
		if n%2 == 0 {
			result = append(result, n)
		}
	}
	return result
}

func Doubled(nums []int) []int {
	result := make([]int, len(nums))
	for i, n := range nums {
		result[i] = n * 2
	}
	return result
}

func DoubleInPlace(nums []int) {
	for i := range nums {
		nums[i] *= 2
	}
}

func IndexOf(items []string, target string) int {
	for i, item := range items {
		if item == target {
			return i
		}
	}
	return -1
}

func Average(values ...float64) float64 {
	if len(values) == 0 {
		return 0
	}
	total := 0.0
	for _, v := range values {
		total += v
	}
	return total / float64(len(values))
}

func Reversed(s []int) []int {
	result := make([]int, 0, len(s))
	for i := len(s) - 1; i >= 0; i-- {
		result = append(result, s[i])
	}
	return result
}

func FirstAndLast(s []string) (string, string, bool) {
	if len(s) == 0 {
		return "", "", false
	}
	return s[0], s[len(s)-1], true
}
