package utils

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"go-fiber-api/config"
	"go-fiber-api/models"
)

func Sign(claims jwt.MapClaims, cfg *config.Config) (string, error) {
	privateKey, err := parsePrivateKey(cfg.JWTPrivateKey)
	if err != nil {
		return "", fmt.Errorf("failed to parse private key: %w", err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodPS256, claims)

	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}

func Verify(text string, cfg *config.Config) (jwt.MapClaims, error) {
	publicKey, err := parsePublicKey(cfg.JWTPublicKey)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse public key: %w", err)
	}

	token, err := jwt.Parse(text, func(token *jwt.Token) (interface{}, error) {
		switch token.Method.(type) {
		case *jwt.SigningMethodRSA:
			// RS256, RS384, RS512
			return publicKey, nil
		case *jwt.SigningMethodRSAPSS:
			// PS256, PS384, PS512 ← FIXED: Now accepts PS256!
			return publicKey, nil
		default:
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
	})

	if err != nil {
		return nil, fmt.Errorf("Failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("Invalid token")
}

func CreateToken(userID uint, opId *uint, userName, email, telephone string, level models.UserLevel, roleAction models.UserRoleAction, cfg *config.Config) (string, error) {
	claims := jwt.MapClaims{
		"sub":         "Json Web Token subject",
		"aud":         "Json Web Token audience",
		"iss":         "Json Web Token issuer",
		"exp":         time.Now().Add(24 * time.Hour).Unix(),
		"iat":         time.Now().Unix(),
		"user_id":     userID,
		"username":    userName,
		"email":       email,
		"telephone":   telephone,
		"level":       level,
		"role_action": roleAction,
		"op_id":       opId,
	}

	return Sign(claims, cfg)
}

func parsePrivateKey(key string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(key))
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

func parsePublicKey(key string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(key))
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	pub, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return pub, nil
}
