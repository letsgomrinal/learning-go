// Lesson 08: Shared state (sync) and cancellation (context).
//
// Run:  go test -race -v ./practice/08_sync_context
//
// ─── CONCEPTS ──────────────────────────────────────────────────────────────
//
// 1. DATA RACES ARE REAL in Go (no GIL). Two goroutines doing m[k]++ on the same map
//    can corrupt it or crash with "fatal error: concurrent map writes".
//    The -race flag detects races at runtime. Use it in tests and CI.
//
// 2. sync.Mutex protects shared state:
//      type Counter struct {
//          mu sync.Mutex   // zero value is an unlocked mutex, ready to use
//          n  int
//      }
//      func (c *Counter) Inc() { c.mu.Lock(); defer c.mu.Unlock(); c.n++ }
//    - NEVER copy a struct that contains a mutex (use pointer receivers). `go vet` catches this.
//    - sync.RWMutex: many readers OR one writer (RLock/RUnlock for reads).
//    - sync/atomic (atomic.Int64 etc.) for simple counters and flags.
//    - Channels vs mutex: channels transfer OWNERSHIP of data; mutexes guard STATE.
//      For a shared map or counter, a mutex is simpler and perfectly idiomatic.
//
// 3. sync.Once runs something exactly once, even under concurrency (lazy init).
//
// 4. context.Context carries cancellation, deadlines and request-scoped values across
//    API boundaries. It's how you stop work nobody needs anymore (client disconnected,
//    timeout hit, sibling task failed).
//      ctx, cancel := context.WithTimeout(parent, 2*time.Second)
//      defer cancel()                       // ALWAYS call cancel to free resources
//      ctx, cancel := context.WithCancel(parent)
//      <-ctx.Done()                         // closed when cancelled or deadline passes
//      ctx.Err()                            // context.Canceled or context.DeadlineExceeded
//    Rules: ctx is the FIRST parameter, named ctx; never store it in a struct;
//    never pass nil (use context.Background() or context.TODO()).
//    Cancellation is COOPERATIVE: your code must check ctx.Done() in a select.
//    Nothing gets killed for you.
//
// 5. Python analogy: asyncio task.cancel() + asyncio.wait_for, except explicit, and it
//    works across goroutines, HTTP requests and DB queries (db.QueryContext(ctx, ...)).
// ───────────────────────────────────────────────────────────────────────────

package syncctx

import (
	"context"
	"time"
)

// SafeCounter counts occurrences per key and is safe for concurrent use.
// Its ZERO VALUE must be usable:  var c SafeCounter; c.Inc("x")
type SafeCounter struct {
	// TODO: add a sync.Mutex and a map[string]int
}

// Inc increments the count for key.
func (c *SafeCounter) Inc(key string) {
	// TODO (lazily create the map if it's nil)
}

// Value returns the count for key.
func (c *SafeCounter) Value(key string) int {
	// TODO
	return 0
}

// Lazy computes a value once, on first use, even if Get is called from many
// goroutines at the same time.
type Lazy[T any] struct {
	// TODO: add a sync.Once and a field for the value
	Fn func() T
}

// Get returns the value, calling Fn exactly once in total.
func (l *Lazy[T]) Get() T {
	// TODO
	var zero T
	return zero
}

// Sleep waits for d, but returns early with ctx.Err() if ctx is cancelled first.
// It returns nil if the full duration elapsed.
//
// Hint: select on time.After(d) (or a time.NewTimer you Stop) and ctx.Done().
func Sleep(ctx context.Context, d time.Duration) error {
	// TODO
	return nil
}

// FetchAll calls fetch for every id CONCURRENTLY and returns the results in id order.
//
// If any fetch returns an error:
//   - FetchAll returns that first error (results can be nil)
//   - the context passed to the OTHER fetches must be cancelled, so they can stop early
//
// This is a mini version of golang.org/x/sync/errgroup, worth building once by hand.
//
// Hint:
//
//	ctx, cancel := context.WithCancel(ctx); defer cancel()
//	results := make([]string, len(ids))
//	var wg sync.WaitGroup; var once sync.Once; var firstErr error
//	for i, id := range ids { wg.Add(1); go func() { ... on error: once.Do(func(){ firstErr = err; cancel() }) }() }
//	wg.Wait()
func FetchAll(ctx context.Context, ids []int, fetch func(ctx context.Context, id int) (string, error)) ([]string, error) {
	// TODO
	return nil, nil
}
