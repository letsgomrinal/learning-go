// Lesson 15 exercises: generics. Read README.md first!
//
// Check one exercise at a time:   go test -run Ex1 -v
package lesson15

import "cmp"

// Ex1: return true if target is in items. Works for any comparable type.
// Contains([]string{"a", "b"}, "b") == true, Contains([]int{1, 2}, 5) == false
func Contains[T comparable](items []T, target T) bool {
	// TODO
	return false
}

// Ex2: return the biggest value and true. For an empty slice, return the zero value of T and false.
// Max([]int{3, 9, 2}) == 9, true      Max([]string{"go", "py"}) == "py", true
//
// Hint: for the empty case, declare  var zero T  and return it.
func Max[T cmp.Ordered](values []T) (T, bool) {
	// TODO
	var zero T
	return zero, false
}

// Ex3: return a new slice with f applied to every item.
// Map([]string{"go", "rust"}, func(s string) int { return len(s) }) == [2 4]
func Map[T, U any](items []T, f func(T) U) []U {
	// TODO
	return nil
}

// Ex4: return a new slice with only the items for which keep returns true, in order.
// Filter([]int{1, 2, 3, 4}, func(n int) bool { return n%2 == 0 }) == [2 4]
func Filter[T any](items []T, keep func(T) bool) []T {
	// TODO
	return nil
}

// Number is a constraint. Right now it only allows exactly int and float64.
//
// Ex5: change it so types BUILT ON int and float64 are allowed too (like Celsius below),
// then write Sum. (README section 3: the ~.)
// When that's done, uncomment the three Celsius lines at the end of TestEx5Sum in exercises_test.go.
// Before the ~ change they wouldn't compile. Try it, and read the error!
type Number interface {
	int | float64
}

// Sum adds up the values. Sum of an empty slice is 0.
// Hint: var total T, then total += v.
func Sum[T Number](values []T) T {
	// TODO
	var total T
	return total
}

// Celsius is a named type built on float64.
type Celsius float64

// Ex6: a stack that holds any type. Fill in the methods.
//
//	var s Stack[string]
//	s.Push("a"); s.Push("b")
//	s.Pop() == "b", true
//	s.Len() == 1
type Stack[T any] struct {
	items []T
}

// Push adds v on top.
func (s *Stack[T]) Push(v T) {
	// TODO
}

// Pop removes and returns the top item, and true. If the stack is empty, return the zero value and false.
func (s *Stack[T]) Pop() (T, bool) {
	// TODO
	var zero T
	return zero, false
}

// Peek returns the top item without removing it (and true), or the zero value and false.
func (s *Stack[T]) Peek() (T, bool) {
	// TODO
	var zero T
	return zero, false
}

// Len returns how many items are on the stack.
func (s *Stack[T]) Len() int {
	// TODO
	return 0
}

// Ex7: return the keys of any map, sorted ascending.
// The keys must be cmp.Ordered (so they can be sorted); the values can be anything.
// SortedKeys(map[string]int{"b": 1, "a": 2}) == [a b]
// SortedKeys(map[int]bool{3: true, 1: false}) == [1 3]
//
// Hint: collect the keys into a []K, then slices.Sort (import "slices").
func SortedKeys[K cmp.Ordered, V any](m map[K]V) []K {
	// TODO
	return nil
}
