package concurrency

import (
	"slices"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

// collect drains ch, failing the test if ch is nil or never gets closed.
func collect(t *testing.T, ch <-chan int) []int {
	t.Helper()
	if ch == nil {
		t.Fatal("got a nil channel (not implemented yet?)")
	}
	var out []int
	timeout := time.After(2 * time.Second)
	for {
		select {
		case v, ok := <-ch:
			if !ok {
				return out
			}
			out = append(out, v)
		case <-timeout:
			t.Fatalf("channel was never closed (got %v so far)", out)
		}
	}
}

func TestParallelSum(t *testing.T) {
	nums := make([]int, 1000)
	for i := range nums {
		nums[i] = i + 1
	}
	for _, w := range []int{1, 3, 8, 2000, 0} {
		if got := ParallelSum(nums, w); got != 500500 {
			t.Errorf("ParallelSum(1..1000, workers=%d) = %d, want 500500", w, got)
		}
	}
	if got := ParallelSum(nil, 4); got != 0 {
		t.Errorf("ParallelSum(nil) = %d", got)
	}
}

func TestPipeline(t *testing.T) {
	got := collect(t, Generate(1, 2, 3, 4))
	if want := []int{1, 2, 3, 4}; !slices.Equal(got, want) {
		t.Errorf("Generate = %v, want %v", got, want)
	}

	sq := Square(Generate(1, 2, 3, 4))
	got = collect(t, sq)
	if want := []int{1, 4, 9, 16}; !slices.Equal(got, want) {
		t.Errorf("Square(Generate) = %v, want %v", got, want)
	}

	// stages compose
	got = collect(t, Square(Square(Generate(2, 3))))
	if want := []int{16, 81}; !slices.Equal(got, want) {
		t.Errorf("Square(Square) = %v, want %v", got, want)
	}
}

func TestMerge(t *testing.T) {
	got := collect(t, Merge(Generate(1, 2, 3), Generate(10, 20), Generate()))
	slices.Sort(got) // merged order is nondeterministic
	if want := []int{1, 2, 3, 10, 20}; !slices.Equal(got, want) {
		t.Errorf("Merge = %v, want %v", got, want)
	}
}

func TestProcessOrdered(t *testing.T) {
	inputs := []string{"a", "b", "c", "d", "e", "f", "g", "h"}
	var active, maxActive atomic.Int32

	fn := func(s string) string {
		n := active.Add(1)
		for { // record the peak number of concurrent calls
			m := maxActive.Load()
			if n <= m || maxActive.CompareAndSwap(m, n) {
				break
			}
		}
		time.Sleep(20 * time.Millisecond)
		active.Add(-1)
		return strings.ToUpper(s)
	}

	start := time.Now()
	got := ProcessOrdered(inputs, 4, fn)
	elapsed := time.Since(start)

	if want := []string{"A", "B", "C", "D", "E", "F", "G", "H"}; !slices.Equal(got, want) {
		t.Fatalf("ProcessOrdered = %v, want %v (order must match input)", got, want)
	}
	if m := maxActive.Load(); m > 4 {
		t.Errorf("max concurrent calls = %d, want <= 4 workers", m)
	}
	if m := maxActive.Load(); m < 2 {
		t.Errorf("max concurrent calls = %d: work isn't running in parallel", m)
	}
	if elapsed > 120*time.Millisecond {
		t.Errorf("took %v; 8 jobs × 20ms on 4 workers should take ~40ms", elapsed)
	}
}

func TestWithTimeout(t *testing.T) {
	v, ok := WithTimeout(func() int { return 42 }, 100*time.Millisecond)
	if !ok || v != 42 {
		t.Errorf("fast fn: got (%d, %v), want (42, true)", v, ok)
	}

	start := time.Now()
	v, ok = WithTimeout(func() int { time.Sleep(300 * time.Millisecond); return 1 }, 20*time.Millisecond)
	if ok || v != 0 {
		t.Errorf("slow fn: got (%d, %v), want (0, false)", v, ok)
	}
	if time.Since(start) > 150*time.Millisecond {
		t.Errorf("WithTimeout waited for the slow function instead of timing out")
	}
}
