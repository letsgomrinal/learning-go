// Lesson 06 exercises: maps. Read README.md first!
//
// Check one exercise at a time:   go test -run Ex1 -v
package lesson06

// Ex1: look up item in prices. Return the price and true, or 0 and false if it isn't there.
// Use the "comma ok" form:  price, ok := prices[item]
//
// (Yes, you can solve it in one line. That's the point: comma ok is short.)
func Lookup(prices map[string]float64, item string) (float64, bool) {
	// TODO
	return 0, false
}

// Ex2: count how many times each word appears.
// WordCount([]string{"go", "is", "go"}) == map[go:2 is:1]
//
// Hint: create the map with map[string]int{}, then counts[w]++ in a loop.
func WordCount(words []string) map[string]int {
	// TODO
	return nil
}

// Ex3: return true if any number appears more than once.
// HasDuplicates([]int{1, 2, 3, 2}) == true
//
// Hint: a set, seen := map[int]bool{}. For each number: if it's already in seen, return true.
// Otherwise add it.
func HasDuplicates(nums []int) bool {
	// TODO
	return false
}

// Ex4: return the words without duplicates, keeping the order of first appearance.
// Unique([]string{"b", "a", "b", "c", "a"}) == [b a c]
//
// Hint: a set to remember what you've seen, and a slice for the result
// (a map alone can't keep the order. Why?).
func Unique(words []string) []string {
	// TODO
	return nil
}

// Ex5: swap keys and values. You can assume no two keys share a value.
// Invert(map[string]int{"one": 1, "two": 2}) == map[1:one 2:two]
func Invert(m map[string]int) map[int]string {
	// TODO
	return nil
}

// Ex6: return the keys of m in alphabetical order.
// SortedKeys(map[string]int{"pear": 1, "apple": 2}) == [apple pear]
//
// Hint: README section 3. You'll need to import "slices".
func SortedKeys(m map[string]int) []string {
	// TODO
	return nil
}

// Ex7: group words by their first letter.
// GroupByFirstLetter([]string{"ana", "bo", "amy"}) == map[a:[ana amy] b:[bo]]
// Skip empty strings (they have no first letter).
//
// Hint: README section 4, "Grouping into slices". word[:1] is the first letter.
func GroupByFirstLetter(words []string) map[string][]string {
	// TODO
	return nil
}

// Ex8: return the word that appears most often. If the slice is empty, return "".
// If there's a tie, return the one that comes first alphabetically (so the answer is predictable).
// MostCommon([]string{"a", "b", "b", "c"}) == "b"
//
// Hint: count first (you wrote WordCount in Ex2!). Then loop over the counts, keeping track of the
// best word and its count. A word beats the best if its count is higher, OR the count is equal and
// the word is alphabetically smaller (strings compare with <, like Python).
func MostCommon(words []string) string {
	// TODO
	return ""
}
