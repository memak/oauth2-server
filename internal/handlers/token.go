package handlers

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/memak/oauth2-server/internal/auth"
	"github.com/memak/oauth2-server/internal/storage"
)

type tokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

func TokenHandler(w http.ResponseWriter, r *http.Request) {
	header := r.Header.Get("Authorization")
	if !strings.HasPrefix(header, "Basic ") {
		http.Error(w, "Invalid auth header", http.StatusUnauthorized)
		return
	}
	authHeader := strings.TrimPrefix(header, "Basic ")
	decoded, _ := base64.StdEncoding.DecodeString(authHeader)
	parts := strings.SplitN(string(decoded), ":", 2)
	if len(parts) != 2 || !storage.ValidateClient(parts[0], parts[1]) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	// Assuming auth is a package with GenerateJWT function
	token, err := auth.GenerateJWT(parts[0])
	if err != nil {
		http.Error(w, "Token generation failed", http.StatusInternalServerError)
		return
	}

	resp := tokenResponse{AccessToken: token, TokenType: "bearer", ExpiresIn: 3600}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
