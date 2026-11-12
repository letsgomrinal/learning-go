# Multi-stage build for lesson 30's app (README section 7).
# Run from the learning-go folder:
#   docker build -f course/part4/30_shipping/Dockerfile -t notes-app .
#   docker run -p 8080:8080 -e DATABASE_PATH=/tmp/app.db notes-app

# Stage 1: compile with the full Go toolchain.
FROM golang:1.27 AS build
WORKDIR /src
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags "-s -w -X main.Version=1.0.0" -o /app ./course/part4/30_shipping/cmd/app

# Stage 2: only the binary, in a tiny image with no shell.
FROM gcr.io/distroless/static-debian12
COPY --from=build /app /app
USER nonroot
EXPOSE 8080
ENTRYPOINT ["/app"]
