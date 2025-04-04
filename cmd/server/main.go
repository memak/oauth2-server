package main

import (
	"net/http"

	_ "github.com/memak/oauth2-server/config"
	"github.com/memak/oauth2-server/internal/handlers"
	"github.com/memak/oauth2-server/internal/middleware"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	r := mux.NewRouter()
	port := viper.GetString("server.port")

	// Create a rate limiter (e.g., 5 requests per second with a burst of 10)
	rateLimiter := middleware.NewRateLimiter(5, 10)

	// Apply the rate limiting middleware
	r.Use(rateLimiter.RateLimitMiddleware)

	r.HandleFunc("/token", handlers.TokenHandler).Methods("POST")
	r.HandleFunc("/jwks", handlers.JWKSHandler).Methods("GET")
	r.HandleFunc("/introspect", handlers.IntrospectHandler).Methods("POST")

	log.Infof("Started OAuth2 server on :%s...", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
