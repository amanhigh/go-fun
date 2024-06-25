package tui

import (
	"github.com/rivo/tview"
)

type Darius struct {
	tviewApp              *tview.Application
	availableServicesList *tview.List
	selectedServicesList  *tview.List
}

func NewApp() *Darius {
	services := []string{"MySQL", "Postgres", "Redis"}
	app := &Darius{
		tviewApp: tview.NewApplication(),
	}
	app.init(services)
	return app
}

func (a *Darius) init(services []string) {
	a.availableServicesList = a.createList("Available Services", services, a.updateContext)
	a.selectedServicesList = a.createList("Selected Services", nil, nil)
}

func (a *Darius) createList(title string, items []string, selectedFunc func(int, string, string, rune)) *tview.List {
	list := tview.NewList()
	list.SetBorder(true).SetTitle(title)

	for _, item := range items {
		list.AddItem(item, "", 0, nil)
	}

	if selectedFunc != nil {
		list.SetSelectedFunc(selectedFunc)
	}

	return list
}

func (a *Darius) updateContext(index int, main string, secondary string, shortcut rune) {
	selectedIndex := a.findItemIndex(a.selectedServicesList, main)
	if selectedIndex == -1 {
		a.selectedServicesList.AddItem(main, "", 0, nil)
	} else {
		a.selectedServicesList.RemoveItem(selectedIndex)
	}
}

func (a *Darius) findItemIndex(list *tview.List, text string) (index int) {
	indexs := list.FindItems(text, "", false, false)
	if len(indexs) > 0 {
		index = indexs[0]
	} else {
		index = -1
	}
	return
}

func (a *Darius) Run() error {
	flex := tview.NewFlex().
		AddItem(a.availableServicesList, 0, 1, true).
		AddItem(a.selectedServicesList, 0, 1, false)

	return a.tviewApp.SetRoot(flex, true).Run()
}
