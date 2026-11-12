// Lesson 07 exercises: structs. Read README.md first!
//
// Check one exercise at a time:   go test -run Ex1 -v
package lesson07

// Book and Rectangle are already defined for you.

type Book struct {
	Title  string
	Author string
	Pages  int
}

type Rectangle struct {
	Width, Height float64
}

// Ex1: create and return a Book with the given title, author and pages.
func NewBook(title, author string, pages int) Book {
	// TODO
	return Book{}
}

// Ex2: return a one-line summary.
// Summary(Book{Title: "Dune", Author: "Frank Herbert", Pages: 412}) == "Dune by Frank Herbert (412 pages)"
//
// Hint: fmt.Sprintf with b.Title, b.Author, b.Pages. Remember to import "fmt".
func Summary(b Book) string {
	// TODO
	return ""
}

// Ex3: return the total number of pages in all books.
func TotalPages(books []Book) int {
	// TODO
	return 0
}

// Ex4: return the book with the most pages, and true.
// For an empty library return Book{} (an empty book) and false.
// If two books have the same number of pages, return the one that comes first.
func LongestBook(books []Book) (Book, bool) {
	// TODO
	return Book{}, false
}

// Ex5: return the books written by author, in their original order.
func BooksBy(books []Book, author string) []Book {
	// TODO
	return nil
}

// Ex6: return the area of the rectangle (width × height).
func Area(r Rectangle) float64 {
	// TODO
	return 0
}

// Ex7: return a NEW rectangle with both sides multiplied by factor.
// Because r is a copy, you can change r and return it, and the caller's rectangle stays the same.
// Scaled(Rectangle{Width: 2, Height: 3}, 2) == Rectangle{Width: 4, Height: 6}
func Scaled(r Rectangle, factor float64) Rectangle {
	// TODO
	return Rectangle{}
}

// Ex8: set Pages to 0 for every book IN the slice (the caller should see the change).
//
// Hint: `for _, b := range books { b.Pages = 0 }` does NOT work. Why? Use the index instead.
func ResetPages(books []Book) {
	// TODO
}

// Ex9: count how many books each author has written.
// For two books by Frank Herbert and one by Jane Austen, the result is
// map["Frank Herbert":2 "Jane Austen":1]
func CountByAuthor(books []Book) map[string]int {
	// TODO
	return nil
}
