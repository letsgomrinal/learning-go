// Reference solution for lesson 08.
package syncctx

import (
	"context"
	"sync"
	"time"
)

type SafeCounter struct {
	mu sync.Mutex
	m  map[string]int
}

func (c *SafeCounter) Inc(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.m == nil {
		c.m = make(map[string]int)
	}
	c.m[key]++
}

func (c *SafeCounter) Value(key string) int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.m[key]
}

type Lazy[T any] struct {
	once sync.Once
	val  T
	Fn   func() T
}

func (l *Lazy[T]) Get() T {
	l.once.Do(func() { l.val = l.Fn() })
	return l.val
}

func Sleep(ctx context.Context, d time.Duration) error {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-t.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func FetchAll(ctx context.Context, ids []int, fetch func(ctx context.Context, id int) (string, error)) ([]string, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	results := make([]string, len(ids))
	var (
		wg       sync.WaitGroup
		once     sync.Once
		firstErr error
	)
	for i, id := range ids {
		wg.Add(1)
		go func() {
			defer wg.Done()
			res, err := fetch(ctx, id)
			if err != nil {
				once.Do(func() {
					firstErr = err
					cancel() // tell everyone else to stop
				})
				return
			}
			results[i] = res
		}()
	}
	wg.Wait()

	if firstErr != nil {
		return nil, firstErr
	}
	return results, nil
}
