// Package mathx has small math helpers. (Reference solution for lesson 09, Ex4.)
package mathx

// Clamp returns v limited to the range lo..hi.
func Clamp(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
