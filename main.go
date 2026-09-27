package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	r.Post("/auth/device", func(w http.ResponseWriter, r *http.Request) {
		// we want to return a auth URL with the device code inside so that they can just open that.
	})

	http.ListenAndServe(":3000", r)
}
