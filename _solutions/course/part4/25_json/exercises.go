// Reference solution for lesson 25.
package lesson25

import (
	"encoding/json"
	"errors"
	"io"
	"strings"
)

type Book struct {
	Title      string   `json:"title"`
	Subtitle   string   `json:"subtitle,omitempty"`
	Author     string   `json:"author"`
	Year       int      `json:"year"`
	Tags       []string `json:"tags"`
	InternalID int      `json:"-"`
}

func BookJSON(b Book) (string, error) {
	data, err := json.Marshal(b)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func ParseBook(text string) (Book, error) {
	var b Book
	err := json.Unmarshal([]byte(text), &b)
	return b, err
}

type Order struct {
	ID     int     `json:"id"`
	Amount float64 `json:"amount"`
	Paid   bool    `json:"paid"`
}

func TotalSpent(text string) (float64, error) {
	var orders []Order
	if err := json.Unmarshal([]byte(text), &orders); err != nil {
		return 0, err
	}
	total := 0.0
	for _, o := range orders {
		if o.Paid {
			total += o.Amount
		}
	}
	return total, nil
}

type Profile struct {
	Name     string   `json:"full_name"`
	Email    string   `json:"email_address"`
	Home     Location `json:"home"`
	Verified bool     `json:"is_verified"`
}

type Location struct {
	City       string `json:"city"`
	PostalCode string `json:"postal_code"`
}

func ParseProfile(text string) (Profile, error) {
	var p Profile
	err := json.Unmarshal([]byte(text), &p)
	return p, err
}

func Pretty(v any) (string, error) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func Lookup(text, key string) (any, bool) {
	var data map[string]any
	if err := json.Unmarshal([]byte(text), &data); err != nil {
		return nil, false
	}
	v, ok := data[key]
	return v, ok
}

type Event struct {
	Level   string `json:"level"`
	Message string `json:"msg"`
}

func ReadEvents(r io.Reader) ([]Event, error) {
	dec := json.NewDecoder(r)
	var events []Event
	for {
		var e Event
		err := dec.Decode(&e)
		if errors.Is(err, io.EOF) {
			return events, nil
		}
		if err != nil {
			return nil, err
		}
		events = append(events, e)
	}
}

type Config struct {
	Port    int    `json:"port"`
	Host    string `json:"host"`
	Verbose bool   `json:"verbose"`
}

func LoadConfig(text string) (Config, error) {
	dec := json.NewDecoder(strings.NewReader(text))
	dec.DisallowUnknownFields()
	var cfg Config
	if err := dec.Decode(&cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}
