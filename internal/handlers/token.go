package handlers

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/memak/oauth2-server/internal/auth"
	"github.com/memak/oauth2-server/internal/storage"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type tokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

func TokenHandler(w http.ResponseWriter, r *http.Request) {
	expiresIn := viper.GetInt("jwt.token_ttl")
	header := r.Header.Get("Authorization")
	if !strings.HasPrefix(header, "Basic ") {
		log.Warn("Invalid auth header")
		http.Error(w, "Invalid auth header", http.StatusUnauthorized)
		return
	}
	authHeader := strings.TrimPrefix(header, "Basic ")
	decoded, err := base64.StdEncoding.DecodeString(authHeader)
	if err != nil {
		log.Errorf("Failed to decode auth header: %v", err)
		http.Error(w, "Invalid auth header", http.StatusUnauthorized)
		return
	}
	parts := strings.SplitN(string(decoded), ":", 2)
	if len(parts) != 2 || !storage.ValidateClient(parts[0], parts[1]) {
		log.Warn("Invalid credentials")
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	// Assuming auth is a package with GenerateJWT function
	token, err := auth.GenerateJWT(parts[0])
	if err != nil {
		log.Errorf("Token generation failed: %v", err)
		http.Error(w, "Token generation failed", http.StatusInternalServerError)
		return
	}

	resp := tokenResponse{AccessToken: token, TokenType: "bearer", ExpiresIn: expiresIn}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Errorf("Failed to encode response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
