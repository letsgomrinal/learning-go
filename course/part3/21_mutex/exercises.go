// Lesson 21 exercises: mutexes and the race detector. Read README.md first!
//
// Check one exercise at a time:   go test -race -run Ex1 -v
//
// ⚠️ Run with -race. These tests call your code from many goroutines at once.
package lesson21

// Ex1: a counter that's safe to use from many goroutines.
// TODO: add a sync.Mutex field, and lock it in BOTH methods (README section 3).
// Its zero value must work:  var c Counter; c.Inc()
type Counter struct {
	n int
}

// Inc adds 1.
func (c *Counter) Inc() {
	c.n++
}

// Value returns the current count.
func (c *Counter) Value() int {
	return c.n
}

// Ex2: start `workers` goroutines. Each adds 1 to a shared total `each` times.
// Return the total (which must be workers × each). Protect the total with a mutex
// declared inside this function, and wait for all goroutines with a WaitGroup.
func CountConcurrently(workers, each int) int {
	// TODO
	return 0
}

// Ex3: counts visits per page, safe for concurrent use.
// Its zero value must work (create the map on first use, inside the lock).
type Visits struct {
	// TODO: a mutex and a map[string]int
}

// Add records one visit to page.
func (v *Visits) Add(page string) {
	// TODO
}

// Count returns the number of visits to page.
func (v *Visits) Count(page string) int {
	// TODO
	return 0
}

// Snapshot returns a COPY of all the counts. Changing the returned map must not affect
// the Visits. (Hint: maps.Clone, inside the lock.)
func (v *Visits) Snapshot() map[string]int {
	// TODO
	return nil
}

// Ex4: a cache. Use a sync.RWMutex: RLock in Get, Lock in Set (README section 5).
type Cache struct {
	// TODO: a sync.RWMutex and a map[string]string
}

// NewCache returns an empty, ready-to-use cache.
func NewCache() *Cache {
	// TODO
	return &Cache{}
}

// Get returns the value for key and true, or "" and false.
func (c *Cache) Get(key string) (string, bool) {
	// TODO
	return "", false
}

// Set stores value under key.
func (c *Cache) Set(key, value string) {
	// TODO
}

// Ex5: a hit counter using sync/atomic instead of a mutex (README section 6).
// TODO: replace the int field with an atomic.Int64 (import "sync/atomic").
type Hits struct {
	n int
}

func (h *Hits) Hit() {
	h.n++
}

func (h *Hits) Total() int64 {
	return int64(h.n)
}

// Ex6: Lazy holds a value that's computed the first time Get is called, and only then.
// Even if 50 goroutines call Get at the same moment, Compute must run exactly once,
// and all of them must get the same value. Use sync.Once (README section 7).
type Lazy struct {
	Compute func() string
	// TODO: a sync.Once and a field to store the value
}

func (l *Lazy) Get() string {
	// TODO
	return l.Compute()
}
