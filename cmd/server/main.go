package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", getRoot)
	http.ListenAndServe(":8080", r)
}

func getRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello, %q", r.URL.Path)
}
