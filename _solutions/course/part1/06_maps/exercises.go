// Reference solution for lesson 06.
package lesson06

import "slices"

func Lookup(prices map[string]float64, item string) (float64, bool) {
	price, ok := prices[item]
	return price, ok
}

func WordCount(words []string) map[string]int {
	counts := map[string]int{}
	for _, w := range words {
		counts[w]++
	}
	return counts
}

func HasDuplicates(nums []int) bool {
	seen := map[int]bool{}
	for _, n := range nums {
		if seen[n] {
			return true
		}
		seen[n] = true
	}
	return false
}

func Unique(words []string) []string {
	seen := map[string]bool{}
	var result []string
	for _, w := range words {
		if !seen[w] {
			seen[w] = true
			result = append(result, w)
		}
	}
	return result
}

func Invert(m map[string]int) map[int]string {
	result := map[int]string{}
	for k, v := range m {
		result[v] = k
	}
	return result
}

func SortedKeys(m map[string]int) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	return keys
}

func GroupByFirstLetter(words []string) map[string][]string {
	groups := map[string][]string{}
	for _, w := range words {
		if w == "" {
			continue
		}
		first := w[:1]
		groups[first] = append(groups[first], w)
	}
	return groups
}

func MostCommon(words []string) string {
	best, bestCount := "", 0
	for w, c := range WordCount(words) {
		if c > bestCount || (c == bestCount && w < best) {
			best, bestCount = w, c
		}
	}
	return best
}
