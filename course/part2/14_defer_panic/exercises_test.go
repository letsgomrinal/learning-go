package lesson14

import (
	"slices"
	"strings"
	"testing"
)

func TestEx1RunTask(t *testing.T) {
	var l Log
	RunTask(&l)
	if want := []string{"open", "work", "close"}; !slices.Equal(l.Lines, want) {
		t.Errorf("Ex1: log = %q, want %q", l.Lines, want)
	}
}

func TestEx2CloseAll(t *testing.T) {
	var l Log
	CloseAll(&l, []string{"a", "b", "c"})
	want := []string{"open a", "open b", "open c", "close c", "close b", "close a"}
	if !slices.Equal(l.Lines, want) {
		t.Errorf("Ex2: log = %q, want %q", l.Lines, want)
	}
}

func TestEx3Measure(t *testing.T) {
	var l Log
	Measure(&l, "fast", func() { l.Add("working") })
	if want := []string{"start fast", "working", "end fast"}; !slices.Equal(l.Lines, want) {
		t.Errorf("Ex3: log = %q, want %q", l.Lines, want)
	}

	var l2 Log
	func() {
		defer func() { recover() }() // the test catches the panic that Measure lets through
		Measure(&l2, "broken", func() { panic("boom") })
	}()
	if want := []string{"start broken", "end broken"}; !slices.Equal(l2.Lines, want) {
		t.Errorf("Ex3: when work panics, log = %q, want %q (use defer for \"end\")", l2.Lines, want)
	}
}

func TestEx4SafeDivide(t *testing.T) {
	if r, err := SafeDivide(10, 3); r != 3 || err != nil {
		t.Errorf("Ex4: SafeDivide(10, 3) = %d, %v; want 3, nil", r, err)
	}
	var r int
	var err error
	func() {
		defer func() {
			if p := recover(); p != nil {
				t.Fatalf("Ex4: SafeDivide(1, 0) panicked (%v). Recover inside SafeDivide", p)
			}
		}()
		r, err = SafeDivide(1, 0)
	}()
	if r != 0 || err == nil || !strings.HasPrefix(err.Error(), "recovered: ") {
		t.Errorf("Ex4: SafeDivide(1, 0) = %d, %v; want 0 and an error starting with \"recovered: \"", r, err)
	}
}

func TestEx5SafeRun(t *testing.T) {
	if err := SafeRun(func() {}); err != nil {
		t.Errorf("Ex5: SafeRun(no panic) = %v, want nil", err)
	}
	var err error
	func() {
		defer func() {
			if p := recover(); p != nil {
				t.Fatalf("Ex5: the panic escaped SafeRun (%v)", p)
			}
		}()
		err = SafeRun(func() { panic("boom") })
	}()
	if err == nil || err.Error() != "panic: boom" {
		t.Errorf("Ex5: SafeRun(panic(\"boom\")) = %v, want \"panic: boom\"", err)
	}
}

func TestEx6MustPositive(t *testing.T) {
	if got := MustPositive(5); got != 5 {
		t.Errorf("Ex6: MustPositive(5) = %d, want 5", got)
	}
	for _, n := range []int{0, -3} {
		func() {
			defer func() {
				r := recover()
				if r == nil {
					t.Errorf("Ex6: MustPositive(%d) didn't panic", n)
					return
				}
				want := "MustPositive: got " + itoa(n)
				if r != want {
					t.Errorf("Ex6: MustPositive(%d) panicked with %v, want %q", n, r, want)
				}
			}()
			MustPositive(n)
		}()
	}
}

func TestEx7RunJobs(t *testing.T) {
	ran := 0
	jobs := []func(){
		func() { ran++ },
		func() { panic("bad job") },
		func() { ran++ },
		func() { var m map[string]int; m["x"] = 1 }, // a runtime panic
		func() { ran++ },
	}
	ok, failed := RunJobs(jobs)
	if ok != 3 || failed != 2 || ran != 3 {
		t.Errorf("Ex7: RunJobs = %d succeeded, %d failed (and %d jobs ran); want 3, 2 (3 ran)", ok, failed, ran)
	}
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	if n < 0 {
		return "-" + itoa(-n)
	}
	s := ""
	for n > 0 {
		s = string(rune('0'+n%10)) + s
		n /= 10
	}
	return s
}
