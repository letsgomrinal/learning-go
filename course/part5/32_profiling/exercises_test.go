package lesson32

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"testing"
	"time"
)

// ── The original slow versions, kept here so the tests can compare. ──

func uniqueWordsSlow(text string) int {
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

func hasCommonSlow(a, b []int) bool {
	for _, x := range a {
		for _, y := range b {
			if x == y {
				return true
			}
		}
	}
	return false
}

func validEmailsSlow(addresses []string) int {
	n := 0
	for _, a := range addresses {
		re := regexp.MustCompile(`^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`)
		if re.MatchString(a) {
			n++
		}
	}
	return n
}

func buildReportSlow(items []string) string {
	report := ""
	for i, item := range items {
		report += fmt.Sprintf("%d. %s\n", i+1, item)
	}
	return report
}

func lookupRegionSlow(city string) string {
	for _, cr := range Regions {
		if cr.City == city {
			return cr.Region
		}
	}
	return "unknown"
}

func summarizeSlow(records []Record) Summary {
	out := Summary{}
	for _, r := range records {
		out[lookupRegionSlow(normalize(r.City))] += r.Cents
	}
	return out
}

// ── Helpers ──

// fasterBy runs both functions a few times and checks that fast is at least `factor` times quicker.
func fasterBy(t *testing.T, name string, factor float64, slow, fast func()) {
	t.Helper()
	if raceEnabled {
		t.Skip("timings aren't meaningful under -race; run plain `go test`")
	}
	measure := func(f func()) time.Duration {
		f() // warm up
		best := time.Duration(1 << 62)
		for range 3 {
			start := time.Now()
			f()
			best = min(best, time.Since(start))
		}
		return best
	}
	s, f := measure(slow), measure(fast)
	ratio := float64(s) / float64(max(f, time.Nanosecond))
	t.Logf("%s: original %v, yours %v (%.1fx faster)", name, s.Round(time.Microsecond), f.Round(time.Microsecond), ratio)
	if ratio < factor {
		t.Errorf("%s is only %.1fx faster than the original, want at least %.0fx", name, ratio, factor)
	}
}

func words(n, distinct int) string {
	r := rand.New(rand.NewSource(1))
	ws := make([]string, n)
	for i := range ws {
		ws[i] = fmt.Sprintf("w%d", r.Intn(distinct))
	}
	return strings.Join(ws, " ")
}

// ── Tests ──

func TestEx1UniqueWords(t *testing.T) {
	for _, text := range []string{"", "a", "go is go", "  the cat  and the hat  ", words(500, 100)} {
		if got, want := UniqueWords(text), uniqueWordsSlow(text); got != want {
			t.Fatalf("Ex1: UniqueWords gives %d, the original gives %d (for %d characters of text)", got, want, len(text))
		}
	}
	big := words(20000, 15000)
	fasterBy(t, "Ex1: UniqueWords", 10, func() { uniqueWordsSlow(big) }, func() { UniqueWords(big) })
}

func TestEx2HasCommon(t *testing.T) {
	cases := [][2][]int{{{1, 2, 3}, {4, 5, 3}}, {{1, 2}, {3, 4}}, {nil, {1}}, {{7}, {7}}}
	for _, c := range cases {
		if got, want := HasCommon(c[0], c[1]), hasCommonSlow(c[0], c[1]); got != want {
			t.Fatalf("Ex2: HasCommon(%v, %v) = %v, want %v", c[0], c[1], got, want)
		}
	}
	a := make([]int, 10000)
	b := make([]int, 10000)
	for i := range a {
		a[i] = i * 2   // even
		b[i] = i*2 + 1 // odd: nothing in common, the slowest case
	}
	fasterBy(t, "Ex2: HasCommon", 10, func() { hasCommonSlow(a, b) }, func() { HasCommon(a, b) })
}

func TestEx3ValidEmails(t *testing.T) {
	sample := []string{"ana@example.com", "bad@", "@bad.com", "bo.smith+go@mail.co", "no-at-sign", "x@y.z"}
	if got, want := ValidEmails(sample), validEmailsSlow(sample); got != want {
		t.Fatalf("Ex3: ValidEmails = %d, want %d", got, want)
	}
	many := make([]string, 3000)
	for i := range many {
		many[i] = sample[i%len(sample)]
	}
	fasterBy(t, "Ex3: ValidEmails", 5, func() { validEmailsSlow(many) }, func() { ValidEmails(many) })
}

func TestEx4BuildReport(t *testing.T) {
	items := []string{"apples", "pears", "plums"}
	if got, want := BuildReport(items), buildReportSlow(items); got != want {
		t.Fatalf("Ex4: BuildReport = %q, want %q", got, want)
	}
	if got := BuildReport(nil); got != "" {
		t.Errorf("Ex4: BuildReport(nil) = %q, want \"\"", got)
	}
	many := make([]string, 10000)
	for i := range many {
		many[i] = "an item with a reasonably long description"
	}
	fasterBy(t, "Ex4: BuildReport", 5, func() { buildReportSlow(many) }, func() { BuildReport(many) })
}

func sampleRecords(n int) []Record {
	r := rand.New(rand.NewSource(2))
	recs := make([]Record, n)
	for i := range recs {
		city := fmt.Sprintf("City-%d ", r.Intn(1100)) // some cities are unknown
		recs[i] = Record{Customer: "c", City: city, Cents: r.Intn(10000)}
	}
	return recs
}

func TestEx5Summarize(t *testing.T) {
	if HotFunction != "lookupRegion" {
		t.Errorf("Ex5: HotFunction = %q. Profile it: go test -bench=Summarize -cpuprofile=cpu.out && go tool pprof -top cpu.out", HotFunction)
	}
	recs := sampleRecords(20000)
	got, want := Summarize(recs), summarizeSlow(recs)
	if len(got) != len(want) {
		t.Fatalf("Ex5: Summarize has %d regions, want %d", len(got), len(want))
	}
	for k, v := range want {
		if got[k] != v {
			t.Fatalf("Ex5: Summarize[%q] = %d, want %d", k, got[k], v)
		}
	}
	fasterBy(t, "Ex5: Summarize", 5, func() { summarizeSlow(recs) }, func() { Summarize(recs) })
}

// BenchmarkSummarize is what you profile for Ex5.
func BenchmarkSummarize(b *testing.B) {
	recs := sampleRecords(20000)
	for b.Loop() {
		Summarize(recs)
	}
}
