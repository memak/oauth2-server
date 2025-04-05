package storage

import "slices"

type Client struct {
	Secret string
	Scopes []string
	DefaultScopes []string
}

var clients = map[string]Client{
	"client_id": {
		Secret: "secret",
		Scopes: []string{"read:products", "write:orders"},
		DefaultScopes: []string{"read:products"},
	},
	"readonly_client": {
		Secret: "readonly",
		Scopes: []string{"read:products"},
		DefaultScopes: []string{"read:products"},
	},
}

func ValidateClient(clientID, secret string) bool {
	client, ok := clients[clientID]
	return ok && client.Secret == secret
}

func ValidateScopes(clientID string, requestedScopes []string) bool {
	client, ok := clients[clientID]
	if !ok {
		return false
	}
	for _, requested := range requestedScopes {
		if !slices.Contains(client.Scopes, requested) {
			return false
		}
	}
	return true
}

func GetClient(clientID string) (Client, bool) {
	client, ok := clients[clientID]
	return client, ok
}