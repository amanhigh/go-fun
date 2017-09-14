package main

import (
	"time"
	"fmt"
)

func main() {
	intChannel := make(chan int)
	go func() {
		for{
			time.Sleep(4 * time.Second)
			intChannel<-1
		}
	}()

	for value := range intChannel {
		fmt.Printf("Value:%+v\n", value)
	}
}
