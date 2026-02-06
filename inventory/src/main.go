package main

import (
	"log"
	"net/http"
)

func main() {

	log.Println("listening on :8085")
	log.Fatal(http.ListenAndServe(":8085", router()))
}
