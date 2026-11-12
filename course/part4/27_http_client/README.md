# Lesson 27: HTTP clients, middleware and logging

**Goal:** call other web services safely (timeouts, status checks, closing bodies), wrap handlers with
middleware for logging, panic recovery and auth, and log in a structured way with `log/slog`.

**Time:** about 2 hours.

```bash
cd ~/Desktop/learning-go/course/part4/27_http_client
```

---

## Part A: calling other services

## 1. The simplest request, and why not to use it

```go
resp, err := http.Get("https://api.github.com")
```

That works, but `http.Get` uses the **default client, which has no timeout**. If the other server hangs,
your program waits forever. In a web server, every such request is a goroutine stuck forever. Python's
`requests` has the same trap (`requests.get(url)` without `timeout=`).

## 2. A client with a timeout

Create your own client once, and reuse it everywhere:

```go
client := &http.Client{Timeout: 10 * time.Second}
```

Reuse matters: a client keeps connections open and reuses them, which is much faster. Don't create a
new client per request.

## 3. A request with a context

For full control, build the request yourself, with a context (lesson 22), so it's cancelled when your
caller gives up:

```go
func getText(ctx context.Context, client *http.Client, url string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err                       // network error, timeout, cancelled...
	}
	defer resp.Body.Close()                  // ⚠️ always close the body

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
```

Three rules, in order:
1. **Check `err`** from `client.Do`. It covers network failures and timeouts.
2. **`defer resp.Body.Close()`**. Not closing it leaks connections, and a long-running program slowly breaks.
3. **Check the status code.** A `404` or `500` is **not** an error from `client.Do`! The request
   "worked"; the server just said no. Python's `requests` is the same (that's what
   `raise_for_status()` is for).

## 4. JSON in, JSON out

Decode straight from the body (lesson 25):

```go
var user User
if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
	return User{}, err
}
```

To send JSON, encode it into a buffer and use it as the body:

```go
data, err := json.Marshal(payload)
if err != nil {
	return err
}
req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(data))
if err != nil {
	return err
}
req.Header.Set("Content-Type", "application/json")
resp, err := client.Do(req)
```

## 5. Testing clients with `httptest.NewServer`

To test code that calls a server, start a **real local server** with fake behaviour:

```go
srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "hello")
}))
defer srv.Close()

text, err := getText(context.Background(), srv.Client(), srv.URL)
```

`srv.URL` is something like `http://127.0.0.1:54321`. No internet needed, and you control exactly what
the "other service" does, including errors and slow responses. This lesson's tests work this way.

---

## Part B: middleware

## 6. What middleware is

You often want the same thing done for **every** request: log it, recover from panics, check an API key,
add a request ID. Copying that code into every handler is a bad idea. Instead, **wrap** the handler:

```go
func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)                                  // call the real handler
		log.Printf("%s %s took %v", r.Method, r.URL.Path, time.Since(start))
	})
}
```

A middleware is a function that **takes a handler and returns a new handler**:
`func(http.Handler) http.Handler`. The new handler does something before and/or after calling `next`.

`http.HandlerFunc(...)` turns an ordinary function into an `http.Handler` (it's a type that has a
`ServeHTTP` method). You'll write this line in every middleware.

Use it by wrapping:

```go
mux := http.NewServeMux()
mux.HandleFunc("GET /hello", hello)
http.ListenAndServe(":8080", logging(mux))     // every request goes through logging first
```

Python comparison: decorators, or Flask's `before_request` / Django middleware. Same idea, written as
plain functions.

## 7. Middleware that stops the request

A middleware doesn't have to call `next`. An auth check can reject the request:

```go
func requireKey(key string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("X-API-Key") != key {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return                               // next is never called
			}
			next.ServeHTTP(w, r)
		})
	}
}
```

When a middleware needs settings (here, the key), you write a function that **returns** the middleware.
That's three nested functions; read it from the outside in. Usage: `requireKey("secret")(mux)`.

## 8. Recovering from panics

A panic in one handler shouldn't take down the whole server. (`net/http` already recovers and closes
the connection, but the client gets no proper response.) A recovery middleware turns a panic into a
clean `500`, using `recover` from lesson 14:

```go
func recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if v := recover(); v != nil {
				http.Error(w, "internal error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
```

## 9. Chaining middleware

```go
handler := recoverer(logging(requireKey("secret")(mux)))
```

The **outermost** wrapper runs first. Here a request goes: recoverer → logging → requireKey → mux.
A small helper makes long chains readable; you'll write one in the exercises.

---

## Part C: structured logging

## 10. `log/slog`

`log.Printf` writes free-form text. For services, **structured** logs (key=value pairs, or JSON) are
much easier to search and filter. The standard library has `log/slog`:

```go
logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
logger.Info("request handled", "method", "GET", "path", "/notes", "status", 200)
// time=2026-11-13T10:00:00.000+04:00 level=INFO msg="request handled" method=GET path=/notes status=200

jsonLogger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
jsonLogger.Error("db down", "err", err)
// {"time":"...","level":"ERROR","msg":"db down","err":"connection refused"}
```

After the message, pass alternating keys and values. Levels: `Debug`, `Info`, `Warn`, `Error`.
Python comparison: the `logging` module plus a JSON formatter, or `structlog`.

Pass the logger to the things that need it (like a middleware), instead of using a global. In tests you
can then give it a `bytes.Buffer` and check what was logged.

### 🤔 Predict the output

```go
mw := func(name string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Print(name, ">")
			next.ServeHTTP(w, r)
			fmt.Print("<", name, " ")
		})
	}
}
final := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { fmt.Print("handler ") })
h := mw("A")(mw("B")(final))
h.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil))
```

<details><summary>Answer</summary>

`A>B>handler <B <A `. The outer middleware (A) runs first on the way in and last on the way out, like
nested function calls.
</details>

---

## 11. Common mistakes

| You wrote | What happens | Fix |
|---|---|---|
| `http.Get(url)` in a server | no timeout: requests can hang forever | `&http.Client{Timeout: ...}` |
| Forgot `defer resp.Body.Close()` | connection leak | close it right after checking `err` |
| Treating `err == nil` as success | 404s and 500s look like success | check `resp.StatusCode` |
| A new `http.Client` for every request | slow, no connection reuse | create one and reuse it |
| Middleware that forgets to call `next` | every request returns an empty 200 | `next.ServeHTTP(w, r)` |
| Auth check with no `return` after the error | the handler runs anyway | `return` after `http.Error` |

---

## 12. Exercises

| | Exercise | Practises |
|---|---|---|
| Ex1 | `NewClient` with a timeout | `http.Client` |
| Ex2 | `GetText`: the three rules | context, `Do`, `Body.Close`, status codes |
| Ex3 | `GetUser`: decode a JSON response | `json.NewDecoder(resp.Body)` |
| Ex4 | `CreateUser`: POST JSON | request bodies and headers |
| Ex5 | `Logging` middleware with `slog` | middleware, structured logs |
| Ex6 | `Recover` middleware | `recover` in middleware |
| Ex7 | `RequireAPIKey` middleware | stopping a request, middleware with settings |
| Ex8 | `Chain` several middlewares | functions of functions |

---

## ✅ Checklist

- [ ] I never use a client without a timeout in real code
- [ ] I always close response bodies and check status codes
- [ ] I can send and receive JSON with an HTTP client, and test it with `httptest.NewServer`
- [ ] I can write a middleware, including one that takes settings
- [ ] I can log with `slog` using key/value pairs

**Next → [Lesson 28: Databases](../28_database/README.md)**
