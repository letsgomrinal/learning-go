package lesson24

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

func TestEx1ParseArgs(t *testing.T) {
	tests := []struct {
		args []string
		want Options
	}{
		{[]string{"-l", "a.txt", "b.txt"}, Options{Lines: true, Files: []string{"a.txt", "b.txt"}}},
		{[]string{}, Options{Lines: true, Words: true, Bytes: true}},
		{[]string{"x.txt"}, Options{Lines: true, Words: true, Bytes: true, Files: []string{"x.txt"}}},
		{[]string{"-w", "-c"}, Options{Words: true, Bytes: true}},
	}
	for _, tt := range tests {
		got, err := ParseArgs(tt.args)
		if err != nil {
			t.Errorf("Ex1: ParseArgs(%q) returned error %v", tt.args, err)
			continue
		}
		if got.Lines != tt.want.Lines || got.Words != tt.want.Words || got.Bytes != tt.want.Bytes ||
			!slices.Equal(got.Files, tt.want.Files) && !(len(got.Files) == 0 && len(tt.want.Files) == 0) {
			t.Errorf("Ex1: ParseArgs(%q) = %+v, want %+v", tt.args, got, tt.want)
		}
	}
	if _, err := ParseArgs([]string{"-z"}); err == nil {
		t.Error("Ex1: ParseArgs([-z]) should return an error for an unknown flag")
	}
}

func TestEx2Count(t *testing.T) {
	tests := []struct {
		in   string
		want Counts
	}{
		{"hello world\nbye\n", Counts{2, 3, 16}},
		{"", Counts{0, 0, 0}},
		{"no newline", Counts{0, 2, 10}},
		{"  spaced   out \n\n", Counts{2, 2, 17}},
		{"héllo\n", Counts{1, 1, 7}},
	}
	for _, tt := range tests {
		got, err := Count(strings.NewReader(tt.in))
		if err != nil || got != tt.want {
			t.Errorf("Ex2: Count(%q) = %+v, %v; want %+v, nil", tt.in, got, err, tt.want)
		}
	}
}

func TestEx3Format(t *testing.T) {
	all := Options{Lines: true, Words: true, Bytes: true}
	tests := []struct {
		c    Counts
		opts Options
		name string
		want string
	}{
		{Counts{2, 3, 16}, all, "a.txt", "       2       3      16 a.txt"},
		{Counts{2, 3, 16}, Options{Words: true}, "", "       3"},
		{Counts{10, 20, 300}, Options{Lines: true, Bytes: true}, "total", "      10     300 total"},
	}
	for _, tt := range tests {
		if got := Format(tt.c, tt.opts, tt.name); got != tt.want {
			t.Errorf("Ex3: Format(%+v, %+v, %q)\n   got:  %q\n   want: %q", tt.c, tt.opts, tt.name, got, tt.want)
		}
	}
}

func runGwc(t *testing.T, args []string, stdin string) (code int, stdout, stderr string) {
	t.Helper()
	var out, errOut strings.Builder
	code = Run(args, strings.NewReader(stdin), &out, &errOut)
	return code, out.String(), errOut.String()
}

func TestEx4Run(t *testing.T) {
	dir := t.TempDir()
	a := filepath.Join(dir, "a.txt")
	b := filepath.Join(dir, "b.txt")
	os.WriteFile(a, []byte("hello world\nbye\n"), 0o644)
	os.WriteFile(b, []byte("one\ntwo\nthree\n"), 0o644)

	t.Run("stdin", func(t *testing.T) {
		code, out, errOut := runGwc(t, []string{"-w"}, "a b c\nd\n")
		if code != 0 || out != "       4\n" || errOut != "" {
			t.Errorf("Ex4: gwc -w < stdin: code %d, stdout %q, stderr %q; want 0, \"       4\\n\", \"\"", code, out, errOut)
		}
	})

	t.Run("one file", func(t *testing.T) {
		code, out, _ := runGwc(t, []string{a}, "")
		want := "       2       3      16 " + a + "\n"
		if code != 0 || out != want {
			t.Errorf("Ex4: gwc a.txt: code %d, stdout\n   got:  %q\n   want: %q", code, out, want)
		}
	})

	t.Run("two files and a total", func(t *testing.T) {
		code, out, _ := runGwc(t, []string{"-l", a, b}, "")
		want := "       2 " + a + "\n       3 " + b + "\n       5 total\n"
		if code != 0 || out != want {
			t.Errorf("Ex4: gwc -l a.txt b.txt: code %d, stdout\n   got:  %q\n   want: %q", code, out, want)
		}
	})

	t.Run("missing file", func(t *testing.T) {
		missing := filepath.Join(dir, "nope.txt")
		code, out, errOut := runGwc(t, []string{"-l", missing, a}, "")
		if code != 1 {
			t.Errorf("Ex4: exit code with a missing file = %d, want 1", code)
		}
		if !strings.HasPrefix(errOut, "gwc: ") || !strings.Contains(errOut, "nope.txt") {
			t.Errorf("Ex4: stderr = %q, want \"gwc: <error mentioning nope.txt>\"", errOut)
		}
		want := "       2 " + a + "\n       2 total\n"
		if out != want {
			t.Errorf("Ex4: stdout should still show the file that worked\n   got:  %q\n   want: %q", out, want)
		}
	})

	t.Run("bad flag", func(t *testing.T) {
		code, out, errOut := runGwc(t, []string{"-z"}, "")
		if code != 2 || out != "" || !strings.HasPrefix(errOut, "gwc: ") {
			t.Errorf("Ex4: gwc -z: code %d, stdout %q, stderr %q; want 2, \"\", \"gwc: ...\"", code, out, errOut)
		}
	})
}
