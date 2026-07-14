package utils

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"log"
	"strings"

	"go-fiber-api/config"

	jose "github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"
)

type Claims struct {
	jwt.Claims
	Role string `json:"role"`
}

func loadPrivateKey(cfg *config.Config) (*rsa.PrivateKey, error) {
	raw := cfg.JWTPrivateKey
	if raw == "" {
		// b64 := os.Getenv("JWT_PRIVATE_KEY_B64")
		// if b64 == "" {
		// 	return nil, errors.New("no private key found in env")
		// }
		// decoded, err := base64.StdEncoding.DecodeString(b64)
		// if err != nil {
		// 	return nil, fmt.Errorf("decode base64 key: %w", err)
		// }
		// raw = string(decoded)
		return nil, errors.New("no private key found in env")
	} else {
		raw = strings.ReplaceAll(raw, `\n`, "\n")
	}

	block, _ := pem.Decode([]byte(raw))
	if block == nil {
		return nil, errors.New("failed to decode PEM block from private key")
	}

	if key, err := x509.ParsePKCS8PrivateKey(block.Bytes); err == nil {
		rsaKey, ok := key.(*rsa.PrivateKey)
		if !ok {
			return nil, errors.New("key is not an RSA private key")
		}
		return rsaKey, nil
	}
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

func loadPublicKey(cfg *config.Config) (*rsa.PublicKey, error) {
	raw := cfg.JWTPublicKey
	if raw == "" {
		// b64 := os.Getenv("JWT_PUBLIC_KEY_B64")
		// if b64 == "" {
		// 	return nil, errors.New("no public key found in env")
		// }
		// decoded, err := base64.StdEncoding.DecodeString(b64)
		// if err != nil {
		// 	return nil, fmt.Errorf("decode base64 key: %w", err)
		// }
		// raw = string(decoded)
		return nil, errors.New("no public key found in env")
	} else {
		raw = strings.ReplaceAll(raw, `\n`, "\n")
	}

	block, _ := pem.Decode([]byte(raw))
	if block == nil {
		return nil, errors.New("failed to decode PEM block from public key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("key is not an RSA public key")
	}
	return rsaPub, nil
}

func SignToken(claim Claims, cfg *config.Config) (string, error) {
	privKey, err := loadPrivateKey(cfg)
	if err != nil {
		log.Fatal("load private key:", err)
		return "", err
	}
	signer, err := jose.NewSigner(jose.SigningKey{
		Algorithm: jose.PS256,
		Key:       privKey,
	}, (&jose.SignerOptions{}).WithType("JWT"))
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	claims := Claims{
		Claims: jwt.Claims{
			Subject: "Json Web Token subject",
			Issuer:  "Json Web Token issuer",
		},
		Role: "admin",
	}

	token, err := jwt.Signed(signer).Claims(claims).Serialize()
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	return token, nil
}

func VerifyToken(text string, cfg *config.Config) (any, error) {
	pubKey, err := loadPublicKey(cfg)
	if err != nil {
		log.Fatal("load public key:", err)
		return nil, err
	}
	parsedToken, err := jwt.ParseSigned(text, []jose.SignatureAlgorithm{jose.PS256})
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	var out Claims
	if err := parsedToken.Claims(pubKey, &out); err != nil {
		log.Fatal("verification failed:", err)
		return nil, err
	}

	return out, nil
}
