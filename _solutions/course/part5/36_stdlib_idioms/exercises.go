// Reference solution for lesson 36.
package lesson36

import (
	"cmp"
	"io"
	"slices"
	"time"
)

type LimitedReader struct {
	R io.Reader
	N int64
}

func (l *LimitedReader) Read(p []byte) (int, error) {
	if l.N <= 0 {
		return 0, io.EOF
	}
	if int64(len(p)) > l.N {
		p = p[:l.N]
	}
	n, err := l.R.Read(p)
	l.N -= int64(n)
	return n, err
}

type multiWriter struct {
	writers []io.Writer
}

func (m *multiWriter) Write(p []byte) (int, error) {
	for _, w := range m.writers {
		n, err := w.Write(p)
		if err != nil {
			return n, err
		}
		if n != len(p) {
			return n, io.ErrShortWrite
		}
	}
	return len(p), nil
}

func MultiWriter(ws ...io.Writer) io.Writer {
	return &multiWriter{writers: ws}
}

type Processor interface {
	Process(s string) string
}

type ProcessorFunc func(string) string

func (f ProcessorFunc) Process(s string) string {
	return f(s)
}

func Pipeline(s string, ps ...Processor) string {
	for _, p := range ps {
		s = p.Process(s)
	}
	return s
}

type QueryError struct {
	Query string
	Err   error
}

func (e *QueryError) Error() string {
	return "query \"" + e.Query + "\": " + e.Err.Error()
}

func (e *QueryError) Unwrap() error {
	return e.Err
}

type Client struct {
	BaseURL   string
	Timeout   time.Duration
	Retries   int
	UserAgent string
}

type Option func(*Client)

func NewClient(baseURL string, opts ...Option) *Client {
	c := &Client{
		BaseURL:   baseURL,
		Timeout:   10 * time.Second,
		Retries:   3,
		UserAgent: "gomastery/1.0",
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func WithTimeout(d time.Duration) Option {
	return func(c *Client) { c.Timeout = d }
}

func WithRetries(n int) Option {
	return func(c *Client) { c.Retries = max(n, 0) }
}

func WithUserAgent(ua string) Option {
	return func(c *Client) { c.UserAgent = ua }
}

type Person struct {
	First, Last string
	Age         int
}

func SortPeople(people []Person) {
	slices.SortFunc(people, func(a, b Person) int {
		return cmp.Or(
			cmp.Compare(a.Last, b.Last),
			cmp.Compare(a.First, b.First),
			cmp.Compare(b.Age, a.Age), // oldest first
		)
	})
}
