package lesson28

import (
	"context"
	"errors"
	"path/filepath"
	"strings"
	"testing"
)

// newStore opens a fresh database file in a temporary folder for each test.
func newStore(t *testing.T) *Store {
	t.Helper()
	s, err := Open(context.Background(), filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("Open failed: %v (do Ex1 first)", err)
	}
	t.Cleanup(func() { s.Close() })
	return s
}

func TestEx1Open(t *testing.T) {
	ctx := context.Background()
	path := filepath.Join(t.TempDir(), "todos.db")
	s, err := Open(ctx, path)
	if err != nil {
		t.Fatalf("Ex1: Open returned %v", err)
	}
	var n int
	if err := s.db.QueryRowContext(ctx, "SELECT count(*) FROM todos").Scan(&n); err != nil {
		t.Errorf("Ex1: the todos table doesn't exist: %v", err)
	}
	s.Close()
	s2, err := Open(ctx, path) // opening an existing database must work too
	if err != nil {
		t.Fatalf("Ex1: opening the same file twice failed: %v (use CREATE TABLE IF NOT EXISTS)", err)
	}
	s2.Close()
}

func TestEx2Add(t *testing.T) {
	s := newStore(t)
	ctx := context.Background()
	id1, err := s.Add(ctx, "buy milk")
	if err != nil || id1 != 1 {
		t.Fatalf("Ex2: first Add = %d, %v; want 1, nil", id1, err)
	}
	id2, _ := s.Add(ctx, "walk dog")
	if id2 != 2 {
		t.Errorf("Ex2: second Add = %d, want 2", id2)
	}
	var title string
	var done bool
	if err := s.db.QueryRowContext(ctx, "SELECT title, done FROM todos WHERE id = 1").Scan(&title, &done); err != nil || title != "buy milk" || done {
		t.Errorf("Ex2: row 1 = %q, done=%v, err=%v; want \"buy milk\", false", title, done, err)
	}
}

func TestEx3Get(t *testing.T) {
	s := newStore(t)
	ctx := context.Background()
	s.db.ExecContext(ctx, "INSERT INTO todos (title, done) VALUES ('read book', 1)")
	got, err := s.Get(ctx, 1)
	if err != nil || got != (Todo{ID: 1, Title: "read book", Done: true}) {
		t.Errorf("Ex3: Get(1) = %+v, %v", got, err)
	}
	if _, err := s.Get(ctx, 99); !errors.Is(err, ErrNotFound) {
		t.Errorf("Ex3: Get(99) error = %v, want ErrNotFound", err)
	}
}

func TestEx4List(t *testing.T) {
	s := newStore(t)
	ctx := context.Background()
	if got, err := s.List(ctx); err != nil || got == nil || len(got) != 0 {
		t.Errorf("Ex4: List on an empty table = %#v, %v; want an empty non-nil slice", got, err)
	}
	for _, title := range []string{"a", "b", "c"} {
		s.db.ExecContext(ctx, "INSERT INTO todos (title) VALUES (?)", title)
	}
	got, err := s.List(ctx)
	if err != nil || len(got) != 3 || got[0].Title != "a" || got[2] != (Todo{ID: 3, Title: "c"}) {
		t.Errorf("Ex4: List = %+v, %v; want 3 todos ordered by id", got, err)
	}
}

func TestEx5MarkDone(t *testing.T) {
	s := newStore(t)
	ctx := context.Background()
	s.db.ExecContext(ctx, "INSERT INTO todos (title) VALUES ('x')")
	if err := s.MarkDone(ctx, 1); err != nil {
		t.Fatalf("Ex5: MarkDone(1) returned %v", err)
	}
	var done bool
	s.db.QueryRowContext(ctx, "SELECT done FROM todos WHERE id = 1").Scan(&done)
	if !done {
		t.Error("Ex5: todo 1 is still not done")
	}
	if err := s.MarkDone(ctx, 42); !errors.Is(err, ErrNotFound) {
		t.Errorf("Ex5: MarkDone(42) = %v, want ErrNotFound", err)
	}
}

func TestEx6Search(t *testing.T) {
	s := newStore(t)
	ctx := context.Background()
	for _, title := range []string{"Buy milk", "walk dog", "milkshake", "call mom"} {
		s.db.ExecContext(ctx, "INSERT INTO todos (title) VALUES (?)", title)
	}
	got, err := s.Search(ctx, "milk")
	if err != nil || len(got) != 2 || got[0].ID != 1 || got[1].ID != 3 {
		t.Errorf("Ex6: Search(milk) = %+v, %v; want todos 1 and 3", got, err)
	}
	// An SQL injection attempt must be treated as plain text, finding nothing and deleting nothing.
	got, err = s.Search(ctx, "x' OR '1'='1")
	if err != nil || len(got) != 0 {
		t.Errorf("Ex6: Search with an injection attempt returned %+v, %v; want nothing. Use a ? placeholder", got, err)
	}
	var n int
	s.db.QueryRowContext(ctx, "SELECT count(*) FROM todos").Scan(&n)
	if n != 4 {
		t.Errorf("Ex6: the table has %d rows after searching, want 4", n)
	}
}

func TestEx7AddAll(t *testing.T) {
	s := newStore(t)
	ctx := context.Background()
	if err := s.AddAll(ctx, []string{"one", "two", "three"}); err != nil {
		t.Fatalf("Ex7: AddAll(valid) returned %v", err)
	}
	err := s.AddAll(ctx, []string{"four", "five", "", "six"})
	if err == nil || !strings.Contains(err.Error(), "position 2") {
		t.Errorf("Ex7: AddAll with an empty title returned %v; want \"empty title at position 2\"", err)
	}
	var n int
	s.db.QueryRowContext(ctx, "SELECT count(*) FROM todos").Scan(&n)
	if n != 3 {
		t.Errorf("Ex7: the table has %d rows, want 3. \"four\" and \"five\" must be rolled back", n)
	}
}
