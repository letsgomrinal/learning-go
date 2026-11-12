package lesson22

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"sync/atomic"
	"testing"
	"time"
)

// Always run these with:  go test -race -v

func within(t *testing.T, limit time.Duration, f func()) {
	t.Helper()
	done := make(chan struct{})
	go func() { f(); close(done) }()
	select {
	case <-done:
	case <-time.After(limit):
		t.Fatalf("did not return within %v. Is the context being checked?", limit)
	}
}

func TestEx1Sleep(t *testing.T) {
	if err := Sleep(context.Background(), 10*time.Millisecond); err != nil {
		t.Errorf("Ex1: full sleep returned %v, want nil", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	var err error
	within(t, 500*time.Millisecond, func() { err = Sleep(ctx, 5*time.Second) })
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("Ex1: cancelled sleep returned %v, want context.DeadlineExceeded", err)
	}
}

func TestEx2SumUntilCancelled(t *testing.T) {
	step := func(n int) int { time.Sleep(10 * time.Millisecond); return n }
	if sum, err := SumUntilCancelled(context.Background(), []int{1, 2, 3}, step); sum != 6 || err != nil {
		t.Errorf("Ex2: not cancelled: got %d, %v; want 6, nil", sum, err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	calls := 0
	cancelAfter2 := func(n int) int {
		calls++
		if calls == 2 {
			cancel()
		}
		return n
	}
	sum, err := SumUntilCancelled(ctx, []int{10, 20, 30, 40}, cancelAfter2)
	if sum != 30 || !errors.Is(err, context.Canceled) {
		t.Errorf("Ex2: cancelled during step 2: got %d, %v; want 30 (10+20), context.Canceled", sum, err)
	}
	if calls != 2 {
		t.Errorf("Ex2: step was called %d times, want 2 (check ctx.Err() BEFORE each step)", calls)
	}
}

func TestEx3Ticker(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ch := Ticker(ctx, 5*time.Millisecond)
	if ch == nil {
		t.Fatal("Ex3: Ticker returned a nil channel")
	}
	for want := 1; want <= 3; want++ {
		select {
		case v := <-ch:
			if v != want {
				t.Fatalf("Ex3: received %d, want %d", v, want)
			}
		case <-time.After(time.Second):
			t.Fatalf("Ex3: no value received (waiting for %d)", want)
		}
	}
	cancel()
	deadline := time.After(time.Second)
	for {
		select {
		case _, ok := <-ch:
			if !ok {
				return
			}
		case <-deadline:
			t.Fatal("Ex3: after cancel(), the channel was never closed")
		}
	}
}

func TestEx4WithTimeout(t *testing.T) {
	err := WithTimeout(time.Second, func(ctx context.Context) error { return nil })
	if err != nil {
		t.Errorf("Ex4: fast work returned %v, want nil", err)
	}
	var sawDeadline bool
	within(t, 500*time.Millisecond, func() {
		err = WithTimeout(20*time.Millisecond, func(ctx context.Context) error {
			_, sawDeadline = ctx.Deadline()
			<-ctx.Done()
			return ctx.Err()
		})
	})
	if !sawDeadline {
		t.Errorf("Ex4: the context passed to work has no deadline. Use context.WithTimeout")
	}
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("Ex4: slow work returned %v, want context.DeadlineExceeded", err)
	}
	boom := errors.New("boom")
	if err := WithTimeout(time.Second, func(context.Context) error { return boom }); !errors.Is(err, boom) {
		t.Errorf("Ex4: return work's error unchanged; got %v", err)
	}
}

func TestEx5FirstResult(t *testing.T) {
	var cancelled atomic.Int32
	worker := func(name string, d time.Duration) func(context.Context) string {
		return func(ctx context.Context) string {
			select {
			case <-time.After(d):
				return name
			case <-ctx.Done():
				cancelled.Add(1)
				return ""
			}
		}
	}
	fns := []func(context.Context) string{
		worker("slow", 500*time.Millisecond),
		worker("fast", 10*time.Millisecond),
		worker("slower", time.Second),
	}
	var got string
	within(t, 300*time.Millisecond, func() { got = FirstResult(context.Background(), fns) })
	if got != "fast" {
		t.Errorf("Ex5: FirstResult = %q, want %q", got, "fast")
	}
	time.Sleep(50 * time.Millisecond)
	if n := cancelled.Load(); n != 2 {
		t.Errorf("Ex5: %d of the 2 slower functions saw the cancellation, want 2. Did you cancel?", n)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	within(t, 300*time.Millisecond, func() { got = FirstResult(ctx, fns[:1]) })
	if got != "" {
		t.Errorf("Ex5: when ctx is cancelled first, got %q, want \"\"", got)
	}
}

func TestEx6FetchAll(t *testing.T) {
	okFetch := func(ctx context.Context, id int) (string, error) {
		time.Sleep(time.Duration(10-id) * time.Millisecond)
		return fmt.Sprintf("item-%d", id), nil
	}
	got, err := FetchAll(context.Background(), []int{1, 2, 3, 4}, okFetch)
	if err != nil || !slices.Equal(got, []string{"item-1", "item-2", "item-3", "item-4"}) {
		t.Fatalf("Ex6: all succeed: got %q, %v; want [item-1 item-2 item-3 item-4], nil", got, err)
	}

	boom := errors.New("boom")
	var cancelled atomic.Int32
	failing := func(ctx context.Context, id int) (string, error) {
		if id == 3 {
			time.Sleep(10 * time.Millisecond)
			return "", boom
		}
		select {
		case <-time.After(2 * time.Second):
			return "slow", nil
		case <-ctx.Done():
			cancelled.Add(1)
			return "", ctx.Err()
		}
	}
	within(t, time.Second, func() { got, err = FetchAll(context.Background(), []int{1, 2, 3, 4, 5}, failing) })
	if !errors.Is(err, boom) || got != nil {
		t.Errorf("Ex6: one fails: got %q, %v; want nil, boom", got, err)
	}
	time.Sleep(50 * time.Millisecond)
	if n := cancelled.Load(); n != 4 {
		t.Errorf("Ex6: %d of the 4 other fetches saw the cancellation, want 4", n)
	}
}
