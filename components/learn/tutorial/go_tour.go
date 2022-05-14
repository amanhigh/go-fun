package tutorial

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
)

type Day int

type rot13Reader struct {
	r io.Reader
}

func GoTour() {
	safeMapFun()
	miscFun()
	errorHandling()
	GoRoutineFun()
	StartCrawl()
}

func (r rot13Reader) Read(b []byte) (n int, e error) {
	/** Named Return Values instead of 'return n,e' */
	n, e = r.r.Read(b)
	for i, c := range b {
		if (c >= 'A' && c < 'N') || (c >= 'a' && c < 'n') {
			b[i] += 13
		} else if (c > 'M' && c <= 'Z') || (c > 'm' && c <= 'z') {
			b[i] -= 13
		}
	}
	return
}

//https://tour.golang.org/moretypes/18
func Pic(dx, dy int) [][]uint8 {
	result := make([][]uint8, dy)
	for i := 0; i < dy; i++ {
		result[i] = make([]uint8, dx)
		for j := 0; j < dx; j++ {
			result[i][j] = uint8(i ^ j)
		}
	}
	return result
}

func safeMapFun() {
	mapV := sync.Map{}
	mapV.Store("Aman", "Preet,Singh")
	mapV.Store("Aman1", "Preet,Done")
	fmt.Println(mapV.Load("Aman1"))
	mapV.Range(func(key, value any) bool {
		stringValue := value.(string)
		fmt.Printf("Key:%v Value:%v\n", key, strings.Split(stringValue, ","))
		return true
	})
}

func swap(a, b string) (string, string) {
	return b, a
}

func miscFun() {
	fmt.Println("\n\nMisc fun")
	a, b := swap("World", "Hello")
	fmt.Println("Multiple Returns:", a, b)

	/** Ascii To Byte */
	x := "A"[0]
	fmt.Printf("%T %v", x, x)

	/** Rot13 */
	fmt.Println("\n\n ")
	s := strings.NewReader("Lbh penpxrq gur pbqr!")
	r := rot13Reader{s}
	io.Copy(os.Stdout, &r)
}

func errorHandling() {
	fmt.Println("\n\nError Handling")
	ff()
	fmt.Println("Returned normally from ff.")
}

func ff() {
	defer func() {
		/** Recovers whatever value is put in Panic */
		if r := recover(); r != nil {
			fmt.Println("Recovered in ff", r)
		}
	}()

	fmt.Println("Calling gg.")
	gg(0)

	/** Below code is not called as after panicking control goes to deferred recover */
	fmt.Println("Returned normally from gg.")
}

func gg(i int) {
	if i > 3 {
		fmt.Println("Panicking!")
		panic(fmt.Sprintf("%v", i*10))
	}
	defer fmt.Println("Defer in gg", i)
	fmt.Println("Printing in gg", i)
	gg(i + 1)
}
