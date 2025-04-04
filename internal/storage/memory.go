package storage

import log "github.com/sirupsen/logrus"

var clients = map[string]string{
	"client_id": "secret",
}

func ValidateClient(clientID, secret string) bool {
	storedSecret, ok := clients[clientID]
	if !ok {
		log.Warnf("Client ID %s not found", clientID)
		return false
	}
	if storedSecret != secret {
		log.Warnf("Invalid secret for client ID %s", clientID)
		return false
	}
	log.Infof("Client ID %s validated successfully", clientID)
	return true
}
