// Reference solution for lesson 13.
package lesson13

import (
	"errors"
	"fmt"
	"strings"
)

var ErrFileNotFound = errors.New("file not found")

var files = map[string]string{
	"app.yaml": "port: 8080",
}

func ReadFile(name string) (string, error) {
	content, ok := files[name]
	if !ok {
		return "", fmt.Errorf("open %s: %w", name, ErrFileNotFound)
	}
	return content, nil
}

func LoadConfig(name string) (string, error) {
	content, err := ReadFile(name)
	if err != nil {
		return "", fmt.Errorf("load config %s: %w", name, err)
	}
	return content, nil
}

type ValidationError struct {
	Field   string
	Problem string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("invalid %s: %s", e.Field, e.Problem)
}

func ValidateUser(name, email string) error {
	if name == "" {
		return &ValidationError{Field: "name", Problem: "is required"}
	}
	if !strings.Contains(email, "@") {
		return &ValidationError{Field: "email", Problem: "missing @"}
	}
	return nil
}

func RegisterUser(name, email string) error {
	if err := ValidateUser(name, email); err != nil {
		return fmt.Errorf("register %s: %w", name, err)
	}
	return nil
}

func BadField(err error) string {
	var ve *ValidationError
	if errors.As(err, &ve) {
		return ve.Field
	}
	return ""
}

type NotFoundError struct {
	ID int
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("user %d not found", e.ID)
}

// Fixed: return a literal nil on success.
func FindUser(id int) error {
	if id < 0 {
		return &NotFoundError{ID: id}
	}
	return nil
}

func StatusCode(err error) int {
	var ve *ValidationError
	switch {
	case err == nil:
		return 200
	case errors.Is(err, ErrFileNotFound):
		return 404
	case errors.As(err, &ve):
		return 400
	default:
		return 500
	}
}
