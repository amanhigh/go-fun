package tui

import (
	"github.com/amanhigh/go-fun/models/config"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ServiceManager", func() {
	var (
		svcMgr     *ServiceManager
		testConfig = config.DariusConfig{
			MakeDir:             "/make",
			SelectedServiceFile: "/tmp/selected.txt",
		}
	)

	BeforeEach(func() {
		svcMgr = newServiceManager(testConfig)
	})

	It("should build", func() {
		Expect(svcMgr).To(Not(BeNil()))
	})
})
