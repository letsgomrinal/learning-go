# Lesson 18: Channels

**Goal:** send values between goroutines with channels, understand when sending and receiving block,
close channels properly, and avoid deadlocks.

**Time:** 1.5 hours.

```bash
cd ~/Desktop/learning-go/course/part3/18_channels
```

---

## 1. What a channel is

A **channel** is a pipe that connects goroutines. One goroutine puts a value in, another takes it out.
It's the main way goroutines talk to each other in Go.

Go's motto for this:

> **Don't communicate by sharing memory; share memory by communicating.**

Instead of several goroutines changing the same variable (the data race from lesson 17), one
goroutine **sends** its result through a channel, and another **receives** it.

Python's closest equivalent is `queue.Queue` (for threads) or `asyncio.Queue`.

---

## 2. Creating, sending, receiving

```go
ch := make(chan string)       // a channel that carries strings

go func() {
	ch <- "hello"             // SEND: put "hello" into the channel
}()

msg := <-ch                   // RECEIVE: take a value out of the channel
fmt.Println(msg)              // hello
```

The arrow `<-` shows which way the data flows:
- `ch <- v`: the value goes **into** the channel (send).
- `v := <-ch`: the value comes **out of** the channel (receive).

A channel carries one type only: `chan int`, `chan string`, `chan Result`...

---

## 3. Unbuffered channels block

This is the key idea. A channel made with `make(chan T)` is **unbuffered**, meaning it can't hold
anything:

- A **send** waits until another goroutine is ready to **receive**.
- A **receive** waits until another goroutine **sends**.

The two goroutines meet at the channel and hand the value over directly, like passing a baton in a
relay race. Nobody continues until the handover has happened.

That's why receiving also works as **waiting**. In the example above, `msg := <-ch` doesn't continue
until the goroutine has sent, so you don't need a `WaitGroup` there.

### Deadlock

What if you send, and nobody ever receives?

```go
func main() {
	ch := make(chan int)
	ch <- 1                 // waits for a receiver... forever
	fmt.Println("never")
}
```

```
fatal error: all goroutines are asleep - deadlock!
```

Go noticed that **every** goroutine is stuck waiting, so nothing can ever happen, and stopped the
program. You'll see this message often while learning. It usually means: a send with no receiver, a
receive with no sender, or a forgotten `wg.Done()`.

---

## 4. Buffered channels

You can give a channel a **buffer**, a number of slots:

```go
ch := make(chan int, 3)      // room for 3 values

ch <- 1                       // doesn't block: there's room
ch <- 2
ch <- 3
// ch <- 4 would block now: the buffer is full

fmt.Println(<-ch, <-ch, <-ch)   // 1 2 3  (first in, first out)
fmt.Println(len(ch), cap(ch))   // 0 3
```

- A send only blocks when the buffer is **full**.
- A receive only blocks when the buffer is **empty**.

When to use which: start with **unbuffered** channels; they're the easiest to reason about. Use a
buffer when you know how many values will be sent and want senders not to wait (you'll see an
important case in lesson 19), or to smooth out bursts.

---

## 5. Closing a channel

When a sender has no more values, it **closes** the channel:

```go
close(ch)
```

Closing says "nothing more is coming". Receivers can detect it:

```go
v, ok := <-ch       // ok is false when the channel is closed AND empty
```

That's the "comma ok" pattern again. After a channel is closed, receives return the remaining
buffered values, then the zero value with `ok == false` forever.

### `range` over a channel

The most common way to receive: loop until the channel is closed.

```go
func numbers(n int) <-chan int {
	ch := make(chan int)
	go func() {
		defer close(ch)            // close when we're done sending
		for i := 1; i <= n; i++ {
			ch <- i
		}
	}()
	return ch
}

func main() {
	for v := range numbers(3) {    // receives until ch is closed
		fmt.Println(v)             // 1, 2, 3
	}
}
```

It looks like a Python generator: values arrive one at a time, and the loop ends when the producer
is done.

If the sender **forgets** to close, the `range` waits forever for the next value: deadlock.

### The closing rules

1. **Only the sender closes**, never the receiver. (The receiver can't know whether more sends are coming.)
2. **Close only once.** Closing twice panics: `close of closed channel`.
3. **Never send on a closed channel**: `panic: send on closed channel`.
4. You **don't have to** close every channel. Close it when the receiver needs to know "no more
   values", for example because it uses `range`.

### 🤔 Predict the output

```go
ch := make(chan int, 5)
for i := range 3 {
	ch <- i * 10
}
close(ch)
sum := 0
for v := range ch {
	sum += v
}
v, ok := <-ch
fmt.Println(sum, v, ok)
```

<details markdown="1">
<summary>Answer</summary>

`30 0 false`. The buffer holds 0, 10 and 20, and `range` still receives them after `close`. Once
they're used up, a receive gives the zero value `0` and `ok == false`.
</details>

---

## 6. Direction in the type

You can restrict a channel to one direction in a function's signature:

```go
func produce(out chan<- int) { ... }   // can only SEND to out
func consume(in <-chan int) { ... }    // can only RECEIVE from in
```

Read the arrow: `chan<-` points **into** the channel (send-only); `<-chan` points **out** (receive-only).

A function that starts a goroutine and returns a channel usually returns `<-chan T`, like `numbers`
above. The caller can only receive from it, so it can't accidentally send or close. A normal
`chan T` converts to either direction automatically.

---

## 7. Collecting results with a channel

In lesson 17 you collected results with one slice slot per goroutine. A channel is the other common way:

```go
func sumOfSquares(nums []int) int {
	results := make(chan int)
	for _, n := range nums {
		go func() {
			results <- n * n          // each goroutine sends its result
		}()
	}

	total := 0
	for range nums {                  // receive exactly len(nums) values
		total += <-results
	}
	return total
}
```

No `WaitGroup` and no data race: only `sumOfSquares` touches `total`, and the goroutines only send.
Receiving exactly `len(nums)` times means we wait for every goroutine.

Results arrive in the order the goroutines **finish**, not the order they started. If you need the
original order, use the slice-slot pattern from lesson 17, or send the index along with the value
(lesson 20).

---

## 8. Common mistakes

| You wrote | What happens | Fix |
|---|---|---|
| `ch <- v` in `main` on an unbuffered channel, with no receiver | deadlock | receive in another goroutine, or add a buffer |
| `for v := range ch` and the sender never closes | deadlock (or a goroutine stuck forever) | `defer close(ch)` in the sender |
| Receiver calls `close(ch)` | the sender may panic with `send on closed channel` | only the sender closes |
| `var ch chan int` then using it | sends and receives on a nil channel block forever | `ch := make(chan int)` |
| Expecting results in start order | they arrive in finish order | slice slots, or send the index too |

---

## 9. Exercises

Run with: `go test -race -run Ex1 -v`. The tests use timeouts, so a mistake shows up as a clear
failure instead of hanging forever.

| | Exercise | Practises |
|---|---|---|
| Ex1 | `Greeting`: one value from a goroutine | send and receive |
| Ex2 | `Numbers`: 1..n, then close | producer + `close` |
| Ex3 | `Sum` everything from a channel | `range` over a channel |
| Ex4 | `SumOfSquares` with one goroutine per number | collecting results with a channel |
| Ex5 | `Doubler`: read from one channel, write to another | receive-only and send-only channels |
| Ex6 | `Queue`: a pre-filled buffered channel | buffered channels without goroutines |

---

## ✅ Checklist

- [ ] I can make a channel, send with `ch <- v` and receive with `<-ch`
- [ ] I know an unbuffered send waits for a receiver, and what "deadlock" means
- [ ] I know when sends and receives block on a buffered channel
- [ ] I close channels from the sender only, and receive with `range` or `v, ok :=`
- [ ] I can read `chan<- T` and `<-chan T`

**Next → [Lesson 19: select and timeouts](../19_select/README.md)**
