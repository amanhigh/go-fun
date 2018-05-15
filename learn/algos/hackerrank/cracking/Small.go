package cracking

func LeftRotate(input []int, rotationCount int) (rotatedArray []int) {
	if len(input) >= rotationCount {
		suffix := input[:rotationCount]
		remaining := input[rotationCount:]
		rotatedArray = append(remaining, suffix...)
	}
	return
}

var mem []int

//mem = make([]int, n+1)

func Fibonacci(n int) (result int) {
	memFib := mem[n]
	if n == 0 || n == 1 {
		result = n
	} else if memFib != 0 {
		result = memFib
	} else {
		result = Fibonacci(n-1) + Fibonacci(n-2)
		mem[n] = result
	}
	return
}

/**
1) Any number xor'd with itself will give zero. 2) Any number xor'd with zero will give the number. 3) We are told there is an odd number of numbers in the array and they are all pairs of the same number, apart from one.

So if we xor all the numbers in the array together then any which are the same will cancel out - and give zero as the result of all the xors.

Then we are left with the unique number, which xor's with zero and so gives the unique number as the answer.

Eg. 1,1,2 -> Answer: 2
*/
func FindLonely(ints []int) int {
	lonelyInt := 0
	for _, i := range ints {
		lonelyInt ^= i
	}
	return lonelyInt
}
