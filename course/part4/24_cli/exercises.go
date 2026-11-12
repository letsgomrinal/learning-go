// Lesson 24 exercises: build the gwc command-line tool. Read README.md first (especially section 7)!
//
// Check one exercise at a time:   go test -run Ex1 -v
// Try the real tool:              go run ./cmd/gwc README.md exercises.go
package lesson24

import "io"

// Options holds what the user asked for on the command line.
type Options struct {
	Lines, Words, Bytes bool     // which columns to show
	Files               []string // the file names (empty means: read standard input)
}

// Ex1: parse the command-line arguments (without the program name).
//
//   - define three bool flags on a flag.FlagSet: -l (lines), -w (words), -c (bytes)
//   - use flag.NewFlagSet("gwc", flag.ContinueOnError) and fs.SetOutput(io.Discard)
//   - if parsing fails (an unknown flag), return Options{} and the error
//   - if none of -l, -w, -c was given, turn all three on
//   - the remaining arguments (fs.Args()) are the file names
//
// ParseArgs([]string{"-l", "a.txt", "b.txt"}) == Options{Lines: true, Files: [a.txt b.txt]}
// ParseArgs([]string{})                       == Options{Lines: true, Words: true, Bytes: true}
func ParseArgs(args []string) (Options, error) {
	// TODO
	return Options{}, nil
}

// Counts holds the three numbers for one input.
type Counts struct {
	Lines, Words, Bytes int
}

// Ex2: count the lines, words and bytes readable from r.
//   - Lines: the number of "\n" characters (that's how wc counts)
//   - Words: groups of non-space characters, like strings.Fields
//   - Bytes: the total number of bytes
//
// Count(strings.NewReader("hello world\nbye\n")) == Counts{Lines: 2, Words: 3, Bytes: 16}
//
// Simplest approach: data, err := io.ReadAll(r), then work on string(data):
// strings.Count(s, "\n"), len(strings.Fields(s)), len(data).
func Count(r io.Reader) (Counts, error) {
	// TODO
	return Counts{}, nil
}

// Ex3: format one output line. For each column that's turned on in opts, in the order
// lines, words, bytes, add fmt.Sprintf("%8d", n). Then, if name isn't "", add " " + name.
//
// Format(Counts{2, 3, 16}, Options{Lines: true, Words: true, Bytes: true}, "a.txt")
//
//	== "       2       3      16 a.txt"
//
// Format(Counts{2, 3, 16}, Options{Words: true}, "") == "       3"
func Format(c Counts, opts Options, name string) string {
	// TODO
	return ""
}

// Ex4: the whole program. (README sections 5 and 7.)
//
//  1. Parse args with ParseArgs. On error: print "gwc: <err>" to stderr and return 2.
//  2. No files: Count(stdin), print Format(counts, opts, "") + "\n" to stdout, return 0.
//  3. Files: for each file, os.Open it (and close it!), Count it, and print its line.
//     If a file can't be opened or read, print "gwc: <err>" to stderr, remember that something
//     failed, and continue with the next file.
//  4. If more than one file was given, print a "total" line with the sums of the files that worked.
//  5. Return 1 if any file failed, otherwise 0.
//
// Use fmt.Fprintln(stdout, ...) and fmt.Fprintln(stderr, ...). Never os.Stdout or os.Exit here!
func Run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	// TODO
	return 0
}
