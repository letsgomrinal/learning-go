// Command cmd runs the todo API from lesson 09 on :8080.
package main

import (
	"log"
	"net/http"

	todoapi "gomastery/practice/09_http_json"
)

func main() {
	log.Println("listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", todoapi.NewServer()))
}
