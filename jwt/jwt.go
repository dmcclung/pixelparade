package jwt

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"strings"
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
		return nil, fmt.Errorf("apple secret key not configured")
	}
	if config.AppleKeyID, ok = os.LookupEnv("APPLE_KEY_ID"); !ok {
		return nil, fmt.Errorf("apple key ID not configured")
	}
	if config.AppleAppID, ok = os.LookupEnv("APPLE_APP_ID"); !ok {
		return nil, fmt.Errorf("apple app ID not configured")
	}
	if config.AppleTeamID, ok = os.LookupEnv("APPLE_TEAM_ID"); !ok {
		return nil, fmt.Errorf("apple team ID not configured")
	}

	return &config, nil
}

func getPrivateKey(keyData string) (*ecdsa.PrivateKey, error) {
	keyData = strings.ReplaceAll(keyData, "\\n", "\n")

	block, _ := pem.Decode([]byte(keyData))
	if block == nil {
		return nil, fmt.Errorf("error decode pem block")
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
		return nil, fmt.Errorf("key is not type *ecdsa.PrivateKey")
	}

	return ecdsaKey, nil
}

type JWK struct {
	Kty string `json:"kty"`
	Kid string `json:"kid"`
	Use string `json:"use"`
	Alg string `json:"alg"`
	N   string `json:"n"`
	E   string `json:"e"`
}

type JWKSet struct {
	Keys []JWK `json:"keys"`
}

func base64URLDecode(data string) ([]byte, error) {
	// Pad string with trailing '=' to make it base64 standard compliant
	if m := len(data) % 4; m != 0 {
		data += strings.Repeat("=", 4-m)
	}
	return base64.URLEncoding.DecodeString(data)
}

func convertJWKToPEM(nStr, eStr string) (string, error) {
	// Decode the base64 URL encoded parameters
	nBytes, err := base64URLDecode(nStr)
	if err != nil {
		return "", err
	}

	eBytes, err := base64URLDecode(eStr)
	if err != nil {
		return "", err
	}

	// Convert eBytes into an integer (big endian)
	e := big.NewInt(0).SetBytes(eBytes)

	// Construct an rsa.PublicKey with the modulus and exponent
	pub := &rsa.PublicKey{
		N: new(big.Int).SetBytes(nBytes),
		E: int(e.Int64()),
	}

	// Marshal the rsa.PublicKey into ASN.1 DER-encoded form
	derPkix, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return "", err
	}

	// Encode the DER-encoded public key to PEM
	pemData := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: derPkix,
	})

	return string(pemData), nil
}

func fetchApplePublicKeys() (*JWKSet, error) {
	res, err := http.Get("https://appleid.apple.com/auth/keys")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var jwks JWKSet
	if err := json.NewDecoder(res.Body).Decode(&jwks); err != nil {
		return nil, err
	}

	return &jwks, nil
}

func keyFunc(jwks *JWKSet) jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, fmt.Errorf("kid header not found")
		}

		for _, key := range jwks.Keys {
			if key.Kid == kid {
				pem, err := convertJWKToPEM(key.N, key.E)
				if err != nil {
					return nil, err
				}
				return jwt.ParseRSAPublicKeyFromPEM([]byte(pem))
			}
		}

		return nil, fmt.Errorf("key %v not found", kid)
	}
}

func GetSubFromJWT(idToken string) (string, error) {
	jwks, err := fetchApplePublicKeys()
	if err != nil {
		return "", fmt.Errorf("fetching apple keys: %w", err)
	}

	token, err := jwt.ParseWithClaims(idToken, jwt.MapClaims{}, keyFunc(jwks))
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
