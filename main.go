package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"crypto/rand"

	"crypto/ecdsa"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/lestrrat-go/jwx/v2/jwk"
)

const PORT = 3000

type AuthInitResponse struct {
	URL string `json:"url"`
}

type AuthVerifyPayload struct {
	Jwt        string `json:"jwt"`
	DeviceCode string `json:"device_code"`
}

type User struct {
	Email     string `json:"email"`
	SessionID string `json:"session_id"`
}

type AuthEntry struct {
	Verified bool `json:"verified"`
	User     User `json:"user"`
}

func generateDeviceSecret() (string, error) {
	// From doc: "Note that no error handling is necessary, as Read always succeeds."
	key := make([]byte, 32)
	rand.Read(key)

	return base64.URLEncoding.EncodeToString(key), nil
}

var authMap map[string]AuthEntry

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	authMap = make(map[string]AuthEntry)
	webAppURL := os.Getenv("WEB_APP_URL")
	jwksJSON := []byte(os.Getenv("JWKS_JSON"))

	set, err := jwk.Parse(jwksJSON)
	if err != nil {
		log.Fatal("Failed to parse JWKS:", err)
	}

	key, ok := set.Key(0)
	if !ok {
		log.Fatal("No key found in JWKS")
	}

	var pubKey ecdsa.PublicKey
	if err := key.Raw(&pubKey); err != nil {
		log.Fatal("Failed to extract public key:", err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	r.Post("/auth/init", func(w http.ResponseWriter, _ *http.Request) {
		deviceCode, _ := generateDeviceSecret()
		url := fmt.Sprintf("%s/auth/login/%s", webAppURL, deviceCode)
		authMap[deviceCode] = AuthEntry{Verified: false}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(201)
		json.NewEncoder(w).Encode(AuthInitResponse{URL: url})
	})

	r.Post("/auth/verify", func(w http.ResponseWriter, r *http.Request) {
		var payload AuthVerifyPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		token, err := jwt.Parse(payload.Jwt, func(_ *jwt.Token) (any, error) {
			return &pubKey, nil
		}, jwt.WithValidMethods([]string{"ES256"}))

		if err != nil {
			fmt.Println(err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			fmt.Println(claims)

			email, ok := claims["email"].(string)
			if !ok {
				http.Error(w, "Invalid payload", http.StatusUnauthorized)
				return
			}
			sessionID, ok := claims["session_id"].(string)

			if !ok {
				http.Error(w, "Invalid payload", http.StatusUnauthorized)
				return
			}
			user := User{Email: email, SessionID: sessionID}

			auth, ok := authMap[payload.DeviceCode]

			if !ok {
				http.Error(w, "we don't expect this", http.StatusUnauthorized)
				return
			}

			auth.Verified = true
			auth.User = user

			authMap[payload.DeviceCode] = auth

			// return 200
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(user)
		} else {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
	})

	r.Get("/auth/verify", func(w http.ResponseWriter, r *http.Request) {
		deviceCode := r.URL.Query().Get("device_code")
		user, ok := authMap[deviceCode]

		if !ok {
			http.Error(w, "Unauthorized!", http.StatusUnauthorized)
			return
		}

		if !user.Verified {
			http.Error(w, "Not yet logged in", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	})

	fmt.Printf("server started at port %d", PORT)
	http.ListenAndServe(fmt.Sprintf(":%d", PORT), r)
}
