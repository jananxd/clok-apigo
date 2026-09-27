package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"crypto/rand"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const PORT = 3000

type AuthDeviceResponse struct {
	DeviceCode string `json:"device_code"`
}

func generateDeviceSecret() (string, error) {
	// From doc: "Note that no error handling is necessary, as Read always succeeds."
	key := make([]byte, 32)
	rand.Read(key)

	return base64.URLEncoding.EncodeToString(key), nil
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	r.Post("/auth/device", func(w http.ResponseWriter, r *http.Request) {
		device_code, _ := generateDeviceSecret()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(201)
		json.NewEncoder(w).Encode(AuthDeviceResponse{DeviceCode: device_code})
	})

	fmt.Printf("server started at port %d", PORT)
	http.ListenAndServe(fmt.Sprintf(":%d", PORT), r)
}
