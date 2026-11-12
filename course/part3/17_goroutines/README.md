# Lesson 17: Goroutines and WaitGroup

**Goal:** run functions at the same time with goroutines, wait for them to finish with
`sync.WaitGroup`, and collect their results safely.

**Time:** 1.5 hours. This is the start of Go's best-known feature, so go slowly.

```bash
cd ~/Desktop/learning-go/course/part3/17_goroutines
```

**From now on, always run with `-race`:** `go run -race ./try` and `go test -race -v`. Lesson 21
explains why; for now, just make it a habit.

---

## 1. Why concurrency?

Many programs spend most of their time **waiting**: for a web page to download, a database to answer,
a file to be read. If you fetch 10 web pages one after another, and each takes 1 second, that's 10
seconds. If you start all 10 at once and wait for them together, it's about 1 second.

Two words you'll hear:
- **Concurrency:** dealing with many things at once (starting them, switching between them).
- **Parallelism:** actually *doing* many things at the same instant, on several CPU cores.

In Python you'd reach for `threading` (limited by the GIL, so only one thread runs Python code at a
time) or `asyncio` (needs `async`/`await` everywhere). Go has neither limit:
- **No GIL.** Goroutines really run in parallel, on all your CPU cores.
- **No async/await.** Any normal function can run concurrently. You don't rewrite your code.

---

## 2. Your first goroutine

Put `go` in front of a function call, and that call runs **in the background**, at the same time as
the rest of your code. That background function is a **goroutine**:

```go
func say(msg string) {
	for i := range 3 {
		fmt.Println(msg, i)
		time.Sleep(100 * time.Millisecond)
	}
}

func main() {
	go say("background")      // starts, and main continues immediately
	say("main")               // runs at the same time
}
```

Run it a few times. The lines from `main` and `background` are mixed together, and the order can
change between runs. That's concurrency: two things progressing at once.

Goroutines are **very cheap**: a few kilobytes each. Starting 100,000 of them is normal. Python
threads cost far more.

---

## 3. ⚠️ When main ends, everything ends

Now try this:

```go
func main() {
	go say("background")
	fmt.Println("main is done")
}
```

You'll probably see only `main is done`. **When `main` returns, the program exits immediately, and
all goroutines are killed**, finished or not. Go doesn't wait for them automatically.

You might be tempted to add `time.Sleep(time.Second)` at the end of `main`. Don't: you'd be guessing
how long the work takes. You need a way to **wait until the goroutines are actually done**.

---

## 4. `sync.WaitGroup`: waiting for goroutines

A `WaitGroup` (package `sync`) is a counter of "goroutines still running":

```go
func main() {
	var wg sync.WaitGroup

	for i := range 3 {
		wg.Add(1)                         // 1. one more goroutine to wait for
		go func() {
			defer wg.Done()               // 2. this one is finished (when the func returns)
			fmt.Println("worker", i, "done")
		}()
	}

	wg.Wait()                             // 3. block until the counter is back to 0
	fmt.Println("all workers done")
}
```

The three steps:
1. **`wg.Add(1)`** before starting each goroutine. Call it in `main`, not inside the goroutine
   (otherwise `Wait` might run before the `Add` happens).
2. **`defer wg.Done()`** as the first line inside the goroutine. `defer` makes sure it's called even if
   the function returns early.
3. **`wg.Wait()`** blocks until every `Done` has been called.

Notice the pattern `go func() { ... }()`: an anonymous function (lesson 02), started as a goroutine.
The `()` at the end calls it. This is how most goroutines are written.

(Go 1.25+ also has `wg.Go(func() { ... })`, which does the `Add` and `Done` for you. This course's
`go.mod` targets Go 1.24 so the exercises work with older versions, so use `Add`/`Done` here.)

Python comparison:

```python
threads = [threading.Thread(target=work, args=(i,)) for i in range(3)]
for t in threads: t.start()
for t in threads: t.join()        # like wg.Wait()
```

### Loop variables are safe (since Go 1.22)

In the example, each goroutine uses `i`. Since Go 1.22, every loop iteration has its **own** `i`, so
each goroutine sees the right number. (Old Go code and old blog posts show a workaround,
`go func(i int) {...}(i)`. You don't need it anymore.)

---

## 5. Getting results back

A goroutine can't `return` a value to the code that started it. `go f()` throws the return value
away. So how do you get results?

The simplest safe way: **give each goroutine its own slot in a slice**, by index.

```go
func squareAll(nums []int) []int {
	results := make([]int, len(nums))
	var wg sync.WaitGroup
	for i, n := range nums {
		wg.Add(1)
		go func() {
			defer wg.Done()
			results[i] = n * n        // each goroutine writes a DIFFERENT index
		}()
	}
	wg.Wait()                         // after this, all results are filled in
	return results
}
```

This is safe because **no two goroutines touch the same element**. And reading `results` after
`wg.Wait()` is safe because `Wait` guarantees all the writes have finished.

What is **not** safe is several goroutines changing the **same** variable:

```go
total := 0
for _, n := range nums {
	wg.Add(1)
	go func() {
		defer wg.Done()
		total += n        // ❌ DATA RACE: many goroutines change total at once
	}()
}
```

`total += n` is really three steps: read `total`, add, write it back. Two goroutines can both read
the same old value, and one update gets lost. The result is randomly wrong. Run it with `-race` and
Go reports `WARNING: DATA RACE`. Lesson 18 (channels) and lesson 21 (mutexes) show the proper fixes.
For now: **one slot per goroutine**, then combine the slots after `Wait()`.

### 🤔 Predict the output

```go
func main() {
	var wg sync.WaitGroup
	results := make([]string, 3)
	for i, name := range []string{"a", "b", "c"} {
		wg.Add(1)
		go func() {
			defer wg.Done()
			time.Sleep(time.Duration(3-i) * 10 * time.Millisecond)
			results[i] = strings.ToUpper(name)
		}()
	}
	wg.Wait()
	fmt.Println(results)
}
```

<details><summary>Answer</summary>

`[A B C]`. The goroutines **finish** in the order c, b, a (because of the different sleeps), but each
writes to its own index, so the slice is in the original order. That's a useful property of the
"one slot per goroutine" pattern.
</details>

---

## 6. How fast is it really?

```go
start := time.Now()
var wg sync.WaitGroup
for range 10 {
	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(100 * time.Millisecond)   // pretend this is a slow network call
	}()
}
wg.Wait()
fmt.Println("took", time.Since(start))      // about 100ms, not 1s
```

Ten "slow calls" of 100ms each take about 100ms in total, because they all wait at the same time.
That's the whole point.

---

## 7. Common mistakes

| You wrote | What happens | Fix |
|---|---|---|
| `go f()` and `main` returns | goroutine killed before finishing | `sync.WaitGroup` |
| `time.Sleep` to "wait for" goroutines | flaky: sometimes too short, always wasteful | `wg.Wait()` |
| `wg.Add(1)` inside the goroutine | `Wait` can return too early | `wg.Add(1)` before `go` |
| Forgot `wg.Done()` | `fatal error: all goroutines are asleep - deadlock!` | `defer wg.Done()` first thing |
| Several goroutines doing `total += n` | wrong totals, `DATA RACE` with `-race` | one slot per goroutine, or lessons 18 / 21 |
| Passing a `WaitGroup` by value to a function | the copy's counter is separate: deadlock | pass `*sync.WaitGroup`, or use a closure |

---

## 8. Exercises

Run with the race detector: `go test -race -run Ex1 -v`

The tests measure time, to check your code really runs things **at the same time** and doesn't just
run them one after another.

| | Exercise | Practises |
|---|---|---|
| Ex1 | `RunBoth`: run two functions at the same time | `go`, `WaitGroup` |
| Ex2 | `RunAll`: run any number of functions | WaitGroup in a loop |
| Ex3 | `SquareAll`: one goroutine per number | one slot per goroutine |
| Ex4 | `FetchAll`: fake slow downloads, in parallel | the same pattern, with timing |
| Ex5 | `ParallelSum`: split the work into parts | partial results combined after `Wait` |

---

## ✅ Checklist

- [ ] I can start a goroutine with `go`, including `go func() { ... }()`
- [ ] I know `main` returning kills all goroutines
- [ ] I can use `wg.Add(1)`, `defer wg.Done()` and `wg.Wait()` correctly
- [ ] I can collect results with one slot per goroutine
- [ ] I can explain why `total += n` from many goroutines is a data race

**Next → [Lesson 18: Channels](../18_channels/README.md)**
