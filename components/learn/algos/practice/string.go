package practice

/*
*
Replace space with %20 character as in encoding.
*/
func ReplaceSpace(input string) string {
	spaceCount := 0
	/* Find number of spaces in this string. */
	for _, c := range input {
		if c == ' ' {
			spaceCount++
		}
	}

	/* Create larger string accomoding for two extra charcter (in %20) */
	out := make([]rune, len(input)+spaceCount*2)

	i := 0
	/* Whenever space is encountered places %20 there */
	for _, c := range input {
		if c == ' ' {
			out[i] = '%'
			out[i+1] = '2'
			out[i+2] = '0'
			i += 3
		} else {
			out[i] = c
			i++
		}
	}

	return string(out)
}

func ReverseString(input string) string {
	var chars = []rune(input)
	var size = len(chars)

	// Place i,j at first and last postion
	// Keep Replacing Chars until the cross.
	for i, j := 0, size-1; i < j; i, j = i+1, j-1 {
		chars[i], chars[j] = chars[j], chars[i]
	}

	return string(chars)
}
