package cracking

func LeftRotate(input []int, rotationCount int) (rotatedArray []int) {
	if len(input) >= rotationCount {
		suffix := input[:rotationCount]
		remaining := input[rotationCount:]
		rotatedArray = append(remaining, suffix...)
	}
	return
}

/**
Find nth Fibonacci number.
*/
var mem []int

func Fibonacci(n int) int {
	mem = make([]int, n+1)
	return FibonacciRecursive(n)
}

func FibonacciRecursive(n int) (result int) {
	memFib := mem[n]
	if n == 0 || n == 1 {
		result = n
	} else if memFib != 0 {
		result = memFib
	} else {
		result = FibonacciRecursive(n-1) + FibonacciRecursive(n-2)
		mem[n] = result
	}
	return
}

/**
1) Any number xor'd with itself will give zero.
2) Any number xor'd with zero will give the number.
3) We are told there is an odd number of numbers in the array and they are all pairs of the same number, apart from one.

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

/**
Complete the function kangaroo which takes starting location and speed of both kangaroos as input,
and return or appropriately.

Can you determine if the kangaroos will ever land at the same location at the same time?

https://www.hackerrank.com/challenges/kangaroo/problem

0 3 4 2 -> True
0 2 5 3 -> False
*/
func KangarooMeet(ints []int) bool {
	x1, v1, x2, v2 := ints[0], ints[1], ints[2], ints[3]
	/* x1<=x2  */
	initialLead := x2 - x1
	speedDifference := v1 - v2
	/* If both have same speed must start at same position */
	if speedDifference == 0 {
		return x1 == x2
	} else {
		/*
			If there is speed difference v1 should have higher speed because x1 <= x2.
			Initial lead should be able to cover in s steps only if lead%speedDiff==0
		*/
		return v1 >= v2 && initialLead%speedDifference == 0
	}
}
