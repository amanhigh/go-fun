package tutorial

import (
	"fmt"
	"time"
)

type Person struct {
	Name string
	Age  int
}

func Syntax() {
	/* Vars */
	var arr = [3]int{1, 2, 3}
	var mapv = map[string]int{}
	fmt.Println("Array", "Map", arr, mapv)

	/* Struct */
	fmt.Println(Person{
		Name: "Aman",
		Age:  27,
	})

	/* Conditionals */
	x := 75
	if x > 50 {
		fmt.Println("x is greater than 50")
	} else if x < 50 {
		fmt.Println("x is less than 50")
	} else {
		fmt.Println("x is equal to 50")
	}

	/* Loop */
	count := 2
	for i := 0; i < count; i++ {
		fmt.Println("ILoop", i)
	}

	numbers := []int{2, 4}

	for index, value := range numbers {
		fmt.Printf("RangeLoop: Index: %d, Value: %d\n", index, value)
	}

	/* Channels */
	c := make(chan int) // ch <-chan int (Recive only Channel)
	// Start a goroutine for the printNumbers function
	go printNumbers(c)

	// Receive values from the channel and print them
	for v := range c {
		fmt.Println("Channel", v)
	}
}

func printNumbers(c chan int) {
	for i := 1; i <= 2; i++ {
		time.Sleep(100 * time.Millisecond)
		c <- i
	}
	close(c)
}
