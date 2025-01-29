package play_test

import (
	"context"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/models"
	"github.com/hashicorp/consul/api"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	"github.com/testcontainers/testcontainers-go"
)

var _ = Describe("Consul", Ordered, Label(models.GINKGO_SLOW), func() {

	var (
		client          *api.Client
		ctx             = context.Background()
		err             error
		consulHost      string
		consulContainer testcontainers.Container
	)

	BeforeAll(func() {
		// Create Consul Test Container
		consulContainer, err = util.ConsulTestContainer(ctx)
		Expect(err).To(BeNil())

		// Get Mapped Port
		consulHost, err = consulContainer.PortEndpoint(ctx, "8500/tcp", "")
		Expect(err).To(BeNil())
		log.Info().Str("Host", consulHost).Msg("Consul Endpoint")

		// Get a new client
		client, err = api.NewClient(&api.Config{Address: consulHost})
		Expect(err).To(BeNil())
	})

	AfterAll(func() {
		log.Warn().Msg("Consul Shutting Down")
		err = consulContainer.Terminate(ctx)
		Expect(err).To(BeNil())
	})

	It("should connect", func() {
		Expect(client).To(Not(BeNil()))

		_, err = client.Agent().Self()
		Expect(err).To(BeNil(), "Failed to connect to Consul")
	})

	Context("Write and Read", func() {
		var (
			kv *api.KV

			// Data
			key   = "aman/1"
			value = []byte("2000")
			p     = api.KVPair{Key: key, Value: value}
		)

		BeforeEach(func() {
			// Get a handle to the KV API
			kv = client.KV()

			_, err = kv.Put(&p, nil)
			Expect(err).To(BeNil())
		})

		AfterEach(func() {
			_, err = kv.Delete(key, nil)
			Expect(err).To(BeNil())
		})

		It("should have Read Value", func() {
			readKV, _, err := kv.Get(key, nil)
			Expect(err).To(BeNil())
			Expect(readKV.Value).To(Equal(value))
		})
	})

})
