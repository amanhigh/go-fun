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
}

func NewUIManager(app *tview.Application, svcManager *ServiceManager) *UIManager {
	return &UIManager{
		app:         app,
		mainFlex:    tview.NewFlex(),
		contextView: createTextView("Context"),
		commandView: createTextView("Command"),
		svcManager:  svcManager,
		svcList:     createList("Services", svcManager.GetAllServices()),
	}
}

func (ui *UIManager) SetupLayout() {
	leftPane := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(ui.svcList, 0, 1, true)
	rightPane := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(ui.contextView, 0, 1, false).
		AddItem(ui.commandView, 0, 1, false)
	ui.mainFlex.AddItem(leftPane, 0, 1, true)
	ui.mainFlex.AddItem(rightPane, 0, 1, false)
	ui.mainFlex.SetTitle("Helm Manager").SetBorder(true)
	ui.app.SetRoot(ui.mainFlex, true)
	ui.updateContext()
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

func (ui *UIManager) updateContext() {
	helpText := "Help:\n" +
		"- Use Arrow keys to navigate\n" +
		"- Enter to select\n" +
		"- Esc to exit\n\n"
	selectedServices := "Selected Services:\n"
	for _, service := range ui.svcManager.GetSelectedServices() {
		selectedServices += fmt.Sprintf("- %s\n", service)
	}
	if len(ui.svcManager.GetSelectedServices()) == 0 {
		selectedServices += "(None)"
	}
	ui.contextView.SetText(fmt.Sprintf("%s%s", helpText, selectedServices))
}
