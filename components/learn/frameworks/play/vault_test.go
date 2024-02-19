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
			value = map[string]interface{}{"password": "abc123", "secret": "correct horse battery staple"}

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
})
