package play_test

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/models"
	vault "github.com/hashicorp/vault-client-go"
	"github.com/hashicorp/vault-client-go/schema"
	"github.com/hashicorp/vault/helper/dhutil"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	"github.com/testcontainers/testcontainers-go"
)

var _ = Describe("Vault", Ordered, Label(models.GINKGO_SLOW), func() {

	var (
		client         *vault.Client
		ctx            = context.Background()
		err            error
		vaultHost      string
		vaultContainer testcontainers.Container
	)

	BeforeAll(func() {
		// Create Vault Test Container
		vaultContainer, err = util.VaultTestContainer(ctx)
		Expect(err).ToNot(HaveOccurred())

		// Get Mapped Port
		vaultHost, err = vaultContainer.PortEndpoint(ctx, "8200/tcp", "")
		vaultHost = "http://" + vaultHost
		Expect(err).ToNot(HaveOccurred())
		log.Info().Str("Host", vaultHost).Msg("Vault Endpoint")

		// Get a new client
		client, err = vault.New(vault.WithAddress(vaultHost), vault.WithRequestTimeout(30*time.Second))
		Expect(err).ToNot(HaveOccurred())

		// Authenticate
		err = client.SetToken(models.VAULT_ROOT_TOKEN)
		Expect(err).ToNot(HaveOccurred())

		// Enable the transit secrets engine
		_, err = client.System.MountsEnableSecretsEngine(ctx, "transit", schema.MountsEnableSecretsEngineRequest{
			Type: "transit",
		})
		Expect(err).ToNot(HaveOccurred())
	})

	AfterAll(func() {
		log.Warn().Msg("Vault Shutting Down")
		err = vaultContainer.Terminate(ctx)
		Expect(err).ToNot(HaveOccurred())
	})

	It("should connect", func() {
		Expect(client).To(Not(BeNil()))

		_, err = client.System.ReadHealthStatus(ctx)
		Expect(err).ToNot(HaveOccurred(), "Failed to connect to Vault")
	})
	Context("Secrets", func() {
		var (
			// Data
			key   = "foo"
			value = map[string]any{"password": "abc123", "secret": "correct horse battery staple"}

			// Vault paths
			mountPath = "secret"
			dataPath  = "myapp"
		)

		BeforeEach(func() {
			_, err = client.Secrets.KvV2Write(ctx, dataPath+"/"+key, schema.KvV2WriteRequest{Data: value}, vault.WithMountPath(mountPath))
			Expect(err).ToNot(HaveOccurred())
		})

		AfterEach(func() {
			_, err = client.Secrets.KvV2Delete(ctx, dataPath+"/"+key, vault.WithMountPath(mountPath))
			Expect(err).ToNot(HaveOccurred())
		})

		It("should give correct Value on Read", func() {
			secret, err := client.Secrets.KvV2Read(ctx, dataPath+"/"+key, vault.WithMountPath(mountPath))
			Expect(err).ToNot(HaveOccurred())
			Expect(secret.Data.Data).To(Equal(value))
		})

		It("should list secrets", func() {

			secrets, err := client.Secrets.KvV2List(ctx, dataPath, vault.WithMountPath(mountPath))
			Expect(err).ToNot(HaveOccurred())
			Expect(secrets.Data.Keys).To(ContainElement(key))
		})
	})

	Context("Transit", func() {
		var (
			// Data
			keyName = "aman-key"
			// rsa-4096 - Asymmetric, aes256-gcm96 - Symmetric
			keyType = "aes256-gcm96"
		)

		BeforeEach(func() {
			// Create Key
			_, err = client.Secrets.TransitCreateKey(ctx, keyName, schema.TransitCreateKeyRequest{
				Exportable: true,
				Type:       keyType,
			})
			Expect(err).ToNot(HaveOccurred())

			// Configure Key (Edit)
			_, err = client.Secrets.TransitConfigureKey(ctx, keyName, schema.TransitConfigureKeyRequest{
				DeletionAllowed:      true,
				AllowPlaintextBackup: true,
			})
			Expect(err).ToNot(HaveOccurred())
		})

		AfterEach(func() {
			// Delete Key
			_, err = client.Secrets.TransitDeleteKey(ctx, keyName)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should read key info", func() {
			key, err := client.Secrets.TransitReadKey(ctx, keyName)
			Expect(err).ToNot(HaveOccurred())
			Expect(key.Data["type"]).To(Equal(keyType))
		})

		It("should list keys", func() {
			keys, err := client.Secrets.TransitListKeys(ctx)
			Expect(err).ToNot(HaveOccurred())
			Expect(keys.Data.Keys).Should(ContainElement(keyName))
		})

		It("should rotate key", func() {
			rotateKey, err := client.Secrets.TransitRotateKey(ctx, keyName, schema.TransitRotateKeyRequest{})
			Expect(err).ToNot(HaveOccurred())
			Expect(rotateKey.Data["latest_version"]).To(Equal(json.Number("2")))
		})

		It("should backup key", func() {
			backupKey, err := client.Secrets.TransitBackUpKey(ctx, keyName)
			Expect(err).ToNot(HaveOccurred())
			Expect(backupKey.Data["backup"]).ToNot(BeNil())
		})

		It("should export hmac key", func() {
			exportKey, err := client.Secrets.TransitExportKey(ctx, keyName, "hmac-key")
			Expect(err).ToNot(HaveOccurred())
			Expect(exportKey.Data["keys"]).To(HaveLen(1))
		})

		Context("Encryption", func() {
			var (
				plainText  = "aman's-secret"
				cipherText string
			)

			BeforeEach(func() {
				// Base 64 Encode
				baseData := base64.StdEncoding.EncodeToString([]byte(plainText))

				// Encrypt Data
				encryptedData, err := client.Secrets.TransitEncrypt(ctx, keyName, schema.TransitEncryptRequest{
					Plaintext: baseData,
				})
				Expect(err).ToNot(HaveOccurred())
				cipherText = encryptedData.Data["ciphertext"].(string)
				Expect(cipherText).ToNot(BeNil())
			})

			It("should decrypt data via vault", func() {
				decryptedData, err := client.Secrets.TransitDecrypt(ctx, keyName, schema.TransitDecryptRequest{
					Ciphertext: cipherText,
				})
				Expect(err).ToNot(HaveOccurred())
				decryptedBaseData := decryptedData.Data["plaintext"].(string)

				// Decode Base64 Data
				decryptedPlainText, err := base64.StdEncoding.DecodeString(decryptedBaseData)
				Expect(err).ToNot(HaveOccurred())
				Expect(string(decryptedPlainText)).To(Equal(plainText))
			})

			Context("Export Encryption Key", func() {
				var (
					key    []byte
					cipher []byte
				)

				BeforeEach(func() {
					//Export Key
					encryptionKey, err := client.Secrets.TransitExportKey(ctx, keyName, "encryption-key")
					Expect(err).ToNot(HaveOccurred())
					Expect(encryptionKey.Data["keys"]).To(HaveLen(1))
					baseKey := encryptionKey.Data["keys"].(map[string]any)["1"].(string)
					// Decode Base64 Key
					key, err = base64.StdEncoding.DecodeString(baseKey)
					Expect(err).ToNot(HaveOccurred())
				})

				Context("AES", func() {
					var (
						AAD    = []byte("additional authenticated data")
						noonce []byte
					)
					BeforeEach(func() {
						// Encrypt Data
						cipher, noonce, err = dhutil.EncryptAES(key, []byte(plainText), AAD)
						Expect(err).ToNot(HaveOccurred())
						Expect(cipher).ToNot(BeNil())
						Expect(noonce).ToNot(BeNil())
					})

					It("should decrypt data", func() {
						decryptedText, err := dhutil.DecryptAES(key, cipher, noonce, AAD)
						Expect(err).ToNot(HaveOccurred())
						Expect(string(decryptedText)).To(Equal(plainText))
					})

				})
			})
		})
	})
})
