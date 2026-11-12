// Lesson 05 exercises: slices. Read README.md first!
//
// Check one exercise at a time:   go test -run Ex1 -v
//
// Write these with loops, NOT with the slices package. Writing them yourself is the point.
package lesson05

// Ex1: return the sum of all numbers. Sum of an empty slice is 0.
// Sum([]int{1, 2, 3}) == 6
func Sum(nums []int) int {
	// TODO
	return 0
}

// Ex2: return the largest number and true.
// If the slice is empty there is no largest number, so return 0 and false.
// Largest([]int{3, 9, 2}) == 9, true
// Largest([]int{-5, -2}) == -2, true      ← careful: don't start your "biggest so far" at 0!
// Largest([]int{}) == 0, false
//
// Hint: check len(nums) == 0 first. Then start with biggest := nums[0].
func Largest(nums []int) (int, bool) {
	// TODO
	return 0, false
}

// Ex3: return a new slice with only the even numbers, in the same order.
// Evens([]int{1, 2, 3, 4, 6}) == [2 4 6]
//
// Hint: var result []int, then result = append(result, n) for each even n.
func Evens(nums []int) []int {
	// TODO
	return nil
}

// Ex4: return a NEW slice with every number doubled. Don't change nums!
// Doubled([]int{1, 2, 3}) == [2 4 6]
//
// Hint: result := make([]int, len(nums)), then fill result[i].
func Doubled(nums []int) []int {
	// TODO
	return nil
}

// Ex5: double every number IN PLACE (change the slice you were given). Returns nothing.
//
//	nums := []int{1, 2, 3}
//	DoubleInPlace(nums)
//	// nums is now [2 4 6]
//
// Hint: README section 4. Changing the range copy won't work; use nums[i].
// (The caller sees the change because a slice points at shared data. Lesson 09 explains more.)
func DoubleInPlace(nums []int) {
	// TODO
}

// Ex6: return the index of the first occurrence of target, or -1 if it isn't there.
// IndexOf([]string{"a", "b", "c", "b"}, "b") == 1
func IndexOf(items []string, target string) int {
	// TODO
	return 0
}

// Ex7: return the average of any number of values. With no values, return 0.
// Average(1, 2, 3, 4) == 2.5
//
// `values ...float64` is variadic: inside the function, values is a []float64.
// Hint: divide by float64(len(values)).
func Average(values ...float64) float64 {
	// TODO
	return 0
}

// Ex8: return a NEW slice with the elements in reverse order. Don't change the input.
// Reversed([]int{1, 2, 3}) == [3 2 1]
//
// Hint: loop i from len(s)-1 down to 0:  for i := len(s) - 1; i >= 0; i-- { ... }
func Reversed(s []int) []int {
	// TODO
	return nil
}

// Ex9: return the first and last element, and true.
// For an empty slice, return "", "", false. For one element, first and last are the same.
// FirstAndLast([]string{"x", "y", "z"}) == "x", "z", true
func FirstAndLast(s []string) (string, string, bool) {
	// TODO
	return "", "", false
}
