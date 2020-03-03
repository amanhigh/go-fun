package security

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

func JwtFun() {
	//Generate RSA Key Pair
	privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	publicKey := privateKey.Public()

	//Marshal to PEM Format
	if pubkey_bytes, err := x509.MarshalPKIXPublicKey(publicKey); err == nil {
		pubkey_pem := pem.EncodeToMemory(
			&pem.Block{
				Type:  "RSA PUBLIC KEY",
				Bytes: pubkey_bytes,
			},
		)
		fmt.Println("Public Pem", string(pubkey_pem))
	}

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
	fmt.Println("JwtToken", tokenString, err)

	token, err = jwt.Parse(tokenString, func(token *jwt.Token) (i interface{}, err error) {
		return publicKey, nil
	})

	claims := token.Claims.(jwt.MapClaims)
	fmt.Println("ParsedToken", claims["purpose"], claims, claims.Valid(), err)
}
