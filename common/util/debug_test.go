package util_test

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/models/config"
)

var _ = Describe("Debug", func() {

	var tempDebugFile string

	BeforeEach(func() {
		// Use a temp file to avoid conflicts with real debug file
		tempDebugFile = "/tmp/test_kohandebug"
		// Clean up any existing test debug file
		os.Remove(tempDebugFile)
	})

	AfterEach(func() {
		// Clean up test debug file
		os.Remove(tempDebugFile)
		// Clean up real debug file in case it was created
		os.Remove(config.DEBUG_FILE)
	})

	Context("DebugControl", func() {
		It("should enable debug mode by creating debug file", func() {
			// Verify file doesn't exist initially
			Expect(util.PathExists(config.DEBUG_FILE)).To(BeFalse())

			util.DebugControl(true)

			// Verify file was created
			Expect(util.PathExists(config.DEBUG_FILE)).To(BeTrue())
		})

		It("should disable debug mode by removing debug file", func() {
			// First create the debug file
			util.DebugControl(true)
			Expect(util.PathExists(config.DEBUG_FILE)).To(BeTrue())

			// Then disable debug mode
			util.DebugControl(false)

			// Verify file was removed
			Expect(util.PathExists(config.DEBUG_FILE)).To(BeFalse())
		})

		It("should handle enabling debug mode when file already exists", func() {
			// Create file manually first
			err := os.WriteFile(config.DEBUG_FILE, []byte{}, 0600)
			Expect(err).NotTo(HaveOccurred())

			// Enable debug mode again
			util.DebugControl(true)

			// File should still exist
			Expect(util.PathExists(config.DEBUG_FILE)).To(BeTrue())
		})

		It("should handle disabling debug mode when file doesn't exist", func() {
			// Ensure file doesn't exist
			os.Remove(config.DEBUG_FILE)
			Expect(util.PathExists(config.DEBUG_FILE)).To(BeFalse())

			// Disable debug mode
			util.DebugControl(false)

			// Should not cause any errors
			Expect(util.PathExists(config.DEBUG_FILE)).To(BeFalse())
		})
	})

	Context("IsDebugMode", func() {
		It("should return true when debug file exists", func() {
			// Create debug file
			err := os.WriteFile(config.DEBUG_FILE, []byte{}, 0600)
			Expect(err).NotTo(HaveOccurred())

			result := util.IsDebugMode()
			Expect(result).To(BeTrue())
		})

		It("should return false when debug file doesn't exist and env var not set", func() {
			// Ensure file doesn't exist
			os.Remove(config.DEBUG_FILE)
			Expect(util.PathExists(config.DEBUG_FILE)).To(BeFalse())

			// Ensure env var is false (it's initialized as false in config)
			result := util.IsDebugMode()
			Expect(result).To(BeFalse())
		})

		It("should return true when KOHAN_DEBUG environment variable is true", func() {
			// Remove debug file to test env var only
			os.Remove(config.DEBUG_FILE)

			// Set the config variable (simulating env var being true)
			originalValue := config.KOHAN_DEBUG
			config.KOHAN_DEBUG = true
			defer func() { config.KOHAN_DEBUG = originalValue }()

			result := util.IsDebugMode()
			Expect(result).To(BeTrue())
		})

		It("should return true when both file exists and env var is true", func() {
			// Create debug file
			err := os.WriteFile(config.DEBUG_FILE, []byte{}, 0600)
			Expect(err).NotTo(HaveOccurred())

			// Set the config variable
			originalValue := config.KOHAN_DEBUG
			config.KOHAN_DEBUG = true
			defer func() { config.KOHAN_DEBUG = originalValue }()

			result := util.IsDebugMode()
			Expect(result).To(BeTrue())
		})
	})
})
