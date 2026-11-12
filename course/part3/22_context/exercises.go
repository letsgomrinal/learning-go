// Lesson 22 exercises: context. Read README.md first!
//
// Check one exercise at a time:   go test -race -run Ex1 -v
package lesson22

import (
	"context"
	"time"
)

// Ex1: wait for d, but stop early if ctx is cancelled.
// Return nil if the full time passed, or ctx.Err() if the context was cancelled first.
// (README section 4.)
func Sleep(ctx context.Context, d time.Duration) error {
	// TODO
	time.Sleep(d)
	return nil
}

// Ex2: add up the numbers one by one. Each step is slow: call step(n) to "process" a number,
// and add what it returns. Before each step, check whether ctx was cancelled; if so, stop and
// return the sum so far together with ctx.Err(). If all numbers are done, return the sum and nil.
func SumUntilCancelled(ctx context.Context, nums []int, step func(int) int) (int, error) {
	// TODO
	return 0, nil
}

// Ex3: return a channel that receives 1, 2, 3, ... (one number every `every`), until ctx is
// cancelled. Then the goroutine must stop and close the channel.
//
// Hint: a ticker, and a select with two cases: case <-ticker.C (send the next number, which
// itself should not block forever: use another select with ctx.Done()), and case <-ctx.Done().
func Ticker(ctx context.Context, every time.Duration) <-chan int {
	// TODO
	return nil
}

// Ex4: run work with a context that times out after d.
// Create it with context.WithTimeout(context.Background(), d), remember to defer cancel(),
// pass it to work, and return whatever work returns.
func WithTimeout(d time.Duration, work func(ctx context.Context) error) error {
	// TODO
	return nil
}

// Ex5: run every function at the same time, each with a context derived from ctx.
// Return the result of whichever finishes FIRST, and cancel the context so the others can stop.
// If ctx is cancelled before any finishes, return "".
//
// Hints:
//   - ctx, cancel := context.WithCancel(ctx); defer cancel()
//   - a BUFFERED results channel, len(fns), so slow goroutines can still send after you've returned
//   - select on the results channel and ctx.Done()
func FirstResult(ctx context.Context, fns []func(ctx context.Context) string) string {
	// TODO
	return ""
}

// Ex6: call fetch(ctx, id) for every id at the same time, and return the results in id order.
// If any fetch returns an error, cancel the others and return nil and that error.
// (README section 6. It's a mini errgroup.)
//
// Hints:
//   - ctx, cancel := context.WithCancel(ctx); defer cancel()
//   - results := make([]string, len(ids)); each goroutine writes its own index
//   - send each goroutine's error (or nil) on a buffered channel; receive len(ids) times;
//     on the first non-nil error, cancel() and return it
func FetchAll(ctx context.Context, ids []int, fetch func(ctx context.Context, id int) (string, error)) ([]string, error) {
	// TODO
	return nil, nil
}
