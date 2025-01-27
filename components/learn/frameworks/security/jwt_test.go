package security

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"

	"github.com/dgrijalva/jwt-go"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Jwt", func() {

	var (
		privateKey *rsa.PrivateKey
		publicKey  crypto.PublicKey
		err        error
	)

	BeforeEach(func() {
		// Generate RSA Key Pair
		privateKey, err = rsa.GenerateKey(rand.Reader, 2048)
		publicKey = privateKey.Public()
		Expect(err).To(BeNil())
	})

	It("should generate Private Key", func() {
		Expect(privateKey).To(Not(BeNil()))
	})

	It("should generate Public Key", func() {
		Expect(publicKey).To(Not(BeNil()))
	})

	Context("PEM", func() {
		var (
			privatePem, publicPem []byte
			unmarshalPrivate      *rsa.PrivateKey
			unmarshalPublic       crypto.PublicKey
		)

		BeforeEach(func() {
			// Marshal to Private Key To PEM Format
			privatePem = MarshalPrivateKey(privateKey)

			// Marshal to Public Key PEM Format
			publicPem = MarshalPublicKey(publicKey)
		})

		It("should Marshal Key", func() {
			Expect(privatePem).To(Not(BeNil()))
			Expect(publicPem).To(Not(BeNil()))
		})

		It("should unmarshal Private Key", func() {
			// Unmarshal Private Key Pem
			unmarshalPrivate, err = UnmarshalRSAPrivateKey(string(privatePem))
			Expect(err).To(BeNil())
			Expect(unmarshalPrivate).To(Equal(privateKey))
		})

		It("should unmarshal Public Key", func() {
			// Unmarshal Private Key Pem
			unmarshalPublic, err = UnmarshalRSAPublicKey(string(publicPem))
			Expect(err).To(BeNil())
			Expect(unmarshalPublic).To(Equal(publicKey))
		})
	})

	Context("Token", func() {
		var (
			token string
		)
		BeforeEach(func() {
			// 	//Generate Token String
			token, err = GenerateToken(privateKey)
			Expect(err).To(BeNil())

		})

		It("should generate", func() {
			Expect(token).To(Not(BeNil()))

		})

		It("should parse", func() {
			// Parse Token
			parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (i any, err error) {
				return publicKey, nil
			})
			Expect(err).To(BeNil())
			Expect(parsedToken).To(Not(BeNil()))

			By("Claims")

			claims := parsedToken.Claims.(jwt.MapClaims)
			Expect(claims).To(HaveKeyWithValue("purpose", "jwtfun"))
			Expect(claims).To(HaveKeyWithValue("iss", "aman"))
			Expect(claims.Valid()).To(BeNil())
		})
	})

})
