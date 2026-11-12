// Reference solution for lesson 09.
package todoapi

import (
	"encoding/json"
	"maps"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"sync"
)

type Todo struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

type Server struct {
	mu     sync.RWMutex
	todos  map[int]Todo
	nextID int
	mux    *http.ServeMux
}

func NewServer() *Server {
	s := &Server{
		todos:  make(map[int]Todo),
		nextID: 1,
		mux:    http.NewServeMux(),
	}
	s.mux.HandleFunc("GET /todos", s.listTodos)
	s.mux.HandleFunc("POST /todos", s.createTodo)
	s.mux.HandleFunc("GET /todos/{id}", s.getTodo)
	s.mux.HandleFunc("PATCH /todos/{id}", s.updateTodo)
	s.mux.HandleFunc("DELETE /todos/{id}", s.deleteTodo)
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

// pathID parses {id}, writing a 400 and returning ok=false if it's invalid.
func pathID(w http.ResponseWriter, r *http.Request) (int, bool) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "id must be an integer")
		return 0, false
	}
	return id, true
}

func (s *Server) listTodos(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	ids := slices.Sorted(maps.Keys(s.todos))
	out := make([]Todo, 0, len(ids)) // non-nil, so it encodes as [] instead of null
	for _, id := range ids {
		out = append(out, s.todos[id])
	}
	s.mu.RUnlock()
	writeJSON(w, http.StatusOK, out)
}

func (s *Server) createTodo(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Title string `json:"title"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}
	title := strings.TrimSpace(body.Title)
	if title == "" {
		writeError(w, http.StatusBadRequest, "title is required")
		return
	}

	s.mu.Lock()
	t := Todo{ID: s.nextID, Title: title}
	s.todos[t.ID] = t
	s.nextID++
	s.mu.Unlock()

	writeJSON(w, http.StatusCreated, t)
}

func (s *Server) getTodo(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	s.mu.RLock()
	t, found := s.todos[id]
	s.mu.RUnlock()
	if !found {
		writeError(w, http.StatusNotFound, "todo not found")
		return
	}
	writeJSON(w, http.StatusOK, t)
}

func (s *Server) updateTodo(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	var body struct {
		Done *bool `json:"done"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}

	s.mu.Lock()
	t, found := s.todos[id]
	if found && body.Done != nil {
		t.Done = *body.Done
		s.todos[id] = t
	}
	s.mu.Unlock()

	if !found {
		writeError(w, http.StatusNotFound, "todo not found")
		return
	}
	writeJSON(w, http.StatusOK, t)
}

func (s *Server) deleteTodo(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	s.mu.Lock()
	_, found := s.todos[id]
	delete(s.todos, id)
	s.mu.Unlock()

	if !found {
		writeError(w, http.StatusNotFound, "todo not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
