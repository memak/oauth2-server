package auth

import (
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	_ "github.com/memak/oauth2-server/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var privateKey *rsa.PrivateKey
var publicKey *rsa.PublicKey
var keyID string

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
	keyID = computeKeyID(publicKey)
}

func PublicKey() *rsa.PublicKey {
	return publicKey
}

func GetKeyID() string {
	return keyID
}

func computeKeyID(pub *rsa.PublicKey) string {
	pubBytes := x509.MarshalPKCS1PublicKey(pub)
	hash := sha256.Sum256(pubBytes)
	return base64.RawURLEncoding.EncodeToString(hash[:])
}

func GenerateJWT(clientID string, scopes []string) (string, error) {
	ttl := viper.GetInt("jwt.token_ttl")

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iss":   "oauth2-server",
		"sub":   clientID,
		"aud":   "api",
		"exp":   time.Now().Add(time.Duration(ttl) * time.Second).Unix(),
		"scope": strings.Join(scopes, " "),
	})

	// Add the "kid" to the header
	token.Header["kid"] = GetKeyID()

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
