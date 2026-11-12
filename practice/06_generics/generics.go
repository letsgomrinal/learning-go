// Lesson 06: Generics and iterators.
//
// Run:  go test -v ./practice/06_generics
//
// ─── CONCEPTS ──────────────────────────────────────────────────────────────
//
// 1. TYPE PARAMETERS (Go 1.18+). Like Python's TypeVar, but enforced and compiled:
//      func Max[T cmp.Ordered](a, b T) T { if a > b { return a }; return b }
//      Max(3, 7)          // T inferred as int
//      Max("a", "b")      // T inferred as string
//
// 2. CONSTRAINTS are interfaces that describe a SET OF TYPES:
//      any               // anything (but you can't do much with it: no +, no <)
//      comparable        // supports == and != (usable as a map key)
//      cmp.Ordered       // supports < <= > >= (ints, floats, strings)
//      type Number interface { ~int | ~int64 | ~float64 }
//    The ~ means "any type whose UNDERLYING type is int", so `type Celsius int` matches ~int.
//
// 3. GENERIC TYPES
//      type Stack[T any] struct { items []T }
//      func (s *Stack[T]) Push(v T) { ... }
//      s := &Stack[string]{}
//    (Go 1.27+ also allows methods with their own extra type parameters.)
//
// 4. WHEN TO USE GENERICS vs INTERFACES
//    - Generics: same algorithm, different element types (containers, Map/Filter, Sum).
//    - Interfaces: different BEHAVIOUR behind the same method set (Shape, io.Reader).
//    Don't reach for generics first. Go code was fine without them for 10 years.
//
// 5. ITERATORS (Go 1.23+): range over functions, Go's answer to Python generators.
//      type Seq[V any] func(yield func(V) bool)   // from package iter
//      func Count(n int) iter.Seq[int] {
//          return func(yield func(int) bool) {
//              for i := range n {
//                  if !yield(i) { return }   // consumer did `break`: stop!
//              }
//          }
//      }
//      for v := range Count(3) { ... }          // 0, 1, 2
//    yield returns false when the loop body breaks. You MUST stop then, or it panics.
//    Compare with Python:  def count(n): for i in range(n): yield i
// ───────────────────────────────────────────────────────────────────────────

package generics

import (
	"iter"
)

// Map applies f to every element.  Map([1 2 3], double) == [2 4 6]
// (Python: [f(x) for x in s])
func Map[T, U any](s []T, f func(T) U) []U {
	// TODO
	return nil
}

// Filter keeps the elements where keep returns true.
func Filter[T any](s []T, keep func(T) bool) []T {
	// TODO
	return nil
}

// Reduce folds s into a single value, starting from init.
// Reduce([1 2 3], 0, add) == 6   (Python: functools.reduce)
func Reduce[T, A any](s []T, init A, f func(A, T) A) A {
	// TODO
	return init
}

// Number is a constraint for numeric types, including named types based on them.
type Number interface {
	~int | ~int32 | ~int64 | ~float32 | ~float64
}

// Sum adds up the numbers. Must work for []int, []float64, and []Celsius.
func Sum[T Number](s []T) T {
	// TODO
	var zero T
	return zero
}

// Celsius is a named type. Sum works with it thanks to ~float64 in Number.
type Celsius float64

// ── A generic Set ────────────────────────────────────────────────────────

// Set is an unordered collection of unique values.
// T must be comparable because it's used as a map key.
type Set[T comparable] struct {
	m map[T]struct{}
}

// NewSet creates a set holding the given items.
func NewSet[T comparable](items ...T) *Set[T] {
	// TODO
	return &Set[T]{}
}

func (s *Set[T]) Add(v T) {
	// TODO (careful: writing to a nil map panics!)
}

func (s *Set[T]) Has(v T) bool {
	// TODO
	return false
}

func (s *Set[T]) Len() int {
	// TODO
	return 0
}

// Union returns a NEW set with the elements of both sets.
func (s *Set[T]) Union(other *Set[T]) *Set[T] {
	// TODO
	return &Set[T]{}
}

// All returns an iterator over the set's elements (in any order), so callers can write:
//
//	for v := range s.All() { ... }
func (s *Set[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		// TODO: loop over s.m and yield each key; stop if yield returns false.
	}
}

// ── Iterators ────────────────────────────────────────────────────────────

// Fibonacci returns an INFINITE iterator over 0, 1, 1, 2, 3, 5, 8...
// It's only safe because the consumer breaks out of the loop, and you stop when
// yield returns false. Python equivalent:
//
//	def fib():
//	    a, b = 0, 1
//	    while True:
//	        yield a
//	        a, b = b, a + b
func Fibonacci() iter.Seq[int] {
	return func(yield func(int) bool) {
		// TODO
	}
}

// Take collects the first n values from any iterator into a slice.
// Hint: for v := range seq { ...; if len(out) == n { break } }
func Take[T any](seq iter.Seq[T], n int) []T {
	// TODO
	return nil
}
