// Lesson 32 exercises: profiling and optimizing. Read README.md first!
//
// Every function here is CORRECT but SLOW. Make each one fast without changing its results.
// Check one exercise at a time (without -race):   go test -run Ex1 -v
package lesson32

import (
	"fmt"
	"regexp"
	"strings"
)

// Ex1: return how many distinct words (split on whitespace) the text has.
// It's slow because `seen` is a slice: checking it is O(n) for every word → O(n²) overall.
// Use a map (or set) instead.
func UniqueWords(text string) int {
	var seen []string
	for _, w := range strings.Fields(text) {
		found := false
		for _, s := range seen {
			if s == w {
				found = true
				break
			}
		}
		if !found {
			seen = append(seen, w)
		}
	}
	return len(seen)
}

// Ex2: report whether a and b have any number in common.
// Nested loops: len(a) × len(b) comparisons. Put one side in a set first.
func HasCommon(a, b []int) bool {
	for _, x := range a {
		for _, y := range b {
			if x == y {
				return true
			}
		}
	}
	return false
}

// Ex3: count how many of the addresses look valid. Something here is done again and again,
// even though its result never changes. Move it out so it happens only once.
func ValidEmails(addresses []string) int {
	n := 0
	for _, a := range addresses {
		re := regexp.MustCompile(`^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`)
		if re.MatchString(a) {
			n++
		}
	}
	return n
}

// Ex4: build a report with one numbered line per item: "1. apples\n2. pears\n..."
// String concatenation in a loop copies everything every time. Use a strings.Builder,
// and write the number with strconv.Itoa instead of fmt.Sprintf.
func BuildReport(items []string) string {
	report := ""
	for i, item := range items {
		report += fmt.Sprintf("%d. %s\n", i+1, item)
	}
	return report
}

// ── Ex5: use the profiler ─────────────────────────────────────────────────

// Record is one sale.
type Record struct {
	Customer string
	City     string
	Cents    int
}

// Summary is the total per region.
type Summary map[string]int

// Regions maps each city to its sales region. (A thousand cities, filled in by init below.)
var Regions []CityRegion

type CityRegion struct {
	City, Region string
}

func init() {
	for i := range 1000 {
		Regions = append(Regions, CityRegion{City: fmt.Sprintf("city-%d", i), Region: fmt.Sprintf("region-%d", i%7)})
	}
}

// Summarize adds up the sales per region. It's slow, and you need to find out WHY with pprof
// (README section 8). Don't guess: profile, then fix only the hot function.
func Summarize(records []Record) Summary {
	out := Summary{}
	for _, r := range records {
		name := normalize(r.City)
		region := lookupRegion(name)
		out[region] += r.Cents
	}
	return out
}

func normalize(city string) string {
	return strings.ToLower(strings.TrimSpace(city))
}

func lookupRegion(city string) string {
	for _, cr := range Regions {
		if cr.City == city {
			return cr.Region
		}
	}
	return "unknown"
}

// HotFunction: after profiling Summarize, write the name of YOUR function that uses most of the time.
// Ignore runtime.* entries (Go internals your code calls): look at flat and cum for the functions in
// this file. Just the function name, for example "normalize".
const HotFunction = ""
