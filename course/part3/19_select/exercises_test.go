package lesson19

import (
	"runtime"
	"testing"
	"time"
)

// Always run these with:  go test -race -v

// finishes runs f and fails if it takes longer than limit.
func finishes(t *testing.T, limit time.Duration, f func()) {
	t.Helper()
	done := make(chan struct{})
	go func() { f(); close(done) }()
	select {
	case <-done:
	case <-time.After(limit):
		t.Fatalf("did not return within %v", limit)
	}
}

func TestEx1First(t *testing.T) {
	a := make(chan string)
	b := make(chan string)
	go func() { time.Sleep(80 * time.Millisecond); a <- "slow a" }()
	go func() { time.Sleep(10 * time.Millisecond); b <- "fast b" }()
	var got string
	finishes(t, time.Second, func() { got = First(a, b) })
	if got != "fast b" {
		t.Errorf("Ex1: First = %q, want %q", got, "fast b")
	}
	go func() { <-a }() // let the slow sender finish
}

func TestEx2ReceiveWithTimeout(t *testing.T) {
	ch := make(chan int, 1)
	ch <- 7
	if v, ok := ReceiveWithTimeout(ch, 50*time.Millisecond); v != 7 || !ok {
		t.Errorf("Ex2: value ready: got %d, %v; want 7, true", v, ok)
	}

	empty := make(chan int)
	start := time.Now()
	var v int
	var ok bool
	finishes(t, time.Second, func() { v, ok = ReceiveWithTimeout(empty, 30*time.Millisecond) })
	if v != 0 || ok {
		t.Errorf("Ex2: nothing sent: got %d, %v; want 0, false", v, ok)
	}
	if e := time.Since(start); e < 25*time.Millisecond {
		t.Errorf("Ex2: returned after %v, before the 30ms timeout", e)
	}

	closed := make(chan int)
	close(closed)
	if v, ok := ReceiveWithTimeout(closed, time.Second); v != 0 || ok {
		t.Errorf("Ex2: closed channel: got %d, %v; want 0, false", v, ok)
	}
}

func TestEx3TrySend(t *testing.T) {
	ch := make(chan int, 1)
	var first, second bool
	finishes(t, time.Second, func() {
		first = TrySend(ch, 1)
		second = TrySend(ch, 2)
	})
	if !first || second {
		t.Errorf("Ex3: TrySend on a 1-slot channel: got %v then %v; want true then false", first, second)
	}
	select {
	case v := <-ch:
		if v != 1 {
			t.Errorf("Ex3: the channel holds %d, want 1", v)
		}
	default:
		t.Errorf("Ex3: nothing was sent on the channel")
	}
}

func TestEx4TryReceive(t *testing.T) {
	ch := make(chan int, 1)
	var v int
	var ok bool
	finishes(t, time.Second, func() { v, ok = TryReceive(ch) })
	if v != 0 || ok {
		t.Errorf("Ex4: empty channel: got %d, %v; want 0, false", v, ok)
	}
	ch <- 5
	if v, ok := TryReceive(ch); v != 5 || !ok {
		t.Errorf("Ex4: value ready: got %d, %v; want 5, true", v, ok)
	}
	close(ch)
	if v, ok := TryReceive(ch); v != 0 || ok {
		t.Errorf("Ex4: closed channel: got %d, %v; want 0, false", v, ok)
	}
}

func TestEx5RunWithTimeout(t *testing.T) {
	if v, ok := RunWithTimeout(func() int { return 42 }, 200*time.Millisecond); v != 42 || !ok {
		t.Errorf("Ex5: fast work: got %d, %v; want 42, true", v, ok)
	}

	before := runtime.NumGoroutine()
	var v int
	var ok bool
	finishes(t, time.Second, func() {
		v, ok = RunWithTimeout(func() int { time.Sleep(60 * time.Millisecond); return 1 }, 10*time.Millisecond)
	})
	if v != 0 || ok {
		t.Errorf("Ex5: slow work: got %d, %v; want 0, false", v, ok)
	}

	time.Sleep(150 * time.Millisecond) // the slow work has finished by now
	if after := runtime.NumGoroutine(); after > before {
		t.Errorf("Ex5: a goroutine is still stuck after the slow work finished: it leaked. Use make(chan int, 1)")
	}
}

func TestEx6CountTicks(t *testing.T) {
	var n int
	finishes(t, time.Second, func() { n = CountTicks(10*time.Millisecond, 55*time.Millisecond) })
	if n < 3 || n > 6 {
		t.Errorf("Ex6: CountTicks(10ms, 55ms) = %d, want about 5", n)
	}
}

func TestEx7Counter(t *testing.T) {
	done := make(chan struct{})
	ch := Counter(done)
	if ch == nil {
		t.Fatal("Ex7: Counter returned a nil channel")
	}
	for want := range 3 {
		select {
		case v := <-ch:
			if v != want {
				t.Fatalf("Ex7: received %d, want %d", v, want)
			}
		case <-time.After(time.Second):
			t.Fatalf("Ex7: no value received (waiting for %d)", want)
		}
	}
	close(done)

	deadline := time.After(time.Second)
	for {
		select {
		case _, ok := <-ch:
			if !ok {
				return // closed: correct
			}
		case <-deadline:
			t.Fatal("Ex7: after close(done), the channel was never closed. Does the goroutine stop?")
		}
	}
}
