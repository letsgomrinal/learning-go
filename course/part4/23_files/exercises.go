// Lesson 23 exercises: files and I/O. Read README.md first!
//
// Check one exercise at a time:   go test -run Ex1 -v
//
// The tests give your functions a temporary folder (t.TempDir()), so nothing touches your real files.
package lesson23

// Ex1: write text to a file called name inside dir. Replace it if it already exists.
// Use filepath.Join(dir, name) and os.WriteFile with permissions 0o644.
func WriteNote(dir, name, text string) error {
	// TODO
	return nil
}

// Ex2: read the file called name inside dir, and return its content as a string.
// If reading fails, return "" and the error wrapped as "read note <name>: <err>".
// errors.Is(err, fs.ErrNotExist) must still work for a missing file. (Which verb?)
func ReadNote(dir, name string) (string, error) {
	// TODO
	return "", nil
}

// Ex3: return the number of lines in the file at path.
// Open it with os.Open, defer Close, and count with a bufio.Scanner. Return sc.Err() at the end.
// If the file can't be opened, return 0 and the error.
func CountLines(path string) (int, error) {
	// TODO
	return 0, nil
}

// Ex4: add line plus "\n" to the end of the file at path, creating it if it doesn't exist.
// Hint: os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644), defer Close,
// then write with f.WriteString or fmt.Fprintln.
func AppendLine(path, line string) error {
	// TODO
	return nil
}

// Ex5: return true if something exists at path (a file or a folder), false otherwise.
// Hint: _, err := os.Stat(path)
func Exists(path string) bool {
	// TODO
	return false
}

// Ex6: return the paths of all files under root (recursively) whose extension is ext (like ".go"),
// RELATIVE to root, sorted alphabetically. Skip any folder named ".git" entirely.
//
//	root/
//	  main.go
//	  util/helpers.go
//	  util/notes.txt
//	  .git/hooks.go
//
//	FindFiles(root, ".go") == ["main.go", "util/helpers.go"]
//
// Hints:
//   - filepath.WalkDir (README section 7)
//   - inside the callback: if d.IsDir() && d.Name() == ".git" { return filepath.SkipDir }
//   - rel, err := filepath.Rel(root, path) gives the path relative to root
//   - sort the result with slices.Sort
func FindFiles(root, ext string) ([]string, error) {
	// TODO
	return nil, nil
}

// Ex7: copy the file at src to dst (create or replace dst), using io.Copy. (README section 8.)
// Close both files. Return any error.
func CopyFile(src, dst string) error {
	// TODO
	return nil
}
