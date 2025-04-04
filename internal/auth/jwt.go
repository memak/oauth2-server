package auth

import (
	"crypto/rsa"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	_ "github.com/memak/oauth2-server/config"
	"github.com/spf13/viper"
)

var privateKey *rsa.PrivateKey
var publicKey *rsa.PublicKey

func init() {
	privKeyData, _ := os.ReadFile(viper.GetString("paths.private_key"))
	pubKeyData, _ := os.ReadFile(viper.GetString("paths.public_key"))
	privateKey, _ = jwt.ParseRSAPrivateKeyFromPEM(privKeyData)
	publicKey, _ = jwt.ParseRSAPublicKeyFromPEM(pubKeyData)
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
	return token.SignedString(privateKey)
}

func ValidateJWT(tokenStr string) (*jwt.Token, error) {
	return jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})
}
