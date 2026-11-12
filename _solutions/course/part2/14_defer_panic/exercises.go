// Reference solution for lesson 14.
package lesson14

import "fmt"

type Log struct {
	Lines []string
}

func (l *Log) Add(line string) {
	l.Lines = append(l.Lines, line)
}

func RunTask(l *Log) {
	l.Add("open")
	defer l.Add("close")
	l.Add("work")
}

func CloseAll(l *Log, names []string) {
	for _, name := range names {
		l.Add("open " + name)
		defer l.Add("close " + name)
	}
}

func Measure(l *Log, name string, work func()) {
	l.Add("start " + name)
	defer l.Add("end " + name)
	work()
}

func SafeDivide(a, b int) (result int, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("recovered: %v", r)
		}
	}()
	return a / b, nil
}

func SafeRun(f func()) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v", r)
		}
	}()
	f()
	return nil
}

func MustPositive(n int) int {
	if n <= 0 {
		panic(fmt.Sprintf("MustPositive: got %d", n))
	}
	return n
}

func RunJobs(jobs []func()) (succeeded, failed int) {
	for _, job := range jobs {
		if err := SafeRun(job); err != nil {
			failed++
		} else {
			succeeded++
		}
	}
	return succeeded, failed
}
