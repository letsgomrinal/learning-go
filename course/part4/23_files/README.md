# Lesson 23: Files and I/O

**Goal:** read and write files, process big files line by line, walk through folders, and handle
"file not found" properly.

**Time:** about 1.5 hours.

```bash
cd ~/Desktop/learning-go/course/part4/23_files
```

Part 4 is about building real programs. Everything from Parts 1–3 comes together here.

---

## 1. Reading a whole file

Python:

```python
with open("notes.txt") as f:
    text = f.read()
```

Go:

```go
data, err := os.ReadFile("notes.txt")
if err != nil {
	return err
}
text := string(data)
```

`os.ReadFile` returns the content as `[]byte` (raw bytes, lesson 11). Convert with `string(data)` to get
text. It opens and closes the file for you.

It reads the **whole file into memory**, which is fine for config files and small data. For big files,
read line by line (section 4).

---

## 2. Writing a whole file

```go
err := os.WriteFile("notes.txt", []byte("hello\n"), 0o644)
```

- It creates the file, or **replaces** it if it exists.
- `[]byte("hello\n")` converts the string to bytes.
- `0o644` is the file's **permissions**, in Unix octal notation: the owner can read and write,
  everyone else can only read. Use `0o644` for normal files. (Python uses the same numbers in
  `os.chmod`.)

---

## 3. Paths: `filepath.Join`

Don't build paths by gluing strings with `"/"`. Use `path/filepath`, which uses the right separator
for the operating system:

```go
p := filepath.Join("data", "2026", "notes.txt")   // data/2026/notes.txt on Mac/Linux
filepath.Base(p)     // "notes.txt"
filepath.Dir(p)      // "data/2026"
filepath.Ext(p)      // ".txt"
```

Python: `os.path.join` / `pathlib.Path`.

Create folders with `os.MkdirAll("data/2026", 0o755)` (like `mkdir -p`; `0o755` is the usual
permission for folders).

---

## 4. Reading line by line

For large files (logs, CSVs), open the file and scan it one line at a time:

```go
func countErrors(path string) (int, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer f.Close()                      // lesson 14: close it when we return

	count := 0
	sc := bufio.NewScanner(f)            // f is an io.Reader (lesson 11)
	for sc.Scan() {
		if strings.Contains(sc.Text(), "ERROR") {
			count++
		}
	}
	return count, sc.Err()               // report any read error that stopped the loop
}
```

Python: `for line in f:`. Only one line is in memory at a time, so this works on a 10 GB file.

Notes:
- **`defer f.Close()` right after a successful `os.Open`.** Every open file must be closed.
- **Check `sc.Err()` after the loop.** `Scan` returns false at the end of the file **and** on errors;
  `Err()` tells you which (it's `nil` at a normal end).
- A `Scanner` handles lines up to 64 KB by default. For longer lines, see `sc.Buffer`.

---

## 5. Writing piece by piece

```go
f, err := os.Create("report.txt")        // create or truncate
if err != nil {
	return err
}
defer f.Close()

w := bufio.NewWriter(f)                  // batches small writes into big ones (faster)
for _, line := range lines {
	fmt.Fprintln(w, line)                // lesson 11: Fprintln writes to any io.Writer
}
return w.Flush()                         // ⚠️ write out whatever is still in the buffer
```

Forgetting `Flush()` is a classic bug: the last few KB never reach the file.

### Appending to a file

`os.Create` wipes the file. To add to the end instead, use `os.OpenFile` with flags:

```go
f, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
```

The flags are combined with `|`: append, create if missing, open for writing only. Python: `open(p, "a")`.

---

## 6. "File not found" and other errors

File errors are normal errors (lesson 12). To check **why** something failed, use `errors.Is` with the
sentinel errors in package `io/fs`:

```go
data, err := os.ReadFile("config.yaml")
if errors.Is(err, fs.ErrNotExist) {
	fmt.Println("no config file, using defaults")
} else if err != nil {
	return err
}
```

| Sentinel | Meaning | Python equivalent |
|---|---|---|
| `fs.ErrNotExist` | the file or folder doesn't exist | `FileNotFoundError` |
| `fs.ErrExist` | it already exists | `FileExistsError` |
| `fs.ErrPermission` | not allowed | `PermissionError` |

Because `errors.Is` looks through wrapping (lesson 13), this still works after you add context with
`fmt.Errorf("load config: %w", err)`.

To check if a file exists without opening it:

```go
_, err := os.Stat(path)
exists := err == nil
```

`os.Stat` also gives you the size, modification time and whether it's a folder.

---

## 7. Walking a folder tree

```go
err := filepath.WalkDir("project", func(path string, d fs.DirEntry, err error) error {
	if err != nil {
		return err                         // couldn't read something: stop
	}
	if d.IsDir() {
		return nil                         // keep going into folders
	}
	if filepath.Ext(path) == ".go" {
		fmt.Println(path)
	}
	return nil
})
```

`WalkDir` calls your function once for every file and folder, recursively. Python: `os.walk`.
Return `filepath.SkipDir` from the function to skip a whole folder (like `.git`).

To list just one folder (not recursive): `entries, err := os.ReadDir(dir)`.

---

## 8. Copying with `io.Copy`

`io.Copy(dst, src)` copies everything from any `io.Reader` to any `io.Writer`, in chunks, without
loading it all into memory:

```go
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}
```

The same `io.Copy` copies an HTTP response body to a file, or a file to an HTTP response. That's
`io.Reader` and `io.Writer` paying off again.

### 🤔 Predict the output

```go
dir := os.TempDir()
p := filepath.Join(dir, "demo.txt")
os.WriteFile(p, []byte("one\n"), 0o644)
os.WriteFile(p, []byte("two\n"), 0o644)
data, _ := os.ReadFile(p)
fmt.Printf("%q\n", string(data))
```

<details markdown="1">
<summary>Answer</summary>

`"two\n"`. `WriteFile` replaces the whole file each time. To keep `one`, you'd append with
`os.OpenFile(..., os.O_APPEND|...)`.
</details>

---

## 9. Testing code that uses files

Never let tests write into your real folders. `t.TempDir()` gives each test a fresh, empty folder that
is **deleted automatically** when the test ends:

```go
func TestSave(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.txt")
	// ... call your code with path, then check the file ...
}
```

All of this lesson's tests use it.

---

## 10. Common mistakes

| You wrote | What happens | Fix |
|---|---|---|
| Opened a file and never closed it | "too many open files" in long-running programs | `defer f.Close()` after a successful open |
| `defer f.Close()` **before** checking `err` | calling Close on a nil file | check `err` first, then `defer` |
| `bufio.Writer` without `Flush()` | the end of the file is missing | `return w.Flush()` |
| `"data/" + name` | breaks on Windows, double slashes | `filepath.Join("data", name)` |
| `err.Error() == "no such file or directory"` | fragile | `errors.Is(err, fs.ErrNotExist)` |
| `os.ReadFile` on a huge file | uses all your memory | `bufio.Scanner` |
| Not checking `sc.Err()` | read errors look like a normal end of file | `return n, sc.Err()` |

---

## 11. Exercises

| | Exercise | Practises |
|---|---|---|
| Ex1 | `WriteNote` | `os.WriteFile`, `filepath.Join` |
| Ex2 | `ReadNote` with a wrapped error | `os.ReadFile`, `%w` |
| Ex3 | `CountLines` of any size of file | `os.Open`, `defer Close`, `bufio.Scanner` |
| Ex4 | `AppendLine` | `os.OpenFile` with flags |
| Ex5 | `Exists` | `os.Stat`, `fs.ErrNotExist` |
| Ex6 | `FindFiles` by extension, recursively | `filepath.WalkDir` |
| Ex7 | `CopyFile` | `io.Copy` |

---

## ✅ Checklist

- [ ] I can read and write whole files with `os.ReadFile` / `os.WriteFile`
- [ ] I can read a file line by line with `bufio.Scanner`, closing it with `defer`
- [ ] I build paths with `filepath.Join`
- [ ] I check for missing files with `errors.Is(err, fs.ErrNotExist)`
- [ ] I can walk a folder tree, and I use `t.TempDir()` in tests

**Next → [Lesson 24: Command-line tools](../24_cli/README.md)**
