package play_test

import (
	"context"
	"time"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/models"
	"github.com/go-zookeeper/zk"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	"github.com/testcontainers/testcontainers-go"
)

// https://github.com/EladLeev/go-zookeeper-examples
var _ = Describe("Zookeeper", Ordered, Label(models.GINKGO_SLOW), func() {
	var (
		connection *zk.Conn
		ctx        = context.Background()
		err        error

		zkContainer testcontainers.Container
	)
	BeforeAll(func() {
		// Create Zookeeper Test Container
		zkContainer, err = util.ZookeeperTestContainer(ctx)
		Expect(err).ToNot(HaveOccurred())

		// Get Mapped Port
		zkHost, err := zkContainer.PortEndpoint(ctx, "2181/tcp", "")
		Expect(err).ToNot(HaveOccurred())
		log.Info().Str("Host", zkHost).Msg("Zookeeper Endpoint")

		connection, _, err = zk.Connect([]string{zkHost}, time.Second)
		Expect(err).ToNot(HaveOccurred())
	})

	AfterAll(func() {
		log.Warn().Msg("Zookeeper Shutting Down")
		err = zkContainer.Terminate(ctx)
		Expect(err).ToNot(HaveOccurred())
	})

	It("should connect", func() {
		Expect(connection).To(Not(BeNil()))
	})

	Context("Write", func() {
		var (
			testPath  = "/testPath"
			testValue = "Test Value"
			readValue []byte
		)
		BeforeEach(func() {
			_, err = connection.Create(testPath, []byte(testValue), 0, zk.WorldACL(zk.PermAll))
			Expect(err).ToNot(HaveOccurred())
		})

		AfterEach(func() {
			err = connection.Delete(testPath, -1)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should have Read Value", func() {
			readValue, _, err = connection.Get(testPath)
			Expect(string(readValue)).To(Equal(testValue))

		})

		It("should check Exists", func() {
			exists, _, err := connection.Exists(testPath)
			Expect(err).ToNot(HaveOccurred())
			Expect(exists).To(BeTrue())
		})

		Context("Update", func() {
			var (
				updateValue = "Update Value"
			)
			BeforeEach(func() {
				_, err = connection.Set(testPath, []byte(updateValue), -1)
				Expect(err).ToNot(HaveOccurred())
			})

			It("should have Updated Value", func() {
				readValue, _, err = connection.Get(testPath)
				Expect(string(readValue)).To(Equal(updateValue))
			})
		})

		Context("Watch", func() {
			var (
				watchValue = "Watch Value"
			)

			It("should get events", func() {
				// Start Watching Path
				data, _, evtChan, err := connection.GetW(testPath)
				Expect(err).ToNot(HaveOccurred())
				Expect(string(data)).To(Equal(testValue))

				// Write to Path
				go connection.Set(testPath, []byte(watchValue), -1)

				Eventually(func() zk.EventType {
					// Select ensures Eventually Go Routine doesn't get stuck.
					// When No Write is done Refer Eventually Documentation as well.
					select {
					case e := <-evtChan:
						data, _, evtChan, _ = connection.GetW(e.Path)
						Expect(string(data)).To(Equal(watchValue))
						return e.Type
					case <-time.After(1 * time.Second): // Timeout if no event is received within 1 second
						return zk.EventNotWatching
					}
				}, "2s").Should(Equal(zk.EventNodeDataChanged)) // Fail the test if no event is received within 2 seconds
			})

			It("should not get events without writes", func() {
				_, _, evtChan, err := connection.GetW(testPath)
				Expect(err).ToNot(HaveOccurred())

				Eventually(evtChan).ShouldNot(Receive())
			})
		})
	})
})
