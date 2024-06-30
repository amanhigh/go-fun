package tui

import (
	"fmt"

	"github.com/rivo/tview"
)

type UIManager struct {
	app        *tview.Application
	svcManager *ServiceManager

	// UI Elements
	mainFlex    *tview.Flex
	svcList     *tview.List
	contextView *tview.TextView
	commandView *tview.TextView
	filterInput *tview.InputField
}

func (ui *UIManager) SetupLayout() {
	leftPane := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(ui.svcList, 0, 1, true).
		AddItem(ui.filterInput, 1, 0, false)
	rightPane := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(ui.contextView, 0, 1, false).
		AddItem(ui.commandView, 0, 1, false)
	ui.mainFlex.AddItem(leftPane, 0, 1, true)
	ui.mainFlex.AddItem(rightPane, 0, 1, false)
	ui.mainFlex.SetTitle("Helm Manager").SetBorder(true)
	ui.app.SetRoot(ui.mainFlex, true)
	ui.UpdateContext()
	ui.setupFilterInput()
}

func createList(title string, items []string) *tview.List {
	list := tview.NewList()
	list.SetBorder(true).SetTitle(title)
	for _, item := range items {
		list.AddItem(item, "", 0, nil)
	}
	return list
}

func createTextView(title string) *tview.TextView {
	tv := tview.NewTextView()
	tv.SetBorder(true).SetTitle(title)
	return tv
}

func (ui *UIManager) ToggleServiceSelection() {
	index := ui.svcList.GetCurrentItem()
	mainText, _ := ui.svcList.GetItemText(index)
	ui.svcManager.ToggleServiceSelection(mainText)
}

func (ui *UIManager) UpdateContext() {
	selectedServices := "Selected Services:\n"
	for _, service := range ui.svcManager.GetSelectedServices() {
		selectedServices += fmt.Sprintf("- %s\n", service)
	}
	if len(ui.svcManager.GetSelectedServices()) == 0 {
		selectedServices += "(None)"
	}
	ui.contextView.SetText(selectedServices)
}

func (ui *UIManager) ShowHelp(helpText string) {
	ui.contextView.SetText(helpText)
}

func (ui *UIManager) ShowError(err error) {
	ui.commandView.SetText(err.Error())
}

func (ui *UIManager) UpdateServicesList(filter string) {
	ui.svcManager.FilterServices(filter)
	services := ui.svcManager.GetFilteredServices()
	ui.svcList.Clear()
	for _, service := range services {
		ui.svcList.AddItem(service, "", 0, nil)
	}
}

func (ui *UIManager) ToggleFilteredServices() {
	ui.svcManager.ToggleFilteredServices()
	ui.UpdateContext()
}

func (ui *UIManager) setupFilterInput() {
	ui.filterInput.SetLabel("Filter: ")
	ui.filterInput.SetChangedFunc(func(text string) {
		ui.UpdateServicesList(text)
	})
}

func (ui *UIManager) clearFilterInput() {
	ui.filterInput.SetText("")
	ui.UpdateServicesList("")
}
