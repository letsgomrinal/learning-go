# Lesson 20: Pipelines and worker pools

**Goal:** combine goroutines and channels into the patterns real Go programs use: pipelines, fan-in,
worker pools, and limiting how much runs at once.

**Time:** 1.5–2 hours. This lesson puts lessons 17–19 together.

```bash
cd ~/Desktop/learning-go/course/part3/20_patterns
```

---

## 1. Pipelines

A **pipeline** is a chain of **stages** connected by channels. Each stage is a goroutine that receives
values from the previous stage, does something, and sends results to the next.

```
generate ──▶ square ──▶ sum
 1,2,3        1,4,9      14
```

You built the pieces in lesson 18. Here they are together:

```go
// Stage 1: produce values, then close.
func generate(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, n := range nums {
			out <- n
		}
	}()
	return out
}

// Stage 2: transform values; close when the input is closed.
func square(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			out <- n * n
		}
	}()
	return out
}

func main() {
	total := 0
	for v := range square(generate(1, 2, 3)) {   // stage 3: consume
		total += v
	}
	fmt.Println(total)   // 14
}
```

Every stage follows the same recipe:
1. Make an output channel.
2. Start a goroutine that reads from the input, writes to the output, and **closes the output when done** (`defer close(out)`).
3. Return the output channel straight away.

Closing propagates down the pipeline: `generate` closes, so `square`'s `range` ends and it closes,
so the final `range` ends. Stages compose like Unix pipes: `square(square(generate(2)))` works.

All the stages run **at the same time**: while `square` works on item 2, `generate` can already be
sending item 3. Python comparison: chained generators, except each stage runs concurrently.

---

## 2. Fan-out and fan-in

**Fan-out:** several goroutines read from the **same** channel to share the work. A channel can have
many receivers, and each value goes to exactly one of them.

**Fan-in:** merge several channels into **one**:

```go
func merge(chans ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup

	for _, ch := range chans {
		wg.Add(1)
		go func() {                     // one forwarder per input channel
			defer wg.Done()
			for v := range ch {
				out <- v
			}
		}()
	}

	go func() {                         // close out once ALL forwarders are done
		wg.Wait()
		close(out)
	}()

	return out
}
```

The interesting part is **when to close `out`**. You can't close it after the loop, because the
forwarders are still sending (sending on a closed channel panics). And you can't call `wg.Wait()`
directly in `merge`, because `merge` must return `out` right away so someone can start receiving
(otherwise the forwarders block forever on `out <- v`). The answer is one extra goroutine that waits,
then closes. You'll see this "wait then close in a goroutine" trick a lot.

---

## 3. Worker pools

Starting one goroutine per item (lesson 17) is fine for 10 items. For a million items, or when each
item hits a database that can only handle 5 connections, you want a **fixed number of workers**:

```
             ┌─ worker 1 ─┐
jobs ───────▶├─ worker 2 ─┤──────▶ results
             └─ worker 3 ─┘
```

```go
func process(items []string, workers int) []string {
	jobs := make(chan int)                    // job = an index into items
	results := make([]string, len(items))
	var wg sync.WaitGroup

	// Start the workers.
	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range jobs {             // each worker takes jobs until jobs is closed
				results[i] = strings.ToUpper(items[i])
			}
		}()
	}

	// Send the jobs.
	for i := range items {
		jobs <- i
	}
	close(jobs)                               // no more jobs: the workers' range loops end

	wg.Wait()                                 // wait for the last jobs to finish
	return results
}
```

The pieces:
- **`jobs` channel:** the work queue. Sending an **index** (not the item) lets each worker write its
  result into the right slot, so the results keep their original order with no locks (lesson 17).
- **N workers:** each loops over `jobs`. Whichever worker is free takes the next job.
- **`close(jobs)`:** tells the workers there's no more work.
- **`wg.Wait()`:** makes sure every job has finished before we return.

Python comparison: `concurrent.futures.ThreadPoolExecutor(max_workers=3).map(...)`. In Go you
build it in about 20 lines, and can see exactly how it works.

### Results through a channel

If you'd rather send results on a channel, include the index so you can put them back in order:

```go
type result struct {
	index int
	value string
}
```

---

## 4. Limiting concurrency with a semaphore

Sometimes you're happy to start a goroutine per item, but want **at most N running at once** (for
example, at most 3 simultaneous downloads). A **buffered channel** makes a neat limiter, called a
**semaphore**:

```go
sem := make(chan struct{}, 3)        // 3 slots

for _, url := range urls {
	wg.Add(1)
	go func() {
		defer wg.Done()
		sem <- struct{}{}            // take a slot (blocks if all 3 are taken)
		defer func() { <-sem }()     // give the slot back when done
		download(url)
	}()
}
wg.Wait()
```

Only 3 goroutines can hold a slot at a time; the rest wait at `sem <- struct{}{}` until a slot frees up.

### 🤔 Predict the output

```go
results := make(chan int)
var wg sync.WaitGroup
for i := 1; i <= 3; i++ {
	wg.Add(1)
	go func() {
		defer wg.Done()
		results <- i * 10
	}()
}
go func() {
	wg.Wait()
	close(results)
}()
sum := 0
for v := range results {
	sum += v
}
fmt.Println(sum)
```

<details markdown="1">
<summary>Answer</summary>

`60`. The values 10, 20 and 30 arrive in some order, but the sum is always 60. The "wait then close
in a goroutine" trick makes the `range` end once all three have sent.
</details>

---

## 5. Choosing a pattern

| Situation | Pattern |
|---|---|
| A handful of independent tasks | one goroutine each + `WaitGroup` (lesson 17) |
| Many tasks, but limited resources (DB connections, API rate limits) | worker pool, or semaphore |
| Data flowing through several processing steps | pipeline |
| Several sources of values, one consumer | fan-in (`merge`) |

---

## 6. Common mistakes

| You wrote | What happens | Fix |
|---|---|---|
| A stage that forgets `defer close(out)` | the next stage's `range` never ends: deadlock | always close your output |
| `close(out)` right after starting the forwarders in `merge` | `panic: send on closed channel` | wait, then close, in a separate goroutine |
| `wg.Wait()` before the code that receives from the results | deadlock: senders wait for a receiver, the receiver waits for senders | wait in a goroutine, or receive first |
| Forgot `close(jobs)` | workers wait for more jobs forever, `wg.Wait()` never returns | close the jobs channel after sending |
| Worker pool with results by arrival order when order matters | shuffled results | send indexes, write `results[i]` |

---

## 7. Exercises

Run with: `go test -race -run Ex1 -v`

| | Exercise | Practises |
|---|---|---|
| Ex1 | `Generate`: the first stage | producer stage |
| Ex2 | `Square`: a middle stage | transform stage |
| Ex3 | `SumSquares`: connect the pipeline | composing stages |
| Ex4 | `Merge` channels into one | fan-in, wait-then-close |
| Ex5 | `WorkerPool` that keeps the order | worker pool with job indexes |
| Ex6 | `RunLimited`: at most N at once | semaphore with a buffered channel |

---

## ✅ Checklist

- [ ] I can write a pipeline stage, and I always close its output channel
- [ ] I can merge channels, and explain why `close(out)` happens in its own goroutine
- [ ] I can write a worker pool from memory, and keep results in order
- [ ] I can limit concurrency with a buffered-channel semaphore
- [ ] I can pick the right pattern for a situation

**Next → [Lesson 21: Mutexes and the race detector](../21_mutex/README.md)**
