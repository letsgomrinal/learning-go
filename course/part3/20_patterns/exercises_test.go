package lesson20

import (
	"slices"
	"sync/atomic"
	"testing"
	"time"
)

// Always run these with:  go test -race -v

func drain(t *testing.T, ch <-chan int) []int {
	t.Helper()
	if ch == nil {
		t.Fatal("got a nil channel (not implemented yet?)")
	}
	var out []int
	deadline := time.After(2 * time.Second)
	for {
		select {
		case v, ok := <-ch:
			if !ok {
				return out
			}
			out = append(out, v)
		case <-deadline:
			t.Fatalf("the channel was never closed (received %v so far)", out)
		}
	}
}

// chanOf builds a closed channel holding vals, so stages can be tested on their own.
func chanOf(vals ...int) <-chan int {
	ch := make(chan int, len(vals))
	for _, v := range vals {
		ch <- v
	}
	close(ch)
	return ch
}

func within(t *testing.T, limit time.Duration, f func()) {
	t.Helper()
	done := make(chan struct{})
	go func() { f(); close(done) }()
	select {
	case <-done:
	case <-time.After(limit):
		t.Fatalf("did not return within %v (deadlock?)", limit)
	}
}

func TestEx1Generate(t *testing.T) {
	if got := drain(t, Generate(3, 1, 2)); !slices.Equal(got, []int{3, 1, 2}) {
		t.Errorf("Ex1: Generate(3, 1, 2) sent %v, want [3 1 2]", got)
	}
	if got := drain(t, Generate()); len(got) != 0 {
		t.Errorf("Ex1: Generate() sent %v, want nothing", got)
	}
}

func TestEx2Square(t *testing.T) {
	if got := drain(t, Square(chanOf(1, 2, 3))); !slices.Equal(got, []int{1, 4, 9}) {
		t.Errorf("Ex2: Square(1, 2, 3) sent %v, want [1 4 9]", got)
	}
	if got := drain(t, Square(Square(chanOf(2, 3)))); !slices.Equal(got, []int{16, 81}) {
		t.Errorf("Ex2: Square(Square(2, 3)) sent %v, want [16 81]", got)
	}
}

func TestEx3SumSquares(t *testing.T) {
	var got int
	within(t, 2*time.Second, func() { got = SumSquares(1, 2, 3, 4) })
	if got != 30 {
		t.Errorf("Ex3: SumSquares(1, 2, 3, 4) = %d, want 30", got)
	}
}

func TestEx4Merge(t *testing.T) {
	got := drain(t, Merge(chanOf(1, 2, 3), chanOf(10, 20), chanOf()))
	slices.Sort(got)
	if want := []int{1, 2, 3, 10, 20}; !slices.Equal(got, want) {
		t.Errorf("Ex4: Merge sent %v (sorted), want %v", got, want)
	}
	if got := drain(t, Merge()); len(got) != 0 {
		t.Errorf("Ex4: Merge() with no inputs sent %v, want nothing (just close)", got)
	}
}

func TestEx5WorkerPool(t *testing.T) {
	var active, peak atomic.Int32
	work := func(n int) int {
		now := active.Add(1)
		for {
			p := peak.Load()
			if now <= p || peak.CompareAndSwap(p, now) {
				break
			}
		}
		time.Sleep(time.Duration(10+(n%3)*10) * time.Millisecond) // jobs finish out of order
		active.Add(-1)
		return n * 2
	}
	jobs := []int{1, 2, 3, 4, 5, 6, 7, 8}

	var got []int
	start := time.Now()
	within(t, 3*time.Second, func() { got = WorkerPool(jobs, 4, work) })
	elapsed := time.Since(start)

	if want := []int{2, 4, 6, 8, 10, 12, 14, 16}; !slices.Equal(got, want) {
		t.Fatalf("Ex5: WorkerPool = %v, want %v (same order as the jobs)", got, want)
	}
	if p := peak.Load(); p > 4 {
		t.Errorf("Ex5: %d jobs ran at once, want at most 4 workers", p)
	}
	if p := peak.Load(); p < 2 {
		t.Errorf("Ex5: jobs never ran at the same time")
	}
	if elapsed > 150*time.Millisecond {
		t.Errorf("Ex5: took %v; with 4 workers it should take well under 150ms", elapsed.Round(time.Millisecond))
	}
	if got := WorkerPool(nil, 3, work); len(got) != 0 {
		t.Errorf("Ex5: WorkerPool(no jobs) = %v, want empty", got)
	}
}

func TestEx6RunLimited(t *testing.T) {
	var active, peak, ran atomic.Int32
	tasks := make([]func(), 12)
	for i := range tasks {
		tasks[i] = func() {
			now := active.Add(1)
			for {
				p := peak.Load()
				if now <= p || peak.CompareAndSwap(p, now) {
					break
				}
			}
			time.Sleep(20 * time.Millisecond)
			active.Add(-1)
			ran.Add(1)
		}
	}
	within(t, 3*time.Second, func() { RunLimited(tasks, 3) })
	if ran.Load() != 12 {
		t.Fatalf("Ex6: %d of 12 tasks had finished when RunLimited returned", ran.Load())
	}
	if p := peak.Load(); p > 3 {
		t.Errorf("Ex6: %d tasks ran at once, want at most 3", p)
	}
	if p := peak.Load(); p < 2 {
		t.Errorf("Ex6: tasks never ran at the same time (peak %d). Start them as goroutines", p)
	}
}
