# Lesson 24: Command-line tools

**Goal:** build a real command-line tool: read arguments and flags, read from files or standard
input, write output and errors, return exit codes, and structure `main` so the whole tool is testable.

**Time:** 1.5–2 hours. You'll build a clone of the Unix `wc` command, called `gwc`.

```bash
cd ~/Desktop/learning-go/course/part4/24_cli
```

Command-line tools are one of Go's sweet spots: one small binary, starts instantly, no interpreter or
virtualenv needed on the machine that runs it. Docker, Kubernetes' `kubectl`, Terraform and the GitHub
CLI `gh` are all Go programs.

---

## 1. Arguments: `os.Args`

```go
func main() {
	fmt.Println(os.Args)
}
```

```bash
go run ./try hello "two words" 3
# [/tmp/go-build.../exe/try hello two words 3]
```

`os.Args[0]` is the program's own path; the real arguments start at `os.Args[1:]`. All of them are
strings. Python: `sys.argv`.

---

## 2. Flags: the `flag` package

For options like `-n 5` or `-verbose`, use the standard `flag` package. Python's equivalent is `argparse`.

```go
func main() {
	name := flag.String("name", "world", "who to greet")   // -name, default "world", help text
	times := flag.Int("n", 1, "how many times")
	shout := flag.Bool("shout", false, "use upper case")
	flag.Parse()                                            // read os.Args

	msg := "hello, " + *name
	if *shout {
		msg = strings.ToUpper(msg)
	}
	for range *times {
		fmt.Println(msg)
	}
	fmt.Println("other arguments:", flag.Args())            // what's left after the flags
}
```

```bash
go run ./try -name Ana -n 2 -shout extra.txt
# HELLO, ANA
# HELLO, ANA
# other arguments: [extra.txt]

go run ./try -help        # free help text, generated from your flag definitions
```

`flag.String` returns a **pointer** (lesson 08), so you read the value with `*name`. Flags come before
the other arguments. Both `-name` and `--name` work.

### Testable flags with a `FlagSet`

`flag.String` and friends use one global set of flags tied to `os.Args`, which is awkward to test.
Instead, create your own `FlagSet` and parse any slice of strings:

```go
fs := flag.NewFlagSet("greet", flag.ContinueOnError)   // return an error instead of exiting
fs.SetOutput(io.Discard)                               // don't print usage text during tests
name := fs.String("name", "world", "who to greet")
if err := fs.Parse(args); err != nil {                 // args is any []string
	return err
}
rest := fs.Args()
```

You'll use this in the exercises.

---

## 3. Standard input, output and error

Every program has three standard streams. In Go they're `*os.File` values, which means they're also
`io.Reader` / `io.Writer` (lesson 11):

| Stream | Go | For |
|---|---|---|
| standard input | `os.Stdin` (a Reader) | data piped in: <code>cat file &#124; gwc</code> |
| standard output | `os.Stdout` (a Writer) | the program's results |
| standard error | `os.Stderr` (a Writer) | error messages and diagnostics |

Write results to stdout and errors to stderr. That way `gwc *.txt > counts.txt` puts only the counts
in the file, while errors still appear on screen.

```go
fmt.Fprintln(os.Stderr, "gwc: something went wrong")
```

---

## 4. Exit codes

When a program ends, it returns a number to the shell: **0 means success**, anything else means
failure. Scripts and CI systems rely on this. Python: `sys.exit(1)`.

```go
os.Exit(1)
```

⚠️ `os.Exit` ends the program **immediately**: deferred functions do **not** run. So call it only
once, at the very end of `main`.

Conventions: `0` success, `1` general failure, `2` wrong usage (bad flags).

---

## 5. The testable `main` pattern

If all your logic lives in `main()` and uses `os.Args`, `os.Stdout` and `os.Exit` directly, you can't
test it. The standard fix is a tiny `main` that calls a `Run` function with everything passed in:

```go
// cmd/gwc/main.go
func main() {
	os.Exit(gwc.Run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}
```

```go
// the real program: no globals, everything passed in
func Run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	...
	return 0
}
```

Now a test can call `Run` with made-up arguments, a `strings.NewReader` as stdin, and
`strings.Builder`s as stdout and stderr, then check the output and the exit code. No real files or
terminals needed. This is the most important idea in this lesson.

---

## 6. Where to put a command

Go convention: the code that builds a binary lives under `cmd/<name>/main.go`, in `package main`.
The real logic lives in a normal package next to it. This lesson is laid out that way:

```
24_cli/
├── exercises.go      ← package lesson24: Options, Count, Format, Run (you write these)
└── cmd/
    └── gwc/
        └── main.go   ← package main: just calls lesson24.Run (already written)
```

Run your tool while developing:

```bash
go run ./cmd/gwc README.md exercises.go
cat README.md | go run ./cmd/gwc -l
```

Build a real binary, and copy it anywhere:

```bash
go build -o gwc ./cmd/gwc
./gwc -w README.md
```

`go install ./cmd/gwc` builds it into `~/go/bin`. Add that folder to your `PATH` and you can run
`gwc` from anywhere.

### 🤔 Predict the output

```go
fs := flag.NewFlagSet("demo", flag.ContinueOnError)
n := fs.Int("n", 10, "count")
v := fs.Bool("v", false, "verbose")
fs.Parse([]string{"-v", "a.txt", "-n", "3"})
fmt.Println(*n, *v, fs.Args())
```

<details markdown="1">
<summary>Answer</summary>

`10 true [a.txt -n 3]`. Flag parsing **stops at the first non-flag argument** (`a.txt`). Everything
after it, including `-n 3`, is treated as a normal argument. Put flags first.
</details>

---

## 7. What you're building: `gwc`

`gwc` counts lines, words and bytes, like Unix `wc`:

```bash
$ gwc README.md exercises.go
     289    1803   10950 README.md
     118     754    4561 exercises.go
     407    2557   15511 total

$ gwc -l README.md
     289 README.md

$ cat README.md | gwc -w
    1803
```

Rules:
- Flags `-l` (lines), `-w` (words), `-c` (bytes). If none is given, show all three.
- Columns are always in the order lines, words, bytes, each right-aligned in 8 characters (`%8d`).
- With no files, read standard input, and print no name.
- With several files, add a `total` line.
- If a file can't be read, print `gwc: <error>` to **stderr**, keep going with the other files, and
  exit with code 1 at the end.
- Bad flags: exit code 2.

You'll build it in four steps (the exercises), each tested separately.

---

## 8. Common mistakes

| You wrote | What happens | Fix |
|---|---|---|
| `name := flag.String(...)` then used `name` | prints an address like `0xc000010250` | `*name` |
| Forgot `flag.Parse()` | every flag has its default value | call it before reading flags |
| Flags after a file name | treated as file names | flags first |
| Errors printed to stdout | they end up mixed into the output file | `fmt.Fprintln(stderr, ...)` |
| `os.Exit` in the middle of a function | deferred cleanup doesn't run | return an exit code, `os.Exit` only in `main` |
| All logic in `main` | can't be tested | the `Run(args, stdin, stdout, stderr) int` pattern |

---

## 9. Exercises

| | Exercise | Practises |
|---|---|---|
| Ex1 | `ParseArgs`: flags and file names | `flag.NewFlagSet` |
| Ex2 | `Count` lines, words and bytes from any reader | `io.Reader`, `bufio` |
| Ex3 | `Format` one output line | `fmt.Sprintf("%8d")` |
| Ex4 | `Run`: the whole tool | stdin/stdout/stderr, exit codes, the testable-main pattern |

When Ex4 passes, try the real thing:

```bash
go run ./cmd/gwc README.md exercises.go
go run ./cmd/gwc -l nope.txt; echo "exit code: $?"
```

---

## ✅ Checklist

- [ ] I can read arguments with `os.Args`, and flags with the `flag` package
- [ ] I know results go to stdout and errors to stderr
- [ ] I know what exit codes mean, and why `os.Exit` belongs only in `main`
- [ ] I can structure a tool as a small `main` plus a testable `Run` function
- [ ] I can build and run a binary with `go build`

**Next → [Lesson 25: JSON](../25_json/README.md)**
