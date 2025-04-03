package main

import (
	"log"
	"net/http"

	"github.com/memak/oauth2-server/internal/handlers"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/token", handlers.TokenHandler).Methods("POST")
	r.HandleFunc("/jwks", handlers.JWKSHandler).Methods("GET")
	r.HandleFunc("/introspect", handlers.IntrospectHandler).Methods("POST")

	log.Println("Started OAuth2 server on :8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
