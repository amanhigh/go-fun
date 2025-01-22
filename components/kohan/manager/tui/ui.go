package tui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type UIManager interface {
	SetupLayout()
	ToggleServiceSelection()
	UpdateServicesList(filter string)
	ToggleFilteredServices()
	UpdateContext()
	ShowOutput(output string)
	ShowError(err error)
	clearFilterInput()

	// Focus Management Methods
	FocusServiceList()     // Focus the service list component
	FocusFilterInput()     // Focus the filter input component
	IsFocusOnFilter() bool // Check if filter input has focus
	IsFocusOnList() bool   // Check if service list has focus

	// New methods for hotkey management
	SetGlobalInputCapture(capture func(*tcell.EventKey) *tcell.EventKey)
	StopApplication()
}

type UIManagerImpl struct {
	app        *tview.Application
	svcManager ServiceManager

	// UI Elements
	mainFlex    *tview.Flex
	svcList     *tview.List
	contextView *tview.TextView
	commandView *tview.TextView
	filterInput *tview.InputField
}

func NewUIManager(app *tview.Application, svcManager ServiceManager) *UIManagerImpl {
	ui := &UIManagerImpl{
		app:         app,
		svcManager:  svcManager,
		mainFlex:    tview.NewFlex(),
		contextView: createTextView("Context"),
		commandView: createTextView("Command"),
		filterInput: tview.NewInputField(),
		svcList:     createList("Services", svcManager.GetAllServices()),
	}
	return ui
}

func (ui *UIManagerImpl) SetupLayout() {
	leftPane := ui.createLeftPane()
	rightPane := ui.createRightPane()
	ui.mainFlex.AddItem(leftPane, 0, 1, true)
	ui.mainFlex.AddItem(rightPane, 0, 1, false)
	ui.mainFlex.SetTitle("Helm Manager").SetBorder(true).SetBorderColor(tcell.ColorBlue)
	ui.app.SetRoot(ui.mainFlex, true)
	ui.UpdateContext()
	ui.setupFilterInput()
	ui.setupCustomKeys()
}

func (ui *UIManagerImpl) createLeftPane() *tview.Flex {
	leftPane := tview.NewFlex().SetDirection(tview.FlexRow)
	leftPane.AddItem(ui.svcList, 0, 1, true)
	leftPane.AddItem(ui.filterInput, 1, 0, false)
	return leftPane
}

func (ui *UIManagerImpl) createRightPane() *tview.Flex {
	rightPane := tview.NewFlex().SetDirection(tview.FlexRow)
	rightPane.AddItem(ui.contextView, 0, 1, false)
	rightPane.AddItem(ui.commandView, 0, 1, false)
	return rightPane
}

func createList(title string, items []string) *tview.List {
	list := tview.NewList()
	list.SetBorder(true).SetTitle(title).SetBorderColor(tcell.ColorGreen)
	for _, item := range items {
		list.AddItem(item, "", 0, nil)
	}
	return list
}

func createTextView(title string) *tview.TextView {
	tv := tview.NewTextView()
	tv.SetBorder(true).SetTitle(title).SetBorderColor(tcell.ColorYellow)

	return tv
}

func (ui *UIManagerImpl) ToggleServiceSelection() {
	if ui.svcList.GetItemCount() > 0 {
		index := ui.svcList.GetCurrentItem()
		name, _ := ui.svcList.GetItemText(index)
		ui.svcManager.ToggleServiceSelection(name)
		ui.UpdateContext()
	}
}

func (ui *UIManagerImpl) UpdateContext() {
	selectedServices := "Selected Services:\n"
	for _, service := range ui.svcManager.GetSelectedServices() {
		selectedServices += fmt.Sprintf("- %s\n", service)
	}
	if len(ui.svcManager.GetSelectedServices()) == 0 {
		selectedServices += "(None)"
	}
	ui.contextView.SetText(selectedServices)
}

func (ui *UIManagerImpl) ShowError(err error) {
	ui.ShowOutput(err.Error())
}

func (ui *UIManagerImpl) ShowOutput(output string) {
	ui.commandView.SetText(output)
}

func (ui *UIManagerImpl) UpdateServicesList(filter string) {
	ui.svcManager.FilterServices(filter)
	services := ui.svcManager.GetFilteredServices()
	ui.svcList.Clear()
	for _, service := range services {
		ui.svcList.AddItem(service, "", 0, nil)
	}
}

func (ui *UIManagerImpl) ToggleFilteredServices() {
	ui.svcManager.ToggleFilteredServices()
	ui.UpdateContext()
}

func (ui *UIManagerImpl) setupFilterInput() {
	ui.filterInput.SetLabel("Filter: ")
	ui.filterInput.SetChangedFunc(func(text string) {
		ui.UpdateServicesList(text)
	})
}
func (ui *UIManagerImpl) clearFilterInput() {
	ui.filterInput.SetText("")
	ui.UpdateServicesList("")
}

func (ui *UIManagerImpl) setupCustomKeys() {
	ui.svcList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() { //nolint:exhaustive
		case tcell.KeyRune:
			switch event.Rune() {
			case 'h':
				return tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModNone)
			case 't':
				return tcell.NewEventKey(tcell.KeyUp, 0, tcell.ModNone)
			}
		}
		return event
	})
}

func (ui *UIManagerImpl) FocusServiceList() {
	ui.app.SetFocus(ui.svcList)
}

func (ui *UIManagerImpl) FocusFilterInput() {
	ui.app.SetFocus(ui.filterInput)
}

func (ui *UIManagerImpl) IsFocusOnFilter() bool {
	return ui.app.GetFocus() == ui.filterInput
}

func (ui *UIManagerImpl) IsFocusOnList() bool {
	return ui.app.GetFocus() == ui.svcList
}

func (ui *UIManagerImpl) SetGlobalInputCapture(capture func(*tcell.EventKey) *tcell.EventKey) {
	ui.app.SetInputCapture(capture)
}

func (ui *UIManagerImpl) StopApplication() {
	ui.app.Stop()
}
