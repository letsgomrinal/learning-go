# Lesson 21: Mutexes and the race detector

**Goal:** understand data races properly, find them with `-race`, and protect shared data with
`sync.Mutex`, `sync.RWMutex`, atomics and `sync.Once`.

**Time:** about 1.5 hours.

```bash
cd ~/Desktop/learning-go/course/part3/21_mutex
```

---

## 1. Seeing a data race

Type this into `try/main.go` and run it **without** `-race` a few times:

```go
func main() {
	count := 0
	var wg sync.WaitGroup
	for range 1000 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			count++
		}()
	}
	wg.Wait()
	fmt.Println(count)
}
```

You expect `1000`. You'll often get something smaller, like `943`, and a different number each time.

Why? `count++` looks like one step, but the CPU does three: **read** `count`, **add 1**, **write it
back**. Two goroutines on different cores can both read `57`, both compute `58`, and both write `58`.
One increment is lost.

This is a **data race**: two goroutines access the same variable at the same time, and at least one
of them writes. The result depends on timing, so the bug can hide for weeks and then appear in
production.

Python's GIL hides many (not all) of these problems. Go has no GIL, so you have to handle them yourself.

---

## 2. The race detector

Now run it with `-race`:

```bash
go run -race ./try
```

```
==================
WARNING: DATA RACE
Read at 0x00c000012345 by goroutine 8:
  main.main.func1()
      /.../try/main.go:14 +0x3c

Previous write at 0x00c000012345 by goroutine 7:
  main.main.func1()
      /.../try/main.go:14 +0x4e
...
```

It tells you **which line** and **which goroutines**. This is why this course has had you use `-race`
since lesson 17. Make it a habit for all tests: `go test -race ./...`. (It makes the program slower,
so don't use it for production builds.)

The race detector only finds races that actually **happen** while it's watching. A passing run with
`-race` is strong evidence, not a proof. Good tests that exercise the concurrent code help it.

---

## 3. `sync.Mutex`: one at a time

A **mutex** ("mutual exclusion") is a lock. Only one goroutine can hold it at a time. Others that try
to `Lock()` wait until it's `Unlock()`ed:

```go
var mu sync.Mutex
count := 0

for range 1000 {
	wg.Add(1)
	go func() {
		defer wg.Done()
		mu.Lock()
		count++          // only one goroutine at a time can be here
		mu.Unlock()
	}()
}
```

Now it always prints `1000`, and `-race` is silent. The code between `Lock` and `Unlock` is called
the **critical section**.

Python comparison: `threading.Lock()` with `with lock:`.

### The standard shape: a struct with a mutex

In real code, the mutex lives **next to the data it protects**, inside a struct, and the methods do
the locking. Callers never touch the lock:

```go
type Counter struct {
	mu sync.Mutex
	n  int
}

func (c *Counter) Inc() {
	c.mu.Lock()
	defer c.mu.Unlock()       // unlock when the method returns, even on early return or panic
	c.n++
}

func (c *Counter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.n
}
```

Things to notice:
- `Lock()` then immediately `defer Unlock()`, exactly like `defer f.Close()` in lesson 14.
- **Reads need the lock too** (`Value`), not just writes. Reading while another goroutine writes is still a race.
- The zero value of `sync.Mutex` is an unlocked mutex: `var c Counter` is ready to use.
- **Pointer receivers**, always. Copying a struct copies its mutex, and then the copy and the original
  have different locks. `go vet` warns about this (`passes lock by value`).

---

## 4. Maps need protection

Go maps are **not** safe for concurrent use. Concurrent writes don't just give wrong answers: they can
crash the whole program:

```
fatal error: concurrent map writes
```

This is a `fatal error`, not a panic: `recover` can't catch it. Protect shared maps with a mutex:

```go
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
```

---

## 5. `sync.RWMutex`: many readers, one writer

If data is read much more often than it's written (a cache, a config), a `sync.RWMutex` allows **many
readers at once**, but a writer gets exclusive access:

```go
type Cache struct {
	mu   sync.RWMutex
	data map[string]string
}

func (c *Cache) Get(key string) (string, bool) {
	c.mu.RLock()              // read lock: other readers can hold it at the same time
	defer c.mu.RUnlock()
	v, ok := c.data[key]
	return v, ok
}

func (c *Cache) Set(key, value string) {
	c.mu.Lock()               // write lock: waits until no readers or writers hold it
	defer c.mu.Unlock()
	c.data[key] = value
}
```

Start with a plain `Mutex`. Switch to `RWMutex` only when lots of concurrent reads are a real bottleneck.

---

## 6. Atomics: for single numbers

For a simple counter or flag, `sync/atomic` is lighter than a mutex. Its operations happen as one
indivisible step:

```go
var count atomic.Int64

for range 1000 {
	wg.Add(1)
	go func() {
		defer wg.Done()
		count.Add(1)
	}()
}
wg.Wait()
fmt.Println(count.Load())   // 1000
```

Use atomics for one number on its own. When several values must change together (a balance and a
transaction history, for example), use a mutex.

---

## 7. `sync.Once`: do something exactly once

Lazy initialisation (load a config or connect to a database the first time it's needed) must happen
**once**, even if many goroutines ask at the same moment:

```go
var (
	once   sync.Once
	config map[string]string
)

func getConfig() map[string]string {
	once.Do(func() {
		fmt.Println("loading config...")   // runs only the first time, ever
		config = map[string]string{"port": "8080"}
	})
	return config
}
```

If 100 goroutines call `getConfig()` at once, the function inside `Do` runs exactly once, and all of
them wait until it has finished.

---

## 8. Channels or mutexes?

Both are fine, and good Go code uses both.

| Use a **mutex** when... | Use a **channel** when... |
|---|---|
| several goroutines share **state** (a counter, a map, a cache) | you're passing **data or ownership** from one goroutine to another |
| the critical section is short and simple | you're coordinating steps: pipelines, worker pools, signals |

A mutex around a map is simpler than a goroutine that owns the map and serves requests over channels.
Don't force everything through channels.

### 🤔 Predict the output

```go
type Box struct {
	mu sync.Mutex
	n  int
}

func (b Box) Add() {       // ← value receiver
	b.mu.Lock()
	defer b.mu.Unlock()
	b.n++
}

func main() {
	var b Box
	b.Add()
	b.Add()
	fmt.Println(b.n)
}
```

<details markdown="1">
<summary>Answer</summary>

`0`. With a value receiver, `Add` works on a **copy** of the box (lesson 08), mutex included, so the
original never changes. `go vet` reports `Add passes lock by value`. Use `func (b *Box) Add()`.
</details>

---

## 9. Common mistakes

| You wrote | What happens | Fix |
|---|---|---|
| Locking for writes but not for reads | still a data race | lock in every method that touches the data |
| `mu.Lock()` without `defer mu.Unlock()`, then an early return | the lock is never released: everyone waits forever | `defer mu.Unlock()` right after `Lock()` |
| Value receivers on a struct with a mutex | each call locks a copy | pointer receivers; `go vet` catches this |
| Calling a locking method from inside another locked method | deadlock: Go's mutexes are not re-entrant | split into a locked public method and an unlocked private helper |
| Concurrent writes to a plain map | `fatal error: concurrent map writes` | protect it with a mutex |
| Returning the internal map from a locked method | callers use it without the lock | return a copy (`maps.Clone`) |

---

## 10. Exercises

Run with: `go test -race -run Ex1 -v`. **The race detector is part of the test here**: a solution
that gets the right number by luck but has a race will be reported.

| | Exercise | Practises |
|---|---|---|
| Ex1 | `Counter` with `Inc` and `Value` | a struct with a mutex |
| Ex2 | `CountConcurrently` | a mutex inside a function |
| Ex3 | `Visits`: counts per page, plus a safe copy | protecting a map |
| Ex4 | `Cache` with `Get` and `Set` | `RWMutex` |
| Ex5 | `Hits` | `atomic.Int64` |
| Ex6 | `Lazy`: compute a value only once | `sync.Once` |

---

## ✅ Checklist

- [ ] I can explain why `count++` from many goroutines loses updates
- [ ] I use `-race` for all concurrent code, and can read its report
- [ ] I can protect a struct's data with a mutex, locking for reads and writes, with `defer Unlock()`
- [ ] I know maps aren't safe for concurrent writes
- [ ] I know when to use `RWMutex`, atomics and `sync.Once`

**Next → [Lesson 22: Context](../22_context/README.md)**
