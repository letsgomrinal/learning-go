package lesson07

import (
	"maps"
	"slices"
	"testing"
)

func sampleLibrary() []Book {
	return []Book{
		{Title: "Dune", Author: "Frank Herbert", Pages: 412},
		{Title: "Emma", Author: "Jane Austen", Pages: 474},
		{Title: "Dune Messiah", Author: "Frank Herbert", Pages: 256},
		{Title: "Persuasion", Author: "Jane Austen", Pages: 474},
		{Title: "Neuromancer", Author: "William Gibson", Pages: 271},
	}
}

func TestEx1NewBook(t *testing.T) {
	got := NewBook("Dune", "Frank Herbert", 412)
	want := Book{Title: "Dune", Author: "Frank Herbert", Pages: 412}
	if got != want {
		t.Errorf("Ex1: NewBook = %+v, want %+v", got, want)
	}
}

func TestEx2Summary(t *testing.T) {
	got := Summary(Book{Title: "Dune", Author: "Frank Herbert", Pages: 412})
	want := "Dune by Frank Herbert (412 pages)"
	if got != want {
		t.Errorf("Ex2: Summary\n   got:  %q\n   want: %q", got, want)
	}
}

func TestEx3TotalPages(t *testing.T) {
	if got := TotalPages(sampleLibrary()); got != 1887 {
		t.Errorf("Ex3: TotalPages = %d, want 1887", got)
	}
	if got := TotalPages(nil); got != 0 {
		t.Errorf("Ex3: TotalPages(empty) = %d, want 0", got)
	}
}

func TestEx4LongestBook(t *testing.T) {
	got, ok := LongestBook(sampleLibrary())
	if !ok || got.Title != "Emma" {
		t.Errorf("Ex4: LongestBook = %+v, %v; want Emma (the first book with 474 pages), true", got, ok)
	}
	if got, ok := LongestBook(nil); ok || got != (Book{}) {
		t.Errorf("Ex4: LongestBook(empty) = %+v, %v; want {}, false", got, ok)
	}
}

func TestEx5BooksBy(t *testing.T) {
	got := BooksBy(sampleLibrary(), "Frank Herbert")
	var titles []string
	for _, b := range got {
		titles = append(titles, b.Title)
	}
	if want := []string{"Dune", "Dune Messiah"}; !slices.Equal(titles, want) {
		t.Errorf("Ex5: BooksBy(Frank Herbert) titles = %q, want %q", titles, want)
	}
	if got := BooksBy(sampleLibrary(), "Nobody"); len(got) != 0 {
		t.Errorf("Ex5: BooksBy(Nobody) = %v, want empty", got)
	}
}

func TestEx6Area(t *testing.T) {
	if got := Area(Rectangle{Width: 3, Height: 4}); got != 12 {
		t.Errorf("Ex6: Area({3 4}) = %v, want 12", got)
	}
}

func TestEx7Scaled(t *testing.T) {
	orig := Rectangle{Width: 2, Height: 3}
	got := Scaled(orig, 2)
	if want := (Rectangle{Width: 4, Height: 6}); got != want {
		t.Errorf("Ex7: Scaled({2 3}, 2) = %+v, want %+v", got, want)
	}
	if orig != (Rectangle{Width: 2, Height: 3}) {
		t.Errorf("Ex7: the original rectangle changed to %+v", orig)
	}
}

func TestEx8ResetPages(t *testing.T) {
	lib := sampleLibrary()
	ResetPages(lib)
	for _, b := range lib {
		if b.Pages != 0 {
			t.Fatalf("Ex8: after ResetPages, %q still has %d pages. Did you change a copy?", b.Title, b.Pages)
		}
	}
}

func TestEx9CountByAuthor(t *testing.T) {
	got := CountByAuthor(sampleLibrary())
	want := map[string]int{"Frank Herbert": 2, "Jane Austen": 2, "William Gibson": 1}
	if !maps.Equal(got, want) {
		t.Errorf("Ex9: CountByAuthor = %v, want %v", got, want)
	}
}
