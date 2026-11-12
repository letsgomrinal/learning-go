// Lesson 09: An HTTP JSON API with only the standard library.
//
// Run:  go test -race -v ./practice/09_http_json
// Try it for real:  go run ./practice/09_http_json/cmd   then
//                   curl -X POST localhost:8080/todos -d '{"title":"learn go"}'
//                   curl localhost:8080/todos
//
// ─── CONCEPTS ──────────────────────────────────────────────────────────────
//
// 1. net/http is production-grade. Many Go services use NO framework (no Flask/FastAPI
//    equivalent needed). The core interface:
//      type Handler interface { ServeHTTP(w http.ResponseWriter, r *http.Request) }
//    http.HandlerFunc lets a plain function act as a Handler.
//
// 2. ROUTING (Go 1.22+ ServeMux supports methods and path params):
//      mux := http.NewServeMux()
//      mux.HandleFunc("GET /todos/{id}", s.getTodo)
//      id := r.PathValue("id")       // a string; convert with strconv.Atoi
//    Wrong method → 405 automatically. Unknown path → 404.
//
// 3. JSON with struct tags (like pydantic field aliases):
//      type Todo struct {
//          ID    int    `json:"id"`
//          Title string `json:"title"`
//          Notes string `json:"notes,omitempty"`  // left out when empty
//      }
//    Only EXPORTED (capitalized) fields are encoded or decoded!
//      json.NewDecoder(r.Body).Decode(&t)      // request body → struct
//      json.NewEncoder(w).Encode(v)            // struct → response
//
// 4. WRITING RESPONSES: set headers FIRST, then WriteHeader(status), then the body.
//      w.Header().Set("Content-Type", "application/json")
//      w.WriteHeader(http.StatusCreated)
//      json.NewEncoder(w).Encode(todo)
//    Calling Write without WriteHeader implies 200.
//
// 5. HANDLERS RUN CONCURRENTLY: every request gets its own goroutine. Shared state
//    (our in-memory store) MUST be protected, e.g. with a sync.RWMutex.
//
// 6. TESTING: net/http/httptest gives you a fake ResponseWriter (NewRecorder) and a
//    real local server (NewServer). No need for requests-mock.
// ───────────────────────────────────────────────────────────────────────────

package todoapi

import (
	"net/http"
)

type Todo struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

// Server is an in-memory todo API.
type Server struct {
	// TODO: a sync.RWMutex, a map[int]Todo (or slice), and a nextID counter
	mux *http.ServeMux
}

// NewServer wires up the routes:
//
//	GET    /todos        → 200, JSON array of all todos sorted by ID (an empty list must be [] not null)
//	POST   /todos        → body {"title": "..."}; 201 with the created todo (ids start at 1)
//	                        400 if the JSON is invalid or the title is empty/whitespace
//	GET    /todos/{id}   → 200 with the todo; 404 if missing; 400 if id isn't a number
//	PATCH  /todos/{id}   → body {"done": true}; 200 with the updated todo; 404 if missing
//	DELETE /todos/{id}   → 204 No Content; 404 if missing
//
// Error responses are JSON:  {"error": "message"}
func NewServer() *Server {
	s := &Server{mux: http.NewServeMux()}
	// TODO: register routes, e.g. s.mux.HandleFunc("GET /todos", s.listTodos)
	return s
}

// ServeHTTP makes *Server an http.Handler.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

// writeJSON is a helper you'll want to use in every handler.
func writeJSON(w http.ResponseWriter, status int, v any) {
	// TODO: set Content-Type, write status, encode v
}

// TODO: write the handlers:
//   func (s *Server) listTodos(w http.ResponseWriter, r *http.Request)
//   func (s *Server) createTodo(w http.ResponseWriter, r *http.Request)
//   func (s *Server) getTodo(w http.ResponseWriter, r *http.Request)
//   func (s *Server) updateTodo(w http.ResponseWriter, r *http.Request)
//   func (s *Server) deleteTodo(w http.ResponseWriter, r *http.Request)
//
// Tip for PATCH: decode into a struct with a POINTER field so you can tell
// "not provided" (nil) from "false":   var body struct{ Done *bool `json:"done"` }
