package play_test

import (
	"github.com/amanhigh/go-fun/models"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/samuel/go-zookeeper/zk"
	"time"
)

//TODO:Complete Migration from zoo.go
var _ = Describe("Zookeeper", Label(models.GINKGO_SETUP), func() {
	var (
		connection *zk.Conn
		err        error
		//dman set zookeeper
		zkHost = "docker"
	)
	BeforeEach(func() {
		connection, _, err = zk.Connect([]string{zkHost}, time.Second) //*10)
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

		It("should have Read Value", func() {
			readValue, _, err = connection.Get(testPath)
			Expect(string(readValue)).To(Equal(testValue))

		})

		Context("Delete Value", func() {
			BeforeEach(func() {
				err = connection.Delete(testPath, -1)
				Expect(err).To(BeNil())
			})

			It("should have No Value", func() {
				readValue, _, err = connection.Get(testPath)
				Expect(string(readValue)).To(Equal(""))
			})
		})
	})

})
