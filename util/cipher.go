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
		/* Do Base 64 Encoding */
		base64Text := base64.StdEncoding.EncodeToString(textBytes)

		/* Build Cipher Text Placeholder */
		ciphertext := make([]byte, aes.BlockSize+len(base64Text))

		/* Encryt using AES */
		iv := ciphertext[:aes.BlockSize]
		if _, err := io.ReadFull(rand.Reader, iv); err == nil {
			cfb := cipher.NewCFBEncrypter(block, iv)
			cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(base64Text))
			encryptedText = string(ciphertext)
		}
	}
	return
}

func Decrypt(key, text string) (decryptedText string, err error) {
	keyBytes := []byte(key)
	textBytes := []byte(text)

	/* Create New Cipher */
	if block, err := aes.NewCipher(keyBytes); err == nil {
		/* Check Minimum Cipher TExt Length */
		if len(textBytes) > aes.BlockSize {
			/* Extract Block of Text to be Encrypted */
			textRight := textBytes[:aes.BlockSize]
			textBytes = textBytes[aes.BlockSize:]

			/* Decrypt Aes */
			cfb := cipher.NewCFBDecrypter(block, textRight)
			cfb.XORKeyStream(textBytes, textBytes)

			/* Decode Base64 */
			if decodedBytes, err := base64.StdEncoding.DecodeString(string(textBytes)); err == nil {
				decryptedText = string(decodedBytes)
			}
		} else {
			err = CIPHER_TOO_SHORT
		}
	}
	return
}
