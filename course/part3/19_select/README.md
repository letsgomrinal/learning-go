# Lesson 19: `select` and timeouts

**Goal:** wait on several channels at once, add timeouts to anything, do non-blocking sends and
receives, and stop a goroutine from the outside.

**Time:** about 1.5 hours.

```bash
cd ~/Desktop/learning-go/course/part3/19_select
```

---

## 1. Waiting on more than one channel

`<-ch` waits on **one** channel. What if results can arrive from two places, and you want whichever
comes first? That's what `select` is for:

```go
fast := make(chan string)
slow := make(chan string)

go func() { time.Sleep(50 * time.Millisecond); fast <- "fast server" }()
go func() { time.Sleep(200 * time.Millisecond); slow <- "slow server" }()

select {
case msg := <-fast:
	fmt.Println("got", msg)
case msg := <-slow:
	fmt.Println("got", msg)
}
// got fast server
```

`select` looks like a `switch`, but each `case` is a **channel operation**. It waits until **one** of
them can go ahead, runs that case, and is done. If several are ready at the same moment, it picks one
at random.

Python's closest equivalent is `asyncio.wait(..., return_when=FIRST_COMPLETED)`.

---

## 2. Timeouts with `time.After`

`time.After(d)` returns a channel that receives a value after the duration `d`. Put it in a `select`,
and you have a timeout:

```go
select {
case result := <-results:
	fmt.Println("result:", result)
case <-time.After(1 * time.Second):
	fmt.Println("gave up after 1 second")
}
```

Whichever happens first wins: the result arrives, or the second runs out. This is how you stop your
program from waiting forever on something slow, like a network call. (Lesson 22 shows `context`, the
standard way to do timeouts across a whole program. Underneath, it works the same way.)

### ⚠️ The goroutine leak trap

Here's a function that runs some work with a timeout:

```go
func withTimeout(work func() int, d time.Duration) (int, bool) {
	ch := make(chan int)            // ⚠️ unbuffered
	go func() {
		ch <- work()
	}()
	select {
	case v := <-ch:
		return v, true
	case <-time.After(d):
		return 0, false
	}
}
```

It looks correct. But think about the timeout case: the function returns, and **nobody will ever
receive from `ch`**. When `work()` finally finishes, the goroutine tries `ch <- ...` and blocks...
forever. That goroutine never ends. It's a **goroutine leak**: it holds memory for as long as the
program runs. Do this on every web request, and the server slowly runs out of memory.

The fix is one character:

```go
ch := make(chan int, 1)            // buffered: the send always succeeds, even with no receiver
```

With a buffer of 1, the goroutine can drop its result in the buffer and exit, whether or not anyone
is still listening.

**For every goroutine you start, ask: "how does this goroutine end?"**

---

## 3. Non-blocking operations with `default`

A `select` with a `default` case never waits. If no channel is ready right now, it runs `default`
instead:

```go
select {
case msg := <-messages:
	fmt.Println("received", msg)
default:
	fmt.Println("no message waiting")
}
```

Same for sending:

```go
select {
case jobs <- job:
	fmt.Println("queued")
default:
	fmt.Println("queue full, dropping the job")
}
```

Use it when you want to "try" without waiting. Don't put a `select` with `default` in a tight loop
to "poll" a channel: that burns a CPU core doing nothing. Just block on the channel instead.

---

## 4. Repeating things with `time.Ticker`

A **ticker** sends on its channel at a regular interval:

```go
ticker := time.NewTicker(500 * time.Millisecond)
defer ticker.Stop()                // always stop tickers when you're done

timeout := time.After(2 * time.Second)
for {
	select {
	case <-ticker.C:
		fmt.Println("tick")
	case <-timeout:
		fmt.Println("done")
		return
	}
}
```

This prints `tick` about four times, then `done`. `for` + `select` is a very common shape: a loop
that reacts to whichever event happens next.

(Important detail: create the `time.After` **before** the loop. Inside the loop, it would create a
new 2-second timer on every pass, and never fire.)

---

## 5. Stopping a goroutine: the "done" channel

How do you tell a running goroutine to stop? You can't kill it from outside. It has to **cooperate**:
it checks a signal, and returns. The classic signal is a channel that gets **closed**:

```go
func counter(done <-chan struct{}) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for i := 0; ; i++ {
			select {
			case out <- i:          // send the next number...
			case <-done:            // ...unless we've been told to stop
				return
			}
		}
	}()
	return out
}

func main() {
	done := make(chan struct{})
	nums := counter(done)
	fmt.Println(<-nums, <-nums, <-nums)   // 0 1 2
	close(done)                           // tell the goroutine to stop
}
```

Why does this work? **A receive from a closed channel never blocks** (it returns immediately with the
zero value). So once `done` is closed, `case <-done` is always ready, and the goroutine returns.
Closing is also a **broadcast**: every goroutine waiting on `done` sees it, however many there are.

`chan struct{}` is the conventional type for pure signals: `struct{}` carries no data and takes no memory.

Lesson 22's `context` is this same idea, standardised: `ctx.Done()` is a done channel.

### 🤔 Predict the output

```go
ch := make(chan int, 1)
for i := range 3 {
	select {
	case ch <- i:
		fmt.Print("sent ", i, "; ")
	default:
		fmt.Print("full ", i, "; ")
	}
}
fmt.Println(<-ch)
```

<details markdown="1">
<summary>Answer</summary>

`sent 0; full 1; full 2; 0`. The buffer has room for one value. `0` fits; `1` and `2` find it full,
so `default` runs. Then we receive the `0`.
</details>

---

## 6. Common mistakes

| You wrote | What happens | Fix |
|---|---|---|
| Unbuffered result channel + timeout | the goroutine leaks when the timeout wins | `make(chan T, 1)` |
| `time.After` created inside a `for`/`select` loop as an overall timeout | it restarts every loop, never fires | create it once, before the loop |
| `select` with `default` in a loop, to wait for a value | 100% CPU doing nothing | remove `default` and block |
| Forgot `ticker.Stop()` | the ticker keeps running | `defer ticker.Stop()` |
| A goroutine with no way to be stopped | leaks when nobody needs it anymore | give it a done channel (or a context) |
| Empty `select {}` by accident | blocks forever | (only use it on purpose) |

---

## 7. Exercises

Run with: `go test -race -run Ex1 -v`

| | Exercise | Practises |
|---|---|---|
| Ex1 | `First`: whichever channel answers first | `select` |
| Ex2 | `ReceiveWithTimeout` | `time.After` |
| Ex3 | `TrySend` | non-blocking send with `default` |
| Ex4 | `TryReceive` | non-blocking receive with `default` |
| Ex5 | `RunWithTimeout`, without leaking | the goroutine leak trap |
| Ex6 | `CountTicks` | `time.Ticker` + `for`/`select` |
| Ex7 | `Counter` that stops when `done` is closed | the done channel |

---

## ✅ Checklist

- [ ] I can wait on several channels with `select`
- [ ] I can add a timeout to any channel receive with `time.After`
- [ ] I can explain the goroutine leak in `withTimeout`, and fix it with a buffer
- [ ] I can do a non-blocking send or receive with `default`
- [ ] I can stop a goroutine by closing a done channel

**Next → [Lesson 20: Pipelines and worker pools](../20_patterns/README.md)**
