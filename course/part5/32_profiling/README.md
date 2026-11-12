# Lesson 32: Profiling and optimizing

**Goal:** find out **where** a program spends its time with profilers instead of guessing, fix the
hot spots, and prove the improvement with numbers.

**Time:** about 2 hours.

```bash
cd ~/Desktop/learning-go/course/part5/32_profiling
```

---

## 1. Don't guess. Measure.

Programmers are famously bad at guessing which part of a program is slow. The usual result of guessing
is complicated code that's no faster. The professional loop is:

1. **Measure**: a benchmark that reproduces the slow case.
2. **Profile**: find which functions use the time (or memory).
3. **Fix** the biggest hot spot, and only that.
4. **Measure again**, and compare. Keep the change only if it really helped.

Go has all the tools built in. Python comparison: `cProfile` + `snakeviz`, but with more detail and
almost no overhead, so you can even profile production servers.

---

## 2. Benchmarks, again

From lesson 29:

```go
func BenchmarkUniqueWords(b *testing.B) {
	text := sampleText()             // setup is not timed: b.Loop starts the clock
	for b.Loop() {
		UniqueWords(text)
	}
}
```

```bash
go test -bench=UniqueWords -benchmem
# BenchmarkUniqueWords-10    48    24310522 ns/op    1203465 B/op    42 allocs/op
```

`24310522 ns/op` is about 24 ms per call. That's the number to beat.

Benchmark results vary a bit between runs (your laptop is doing other things). Run several times with
`-count=10`, and compare old and new with **`benchstat`**:

```bash
go install golang.org/x/perf/cmd/benchstat@latest

go test -bench=UniqueWords -count=10 > old.txt
# ... make your change ...
go test -bench=UniqueWords -count=10 > new.txt
benchstat old.txt new.txt
```

```
              │   old.txt    │               new.txt               │
              │    sec/op    │   sec/op     vs base                │
UniqueWords-10   24.31m ± 1%   0.41m ± 2%  -98.31% (p=0.000 n=10)
```

`benchstat` tells you whether the difference is real or just noise (`p=0.000` means "definitely real").

---

## 3. CPU profiling

A **CPU profile** samples the program about 100 times per second and records which function was
running. Functions that appear in many samples are where the time goes.

```bash
go test -bench=Summarize -cpuprofile=cpu.out
go tool pprof -top cpu.out
```

```
      flat  flat%   sum%        cum   cum%
     0.97s 43.11% 43.11%      1.83s 81.33%  lesson32.findCustomer (inline)
     0.87s 38.67% 81.78%      0.87s 38.67%  runtime.memequal
     0.21s  9.33% 91.11%      0.21s  9.33%  runtime.madvise
     ...
```

How to read it:
- **flat**: time spent **in this function itself**.
- **cum** (cumulative): time in this function **plus everything it calls**.
- A function with a big **flat** is a hot spot. A function with a big **cum** but small **flat** is
  just a caller of hot spots.
- **`runtime.*` functions** are Go's own internals, called on your behalf. `runtime.memequal` compares
  memory (string comparisons, for example). Ask "which of **my** functions is calling it?". The `cum`
  column and `-list` (below) tell you.
- `(inline)` means the compiler pasted that small function into its caller (lesson 31).

Here `findCustomer` uses 43% of the time itself, and 81% including the string comparisons it triggers.
That's where to look.

### Looking at the code line by line

```bash
go tool pprof -list=lookupRegion cpu.out
```

shows the function's source with the time spent on each line.

### The web view

```bash
go tool pprof -http=:8080 cpu.out
```

opens an interactive page in your browser with a call graph and a **flame graph** (View → Flame Graph):
each bar is a function, its width is the time spent, and callers sit above their callees. Wide bars
are expensive. This is the view most people use day to day.

---

## 4. Memory profiling

```bash
go test -bench=BuildReport -memprofile=mem.out
go tool pprof -sample_index=alloc_space -top mem.out
```

This shows **which functions allocate the most memory**. Combined with lesson 31, it tells you where
allocations are worth removing.

---

## 5. Profiling a running server

For a long-running service, import `net/http/pprof` for its side effect (the same `_` import as database
drivers in lesson 28). It adds debug endpoints:

```go
import _ "net/http/pprof"

go http.ListenAndServe("localhost:6060", nil)   // a separate, private port
```

```bash
go tool pprof "http://localhost:6060/debug/pprof/profile?seconds=30"   # 30s CPU profile
go tool pprof http://localhost:6060/debug/pprof/heap                   # memory in use
curl "http://localhost:6060/debug/pprof/goroutine?debug=1"              # every goroutine's stack
```

The goroutine dump is invaluable for finding leaks and deadlocks (lesson 33). ⚠️ Never expose these
endpoints to the internet. Bind them to `localhost`, or protect them.

---

## 6. Execution traces

`pprof` tells you **what** uses CPU. A **trace** shows **when** things happen: which goroutines ran on
which CPU, when they were blocked, when the GC ran. Use it for "it's slow but the CPU isn't busy"
problems (waiting on locks or channels, too little parallelism):

```bash
go test -bench=. -trace=trace.out
go tool trace trace.out
```

---

## 7. The usual suspects

Most real-world Go slowness comes from a short list. Learn to recognize these in a profile:

| Hot spot in the profile | Typical cause | Typical fix |
|---|---|---|
| A search function, `runtime.memequal`, string comparisons | looking things up in a slice in a loop: O(n²) | build a map once: O(n) |
| `regexp.Compile` / `regexp.MustCompile` | compiling a regex on every call | compile once into a package-level variable |
| `runtime.mallocgc`, `runtime.growslice`, `concatstrings` | allocations in a hot loop | preallocate, `strings.Builder` (lesson 31) |
| `fmt.Sprintf` | formatting in a tight loop | `strconv.Itoa`, `strconv.AppendInt`, direct writes |
| `encoding/json` reflection functions | JSON in a hot path | decode once, cache, or decode only what's needed |
| `syscall`, `runtime.futex` | unbuffered I/O, lock contention | `bufio`, fewer/smaller critical sections |

The biggest wins are almost always **algorithmic** (fixing an O(n²)), not micro-tuning.

### 🤔 Predict the output

A profile shows:

```
      flat  flat%        cum   cum%
     0.05s  1.20%      4.10s 98.50%  main.processAll
     3.90s 93.70%      3.90s 93.70%  main.findUser
```

Which function do you optimize?

<details><summary>Answer</summary>

`findUser`: it spends 3.9 s **itself** (flat). `processAll` has a high **cum** only because it calls
`findUser`. Its own work (flat 0.05 s) is tiny, so optimizing it wouldn't help.
</details>

---

## 8. Exercises

Each exercise gives you a **slow but correct** function. Make it fast without changing what it returns.
The tests check both: same results, and much faster than the original (the original is kept in the test
file for comparison).

| | Exercise | The problem to find |
|---|---|---|
| Ex1 | `UniqueWords` | O(n²) search in a slice |
| Ex2 | `HasCommon` | nested loops |
| Ex3 | `ValidEmails` | something expensive repeated on every call |
| Ex4 | `BuildReport` | string building |
| Ex5 | `Summarize` | **use the profiler** to find which helper is slow, then fix it |

For Ex5, follow these steps:

```bash
go test -bench=Summarize -cpuprofile=cpu.out
go tool pprof -top cpu.out
go tool pprof -list=Summarize cpu.out
```

Write the name of the hot function into the `HotFunction` constant, then fix it.

Run the tests without `-race` (the race detector slows code down unevenly and would distort the timings):
`go test -run Ex1 -v`

---

## ✅ Checklist

- [ ] I measure before and after any optimization, with benchmarks and `benchstat`
- [ ] I can take a CPU profile and read `flat` vs `cum` in `pprof -top`
- [ ] I can open the flame graph with `go tool pprof -http`
- [ ] I know how to profile a live server with `net/http/pprof`
- [ ] I recognize the usual suspects, and look for algorithmic wins first

**Next → [Lesson 33: The scheduler, leaks and deadlocks](../33_scheduler/README.md)**
