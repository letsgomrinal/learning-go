// Reference solution for lesson 26.
package lesson26

import (
	"encoding/json"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"sync"
)

type Note struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}

type Server struct {
	mu     sync.Mutex
	notes  map[int]Note
	nextID int
	mux    *http.ServeMux
}

func NewServer() *Server {
	s := &Server{
		notes:  map[int]Note{},
		nextID: 1,
		mux:    http.NewServeMux(),
	}
	s.mux.HandleFunc("GET /health", s.health)
	s.mux.HandleFunc("POST /notes", s.createNote)
	s.mux.HandleFunc("GET /notes", s.listNotes)
	s.mux.HandleFunc("GET /notes/{id}", s.getNote)
	s.mux.HandleFunc("DELETE /notes/{id}", s.deleteNote)
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

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) createNote(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Text string `json:"text"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	text := strings.TrimSpace(input.Text)
	if text == "" {
		writeError(w, http.StatusBadRequest, "text is required")
		return
	}
	s.mu.Lock()
	note := Note{ID: s.nextID, Text: text}
	s.notes[note.ID] = note
	s.nextID++
	s.mu.Unlock()
	writeJSON(w, http.StatusCreated, note)
}

// pathID reads {id}, writing a 400 and returning ok=false if it isn't a number.
func pathID(w http.ResponseWriter, r *http.Request) (int, bool) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "id must be a number")
		return 0, false
	}
	return id, true
}

func (s *Server) getNote(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	s.mu.Lock()
	note, found := s.notes[id]
	s.mu.Unlock()
	if !found {
		writeError(w, http.StatusNotFound, "note not found")
		return
	}
	writeJSON(w, http.StatusOK, note)
}

func (s *Server) listNotes(w http.ResponseWriter, r *http.Request) {
	q := strings.ToLower(r.URL.Query().Get("q"))
	notes := []Note{}
	s.mu.Lock()
	for _, n := range s.notes {
		if q == "" || strings.Contains(strings.ToLower(n.Text), q) {
			notes = append(notes, n)
		}
	}
	s.mu.Unlock()
	slices.SortFunc(notes, func(a, b Note) int { return a.ID - b.ID })
	writeJSON(w, http.StatusOK, notes)
}

func (s *Server) deleteNote(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	s.mu.Lock()
	_, found := s.notes[id]
	delete(s.notes, id)
	s.mu.Unlock()
	if !found {
		writeError(w, http.StatusNotFound, "note not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
