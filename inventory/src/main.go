package main

import (
	"log"
	"net/http"
)

func main() {

	log.Println("listening on :80")
	log.Fatal(http.ListenAndServe(":80", router()))
}
