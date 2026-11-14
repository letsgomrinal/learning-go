# Lesson 34: The memory model and advanced sync

**Goal:** understand *why* unsynchronized code is broken even when it "looks fine", what
"happens-before" means, and build the advanced tools real services use: atomic pointers,
singleflight, an errgroup, and a blocking queue with `sync.Cond`.

**Time:** about 2 hours.

```bash
cd ~/Desktop/learning-go/course/part5/34_memory_model
```

---

## 1. A program that "should" work

```go
var ready bool
var message string

func main() {
	go func() {
		message = "hello"
		ready = true
	}()
	for !ready {                 // wait until the goroutine sets ready
	}
	fmt.Println(message)
}
```

It seems obvious that this prints `hello`. But it's **broken**, and the Go memory model allows all of
these outcomes:
- it prints `hello`;
- it prints an **empty** message: the other core sees `ready = true` before it sees `message = "hello"`;
- it **never ends**: the compiler may read `ready` once, keep it in a register, and loop forever.

Why? To be fast, compilers **reorder and cache** memory operations, and CPUs have per-core caches
that don't instantly agree with each other. Within one goroutine, everything still *appears* to happen
in order. But **another goroutine can see the writes in a different order, late, or not at all**,
unless you synchronize.

And of course, `go run -race` reports it as a data race.

---

## 2. Happens-before

The [Go memory model](https://go.dev/ref/mem) answers one question: **when is a goroutine guaranteed
to see another goroutine's write?** The answer: only when there's a **happens-before** relationship
between them, created by synchronization. The main ones:

| Synchronization | Guarantee |
|---|---|
| `go f()` | everything before the `go` statement happens before `f` starts |
| channel send → receive | a send happens before the matching receive completes |
| `close(ch)` → a receive that sees it's closed | the close happens before that receive |
| `mu.Unlock()` → next `mu.Lock()` | everything before the unlock is visible after the lock |
| `wg.Done()` → `wg.Wait()` returns | the goroutine's work is visible after `Wait` |
| `once.Do(f)` | `f`'s writes are visible to everyone after `Do` returns |
| atomic operations | an atomic write is visible to an atomic read that sees its value |

Everything you learned in Part 3 works **because** of these rules. "One slot per goroutine, read after
`wg.Wait()`" (lesson 17) is correct because `Done` → `Wait` creates a happens-before edge.

The fix for section 1, with a channel:

```go
done := make(chan struct{})
go func() {
	message = "hello"
	close(done)          // happens before...
}()
<-done                   // ...this receive completes, so message is visible
fmt.Println(message)
```

The memory model's own advice: **"Don't be clever."** Use channels, mutexes, `sync.Once`, `WaitGroup`
and atomics, and run `-race`.

---

## 3. Atomics, properly

`sync/atomic` (lesson 21) gives single values that can be read and written safely without a mutex.
The typed versions are the easiest to use correctly:

```go
var stopped atomic.Bool
stopped.Store(true)
if stopped.Load() { ... }

var hits atomic.Int64
hits.Add(1)

old := counter.Swap(0)                       // set, returning the previous value
ok := state.CompareAndSwap(idle, running)    // set only if it currently equals idle
```

### `atomic.Pointer`: swapping a whole config at once

A very useful pattern: a configuration that's **read constantly** and **replaced occasionally** (a
"hot reload"). Readers mustn't see a half-updated config, and must not wait on a lock:

```go
type Config struct {
	RateLimit int
	Banned    map[string]bool
}

var current atomic.Pointer[Config]

func reload(c *Config) { current.Store(c) }      // publish a fully built config
func get() *Config      { return current.Load() } // always a complete config
```

The rule that makes it safe: **never modify a `Config` after storing it.** Build a new one and `Store`
it. Readers that already `Load`ed the old one keep using it safely.

---

## 4. Singleflight: deduplicating concurrent calls

Imagine 500 requests arrive at the same moment, all for the same uncached item. Without care, you send
500 identical queries to the database (a "thundering herd" or "cache stampede"). **Singleflight**
makes concurrent calls with the **same key** share **one** execution:

```go
v, err, shared := group.Do("user:42", func() (any, error) {
	return db.LoadUser(42)           // runs ONCE, even if 500 goroutines call Do("user:42") at once
})
```

The first caller runs the function. Callers that arrive while it's running **wait for it**, and all get
the same result. Once it finishes, the next call for that key runs the function again.

The standard implementation is `golang.org/x/sync/singleflight`. You'll build your own in Ex2 with a
mutex, a map of in-progress calls, and a `WaitGroup` per call.

---

## 5. Errgroup: goroutines that return errors

You've written "run tasks concurrently, cancel the rest on the first error" a few times now (lessons 22
and practice 08). `golang.org/x/sync/errgroup` packages it up:

```go
g, ctx := errgroup.WithContext(ctx)
for _, url := range urls {
	g.Go(func() error {
		return fetch(ctx, url)       // ctx is cancelled as soon as any task fails
	})
}
if err := g.Wait(); err != nil {     // waits for all; returns the first error
	return err
}
```

It's one of the most used concurrency helpers in Go. You'll write your own `Group` with the same API
in Ex4, from pieces you already know.

---

## 6. `sync.Cond`: waiting for a condition

Sometimes goroutines must wait until some **shared state** changes: "until the queue isn't empty",
"until the queue isn't full". A `sync.Cond` is tied to a mutex and lets goroutines sleep until another
goroutine signals that the state changed:

```go
type Queue struct {
	mu       sync.Mutex
	notEmpty *sync.Cond
	items    []int
}

func NewQueue() *Queue {
	q := &Queue{}
	q.notEmpty = sync.NewCond(&q.mu)
	return q
}

func (q *Queue) Put(v int) {
	q.mu.Lock()
	q.items = append(q.items, v)
	q.mu.Unlock()
	q.notEmpty.Signal()               // wake one waiting Get
}

func (q *Queue) Get() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	for len(q.items) == 0 {           // ⚠️ always a loop, never an if
		q.notEmpty.Wait()             // unlocks mu, sleeps, re-locks mu when woken
	}
	v := q.items[0]
	q.items = q.items[1:]
	return v
}
```

Key points:
- `Wait()` **unlocks** the mutex while sleeping and **re-locks** it before returning.
- Always check the condition in a **`for` loop**: when you wake up, another goroutine may already have
  taken the item.
- `Signal()` wakes one waiter; `Broadcast()` wakes all of them.

Honestly, a **buffered channel** is a simpler bounded queue, and most Go code uses one. `sync.Cond` is
for conditions more complex than a channel can express. You'll build one in Ex3 anyway, because
understanding it is part of understanding concurrency deeply.

### 🤔 Predict the output

```go
var counter int64
var wg sync.WaitGroup
for range 1000 {
	wg.Add(1)
	go func() {
		defer wg.Done()
		v := atomic.LoadInt64(&counter)
		atomic.StoreInt64(&counter, v+1)
	}()
}
wg.Wait()
fmt.Println(counter == 1000)
```

<details markdown="1">
<summary>Answer</summary>

Usually `false`. Each load and store is atomic, but **load-then-store is two steps**: two goroutines can
load the same value and both store `v+1`. The race detector won't complain (every access is atomic),
but the logic is still broken. Use `atomic.AddInt64(&counter, 1)`, which is one indivisible step, or
a `CompareAndSwap` loop. Atomics make single operations safe, not sequences of them.
</details>

---

## 7. Common mistakes

| You wrote | Problem | Fix |
|---|---|---|
| A plain `bool` flag shared between goroutines | may never be seen, or seen out of order | `atomic.Bool`, a channel, or a mutex |
| Load, then Store, with atomics | lost updates (not one step) | `Add`, `Swap`, or `CompareAndSwap` |
| Modifying a struct after publishing it through `atomic.Pointer` | readers see half-updated data | build a new one, then `Store` |
| `if` around `cond.Wait()` | wakes up when the condition is false again | `for` loop |
| `cond.Wait()` without holding the lock | panic or undefined behaviour | `Lock` first |
| "It works on my machine, so it's fine" | memory-model bugs depend on CPU and timing | `-race`, and only rely on happens-before |

---

## 8. Exercises

Run with: `go test -race -run Ex1 -v`. **The race detector is part of these tests.**

| | Exercise | Practises |
|---|---|---|
| Ex1 | Fix `Worker`'s racy stop flag | the memory model, `atomic.Bool` |
| Ex2 | `Group.Do`: your own singleflight | mutex + map + `WaitGroup` |
| Ex3 | `BoundedQueue` with `sync.Cond` | `Wait` in a loop, `Signal` / `Broadcast` |
| Ex4 | `ErrGroup`: your own errgroup | `WaitGroup`, `sync.Once`, context |
| Ex5 | `Settings` with hot reload | `atomic.Pointer` |

---

## ✅ Checklist

- [ ] I can explain why the "ready flag" program is broken
- [ ] I can name the operations that create happens-before edges
- [ ] I know atomics make single operations safe, not sequences
- [ ] I can explain singleflight and errgroup, and build simple versions
- [ ] I can use `sync.Cond` correctly, and know when a channel is simpler

**Next → [Lesson 35: Reflection and struct tags](../35_reflection/README.md)**
