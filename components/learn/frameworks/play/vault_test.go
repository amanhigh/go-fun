package play_test

import (
	"context"
	"time"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/models"
	"github.com/fatih/color"
	vault "github.com/hashicorp/vault-client-go"
	"github.com/hashicorp/vault-client-go/schema"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/testcontainers/testcontainers-go"
)

var _ = FDescribe("Vault", Ordered, Label(models.GINKGO_SLOW), func() {

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

			// Transit path
			path = "transit"
		)

		BeforeEach(func() {
			// Enable the transit secrets engine
			_, err = client.System.MountsEnableSecretsEngine(ctx, path, schema.MountsEnableSecretsEngineRequest{
				Type: "transit",
			})
			Expect(err).To(BeNil())

			// Create Key
			_, err = client.Secrets.TransitCreateKey(ctx, keyName, schema.TransitCreateKeyRequest{
				Exportable: true,
				Type:       keyType,
			})
			Expect(err).To(BeNil())

			// Configure Key to allow deletion
			_, err = client.Secrets.TransitConfigureKey(ctx, keyName, schema.TransitConfigureKeyRequest{
				DeletionAllowed: true,
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
	})

})
