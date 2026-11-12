# Learning Go: zero to hero

A step-by-step Go course for Python developers who have **never written any Go**.

You don't need to know anything about Go to start. Every lesson explains one idea at a time,
shows you a small example to type and run, asks you to predict what code does, and only then
gives you exercises. The exercises start very easy and get harder slowly.

**Start here → [course/part1/00_start/README.md](course/part1/00_start/README.md)**

## Get the course

Read it online at **https://letsgomrinal.github.io/learning-go/**, but do the exercises on your own
computer. You need Go installed ([go.dev/dl](https://go.dev/dl/), or `brew install go` on a Mac):

```bash
git clone https://github.com/letsgomrinal/learning-go.git
cd learning-go
go version
```

The lessons say `~/Desktop/learning-go`. Read that as "wherever you cloned the course".

---

## How every lesson works

Each lesson is a folder with the same four things:

| File | What it's for |
|---|---|
| `README.md` | **Read this first.** The explanation, with Python comparisons, small examples and "predict the output" questions. |
| `try/main.go` | Your playground. Type the README's examples here and run them with `go run ./try`. |
| `exercises.go` | The exercises. Each one is a function with a `TODO` inside for you to fill in. |
| `exercises_test.go` | Tests that check your exercises. You don't need to edit this file (but reading it is good practice). |

The rhythm for a lesson:

1. Open a terminal **in the lesson folder**, for example: `cd ~/Desktop/learning-go/course/part1/00_start`
2. Read `README.md` from top to bottom. Type every example into `try/main.go` and run it. **Type it, don't copy-paste.** Typing is how your fingers learn the syntax.
3. Before you reveal a "predict the output" answer, actually decide what you think it prints.
4. Do the exercises one at a time:
   ```bash
   go test -run Ex1 -v      # check exercise 1 only
   go test -run Ex2 -v      # then exercise 2, and so on
   go test -v               # finally, everything in the lesson
   ```
5. Stuck for more than 20 minutes? Each exercise has hints. If you're still stuck, the answer is in
   `_solutions/` at the same path. Read it, close it, then write it again yourself from memory.

Expect **one lesson per 1–2 sessions** of about an hour. Going slowly is fine. Skipping the typing is not.

---

## The course map

### Part 1: Foundations (weeks 1–4) ✅ written

| # | Lesson | You'll learn |
|---|---|---|
| 00 | [Your first program](course/part1/00_start/README.md) | Running Go, reading compiler errors, how the exercises and tests work |
| 01 | [Variables and types](course/part1/01_variables/README.md) | `var`, `:=`, basic types, zero values, constants, converting types, printing |
| 02 | [Functions](course/part1/02_functions/README.md) | Parameters, return values, returning more than one value, functions as values |
| 03 | [Control flow](course/part1/03_control_flow/README.md) | `if`, the one loop (`for`), `switch` |
| 04 | [Strings](course/part1/04_strings/README.md) | The `strings` and `strconv` packages, bytes vs characters |
| 05 | [Slices](course/part1/05_slices/README.md) | Go's version of a Python list |
| 06 | [Maps](course/part1/06_maps/README.md) | Go's version of a Python dict, and sets |
| 07 | [Structs](course/part1/07_structs/README.md) | Grouping data together (where Python would use a class or dataclass) |
| 08 | [Pointers and methods](course/part1/08_pointers_methods/README.md) | What a pointer is, and attaching methods to your types |
| 09 | [Packages and tests](course/part1/09_packages_tests/README.md) | Splitting code into packages, public vs private, writing your own tests |

After Part 1: do the practice sets `practice/01_basics`, `practice/02_slices_maps` and
`practice/03_structs_methods`. They're harder on purpose and mix everything together.

### Part 2: The Go way (weeks 5–8) ✅ written

| # | Lesson | You'll learn |
|---|---|---|
| 10 | [Interfaces](course/part2/10_interfaces/README.md) | Duck typing checked by the compiler, implicit interfaces, pointer receivers |
| 11 | [Interfaces in practice](course/part2/11_interfaces_practice/README.md) | `fmt.Stringer`, `any`, type assertions and switches, `io.Writer`, `io.Reader` |
| 12 | [Errors](course/part2/12_errors/README.md) | Returning and handling errors, `if err != nil`, sentinel errors, `errors.Is` |
| 13 | [Custom errors and wrapping](course/part2/13_error_types/README.md) | `%w`, error chains, custom error types, `errors.As`, the nil interface trap |
| 14 | [defer, panic and recover](course/part2/14_defer_panic/README.md) | Reliable cleanup, what panics are for, turning a panic into an error |
| 15 | [Generics](course/part2/15_generics/README.md) | Type parameters, constraints, generic types, generics vs interfaces |
| 16 | [Iterators](course/part2/16_iterators/README.md) | `iter.Seq`, `yield`, Go's version of Python generators |

After Part 2: do the practice sets `practice/04_interfaces`, `practice/05_errors` and `practice/06_generics`.

### Part 3: Concurrency (weeks 9–13) ✅ written

| # | Lesson | You'll learn |
|---|---|---|
| 17 | [Goroutines and WaitGroup](course/part3/17_goroutines/README.md) | Running things at the same time, waiting for them, collecting results safely |
| 18 | [Channels](course/part3/18_channels/README.md) | Sending values between goroutines, blocking, closing, deadlocks |
| 19 | [select and timeouts](course/part3/19_select/README.md) | Waiting on several channels, timeouts, non-blocking operations, stopping goroutines |
| 20 | [Pipelines and worker pools](course/part3/20_patterns/README.md) | Pipelines, fan-in, worker pools, limiting concurrency |
| 21 | [Mutexes and the race detector](course/part3/21_mutex/README.md) | Data races, `-race`, `Mutex`, `RWMutex`, atomics, `sync.Once` |
| 22 | [Context](course/part3/22_context/README.md) | Cancellation, timeouts, and cancelling the rest when one task fails |

After Part 3: do the practice sets `practice/07_concurrency` and `practice/08_sync_context`.

### Part 4: Building real programs (weeks 14–19) ✅ written

| # | Lesson | You'll learn |
|---|---|---|
| 23 | [Files and I/O](course/part4/23_files/README.md) | Reading and writing files, line by line, walking folders, "file not found" |
| 24 | [Command-line tools](course/part4/24_cli/README.md) | Arguments, flags, stdin/stdout/stderr, exit codes. You build a `wc` clone |
| 25 | [JSON](course/part4/25_json/README.md) | Struct tags, nested JSON, unknown JSON, streaming |
| 26 | [HTTP servers](course/part4/26_http_server/README.md) | Routes, path parameters, JSON APIs, status codes, `httptest`. You build a notes API |
| 27 | [HTTP clients and middleware](course/part4/27_http_client/README.md) | Timeouts, calling APIs safely, middleware, `log/slog` |
| 28 | [Databases](course/part4/28_database/README.md) | `database/sql` with SQLite, placeholders, rows, transactions (its own module) |
| 29 | [Testing in depth](course/part4/29_testing/README.md) | You write the tests: subtests, fakes, benchmarks, fuzzing, examples |
| 30 | [Shipping](course/part4/30_shipping/README.md) | Config, `go:embed`, versions, graceful shutdown, cross-compiling, Docker |

After Part 4: do the practice set `practice/09_http_json`, then build the URL shortener milestone.

### Part 5: Expert topics (weeks 20–26) ✅ written

| # | Lesson | You'll learn |
|---|---|---|
| 31 | [Memory, allocations and escape analysis](course/part5/31_memory/README.md) | Stack vs heap, `-gcflags=-m`, allocating less, struct padding, the GC |
| 32 | [Profiling and optimizing](course/part5/32_profiling/README.md) | Benchmarks, `benchstat`, `pprof`, flame graphs. You speed up slow code |
| 33 | [The scheduler, leaks and deadlocks](course/part5/33_scheduler/README.md) | G-M-P, `GOMAXPROCS`, worker counts, finding leaks, lock ordering |
| 34 | [The memory model and advanced sync](course/part5/34_memory_model/README.md) | Happens-before, atomics, singleflight, errgroup, `sync.Cond` |
| 35 | [Reflection and struct tags](course/part5/35_reflection/README.md) | `reflect`, custom tags, how `encoding/json` works inside |
| 36 | [Idioms from the standard library](course/part5/36_stdlib_idioms/README.md) | Reading stdlib source, composition, functional options, error types |

After Part 5: the expert milestone projects in roadmap.html, and the reading list at the end of lesson 36.

Open [roadmap.html](roadmap.html) in a browser for the full plan, with a progress tracker, reading lists
and the milestone projects for each part.

---

## A note on LLMs

You're learning the concepts, so for the lessons: **turn off AI autocomplete** in your editor. It's fine
to ask an AI to *explain* a concept or an error message. Don't ask it for the code. After your
exercise passes, asking "how would an experienced Go developer improve this?" is a great habit.
