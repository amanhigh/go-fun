package core_test

import (
	"path/filepath"

	"github.com/amanhigh/go-fun/common/util"
	kohan "github.com/amanhigh/go-fun/models/config"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Tax Integration", Label("it"), func() {
	var (
		config  kohan.KohanConfig
		tempDir string
	)

	BeforeEach(func() {
		tempDir = GinkgoT().TempDir()
		err := setupTestFiles(tempDir)
		Expect(err).ToNot(HaveOccurred())

		config = kohan.KohanConfig{
			Tax: kohan.TaxConfig{
				DownloadsDir: tempDir,
			},
		}
	})

	Context("Tax Processing", func() {
		// TODO: Implement Integration Test
		It("should process complete tax flow", func() {
			// Placeholder for future implementation
			Expect(config.Tax.DownloadsDir).To(Equal(tempDir))
		})
	})
})

func setupTestFiles(tempDir string) error {
	// Source directory for test files
	srcDir := filepath.Join("testdata", "tax")

	// Get all files from source directory
	files, err := filepath.Glob(filepath.Join(srcDir, "*.csv"))
	if err != nil {
		return err
	}

	// Copy each file to temp directory
	for _, file := range files {
		destPath := filepath.Join(tempDir, filepath.Base(file))
		if err := util.Copy(file, destPath); err != nil {
			return err
		}
	}

	return nil
}
