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
