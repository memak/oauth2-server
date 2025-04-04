package main

import (
	"log"
	"net/http"

	_ "github.com/memak/oauth2-server/config"
	"github.com/memak/oauth2-server/internal/handlers"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

func main() {
	r := mux.NewRouter()
	port := viper.GetString("server.port")

	r.HandleFunc("/token", handlers.TokenHandler).Methods("POST")
	r.HandleFunc("/jwks", handlers.JWKSHandler).Methods("GET")
	r.HandleFunc("/introspect", handlers.IntrospectHandler).Methods("POST")	

	log.Printf("Started OAuth2 server on :%s...\n", port)
	log.Fatal(http.ListenAndServe(":" + port, r))
}
