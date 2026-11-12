package lesson21

import (
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// Always run these with:  go test -race -v
// A data race makes the test fail with "race detected during execution of test".

func TestEx1Counter(t *testing.T) {
	var c Counter
	var wg sync.WaitGroup
	for range 50 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range 100 {
				c.Inc()
				_ = c.Value()
			}
		}()
	}
	wg.Wait()
	if got := c.Value(); got != 5000 {
		t.Errorf("Ex1: Value = %d, want 5000 (updates were lost: lock in Inc and Value)", got)
	}
}

func TestEx2CountConcurrently(t *testing.T) {
	if got := CountConcurrently(20, 500); got != 10000 {
		t.Errorf("Ex2: CountConcurrently(20, 500) = %d, want 10000", got)
	}
	if got := CountConcurrently(1, 3); got != 3 {
		t.Errorf("Ex2: CountConcurrently(1, 3) = %d, want 3", got)
	}
}

func TestEx3Visits(t *testing.T) {
	var v Visits
	var wg sync.WaitGroup
	pages := []string{"/", "/about", "/blog"}
	for i := range 30 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range 10 {
				v.Add(pages[i%3])
				_ = v.Count("/")
			}
		}()
	}
	wg.Wait()
	for _, p := range pages {
		if got := v.Count(p); got != 100 {
			t.Errorf("Ex3: Count(%q) = %d, want 100", p, got)
		}
	}
	if got := v.Count("/missing"); got != 0 {
		t.Errorf("Ex3: Count(missing) = %d, want 0", got)
	}

	snap := v.Snapshot()
	if len(snap) != 3 || snap["/"] != 100 {
		t.Fatalf("Ex3: Snapshot = %v, want 3 pages with 100 visits each", snap)
	}
	snap["/"] = -1
	if v.Count("/") != 100 {
		t.Errorf("Ex3: changing the snapshot changed the Visits. Return a copy")
	}
}

func TestEx4Cache(t *testing.T) {
	c := NewCache()
	var wg sync.WaitGroup
	for i := range 20 {
		wg.Add(2)
		go func() {
			defer wg.Done()
			c.Set("key"+strconv.Itoa(i), "value"+strconv.Itoa(i))
		}()
		go func() {
			defer wg.Done()
			c.Get("key" + strconv.Itoa(i))
		}()
	}
	wg.Wait()
	for i := range 20 {
		k, want := "key"+strconv.Itoa(i), "value"+strconv.Itoa(i)
		if got, ok := c.Get(k); got != want || !ok {
			t.Errorf("Ex4: Get(%q) = %q, %v; want %q, true", k, got, ok, want)
		}
	}
	if got, ok := c.Get("nope"); got != "" || ok {
		t.Errorf("Ex4: Get(missing) = %q, %v; want \"\", false", got, ok)
	}
}

func TestEx5Hits(t *testing.T) {
	var h Hits
	var wg sync.WaitGroup
	for range 40 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range 250 {
				h.Hit()
			}
		}()
	}
	wg.Wait()
	if got := h.Total(); got != 10000 {
		t.Errorf("Ex5: Total = %d, want 10000", got)
	}
}

func TestEx6Lazy(t *testing.T) {
	var calls atomic.Int32
	l := &Lazy{Compute: func() string {
		calls.Add(1)
		time.Sleep(20 * time.Millisecond) // slow, so many goroutines arrive while it runs
		return "config loaded"
	}}
	var wg sync.WaitGroup
	results := make([]string, 50)
	for i := range 50 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			results[i] = l.Get()
		}()
	}
	wg.Wait()
	if n := calls.Load(); n != 1 {
		t.Errorf("Ex6: Compute ran %d times, want exactly 1", n)
	}
	for i, r := range results {
		if r != "config loaded" {
			t.Fatalf("Ex6: goroutine %d got %q, want \"config loaded\"", i, r)
		}
	}
}
