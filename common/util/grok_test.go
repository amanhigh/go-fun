package util_test

import (
	"github.com/amanhigh/go-fun/common/util"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Grok", func() {
	Context("ReplaceRegEx", func() {
		It("should replace text using regex patterns (SDK wrapper)", func() {
			// Basic replacement
			result := util.ReplaceRegEx("Hello World", "World", "Universe")
			Expect(result).To(Equal("Hello Universe"))

			// Multiple occurrences
			result = util.ReplaceRegEx("foo bar foo", "foo", "qux")
			Expect(result).To(Equal("qux bar qux"))

			// Regex patterns
			result = util.ReplaceRegEx("abc123def456", `\d+`, "NUM")
			Expect(result).To(Equal("abcNUMdefNUM"))

			// Groups
			result = util.ReplaceRegEx("John Smith", `(\w+) (\w+)`, "$2, $1")
			Expect(result).To(Equal("Smith, John"))

			// Edge cases
			Expect(util.ReplaceRegEx("", "anything", "something")).To(Equal(""))
			Expect(util.ReplaceRegEx("no match", "xyz", "abc")).To(Equal("no match"))
		})
	})

	Context("GoGrep", func() {
		It("should extract matching portions from lines (not full lines)", func() {
			// GoGrep returns only the matched portion, not the full line
			input := "line1\nline2 match\nline3\nline4 match"
			result := util.GoGrep(input, "match")
			Expect(result).To(Equal("match\nmatch\n"))

			// No matches
			result = util.GoGrep("line1\nline2\nline3", "nomatch")
			Expect(result).To(Equal(""))

			// Empty input
			result = util.GoGrep("", "anything")
			Expect(result).To(Equal(""))

			// Single line
			result = util.GoGrep("single line with match", "match")
			Expect(result).To(Equal("match\n"))

			// Regex patterns - returns first match per line
			input = "line1\nline with 123\nline with abc\nline with 456"
			result = util.GoGrep(input, `\d+`)
			Expect(result).To(Equal("1\n123\n456\n"))

			// Word boundaries
			input = "cat\ncatch\nthe cat is here\ncatnip"
			result = util.GoGrep(input, `\bcat\b`)
			Expect(result).To(Equal("cat\ncat\n"))

			// Anchors
			input = "start here\nthis starts\nend here\nhere ends"
			result = util.GoGrep(input, `^start`)
			Expect(result).To(Equal("start\n"))

			result = util.GoGrep(input, `ends$`)
			Expect(result).To(Equal("ends\n"))
		})

		It("should handle special cases", func() {
			// Special regex characters
			input := "line1\nline.2\nline*3\nline+4"
			result := util.GoGrep(input, `line\.`)
			Expect(result).To(Equal("line.\n"))

			// Unicode
			input = "English line\nこんにちは世界\nAnother line"
			result = util.GoGrep(input, "こんにちは")
			Expect(result).To(Equal("こんにちは\n"))
		})
	})

	Context("Integration", func() {
		It("should work together for text processing", func() {
			input := "Error: 404 - Not found\nInfo: 200 - OK\nWarning: 500 - Server error"

			// Extract error/warning lines (get matching portions)
			errorLines := util.GoGrep(input, `(Error|Warning):`)
			Expect(errorLines).To(Equal("Error:\nWarning:\n"))

			// Replace status codes
			processed := util.ReplaceRegEx(errorLines, `\d{3}`, "[STATUS]")
			Expect(processed).To(Equal("Error:\nWarning:\n")) // No status codes in extracted portions
		})
	})
})
