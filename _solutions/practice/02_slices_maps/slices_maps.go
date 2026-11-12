// Reference solution for lesson 02.
package slicesmaps

import (
	"maps"
	"slices"
)

func RemoveAt(s []int, i int) []int {
	if i < 0 || i >= len(s) {
		return slices.Clone(s)
	}
	out := make([]int, 0, len(s)-1)
	out = append(out, s[:i]...)
	out = append(out, s[i+1:]...)
	return out
}

func SafeAppend(base []int, v int) []int {
	out := make([]int, len(base), len(base)+1)
	copy(out, base)
	return append(out, v)
	// Alternative one-liner: append(base[:len(base):len(base)], v)
	// The "full slice expression" s[low:high:max] caps the capacity, forcing append to allocate.
}

func Chunk(s []int, size int) [][]int {
	if size <= 0 {
		return nil
	}
	var out [][]int
	for i := 0; i < len(s); i += size {
		out = append(out, s[i:min(i+size, len(s))])
	}
	return out
}

func Dedupe(s []string) []string {
	seen := make(map[string]struct{}, len(s))
	var out []string
	for _, x := range s {
		if _, ok := seen[x]; ok {
			continue
		}
		seen[x] = struct{}{}
		out = append(out, x)
	}
	return out
}

func GroupByLen(words []string) map[int][]string {
	groups := make(map[int][]string)
	for _, w := range words {
		groups[len(w)] = append(groups[len(w)], w)
	}
	return groups
}

func SortedKeys(m map[string]int) []string {
	return slices.Sorted(maps.Keys(m))
}

func Matrix(rows, cols int) [][]int {
	grid := make([][]int, rows)
	for i := range grid {
		grid[i] = make([]int, cols)
	}
	return grid
}
