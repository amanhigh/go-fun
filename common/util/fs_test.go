package util_test

import (
	"os"
	"path/filepath"

	"github.com/amanhigh/go-fun/common/util"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("FileSystem", func() {
	var (
		tempDir  string
		testFile string
	)

	BeforeEach(func() {
		var err error
		tempDir, err = os.MkdirTemp("", "fs_test_*")
		Expect(err).NotTo(HaveOccurred())
		testFile = filepath.Join(tempDir, "test.txt")
	})

	AfterEach(func() {
		os.RemoveAll(tempDir)
	})

	Context("OpenOrCreateFile", func() {
		It("should create new file if it doesn't exist", func() {
			file, err := util.OpenOrCreateFile(testFile)
			Expect(err).NotTo(HaveOccurred())
			Expect(file).NotTo(BeNil())
			defer file.Close()

			// Check file exists
			Expect(util.PathExists(testFile)).To(BeTrue())
		})

		It("should open existing file", func() {
			// Create file first
			err := os.WriteFile(testFile, []byte("existing content"), 0600)
			Expect(err).NotTo(HaveOccurred())

			file, err := util.OpenOrCreateFile(testFile)
			Expect(err).NotTo(HaveOccurred())
			Expect(file).NotTo(BeNil())
			defer file.Close()
		})

		It("should return error for invalid path", func() {
			invalidPath := "/invalid/nonexistent/path/file.txt"
			file, err := util.OpenOrCreateFile(invalidPath)
			Expect(err).To(HaveOccurred())
			Expect(file).To(BeNil())
		})
	})

	Context("AppendFile", func() {
		It("should append content to existing file", func() {
			// Create initial file
			err := os.WriteFile(testFile, []byte("initial"), 0600)
			Expect(err).NotTo(HaveOccurred())

			util.AppendFile(testFile, " appended")

			content, err := os.ReadFile(testFile)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(content)).To(Equal("initial appended"))
		})

		It("should handle non-existent file gracefully", func() {
			nonExistentFile := filepath.Join(tempDir, "nonexistent.txt")
			// Should not panic - logs error internally
			util.AppendFile(nonExistentFile, "content")
		})
	})

	Context("ReadAllLines", func() {
		It("should read all lines from file", func() {
			content := "line1\nline2\nline3"
			err := os.WriteFile(testFile, []byte(content), 0600)
			Expect(err).NotTo(HaveOccurred())

			lines := util.ReadAllLines(testFile)
			Expect(lines).To(HaveLen(3))
			Expect(lines[0]).To(Equal("line1"))
			Expect(lines[1]).To(Equal("line2"))
			Expect(lines[2]).To(Equal("line3"))
		})

		It("should return empty slice for non-existent file", func() {
			lines := util.ReadAllLines("nonexistent.txt")
			Expect(lines).To(BeEmpty())
		})

		It("should handle empty file", func() {
			err := os.WriteFile(testFile, []byte(""), 0600)
			Expect(err).NotTo(HaveOccurred())

			lines := util.ReadAllLines(testFile)
			Expect(lines).To(BeEmpty())
		})
	})

	Context("WriteLines", func() {
		It("should write lines to file", func() {
			lines := []string{"line1", "line2", "line3"}
			err := util.WriteLines(testFile, lines)
			Expect(err).NotTo(HaveOccurred())

			content, err := os.ReadFile(testFile)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(content)).To(Equal("line1\nline2\nline3"))
		})

		It("should handle empty slice", func() {
			lines := []string{}
			err := util.WriteLines(testFile, lines)
			Expect(err).NotTo(HaveOccurred())

			content, err := os.ReadFile(testFile)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(content)).To(Equal(""))
		})

		It("should return error for invalid path", func() {
			lines := []string{"test"}
			err := util.WriteLines("/invalid/path/file.txt", lines)
			Expect(err).To(HaveOccurred())
		})
	})

	Context("PathExists", func() {
		It("should return true for existing file", func() {
			err := os.WriteFile(testFile, []byte("test"), 0600)
			Expect(err).NotTo(HaveOccurred())

			Expect(util.PathExists(testFile)).To(BeTrue())
		})

		It("should return true for existing directory", func() {
			Expect(util.PathExists(tempDir)).To(BeTrue())
		})

		It("should return false for non-existent path", func() {
			Expect(util.PathExists("/nonexistent/path")).To(BeFalse())
		})
	})

	Context("ListFiles", func() {
		It("should list files in directory", func() {
			// Create test files
			file1 := filepath.Join(tempDir, "file1.txt")
			file2 := filepath.Join(tempDir, "file2.txt")
			err := os.WriteFile(file1, []byte("test1"), 0600)
			Expect(err).NotTo(HaveOccurred())
			err = os.WriteFile(file2, []byte("test2"), 0600)
			Expect(err).NotTo(HaveOccurred())

			files := util.ListFiles(tempDir)
			Expect(files).To(HaveLen(2))
			Expect(files).To(ContainElement(file1))
			Expect(files).To(ContainElement(file2))
		})

		It("should return empty slice for non-existent directory", func() {
			files := util.ListFiles("/nonexistent/directory")
			Expect(files).To(BeEmpty())
		})
	})

	Context("Copy", func() {
		It("should copy file content", func() {
			// Create source file
			sourceContent := "source file content"
			err := os.WriteFile(testFile, []byte(sourceContent), 0600)
			Expect(err).NotTo(HaveOccurred())

			// Copy to destination
			destFile := filepath.Join(tempDir, "dest.txt")
			err = util.Copy(testFile, destFile)
			Expect(err).NotTo(HaveOccurred())

			// Verify copy
			destContent, err := os.ReadFile(destFile)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(destContent)).To(Equal(sourceContent))
		})

		It("should return error for non-existent source", func() {
			err := util.Copy("nonexistent.txt", testFile)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failed to read source file"))
		})

		It("should return error for invalid destination", func() {
			err := os.WriteFile(testFile, []byte("test"), 0600)
			Expect(err).NotTo(HaveOccurred())

			err = util.Copy(testFile, "/invalid/path/dest.txt")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failed to write to destination file"))
		})
	})

	Context("FindReplaceFile", func() {
		It("should replace regex pattern in file", func() {
			content := "Hello World, this is a test World"
			err := os.WriteFile(testFile, []byte(content), 0600)
			Expect(err).NotTo(HaveOccurred())

			err = util.FindReplaceFile(testFile, "World", "Universe")
			Expect(err).NotTo(HaveOccurred())

			newContent, err := os.ReadFile(testFile)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(newContent)).To(Equal("Hello Universe, this is a test Universe"))
		})

		It("should handle invalid regex", func() {
			err := os.WriteFile(testFile, []byte("test"), 0600)
			Expect(err).NotTo(HaveOccurred())

			err = util.FindReplaceFile(testFile, "[invalid", "replacement")
			Expect(err).To(HaveOccurred())
		})

		It("should handle non-existent file", func() {
			err := util.FindReplaceFile("nonexistent.txt", "pattern", "replacement")
			Expect(err).To(HaveOccurred())
		})
	})

	Context("ReplaceContent", func() {
		It("should replace content using regex", func() {
			content := "version: 1.0.0\nother content"
			err := os.WriteFile(testFile, []byte(content), 0600)
			Expect(err).NotTo(HaveOccurred())

			util.ReplaceContent(testFile, `version: \d+\.\d+\.\d+`, "version: 2.0.0")

			newContent, err := os.ReadFile(testFile)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(newContent)).To(Equal("version: 2.0.0\nother content"))
		})

		It("should handle non-existent file gracefully", func() {
			// Should not panic - logs error internally
			util.ReplaceContent("nonexistent.txt", "pattern", "replacement")
		})
	})

	Context("RecreateDir", func() {
		It("should remove and recreate directory", func() {
			// Create subdirectory with files
			subDir := filepath.Join(tempDir, "subdir")
			err := os.MkdirAll(subDir, 0755)
			Expect(err).NotTo(HaveOccurred())

			testFileInSub := filepath.Join(subDir, "test.txt")
			err = os.WriteFile(testFileInSub, []byte("test"), 0600)
			Expect(err).NotTo(HaveOccurred())

			// Recreate directory
			util.RecreateDir(subDir)

			// Directory should exist but be empty
			Expect(util.PathExists(subDir)).To(BeTrue())
			files := util.ListFiles(subDir)
			Expect(files).To(BeEmpty())
		})
	})

	Context("ClearDirectory", func() {
		It("should remove all files from directory", func() {
			// Create test files
			file1 := filepath.Join(tempDir, "file1.txt")
			file2 := filepath.Join(tempDir, "file2.txt")
			err := os.WriteFile(file1, []byte("test1"), 0600)
			Expect(err).NotTo(HaveOccurred())
			err = os.WriteFile(file2, []byte("test2"), 0600)
			Expect(err).NotTo(HaveOccurred())

			util.ClearDirectory(tempDir)

			files := util.ListFiles(tempDir)
			Expect(files).To(BeEmpty())
		})

		It("should handle non-existent directory gracefully", func() {
			// Should not panic
			util.ClearDirectory("/nonexistent/directory")
		})
	})

	Context("ReadFileMap", func() {
		It("should read files into map", func() {
			// Create test files
			file1 := filepath.Join(tempDir, "file1.txt")
			file2 := filepath.Join(tempDir, "file2.txt")
			err := os.WriteFile(file1, []byte("line1\nline2"), 0600)
			Expect(err).NotTo(HaveOccurred())
			err = os.WriteFile(file2, []byte("content"), 0600)
			Expect(err).NotTo(HaveOccurred())

			fileMap := util.ReadFileMap(tempDir, false)

			Expect(fileMap).To(HaveLen(2))
			Expect(fileMap[file1]).To(Equal([]string{"line1", "line2"}))
			Expect(fileMap[file2]).To(Equal([]string{"content"}))
		})

		It("should include empty files when readEmpty is true", func() {
			emptyFile := filepath.Join(tempDir, "empty.txt")
			err := os.WriteFile(emptyFile, []byte(""), 0600)
			Expect(err).NotTo(HaveOccurred())

			fileMap := util.ReadFileMap(tempDir, true)
			Expect(fileMap).To(HaveKey(emptyFile))
			Expect(fileMap[emptyFile]).To(BeEmpty())
		})

		It("should exclude empty files when readEmpty is false", func() {
			emptyFile := filepath.Join(tempDir, "empty.txt")
			nonEmptyFile := filepath.Join(tempDir, "nonempty.txt")
			err := os.WriteFile(emptyFile, []byte(""), 0600)
			Expect(err).NotTo(HaveOccurred())
			err = os.WriteFile(nonEmptyFile, []byte("content"), 0600)
			Expect(err).NotTo(HaveOccurred())

			fileMap := util.ReadFileMap(tempDir, false)
			Expect(fileMap).NotTo(HaveKey(emptyFile))
			Expect(fileMap).To(HaveKey(nonEmptyFile))
		})
	})

	Context("ReadAllFiles", func() {
		It("should read content from all files", func() {
			// Create test files
			file1 := filepath.Join(tempDir, "file1.txt")
			file2 := filepath.Join(tempDir, "file2.txt")
			err := os.WriteFile(file1, []byte("line1\nline2"), 0600)
			Expect(err).NotTo(HaveOccurred())
			err = os.WriteFile(file2, []byte("content"), 0600)
			Expect(err).NotTo(HaveOccurred())

			allLines := util.ReadAllFiles(tempDir)

			Expect(allLines).To(HaveLen(3)) // 2 lines from file1 + 1 line from file2
			Expect(allLines).To(ContainElement("line1"))
			Expect(allLines).To(ContainElement("line2"))
			Expect(allLines).To(ContainElement("content"))
		})
	})

	Context("Constants", func() {
		It("should have proper file permission constants", func() {
			Expect(util.DEFAULT_PERM).To(Equal(os.FileMode(0644)))
			Expect(util.DIR_DEFAULT_PERM).To(Equal(os.FileMode(0755)))
			Expect(util.APPEND_PERM).To(Equal(os.FileMode(0600)))
		})
	})
})
