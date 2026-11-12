// Reference solution for lesson 31. Compare with yours AFTER your tests pass.
package lesson31

import (
	"bytes"
	"math"
	"strings"
	"sync"
)

// Largest fields first: 8 + 8 + 4 + 1 + 1 + 1 (+1 padding) = 24 bytes.
type Event struct {
	Timestamp int64
	Code      int64
	Source    int32
	Urgent    bool
	Retried   bool
	Done      bool
}

func Squares(n int) []int {
	out := make([]int, 0, n)
	for i := range n {
		out = append(out, i*i)
	}
	return out
}

func JoinCSV(fields []string) string {
	if len(fields) == 0 {
		return ""
	}
	total := len(fields) - 1 // the commas
	for _, f := range fields {
		total += len(f)
	}
	var b strings.Builder
	b.Grow(total)
	for i, f := range fields {
		if i > 0 {
			b.WriteByte(',')
		}
		b.WriteString(f)
	}
	return b.String()
}

type Vec struct {
	X, Y float64
}

func NewVec(x, y float64) Vec {
	return Vec{X: x, Y: y}
}

func Distance(x1, y1, x2, y2 float64) float64 {
	a := NewVec(x1, y1)
	b := NewVec(x2, y2)
	return math.Hypot(b.X-a.X, b.Y-a.Y)
}

var bufPool = sync.Pool{
	New: func() any { return new(bytes.Buffer) },
}

func Greet(name string) string {
	buf := bufPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer bufPool.Put(buf)

	buf.WriteString("Hello, ")
	buf.WriteString(name)
	buf.WriteString("! Welcome back.")
	return buf.String()
}

const (
	AnswerNewPointP  = "heap"  // "moved to heap: p"
	AnswerSumNums    = "stack" // "nums does not escape"
	AnswerMakePointP = "stack" // returned by value: copied, no heap allocation
)
