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

// TODO:Complete Migration from zoo.go
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
	})

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
