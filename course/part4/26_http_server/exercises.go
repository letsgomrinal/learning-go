// Lesson 26 exercises: a notes API with net/http. Read README.md first!
//
// Check one exercise at a time:   go test -race -run Ex1 -v
// Run the real server:            go run ./cmd/notes
package lesson26

import (
	"net/http"
	"sync"
)

// Note is one note. It's sent and received as JSON.
type Note struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}

// Server holds the notes in memory. Handlers run concurrently, so mu protects notes and nextID.
type Server struct {
	mu     sync.Mutex
	notes  map[int]Note
	nextID int
	mux    *http.ServeMux
}

// NewServer creates the server and registers the routes.
// TODO: register each route as you write its handler, for example:
//
//	s.mux.HandleFunc("GET /health", s.health)
func NewServer() *Server {
	s := &Server{
		notes:  map[int]Note{},
		nextID: 1,
		mux:    http.NewServeMux(),
	}
	// TODO: routes go here
	return s
}

// ServeHTTP makes *Server an http.Handler (README section 7). Already done.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

// Ex1: write v as JSON with the given status code.
// Set the "Content-Type" header to "application/json" FIRST, then WriteHeader(status), then encode.
// (You'll need "encoding/json".)
func writeJSON(w http.ResponseWriter, status int, v any) {
	// TODO
}

// Ex1: write {"error": msg} with the given status code, using writeJSON.
func writeError(w http.ResponseWriter, status int, msg string) {
	// TODO
}

// Ex2: GET /health → 200 with {"status":"ok"}
// Write the handler as a method:  func (s *Server) health(w http.ResponseWriter, r *http.Request)
// and register it in NewServer.

// Ex3: POST /notes
//   - body: {"text": "buy milk"}
//   - invalid JSON → 400 {"error": "..."}
//   - text empty or only spaces → 400 {"error": "..."}
//   - otherwise: give it the next id (starting at 1), store it, respond 201 with the note as JSON
//     (store the text trimmed of surrounding spaces)
//
// Lock s.mu while you touch s.notes and s.nextID.

// Ex4: GET /notes/{id}
//   - id isn't a number → 400
//   - no note with that id → 404
//   - otherwise → 200 with the note
//
// Hint: id, err := strconv.Atoi(r.PathValue("id"))

// Ex5: GET /notes → 200 with a JSON array of all notes, sorted by id.
// With no notes, send [] (not null): start from notes := []Note{}.
// Hint: collect them under the lock, then slices.SortFunc(notes, func(a, b Note) int { return a.ID - b.ID })

// Ex6: DELETE /notes/{id}
//   - not a number → 400, missing → 404
//   - otherwise delete it and respond with w.WriteHeader(http.StatusNoContent) (204) and NO body

// Ex7: extend GET /notes: if the query parameter q is set (GET /notes?q=milk), only include notes
// whose text contains q, ignoring upper/lower case.
// Hint: q := r.URL.Query().Get("q"); strings.Contains(strings.ToLower(n.Text), strings.ToLower(q))
