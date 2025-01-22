package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

var (
	CIPHER_TOO_SHORT = errors.New("ciphertext too short")
)

func Encrypt(key, text string) (encryptedText string, err error) {
	keyBytes := []byte(key)
	textBytes := []byte(text)

	/* Create New Cipher */
	if block, err := aes.NewCipher(keyBytes); err == nil {
		/* Build Cipher Text Placeholder */
		ciphertext := make([]byte, aes.BlockSize+len(textBytes))

		/* Encrypt using AES */
		iv := ciphertext[:aes.BlockSize]
		if _, err := io.ReadFull(rand.Reader, iv); err == nil {
			cfb := cipher.NewCFBEncrypter(block, iv)
			cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(textBytes))

			/* Do Base64 Encoding on Encrypted Text */
			encryptedText = base64.URLEncoding.EncodeToString(ciphertext)
		}
	}
	return
}

func Decrypt(key, text string) (string, error) {
	keyBytes := []byte(key)

	// Decode Base64 to get back encrypted text
	textBytes, err := base64.URLEncoding.DecodeString(text)
	if err != nil {
		return "", err
	}

	// Check minimum cipher text length
	if len(textBytes) < aes.BlockSize {
		return "", CIPHER_TOO_SHORT
	}

	// Create new cipher
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	// Extract block of text to be encrypted
	textRight := textBytes[:aes.BlockSize]
	textBytes = textBytes[aes.BlockSize:]

	// Decrypt AES
	cfb := cipher.NewCFBDecrypter(block, textRight)
	cfb.XORKeyStream(textBytes, textBytes)

	// Convert decrypted bytes to string
	return string(textBytes), nil
}
