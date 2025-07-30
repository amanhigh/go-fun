package repository_test

import (
	"os"
	"path/filepath"

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
		tempDir, err = os.MkdirTemp("", "tui-repo-test-*")
		Expect(err).NotTo(HaveOccurred())

		// Create a dummy services directory for makeDir
		dummyMakeDir = filepath.Join(tempDir, "makefiles")
		util.RecreateDir(dummyMakeDir)
		util.RecreateDir(filepath.Join(dummyMakeDir, "services"))

		selectedServicesFile = filepath.Join(tempDir, "selected_services.txt")
		repo = repository.NewTuiServiceRepository(selectedServicesFile)
	})

	AfterEach(func() {
		err = os.RemoveAll(tempDir)
		Expect(err).NotTo(HaveOccurred())
	})

	Context("LoadSelectedServices", func() {
		Context("Valid Content", func() {
			BeforeEach(func() {
				servicesToTest = []string{"service1", "service2", "serviceAlpha"}
				err = repo.SaveSelectedServices(servicesToTest)
				Expect(err).NotTo(HaveOccurred())
			})

			It("should load the services correctly", func() {
				loadedServices, loadErr := repo.LoadSelectedServices()
				Expect(loadErr).NotTo(HaveOccurred())
				Expect(loadedServices).To(Equal(servicesToTest))
			})
		})

		Context("Empty Content", func() {
			BeforeEach(func() {
				servicesToTest = []string{}
				err = repo.SaveSelectedServices(servicesToTest)
				Expect(err).NotTo(HaveOccurred())
			})

			It("should return an empty slice", func() {
				loadedServices, loadErr := repo.LoadSelectedServices()
				Expect(loadErr).NotTo(HaveOccurred())
				Expect(loadedServices).To(BeEmpty())
			})
		})
	})

	Describe("SaveSelectedServices", func() {
		Context("Overwriting", func() {
			BeforeEach(func() {
				initialServices := []string{"old_service_one", "old_service_two"}
				err = repo.SaveSelectedServices(initialServices)
				Expect(err).NotTo(HaveOccurred())
			})

			It("should replace the content entirely with the new services", func() {
				newServices := []string{"new_service_A", "new_service_B"}
				err = repo.SaveSelectedServices(newServices)
				Expect(err).NotTo(HaveOccurred())

				loadedServices, loadErr := repo.LoadSelectedServices()
				Expect(loadErr).NotTo(HaveOccurred())
				Expect(loadedServices).To(Equal(newServices))
			})
		})
	})
})
