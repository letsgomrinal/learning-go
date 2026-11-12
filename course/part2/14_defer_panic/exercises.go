// Lesson 14 exercises: defer, panic and recover. Read README.md first!
//
// Check one exercise at a time:   go test -run Ex1 -v
package lesson14

import "fmt"

// Log records lines, so the tests can check what happened and in which order.
type Log struct {
	Lines []string
}

// Add appends a line to the log.
func (l *Log) Add(line string) {
	l.Lines = append(l.Lines, line)
}

// Ex1: log "open", then log "work", and make sure "close" is logged LAST, using defer.
// Write the defer right after "open", the way you'd write defer f.Close() right after opening a file.
// The final log must be: open, work, close
func RunTask(l *Log) {
	// TODO
}

// Ex2: for each name, log "open <name>" and defer logging "close <name>".
// CloseAll(l, []string{"a", "b"}) must produce: open a, open b, close b, close a
//
// Hint: defer l.Add("close " + name) inside the loop. The argument is evaluated immediately.
func CloseAll(l *Log, names []string) {
	// TODO
}

// Ex3: log "start <name>", then call work(), and make sure "end <name>" is ALWAYS logged,
// even if work() panics. (Don't recover here. Let the panic continue. Just make sure "end" is logged.)
func Measure(l *Log, name string, work func()) {
	// TODO
}

// Ex4: return a / b. Integer division by zero panics; recover from it and return
// 0 and an error whose message starts with "recovered: ".
//
// Follow README section 6, "Turning a panic into an error". Don't check b == 0 yourself:
// the point is to practise recover.
func SafeDivide(a, b int) (result int, err error) {
	// TODO
	return a / b, nil
}

// Ex5: call f. If f panics, recover and return an error "panic: <value>"
// (that's fmt.Errorf("panic: %v", r)). If f finishes normally, return nil.
//
//	SafeRun(func() { panic("boom") })  → error "panic: boom"
//	SafeRun(func() {})                 → nil
func SafeRun(f func()) (err error) {
	// TODO
	f()
	return nil
}

// Ex6: MustPositive is meant for values written in the program's own code (like a config constant),
// where a bad value means the programmer made a mistake.
// Return n if it's greater than 0. Otherwise panic with the message "MustPositive: got <n>".
//
// Hint: panic(fmt.Sprintf(...))
func MustPositive(n int) int {
	// TODO
	return n
}

// Ex7: run every job. A job may panic; that must NOT stop the other jobs.
// Return how many jobs succeeded and how many panicked.
//
// Hint: use your SafeRun from Ex5 on each job.
func RunJobs(jobs []func()) (succeeded, failed int) {
	// TODO
	return 0, 0
}

var _ = fmt.Sprintf // keeps the "fmt" import used until you use it yourself
