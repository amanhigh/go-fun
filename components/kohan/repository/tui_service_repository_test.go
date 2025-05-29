package repository_test

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/components/kohan/repository"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("TuiServiceRepository", func() {
	var (
		repo                 repository.TuiServiceRepository
		tempDir              string
		err                  error
		dummyMakeDir         string
		selectedServicesFile string
		servicesToTest       []string
	)

	BeforeEach(func() {
		repo = repository.NewTuiServiceRepository()
		tempDir, err = os.MkdirTemp("", "tui-repo-test-*")
		Expect(err).NotTo(HaveOccurred())

		// Create a dummy services directory for makeDir
		dummyMakeDir = filepath.Join(tempDir, "makefiles")
		util.RecreateDir(dummyMakeDir)
		util.RecreateDir(filepath.Join(dummyMakeDir, "services"))

		selectedServicesFile = filepath.Join(tempDir, "selected_services.txt")
	})

	AfterEach(func() {
		err = os.RemoveAll(tempDir)
		Expect(err).NotTo(HaveOccurred())
	})

	Context("LoadSelectedServices", func() {
		Context("Valid Content", func() {
			BeforeEach(func() {
				servicesToTest = []string{"service1", "service2", "serviceAlpha"}
				err = repo.SaveSelectedServices(selectedServicesFile, servicesToTest)
				Expect(err).NotTo(HaveOccurred())
			})

			It("should load the services correctly", func() {
				loadedServices, loadErr := repo.LoadSelectedServices(selectedServicesFile)
				Expect(loadErr).NotTo(HaveOccurred())
				Expect(loadedServices).To(Equal(servicesToTest))
			})
		})

		Context("Empty Content", func() {
			BeforeEach(func() {
				servicesToTest = []string{}
				err = repo.SaveSelectedServices(selectedServicesFile, servicesToTest)
				Expect(err).NotTo(HaveOccurred())
			})

			It("should return an empty slice", func() {
				loadedServices, loadErr := repo.LoadSelectedServices(selectedServicesFile)
				Expect(loadErr).NotTo(HaveOccurred())
				Expect(loadedServices).To(BeEmpty())
			})
		})

		Context("No file", func() {
			It("should return an empty slice and no error", func() {
				// Ensure the file truly doesn't exist for this specific test
				nonExistentFile := filepath.Join(tempDir, "definitely_not_there.txt")
				loadedServices, loadErr := repo.LoadSelectedServices(nonExistentFile)
				Expect(loadErr).NotTo(HaveOccurred())
				Expect(loadedServices).To(BeEmpty())
			})
		})
	})

	Describe("SaveSelectedServices", func() {
		Context("Overwriting", func() {
			BeforeEach(func() {
				initialServices := []string{"old_service_one", "old_service_two"}
				err = repo.SaveSelectedServices(selectedServicesFile, initialServices)
				Expect(err).NotTo(HaveOccurred())
			})

			It("should replace the content entirely with the new services", func() {
				newServices := []string{"new_service_A", "new_service_B"}
				err = repo.SaveSelectedServices(selectedServicesFile, newServices)
				Expect(err).NotTo(HaveOccurred())

				loadedServices, loadErr := repo.LoadSelectedServices(selectedServicesFile)
				Expect(loadErr).NotTo(HaveOccurred())
				Expect(loadedServices).To(Equal(newServices))
			})
		})

		Context("Invalid Path", func() {
			It("should return an error", func() {
				// tempDir is a directory, SaveSelectedServices should fail.
				err = repo.SaveSelectedServices(tempDir, []string{"service1"})
				Expect(err).To(HaveOccurred())
				// Check for a more specific part of the error message if possible,
				// but the wrapper in SaveSelectedServices adds "SaveSelectedServices: writing lines:"
				Expect(strings.Contains(err.Error(), "SaveSelectedServices: writing lines:")).To(BeTrue())
			})
		})
	})
})
