package tui

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rivo/tview"
)

var _ = Describe("HotkeyManager", func() {
	var (
		hotkeyMgr *HotkeyManagerImpl
		app       *tview.Application
		uiMgr     *UIManagerImpl
		svcMgr    *ServiceManagerImpl
	)

	BeforeEach(func() {
		app = tview.NewApplication()
		uiMgr = &UIManagerImpl{app: app}
		svcMgr = &ServiceManagerImpl{}
		hotkeyMgr = &HotkeyManagerImpl{
			app:            app,
			uiManager:      uiMgr,
			serviceManager: svcMgr,
		}
	})
	It("should build", func() {
		Expect(hotkeyMgr).NotTo(BeNil())
		Expect(hotkeyMgr.app).To(Equal(app))
		Expect(hotkeyMgr.uiManager).To(Equal(uiMgr))
		Expect(hotkeyMgr.serviceManager).To(Equal(svcMgr))
	})

	Describe("Generate Help", func() {
		BeforeEach(func() {
			hotkeyMgr.SetupHotkeys()
		})

		It("should set up hotkeys", func() {
			Expect(hotkeyMgr.hotkeys).NotTo(BeEmpty())
			Expect(len(hotkeyMgr.hotkeys)).To(BeNumerically(">=", 8)) // Assuming at least 8 hotkeys as per the original code
		})

		It("should include all hotkeys in help text", func() {
			helpText := hotkeyMgr.GenerateHelpText()
			for _, hotkey := range hotkeyMgr.hotkeys {
				Expect(helpText).To(ContainSubstring(string(hotkey.Key)))
				Expect(helpText).To(ContainSubstring(hotkey.Description))
			}
		})
	})
})
