# Lesson 36: Idioms from the standard library

**Goal:** learn to read the Go standard library's source code, and use the design patterns it's built
on in your own code. This is the last lesson of the course.

**Time:** about 2 hours.

```bash
cd ~/Desktop/learning-go/course/part5/36_stdlib_idioms
```

---

## 1. Read the source

Go's standard library is some of the best Go code ever written, and it's right there on your machine.
Reading it is one of the fastest ways to get better, and much easier than you might expect: it's plain
Go, with good comments.

```bash
go doc -src io.LimitReader         # show the source of one function or type
go doc -src strings.Builder.WriteString
go env GOROOT                      # where the standard library lives: $(go env GOROOT)/src
```

In VS Code or GoLand: **Cmd+click** (or F12) on any standard library name jumps to its source.

Good files to read, in roughly this order:
- `io/io.go`: `Reader`, `Writer`, `Copy`, `LimitReader`, `MultiWriter`. Small and full of ideas.
- `strings/builder.go`: how `strings.Builder` avoids copying.
- `sync/once.go`: a short masterpiece on correctness and speed.
- `errors/wrap.go`: how `errors.Is` and `errors.As` walk the chain.
- `context/context.go`: everything from lesson 22, in about 800 lines.
- `net/http/server.go`: big, but you'll recognise `ServeMux`, `HandlerFunc` and `Server.Shutdown`.

This lesson's exercises re-build small pieces of the standard library, using its idioms.

---

## 2. Idiom: small interfaces that compose

`io.Reader` and `io.Writer` have one method each. That's what lets you stack them like Lego:

```go
f, _ := os.Open("data.csv.gz")                  // a Reader (a file)
gz, _ := gzip.NewReader(f)                      // a Reader that decompresses another Reader
limited := io.LimitReader(gz, 1<<20)            // a Reader that stops after 1 MB
sc := bufio.NewScanner(limited)                 // reads lines from any Reader
```

Each layer only knows "I have an `io.Reader`". Here's `io.LimitReader` from the standard library, almost
exactly:

```go
type LimitedReader struct {
	R io.Reader   // the underlying reader
	N int64       // bytes still allowed
}

func (l *LimitedReader) Read(p []byte) (n int, err error) {
	if l.N <= 0 {
		return 0, io.EOF
	}
	if int64(len(p)) > l.N {
		p = p[0:l.N]              // never read more than allowed
	}
	n, err = l.R.Read(p)
	l.N -= int64(n)
	return
}
```

You'll write `MultiWriter` (write to several writers at once) in the same style.

---

## 3. Idiom: a function type that implements an interface

`http.Handler` is an interface with one method. But most handlers are just functions. The standard library
bridges the two with a **function type that has a method**:

```go
type HandlerFunc func(ResponseWriter, *Request)

func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request) {
	f(w, r)                       // the method just calls the function itself
}
```

Now `http.HandlerFunc(myFunc)` is a `Handler`. You used this in lesson 27's middleware. It works for any
one-method interface you design: users can pass a full type, **or** a plain function.

---

## 4. Idiom: functional options

How do you build a type with many optional settings, given that Go has no default parameter values
or keyword arguments (lesson 02)? The popular answer is **functional options**:

```go
type Server struct {
	addr     string
	timeout  time.Duration
	maxConns int
}

type Option func(*Server)

func WithTimeout(d time.Duration) Option {
	return func(s *Server) { s.timeout = d }
}

func WithMaxConns(n int) Option {
	return func(s *Server) { s.maxConns = n }
}

func NewServer(addr string, opts ...Option) *Server {
	s := &Server{addr: addr, timeout: 30 * time.Second, maxConns: 100}   // defaults
	for _, opt := range opts {
		opt(s)                                                           // apply each option
	}
	return s
}
```

Callers only mention what they change, and it reads well:

```go
s := NewServer(":8080")
s := NewServer(":8080", WithTimeout(5*time.Second), WithMaxConns(1000))
```

It's the closest Go gets to Python's `Server(":8080", timeout=5, max_conns=1000)`. You'll see it in
gRPC, Zap, and many other libraries. (For simple cases, a plain `Config` struct is also fine and simpler.)

---

## 5. Idiom: error types that wrap

`os.Open` on a missing file returns an `*fs.PathError`:

```go
type PathError struct {
	Op   string
	Path string
	Err  error
}

func (e *PathError) Error() string { return e.Op + " " + e.Path + ": " + e.Err.Error() }
func (e *PathError) Unwrap() error { return e.Err }
```

The **`Unwrap() error`** method is what lets `errors.Is(err, fs.ErrNotExist)` see **through** the
`PathError` to the cause inside it (lesson 13 used `%w`, which builds the same thing for you). So a
custom error type can carry structured data (`Op`, `Path`) **and** keep the chain working. You'll write
one in Ex4.

---

## 6. Idiom: comparison functions with `cmp`

Sorting by several keys used to need the clunky `sort.Interface`. Now it's a comparison function with
`slices.SortFunc`, and `cmp.Compare` / `cmp.Or`:

```go
slices.SortFunc(people, func(a, b Person) int {
	return cmp.Or(
		cmp.Compare(a.Last, b.Last),      // first by last name
		cmp.Compare(a.First, b.First),    // then by first name
		cmp.Compare(b.Age, a.Age),        // then oldest first (note b, a)
	)
})
```

A comparison function returns negative (a first), zero (equal) or positive (b first). `cmp.Or` returns
its first non-zero argument, so it naturally means "compare by this, then by that". Python comparison:
`people.sort(key=lambda p: (p.last, p.first, -p.age))`.

---

## 7. Idiom: make the zero value useful

You've seen it all course long: `var mu sync.Mutex`, `var sb strings.Builder`, `var buf bytes.Buffer`,
`var wg sync.WaitGroup`. All are ready to use with **no constructor**. When you design types, aim for
the same: lazily create maps on first use (lesson 21's `Visits`), pick defaults that make sense as zero
values. Callers can't then forget to initialise anything.

### 🤔 Predict the output

```go
type Op func(int) int

func (o Op) Then(next Op) Op {
	return func(x int) int { return next(o(x)) }
}

double := Op(func(x int) int { return x * 2 })
inc := Op(func(x int) int { return x + 1 })
fmt.Println(double.Then(inc)(5), inc.Then(double)(5))
```

<details markdown="1">
<summary>Answer</summary>

`11 12`. A function type with a method (section 3) lets you chain behaviour: `double.Then(inc)` doubles
then adds one (5 → 10 → 11); `inc.Then(double)` adds then doubles (5 → 6 → 12).
</details>

---

## 8. Exercises

| | Exercise | The idiom |
|---|---|---|
| Ex1 | `LimitedReader` | composing `io.Reader`s |
| Ex2 | `MultiWriter` | composing `io.Writer`s |
| Ex3 | `ProcessorFunc` | a function type that implements an interface |
| Ex4 | `QueryError` with `Unwrap` | error types that keep the chain |
| Ex5 | `NewClient` with functional options | functional options and defaults |
| Ex6 | `SortPeople` by several keys | `slices.SortFunc` + `cmp.Or` |

---

## ✅ Checklist: end of the course 🎉

- [ ] I can open and read standard library source (`go doc -src`, Cmd+click)
- [ ] I can write types that compose through small interfaces, like `io.Reader`
- [ ] I can make a function type implement an interface
- [ ] I can design an API with functional options
- [ ] I can write a custom error type with `Unwrap`, and sort by several keys with `cmp.Or`

---

## You've finished the course

You've gone from your first `fmt.Println` to the scheduler, the memory model and profiling. Becoming an
expert from here is about **building real things** and **reading real code**. In roughly this order:

1. **Finish the milestone projects** in [roadmap.html](../../../roadmap.html), especially the URL shortener
   (Part 4) and one Part 5 project. Profile it. Make it fast. Write about what you found.
2. **Read** *"100 Go Mistakes and How to Avoid Them"* cover to cover. You'll now recognise almost all of them.
3. **Build something hard:**
   - a key-value store with a TCP protocol, a write-ahead log and crash recovery;
   - Raft consensus from the paper (the MIT 6.5840 distributed systems labs are in Go);
   - a job queue with retries, timeouts and at-least-once delivery.
4. **Contribute to open source Go projects** you use: Prometheus, Hugo, Caddy, CockroachDB, the
   Kubernetes ecosystem. Start with issues labelled "good first issue". Code review from experienced
   maintainers is the fastest teacher there is.
5. **Keep reading the standard library.** Every time you use a package, Cmd+click into it.

And now that your concepts are strong, use LLMs as the powerful assistants they are. You'll be able to
tell when their Go is idiomatic, when it leaks a goroutine, and when it's simply wrong.

**You've finished the course. Next → the expert milestone projects in [roadmap.html](../../../roadmap.html)**
