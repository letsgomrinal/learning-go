# Lesson 30: Shipping your program

**Goal:** get a Go program ready for the real world: configuration from environment variables, files
embedded in the binary, a version number, graceful shutdown, building for other systems, and a small
Docker image.

**Time:** about 1.5 hours.

```bash
cd ~/Desktop/learning-go/course/part4/30_shipping
```

---

## 1. One binary, no runtime

Shipping a Python service means shipping Python, a virtualenv, and every dependency, usually inside a
large Docker image. A Go program compiles into **one file** with everything inside it:

```bash
go build -o app ./cmd/app
ls -lh app          # about 8 MB, and it's all you need
./app
```

Copy that file to a server with the same OS and CPU type, and it runs. That's a big reason teams pick
Go for services and tools.

---

## 2. Configuration from environment variables

The common convention for services ("12-factor apps"): settings come from **environment variables**,
so the same binary runs on your laptop, in staging and in production with different settings.

```go
port := os.Getenv("PORT")      // "" if not set
if port == "" {
	port = "8080"
}
```

`os.LookupEnv("PORT")` also tells you whether the variable is set at all (`value, ok`).

**Make it testable:** instead of calling `os.Getenv` deep inside your code, write a function that takes
a "get" function, and pass `os.Getenv` in `main`:

```go
func LoadConfig(getenv func(string) string) (Config, error) { ... }

cfg, err := LoadConfig(os.Getenv)                   // in main
cfg, err := LoadConfig(func(k string) string {...}) // in a test: any values you like
```

That's dependency injection again (lesson 29), with a function instead of an interface.

**Fail fast:** check the whole config at startup and exit with a clear message if something is wrong.
That's much better than discovering a missing setting at 3 a.m. when the first request needs it.

---

## 3. Embedding files with `//go:embed`

HTML templates, default config files, SQL migrations, images: you can put them **inside the binary**,
so you still ship just one file:

```go
import _ "embed"

//go:embed static/index.html
var indexHTML string
```

- The `//go:embed` comment (no space after `//`) must be **directly above** a package-level variable.
- The path is relative to the folder of the `.go` file.
- The variable can be a `string`, a `[]byte`, or an `embed.FS` for a whole folder:

```go
import "embed"

//go:embed static
var staticFiles embed.FS

http.Handle("GET /static/", http.FileServerFS(staticFiles))
```

The file is read at **build time**. Change it, and you need to rebuild.

---

## 4. A version number with `-ldflags`

When something breaks in production, the first question is "which version is running?". Declare a
variable:

```go
var Version = "dev"
```

and set it when building, without editing code:

```bash
go build -ldflags "-X main.Version=1.4.2" -o app ./cmd/app
./app -version      # 1.4.2
```

`-X package.Variable=value` overwrites a string variable at link time. Build scripts and CI usually pass
the git tag: `-X main.Version=$(git describe --tags)`. Expose it in a health check endpoint, and you
can always tell what's deployed.

(Go also records module and git information automatically. See `runtime/debug.ReadBuildInfo()` and
`go version -m ./app`.)

---

## 5. Graceful shutdown

When you deploy a new version, the old process gets a signal to stop (`SIGTERM`, or Ctrl+C locally).
If it just dies, requests in progress are cut off and users get errors. A **graceful shutdown**:
1. stops accepting new connections,
2. waits for requests in progress to finish (with a time limit),
3. then exits.

`http.Server.Shutdown` does exactly that. Combined with `signal.NotifyContext` (lesson 22):

```go
func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	srv := &http.Server{Addr: ":8080", Handler: mux}

	go func() {
		<-ctx.Done()                                        // wait for Ctrl+C / SIGTERM
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		srv.Shutdown(shutdownCtx)                           // finish in-flight requests
	}()

	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
	log.Println("bye")
}
```

Notes:
- After `Shutdown` is called, `ListenAndServe` returns **`http.ErrServerClosed`**. That's the normal,
  expected result, not a failure, so check for it with `errors.Is`.
- The shutdown context needs a **fresh** timeout: the original `ctx` is already cancelled at that point.
- Also set timeouts on the server itself in production: `ReadHeaderTimeout`, `ReadTimeout`,
  `WriteTimeout`, `IdleTimeout`. The zero values mean "wait forever".

---

## 6. Building for other systems

Go can build for another operating system or CPU from your Mac, by setting two environment variables:

```bash
GOOS=linux   GOARCH=amd64 go build -o app-linux-amd64 ./cmd/app     # most cloud servers
GOOS=linux   GOARCH=arm64 go build -o app-linux-arm64 ./cmd/app     # ARM servers, Raspberry Pi
GOOS=windows GOARCH=amd64 go build -o app.exe ./cmd/app
```

No extra tools needed. `go tool dist list` shows every supported combination.

Add `CGO_ENABLED=0` to get a fully static binary that doesn't depend on any system libraries. That's
what you want for tiny Docker images. Smaller binaries: `-ldflags "-s -w"` removes debug information.

---

## 7. A small Docker image

The `Dockerfile` in this folder uses a **multi-stage build**: one stage with the full Go toolchain to
compile, and a final stage containing **only the binary**:

```dockerfile
FROM golang:1.27 AS build
WORKDIR /src
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags "-s -w -X main.Version=1.0.0" -o /app ./course/part4/30_shipping/cmd/app

FROM gcr.io/distroless/static-debian12
COPY --from=build /app /app
USER nonroot
ENTRYPOINT ["/app"]
```

The final image is around **10–15 MB**. A typical Python service image is several hundred MB. It
also contains no shell and no package manager, so there's much less to attack.

(You don't need Docker installed for this lesson. If you have it:
`docker build -f course/part4/30_shipping/Dockerfile -t notes-app .` from the `learning-go` folder.)

### 🤔 Predict the output

```go
var Version = "dev"

func main() {
	fmt.Println("version", Version)
}
```

```bash
go run .
go run -ldflags "-X main.Version=2.0.0" .
```

<details><summary>Answer</summary>

`version dev`, then `version 2.0.0`. The default value applies unless the linker overwrites it with `-X`.
</details>

---

## 8. A release checklist

Before you deploy a Go service:

- [ ] `go vet ./...` and `go test -race ./...` pass
- [ ] `staticcheck ./...` or `golangci-lint run` (popular linters) are clean
- [ ] `govulncheck ./...` reports no known vulnerabilities in your dependencies
- [ ] Config comes from the environment, and is checked at startup
- [ ] The HTTP server has timeouts and shuts down gracefully
- [ ] A health endpoint reports the version
- [ ] Logs are structured (`slog`), errors carry context (`%w`)
- [ ] The binary is built with `CGO_ENABLED=0` into a small image, running as a non-root user

---

## 9. Common mistakes

| You wrote | What happens | Fix |
|---|---|---|
| `// go:embed file` (with a space) | it's just a comment; the variable stays empty | `//go:embed file` |
| `//go:embed` above a variable inside a function | compile error | package-level variables only |
| Treating `http.ErrServerClosed` as a failure | an error message on every clean shutdown | `!errors.Is(err, http.ErrServerClosed)` |
| Reusing the cancelled ctx for `Shutdown` | shutdown gives up immediately | `context.WithTimeout(context.Background(), ...)` |
| `os.Getenv` scattered all over the code | hard to test, hard to see all settings | one `LoadConfig`, checked at startup |
| `-X main.version=...` with a different name or package | silently does nothing | the full, exact `package.Variable` name |

---

## 10. Exercises

| | Exercise | Practises |
|---|---|---|
| Ex1 | `LoadConfig` from environment variables | defaults, validation, testable config |
| Ex2 | Embed `static/index.html` and serve it | `//go:embed` |
| Ex3 | `Health` handler with the version | `Version` variable, `-ldflags` |
| Ex4 | `Serve` with graceful shutdown | `http.Server`, `Shutdown`, `ErrServerClosed` |

Then run the whole app (it uses all four):

```bash
DATABASE_PATH=/tmp/x.db go run -ldflags "-X main.Version=0.1.0" ./cmd/app
# in another terminal:
curl localhost:8080/          # your embedded page
curl localhost:8080/health    # {"status":"ok","version":"0.1.0"}
# back in the first terminal: Ctrl+C → graceful shutdown
```

---

## ✅ Checklist: end of Part 4 🎉

- [ ] I load config from environment variables in one testable function, and fail fast
- [ ] I can embed files into the binary
- [ ] I can set a version at build time with `-ldflags -X`
- [ ] I can shut an HTTP server down gracefully
- [ ] I can cross-compile, and I understand how a multi-stage Dockerfile gets a tiny image

---

## What's next

You've finished Part 4: you can build real programs in Go.

1. **Practise.** Do `practice/09_http_json`.
2. **Build the portfolio project.** The Part 4 milestone in [roadmap.html](../../../roadmap.html): a URL
   shortener API that combines lessons 23–30: HTTP + JSON + SQLite + middleware + graceful shutdown +
   tests + a Docker image. Then put it on GitHub.
3. **Read.** *"Let's Go"* and *"Let's Go Further"* by Alex Edwards walk through building complete web
   applications in exactly this style.

Then Part 5: expert topics (how the scheduler, memory and garbage collector work; profiling and
performance).
