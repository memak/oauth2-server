package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/memak/oauth2-server/internal/auth"
)

type introspectResponse struct {
	Active bool `json:"active"`
}

func IntrospectHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	token := r.FormValue("token")
	if token == "" {
		http.Error(w, "Missing token", http.StatusBadRequest)
		return
	}

	parsedToken, err := auth.ValidateJWT(token)
	resp := introspectResponse{Active: err == nil && parsedToken.Valid}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
