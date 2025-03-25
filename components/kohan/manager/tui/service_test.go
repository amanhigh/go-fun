package tui

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ServiceManager", func() {
	var (
		makeDir = "/make"
		svcPath = "/selected.txt"
		svcMgr  *ServiceManagerImpl
	)

	BeforeEach(func() {
		svcMgr = NewServiceManager(makeDir, svcPath)
	})

	It("should build", func() {
		Expect(svcMgr).To(Not(BeNil()))
		Expect(svcMgr.GetAllServices()).Should(ContainElement("dummy"))
		Expect(svcMgr.GetSelectedServices()).To(BeEmpty())
	})

	Context("Services Basics", func() {
		var (
			dummyAvailableServices = []string{"dummy", "dummy2"}
			dummySelectedServices  = []string{"dummy"}
		)

		BeforeEach(func() {
			svcMgr.allServices = dummyAvailableServices
			svcMgr.selectedServices = dummySelectedServices
		})

		It("should load available services", func() {
			Expect(svcMgr.GetAllServices()).To(Equal(dummyAvailableServices))
		})

		It("should load selected services", func() {
			Expect(svcMgr.GetSelectedServices()).To(Equal(dummySelectedServices))
		})
	})

	Context("Service Selection", func() {
		BeforeEach(func() {
			svcMgr.allServices = []string{"service1", "service2", "service3"}
			svcMgr.selectedServices = []string{"service1"}
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
			svcMgr.ClearSelectedServices()
			Expect(svcMgr.GetSelectedServices()).To(BeEmpty())
		})
	})

	Context("Service Filtering", func() {
		BeforeEach(func() {
			svcMgr.allServices = []string{"service1", "service2", "microservice3"}
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
			svcMgr.ToggleFilteredServices()
			Expect(svcMgr.GetSelectedServices()).To(ConsistOf("service1", "service2", "microservice3"))

			svcMgr.ClearSelectedServices()
			svcMgr.FilterServices("micro")
			svcMgr.ToggleFilteredServices()
			Expect(svcMgr.GetSelectedServices()).To(ConsistOf("microservice3"))
		})
	})

})
