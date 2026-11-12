package lesson31

import (
	"os"
	"strings"
)

// usesPool checks that exercises.go declares and uses a sync.Pool for Ex5 (ignoring comments).
func usesPool() bool {
	src, err := os.ReadFile("exercises.go")
	if err != nil {
		return false
	}
	var code strings.Builder
	for _, line := range strings.Split(string(src), "\n") {
		if !strings.HasPrefix(strings.TrimSpace(line), "//") {
			code.WriteString(line + "\n")
		}
	}
	c := code.String()
	return strings.Contains(c, "sync.Pool") && strings.Contains(c, ".Get()") && strings.Contains(c, ".Put(")
}
