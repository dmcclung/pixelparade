package jwt

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

type config struct {
	AppleKey    string
	AppleKeyID  string
	AppleAppID  string
	AppleTeamID string
}

func loadConfig() (*config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("loading .env file: %v", err)
	}

	config := config{}

	var ok bool
	if config.AppleKey, ok = os.LookupEnv("APPLE_KEY"); !ok {
		return nil, fmt.Errorf("Apple secret key not configured")
	}
	if config.AppleKeyID, ok = os.LookupEnv("APPLE_KEY_ID"); !ok {
		return nil, fmt.Errorf("Apple key ID not configured")
	}
	if config.AppleAppID, ok = os.LookupEnv("APPLE_APP_ID"); !ok {
		return nil, fmt.Errorf("Apple app ID not configured")
	}
	if config.AppleTeamID, ok = os.LookupEnv("APPLE_TEAM_ID"); !ok {
		return nil, fmt.Errorf("Apple team ID not configured")
	}

	return &config, nil
}

func getPrivateKey(keyPath string) (*ecdsa.PrivateKey, error) {
	keyData, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("reading key file: %w", err)
	}

	block, _ := pem.Decode(keyData)
	if block == nil {
		return nil, fmt.Errorf("decode pem block: %w", err)
	}

	if block.Type != "PRIVATE KEY" {
		return nil, fmt.Errorf("PEM block type is not PRIVATE KEY, got %s", block.Type)
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse private key: %w", err)
	}

	ecdsaKey, ok := privateKey.(*ecdsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("Key is not type *ecdsa.PrivateKey")
	}

	return ecdsaKey, nil
}

func GetSubFromJWT(idToken string) (string, error) {
	token, _, err := jwt.NewParser().ParseUnverified(idToken, jwt.MapClaims{})
	if err != nil {
		return "", fmt.Errorf("parsing identity token: %w", err)
	}

	subject, err := token.Claims.GetSubject()
	if err != nil {
		return "", fmt.Errorf("subject claim: %w", err)
	}

	return subject, nil
}

func GenerateJWT() (string, error) {
	config, err := loadConfig()
	if err != nil {
		return "", fmt.Errorf("loading config: %w", err)
	}

	claims := jwt.MapClaims{
		"iss": config.AppleTeamID,
		"sub": config.AppleAppID,
		"aud": "https://appleid.apple.com",
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	token.Header["kid"] = config.AppleKeyID

	ecdsaKey, err := getPrivateKey(config.AppleKey)
	if err != nil {
		return "", fmt.Errorf("private key: %w", err)
	}

	signed, err := token.SignedString(ecdsaKey)
	if err != nil {
		return "", fmt.Errorf("error signing token: %w", err)
	}

	log.Printf("Signed token: %v\n", signed)
	return signed, nil
}
