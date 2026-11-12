package lesson23

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

func write(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func read(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("could not read %s: %v", filepath.Base(path), err)
	}
	return string(data)
}

func TestEx1WriteNote(t *testing.T) {
	dir := t.TempDir()
	if err := WriteNote(dir, "todo.txt", "buy milk\n"); err != nil {
		t.Fatalf("Ex1: WriteNote returned %v", err)
	}
	if got := read(t, filepath.Join(dir, "todo.txt")); got != "buy milk\n" {
		t.Errorf("Ex1: file contains %q, want %q", got, "buy milk\n")
	}
	WriteNote(dir, "todo.txt", "replaced")
	if got := read(t, filepath.Join(dir, "todo.txt")); got != "replaced" {
		t.Errorf("Ex1: writing again should replace the file; it contains %q", got)
	}
}

func TestEx2ReadNote(t *testing.T) {
	dir := t.TempDir()
	write(t, filepath.Join(dir, "a.txt"), "hello")
	if got, err := ReadNote(dir, "a.txt"); got != "hello" || err != nil {
		t.Errorf("Ex2: ReadNote(a.txt) = %q, %v; want \"hello\", nil", got, err)
	}
	got, err := ReadNote(dir, "missing.txt")
	if err == nil {
		t.Fatal("Ex2: ReadNote(missing.txt) returned no error")
	}
	if got != "" || !strings.HasPrefix(err.Error(), "read note missing.txt: ") {
		t.Errorf("Ex2: got %q, %q; want \"\" and an error starting with \"read note missing.txt: \"", got, err)
	}
	if !errors.Is(err, fs.ErrNotExist) {
		t.Errorf("Ex2: errors.Is(err, fs.ErrNotExist) = false. Wrap with %%w")
	}
}

func TestEx3CountLines(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "log.txt")
	write(t, p, strings.Repeat("a line\n", 5000))
	if n, err := CountLines(p); n != 5000 || err != nil {
		t.Errorf("Ex3: CountLines(5000 lines) = %d, %v; want 5000, nil", n, err)
	}
	write(t, p, "")
	if n, err := CountLines(p); n != 0 || err != nil {
		t.Errorf("Ex3: CountLines(empty file) = %d, %v; want 0, nil", n, err)
	}
	if _, err := CountLines(filepath.Join(dir, "nope.txt")); !errors.Is(err, fs.ErrNotExist) {
		t.Errorf("Ex3: CountLines(missing) error = %v, want a not-exist error", err)
	}
}

func TestEx4AppendLine(t *testing.T) {
	p := filepath.Join(t.TempDir(), "app.log")
	for _, line := range []string{"started", "working", "stopped"} {
		if err := AppendLine(p, line); err != nil {
			t.Fatalf("Ex4: AppendLine(%q) returned %v", line, err)
		}
	}
	if got, want := read(t, p), "started\nworking\nstopped\n"; got != want {
		t.Errorf("Ex4: file contains %q, want %q", got, want)
	}
}

func TestEx5Exists(t *testing.T) {
	dir := t.TempDir()
	write(t, filepath.Join(dir, "here.txt"), "x")
	if !Exists(filepath.Join(dir, "here.txt")) {
		t.Error("Ex5: Exists(existing file) = false")
	}
	if !Exists(dir) {
		t.Error("Ex5: Exists(existing folder) = false")
	}
	if Exists(filepath.Join(dir, "not-here.txt")) {
		t.Error("Ex5: Exists(missing file) = true")
	}
}

func TestEx6FindFiles(t *testing.T) {
	root := t.TempDir()
	for _, f := range []string{"main.go", "util/helpers.go", "util/notes.txt", "util/deep/more.go", ".git/hooks.go", "README.md"} {
		write(t, filepath.Join(root, f), "x")
	}
	got, err := FindFiles(root, ".go")
	if err != nil {
		t.Fatalf("Ex6: FindFiles returned %v", err)
	}
	want := []string{"main.go", "util/deep/more.go", "util/helpers.go"}
	if !slices.Equal(got, want) {
		t.Errorf("Ex6: FindFiles(root, \".go\") = %q, want %q", got, want)
	}
	if got, _ := FindFiles(root, ".md"); !slices.Equal(got, []string{"README.md"}) {
		t.Errorf("Ex6: FindFiles(root, \".md\") = %q, want [README.md]", got)
	}
}

func TestEx7CopyFile(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "src.bin")
	content := strings.Repeat("0123456789", 100000) // 1 MB
	write(t, src, content)
	dst := filepath.Join(dir, "dst.bin")
	if err := CopyFile(src, dst); err != nil {
		t.Fatalf("Ex7: CopyFile returned %v", err)
	}
	if got := read(t, dst); got != content {
		t.Errorf("Ex7: the copy has %d bytes, want %d with the same content", len(got), len(content))
	}
	if err := CopyFile(filepath.Join(dir, "nope"), dst); !errors.Is(err, fs.ErrNotExist) {
		t.Errorf("Ex7: copying a missing file: error = %v, want a not-exist error", err)
	}
}
