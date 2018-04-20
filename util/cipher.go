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

func Decrypt(key, text string) (decryptedText string, err error) {
	var textBytes[]byte
	keyBytes := []byte(key)

	/* Decode Base64 to get back Encrpted Text */
	if textBytes, err = base64.URLEncoding.DecodeString(text); err == nil {
		/* Create New Cipher */
		if block, err := aes.NewCipher(keyBytes); err == nil {
			/* Check Minimum Cipher Text Length */
			if len(textBytes) > aes.BlockSize {
				/* Extract Block of Text to be Encrypted */
				textRight := textBytes[:aes.BlockSize]
				textBytes = textBytes[aes.BlockSize:]

				/* Decrypt Aes */
				cfb := cipher.NewCFBDecrypter(block, textRight)
				cfb.XORKeyStream(textBytes, textBytes)

				/* Convert Decrypted Bytes to String */
				decryptedText = string(textBytes)
			} else {
				err = CIPHER_TOO_SHORT
			}
		}
	}
	return
}
