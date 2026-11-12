package syncctx

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestSafeCounter(t *testing.T) {
	var c SafeCounter // zero value
	var wg sync.WaitGroup
	for range 100 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range 100 {
				c.Inc("hits")
			}
		}()
	}
	wg.Wait()
	if got := c.Value("hits"); got != 10_000 {
		t.Errorf("Value(hits) = %d, want 10000 (run with -race to see why)", got)
	}
	if got := c.Value("missing"); got != 0 {
		t.Errorf("Value(missing) = %d", got)
	}
}

func TestLazy(t *testing.T) {
	var calls atomic.Int32
	l := &Lazy[string]{Fn: func() string {
		calls.Add(1)
		time.Sleep(10 * time.Millisecond)
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
		t.Errorf("Fn called %d times, want exactly 1", n)
	}
	for i, r := range results {
		if r != "config loaded" {
			t.Fatalf("results[%d] = %q", i, r)
		}
	}
}

func TestSleep(t *testing.T) {
	if err := Sleep(context.Background(), 10*time.Millisecond); err != nil {
		t.Errorf("full sleep: err = %v, want nil", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	start := time.Now()
	err := Sleep(ctx, 5*time.Second)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("err = %v, want context.DeadlineExceeded", err)
	}
	if time.Since(start) > 500*time.Millisecond {
		t.Errorf("Sleep ignored ctx cancellation")
	}
}

func TestFetchAllSuccess(t *testing.T) {
	fetch := func(ctx context.Context, id int) (string, error) {
		time.Sleep(time.Duration(10-id) * time.Millisecond) // finish out of order
		return fmt.Sprintf("item-%d", id), nil
	}
	got, err := FetchAll(context.Background(), []int{1, 2, 3, 4}, fetch)
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"item-1", "item-2", "item-3", "item-4"}; !slices.Equal(got, want) {
		t.Errorf("FetchAll = %v, want %v", got, want)
	}
}

func TestFetchAllCancelsOthers(t *testing.T) {
	boom := errors.New("boom")
	var cancelled atomic.Int32

	fetch := func(ctx context.Context, id int) (string, error) {
		if id == 3 {
			time.Sleep(10 * time.Millisecond)
			return "", boom
		}
		select {
		case <-time.After(3 * time.Second): // a slow request...
			return "slow", nil
		case <-ctx.Done(): // ...that should be cancelled
			cancelled.Add(1)
			return "", ctx.Err()
		}
	}

	start := time.Now()
	_, err := FetchAll(context.Background(), []int{1, 2, 3, 4, 5}, fetch)
	if !errors.Is(err, boom) {
		t.Errorf("err = %v, want boom (the FIRST error)", err)
	}
	if time.Since(start) > time.Second {
		t.Errorf("FetchAll took %v: other fetches weren't cancelled", time.Since(start))
	}
	if n := cancelled.Load(); n != 4 {
		t.Errorf("%d fetches saw cancellation, want 4", n)
	}
}
