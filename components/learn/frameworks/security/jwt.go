package security

import (
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func GenerateToken(privateKey *rsa.PrivateKey) (string, error) {
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iss":     "aman",
		"purpose": "jwtfun",
		"iat":     time.Now().Unix(),
		"exp":     time.Now().Add(time.Hour).Unix(),
		"sub":     "aman",
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(privateKey)
	return tokenString, err
}

func MarshalPublicKey(publicKey crypto.PublicKey) []byte {
	if pubkeyBytes, err := x509.MarshalPKIXPublicKey(publicKey); err == nil {
		pubkeyPem := pem.EncodeToMemory(
			&pem.Block{
				Type:  "RSA PUBLIC KEY",
				Bytes: pubkeyBytes,
			},
		)
		return pubkeyPem
	}
	return nil
}

func MarshalPrivateKey(privateKey *rsa.PrivateKey) []byte {
	privatekeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privatekeyPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: privatekeyBytes,
		},
	)
	return privatekeyPem
}

func UnmarshalRSAPublicKey(pemPubKey string) (key *rsa.PublicKey, err error) {
	var block *pem.Block
	var pub any
	if block, _ = pem.Decode([]byte(pemPubKey)); block != nil {
		if pub, err = x509.ParsePKIXPublicKey(block.Bytes); err == nil {
			switch key := pub.(type) {
			case *rsa.PublicKey:
				return key, nil
			default:
				err = errors.New("unknown type of public key")
			}
		}
	}
	return
}

func UnmarshalRSAPrivateKey(pemPrivateKey string) (key *rsa.PrivateKey, err error) {
	var block *pem.Block
	var private any
	if block, _ = pem.Decode([]byte(pemPrivateKey)); block != nil {
		if private, err = x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
			switch key := private.(type) {
			case *rsa.PrivateKey:
				return key, nil
			default:
				err = errors.New("unknown type of private key")
			}
		}
	}
	return
}
