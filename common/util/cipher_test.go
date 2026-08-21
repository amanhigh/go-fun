package util_test

import (
	"encoding/base64"
	"strings"

	"github.com/amanhigh/go-fun/common/util"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cipher", func() {
	Context("Encrypt and Decrypt", func() {
		var (
			validKey16   = "1234567890123456"                 // 16 bytes for AES-128
			validKey24   = "123456789012345678901234"         // 24 bytes for AES-192
			validKey32   = "12345678901234567890123456789012" // 32 bytes for AES-256
			plaintext    = "Hello, World! This is a test message."
			emptyText    = ""
			specialChars = "!@#$%^&*()_+-=[]{}|;':\",./<>?"
		)

		Context("Round Trip Encryption/Decryption", func() {
			It("should encrypt and decrypt with 16-byte key", func() {
				encrypted, err := util.Encrypt(validKey16, plaintext)
				Expect(err).NotTo(HaveOccurred())
				Expect(encrypted).NotTo(BeEmpty())
				Expect(encrypted).NotTo(Equal(plaintext))

				decrypted, err := util.Decrypt(validKey16, encrypted)
				Expect(err).NotTo(HaveOccurred())
				Expect(decrypted).To(Equal(plaintext))
			})

			It("should encrypt and decrypt with 24-byte key", func() {
				encrypted, err := util.Encrypt(validKey24, plaintext)
				Expect(err).NotTo(HaveOccurred())
				Expect(encrypted).NotTo(BeEmpty())

				decrypted, err := util.Decrypt(validKey24, encrypted)
				Expect(err).NotTo(HaveOccurred())
				Expect(decrypted).To(Equal(plaintext))
			})

			It("should encrypt and decrypt with 32-byte key", func() {
				encrypted, err := util.Encrypt(validKey32, plaintext)
				Expect(err).NotTo(HaveOccurred())
				Expect(encrypted).NotTo(BeEmpty())

				decrypted, err := util.Decrypt(validKey32, encrypted)
				Expect(err).NotTo(HaveOccurred())
				Expect(decrypted).To(Equal(plaintext))
			})

			It("should handle empty text", func() {
				encrypted, err := util.Encrypt(validKey16, emptyText)
				Expect(err).NotTo(HaveOccurred())
				Expect(encrypted).NotTo(BeEmpty()) // nonce is always present

				decrypted, err := util.Decrypt(validKey16, encrypted)
				Expect(err).NotTo(HaveOccurred())
				Expect(decrypted).To(Equal(emptyText))
			})

			It("should handle special characters", func() {
				encrypted, err := util.Encrypt(validKey16, specialChars)
				Expect(err).NotTo(HaveOccurred())
				Expect(encrypted).NotTo(BeEmpty())

				decrypted, err := util.Decrypt(validKey16, encrypted)
				Expect(err).NotTo(HaveOccurred())
				Expect(decrypted).To(Equal(specialChars))
			})

			It("should handle unicode characters", func() {
				unicode := "こんにちは世界 🌍 ñáéíóú"
				encrypted, err := util.Encrypt(validKey16, unicode)
				Expect(err).NotTo(HaveOccurred())

				decrypted, err := util.Decrypt(validKey16, encrypted)
				Expect(err).NotTo(HaveOccurred())
				Expect(decrypted).To(Equal(unicode))
			})
		})

		Context("Encryption Properties", func() {
			It("should produce different ciphertext for same plaintext", func() {
				encrypted1, err1 := util.Encrypt(validKey16, plaintext)
				encrypted2, err2 := util.Encrypt(validKey16, plaintext)

				Expect(err1).NotTo(HaveOccurred())
				Expect(err2).NotTo(HaveOccurred())
				Expect(encrypted1).NotTo(Equal(encrypted2)) // Different nonce each time
			})

			It("should produce base64 encoded output", func() {
				encrypted, err := util.Encrypt(validKey16, plaintext)
				Expect(err).NotTo(HaveOccurred())

				// Should be valid base64
				_, err = base64.URLEncoding.DecodeString(encrypted)
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("Error Conditions", func() {
			Context("Encrypt Error Handling", func() {
				It("should return error with invalid key sizes", func() {
					shortKey := "short"
					encrypted, err := util.Encrypt(shortKey, plaintext)
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("failed to create AES cipher"))
					Expect(encrypted).To(BeEmpty())
				})

				It("should return error with wrong key size", func() {
					key13Bytes := "1234567890123" // 13 bytes - invalid for AES
					encrypted, err := util.Encrypt(key13Bytes, plaintext)
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("failed to create AES cipher"))
					Expect(encrypted).To(BeEmpty())
				})
			})

			Context("Decrypt Errors", func() {
				It("should return error for invalid base64", func() {
					invalidBase64 := "not-valid-base64!@#$"
					decrypted, err := util.Decrypt(validKey16, invalidBase64)
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("failed to decode base64"))
					Expect(decrypted).To(BeEmpty())
				})

				It("should return error for ciphertext too short", func() {
					// Create valid base64 but with content shorter than the minimum
					// GCM nonce + tag length.
					shortCipher := base64.URLEncoding.EncodeToString([]byte("short"))
					decrypted, err := util.Decrypt(validKey16, shortCipher)
					Expect(err).To(HaveOccurred())
					Expect(err).To(Equal(util.ErrCipherTooShort))
					Expect(decrypted).To(BeEmpty())
				})

				It("should return error for invalid key during decryption", func() {
					// First encrypt with valid key
					encrypted, err := util.Encrypt(validKey16, plaintext)
					Expect(err).NotTo(HaveOccurred())

					// Try to decrypt with invalid key
					invalidKey := "wrong"
					decrypted, err := util.Decrypt(invalidKey, encrypted)
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("failed to create AES cipher"))
					Expect(decrypted).To(BeEmpty())
				})

				It("should fail to decrypt with wrong key", func() {
					// Encrypt with one key
					encrypted, err := util.Encrypt(validKey16, plaintext)
					Expect(err).NotTo(HaveOccurred())

					// Try to decrypt with a different valid key. AES-GCM authentication
					// must fail because the tag was produced under a different key.
					decrypted, err := util.Decrypt(validKey24, encrypted)
					Expect(err).To(HaveOccurred())
					Expect(decrypted).To(BeEmpty())
				})

				It("should fail to decrypt tampered ciphertext", func() {
					// Encrypt with a valid key
					encrypted, err := util.Encrypt(validKey16, plaintext)
					Expect(err).NotTo(HaveOccurred())

					// Tamper with the ciphertext (flip the last byte of the data).
					decoded, err := base64.URLEncoding.DecodeString(encrypted)
					Expect(err).NotTo(HaveOccurred())
					decoded[len(decoded)-1] ^= 0xFF
					tampered := base64.URLEncoding.EncodeToString(decoded)

					// AES-GCM authentication must detect the tampering.
					decrypted, err := util.Decrypt(validKey16, tampered)
					Expect(err).To(HaveOccurred())
					Expect(decrypted).To(BeEmpty())
				})

				It("should fail to decrypt random ciphertext", func() {
					// Random bytes are not a valid AES-GCM ciphertext, so
					// authentication must fail rather than returning garbage.
					random := make([]byte, 32)
					for i := range random {
						random[i] = byte(i)
					}
					randomCipher := base64.URLEncoding.EncodeToString(random)

					decrypted, err := util.Decrypt(validKey16, randomCipher)
					Expect(err).To(HaveOccurred())
					Expect(decrypted).To(BeEmpty())
				})
			})
		})

		Context("Edge Cases", func() {
			It("should handle very long text", func() {
				longText := strings.Repeat("A", 10000)
				encrypted, err := util.Encrypt(validKey32, longText)
				Expect(err).NotTo(HaveOccurred())

				decrypted, err := util.Decrypt(validKey32, encrypted)
				Expect(err).NotTo(HaveOccurred())
				Expect(decrypted).To(Equal(longText))
			})

			It("should handle single character", func() {
				singleChar := "X"
				encrypted, err := util.Encrypt(validKey16, singleChar)
				Expect(err).NotTo(HaveOccurred())

				decrypted, err := util.Decrypt(validKey16, encrypted)
				Expect(err).NotTo(HaveOccurred())
				Expect(decrypted).To(Equal(singleChar))
			})
		})
	})

	Context("Error Constants", func() {
		It("should have ErrCipherTooShort defined", func() {
			Expect(util.ErrCipherTooShort).To(HaveOccurred())
			Expect(util.ErrCipherTooShort.Error()).To(Equal("ciphertext too short"))
		})
	})
})
