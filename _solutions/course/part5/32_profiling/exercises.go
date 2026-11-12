// Reference solution for lesson 32.
package lesson32

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func UniqueWords(text string) int {
	seen := map[string]struct{}{}
	for _, w := range strings.Fields(text) {
		seen[w] = struct{}{}
	}
	return len(seen)
}

func HasCommon(a, b []int) bool {
	set := make(map[int]struct{}, len(a))
	for _, x := range a {
		set[x] = struct{}{}
	}
	for _, y := range b {
		if _, ok := set[y]; ok {
			return true
		}
	}
	return false
}

// Compiled once, when the program starts.
var emailRe = regexp.MustCompile(`^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`)

func ValidEmails(addresses []string) int {
	n := 0
	for _, a := range addresses {
		if emailRe.MatchString(a) {
			n++
		}
	}
	return n
}

func BuildReport(items []string) string {
	var b strings.Builder
	for i, item := range items {
		b.WriteString(strconv.Itoa(i + 1))
		b.WriteString(". ")
		b.WriteString(item)
		b.WriteByte('\n')
	}
	return b.String()
}

type Record struct {
	Customer string
	City     string
	Cents    int
}

type Summary map[string]int

var Regions []CityRegion

type CityRegion struct {
	City, Region string
}

// regionByCity is built once from Regions, so lookups are O(1) instead of a linear search.
var regionByCity map[string]string

func init() {
	for i := range 1000 {
		Regions = append(Regions, CityRegion{City: fmt.Sprintf("city-%d", i), Region: fmt.Sprintf("region-%d", i%7)})
	}
	regionByCity = make(map[string]string, len(Regions))
	for _, cr := range Regions {
		if _, dup := regionByCity[cr.City]; !dup { // keep the first, like the linear search did
			regionByCity[cr.City] = cr.Region
		}
	}
}

func Summarize(records []Record) Summary {
	out := Summary{}
	for _, r := range records {
		out[lookupRegion(normalize(r.City))] += r.Cents
	}
	return out
}

func normalize(city string) string {
	return strings.ToLower(strings.TrimSpace(city))
}

func lookupRegion(city string) string {
	if region, ok := regionByCity[city]; ok {
		return region
	}
	return "unknown"
}

const HotFunction = "lookupRegion"
