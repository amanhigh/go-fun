package tui

import (
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
}

func NewApp() *Darius {
	services := []string{"MySQL", "Postgres", "Redis", "Mongo"}
	app := &Darius{
		tviewApp:    tview.NewApplication(),
		allServices: services,
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

	a.commandView.SetText("Select a service & perform operation to see the command to run.")

	// Set up key bindings
	a.tviewApp.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEnter:
			if a.tviewApp.GetFocus() == a.availableServicesList {
				a.toggleServiceSelection(a.availableServicesList.GetCurrentItem(), "", "", 0)
				return nil
			} else if a.tviewApp.GetFocus() == a.filterInput {
				filterText := strings.ToLower(a.filterInput.GetText()) // Step 1: Get and lowercase the filter text
				for i := 0; i < a.availableServicesList.GetItemCount(); i++ {
					itemName, _ := a.availableServicesList.GetItemText(i)
					itemNameLower := strings.ToLower(itemName)       // Lowercase the item name for case-insensitive comparison
					if strings.Contains(itemNameLower, filterText) { // Step 3: Check if the item name contains the filter text
						a.toggleServiceSelection(i, "", "", 0) // Step 4: Toggle the selection
					}
				}
				return nil
			}
		case tcell.KeyEscape:
			if a.tviewApp.GetFocus() == a.filterInput {
				a.filterInput.SetText("")                    // Clear the filter input field
				a.filterAvailableServices("")                // Reset the available services list
				a.tviewApp.SetFocus(a.availableServicesList) // Set focus back to the available services list
				return nil                                   // Prevent further handling of the Escape key
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
