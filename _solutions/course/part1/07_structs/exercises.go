// Reference solution for lesson 07.
package lesson07

import "fmt"

type Book struct {
	Title  string
	Author string
	Pages  int
}

type Rectangle struct {
	Width, Height float64
}

func NewBook(title, author string, pages int) Book {
	return Book{Title: title, Author: author, Pages: pages}
}

func Summary(b Book) string {
	return fmt.Sprintf("%s by %s (%d pages)", b.Title, b.Author, b.Pages)
}

func TotalPages(books []Book) int {
	total := 0
	for _, b := range books {
		total += b.Pages
	}
	return total
}

func LongestBook(books []Book) (Book, bool) {
	if len(books) == 0 {
		return Book{}, false
	}
	longest := books[0]
	for _, b := range books {
		if b.Pages > longest.Pages { // > (not >=) keeps the first one on a tie
			longest = b
		}
	}
	return longest, true
}

func BooksBy(books []Book, author string) []Book {
	var result []Book
	for _, b := range books {
		if b.Author == author {
			result = append(result, b)
		}
	}
	return result
}

func Area(r Rectangle) float64 {
	return r.Width * r.Height
}

func Scaled(r Rectangle, factor float64) Rectangle {
	r.Width *= factor
	r.Height *= factor
	return r
}

func ResetPages(books []Book) {
	for i := range books {
		books[i].Pages = 0
	}
}

func CountByAuthor(books []Book) map[string]int {
	counts := map[string]int{}
	for _, b := range books {
		counts[b.Author]++
	}
	return counts
}
