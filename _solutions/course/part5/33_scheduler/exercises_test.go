package lesson33

import (
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// Always run these with:  go test -race -v

// settledGoroutines waits a little for goroutines to finish, then counts them.
func settledGoroutines(baseline int) int {
	n := runtime.NumGoroutine()
	for range 50 {
		if n <= baseline {
			return n
		}
		time.Sleep(10 * time.Millisecond)
		n = runtime.NumGoroutine()
	}
	return n
}

func within(t *testing.T, limit time.Duration, what string, f func()) {
	t.Helper()
	done := make(chan struct{})
	go func() { f(); close(done) }()
	select {
	case <-done:
	case <-time.After(limit):
		t.Fatalf("%s did not finish within %v (deadlock or leak?)", what, limit)
	}
}

func TestEx1FirstMatch(t *testing.T) {
	items := []string{"apple", "banana", "cherry", "avocado", "apricot", "blueberry"}
	startsWithA := func(s string) bool { return strings.HasPrefix(s, "a") }

	before := runtime.NumGoroutine()
	var got string
	within(t, 2*time.Second, "Ex1: FirstMatch", func() { got = FirstMatch(items, startsWithA) })
	if !startsWithA(got) {
		t.Errorf("Ex1: FirstMatch = %q, want one of the items starting with a", got)
	}
	if after := settledGoroutines(before); after > before {
		t.Errorf("Ex1: %d goroutines leaked (still running after FirstMatch returned)", after-before)
	}

	within(t, 2*time.Second, "Ex1: FirstMatch with no match", func() {
		got = FirstMatch(items, func(string) bool { return false })
	})
	if got != "" {
		t.Errorf("Ex1: no match: got %q, want \"\"", got)
	}
}

func TestEx2StartHeartbeat(t *testing.T) {
	before := runtime.NumGoroutine()
	var beats atomic.Int32
	stop := StartHeartbeat(5*time.Millisecond, func() { beats.Add(1) })
	time.Sleep(60 * time.Millisecond)
	within(t, time.Second, "Ex2: stop()", stop)

	n := beats.Load()
	if n < 3 {
		t.Fatalf("Ex2: only %d beats in 60ms with a 5ms interval", n)
	}
	time.Sleep(30 * time.Millisecond)
	if later := beats.Load(); later != n {
		t.Errorf("Ex2: beat was called %d more times after stop() returned", later-n)
	}
	if after := settledGoroutines(before); after > before {
		t.Errorf("Ex2: the heartbeat goroutine is still running after stop()")
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Ex2: calling stop() twice panicked: %v", r)
			}
		}()
		stop()
	}()
}

func TestEx3Transfer(t *testing.T) {
	a := &Account{ID: 1, Balance: 1000}
	b := &Account{ID: 2, Balance: 1000}
	within(t, 5*time.Second, "Ex3: 200 transfers in both directions", func() {
		var wg sync.WaitGroup
		for range 100 {
			wg.Add(2)
			go func() { defer wg.Done(); Transfer(a, b, 3) }()
			go func() { defer wg.Done(); Transfer(b, a, 1) }()
		}
		wg.Wait()
	})
	if a.Balance != 800 || b.Balance != 1200 {
		t.Errorf("Ex3: balances %d and %d, want 800 and 1200", a.Balance, b.Balance)
	}
}

func TestEx4CountPrimes(t *testing.T) {
	for limit, want := range map[int]int{0: 0, 2: 0, 3: 1, 10: 4, 100: 25, 1000: 168} {
		if got := CountPrimes(limit); got != want {
			t.Errorf("Ex4: CountPrimes(%d) = %d, want %d", limit, got, want)
		}
	}
	if runtime.GOMAXPROCS(0) < 2 {
		t.Skip("only one CPU: can't check the speed-up")
	}
	sequential := func() {
		n := 0
		for i := range 3_000_000 {
			if IsPrime(i) {
				n++
			}
		}
	}
	measure := func(f func()) time.Duration {
		best := time.Duration(1 << 62)
		for range 2 {
			start := time.Now()
			f()
			best = min(best, time.Since(start))
		}
		return best
	}
	seq := measure(sequential)
	par := measure(func() { CountPrimes(3_000_000) })
	t.Logf("sequential %v, yours %v", seq.Round(time.Millisecond), par.Round(time.Millisecond))
	if float64(seq)/float64(par) < 1.5 {
		t.Errorf("Ex4: only %.1fx faster than one goroutine on %d cores. Split the work between GOMAXPROCS workers",
			float64(seq)/float64(par), runtime.GOMAXPROCS(0))
	}
}
