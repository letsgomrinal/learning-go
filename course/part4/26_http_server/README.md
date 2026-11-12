# Lesson 26: HTTP servers

**Goal:** build a JSON web API with only the standard library: routes, path parameters, reading
request bodies, writing JSON responses with the right status codes, and testing it all.

**Time:** about 2 hours.

```bash
cd ~/Desktop/learning-go/course/part4/26_http_server
```

---

## 1. A web server in 10 lines

```go
package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("GET /hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello from Go!")
	})
	fmt.Println("listening on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
```

Run it with `go run ./try`, then open http://localhost:8080/hello in your browser, or in another
terminal run `curl localhost:8080/hello`. Stop it with Ctrl+C.

No Flask, no FastAPI, no uvicorn. `net/http` is production-quality: many companies run large
services on it directly. Every request runs in **its own goroutine** automatically (Part 3), so the
server handles many requests at once without any extra work from you.

---

## 2. Handlers

A **handler** is a function with this exact signature:

```go
func(w http.ResponseWriter, r *http.Request)
```

- `r *http.Request` is everything about the incoming request: method, URL, headers, body.
- `w http.ResponseWriter` is how you send the response back. It's an `io.Writer` (lesson 11), so
  `fmt.Fprintln(w, ...)` and `json.NewEncoder(w)` work on it.

Python/Flask comparison:

```python
@app.get("/hello")
def hello():
    return "Hello from Python!"
```

In Go you don't return the response; you **write** it to `w`.

---

## 3. Routing with `ServeMux`

A `ServeMux` (a "router") sends each request to the right handler. Since Go 1.22, patterns can include
the HTTP method and path parameters:

```go
mux := http.NewServeMux()
mux.HandleFunc("GET /notes", listNotes)
mux.HandleFunc("POST /notes", createNote)
mux.HandleFunc("GET /notes/{id}", getNote)
mux.HandleFunc("DELETE /notes/{id}", deleteNote)

http.ListenAndServe(":8080", mux)
```

- `GET /notes/{id}` matches `/notes/7`, `/notes/abc`, etc.
- Read the parameter with `r.PathValue("id")`. It's always a string; convert with `strconv.Atoi`.
- A path that matches no pattern gets **404 Not Found** automatically.
- The right path with the wrong method (say `PUT /notes`) gets **405 Method Not Allowed** automatically.

Query parameters (`/notes?limit=10`) come from `r.URL.Query().Get("limit")`.

---

## 4. Writing responses: status, headers, body

The order matters: **headers first, then the status code, then the body**.

```go
func getNote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")   // 1. headers
	w.WriteHeader(http.StatusOK)                         // 2. status code (200)
	json.NewEncoder(w).Encode(note)                      // 3. body
}
```

- If you never call `WriteHeader`, the first write sends `200 OK`.
- After the body has started, you can't change the status or headers anymore.
- Use the named constants: `http.StatusOK` (200), `http.StatusCreated` (201), `http.StatusNoContent`
  (204), `http.StatusBadRequest` (400), `http.StatusNotFound` (404), `http.StatusInternalServerError` (500).

You'll write the same few lines in every handler, so make helpers:

```go
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
```

---

## 5. Reading a JSON request body

```go
func createNote(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Text string `json:"text"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return                                   // ⚠️ don't forget to return!
	}
	if strings.TrimSpace(input.Text) == "" {
		writeError(w, http.StatusBadRequest, "text is required")
		return
	}
	// ... save it, then:
	writeJSON(w, http.StatusCreated, note)
}
```

`r.Body` is an `io.Reader`, so `json.NewDecoder` reads straight from it (lesson 25).

The `var input struct { ... }` is an **anonymous struct**: a one-off type just for this request's
shape. It's common for request bodies.

⚠️ After writing an error, **`return`**. Otherwise the handler keeps going and writes a second
response on top (and Go logs `superfluous response.WriteHeader call`).

---

## 6. Shared state needs a mutex

Handlers run concurrently, one goroutine per request. If they share data (an in-memory list of notes),
protect it with a mutex (lesson 21):

```go
type Server struct {
	mu     sync.Mutex
	notes  map[int]Note
	nextID int
}

func (s *Server) createNote(w http.ResponseWriter, r *http.Request) {
	// ... decode and validate input ...
	s.mu.Lock()
	note := Note{ID: s.nextID, Text: input.Text}
	s.notes[note.ID] = note
	s.nextID++
	s.mu.Unlock()
	writeJSON(w, http.StatusCreated, note)
}
```

Methods can be handlers too: `mux.HandleFunc("POST /notes", s.createNote)`. That's the usual way to
give handlers access to shared state (the database, config, a logger) without global variables.

---

## 7. Your server as an `http.Handler`

`http.ListenAndServe` takes an **`http.Handler`**, an interface (lesson 10):

```go
type Handler interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}
```

`*http.ServeMux` is a Handler. Your own type can be one too, by forwarding to its mux:

```go
func NewServer() *Server {
	s := &Server{notes: map[int]Note{}, nextID: 1, mux: http.NewServeMux()}
	s.mux.HandleFunc("GET /notes", s.listNotes)
	// ...
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}
```

Now `http.ListenAndServe(":8080", NewServer())` works, and so do tests.

---

## 8. Testing handlers with `httptest`

You don't need to start a real server to test. `net/http/httptest` gives you a fake request and a
**recorder** that captures the response:

```go
func TestCreateNote(t *testing.T) {
	s := NewServer()

	req := httptest.NewRequest("POST", "/notes", strings.NewReader(`{"text":"buy milk"}`))
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201", rec.Code)
	}
	var got Note
	json.Unmarshal(rec.Body.Bytes(), &got)
	if got.Text != "buy milk" {
		t.Errorf("text = %q", got.Text)
	}
}
```

It's fast (no network) and fully reliable. All of this lesson's tests work this way; read them as examples.

### 🤔 Predict the output

```go
mux := http.NewServeMux()
mux.HandleFunc("GET /users/{name}", func(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "hi ", r.PathValue("name"))
})

cases := []struct{ method, path string }{
	{"GET", "/users/ana"},
	{"POST", "/users/ana"},
	{"GET", "/teams"},
}
for _, c := range cases {
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(c.method, c.path, nil))
	fmt.Println(rec.Code, strings.TrimSpace(rec.Body.String()))
}
```

<details><summary>Answer</summary>

```
200 hi ana
405 Method Not Allowed
404 404 page not found
```

The router handles wrong methods and unknown paths for you.
</details>

---

## 9. Common mistakes

| You wrote | What happens | Fix |
|---|---|---|
| `writeError(...)` with no `return` after it | the handler continues, writes twice | `return` after every error response |
| `w.WriteHeader(201)` then `w.Header().Set(...)` | the header is ignored (already sent) | headers, then status, then body |
| Shared map in handlers without a lock | random crashes: `concurrent map writes` | protect with a mutex |
| `strconv.Atoi(r.PathValue("id"))` without checking the error | id 0 for `/notes/abc` | return 400 on error |
| Returning `nil` slice as JSON | `null` instead of `[]` | `[]Note{}` |
| Global variables for state | hard to test, shared between tests | a `Server` struct with methods |

---

## 10. Exercises: build a notes API

You'll build a small in-memory notes API, one endpoint at a time. Each exercise has its own test.

| | Endpoint | Behaviour |
|---|---|---|
| Ex1 | `writeJSON` / `writeError` helpers | JSON content type, status, body |
| Ex2 | `GET /health` | `200 {"status":"ok"}` |
| Ex3 | `POST /notes` | create from `{"text": "..."}` → `201` with the note; `400` on bad JSON or empty text |
| Ex4 | `GET /notes/{id}` | `200` with the note; `404` if missing; `400` if id isn't a number |
| Ex5 | `GET /notes` | `200` with all notes sorted by id, `[]` when empty |
| Ex6 | `DELETE /notes/{id}` | `204` with no body; `404` if missing |
| Ex7 | `GET /notes?q=milk` | filter by text (case-insensitive) |

Try the real server once your tests pass:

```bash
go run ./cmd/notes
curl -X POST localhost:8080/notes -d '{"text":"learn Go"}'
curl localhost:8080/notes
```

---

## ✅ Checklist

- [ ] I can write a handler and register it with a method + path pattern
- [ ] I can read path parameters, query parameters and a JSON body
- [ ] I write headers, then status, then body, and `return` after errors
- [ ] I protect shared state in handlers with a mutex
- [ ] I can test handlers with `httptest.NewRequest` and `httptest.NewRecorder`

**Next → [Lesson 27: HTTP clients and middleware](../27_http_client/README.md)**
