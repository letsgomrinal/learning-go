package lesson34

import (
	"context"
	"errors"
	"slices"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// Always run these with:  go test -race -v

func within(t *testing.T, limit time.Duration, what string, f func()) {
	t.Helper()
	done := make(chan struct{})
	go func() { f(); close(done) }()
	select {
	case <-done:
	case <-time.After(limit):
		t.Fatalf("%s did not finish within %v", what, limit)
	}
}

func TestEx1Worker(t *testing.T) {
	w := &Worker{}
	finished := make(chan struct{})
	go func() {
		w.Run()
		close(finished)
	}()
	time.Sleep(10 * time.Millisecond)
	w.Stop()
	select {
	case <-finished:
	case <-time.After(2 * time.Second):
		t.Fatal("Ex1: Run never noticed Stop (the compiler may have cached the flag). Use atomic.Bool")
	}
	if w.Rounds == 0 {
		t.Error("Ex1: Run didn't loop at all")
	}
}

func TestEx2Singleflight(t *testing.T) {
	var g Group
	var calls atomic.Int32
	fn := func() (any, error) {
		calls.Add(1)
		time.Sleep(50 * time.Millisecond)
		return "loaded", nil
	}
	var wg sync.WaitGroup
	var sharedCount atomic.Int32
	results := make([]any, 20)
	for i := range 20 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			v, err, shared := g.Do("user:42", fn)
			if err != nil {
				t.Errorf("Ex2: unexpected error %v", err)
			}
			if shared {
				sharedCount.Add(1)
			}
			results[i] = v
		}()
	}
	within(t, 2*time.Second, "Ex2: 20 concurrent Do calls", wg.Wait)
	if n := calls.Load(); n != 1 {
		t.Errorf("Ex2: fn ran %d times for 20 concurrent calls with the same key, want 1", n)
	}
	if s := sharedCount.Load(); s != 19 {
		t.Errorf("Ex2: %d callers got shared == true, want 19", s)
	}
	for i, r := range results {
		if r != "loaded" {
			t.Fatalf("Ex2: caller %d got %v", i, r)
		}
	}
	// After it finished, a new call runs fn again.
	g.Do("user:42", fn)
	if n := calls.Load(); n != 2 {
		t.Errorf("Ex2: a later call should run fn again (total %d, want 2). Delete the key when done", n)
	}
	// Errors are shared too.
	boom := errors.New("boom")
	if _, err, _ := g.Do("bad", func() (any, error) { return nil, boom }); !errors.Is(err, boom) {
		t.Errorf("Ex2: error = %v, want boom", err)
	}
}

func TestEx3BoundedQueue(t *testing.T) {
	q := NewBoundedQueue(2)
	within(t, time.Second, "Ex3: Put, Put", func() { q.Put(1); q.Put(2) })

	putDone := make(chan struct{})
	go func() { q.Put(3); close(putDone) }() // must wait: the queue is full
	select {
	case <-putDone:
		t.Fatal("Ex3: Put on a full queue didn't wait")
	case <-time.After(50 * time.Millisecond):
	}
	var first int
	within(t, time.Second, "Ex3: Get", func() { first = q.Get() })
	if first != 1 {
		t.Errorf("Ex3: Get = %d, want 1 (first in, first out)", first)
	}
	within(t, time.Second, "Ex3: the waiting Put after a Get", func() { <-putDone })

	// Many producers and consumers.
	q2 := NewBoundedQueue(3)
	var got []int
	var mu sync.Mutex
	within(t, 3*time.Second, "Ex3: 4 producers and 4 consumers", func() {
		var wg sync.WaitGroup
		for p := range 4 {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for i := range 25 {
					q2.Put(p*100 + i)
				}
			}()
		}
		for range 4 {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for range 25 {
					v := q2.Get()
					mu.Lock()
					got = append(got, v)
					mu.Unlock()
				}
			}()
		}
		wg.Wait()
	})
	if len(got) != 100 {
		t.Fatalf("Ex3: consumers got %d items, want 100", len(got))
	}
	slices.Sort(got)
	if got[0] != 0 || got[99] != 324 {
		t.Errorf("Ex3: items were lost or duplicated")
	}
}

func TestEx4ErrGroup(t *testing.T) {
	g, ctx := WithContext(context.Background())
	var ran atomic.Int32
	for range 5 {
		g.Go(func() error { time.Sleep(10 * time.Millisecond); ran.Add(1); return nil })
	}
	var err error
	within(t, 2*time.Second, "Ex4: Wait", func() { err = g.Wait() })
	if err != nil || ran.Load() != 5 {
		t.Errorf("Ex4: all succeed: err %v, %d ran; want nil, 5 (Wait must wait for all)", err, ran.Load())
	}
	if ctx.Err() == nil {
		t.Errorf("Ex4: the context should be cancelled after Wait returns")
	}

	boom := errors.New("boom")
	g, ctx = WithContext(context.Background())
	var cancelled atomic.Int32
	for i := range 5 {
		g.Go(func() error {
			if i == 2 {
				time.Sleep(10 * time.Millisecond)
				return boom
			}
			select {
			case <-ctx.Done():
				cancelled.Add(1)
				return ctx.Err()
			case <-time.After(3 * time.Second):
				return nil
			}
		})
	}
	within(t, 2*time.Second, "Ex4: Wait after a failure", func() { err = g.Wait() })
	if !errors.Is(err, boom) {
		t.Errorf("Ex4: Wait = %v, want the first error (boom)", err)
	}
	if cancelled.Load() != 4 {
		t.Errorf("Ex4: %d of 4 other functions saw the cancellation", cancelled.Load())
	}
}

func TestEx5Settings(t *testing.T) {
	var s Settings
	s.Store(&Config{RateLimit: 10, Banned: map[string]bool{}})
	var wg sync.WaitGroup
	stop := make(chan struct{})
	for range 8 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-stop:
					return
				default:
					if c := s.Current(); c == nil || c.RateLimit < 10 {
						t.Error("Ex5: read an incomplete config")
						return
					}
				}
			}
		}()
	}
	for i := range 100 {
		s.Store(&Config{RateLimit: 10 + i, Banned: map[string]bool{"x": true}})
	}
	close(stop)
	wg.Wait()
	if got := s.Current().RateLimit; got != 109 {
		t.Errorf("Ex5: final RateLimit = %d, want 109", got)
	}
}
