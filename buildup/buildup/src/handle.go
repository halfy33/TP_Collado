package main

import (
	"fmt"
	"net/http"
)

// health URL
func HealthHandle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ok")
}
