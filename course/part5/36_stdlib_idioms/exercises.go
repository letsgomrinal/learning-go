// Lesson 36 exercises: idioms from the standard library. Read README.md first!
//
// Check one exercise at a time:   go test -run Ex1 -v
package lesson36

import (
	"io"
	"time"
)

// Ex1: LimitedReader reads from R, but returns io.EOF after N bytes in total.
// Rebuild io.LimitReader (README section 2). Try `go doc -src io.LimitedReader.Read` AFTER you've
// written your own, to compare.
type LimitedReader struct {
	R io.Reader
	N int64
}

func (l *LimitedReader) Read(p []byte) (int, error) {
	// TODO
	return l.R.Read(p)
}

// Ex2: return a Writer that writes everything to ALL of the given writers, in order.
// Write returns len(p), nil if every writer succeeded. If a writer returns an error, stop and return
// that error. If a writer writes fewer bytes than len(p) without an error, return io.ErrShortWrite.
//
// Hint: define a type multiWriter struct{ writers []io.Writer } with a Write method,
// and return &multiWriter{writers: ws}.
func MultiWriter(ws ...io.Writer) io.Writer {
	// TODO
	if len(ws) == 0 {
		return io.Discard
	}
	return ws[0]
}

// Processor transforms a string.
type Processor interface {
	Process(s string) string
}

// Ex3: ProcessorFunc lets a plain function be used as a Processor (README section 3).
// TODO: add the method  func (f ProcessorFunc) Process(s string) string  that calls f.
type ProcessorFunc func(string) string

// Pipeline runs s through every processor, in order. Already done.
func Pipeline(s string, ps ...Processor) string {
	for _, p := range ps {
		s = p.Process(s)
	}
	return s
}

// Ex4: QueryError describes a failed database query and wraps its cause (README section 5).
// Error() must return:  query "<Query>": <Err's message>
// Unwrap() must return Err, so errors.Is / errors.As can see through it.
type QueryError struct {
	Query string
	Err   error
}

// TODO: write the Error and Unwrap methods (pointer receivers).

// Ex5: a client with functional options (README section 4).
//
// Defaults: Timeout 10s, Retries 3, UserAgent "gomastery/1.0".
//
//	c := NewClient("https://api.example.com")
//	c := NewClient("https://api.example.com", WithTimeout(2*time.Second), WithRetries(0))
type Client struct {
	BaseURL   string
	Timeout   time.Duration
	Retries   int
	UserAgent string
}

// Option changes one setting of a Client.
type Option func(*Client)

func NewClient(baseURL string, opts ...Option) *Client {
	// TODO: set the defaults, then apply every option
	return &Client{BaseURL: baseURL}
}

// WithTimeout sets the timeout.
func WithTimeout(d time.Duration) Option {
	// TODO
	return func(c *Client) {}
}

// WithRetries sets the number of retries. Negative values are treated as 0.
func WithRetries(n int) Option {
	// TODO
	return func(c *Client) {}
}

// WithUserAgent sets the user agent.
func WithUserAgent(ua string) Option {
	// TODO
	return func(c *Client) {}
}

// Person is someone to sort.
type Person struct {
	First, Last string
	Age         int
}

// Ex6: sort people IN PLACE by last name, then first name, then age with the OLDEST first.
// Use slices.SortFunc with cmp.Or and cmp.Compare (README section 6). Import "cmp" and "slices".
func SortPeople(people []Person) {
	// TODO
}
