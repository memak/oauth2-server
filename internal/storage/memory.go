package storage

var clients = map[string]string{
	"client_id": "secret",
}

func ValidateClient(clientID, secret string) bool {
	storedSecret, ok := clients[clientID]
	return ok && storedSecret == secret
}
