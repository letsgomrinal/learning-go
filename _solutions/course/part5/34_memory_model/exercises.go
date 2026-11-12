// Reference solution for lesson 34.
package lesson34

import (
	"context"
	"sync"
	"sync/atomic"
)

type Worker struct {
	stopped atomic.Bool
	Rounds  int
}

func (w *Worker) Run() {
	for !w.stopped.Load() {
		w.Rounds++
	}
}

func (w *Worker) Stop() {
	w.stopped.Store(true)
}

type call struct {
	wg  sync.WaitGroup
	val any
	err error
}

type Group struct {
	mu    sync.Mutex
	calls map[string]*call
}

func (g *Group) Do(key string, fn func() (any, error)) (v any, err error, shared bool) {
	g.mu.Lock()
	if g.calls == nil {
		g.calls = map[string]*call{}
	}
	if c, ok := g.calls[key]; ok {
		g.mu.Unlock()
		c.wg.Wait()
		return c.val, c.err, true
	}
	c := &call{}
	c.wg.Add(1)
	g.calls[key] = c
	g.mu.Unlock()

	c.val, c.err = fn()
	c.wg.Done()

	g.mu.Lock()
	delete(g.calls, key)
	g.mu.Unlock()
	return c.val, c.err, false
}

type BoundedQueue struct {
	mu       sync.Mutex
	notEmpty *sync.Cond
	notFull  *sync.Cond
	items    []int
	capacity int
}

func NewBoundedQueue(capacity int) *BoundedQueue {
	q := &BoundedQueue{capacity: capacity}
	q.notEmpty = sync.NewCond(&q.mu)
	q.notFull = sync.NewCond(&q.mu)
	return q
}

func (q *BoundedQueue) Put(v int) {
	q.mu.Lock()
	for len(q.items) >= q.capacity {
		q.notFull.Wait()
	}
	q.items = append(q.items, v)
	q.mu.Unlock()
	q.notEmpty.Signal()
}

func (q *BoundedQueue) Get() int {
	q.mu.Lock()
	for len(q.items) == 0 {
		q.notEmpty.Wait()
	}
	v := q.items[0]
	q.items = q.items[1:]
	q.mu.Unlock()
	q.notFull.Signal()
	return v
}

func (q *BoundedQueue) Len() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.items)
}

type ErrGroup struct {
	wg     sync.WaitGroup
	once   sync.Once
	err    error
	cancel context.CancelFunc
}

func WithContext(ctx context.Context) (*ErrGroup, context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	return &ErrGroup{cancel: cancel}, ctx
}

func (g *ErrGroup) Go(f func() error) {
	g.wg.Add(1)
	go func() {
		defer g.wg.Done()
		if err := f(); err != nil {
			g.once.Do(func() {
				g.err = err
				if g.cancel != nil {
					g.cancel()
				}
			})
		}
	}()
}

func (g *ErrGroup) Wait() error {
	g.wg.Wait()
	if g.cancel != nil {
		g.cancel()
	}
	return g.err
}

type Config struct {
	RateLimit int
	Banned    map[string]bool
}

type Settings struct {
	current atomic.Pointer[Config]
}

func (s *Settings) Store(c *Config) {
	s.current.Store(c)
}

func (s *Settings) Current() *Config {
	return s.current.Load()
}
