package tutorial

import (
	"fmt"
	"io"
	"math"
	"math/rand"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

var global, second_global = 5, 10

var (
	r, i rune = 8, 9
	g         = 0.867 + 0.5i // complex128
)

const (
	// Create a huge number by shifting a 1 bit left 100 places.
	// In other words, the binary number that is 1 followed by 100 zeroes.
	Big = 1 << 100
	// Shift it right again 99 places, so we end up with 1<<1, or 2.
	Small = Big >> 99
)

type Day int

const (
	MONDAY = 1 + iota
	TUESDAY
	WEDNESDAY
	THURSDAY
	FRIDAY
	SATURDAY
	SUNDAY
)

type Vertex struct {
	X, Y float64
}

type rot13Reader struct {
	r io.Reader
}

func GoTour() {
	fmt.Println("Yay :D :D !")
	fmt.Println("The time is", time.Now())

	variableFun()
	safeMapFun()
	mathFun()
	miscFun()
	regexFun()
	collectionFun()
	loopFun()
	switchFun()
	pointerFun()
	errorHandling()
	lambdaFun()
	GoRoutineFun()
	StartCrawl()
}

func regexFun() {
	fmt.Println("\n\nRegex Fun")
	s := "aman:aman@tcp(mysql:3306)/compute?charset=utf8&parseTime=True&loc=Local"
	m := regexp.MustCompile("^(.*)\\((.*)\\)(.*)$")
	fmt.Println(m.FindAllStringSubmatch(s, 5))
	fmt.Println(m.ReplaceAllString(s, `$1#$2#$3`))
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

func print(convert_function convert, x int) {
	fmt.Println("Lamda Sum:", convert_function(x))
}

type convert func(int) int

// convert types take an int and return a string value.
func lambdaFun() {
	fmt.Println("\n\nLamda Fun")
	print(double, 5)
	print(triple, 5)
	closure()
}
func closure() {
	fmt.Println("Closure")
	pos, neg := adder(), adder()
	for i := 0; i < 10; i++ {
		fmt.Println(pos(i), neg(-2*i))
	}

	fmt.Println("Fibonacci")
	f := fibonacciRecurse()
	for i := 0; i < 10; i++ {
		fmt.Print(f(), ",")
	}
}

// fibonacci is a function that returns
// a function that returns an int.
func fibonacciRecurse() func() int {
	lastFibBeforeUpdate, lastFib := 0, 0
	fib := 1
	return func() int {
		lastFibBeforeUpdate, lastFib, fib = lastFib, fib, lastFib+fib // Simultaneous Assignment :D
		return lastFibBeforeUpdate
	}
}

func adder() func(int) int {
	sum := 0
	return func(x int) int {
		sum += x
		return sum
	}
}

func double(i int) int {
	return i + i
}

func triple(i int) int {
	return i * 3
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

func loopFun() {
	fmt.Println("\n\nLoop Fun")
	sum := 0
	for i := 0; i < 10; i++ {
		sum += i
	}

	//Infinite Loop
	for {
		/**Precondition Before If**/
		if sum += 50; sum > 300 {
			break
		} else {
			//Precondition Variable available Here
			fmt.Println("Looping Infinite:", sum)
		}
	}

	fmt.Println("Sum:", sum)
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

func printRandom() {
	rand.Seed(time.Now().UnixNano())
	fmt.Println("My Random number is", rand.Intn(10))
}

func add(x, y int) int {
	return x + y
}

func swap(a, b string) (string, string) {
	return b, a
}

func variableFun() {
	fmt.Println("\n\nStarting Variable Fun")
	var local string = "localvariable"
	shortHand := "Shorthand Variable"
	const c string = "Constant"
	fmt.Println("Variables:", global, second_global, local, shortHand)
	fmt.Println("Constants:", c)

	fmt.Printf("Variables Type: %T Value: %v\n", r, r)
	fmt.Printf("Variables Type: %T Value: %v\n", i, i)
	fmt.Printf("Variables Type: %T Value: %v\n", g, g)

	vertexFun()

	fmt.Println("Exported Name Test:", math.Pi)

	typeCheck()

	fmt.Println("Enums")
	fmt.Println(MONDAY, MONDAY == 1)

	/** Returns Control */
	i, err := strconv.Atoi("XX")
	if err != nil {
		fmt.Printf("couldn't convert number: %v\n", err)
		return
	}
	fmt.Println("Converted integer:", i)
}

func vertexFun() {
	fmt.Println("\n\nVertexFun")
	vertex := Vertex{1, 2}
	ver := vertex
	pver := &vertex
	ver.Y = 7
	pver.Y = 9
	fmt.Println("Vertex:", vertex)
	fmt.Println("Vertex (By Value):", ver)
	fmt.Println("Vertex (By Reference):", pver)

	/** Interface */
	var a Abser
	//a=vertex /** Gives Error as Abs takes only Pointer */
	a = pver

	/* While methods with pointer receivers take either a value or a pointer as the receiver when they are called: */
	fmt.Println("Vertex Method: ", vertex.Abs(), pver.Abs(), a.Abs()) //Methods linked to Struct

	/** Null Handling pver.Abs() Would still work but when Abs will try to access X,Y Null Pointer would come. */
	//pver=nil;pver.Abs();
	/** Null Handling on Type would be error even if error is called on concrete type */
	//vertex=nil;vertex.Abs()
}

type Abser interface {
	Abs() float64
}

func (v *Vertex) Abs() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func mathFun() {
	fmt.Println("\n\nMath Fun")
	printRandom()
	fmt.Println("Adding:", add(5, 8))

	fmt.Println("An untyped constant takes the type needed by its context")
	//Small is 2 and Big is 1^100
	fmt.Println(needInt(Small))
	fmt.Println(needFloat(Small))
	fmt.Println(needFloat(Big))

	psqrt(-2)
	psqrt(8)
	psqrt(64)

}

func typeCheck() {

	fmt.Println("\n\nType Check")
	/** Empty Interface */
	var i any = "hello"

	s := i.(string) //Type Casting
	fmt.Println(s)

	//f := i.(float64) would raise panic as mismatching type.
	f, ok := i.(float64) // Test Statement
	fmt.Println(f, ok)
}

func needInt(x int) int {
	return x*10 + 1
}
func needFloat(x float64) float64 {
	return x * 0.1
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

type ErrNegativeSqrt float64

func psqrt(x int) {
	z, e := sqrt(x)
	if e != nil {
		fmt.Println("Error Computing Sqrt for", x)
		return
	}
	fmt.Println("Square Root of", x, "is", z, "Actual Root:", math.Sqrt(float64(x)))
}

func sqrt(x int) (float64, error) {
	if x < 0 {
		return 0, ErrNegativeSqrt(x)
	}
	fX := float64(x)
	z := float64(1)
	z = 1.0
	for i := 0; i < 10; i++ {
		z = z - ((math.Pow(z, 2) - fX) / (2 * z))
	}
	return z, nil
}

func (e ErrNegativeSqrt) Error() string {
	return fmt.Sprintf("cannot Sqrt negativ number: %g", float64(e))
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
