// Reference solution for lesson 28.
package lesson28

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	_ "modernc.org/sqlite"
)

var ErrNotFound = errors.New("todo not found")

type Todo struct {
	ID    int64
	Title string
	Done  bool
}

const schema = `
CREATE TABLE IF NOT EXISTS todos (
	id    INTEGER PRIMARY KEY AUTOINCREMENT,
	title TEXT    NOT NULL,
	done  BOOLEAN NOT NULL DEFAULT 0
)`

type Store struct {
	db *sql.DB
}

func (s *Store) Close() error {
	return s.db.Close()
}

func Open(ctx context.Context, path string) (*Store, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, err
	}
	if _, err := db.ExecContext(ctx, schema); err != nil {
		db.Close()
		return nil, err
	}
	return &Store{db: db}, nil
}

func (s *Store) Add(ctx context.Context, title string) (int64, error) {
	res, err := s.db.ExecContext(ctx, "INSERT INTO todos (title) VALUES (?)", title)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (s *Store) Get(ctx context.Context, id int64) (Todo, error) {
	var t Todo
	err := s.db.QueryRowContext(ctx, "SELECT id, title, done FROM todos WHERE id = ?", id).Scan(&t.ID, &t.Title, &t.Done)
	if errors.Is(err, sql.ErrNoRows) {
		return Todo{}, ErrNotFound
	}
	if err != nil {
		return Todo{}, err
	}
	return t, nil
}

func (s *Store) query(ctx context.Context, q string, args ...any) ([]Todo, error) {
	rows, err := s.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	todos := []Todo{}
	for rows.Next() {
		var t Todo
		if err := rows.Scan(&t.ID, &t.Title, &t.Done); err != nil {
			return nil, err
		}
		todos = append(todos, t)
	}
	return todos, rows.Err()
}

func (s *Store) List(ctx context.Context) ([]Todo, error) {
	return s.query(ctx, "SELECT id, title, done FROM todos ORDER BY id")
}

func (s *Store) MarkDone(ctx context.Context, id int64) error {
	res, err := s.db.ExecContext(ctx, "UPDATE todos SET done = 1 WHERE id = ?", id)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *Store) Search(ctx context.Context, text string) ([]Todo, error) {
	return s.query(ctx, "SELECT id, title, done FROM todos WHERE title LIKE ? ORDER BY id", "%"+text+"%")
}

func (s *Store) AddAll(ctx context.Context, titles []string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for i, title := range titles {
		if title == "" {
			return fmt.Errorf("empty title at position %d", i)
		}
		if _, err := tx.ExecContext(ctx, "INSERT INTO todos (title) VALUES (?)", title); err != nil {
			return err
		}
	}
	return tx.Commit()
}
