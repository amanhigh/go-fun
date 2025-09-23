package util_test

import (
	"crypto/md5"
	"encoding/hex"

	"github.com/amanhigh/go-fun/common/util"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("MD5", func() {
	Context("GetMD5Hash", func() {
		It("should generate correct MD5 hashes (SDK wrapper)", func() {
			// Test known hash values
			Expect(util.GetMD5Hash("")).To(Equal("d41d8cd98f00b204e9800998ecf8427e"))
			Expect(util.GetMD5Hash("hello")).To(Equal("5d41402abc4b2a76b9719d911017c592"))

			// Verify it matches standard library
			testInput := "test string"
			hasher := md5.New()
			hasher.Write([]byte(testInput))
			expected := hex.EncodeToString(hasher.Sum(nil))
			Expect(util.GetMD5Hash(testInput)).To(Equal(expected))

			// Deterministic behavior
			hash1 := util.GetMD5Hash("consistent")
			hash2 := util.GetMD5Hash("consistent")
			Expect(hash1).To(Equal(hash2))

			// Format verification
			result := util.GetMD5Hash("any input")
			Expect(result).To(HaveLen(32))
			Expect(result).To(MatchRegexp("^[a-f0-9]{32}$"))
		})
	})

	Context("Md5Info Struct", func() {
		var md5Info *util.Md5Info

		BeforeEach(func() {
			md5Info = &util.Md5Info{
				Hash:     "test_hash",
				FileList: []string{},
				Count:    0,
			}
		})

		It("should handle file path additions correctly", func() {
			// Add single file
			md5Info.Add("/path/to/file.txt")
			Expect(md5Info.FileList).To(HaveLen(1))
			Expect(md5Info.FileList[0]).To(Equal("/path/to/file.txt"))
			Expect(md5Info.Count).To(Equal(1))

			// Add multiple files
			md5Info.Add("/another/file.txt")
			md5Info.Add("/third/file.txt")
			Expect(md5Info.FileList).To(HaveLen(3))
			Expect(md5Info.Count).To(Equal(3))

			// Maintains order
			Expect(md5Info.FileList[0]).To(Equal("/path/to/file.txt"))
			Expect(md5Info.FileList[1]).To(Equal("/another/file.txt"))
			Expect(md5Info.FileList[2]).To(Equal("/third/file.txt"))

			// Hash remains unchanged
			Expect(md5Info.Hash).To(Equal("test_hash"))
		})

		It("should handle edge cases", func() {
			// Empty paths
			md5Info.Add("")
			Expect(md5Info.FileList).To(ContainElement(""))
			Expect(md5Info.Count).To(Equal(1))

			// Duplicate paths
			md5Info.Add("duplicate")
			md5Info.Add("duplicate")
			Expect(md5Info.FileList).To(HaveLen(3)) // including empty from above
			Expect(md5Info.Count).To(Equal(3))
		})
	})
})
