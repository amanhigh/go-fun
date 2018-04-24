package hackerrank

func LeftRotate(input []int, rotationCount int) (rotatedArray []int) {
	if len(input) >= rotationCount {
		suffix := input[:rotationCount]
		remaining := input[rotationCount:]
		rotatedArray = append(remaining, suffix...)
	}
	return
}
