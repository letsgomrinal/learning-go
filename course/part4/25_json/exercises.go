// Lesson 25 exercises: JSON. Read README.md first!
//
// Check one exercise at a time:   go test -run Ex1 -v
package lesson25

import (
	"encoding/json"
	"io"
)

// Ex1: add struct tags so a Book turns into JSON like this:
//
//		{"title":"Dune","author":"Frank Herbert","year":1965,"tags":["sci-fi"]}
//
//	  - keys are lowercase: title, author, year, subtitle, tags
//	  - Subtitle is left out when it's empty (omitempty)
//	  - InternalID must NEVER appear in the JSON (json:"-")
//
// Then write BookJSON below.
type Book struct {
	Title      string
	Subtitle   string
	Author     string
	Year       int
	Tags       []string
	InternalID int
}

// BookJSON returns the book as a JSON string. Hint: data, err := json.Marshal(b)
func BookJSON(b Book) (string, error) {
	// TODO
	return "", nil
}

// Ex2: parse a JSON object into a Book. Return any error from json.Unmarshal.
// (Don't forget the &.)
func ParseBook(text string) (Book, error) {
	// TODO
	return Book{}, nil
}

// Order is one order in a JSON array like:
//
//	[{"id": 1, "amount": 9.5, "paid": true}, {"id": 2, "amount": 20, "paid": false}]
type Order struct {
	ID     int     `json:"id"`
	Amount float64 `json:"amount"`
	Paid   bool    `json:"paid"`
}

// Ex3: decode the JSON array of orders, and return the sum of the amounts of the PAID orders.
// Return 0 and the error if the JSON is invalid.
func TotalSpent(text string) (float64, error) {
	// TODO
	return 0, nil
}

// Ex4: an API returns profiles like this (look carefully at the key names):
//
//	{
//	  "full_name": "Ana Silva",
//	  "email_address": "ana@example.com",
//	  "home": {"city": "Lisbon", "postal_code": "1000-001"},
//	  "is_verified": true
//	}
//
// Add struct tags to Profile and Location so json.Unmarshal fills every field.
// (Nothing else to write: ParseProfile is already done.)
type Profile struct {
	Name     string
	Email    string
	Home     Location
	Verified bool
}

type Location struct {
	City       string
	PostalCode string
}

func ParseProfile(text string) (Profile, error) {
	var p Profile
	err := json.Unmarshal([]byte(text), &p)
	return p, err
}

// Ex5: return v as indented JSON: no prefix, two spaces per level. (README section 6.)
func Pretty(v any) (string, error) {
	// TODO
	return "", nil
}

// Ex6: decode a JSON object of unknown shape into a map[string]any, and return the value under key.
// Return nil and false if the key is missing or the JSON is invalid.
//
//	Lookup(`{"name": "Ana", "age": 30}`, "age") == float64(30), true   ← numbers are float64!
func Lookup(text, key string) (any, bool) {
	// TODO
	return nil, false
}

// Event is one line of a JSON Lines log.
type Event struct {
	Level   string `json:"level"`
	Message string `json:"msg"`
}

// Ex7: read every JSON object from r (one after another, like JSON Lines), and return them all.
// Stop at io.EOF. If any value is invalid JSON, return nil and the error. (README section 8.)
func ReadEvents(r io.Reader) ([]Event, error) {
	// TODO
	return nil, nil
}

// Config is an app's configuration file.
type Config struct {
	Port    int    `json:"port"`
	Host    string `json:"host"`
	Verbose bool   `json:"verbose"`
}

// Ex8: decode a config, REJECTING unknown keys (so a typo like "prot" is an error, not silently
// ignored). Use a json.Decoder with DisallowUnknownFields. Return Config{} and the error on failure.
//
// Hint: dec := json.NewDecoder(strings.NewReader(text))   (import "strings")
func LoadConfig(text string) (Config, error) {
	// TODO
	return Config{}, nil
}
