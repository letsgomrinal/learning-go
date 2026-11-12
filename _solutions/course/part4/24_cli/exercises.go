// Reference solution for lesson 24.
package lesson24

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

type Options struct {
	Lines, Words, Bytes bool
	Files               []string
}

func ParseArgs(args []string) (Options, error) {
	fs := flag.NewFlagSet("gwc", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	lines := fs.Bool("l", false, "count lines")
	words := fs.Bool("w", false, "count words")
	bytes := fs.Bool("c", false, "count bytes")
	if err := fs.Parse(args); err != nil {
		return Options{}, err
	}
	opts := Options{Lines: *lines, Words: *words, Bytes: *bytes, Files: fs.Args()}
	if !opts.Lines && !opts.Words && !opts.Bytes {
		opts.Lines, opts.Words, opts.Bytes = true, true, true
	}
	return opts, nil
}

type Counts struct {
	Lines, Words, Bytes int
}

func Count(r io.Reader) (Counts, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return Counts{}, err
	}
	s := string(data)
	return Counts{
		Lines: strings.Count(s, "\n"),
		Words: len(strings.Fields(s)),
		Bytes: len(data),
	}, nil
}

func Format(c Counts, opts Options, name string) string {
	var b strings.Builder
	if opts.Lines {
		fmt.Fprintf(&b, "%8d", c.Lines)
	}
	if opts.Words {
		fmt.Fprintf(&b, "%8d", c.Words)
	}
	if opts.Bytes {
		fmt.Fprintf(&b, "%8d", c.Bytes)
	}
	if name != "" {
		b.WriteString(" " + name)
	}
	return b.String()
}

func countFile(name string) (Counts, error) {
	f, err := os.Open(name)
	if err != nil {
		return Counts{}, err
	}
	defer f.Close()
	return Count(f)
}

func Run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	opts, err := ParseArgs(args)
	if err != nil {
		fmt.Fprintln(stderr, "gwc:", err)
		return 2
	}

	if len(opts.Files) == 0 {
		c, err := Count(stdin)
		if err != nil {
			fmt.Fprintln(stderr, "gwc:", err)
			return 1
		}
		fmt.Fprintln(stdout, Format(c, opts, ""))
		return 0
	}

	code := 0
	var total Counts
	for _, name := range opts.Files {
		c, err := countFile(name)
		if err != nil {
			fmt.Fprintln(stderr, "gwc:", err)
			code = 1
			continue
		}
		total.Lines += c.Lines
		total.Words += c.Words
		total.Bytes += c.Bytes
		fmt.Fprintln(stdout, Format(c, opts, name))
	}
	if len(opts.Files) > 1 {
		fmt.Fprintln(stdout, Format(total, opts, "total"))
	}
	return code
}
