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
	Scope   	string `json:"scope"`
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

	// Check grant_type
	if err := r.ParseForm(); err != nil {
		log.Warn("Failed to parse form data")
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	grantType := r.FormValue("grant_type")
	if grantType != "client_credentials" {
		log.Warnf("Unsupported grant_type: %s", grantType)
		http.Error(w, "Unsupported grant_type", http.StatusBadRequest)
		return
	}

	parts := strings.SplitN(string(decoded), ":", 2)
	if len(parts) != 2 || !storage.ValidateClient(parts[0], parts[1]) {
		log.Warn("Invalid credentials")
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	requestedScopes := strings.Fields(r.FormValue("scope"))
	if len(requestedScopes) == 0 {
		log.Info("No scopes provided, attempting to use default scopes")
		client, ok := storage.GetClient(parts[0])
		
		if !ok {
			log.Warn("Failed to retrieve client")
			http.Error(w, "Failed to retrieve client", http.StatusBadRequest)
			return
		}

		requestedScopes = client.DefaultScopes		
		if len(requestedScopes) == 0 {
			log.Warn("No default scopes found")
			http.Error(w, "No default scopes found", http.StatusBadRequest)
			return
		}
	}

	if !storage.ValidateScopes(parts[0], requestedScopes) {
		log.Warnf("Invalid scopes requested: %v", requestedScopes)
		http.Error(w, "Invalid scope", http.StatusBadRequest)
		return
	}

	token, err := auth.GenerateJWT(parts[0], requestedScopes)
	if err != nil {
		log.Errorf("Token generation failed: %v", err)
		http.Error(w, "Token generation failed", http.StatusInternalServerError)
		return
	}

	resp := tokenResponse{AccessToken: token, TokenType: "bearer", ExpiresIn: expiresIn, Scope: strings.Join(requestedScopes, " ")}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Errorf("Failed to encode response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
