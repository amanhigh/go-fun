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

var (
	ErrCipherTooShort = errors.New("ciphertext too short")
)

func Encrypt(key, text string) (encryptedText string, err error) {
	keyBytes := []byte(key)
	textBytes := []byte(text)

	/* Create New Cipher */
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", fmt.Errorf("failed to create AES cipher: %w", err)
	}

	/* Build Cipher Text Placeholder */
	ciphertext := make([]byte, aes.BlockSize+len(textBytes))

	/* Generate IV */
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", fmt.Errorf("failed to generate IV: %w", err)
	}

	/* Encrypt using AES-CTR (recommended replacement for CFB) */
	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], textBytes)

	/* Do Base64 Encoding on Encrypted Text */
	encryptedText = base64.URLEncoding.EncodeToString(ciphertext)
	return encryptedText, nil
}

func Decrypt(key, text string) (string, error) {
	keyBytes := []byte(key)

	// Decode Base64 to get back encrypted text
	textBytes, err := base64.URLEncoding.DecodeString(text)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64 string: %w", err)
	}

	// Check minimum cipher text length
	if len(textBytes) < aes.BlockSize {
		return "", ErrCipherTooShort
	}

	// Create new cipher
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", fmt.Errorf("failed to create AES cipher: %w", err)
	}

	// Extract IV from the beginning
	iv := textBytes[:aes.BlockSize]
	textBytes = textBytes[aes.BlockSize:]

	// Decrypt using AES-CTR (recommended replacement for CFB)
	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(textBytes, textBytes)

	// Convert decrypted bytes to string
	return string(textBytes), nil
}
