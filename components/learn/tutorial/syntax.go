package tutorial

import (
	"fmt"
	"sync"
	"time"
)

type PersonInterface interface {
	Name() string
	Age() int
}

type Person struct {
	Name string
	Age  int
	mu   sync.Mutex
}

func (p *Person) GetName() string {
	defer p.mu.Unlock()
	p.mu.Lock()
	return p.Name
}

func Syntax() {
	demonstrateVariables()
	demonstrateStruct()
	demonstrateConditionals()
	demonstrateLoops()
	demonstrateChannels()
}

func demonstrateVariables() {
	var arr = [3]int{1, 2, 3}   //make([]int, 3)
	var mapv = map[string]int{} //make(map[string]int)
	fmt.Println("Array", "Map", arr, mapv)
}

func demonstrateStruct() {
	fmt.Println(Person{
		Name: "Aman",
		Age:  27,
	})
}

func demonstrateConditionals() {
	x := 75
	if x > 50 {
		fmt.Println("x is greater than 50")
	} else if x < 50 {
		fmt.Println("x is less than 50")
	} else {
		fmt.Println("x is equal to 50")
	}
}

func demonstrateLoops() {
	count := 2
	for i := 0; i < count; i++ {
		fmt.Println("ILoop", i)
	}

	arr := [3]int{1, 2, 3}
	numbers := arr[1:]
	for index, value := range numbers {
		fmt.Printf("RangeLoop: Index: %d, Value: %d\n", index, value)
	}
}

func demonstrateChannels() {
	// ch <-chan int (Recive only Channel)
	c := make(chan int)
	// Start a goroutine for the printNumbers function
	go printNumbers(c)

	// Receive values from the channel and print them
	for v := range c {
		fmt.Println("Channel", v)
	}
}

func printNumbers(c chan int) {
	defer close(c)

	for i := 1; i <= 2; i++ {
		time.Sleep(100 * time.Millisecond)
		c <- i
	}
}
