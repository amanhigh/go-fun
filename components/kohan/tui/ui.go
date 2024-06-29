package tui

import (
	"github.com/rivo/tview"
)

type UIManager struct {
	app      *tview.Application
	mainFlex *tview.Flex
	svcList  *tview.List
}

func NewUIManager(app *tview.Application, services []string) *UIManager {
	return &UIManager{
		app:      app,
		mainFlex: tview.NewFlex(),
		svcList:  createList("Services", services),
	}
}

func (ui *UIManager) SetupLayout() {
	ui.mainFlex.AddItem(ui.svcList, 0, 1, true)
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
