package main

import "net/http"

// Routeur pour l'ensemble de nos routes

func router() http.Handler {
	mux := http.NewServeMux()
	// health URL
	mux.HandleFunc("GET /health", HealthHandle)

	// Autres cas : fichiers statiques
	fs := http.FileServer(http.Dir("www"))
	mux.Handle("/", fs)

	return mux
}
