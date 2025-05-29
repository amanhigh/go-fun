package tui

import (
	"fmt"

	"github.com/amanhigh/go-fun/components/kohan/repository/mocks"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
)

var _ = Describe("ServiceManager", func() {
	var (
		makeDir  = "/make"
		svcMgr   *ServiceManagerImpl
		mockRepo *mocks.TuiServiceRepository
	)

	BeforeEach(func() {
		mockRepo = new(mocks.TuiServiceRepository)
		// Default mocks for initial load in NewServiceManager
		mockRepo.On("LoadAvailableServices", makeDir).Return([]string{"fallback"}, nil).Maybe()
		mockRepo.On("LoadSelectedServices").Return([]string{}, nil).Maybe()
		svcMgr = NewServiceManager(makeDir, mockRepo)
	})

	It("should build", func() {
		Expect(svcMgr).To(Not(BeNil()))
		Expect(svcMgr.GetAllServices()).To(Equal([]string{"fallback"}))
		Expect(svcMgr.GetSelectedServices()).To(BeEmpty())
	})

	Context("Services Basics", func() {
		var (
			dummyAvailableServices = []string{"dummy", "dummy2"}
			dummySelectedServices  = []string{"dummy"}
		)

		BeforeEach(func() {
			mockRepo = new(mocks.TuiServiceRepository) // New mock for specific context
			mockRepo.On("LoadAvailableServices", makeDir).Return(dummyAvailableServices, nil)
			mockRepo.On("LoadSelectedServices").Return(dummySelectedServices, nil)
			svcMgr = NewServiceManager(makeDir, mockRepo)
		})

		It("should load available services", func() {
			Expect(svcMgr.GetAllServices()).To(Equal(dummyAvailableServices))
		})

		It("should load selected services", func() {
			Expect(svcMgr.GetSelectedServices()).To(Equal(dummySelectedServices))
		})
	})

	Context("Service Selection", func() {
		var initialSelectedServices = []string{"service1"}
		var allServices = []string{"service1", "service2", "service3"}
		BeforeEach(func() {
			mockRepo = new(mocks.TuiServiceRepository) // New mock
			mockRepo.On("LoadAvailableServices", makeDir).Return(allServices, nil)
			mockRepo.On("LoadSelectedServices").Return(initialSelectedServices, nil)
			// Allow SaveSelectedServices to be called multiple times
			mockRepo.On("SaveSelectedServices", mock.Anything).Return(nil).Maybe()
			svcMgr = NewServiceManager(makeDir, mockRepo)
		})

		It("should toggle service selection", func() {
			Expect(svcMgr.IsServiceSelected("service1")).To(BeTrue())
			Expect(svcMgr.IsServiceSelected("service2")).To(BeFalse())

			svcMgr.ToggleServiceSelection("service2")
			Expect(svcMgr.IsServiceSelected("service2")).To(BeTrue())
			Expect(svcMgr.GetSelectedServices()).To(HaveLen(2))

			svcMgr.ToggleServiceSelection("service1")
			Expect(svcMgr.IsServiceSelected("service1")).To(BeFalse())
			Expect(svcMgr.GetSelectedServices()).To(HaveLen(1))
		})

		It("should clear selected services", func() {
			mockRepo.On("SaveSelectedServices", []string{}).Return(nil).Once()
			svcMgr.ClearSelectedServices()
			Expect(svcMgr.GetSelectedServices()).To(BeEmpty())
		})
	})

	Context("Service Filtering", func() {
		var allServices = []string{"service1", "service2", "microservice3"}
		BeforeEach(func() {
			mockRepo = new(mocks.TuiServiceRepository) // New mock
			mockRepo.On("LoadAvailableServices", makeDir).Return(allServices, nil)
			mockRepo.On("LoadSelectedServices").Return([]string{}, nil) // Start with no selected services for filter toggle test
			svcMgr = NewServiceManager(makeDir, mockRepo)
		})

		It("should filter services based on keyword", func() {
			svcMgr.FilterServices("service")
			Expect(svcMgr.GetFilteredServices()).To(ConsistOf("service1", "service2", "microservice3"))

			svcMgr.FilterServices("micro")
			Expect(svcMgr.GetFilteredServices()).To(ConsistOf("microservice3"))

			svcMgr.FilterServices("nonexistent")
			Expect(svcMgr.GetFilteredServices()).To(BeEmpty())
		})

		It("should toggle filtered services", func() {
			svcMgr.FilterServices("service")
			// SaveSelectedServices will be called for each toggled service
			// For "service1", "service2", "microservice3"
			mockRepo.On("SaveSelectedServices", []string{"service1"}).Return(nil).Once()
			mockRepo.On("SaveSelectedServices", []string{"service1", "service2"}).Return(nil).Once()
			mockRepo.On("SaveSelectedServices", []string{"service1", "service2", "microservice3"}).Return(nil).Once()

			svcMgr.ToggleFilteredServices()
			Expect(svcMgr.GetSelectedServices()).To(ConsistOf("service1", "service2", "microservice3"))

			// Reset selected services and mock for the next part of the test
			mockRepo.On("SaveSelectedServices", []string{}).Return(nil).Once() // For ClearSelectedServices
			svcMgr.ClearSelectedServices()

			mockRepo.On("SaveSelectedServices", []string{"microservice3"}).Return(nil).Once()
			svcMgr.FilterServices("micro")
			svcMgr.ToggleFilteredServices()
			Expect(svcMgr.GetSelectedServices()).To(ConsistOf("microservice3"))

			mockRepo.AssertExpectations(GinkgoT())
		})
	})

	Context("Service Operations", func() {
		var selectedServices = []string{"service1"}
		var serviceMakeDir = makeDir + "/services" // As per getServiceMakeDir in manager

		BeforeEach(func() {
			mockRepo = new(mocks.TuiServiceRepository)
			mockRepo.On("LoadAvailableServices", makeDir).Return([]string{"service1", "service2"}, nil)
			mockRepo.On("LoadSelectedServices").Return(selectedServices, nil)
			svcMgr = NewServiceManager(makeDir, mockRepo)
		})

		It("should execute CleanServices", func() {
			expectedOutput := []string{"clean output"}
			mockRepo.On("ExecuteMakeCommand", serviceMakeDir, "Makefile", "clean").Return(expectedOutput, nil).Once()
			output, err := svcMgr.CleanServices()
			Expect(err).NotTo(HaveOccurred())
			Expect(output).To(Equal("clean output"))
		})

		It("should handle error from CleanServices", func() {
			mockRepo.On("ExecuteMakeCommand", serviceMakeDir, "Makefile", "clean").Return(nil, fmt.Errorf("make error")).Once()
			_, err := svcMgr.CleanServices()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("make error"))
		})

		It("should execute SetupServices", func() {
			expectedOutput := []string{"setup output"}
			mockRepo.On("ExecuteMakeCommand", serviceMakeDir, "Makefile", "setup").Return(expectedOutput, nil).Once()
			output, err := svcMgr.SetupServices()
			Expect(err).NotTo(HaveOccurred())
			Expect(output).To(Equal("setup output"))
		})

		It("should execute UpdateServices", func() {
			expectedOutput := []string{"update output"}
			mockRepo.On("ExecuteMakeCommand", serviceMakeDir, "Makefile", "update").Return(expectedOutput, nil).Once()
			output, err := svcMgr.UpdateServices()
			Expect(err).NotTo(HaveOccurred())
			Expect(output).To(Equal("update output"))
		})
	})
})
