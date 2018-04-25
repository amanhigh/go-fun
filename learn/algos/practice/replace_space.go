package practice

func ReplaceSpace(input string) string {
	spaceCount := 0
	for _, c := range input {
		if c == ' ' {
			spaceCount++
		}
	}

	out := make([]rune, len(input)+spaceCount*2)

	i := 0
	for _, c := range input {
		if c == ' ' {
			out[i] = '%'
			out[i+1] = '2'
			out[i+2] = '0'
			i = i + 3
		} else {
			out[i] = rune(c)
			i++
		}
	}

	return string(out)
}
