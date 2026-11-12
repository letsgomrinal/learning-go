// Lesson 28 exercises: a todo store with database/sql. Read README.md first!
//
// This folder is its own module (see go.mod). Run tests from THIS folder:
//
//	go test -race -run Ex1 -v
package lesson28

import (
	"context"
	"database/sql"
	"errors"

	_ "modernc.org/sqlite" // registers the "sqlite" driver (README section 1)
)

// ErrNotFound is returned when no todo has the requested id.
var ErrNotFound = errors.New("todo not found")

// Todo is one row of the todos table.
type Todo struct {
	ID    int64
	Title string
	Done  bool
}

// schema creates the table if it doesn't exist yet.
const schema = `
CREATE TABLE IF NOT EXISTS todos (
	id    INTEGER PRIMARY KEY AUTOINCREMENT,
	title TEXT    NOT NULL,
	done  BOOLEAN NOT NULL DEFAULT 0
)`

// Store wraps the database. The rest of a program uses its methods and never sees SQL.
type Store struct {
	db *sql.DB
}

// Close closes the database. Already done.
func (s *Store) Close() error {
	return s.db.Close()
}

// Ex1: open the SQLite database file at path and prepare it.
//  1. db, err := sql.Open("sqlite", path)
//  2. db.PingContext(ctx) to check the connection works
//  3. db.ExecContext(ctx, schema) to create the table
//  4. return &Store{db: db}
//
// If step 2 or 3 fails, close db before returning the error.
func Open(ctx context.Context, path string) (*Store, error) {
	// TODO
	return nil, errors.New("not implemented")
}

// Ex2: insert a new todo (not done) and return its id.
// Use a placeholder for the title (README section 3), and res.LastInsertId().
func (s *Store) Add(ctx context.Context, title string) (int64, error) {
	// TODO
	return 0, nil
}

// Ex3: return the todo with this id.
// If there's no such row, return Todo{} and ErrNotFound (translate sql.ErrNoRows, README section 4).
func (s *Store) Get(ctx context.Context, id int64) (Todo, error) {
	// TODO
	return Todo{}, nil
}

// Ex4: return all todos, ordered by id. With no todos, return an empty (non-nil) slice.
// Remember: defer rows.Close(), and return rows.Err() at the end. (README section 5.)
func (s *Store) List(ctx context.Context) ([]Todo, error) {
	// TODO
	return nil, nil
}

// Ex5: mark the todo as done. If no row was changed (no such id), return ErrNotFound.
// Hint: res.RowsAffected()
func (s *Store) MarkDone(ctx context.Context, id int64) error {
	// TODO
	return nil
}

// Ex6: return the todos whose title contains text (case-insensitive in SQLite for ASCII), ordered by id.
// Use LIKE with a placeholder:  WHERE title LIKE ?  and pass "%" + text + "%" as the argument.
// NEVER put text into the SQL string itself. The test tries an SQL injection!
func (s *Store) Search(ctx context.Context, text string) ([]Todo, error) {
	// TODO
	return nil, nil
}

// Ex7: add all the titles in ONE transaction. If any title is empty, add NONE of them:
// return an error "empty title at position <i>" and roll back. (README section 6.)
//
//	tx, err := s.db.BeginTx(ctx, nil)
//	defer tx.Rollback()
//	... tx.ExecContext for each title ...
//	return tx.Commit()
func (s *Store) AddAll(ctx context.Context, titles []string) error {
	// TODO
	return nil
}
