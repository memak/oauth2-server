package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/memak/oauth2-server/internal/auth"
	log "github.com/sirupsen/logrus"
)

func IntrospectHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
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

		if scope, ok := claims["scope"]; ok {
			resp["scope"] = scope
		}
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.WithError(err).Error("Failed to encode response")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}