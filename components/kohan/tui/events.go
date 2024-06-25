package tui

import (
	"fmt"
	"strings"
)

func (a *Darius) toggleServiceSelection(index int, main string, secondary string, shortcut rune) {
	item, _ := a.availableServicesList.GetItemText(index)
	if strings.HasPrefix(item, "[x]") {
		a.availableServicesList.SetItemText(index, strings.TrimPrefix(item, "[x] "), "")
	} else {
		a.availableServicesList.SetItemText(index, "[x] "+item, "")
	}
	a.updateContext()
}

func (a *Darius) updateContext() {
	a.contextView.Clear()
	fmt.Fprintf(a.contextView, "Selected Services:\n")
	for i := 0; i < a.availableServicesList.GetItemCount(); i++ {
		item, _ := a.availableServicesList.GetItemText(i)
		if strings.HasPrefix(item, "[x]") {
			fmt.Fprintf(a.contextView, "- %s\n", strings.TrimPrefix(item, "[x] "))
		}
	}
	fmt.Fprintf(a.contextView, "\nOperations:\n")
	fmt.Fprintf(a.contextView, "I - Install (make setup)\n")
	fmt.Fprintf(a.contextView, "C - Clean (make clean)\n")
}

func (a *Darius) filterAvailableServices(filter string) {
	a.availableServicesList.Clear()
	for _, service := range a.allServices {
		if strings.Contains(strings.ToLower(service), strings.ToLower(filter)) {
			a.availableServicesList.AddItem(service, "", 0, nil)
		}
	}
}
