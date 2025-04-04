package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/memak/oauth2-server/internal/auth"
	log "github.com/sirupsen/logrus"
)

type introspectResponse struct {
	Active bool `json:"active"`
}

func IntrospectHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	token := r.FormValue("token")
	if token == "" {
		log.Warn("Missing token in request")
		http.Error(w, "Missing token", http.StatusBadRequest)
		return
	}

	parsedToken, err := auth.ValidateJWT(token)
	if err != nil {
		log.Errorf("Failed to validate token: %v", err)
	}

	resp := introspectResponse{Active: err == nil && parsedToken.Valid}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Errorf("Failed to encode response: %v", err)
	}
}
