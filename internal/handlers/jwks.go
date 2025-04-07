package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/memak/oauth2-server/internal/auth"
	log "github.com/sirupsen/logrus"
)

func JWKSHandler(w http.ResponseWriter, r *http.Request) {
	key, err := jwk.FromRaw(auth.PublicKey())
	if err != nil {
		log.Errorf("failed to create JWK from raw key: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if err := key.Set(jwk.KeyIDKey, auth.GetKeyID()); err != nil {
		log.Errorf("failed to set key ID: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if err := key.Set(jwk.AlgorithmKey, "RS256"); err != nil {
		log.Errorf("failed to set alg: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if err := key.Set(jwk.KeyUsageKey, "sig"); err != nil {
		log.Errorf("failed to set use: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}


	set := jwk.NewSet()
	set.AddKey(key)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(set); err != nil {
		log.Errorf("failed to encode JWK set to JSON: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
