package auth

import (
	"crypto/rsa"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	_ "github.com/memak/oauth2-server/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var privateKey *rsa.PrivateKey
var publicKey *rsa.PublicKey

func init() {
	privKeyData, err := os.ReadFile(viper.GetString("paths.private_key"))
	if err != nil {
		log.Fatalf("Failed to read private key: %v", err)
	}
	pubKeyData, err := os.ReadFile(viper.GetString("paths.public_key"))
	if err != nil {
		log.Fatalf("Failed to read public key: %v", err)
	}
	privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privKeyData)
	if err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}
	publicKey, err = jwt.ParseRSAPublicKeyFromPEM(pubKeyData)
	if err != nil {
		log.Fatalf("Failed to parse public key: %v", err)
	}
}

func PublicKey() *rsa.PublicKey {
	return publicKey
}

func GenerateJWT(clientID string) (string, error) {
	ttl := viper.GetInt("jwt.token_ttl")

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iss": "oauth2-server",
		"sub": clientID,
		"aud": "api",
		"exp": time.Now().Add(time.Duration(ttl) * time.Second).Unix(),
	})
	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		log.Errorf("Failed to sign token: %v", err)
		return "", err
	}
	return signedToken, nil
}

func ValidateJWT(tokenStr string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})
	if err != nil {
		log.Errorf("Failed to validate token: %v", err)
		return nil, err
	}
	return token, nil
}
