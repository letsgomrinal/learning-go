// Lesson 31 exercises: memory and allocations. Read README.md first!
//
// Check one exercise at a time:   go test -run Ex1 -v
// (Use plain go test here, not -race: the race detector changes allocation counts.)
package lesson31

import (
	"math"
	"strings"
)

// Ex1: Event is 40 bytes because of padding. Reorder its fields (don't add, remove or rename
// any) so unsafe.Sizeof(Event{}) is 24. (README section 5: largest first.)
type Event struct {
	Urgent    bool
	Timestamp int64
	Retried   bool
	Source    int32
	Code      int64
	Done      bool
}

// Ex2: return the squares 0², 1², ..., (n-1)², with exactly ONE allocation.
// The current version grows the slice again and again. Preallocate it.
func Squares(n int) []int {
	var out []int
	for i := range n {
		out = append(out, i*i)
	}
	return out
}

// Ex3: join the fields with commas: JoinCSV([]string{"a", "b", "c"}) == "a,b,c"
// With AT MOST one allocation. The current version allocates for every field.
//
// Hint: work out the final length first (the fields' lengths plus the commas),
// then var b strings.Builder; b.Grow(total); write the pieces; return b.String()
func JoinCSV(fields []string) string {
	s := ""
	for i, f := range fields {
		if i > 0 {
			s += ","
		}
		s += f
	}
	return s
}

var _ = strings.Builder{}

// Ex4: Distance returns the distance between two points. It must not allocate at all.
// Right now NewVec returns a pointer, which forces heap allocations (README section 4c).
// Change NewVec to return a Vec value (not a pointer), and update Distance to match.
type Vec struct {
	X, Y float64
}

func NewVec(x, y float64) *Vec {
	return &Vec{X: x, Y: y}
}

// sink stops the compiler from optimising the pointers away, so you can see the effect.
var sink *Vec

func Distance(x1, y1, x2, y2 float64) float64 {
	a := NewVec(x1, y1)
	b := NewVec(x2, y2)
	sink = a
	return math.Hypot(b.X-a.X, b.Y-a.Y)
}

// Ex5: Greet returns "Hello, <name>! Welcome back." building it in a bytes.Buffer that comes from a
// sync.Pool, so the buffer is reused between calls. Only the final string may allocate (1 allocation).
//
// Steps (README section 4d):
//  1. declare a package-level  var bufPool = sync.Pool{New: func() any { return new(bytes.Buffer) }}
//  2. in Greet: buf := bufPool.Get().(*bytes.Buffer); buf.Reset(); defer bufPool.Put(buf)
//  3. write the pieces with buf.WriteString, return buf.String()
func Greet(name string) string {
	// TODO
	return "Hello, " + name + "! Welcome back."
}

// Ex6: run  go build -gcflags=-m ./try  on this lesson's try/main.go after copying the README's
// section 2 example into it (newPoint and sum). Then fill in the answers below with what the
// compiler says: "heap" or "stack".
const (
	AnswerNewPointP  = "" // where does p in newPoint go?
	AnswerSumNums    = "" // where does the nums slice passed to sum stay?
	AnswerMakePointP = "" // and Point{1, 2} returned BY VALUE from makePoint (section 4c)?
)
