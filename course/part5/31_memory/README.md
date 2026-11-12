# Lesson 31: Memory, allocations and escape analysis

**Goal:** understand where Go puts your data (stack or heap), why that matters for speed, how to see
the compiler's decisions, and simple ways to allocate less.

**Time:** about 2 hours.

```bash
cd ~/Desktop/learning-go/course/part5/31_memory
```

**Welcome to Part 5.** Parts 1–4 taught you to write correct Go. Part 5 is about **why** Go behaves the
way it does: memory, performance, the scheduler, the memory model. This is what separates someone who
writes Go from someone who can debug a slow or misbehaving Go service in production.

A word of caution, which the Go community repeats often: **write clear code first, measure, and only
then optimize the parts the measurements point to.** Everything in this part is about knowing *how* to
make code faster, and lesson 32 is about knowing *where*.

---

## 1. Stack and heap

Your program's memory has two main areas.

**The stack.** Every goroutine has its own stack. When a function is called, it gets a "frame" on the
stack for its local variables; when it returns, the frame is thrown away instantly. It's extremely
cheap: no bookkeeping, no cleanup.

**The heap.** A shared area for data that must **outlive** the function that created it, or whose size
isn't known in advance. Allocating there costs more, and the **garbage collector** (GC) has to find
and free it later, which uses CPU time.

In Python, almost everything is on the heap (every object is). In Go, the compiler decides for each
value, and tries hard to keep things on the stack.

**Allocations** usually means heap allocations. Fewer allocations means less work for the GC, which
usually means faster programs. In hot code (code that runs millions of times), it can make a big difference.

---

## 2. Escape analysis

The compiler decides stack vs heap with **escape analysis**: if a value might still be used after its
function returns (it "escapes"), it goes on the heap. Otherwise it stays on the stack.

```go
func newPoint() *Point {
	p := Point{1, 2}
	return &p                     // p is used after newPoint returns → heap
}

func sum(nums []int) int {        // nums is only read here → doesn't escape
	t := 0
	for _, n := range nums {
		t += n
	}
	return t
}
```

Remember lesson 08: returning a pointer to a local variable is **safe** in Go. Escape analysis is
*why* it's safe: the compiler sees the pointer escapes and puts `p` on the heap.

### Seeing the compiler's decisions

```bash
go build -gcflags=-m ./try
```

```
./main.go:10:26: moved to heap: p
./main.go:11:10: nums does not escape
./main.go:12:33: []int{...} does not escape
./main.go:10:6: can inline newPoint
```

- **`moved to heap: p`**: `p` escapes, so it's allocated on the heap.
- **`does not escape`**: stays on the stack, free.
- **`can inline`**: the compiler can paste this small function's body into its callers, removing the
  call overhead (and often helping escape analysis too).
- **`escapes to heap`** next to an argument of `fmt.Println`: values passed as `any` (interfaces) often
  escape, because the compiler can't see what the callee will do with them.

Common reasons for escaping:
- returning a pointer to a local variable;
- storing a pointer in something that outlives the function (a global, a struct on the heap, a channel);
- converting a value to an interface (`any`, `error`, `fmt.Stringer`...), in many cases;
- a slice or map whose size isn't known at compile time, or that's too big for the stack;
- a closure that captures a variable and outlives the function.

Don't memorise the rules. Use `-gcflags=-m` and look.

---

## 3. Measuring allocations

Two tools.

**In benchmarks** (lesson 29), add `-benchmem`:

```bash
go test -bench=. -benchmem
# BenchmarkJoin-10   5000000   245 ns/op   112 B/op   4 allocs/op
```

`B/op` is bytes allocated per call, `allocs/op` is the number of heap allocations per call.

**In tests**, `testing.AllocsPerRun` counts the average allocations of a function:

```go
allocs := testing.AllocsPerRun(100, func() {
	Squares(1000)
})
if allocs > 1 {
	t.Errorf("Squares allocated %v times, want at most 1", allocs)
}
```

This lesson's tests use `AllocsPerRun` to check your optimizations. (Allocation counts change under
`-race`, so these tests are skipped when the race detector is on. Run them with plain `go test`.)

---

## 4. Four easy ways to allocate less

### a) Preallocate slices

```go
// grows several times: allocate, copy, allocate bigger, copy...
var out []int
for i := range n {
	out = append(out, i*i)
}

// one allocation: we know the final size
out := make([]int, 0, n)
for i := range n {
	out = append(out, i*i)
}
```

When `append` runs out of capacity it allocates a bigger array (about double) and copies everything.
Giving the right capacity up front avoids all of that. The same goes for maps: `make(map[string]int, n)`.

### b) Build strings with `strings.Builder` (and `Grow`)

```go
s := ""
for _, w := range words {
	s += w + ","             // a new string every time: n allocations, and lots of copying
}

var b strings.Builder
b.Grow(totalLen)             // optional: reserve the space once
for _, w := range words {
	b.WriteString(w)
	b.WriteByte(',')
}
s := b.String()              // no extra copy
```

### c) Return values instead of pointers (for small types)

```go
func newPoint() *Point { return &Point{1, 2} }   // allocates on the heap
func makePoint() Point  { return Point{1, 2} }   // usually no allocation at all
```

Copying a small struct is very cheap. Pointers aren't automatically "faster". For small structs they're
often slower, because they force heap allocations.

### d) Reuse buffers with `sync.Pool`

For temporary objects used very often (buffers in a busy server), `sync.Pool` keeps used objects around
so they can be reused instead of allocated again:

```go
var bufPool = sync.Pool{
	New: func() any { return new(bytes.Buffer) },
}

func render(name string) string {
	buf := bufPool.Get().(*bytes.Buffer)     // get one (new or reused)
	buf.Reset()                               // ⚠️ it may contain old data
	defer bufPool.Put(buf)                    // give it back when done

	buf.WriteString("hello, ")
	buf.WriteString(name)
	return buf.String()
}
```

Use it only when profiling shows allocations of a temporary object really matter. It adds complexity.

---

## 5. Struct layout and padding

The CPU reads memory fastest when values sit at addresses that are multiples of their size (an `int64`
at a multiple of 8). So the compiler inserts invisible **padding** between fields. Field **order**
changes the size of the struct:

```go
type Bad struct {
	A bool    // 1 byte + 7 padding (so B lines up at 8)
	B int64   // 8
	C bool    // 1 byte + 3 padding
	D int32   // 4
	E bool    // 1 byte + 7 padding at the end
}             // 32 bytes

type Good struct {
	B int64   // 8
	D int32   // 4
	A bool    // 1
	C bool    // 1
	E bool    // 1 + 1 padding
}             // 16 bytes

fmt.Println(unsafe.Sizeof(Bad{}), unsafe.Sizeof(Good{}))   // 32 16
```

Same data, half the memory. **Rule of thumb: order fields from largest to smallest.** It only matters
for structs you have millions of (in a big slice or cache), but then it matters a lot. (Tools like
`fieldalignment` from `golang.org/x/tools` find these automatically.)

`unsafe.Sizeof` is safe to call; the `unsafe` package name is a warning about its other functions.

---

## 6. The garbage collector, briefly

Go's GC runs **concurrently** with your program: most of its work happens alongside your goroutines,
with only very short pauses (usually well under a millisecond). Two settings exist:

- **`GOGC`** (default `100`): start a collection when the heap has grown 100% since the last one.
  Higher = less GC work but more memory; lower = the opposite.
- **`GOMEMLIMIT`** (for example `GOMEMLIMIT=1GiB`): a soft memory limit. The GC works harder as the
  program approaches it. Useful in containers with a fixed memory limit.

You can watch the GC work: `GODEBUG=gctrace=1 go run ./try` prints a line per collection.

Usually the best GC tuning is to **allocate less** (section 4). Read the official
[GC guide](https://go.dev/doc/gc-guide) when you need more.

### 🤔 Predict the output

```go
type Item struct {
	Active bool
	ID     int64
	Kind   bool
}

items := make([]Item, 1_000_000)
fmt.Println(unsafe.Sizeof(Item{}), "bytes each,", unsafe.Sizeof(Item{})*1_000_000/1_000_000, "MB total")
```

<details><summary>Answer</summary>

`24 bytes each, 24 MB total`. `Active` gets 7 bytes of padding so `ID` lines up, and `Kind` gets 7 more
at the end. Reordering to `ID, Active, Kind` gives 16 bytes each, saving 8 MB for a million items.
</details>

---

## 7. Common mistakes

| You did | Problem | Better |
|---|---|---|
| Optimized code without measuring | complexity for no gain, often slower | benchmark first (lesson 32) |
| Pointers everywhere "for speed" | more heap allocations, more GC work | values for small structs |
| `s += x` in a big loop | quadratic copying | `strings.Builder` |
| `append` in a loop with a known final size | repeated grow-and-copy | `make([]T, 0, n)` |
| Forgot `buf.Reset()` after `pool.Get()` | leftover data from a previous user | always reset |
| Storing things in a `sync.Pool` as a cache | the GC can empty the pool at any time | pools are for temporary objects only |

---

## 8. Exercises

Run with plain `go test -v` (not `-race`, because allocation counts change under the race detector).

| | Exercise | Practises |
|---|---|---|
| Ex1 | Reorder `Event`'s fields to make it 24 bytes | struct padding |
| Ex2 | `Squares` with exactly one allocation | preallocating slices |
| Ex3 | `JoinCSV` with at most one allocation | `strings.Builder` + `Grow` |
| Ex4 | `Distance` without any heap allocation | values instead of pointers |
| Ex5 | `Greet` using a `sync.Pool` buffer | `sync.Pool` |
| Ex6 | Explain the escape analysis of `try/main.go` | `-gcflags=-m` (written answer, checked) |

---

## ✅ Checklist

- [ ] I can explain stack vs heap, and why heap allocations cost more
- [ ] I can read `go build -gcflags=-m` output
- [ ] I can measure allocations with `-benchmem` and `testing.AllocsPerRun`
- [ ] I know four ways to allocate less, and when not to bother
- [ ] I can shrink a struct by reordering its fields

**Next → [Lesson 32: Profiling and optimizing](../32_profiling/README.md)**
