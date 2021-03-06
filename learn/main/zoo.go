package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

func main() {
	c, _, err := zk.Connect([]string{"127.0.0.1"}, time.Second) //*10)
	if err != nil {
		panic(err)
	}

	/** Add & Remove Test Path */
	tPath := "/testPath"
	c.Create(tPath, []byte("TestValue"), 0, zk.WorldACL(zk.PermAll))
	printPath(tPath, c)
	c.Delete(tPath, -1)
	printPath(tPath, c)

	/** Wait for 3 Events to Be Processed & Exit */
	fmt.Println("Wating for Events")
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(3)
	go pathWatcherSelect("/aman", c, &waitGroup)
	waitGroup.Wait()
}

func printPath(path string, c *zk.Conn) {
	o, _, e := c.Get(path)
	fmt.Println("TestValue:", string(o), "Error:", e)
}

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
