package main

import (
	"time"
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"sync"
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
	go pathWatcher("/aman", c, &waitGroup)
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
