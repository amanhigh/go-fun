package play_test

import (
	"context"
	"time"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/models"
	"github.com/fatih/color"
	"github.com/go-zookeeper/zk"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
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
		//Create Zookeeper Test Container
		zkContainer, err = util.ZookeeperTestContainer(ctx)
		Expect(err).To(BeNil())

		//Get Mapped Port
		zkHost, err := zkContainer.Endpoint(ctx, "")
		Expect(err).To(BeNil())
		color.Green("Zookeeper Endpoint: %s", zkHost)

		connection, _, err = zk.Connect([]string{zkHost}, time.Second)
		Expect(err).To(BeNil())
	})

	AfterAll(func() {
		color.Red("Zookeeper Shutting Down")
		err = zkContainer.Terminate(ctx)
		Expect(err).To(BeNil())
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
			Expect(err).To(BeNil())
		})

		AfterEach(func() {
			err = connection.Delete(testPath, -1)
			Expect(err).To(BeNil())
		})

		It("should have Read Value", func() {
			readValue, _, err = connection.Get(testPath)
			Expect(string(readValue)).To(Equal(testValue))

		})

		It("should check Exists", func() {
			exists, _, err := connection.Exists(testPath)
			Expect(err).To(BeNil())
			Expect(exists).To(BeTrue())
		})

		Context("Update", func() {
			var (
				updateValue = "Update Value"
			)
			BeforeEach(func() {
				_, err = connection.Set(testPath, []byte(updateValue), -1)
				Expect(err).To(BeNil())
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
				//Start Watching Path
				data, _, evtChan, err := connection.GetW(testPath)
				Expect(err).To(BeNil())
				Expect(string(data)).To(Equal(testValue))

				//Write to Path
				go connection.Set(testPath, []byte(watchValue), -1)

				Eventually(func() zk.EventType {
					// BUG: #C Eventually doesnt' timeout if write is ommitted.
					e := <-evtChan
					data, _, evtChan, _ = connection.GetW(e.Path)
					Expect(string(data)).To(Equal(watchValue))
					return e.Type
				}).Should(Equal(zk.EventNodeDataChanged))
			})

			It("should not get events without writes", func() {
				_, _, evtChan, err := connection.GetW(testPath)
				Expect(err).To(BeNil())

				Eventually(evtChan).ShouldNot(Receive())
			})
		})
	})
})
