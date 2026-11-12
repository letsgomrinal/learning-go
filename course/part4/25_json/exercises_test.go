package lesson25

import (
	"strings"
	"testing"
)

func TestEx1BookJSON(t *testing.T) {
	b := Book{Title: "Dune", Author: "Frank Herbert", Year: 1965, Tags: []string{"sci-fi"}, InternalID: 99}
	got, err := BookJSON(b)
	want := `{"title":"Dune","author":"Frank Herbert","year":1965,"tags":["sci-fi"]}`
	if err != nil || got != want {
		t.Errorf("Ex1: BookJSON\n   got:  %s (err %v)\n   want: %s", got, err, want)
	}
	b.Subtitle = "Book One"
	got, _ = BookJSON(b)
	if !strings.Contains(got, `"subtitle":"Book One"`) {
		t.Errorf("Ex1: with a subtitle, got %s; want it to include \"subtitle\":\"Book One\"", got)
	}
	if strings.Contains(got, "99") || strings.Contains(strings.ToLower(got), "internal") {
		t.Errorf("Ex1: InternalID leaked into the JSON: %s", got)
	}
}

func TestEx2ParseBook(t *testing.T) {
	b, err := ParseBook(`{"title": "Emma", "author": "Jane Austen", "year": 1815, "tags": ["classic", "romance"]}`)
	if err != nil || b.Title != "Emma" || b.Author != "Jane Austen" || b.Year != 1815 || len(b.Tags) != 2 {
		t.Errorf("Ex2: ParseBook = %+v, %v", b, err)
	}
	if _, err := ParseBook(`{"title": 42}`); err == nil {
		t.Error("Ex2: a number where a string belongs should be an error")
	}
	if _, err := ParseBook(`not json`); err == nil {
		t.Error("Ex2: invalid JSON should be an error")
	}
}

func TestEx3TotalSpent(t *testing.T) {
	text := `[{"id": 1, "amount": 9.5, "paid": true}, {"id": 2, "amount": 20, "paid": false}, {"id": 3, "amount": 0.5, "paid": true}]`
	if got, err := TotalSpent(text); got != 10 || err != nil {
		t.Errorf("Ex3: TotalSpent = %v, %v; want 10, nil", got, err)
	}
	if got, err := TotalSpent(`[]`); got != 0 || err != nil {
		t.Errorf("Ex3: TotalSpent([]) = %v, %v; want 0, nil", got, err)
	}
	if _, err := TotalSpent(`[{"amount": "lots"}]`); err == nil {
		t.Error("Ex3: invalid data should return an error")
	}
}

func TestEx4Profile(t *testing.T) {
	text := `{
	  "full_name": "Ana Silva",
	  "email_address": "ana@example.com",
	  "home": {"city": "Lisbon", "postal_code": "1000-001"},
	  "is_verified": true
	}`
	p, err := ParseProfile(text)
	want := Profile{Name: "Ana Silva", Email: "ana@example.com", Home: Location{City: "Lisbon", PostalCode: "1000-001"}, Verified: true}
	if err != nil || p != want {
		t.Errorf("Ex4: ParseProfile\n   got:  %+v (err %v)\n   want: %+v", p, err, want)
	}
}

func TestEx5Pretty(t *testing.T) {
	got, err := Pretty(map[string]any{"b": []int{1, 2}, "a": 1})
	want := "{\n  \"a\": 1,\n  \"b\": [\n    1,\n    2\n  ]\n}"
	if err != nil || got != want {
		t.Errorf("Ex5: Pretty\n   got:\n%s\n   want:\n%s", got, want)
	}
}

func TestEx6Lookup(t *testing.T) {
	text := `{"name": "Ana", "age": 30, "langs": ["go", "py"], "admin": false}`
	if v, ok := Lookup(text, "name"); v != "Ana" || !ok {
		t.Errorf("Ex6: Lookup(name) = %v, %v; want \"Ana\", true", v, ok)
	}
	if v, ok := Lookup(text, "age"); v != float64(30) || !ok {
		t.Errorf("Ex6: Lookup(age) = %v (%T), %v; want float64(30), true", v, v, ok)
	}
	if v, ok := Lookup(text, "admin"); v != false || !ok {
		t.Errorf("Ex6: Lookup(admin) = %v, %v; want false, true", v, ok)
	}
	if v, ok := Lookup(text, "missing"); v != nil || ok {
		t.Errorf("Ex6: Lookup(missing) = %v, %v; want nil, false", v, ok)
	}
	if _, ok := Lookup(`{broken`, "name"); ok {
		t.Error("Ex6: invalid JSON should give ok == false")
	}
}

func TestEx7ReadEvents(t *testing.T) {
	input := `{"level": "info", "msg": "started"}
{"level": "warn", "msg": "disk 90% full"}
{"level": "info", "msg": "stopped"}
`
	got, err := ReadEvents(strings.NewReader(input))
	if err != nil || len(got) != 3 || got[1] != (Event{Level: "warn", Message: "disk 90% full"}) {
		t.Errorf("Ex7: ReadEvents = %+v, %v; want 3 events, the second {warn disk 90%% full}", got, err)
	}
	if got, err := ReadEvents(strings.NewReader("")); err != nil || len(got) != 0 {
		t.Errorf("Ex7: empty input = %v, %v; want no events, nil", got, err)
	}
	if _, err := ReadEvents(strings.NewReader(`{"level": "info"} {oops`)); err == nil {
		t.Error("Ex7: broken JSON should return an error")
	}
}

func TestEx8LoadConfig(t *testing.T) {
	cfg, err := LoadConfig(`{"port": 8080, "host": "localhost"}`)
	if err != nil || cfg != (Config{Port: 8080, Host: "localhost"}) {
		t.Errorf("Ex8: LoadConfig(valid) = %+v, %v", cfg, err)
	}
	_, err = LoadConfig(`{"prot": 8080}`)
	if err == nil || !strings.Contains(err.Error(), "unknown field") {
		t.Errorf("Ex8: LoadConfig with a typo returned %v; want an \"unknown field\" error", err)
	}
}
