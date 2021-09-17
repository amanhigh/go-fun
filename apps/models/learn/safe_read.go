package learn

import (
	"fmt"
	"time"
)

type SafeReadWrite struct {
	I    int
	Intc chan int
}

func (self *SafeReadWrite) Write(i int) {
	self.Intc <- i
}

func (self *SafeReadWrite) Close() {
	close(self.Intc)
}

func (self *SafeReadWrite) Read() (val int) {
	select {
	case v, ok := <-self.Intc:
		//If Channel is Not Closed Update I
		if ok {
			fmt.Println("Channel Written (400 ms)", time.Now().UnixMilli(), v)
			//Update New Value in Cache
			self.I = v
		}
		//Serve Updated I
		val = self.I
	default:
		//Runs if no other event is ready run
		time.Sleep(100 * time.Millisecond)
		//Serve Cached I
		val = self.I
	}
	return
}
