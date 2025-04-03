package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/memak/oauth2-server/internal/auth"

	"github.com/lestrrat-go/jwx/v2/jwk"
)

func JWKSHandler(w http.ResponseWriter, r *http.Request) {
	key, _ := jwk.FromRaw(auth.PublicKey())
	key.Set(jwk.KeyIDKey, "default")
	set := jwk.NewSet()
	set.AddKey(key)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(set)
}
