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
var iss string;
var aud string;
var alg string;

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
	iss = "https://oauth2.com"
	aud = "https://api.oauth2.com"
	alg = "RS256"
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
		"iss":   iss,
		"sub":   clientID,
		"aud":   aud,
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

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Errorf("Invalid token claims")
		return nil, jwt.ErrTokenInvalidClaims
	}

	if token.Method.Alg() != alg {
		log.Errorf("Invalid algorithm: expected %s, got %s", alg, token.Method.Alg())
		return nil, jwt.ErrTokenUnverifiable
	}

	if claims["iss"] != iss {
		log.Errorf("Invalid issuer: expected %s, got %s", iss, claims["iss"])
		return nil, jwt.ErrTokenInvalidIssuer
	}

	if claims["aud"] != aud {
		log.Errorf("Invalid audience: expected %s, got %s", aud, claims["aud"])
		return nil, jwt.ErrTokenInvalidAudience
	}

	return token, nil
}
