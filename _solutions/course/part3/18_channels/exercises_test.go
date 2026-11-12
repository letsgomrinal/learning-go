package lesson18

import (
	"runtime"
	"slices"
	"testing"
	"time"
)

// Always run these with:  go test -race -v

// receive gets one value, failing instead of hanging if nothing arrives.
func receive[T any](t *testing.T, ch <-chan T) (T, bool) {
	t.Helper()
	if ch == nil {
		t.Fatal("got a nil channel (not implemented yet?)")
	}
	select {
	case v, ok := <-ch:
		return v, ok
	case <-time.After(2 * time.Second):
		t.Fatal("nothing was sent on the channel within 2 seconds")
	}
	panic("unreachable")
}

// drain receives until the channel is closed, failing if it never closes.
func drain[T any](t *testing.T, ch <-chan T) []T {
	t.Helper()
	var out []T
	for {
		v, ok := receive(t, ch)
		if !ok {
			return out
		}
		out = append(out, v)
		if len(out) > 10000 {
			t.Fatal("received more than 10000 values: the channel never closed")
		}
	}
}

// chanOf builds a closed channel holding vals, for testing stages without your own code.
func chanOf(vals ...int) <-chan int {
	ch := make(chan int, len(vals))
	for _, v := range vals {
		ch <- v
	}
	close(ch)
	return ch
}

func TestEx1Greeting(t *testing.T) {
	if got, _ := receive(t, Greeting("Ana")); got != "hello, Ana" {
		t.Errorf("Ex1: received %q, want %q", got, "hello, Ana")
	}
}

func TestEx2Numbers(t *testing.T) {
	if got := drain(t, Numbers(5)); !slices.Equal(got, []int{1, 2, 3, 4, 5}) {
		t.Errorf("Ex2: Numbers(5) sent %v, want [1 2 3 4 5]", got)
	}
	if got := drain(t, Numbers(0)); len(got) != 0 {
		t.Errorf("Ex2: Numbers(0) sent %v, want nothing (just close)", got)
	}
}

func TestEx3Sum(t *testing.T) {
	done := make(chan int)
	go func() { done <- Sum(chanOf(1, 2, 3, 4)) }()
	select {
	case got := <-done:
		if got != 10 {
			t.Errorf("Ex3: Sum(1, 2, 3, 4) = %d, want 10", got)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Ex3: Sum never returned. Does it stop when the channel is closed?")
	}
}

func TestEx4SumOfSquares(t *testing.T) {
	before := runtime.NumGoroutine()
	if got := SumOfSquares([]int{1, 2, 3, 4}); got != 30 {
		t.Errorf("Ex4: SumOfSquares([1 2 3 4]) = %d, want 30", got)
	}
	if got := SumOfSquares(nil); got != 0 {
		t.Errorf("Ex4: SumOfSquares(empty) = %d, want 0", got)
	}
	time.Sleep(20 * time.Millisecond)
	if after := runtime.NumGoroutine(); after > before {
		t.Errorf("Ex4: %d goroutines are still running afterwards (stuck sending?)", after-before)
	}
}

func TestEx5Doubler(t *testing.T) {
	if got := drain(t, Doubler(chanOf(1, 2, 3))); !slices.Equal(got, []int{2, 4, 6}) {
		t.Errorf("Ex5: Doubler(1, 2, 3) sent %v, want [2 4 6]", got)
	}
	if got := drain(t, Doubler(chanOf())); len(got) != 0 {
		t.Errorf("Ex5: Doubler(empty) sent %v, want nothing", got)
	}
}

func TestEx6Queue(t *testing.T) {
	before := runtime.NumGoroutine()
	q := Queue([]string{"a", "b", "c"})
	if runtime.NumGoroutine() > before {
		t.Errorf("Ex6: Queue started a goroutine. Use a buffered channel instead")
	}
	if q != nil && cap(q) < 3 {
		t.Errorf("Ex6: the channel's buffer holds %d, want room for all 3 items", cap(q))
	}
	if got := drain(t, q); !slices.Equal(got, []string{"a", "b", "c"}) {
		t.Errorf("Ex6: Queue sent %q, want [a b c]", got)
	}
}
