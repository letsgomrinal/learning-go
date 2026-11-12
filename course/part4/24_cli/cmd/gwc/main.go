// Command gwc counts lines, words and bytes, like Unix wc. (Lesson 24.)
//
//	go run ./cmd/gwc [-l] [-w] [-c] [files...]
package main

import (
	"os"

	lesson24 "gomastery/course/part4/24_cli"
)

func main() {
	os.Exit(lesson24.Run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}
