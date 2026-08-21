package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

// ErrCipherTooShort is returned when a decoded ciphertext is too short to
// contain a valid AES-GCM nonce and authentication tag.
var ErrCipherTooShort = errors.New("ciphertext too short")

// Encrypt encrypts plaintext with the given AES key (16, 24, or 32 bytes) using
// AES-GCM. A fresh random nonce is generated for each call, prepended to the
// sealed ciphertext, and the result is returned as base64.URLEncoding.
func Encrypt(key, text string) (encryptedText string, err error) {
	keyBytes := []byte(key)
	textBytes := []byte(text)

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", fmt.Errorf("failed to create AES cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	sealed := gcm.Seal(nil, nonce, textBytes, nil)
	ciphertext := make([]byte, 0, len(nonce)+len(sealed))
	ciphertext = append(ciphertext, nonce...)
	ciphertext = append(ciphertext, sealed...)

	encryptedText = base64.URLEncoding.EncodeToString(ciphertext)
	return encryptedText, nil
}

// Decrypt decodes a base64.URLEncoding ciphertext produced by Encrypt and
// authenticates and decrypts it with AES-GCM using the given key. It returns
// ErrCipherTooShort for payloads too short to contain a nonce and tag, and
// rejects any ciphertext that fails authentication.
func Decrypt(key, text string) (string, error) {
	keyBytes := []byte(key)

	textBytes, err := base64.URLEncoding.DecodeString(text)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64 string: %w", err)
	}

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", fmt.Errorf("failed to create AES cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	if len(textBytes) < gcm.NonceSize()+gcm.Overhead() {
		return "", ErrCipherTooShort
	}

	nonce, sealed := textBytes[:gcm.NonceSize()], textBytes[gcm.NonceSize():]

	plaintext, err := gcm.Open(nil, nonce, sealed, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt ciphertext: %w", err)
	}

	return string(plaintext), nil
}
