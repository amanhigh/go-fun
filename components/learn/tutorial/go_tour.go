package tutorial

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

type Day int

type rot13Reader struct {
	r io.Reader
}

func GoTour() {
	safeMapFun()
	miscFun()
	collectionFun()
	switchFun()
	pointerFun()
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

func collectionFun() {
	var a [2]string
	a[0] = "Hello"
	a[1] = "World"
	fmt.Println(a)
	primes := [6]int{2, 3, 5, 7, 11, 13}
	fmt.Println(primes[1:4]) // Prints 3,5,7 (Its a Reference not Value Copy)

	/** Slice referencing the Array as no Size is Specified for Struct Array */
	s := []struct {
		i int
		b bool
	}{{2, true}, {3, false}, {5, true}, {7, true}, {11, false}, {13, true}}
	fmt.Println(s)
	s = s[2:4] // The capacity of a slice is the number of elements in the underlying array, counting from the first element in the slice.
	fmt.Printf("len=%d cap=%d %v\n", len(s), cap(s), s)

	/** Ranges where i is optional can use _,v */
	for i, v := range primes {
		fmt.Printf("2**%d = %d\n", i, v)
	}

	/** Two Dimensional */
	var twod [5][5]uint8 //Array 5x5
	twod[1][1] = 5
	fmt.Println("Two Dimensional:", twod[1][1])

	fmt.Printf("Len: %d, Cap: %d\n", len(twod), cap(twod))

	hashMap := map[string]int{"One": 1, "Two": 2}
	v2, ok := hashMap["Two"] //Ok Holds if element is present or not.
	fmt.Println("HashMap:", hashMap["One"], "-", v2, ok)

	fmt.Println("Make vs New", len(make([]int, 50, 100)), len(new([100]int)[0:50]))

	fmt.Println("WordCount:", WordCount("Hello World Hello Aman"))
}

func WordCount(input string) map[string]int {
	countMap := make(map[string]int)
	fmt.Println(countMap)
	fields := strings.Fields(input)
	for _, f := range fields {
		countMap[f] += 1 //No NPE :), No Init Required because entry value is primitive
	}
	return countMap
}

func pointerFun() {
	fmt.Println("\n\n Pointer Fun")
	i, j := 42, 2701

	p := &i         // point to i
	fmt.Println(p)  // Address of i (Value of p)
	fmt.Println(*p) // read i through the pointer
	*p = 21         // set i through the pointer
	fmt.Println(i)  // see the new value of i

	p = &j         // point to j
	*p = *p / 37   // divide j through the pointer
	fmt.Println(j) // see the new value of j
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

func switchFun() {
	fmt.Println("\n\nSwitch Fun")
	fmt.Print("Go runs on ")
	switch os := runtime.GOOS; os {
	case "darwin":
		fmt.Println("OS X.")
	//fallthrough
	case "linux":
		fmt.Println("Linux.")
	default:
		// freebsd, openbsd,
		// plan9, windows...
		fmt.Printf("%s.", os)
	}

	fmt.Println("When's Saturday?")
	today := time.Now().Weekday()
	switch time.Saturday {
	case today + 0:
		fmt.Println("Today.")
	case today + 1:
		fmt.Println("Tomorrow.")
	case today + 2:
		fmt.Println("In two days.")
	default:
		fmt.Println("Too far away.")
	}

	/** Emulates long if/else chains */
	t := time.Now()
	switch {
	case t.Hour() < 12:
		fmt.Println("Good morning!")
	case t.Hour() < 17:
		fmt.Println("Good afternoon.")
	default:
		fmt.Println("Good evening.")
	}
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

	deferFun()
}

func deferFun() {
	/** Defer */
	message := "Captured Argument"
	defer fmt.Println(message)
	message = "Now Changed"
	fmt.Println("Arguments were captured but function excutes post this. Current Value of Message:", message)

	fmt.Println("counting")
	for i := 0; i < 10; i++ {
		defer fmt.Println(i)
	}
	fmt.Println("done")

	fmt.Println("Value Returned from Defered Function is", deferReturn())
}

func deferReturn() (i int) {
	defer func() {
		i++
	}()
	return 4
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
