# Lesson 33: The scheduler, leaks and deadlocks

**Goal:** understand how Go runs millions of goroutines on a few OS threads, choose the right number of
workers for CPU-bound and I/O-bound work, and find and fix goroutine leaks and deadlocks.

**Time:** about 2 hours.

```bash
cd ~/Desktop/learning-go/course/part5/33_scheduler
```

---

## 1. Goroutines are not threads

Operating system **threads** are heavy: each has a large fixed stack (often 1–8 MB), and switching
between them goes through the OS kernel. A few thousand threads is a lot.

**Goroutines** are managed by Go's own **scheduler**, inside your program:
- they start with a tiny stack (a few KB) that grows and shrinks as needed;
- switching between them doesn't involve the kernel, so it's very cheap;
- a million goroutines is fine on a laptop.

The scheduler runs your goroutines **on top of** a small number of OS threads.

---

## 2. The G-M-P model

The scheduler is built from three kinds of things:

| Letter | Name | What it is |
|---|---|---|
| **G** | goroutine | your function, with its stack and state |
| **M** | machine | an OS thread |
| **P** | processor | a "permission to run Go code", with its own queue of ready goroutines |

```
   P0                 P1                 P2                 P3
 [G G G]            [G G]              [G G G G]          [G]       ← local run queues
    │                  │                  │                  │
    M0                 M1                 M2                 M3     ← OS threads
    │                  │                  │                  │
  CPU 0              CPU 1              CPU 2              CPU 3
```

- There are **`GOMAXPROCS`** Ps. By default that's the number of CPU cores (Go 1.25+ also respects
  container CPU limits). So at most `GOMAXPROCS` goroutines run Go code **at the same instant**.
- Each P has a queue of runnable goroutines. An M (thread) picks a P and runs goroutines from its queue.
- **Work stealing:** when a P's queue is empty, it steals half of another P's queue. That keeps all
  cores busy without a central bottleneck.

You can read and set it:

```go
fmt.Println(runtime.GOMAXPROCS(0))   // 0 means "just tell me": e.g. 10 on a 10-core Mac
fmt.Println(runtime.NumCPU())
fmt.Println(runtime.NumGoroutine())  // how many goroutines exist right now
```

---

## 3. What happens when a goroutine blocks

This is where Go shines compared with Python's threads and asyncio:

- **Waiting on a channel, mutex, `time.Sleep`, or a network socket:** the goroutine is parked, and the
  P immediately runs another goroutine. No thread is wasted. Network I/O uses the OS's event system
  (like epoll/kqueue) under the hood, so 10,000 connections waiting for data cost almost nothing.
  This is what asyncio does, except you write normal blocking code with no `async`/`await`.
- **A blocking system call (a file read, cgo):** the thread itself is stuck in the kernel, so the
  scheduler hands the P to **another thread** and keeps running other goroutines. When the call
  returns, the goroutine gets back in a queue.
- **A long-running loop with no function calls:** since Go 1.14 the scheduler **preempts** it (pauses it
  from the outside) after about 10ms, so one busy goroutine can't starve the others.

---

## 4. How many workers?

The answer depends on what the work waits for:

| Work type | Example | Workers |
|---|---|---|
| **CPU-bound** | resizing images, hashing, compressing | about `runtime.GOMAXPROCS(0)`. More doesn't help: there are only that many cores |
| **I/O-bound** | HTTP calls, database queries | many more (tens to thousands), limited by what the **other side** can handle |

More goroutines than cores for CPU-bound work just adds switching overhead. For I/O-bound work, the
limit is usually the remote service, a rate limit, or the database's connection pool.

---

## 5. Goroutine leaks

A **leak** is a goroutine that never ends. It holds its stack and everything it references, forever.
Leak a few per request, and a server's memory grows until it crashes. The usual causes:

1. **Sending on a channel nobody will ever read** (lesson 19's timeout trap).
2. **Receiving from a channel nobody will ever send to or close.**
3. **A loop with no way to stop** (no done channel, no context).

### Example: "first result wins"

```go
func firstMatch(items []string, match func(string) bool) string {
	results := make(chan string)               // ⚠️ unbuffered
	for _, it := range items {
		go func() {
			if match(it) {
				results <- it                  // only the first send is received...
			}
		}()
	}
	return <-results                           // ...the other matching goroutines block forever
}
```

Every goroutine after the first match is stuck on `results <- it` forever. The fixes: a **buffered**
channel with room for everyone (`make(chan string, len(items))`), or a **context** that tells the
others to give up.

### Finding leaks

- In tests, compare `runtime.NumGoroutine()` before and after (allowing a moment for goroutines to
  finish). This lesson's tests do that.
- In a running server, the goroutine profile (lesson 32) lists every goroutine and **where it's
  stuck**: `curl "localhost:6060/debug/pprof/goroutine?debug=1"`. Thousands of goroutines stuck on the
  same line is a leak.
- The package `go.uber.org/goleak` checks for leaks automatically at the end of tests.

---

## 6. Deadlocks

A **deadlock** is when goroutines wait for each other in a circle, and none can continue.

If **every** goroutine is stuck, Go notices and crashes with the message you know:

```
fatal error: all goroutines are asleep - deadlock!
```

But if **some** goroutines are still running (an HTTP server always has some), Go **can't** detect it.
The stuck goroutines just hang silently, and requests time out. That's the dangerous kind.

### The classic: lock ordering

```go
func transfer(from, to *Account, amount int) {
	from.mu.Lock()
	defer from.mu.Unlock()
	to.mu.Lock()                 // ⚠️
	defer to.mu.Unlock()
	from.Balance -= amount
	to.Balance += amount
}
```

Goroutine 1 runs `transfer(a, b, 10)`: it locks `a`, then wants `b`.
Goroutine 2 runs `transfer(b, a, 5)`: it locks `b`, then wants `a`.
Each holds what the other needs. Both wait forever.

**The fix: always take locks in the same global order**, for example by account ID:

```go
first, second := from, to
if to.ID < from.ID {
	first, second = to, from
}
first.mu.Lock()
defer first.mu.Unlock()
second.mu.Lock()
defer second.mu.Unlock()
```

Now both goroutines lock the lower ID first, so they can never hold one lock each.

Other deadlock causes: locking the same mutex twice in one goroutine (Go's mutexes aren't re-entrant,
lesson 21), or `wg.Wait()` before the code that would let the goroutines finish.

### 🤔 Predict the output

```go
func main() {
	ch := make(chan int)
	go func() {
		for v := range ch {
			fmt.Println(v)
		}
	}()
	ch <- 1
	ch <- 2
	time.Sleep(10 * time.Millisecond)
	fmt.Println("goroutines:", runtime.NumGoroutine())
}
```

<details><summary>Answer</summary>

```
1
2
goroutines: 2
```

`main` plus the receiver goroutine, which is still waiting in `range ch` because nobody closed the
channel. In a long-running program, that's a leak. `close(ch)` after the sends lets it finish.
</details>

---

## 7. Common mistakes

| You did | Problem | Fix |
|---|---|---|
| 1000 workers for CPU-bound work | overhead, no speed-up | `runtime.GOMAXPROCS(0)` workers |
| Unbuffered result channel, only the first result read | the rest of the goroutines leak | buffer it, or cancel with a context |
| A background loop with no way to stop | leaks when its owner is done | done channel or context, plus a `Stop` function |
| Two locks taken in different orders | deadlock under load | one global lock order |
| Relying on "all goroutines are asleep" to catch deadlocks | it never fires in a server | timeouts in tests, goroutine dumps in production |

---

## 8. Exercises

Run with: `go test -race -run Ex1 -v`. The tests use timeouts and goroutine counts, so leaks and
deadlocks show up as clear failures.

| | Exercise | Practises |
|---|---|---|
| Ex1 | Fix the leak in `FirstMatch` | goroutine leaks |
| Ex2 | `StartHeartbeat` that can be stopped | goroutine lifetime, stop functions |
| Ex3 | Fix the deadlock in `Transfer` | lock ordering |
| Ex4 | `CountPrimes` using all CPU cores | CPU-bound worker counts, `GOMAXPROCS` |

---

## ✅ Checklist

- [ ] I can explain G, M and P, and what `GOMAXPROCS` controls
- [ ] I know what happens to the thread when a goroutine waits on a channel, a socket, or a system call
- [ ] I choose worker counts differently for CPU-bound and I/O-bound work
- [ ] I can find a goroutine leak (`NumGoroutine`, goroutine profile) and fix it
- [ ] I can explain the lock-ordering deadlock and fix it with a global order

**Next → [Lesson 34: The memory model and advanced sync](../34_memory_model/README.md)**
