# Lesson 22: Context

**Goal:** cancel work that nobody needs anymore, put deadlines on slow operations, and pass a
`context.Context` through your code the way real Go programs do.

**Time:** about 1.5 hours.

```bash
cd ~/Desktop/learning-go/course/part3/22_context
```

---

## 1. The problem

A user opens a web page. Your server starts three things for that request: a database query, a call
to another service, and a cache lookup. Then the user closes the tab. All three are still running,
using resources for a page nobody will see.

Or: one of the three fails, so the page can't be built anyway. The other two keep running for nothing.

Or: the other service is having a bad day and takes 30 seconds. You'd rather give up after 2.

In lesson 19 you stopped a goroutine by closing a "done" channel. `context.Context` is the standard,
built-in version of that idea, used by the whole standard library and almost every Go package.

---

## 2. What a context is

A `context.Context` (package `context`) carries:
- a **cancellation signal**: a `Done()` channel that gets closed when the work should stop;
- an optional **deadline**: a time after which it cancels itself;
- the **reason** it stopped: `Err()`.

You don't create contexts from nothing. You start from a root, and derive children from it:

```go
ctx := context.Background()     // the root: never cancelled. Use it in main and in tests.
```

### Cancelling by hand: `WithCancel`

```go
ctx, cancel := context.WithCancel(context.Background())

go func() {
	<-ctx.Done()                  // waits until the context is cancelled
	fmt.Println("stopped because:", ctx.Err())
}()

time.Sleep(50 * time.Millisecond)
cancel()                          // closes ctx.Done()
time.Sleep(10 * time.Millisecond)
// stopped because: context canceled
```

`WithCancel` returns a new context **and** a `cancel` function. Calling `cancel()` closes the
context's `Done()` channel, and every goroutine watching it gets the signal.

### Cancelling by time: `WithTimeout`

```go
ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
defer cancel()

// ... the context cancels itself after 2 seconds ...
```

When the time runs out, `Done()` closes and `ctx.Err()` returns `context.DeadlineExceeded`.

(There's also `WithDeadline(parent, someTime)`, the same thing with a clock time instead of a duration.)

### Always call `cancel`

Even with `WithTimeout`, **always `defer cancel()`** right after creating the context. It releases the
timer and resources as soon as you're done, instead of waiting for the timeout. `go vet` warns if you
throw it away (`the cancel function returned by context.WithTimeout should be called, not discarded`).

---

## 3. Contexts form a tree

Each derived context is a **child** of its parent. **Cancelling a parent cancels all its children**,
and their children, all the way down:

```
Background
 └── request (WithCancel)            ← user closes the tab: cancel()
      ├── db query (WithTimeout 2s)        cancelled too
      └── api call (WithTimeout 5s)        cancelled too
```

Cancelling a child does **not** affect its parent or siblings. A child's timeout can only make the
deadline **earlier**, never later than the parent's.

---

## 4. Writing functions that respect a context

Cancellation is **cooperative**: nothing gets killed for you. Your code has to check `ctx.Done()`,
usually in a `select`, and return when it's closed:

```go
func slowSquare(ctx context.Context, n int) (int, error) {
	select {
	case <-time.After(100 * time.Millisecond):    // pretend this is slow work
		return n * n, nil
	case <-ctx.Done():
		return 0, ctx.Err()                       // give up: context canceled / deadline exceeded
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	_, err := slowSquare(ctx, 4)
	fmt.Println(err)                              // context deadline exceeded
}
```

For loops that do work in steps, check between steps:

```go
for _, item := range items {
	if err := ctx.Err(); err != nil {     // non-nil once cancelled
		return err
	}
	process(item)
}
```

Check for a specific reason with `errors.Is` (lesson 13):

```go
if errors.Is(err, context.DeadlineExceeded) {
	fmt.Println("too slow")
}
```

---

## 5. The rules for using context

These are conventions every Go programmer follows:

1. **`ctx` is the first parameter**, named `ctx`: `func Fetch(ctx context.Context, url string) error`.
2. **Pass it down** to everything that might be slow: your own functions, database calls,
   HTTP requests. The standard library accepts it: `http.NewRequestWithContext(ctx, ...)`,
   `db.QueryContext(ctx, ...)`.
3. **Don't store it in a struct.** Pass it as a parameter to each call.
4. **Never pass `nil`.** If you don't have one yet, use `context.Background()` (or `context.TODO()` as
   a reminder to add a real one later).
5. **`defer cancel()`** right after every `WithCancel` / `WithTimeout`.

You'll also see `context.WithValue` for carrying request-scoped data, like a request ID. Use it
sparingly: only for data that crosses API boundaries, never for normal function parameters.

---

## 6. Cancelling siblings when one fails

A very common need: run several tasks at once; if **any** fails, cancel the rest and return the error.

```go
func fetchAll(ctx context.Context, urls []string) error {
	ctx, cancel := context.WithCancel(ctx)   // a child we can cancel
	defer cancel()

	errs := make(chan error, len(urls))      // buffered: no goroutine gets stuck sending
	for _, url := range urls {
		go func() {
			errs <- fetch(ctx, url)          // each fetch watches ctx
		}()
	}

	for range urls {
		if err := <-errs; err != nil {
			cancel()                         // tell the others to stop
			return err
		}
	}
	return nil
}
```

This is exactly what the popular package `golang.org/x/sync/errgroup` does. You'll write your own
version in the exercises, which makes `errgroup` easy to understand when you meet it.

### Ctrl+C for command-line programs

```go
ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
defer stop()
// ctx is cancelled when the user presses Ctrl+C. Pass it to everything.
```

### 🤔 Predict the output

```go
parent, cancelParent := context.WithCancel(context.Background())
child, cancelChild := context.WithTimeout(parent, time.Hour)
defer cancelChild()

cancelParent()
<-child.Done()
fmt.Println(child.Err(), parent.Err())
```

<details markdown="1">
<summary>Answer</summary>

`context canceled context canceled`. The child had an hour left on its own timeout, but cancelling
the parent cancels the child immediately. Its reason is "canceled", not "deadline exceeded".
</details>

---

## 7. Common mistakes

| You wrote | What happens | Fix |
|---|---|---|
| `ctx, _ := context.WithTimeout(...)` | the timer leaks until it fires; `go vet` warns | `ctx, cancel := ...; defer cancel()` |
| A slow function that never checks `ctx.Done()` | cancelling does nothing | `select` on `ctx.Done()`, or check `ctx.Err()` between steps |
| `ctx` stored in a struct field | hard to follow which request it belongs to | pass it as the first parameter |
| Passing `nil` as a context | panics deep inside a library | `context.Background()` |
| Creating a new `Background()` deep inside your code | loses the caller's cancellation | pass the caller's `ctx` down |
| Unbuffered error channel + early return | the remaining goroutines block forever | buffer it with `len(tasks)` |

---

## 8. Exercises

Run with: `go test -race -run Ex1 -v`

| | Exercise | Practises |
|---|---|---|
| Ex1 | `Sleep` that can be cancelled | `select` on `ctx.Done()` |
| Ex2 | `SumUntilCancelled`: stop between steps | `ctx.Err()` in a loop |
| Ex3 | `Ticker`: a producer that stops on cancel | a goroutine watching a context |
| Ex4 | `WithTimeout`: run work with a deadline | `context.WithTimeout`, `defer cancel()` |
| Ex5 | `FirstResult`: take the fastest, cancel the rest | racing goroutines + cancel |
| Ex6 | `FetchAll`: first error cancels everything | a mini `errgroup` |

---

## ✅ Checklist: end of Part 3 🎉

- [ ] I can create contexts with `WithCancel` and `WithTimeout`, and always `defer cancel()`
- [ ] I know cancelling a parent cancels its children
- [ ] I can write a function that stops when its context is cancelled
- [ ] I follow the rules: `ctx` first, passed down, never stored, never nil
- [ ] I can run tasks concurrently and cancel the rest when one fails

---

## What's next

You've finished Part 3. Concurrency is where Go is most different from Python, so it's worth
consolidating before moving on:

1. **Practise.** Do the practice sets, in order:
   - `practice/07_concurrency`
   - `practice/08_sync_context`
2. **Build something.** The Part 3 milestone in [roadmap.html](../../../roadmap.html): a concurrent web
   crawler with a worker pool, a protected "seen" set, a global timeout via `context`, and graceful Ctrl+C.
3. **Read.** Now is a great time to start *"100 Go Mistakes and How to Avoid Them"*. Many of the
   mistakes will make sense now.

Then Part 4: building real programs (files, HTTP servers, JSON, databases).

**Next → [Part 4, Lesson 23: Files and I/O](../../part4/23_files/README.md)**
