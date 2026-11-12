// Lesson 34 exercises: the memory model and advanced sync. Read README.md first!
//
// Check one exercise at a time:   go test -race -run Ex1 -v
// The race detector is part of these tests: always use -race.
package lesson34

import (
	"context"
	"sync"
)

// Ex1: a worker that counts loop rounds until it's stopped.
// BUG: stopped is a plain bool read and written by different goroutines: a data race
// (README section 1). Make it safe with an atomic.Bool (import "sync/atomic").
type Worker struct {
	stopped bool
	Rounds  int
}

// Run loops until Stop is called. Only the goroutine running Run touches Rounds.
func (w *Worker) Run() {
	for !w.stopped {
		w.Rounds++
	}
}

// Stop asks Run to finish. It's called from a different goroutine.
func (w *Worker) Stop() {
	w.stopped = true
}

// ── Ex2: singleflight ─────────────────────────────────────────────────────

// call is one in-progress (or finished) execution of a function.
type call struct {
	wg  sync.WaitGroup
	val any
	err error
}

// Group deduplicates concurrent calls with the same key (README section 4).
// Its zero value must be ready to use.
type Group struct {
	mu    sync.Mutex
	calls map[string]*call
}

// Do runs fn for key, unless a call for the same key is already in progress: then it waits for that
// call and returns its result, with shared == true. Once a call finishes, the key is forgotten, so a
// later Do runs fn again.
//
// Steps:
//  1. lock; if calls is nil, create it
//  2. if calls[key] exists: unlock, c.wg.Wait(), return c.val, c.err, true
//  3. otherwise: c := &call{}; c.wg.Add(1); calls[key] = c; unlock
//  4. c.val, c.err = fn(); c.wg.Done()
//  5. lock; delete(calls, key); unlock; return c.val, c.err, false
func (g *Group) Do(key string, fn func() (any, error)) (v any, err error, shared bool) {
	// TODO
	v, err = fn()
	return v, err, false
}

// ── Ex3: a bounded blocking queue with sync.Cond ──────────────────────────

// BoundedQueue holds at most capacity items. Put waits while it's full; Get waits while it's empty.
type BoundedQueue struct {
	mu       sync.Mutex
	notEmpty *sync.Cond
	notFull  *sync.Cond
	items    []int
	capacity int
}

// NewBoundedQueue creates a queue. Already done: both conditions share the same mutex.
func NewBoundedQueue(capacity int) *BoundedQueue {
	q := &BoundedQueue{capacity: capacity}
	q.notEmpty = sync.NewCond(&q.mu)
	q.notFull = sync.NewCond(&q.mu)
	return q
}

// Put adds v at the end, waiting while the queue is full. Then signal notEmpty.
// (README section 6: lock, `for` the condition { Wait }, change, Signal.)
func (q *BoundedQueue) Put(v int) {
	// TODO
}

// Get removes and returns the first item, waiting while the queue is empty. Then signal notFull.
func (q *BoundedQueue) Get() int {
	// TODO
	return 0
}

// Len returns the number of items (under the lock).
func (q *BoundedQueue) Len() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.items)
}

// ── Ex4: errgroup ─────────────────────────────────────────────────────────

// ErrGroup runs functions in goroutines, waits for all of them, and remembers the FIRST error.
// When created with WithContext, the first error also cancels the context (README section 5).
type ErrGroup struct {
	wg     sync.WaitGroup
	once   sync.Once
	err    error
	cancel context.CancelFunc
}

// WithContext returns a new group and a context derived from ctx that is cancelled when any function
// returns an error, or when Wait returns.
func WithContext(ctx context.Context) (*ErrGroup, context.Context) {
	// TODO
	return &ErrGroup{}, ctx
}

// Go runs f in a new goroutine. If f returns an error, record it (only the first one, using once)
// and call cancel if it isn't nil.
func (g *ErrGroup) Go(f func() error) {
	// TODO
	go f()
}

// Wait waits for all functions, calls cancel (if not nil), and returns the first error (or nil).
func (g *ErrGroup) Wait() error {
	// TODO
	return nil
}

// ── Ex5: hot-reloadable settings ──────────────────────────────────────────

// Config must never be modified after it's been published with Store.
type Config struct {
	RateLimit int
	Banned    map[string]bool
}

// Settings holds the current Config. Many goroutines call Current(); occasionally one calls Store().
// TODO: replace the plain pointer with an atomic.Pointer[Config] so this is race-free (README section 3).
type Settings struct {
	current *Config
}

func (s *Settings) Store(c *Config) {
	s.current = c
}

func (s *Settings) Current() *Config {
	return s.current
}
