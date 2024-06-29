package tui

import (
	"fmt"
	"strings"

	"github.com/rivo/tview"
)

func (d *Darius) filterServices(filter string) {
	d.availableServices.Clear()
	for _, service := range d.allServices {
		if strings.Contains(strings.ToLower(service), strings.ToLower(filter)) {
			d.availableServices.AddItem(service, "", 0, nil)
		}
	}
}

func (d *Darius) showMinikubePanel() {
	if !d.isMinikubeVisible {
		d.isMinikubeVisible = true
		d.mainFlex.GetItem(0).(*tview.Flex).AddItem(d.minikubeModal, 0, 1, true)
	}
	d.app.SetFocus(d.minikubeModal)
}

func (d *Darius) hideMinikubePanel() {
	if d.isMinikubeVisible {
		d.isMinikubeVisible = false
		d.mainFlex.GetItem(0).(*tview.Flex).RemoveItem(d.minikubeModal)
	}
	d.app.SetFocus(d.availableServices)
}

func (d *Darius) createServiceList(title string, items []string, selectedFunc func(int, string, string, rune)) *tview.List {
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

func (d *Darius) createTextView(title string, dynamicColors bool) *tview.TextView {
	textView := tview.NewTextView()
	textView.SetDynamicColors(dynamicColors).SetRegions(true).SetWrap(true).SetTitle(title).SetBorder(true)
	return textView
}

func (d *Darius) createFilterInput() *tview.InputField {
	return tview.NewInputField().
		SetLabel(FilterLabel).
		SetChangedFunc(func(text string) {
			d.filterServices(text)
		})
}

func (d *Darius) setupLayout() {
	leftPane := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(d.filterInput, 1, 0, false).
		AddItem(d.availableServices, 0, 1, true)

	rightPane := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(d.contextView, 0, 2, false).
		AddItem(d.commandView, 0, 1, false).
		AddItem(d.runOutputView, 0, 4, false)

	d.mainFlex = tview.NewFlex().
		AddItem(leftPane, 0, 1, true).
		AddItem(rightPane, 0, 1, false)

	d.mainFlex.SetTitle("Helm Manager").SetBorder(true)

	d.minikubeModal = tview.NewModal().
		SetText("Minikube Operations").
		AddButtons([]string{"Start", "Stop", "Reset", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			switch buttonLabel {
			case "Stop":
				d.hideMinikubePanel()
				d.executeCommand(MinikubeStop)
			case "Reset":
				d.hideMinikubePanel()
				d.executeCommand(MinikubeReset)
			case "Cancel":
				d.hideMinikubePanel()
			}
		})

	d.app.SetRoot(d.mainFlex, true)
}

func (d *Darius) updateContext() {
	d.contextView.Clear()
	fmt.Fprintln(d.contextView, "Selected Services:")
	for service := range d.selectedServices {
		fmt.Fprintf(d.contextView, "- %s\n", service)
	}
	fmt.Fprintln(d.contextView, "\nOperations:")
	fmt.Fprintln(d.contextView, "I - Install (make setup)")
	fmt.Fprintln(d.contextView, "C - Clean (make clean)")
}
