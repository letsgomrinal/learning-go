// Reference solution for lesson 21.
package lesson21

import (
	"maps"
	"sync"
	"sync/atomic"
)

type Counter struct {
	mu sync.Mutex
	n  int
}

func (c *Counter) Inc() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.n++
}

func (c *Counter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.n
}

func CountConcurrently(workers, each int) int {
	var (
		mu    sync.Mutex
		wg    sync.WaitGroup
		total int
	)
	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range each {
				mu.Lock()
				total++
				mu.Unlock()
			}
		}()
	}
	wg.Wait()
	return total
}

type Visits struct {
	mu     sync.Mutex
	counts map[string]int
}

func (v *Visits) Add(page string) {
	v.mu.Lock()
	defer v.mu.Unlock()
	if v.counts == nil {
		v.counts = map[string]int{}
	}
	v.counts[page]++
}

func (v *Visits) Count(page string) int {
	v.mu.Lock()
	defer v.mu.Unlock()
	return v.counts[page]
}

func (v *Visits) Snapshot() map[string]int {
	v.mu.Lock()
	defer v.mu.Unlock()
	return maps.Clone(v.counts)
}

type Cache struct {
	mu   sync.RWMutex
	data map[string]string
}

func NewCache() *Cache {
	return &Cache{data: map[string]string{}}
}

func (c *Cache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	v, ok := c.data[key]
	return v, ok
}

func (c *Cache) Set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = value
}

type Hits struct {
	n atomic.Int64
}

func (h *Hits) Hit() {
	h.n.Add(1)
}

func (h *Hits) Total() int64 {
	return h.n.Load()
}

type Lazy struct {
	Compute func() string
	once    sync.Once
	value   string
}

func (l *Lazy) Get() string {
	l.once.Do(func() {
		l.value = l.Compute()
	})
	return l.value
}
