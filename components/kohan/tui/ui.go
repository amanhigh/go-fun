package tui

import (
	"github.com/rivo/tview"
)

type UIManager struct {
	app         *tview.Application
	mainFlex    *tview.Flex
	svcList     *tview.List
	contextView *tview.TextView
	commandView *tview.TextView
}

func NewUIManager(app *tview.Application, services []string) *UIManager {
	return &UIManager{
		app:         app,
		mainFlex:    tview.NewFlex(),
		contextView: createTextView("Context"),
		commandView: createTextView("Command"),
		svcList:     createList("Services", services),
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
