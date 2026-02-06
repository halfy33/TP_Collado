package main

import "net/http"

// Routeur pour l'ensemble de nos routes

func router() http.Handler {
	mux := http.NewServeMux()
	// health URL
	mux.HandleFunc("GET /health", HealthHandle)
	mux.HandleFunc("GET /avg", ReadLoadHandler)
	mux.HandleFunc("GET /disk", ReadDiskhandler)
	mux.HandleFunc("GET /cpu", ReadCpuHandler)
	mux.HandleFunc("GET /mem", ReadMemoryHandler)
	mux.HandleFunc("GET /net", ReadNetcardHandler)

	// Autres cas : fichiers statiques
	fs := http.FileServer(http.Dir("www"))
	mux.Handle("/", fs)

	return mux
}
