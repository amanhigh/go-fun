package play_test

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/models"
	"github.com/fatih/color"
	vault "github.com/hashicorp/vault-client-go"
	"github.com/hashicorp/vault-client-go/schema"
	"github.com/hashicorp/vault/helper/dhutil"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/testcontainers/testcontainers-go"
)

var _ = Describe("Vault", Ordered, Label(models.GINKGO_SLOW), func() {

	var (
		client *vault.Client
		ctx    = context.Background()
		err    error

		vaultContainer testcontainers.Container
	)

	BeforeAll(func() {
		// Create Vault Test Container
		vaultContainer, err = util.VaultTestContainer(ctx)
		Expect(err).To(BeNil())

		// Get Mapped Port
		vaultHost, err := vaultContainer.PortEndpoint(ctx, "8200/tcp", "")
		vaultHost = "http://" + vaultHost
		Expect(err).To(BeNil())
		color.Green("Vault Endpoint: %s", vaultHost)

		// Get a new client
		client, err = vault.New(vault.WithAddress(vaultHost), vault.WithRequestTimeout(30*time.Second))
		Expect(err).To(BeNil())

		//Authenticate
		err = client.SetToken(models.VAULT_ROOT_TOKEN)
		Expect(err).To(BeNil())

		// Enable the transit secrets engine
		_, err = client.System.MountsEnableSecretsEngine(ctx, "transit", schema.MountsEnableSecretsEngineRequest{
			Type: "transit",
		})
		Expect(err).To(BeNil())
	})

	AfterAll(func() {
		color.Red("Vault Shutting Down")
		err = vaultContainer.Terminate(ctx)
		Expect(err).To(BeNil())
	})

	It("should connect", func() {
		Expect(client).To(Not(BeNil()))

		_, err = client.System.ReadHealthStatus(ctx)
		Expect(err).To(BeNil(), "Failed to connect to Vault")
	})

	Context("Secrets", func() {
		var (
			//Data
			key   = "foo"
			value = map[string]any{"password": "abc123", "secret": "correct horse battery staple"}

			// Vault path
			path = "secret/data/myapp"
		)

		BeforeEach(func() {
			_, err = client.Secrets.KvV2Write(ctx, key, schema.KvV2WriteRequest{Data: value}, vault.WithMountPath(path))
			Expect(err).To(BeNil())

		})

		AfterEach(func() {
			_, err = client.Secrets.KvV2Delete(ctx, key, vault.WithMountPath(path))
			Expect(err).To(BeNil())
		})

		// FIXME: Add list client.Logical().List("/secret/kv")
		It("should give correct Value on Read", func() {
			secret, err := client.Secrets.KvV2Read(ctx, key, vault.WithMountPath(path))
			Expect(err).To(BeNil())
			Expect(secret.Data.Data).To(Equal(value))
		})
	})

	Context("Transit", func() {
		var (
			//Data
			keyName = "aman-key"
			//rsa-4096 - Asymmetric, aes256-gcm96 - Symmetric
			keyType = "aes256-gcm96"
		)

		BeforeEach(func() {
			// Create Key
			_, err = client.Secrets.TransitCreateKey(ctx, keyName, schema.TransitCreateKeyRequest{
				Exportable: true,
				Type:       keyType,
			})
			Expect(err).To(BeNil())

			// Configure Key (Edit)
			_, err = client.Secrets.TransitConfigureKey(ctx, keyName, schema.TransitConfigureKeyRequest{
				DeletionAllowed:      true,
				AllowPlaintextBackup: true,
			})
			Expect(err).To(BeNil())
		})

		AfterEach(func() {
			// Delete Key
			_, err = client.Secrets.TransitDeleteKey(ctx, keyName)
			Expect(err).To(BeNil())
		})

		It("should read key info", func() {
			key, err := client.Secrets.TransitReadKey(ctx, keyName)
			Expect(err).To(BeNil())
			Expect(key.Data["type"]).To(Equal(keyType))
		})

		It("should list keys", func() {
			keys, err := client.Secrets.TransitListKeys(ctx)
			Expect(err).To(BeNil())
			Expect(keys.Data.Keys).Should(ContainElement(keyName))
		})

		It("should rotate key", func() {
			rotateKey, err := client.Secrets.TransitRotateKey(ctx, keyName, schema.TransitRotateKeyRequest{})
			Expect(err).To(BeNil())
			Expect(rotateKey.Data["latest_version"]).To(Equal(json.Number("2")))
		})

		It("should backup key", func() {
			backupKey, err := client.Secrets.TransitBackUpKey(ctx, keyName)
			Expect(err).To(BeNil())
			Expect(backupKey.Data["backup"]).ToNot(BeNil())
		})

		It("should export hmac key", func() {
			exportKey, err := client.Secrets.TransitExportKey(ctx, keyName, "hmac-key")
			Expect(err).To(BeNil())
			Expect(exportKey.Data["keys"]).To(HaveLen(1))
		})

		Context("Encryption", func() {
			var (
				plainText  = "aman's-secret"
				cipherText string
			)

			BeforeEach(func() {
				//Base 64 Encode
				baseData := base64.StdEncoding.EncodeToString([]byte(plainText))

				//Encrypt Data
				encryptedData, err := client.Secrets.TransitEncrypt(ctx, keyName, schema.TransitEncryptRequest{
					Plaintext: baseData,
				})
				Expect(err).To(BeNil())
				cipherText = encryptedData.Data["ciphertext"].(string)
				Expect(cipherText).ToNot(BeNil())
			})

			It("should decrypt data via vault", func() {
				decryptedData, err := client.Secrets.TransitDecrypt(ctx, keyName, schema.TransitDecryptRequest{
					Ciphertext: cipherText,
				})
				Expect(err).To(BeNil())
				decryptedBaseData := decryptedData.Data["plaintext"].(string)

				//Decode Base64 Data
				decryptedPlainText, err := base64.StdEncoding.DecodeString(decryptedBaseData)
				Expect(err).To(BeNil())
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
					Expect(err).To(BeNil())
					Expect(encryptionKey.Data["keys"]).To(HaveLen(1))
					baseKey := encryptionKey.Data["keys"].(map[string]any)["1"].(string)
					//Decode Base64 Key
					key, err = base64.StdEncoding.DecodeString(baseKey)
					Expect(err).To(BeNil())
				})

				Context("AES", func() {
					var (
						AAD    = []byte("additional authenticated data")
						noonce []byte
					)
					BeforeEach(func() {
						//Encrypt Data
						cipher, noonce, err = dhutil.EncryptAES(key, []byte(plainText), AAD)
						Expect(err).To(BeNil())
						Expect(cipher).ToNot(BeNil())
						Expect(noonce).ToNot(BeNil())
					})

					It("should decrypt data", func() {
						decryptedText, err := dhutil.DecryptAES(key, cipher, noonce, AAD)
						Expect(err).To(BeNil())
						Expect(string(decryptedText)).To(Equal(plainText))
					})

				})
			})
		})
	})
})
