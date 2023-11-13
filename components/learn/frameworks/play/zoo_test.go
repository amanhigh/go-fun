package play_test

import (
	"time"

	"github.com/amanhigh/go-fun/models"
	"github.com/go-zookeeper/zk"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// TODO:Complete Migration from zoo.go
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

	// Context("Write", func() {
	// 	var (
	// 		testPath  = "/testPath"
	// 		testValue = "Test Value"
	// 		readValue []byte
	// 	)
	// 	BeforeEach(func() {
	// 		_, err = connection.Create(testPath, []byte(testValue), 0, zk.WorldACL(zk.PermAll))
	// 		Expect(err).To(BeNil())
	// 	})

	// 	It("should have Read Value", func() {
	// 		readValue, _, err = connection.Get(testPath)
	// 		Expect(string(readValue)).To(Equal(testValue))

	// 	})

	// 	Context("Delete Value", func() {
	// 		BeforeEach(func() {
	// 			err = connection.Delete(testPath, -1)
	// 			Expect(err).To(BeNil())
	// 		})

	// 		It("should have No Value", func() {
	// 			readValue, _, err = connection.Get(testPath)
	// 			Expect(string(readValue)).To(Equal(""))
	// 		})
	// 	})
	// })

})

/*

Not Working due to library disconnection unclear error no help

	fmt.Println("Wating for Events")
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(3)
	go pathWatcherSelect("/aman", c, &waitGroup)
	waitGroup.Wait()

func pathWatcher(path string, c *zk.Conn, wg *sync.WaitGroup) {
	for o, _, cha, _ := c.GetW(path); ; wg.Done() {
		o, _, cha, _ = c.GetW((<-cha).Path)
		fmt.Println("Event Processed:", string(o))
	}
}

func pathWatcherSelect(path string, c *zk.Conn, wg *sync.WaitGroup) {
	o, _, cha, _ := c.GetW(path)
	timeout := time.After(5 * time.Second)
	for {
		select {
		case e := <-cha:
			o, _, cha, _ = c.GetW(e.Path)
			fmt.Println("Event Processed:", string(o))
		case t := <-time.After(1 * time.Second):
			fmt.Println("No Event Recived:", t)
		case <-timeout:
			fmt.Println("Timeout Out")
		}
		wg.Done()
	}
}

*/
