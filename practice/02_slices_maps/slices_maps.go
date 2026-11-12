// Lesson 02: Slices and maps. Where Python devs get bitten first.
//
// Run:  go test -v ./practice/02_slices_maps
//
// ─── CONCEPTS ──────────────────────────────────────────────────────────────
//
// 1. ARRAY vs SLICE
//      var a [3]int        // ARRAY: fixed size, part of the type. Rarely used directly.
//      s := []int{1, 2, 3} // SLICE: what you use 99% of the time.
//
// 2. A SLICE IS A SMALL STRUCT (the "slice header"):
//      type slice struct { ptr *T; len int; cap int }
//    It POINTS to a backing array. Copying a slice copies the header, NOT the data:
//      a := []int{1, 2, 3}
//      b := a          // b shares a's backing array
//      b[0] = 99       // a[0] is now 99 too (same as Python lists, so far so good)
//
// 3. SLICING SHARES MEMORY (unlike Python, where lst[1:3] makes a COPY!)
//      s := []int{0, 1, 2, 3, 4}
//      t := s[1:3]     // [1 2], len 2, cap 4. Shares memory with s!
//      t[0] = 100      // s is now [0 100 2 3 4]
//
// 4. APPEND MAY OR MAY NOT ALLOCATE. This is the #1 Go slice trap.
//      s = append(s, x)   // ALWAYS reassign the result
//    If len < cap, append writes into the EXISTING backing array (visible to anyone
//    else sharing it). If len == cap, it allocates a new, bigger array and copies.
//    So two appends to the same base slice can overwrite each other.
//    Read: https://go.dev/blog/slices-intro
//
// 5. MAKING AND COPYING
//      s := make([]int, 5)        // len 5, cap 5, all zeros
//      s := make([]int, 0, 100)   // len 0, cap 100 (preallocate: fast appends)
//      n := copy(dst, src)        // copies min(len(dst), len(src)) elements
//      c := slices.Clone(s)       // Go 1.21+: independent copy
//
// 6. MAPS
//      m := map[string]int{}      // or make(map[string]int)
//      var m map[string]int       // nil map: reads OK, WRITES PANIC!
//      v, ok := m["k"]            // ok is false if missing (v is the zero value)
//      delete(m, "k")
//      Iteration order is deliberately RANDOM. If you need order, sort the keys.
//    Sets don't exist. Use map[T]struct{} (struct{} takes zero bytes) or map[T]bool.
//
// 7. USEFUL PACKAGES (Go 1.21+): slices (Sort, Contains, Index, Equal, Clone, Reverse...)
//    and maps (Keys, Values, Clone, Equal). maps.Keys returns an ITERATOR (lesson 06);
//    slices.Collect turns an iterator into a slice. slices.Sorted(maps.Keys(m)) is handy.
// ───────────────────────────────────────────────────────────────────────────

package slicesmaps

// RemoveAt returns a NEW slice with the element at index i removed.
// The input slice must NOT be modified. Callers rely on it staying intact.
// If i is out of range, return a copy of s unchanged.
//
// Trap: `append(s[:i], s[i+1:]...)` looks right but OVERWRITES s's backing array.
// Build into a fresh slice instead: make([]int, 0, len(s)-1).
func RemoveAt(s []int, i int) []int {
	// TODO
	return nil
}

// SafeAppend returns base with v appended, guaranteeing the result does NOT share
// memory with base. Calling SafeAppend(base, 1) and SafeAppend(base, 2) must give two
// independent slices, even when base has spare capacity.
//
// Read the test TestSafeAppend to see the bug that a plain append(base, v) would cause.
func SafeAppend(base []int, v int) []int {
	// TODO
	return nil
}

// Chunk splits s into consecutive chunks of at most `size` elements.
// Chunk([1 2 3 4 5], 2) == [[1 2] [3 4] [5]]. If size <= 0, return nil.
//
// Hint: a slice of slices is [][]int. s[i:min(i+size, len(s))] is a chunk view.
func Chunk(s []int, size int) [][]int {
	// TODO
	return nil
}

// Dedupe returns the unique strings of s in order of first appearance.
// Dedupe(["b","a","b","c","a"]) == ["b","a","c"]
//
// Hint: seen := make(map[string]struct{}); seen[x] = struct{}{}; _, ok := seen[x]
func Dedupe(s []string) []string {
	// TODO
	return nil
}

// GroupByLen groups words by their length in bytes.
// GroupByLen(["go","is","fun","yes"]) == {2: ["go","is"], 3: ["fun","yes"]}
//
// Hint: appending to a nil slice works, so m[k] = append(m[k], w) is all you need.
func GroupByLen(words []string) map[int][]string {
	// TODO
	return nil
}

// SortedKeys returns the map's keys sorted ascending.
// (Because map iteration order is random, this is a very common helper.)
func SortedKeys(m map[string]int) []string {
	// TODO
	return nil
}

// Matrix creates a rows×cols grid filled with zeros.
// Each row must be independent: changing grid[0][0] must not change grid[1][0].
// (Python has the same trap with [[0]*cols]*rows!)
func Matrix(rows, cols int) [][]int {
	// TODO
	return nil
}
