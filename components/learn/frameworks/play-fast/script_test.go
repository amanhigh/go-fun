package play_fast

import (
	"regexp"
	"strings"

	"github.com/bitfield/script"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// https://github.com/bitfield/script
var _ = Describe("Script", func() {

	var (
		resourceDir = "../res"
		filePath    = resourceDir + "/play.md"
	)

	It("should read file", func() {
		contents, err := script.File(filePath).String()
		Expect(err).To(BeNil())
		Expect(contents).Should(ContainSubstring("Sample Markdown"))
	})

	It("should count lines", func() {
		lines, err := script.File(filePath).First(10).CountLines()
		Expect(err).To(BeNil())
		Expect(lines).Should(Equal(10))
	})

	It("should filter", func() {
		contents, err := script.File(filePath).Match("Sample").Tee().FilterLine(strings.ToUpper).String()
		Expect(err).To(BeNil())
		Expect(contents).Should(ContainSubstring("SAMPLE MARKDOWN"))
	})

	It("should curl", func() {
		contents, err := script.Get("https://google.com").String()
		Expect(err).To(BeNil())
		Expect(contents).Should(ContainSubstring("Google"))
	})

	It("should get first line", func() {
		line, err := script.File(filePath).First(1).String()
		Expect(err).To(BeNil())
		Expect(line).Should(Equal("# Sample Markdown File\n"))
	})

	It("should reject lines", func() {
		contents, err := script.File(filePath).Reject("#").String()
		Expect(err).To(BeNil())
		Expect(contents).ShouldNot(ContainSubstring("# Sample Markdown File"))
	})

	It("should get last line", func() {
		line, err := script.File(filePath).Last(1).String()
		Expect(err).To(BeNil())
		Expect(line).Should(ContainSubstring("3. Level 1 Item 3"))
	})

	It("should get column", func() {
		column, err := script.File(filePath).First(1).Column(3).String()
		Expect(err).To(BeNil())
		Expect(column).Should(Equal("Markdown\n"))
	})

	It("should reject lines with regexp", func() {
		contents, err := script.File(filePath).RejectRegexp(regexp.MustCompile("^#")).String()
		Expect(err).To(BeNil())
		Expect(contents).ShouldNot(ContainSubstring("# Sample Markdown File"))
	})

	It("should get SHA256sum", func() {
		sum, err := script.IfExists(filePath).SHA256Sum()
		Expect(err).To(BeNil())
		Expect(sum).ShouldNot(BeNil())
	})

	It("should get directory name", func() {
		dir, err := script.Echo(filePath).Dirname().String()
		Expect(err).To(BeNil())
		// Replace with the actual directory of your file
		Expect(dir).Should(Equal(resourceDir + "\n"))
	})

	It("should join lines", func() {
		joined, err := script.File(filePath).First(3).Join().String()
		Expect(err).To(BeNil())
		Expect(joined).Should(Equal("# Sample Markdown File Markdown Play ## Headers\n"))
	})

	It("should filter with custom function", func() {
		filtered, err := script.File(filePath).FilterLine(func(line string) (result string) {
			if strings.HasPrefix(line, "#") {
				result = line
			}
			return
		}).String()
		Expect(err).To(BeNil())
		Expect(filtered).Should(ContainSubstring("# Sample Markdown File"))
	})

	It("should concatenate files", func() {
		contents, err := script.ListFiles(resourceDir).Concat().String()
		Expect(err).To(BeNil())
		Expect(contents).Should(ContainSubstring("Sample Markdown"))
	})

	It("should execute command", func() {
		output, err := script.Exec("echo Hello, World!").String()
		Expect(err).To(BeNil())
		Expect(output).Should(Equal("Hello, World!\n"))
	})

	It("should list files", func() {
		files, err := script.ListFiles(resourceDir).Slice()
		Expect(err).To(BeNil())
		// Replace with the actual list of files in your directory
		Expect(files).Should(ContainElements([]string{filePath}))
	})

	It("should find files", func() {
		files, err := script.FindFiles(resourceDir).Slice()
		Expect(err).To(BeNil())
		// Replace with the actual list of files in your directory
		Expect(files).Should(ContainElements([]string{filePath}))
	})

	It("should wait for parallel commands", func() {
		script.Exec("echo Hello, World!").Wait()
		script.Exec("echo Foo, Bar!").Wait()
	})
})
