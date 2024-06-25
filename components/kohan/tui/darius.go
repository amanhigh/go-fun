package tui

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Darius struct {
	tviewApp              *tview.Application
	availableServicesList *tview.List
	contextView           *tview.TextView
	commandView           *tview.TextView
	runView               *tview.TextView
	filterInput           *tview.InputField
	allServices           []string
	selectedServices      map[string]bool // Map to track selected services
}

func NewApp() *Darius {
	services := []string{"MySQL", "Postgres", "Redis", "Mongo"}
	app := &Darius{
		tviewApp:         tview.NewApplication(),
		allServices:      services,
		selectedServices: make(map[string]bool), // Initialize the map
	}
	app.init()
	return app
}

func (a *Darius) init() {
	a.availableServicesList = a.createList("Available Services", a.allServices, a.toggleServiceSelection)

	a.contextView = tview.NewTextView()
	a.contextView.SetDynamicColors(true)
	a.contextView.SetRegions(true)
	a.contextView.SetWrap(true)
	a.contextView.SetTitle("Context")
	a.contextView.SetBorder(true)

	a.commandView = tview.NewTextView()
	a.commandView.SetDynamicColors(true)
	a.commandView.SetTitle("Command")
	a.commandView.SetBorder(true)
	a.commandView.SetText("Select a service & perform operation to see the command to run.")

	a.runView = tview.NewTextView()
	a.runView.SetDynamicColors(true)
	a.runView.SetRegions(true)
	a.runView.SetWrap(true)
	a.runView.SetTitle("Run")
	a.runView.SetBorder(true)

	a.filterInput = tview.NewInputField().
		SetLabel("Filter: ").
		SetChangedFunc(func(text string) {
			a.filterAvailableServices(text)
		})
}

func (a *Darius) createList(title string, items []string, selectedFunc func(int, string, string, rune)) *tview.List {
	list := tview.NewList()
	list.SetBorder(true)
	list.SetTitle(title)

	for _, item := range items {
		list.AddItem(item, "", 0, nil)
	}

	if selectedFunc != nil {
		list.SetSelectedFunc(selectedFunc)
	}

	return list
}

func (a *Darius) executeCommand(command string) {
	a.commandView.SetText("Running: make " + command)
	cmd := exec.Command("make", command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		a.runView.SetText("Error: " + err.Error() + "\n" + string(output))
	} else {
		a.runView.SetText(string(output))
	}
}

func (a *Darius) Run() error {
	// Create the layout
	leftSide := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(a.filterInput, 1, 0, false).
		AddItem(a.availableServicesList, 0, 1, true)

	rightSide := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(a.contextView, 0, 2, false).
		AddItem(a.commandView, 0, 1, false).
		AddItem(a.runView, 0, 4, false)

	mainFlex := tview.NewFlex().
		AddItem(leftSide, 0, 1, true).
		AddItem(rightSide, 0, 1, false)

	// Set up key bindings
	a.tviewApp.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEnter:
			if a.tviewApp.GetFocus() == a.availableServicesList {
				a.toggleServiceSelection(a.availableServicesList.GetCurrentItem(), "", "", 0)
				return nil
			} else if a.tviewApp.GetFocus() == a.filterInput {
				filterText := strings.ToLower(a.filterInput.GetText())
				for i := 0; i < a.availableServicesList.GetItemCount(); i++ {
					itemName, _ := a.availableServicesList.GetItemText(i)
					itemNameLower := strings.ToLower(itemName)
					if strings.Contains(itemNameLower, filterText) {
						a.toggleServiceSelection(i, "", "", 0)
					}
				}
				return nil
			}
		case tcell.KeyEscape:
			if a.tviewApp.GetFocus() == a.filterInput {
				a.filterInput.SetText("")
				a.filterAvailableServices("")
				a.tviewApp.SetFocus(a.availableServicesList)
				return nil
			}
		case tcell.KeyRune:
			switch event.Rune() {
			case 'i', 'I':
				a.executeCommand("setup")
			case 'c', 'C':
				a.executeCommand("clean")
			case '/':
				a.tviewApp.SetFocus(a.filterInput)
				return nil // Prevent '/' from being typed in the filter bar
			case 'q', 'Q':
				a.tviewApp.Stop()
			}
		}
		return event
	})

	// Set the title
	mainFlex.SetTitle("Helm Manager")
	mainFlex.SetBorder(true)

	// Update context initially
	a.updateContext()

	return a.tviewApp.SetRoot(mainFlex, true).Run()
}

func (a *Darius) toggleServiceSelection(index int, main string, secondary string, shortcut rune) {
	item, _ := a.availableServicesList.GetItemText(index)
	if a.selectedServices[item] {
		delete(a.selectedServices, item)
	} else {
		a.selectedServices[item] = true
	}
	a.updateContext()
}

func (a *Darius) updateContext() {
	a.contextView.Clear()
	fmt.Fprintf(a.contextView, "Selected Services:\n")
	for service := range a.selectedServices {
		fmt.Fprintf(a.contextView, "- %s\n", service)
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
