// Lesson 13 exercises: custom errors and wrapping. Read README.md first!
//
// Check one exercise at a time:   go test -run Ex1 -v
package lesson13

import (
	"errors"
	"fmt"
)

// ── A pretend file system, so the exercises don't need real files ─────────

var ErrFileNotFound = errors.New("file not found")

var files = map[string]string{
	"app.yaml": "port: 8080",
}

// ReadFile is already written. It returns the file's content, or a wrapped ErrFileNotFound.
func ReadFile(name string) (string, error) {
	content, ok := files[name]
	if !ok {
		return "", fmt.Errorf("open %s: %w", name, ErrFileNotFound)
	}
	return content, nil
}

// Ex1: read the config file with ReadFile. If it fails, return "" and the error WRAPPED
// with the context "load config <name>: ".
//
// LoadConfig("missing.yaml") should return an error whose message is:
//
//	load config missing.yaml: open missing.yaml: file not found
//
// and errors.Is(err, ErrFileNotFound) must still be true. (Which verb keeps the chain?)
func LoadConfig(name string) (string, error) {
	// TODO
	return "", nil
}

// ── Custom error type ─────────────────────────────────────────────────────

// ValidationError says which field was invalid, and why.
type ValidationError struct {
	Field   string
	Problem string
}

// Ex2: return "invalid <field>: <problem>", for example "invalid email: missing @".
func (e *ValidationError) Error() string {
	// TODO
	return ""
}

// Ex3: validate a new user's details. Check in this order:
//   - name is ""                → return &ValidationError{Field: "name", Problem: "is required"}
//   - email doesn't contain "@" → return &ValidationError{Field: "email", Problem: "missing @"}
//   - otherwise                 → return nil
//
// You'll need strings.Contains (import "strings").
func ValidateUser(name, email string) error {
	// TODO
	return nil
}

// Ex4: register a user. Call ValidateUser; if it fails, return the error wrapped as
// "register <name>: <the error>". Otherwise return nil.
//
// RegisterUser("ana", "nope") → "register ana: invalid email: missing @"
func RegisterUser(name, email string) error {
	// TODO
	return nil
}

// Ex5: if err (or anything wrapped inside it) is a *ValidationError, return its Field.
// Otherwise return "".
//
// Hint: var ve *ValidationError; if errors.As(err, &ve) { ... }
func BadField(err error) string {
	// TODO
	return ""
}

// ── Ex6: the nil interface trap ───────────────────────────────────────────

type NotFoundError struct {
	ID int
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("user %d not found", e.ID)
}

// FindUser should return nil for valid ids (0 or more), and a *NotFoundError for negative ids.
// BUG: it has the nil interface bug from README section 6, so callers think valid ids failed.
// Run the test to see it, then fix it.
func FindUser(id int) error {
	var err *NotFoundError
	if id < 0 {
		err = &NotFoundError{ID: id}
	}
	return err
}

// Ex7: choose the HTTP status code for an error:
//   - nil                                  → 200
//   - ErrFileNotFound anywhere in the chain → 404
//   - a *ValidationError anywhere in chain → 400
//   - anything else                        → 500
func StatusCode(err error) int {
	// TODO
	return 0
}
