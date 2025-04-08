package handlers

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/memak/oauth2-server/internal/auth"
	"github.com/memak/oauth2-server/internal/storage"
	log "github.com/sirupsen/logrus"
)

func IntrospectHandler(w http.ResponseWriter, r *http.Request) {
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
	
	err = r.ParseForm()
	if err != nil {
		log.WithError(err).Error("Failed to parse form")
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	token := r.FormValue("token")
	if token == "" {
		log.Warn("Missing token in request")
		http.Error(w, "Missing token", http.StatusBadRequest)
		return
	}

	log.WithField("token", token).Info("Validating token")
	parsedToken, err := auth.ValidateJWT(token)
	active := err == nil && parsedToken.Valid

	if err != nil {
		log.WithError(err).Warn("Token validation failed")
	} 

	resp := map[string]interface{}{
		"active": active,
	}

	if active {
		claims := parsedToken.Claims.(jwt.MapClaims)
		resp["sub"] = claims["sub"]
		resp["exp"] = claims["exp"]
		resp["client_id"] = claims["sub"]
		resp["token_type"] = "access_token"
		resp["scope"] = claims["scope"]
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.WithError(err).Error("Failed to encode response")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}