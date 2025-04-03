package auth

import (
	"crypto/rsa"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var privateKey *rsa.PrivateKey
var publicKey *rsa.PublicKey

func init() {
	privKeyData, _ := os.ReadFile("keys/private.pem")
	pubKeyData, _ := os.ReadFile("keys/public.pem")
	privateKey, _ = jwt.ParseRSAPrivateKeyFromPEM(privKeyData)
	publicKey, _ = jwt.ParseRSAPublicKeyFromPEM(pubKeyData)
}

func PublicKey() *rsa.PublicKey {
	return publicKey
}

func GenerateJWT(clientID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iss": "simple-oauth2-server",
		"sub": clientID,
		"aud": "api",
		"exp": time.Now().Add(time.Hour).Unix(),
	})
	return token.SignedString(privateKey)
}

func ValidateJWT(tokenStr string) (*jwt.Token, error) {
	return jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})
}
