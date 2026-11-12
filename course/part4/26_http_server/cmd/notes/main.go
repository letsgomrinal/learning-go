// Command notes runs the lesson 26 notes API on :8080.
package main

import (
	"log"
	"net/http"

	lesson26 "gomastery/course/part4/26_http_server"
)

func main() {
	log.Println("listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", lesson26.NewServer()))
}
