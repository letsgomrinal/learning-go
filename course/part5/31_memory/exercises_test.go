package lesson31

import (
	"math"
	"slices"
	"strings"
	"testing"
	"unsafe"
)

// Results are stored in these package-level variables so the compiler can't optimize the work away.
var (
	sinkInts []int
	sinkStr  string
)

func noRace(t *testing.T) {
	t.Helper()
	if raceEnabled {
		t.Skip("allocation counts aren't reliable under -race; run plain `go test`")
	}
}

func TestEx1EventSize(t *testing.T) {
	if got := unsafe.Sizeof(Event{}); got != 24 {
		t.Errorf("Ex1: unsafe.Sizeof(Event{}) = %d, want 24. Order the fields largest first", got)
	}
}

func TestEx2Squares(t *testing.T) {
	if got := Squares(5); !slices.Equal(got, []int{0, 1, 4, 9, 16}) {
		t.Fatalf("Ex2: Squares(5) = %v", got)
	}
	noRace(t)
	if a := testing.AllocsPerRun(100, func() { sinkInts = Squares(1000) }); a > 1 {
		t.Errorf("Ex2: Squares(1000) makes %v allocations, want 1. Preallocate the slice", a)
	}
}

func TestEx3JoinCSV(t *testing.T) {
	fields := strings.Fields("alpha beta gamma delta epsilon zeta eta theta")
	if got := JoinCSV(fields); got != "alpha,beta,gamma,delta,epsilon,zeta,eta,theta" {
		t.Fatalf("Ex3: JoinCSV = %q", got)
	}
	if got := JoinCSV(nil); got != "" {
		t.Errorf("Ex3: JoinCSV(nil) = %q, want \"\"", got)
	}
	if got := JoinCSV([]string{"solo"}); got != "solo" {
		t.Errorf("Ex3: JoinCSV([solo]) = %q", got)
	}
	noRace(t)
	if a := testing.AllocsPerRun(100, func() { sinkStr = JoinCSV(fields) }); a > 1 {
		t.Errorf("Ex3: JoinCSV makes %v allocations, want at most 1", a)
	}
}

func TestEx4Distance(t *testing.T) {
	if got := Distance(0, 0, 3, 4); math.Abs(got-5) > 1e-9 {
		t.Fatalf("Ex4: Distance(0,0,3,4) = %v, want 5", got)
	}
	noRace(t)
	if a := testing.AllocsPerRun(100, func() { Distance(1, 2, 3, 4) }); a != 0 {
		t.Errorf("Ex4: Distance makes %v allocations, want 0. Return Vec values, not pointers", a)
	}
}

func TestEx5Greet(t *testing.T) {
	if got := Greet("Ana"); got != "Hello, Ana! Welcome back." {
		t.Fatalf("Ex5: Greet(Ana) = %q", got)
	}
	if got := Greet("Bo"); got != "Hello, Bo! Welcome back." {
		t.Fatalf("Ex5: Greet(Bo) = %q (did you Reset the buffer?)", got)
	}
	noRace(t)
	long := strings.Repeat("x", 200)
	Greet(long) // warm the pool up with a big enough buffer
	if a := testing.AllocsPerRun(100, func() { sinkStr = Greet(long) }); a > 1 {
		t.Errorf("Ex5: Greet makes %v allocations, want at most 1 (just the result string). Use the pool", a)
	}
	if a := testing.AllocsPerRun(100, func() { sinkStr = Greet(long) }); a == 1 && !usesPool() {
		t.Errorf("Ex5: define a package-level sync.Pool and use it in Greet")
	}
}

func TestEx6EscapeAnswers(t *testing.T) {
	answers := map[string][2]string{
		"AnswerNewPointP":  {AnswerNewPointP, "heap"},
		"AnswerSumNums":    {AnswerSumNums, "stack"},
		"AnswerMakePointP": {AnswerMakePointP, "stack"},
	}
	for name, a := range answers {
		if strings.ToLower(strings.TrimSpace(a[0])) != a[1] {
			t.Errorf("Ex6: %s = %q. Run go build -gcflags=-m ./try and check again", name, a[0])
		}
	}
}
