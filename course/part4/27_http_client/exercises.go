// Lesson 27 exercises: HTTP clients, middleware and logging. Read README.md first!
//
// Check one exercise at a time:   go test -race -run Ex1 -v
package lesson27

import (
	"context"
	"log/slog"
	"net/http"
)

// Ex1: return a new *http.Client with a 5 second timeout.
func NewClient() *http.Client {
	// TODO
	return nil
}

// Ex2: GET url and return the response body as a string. (README section 3: the three rules.)
//   - build the request with http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
//   - send it with client.Do; return any error
//   - defer resp.Body.Close()
//   - if the status isn't 200, return "" and an error with the message "unexpected status <code>"
//   - read the body with io.ReadAll
func GetText(ctx context.Context, client *http.Client, url string) (string, error) {
	// TODO
	return "", nil
}

// User is what the user API sends and receives.
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Ex3: GET baseURL + "/users/" + id (for example http://127.0.0.1:1234/users/7) and decode the JSON
// response into a User. Same rules as Ex2: context, close the body, status must be 200.
// Hint: strconv.Itoa(id), and json.NewDecoder(resp.Body).Decode(&u)
func GetUser(ctx context.Context, client *http.Client, baseURL string, id int) (User, error) {
	// TODO
	return User{}, nil
}

// Ex4: POST u as JSON to baseURL + "/users". Set the header "Content-Type: application/json".
// The server answers 201 with the created user (with its new ID) as JSON: decode and return it.
// Any status other than 201 → return User{} and an error "unexpected status <code>".
// Hint: data, err := json.Marshal(u); body := bytes.NewReader(data)
func CreateUser(ctx context.Context, client *http.Client, baseURL string, u User) (User, error) {
	// TODO
	return User{}, nil
}

// Ex5: return a middleware that logs every request with logger, AFTER calling next, as:
//
//	logger.Info("request", "method", r.Method, "path", r.URL.Path)
//
// The structure (README sections 6 and 7):
//
//	return func(next http.Handler) http.Handler {
//		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//			...
//		})
//	}
func Logging(logger *slog.Logger) func(http.Handler) http.Handler {
	// TODO
	return func(next http.Handler) http.Handler { return next }
}

// Ex6: a middleware that recovers from a panic in next and responds with status 500 and the
// JSON body {"error":"internal error"} (set the Content-Type to application/json). (README section 8.)
func Recover(next http.Handler) http.Handler {
	// TODO
	return next
}

// Ex7: return a middleware that only lets a request through if its "X-API-Key" header equals key.
// Otherwise respond 401 with the JSON body {"error":"unauthorized"} and do NOT call next.
func RequireAPIKey(key string) func(http.Handler) http.Handler {
	// TODO
	return func(next http.Handler) http.Handler { return next }
}

// Ex8: wrap h in all the middlewares, so that the FIRST one in the list is the OUTERMOST
// (it runs first).
//
//	Chain(h, a, b, c) == a(b(c(h)))
//
// Hint: loop over mws backwards, wrapping as you go: h = mws[i](h)
func Chain(h http.Handler, mws ...func(http.Handler) http.Handler) http.Handler {
	// TODO
	return h
}
