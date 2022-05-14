package tutorial

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type Day int

type rot13Reader struct {
	r io.Reader
}

func GoTour() {
	miscFun()
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
