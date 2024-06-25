package tui

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	CommandSetup  = "setup"
	CommandClean  = "clean"
	FilterLabel   = "Filter: "
	CommandPrefix = "Running: make "
)

type Darius struct {
	app               *tview.Application
	availableServices *tview.List
	contextView       *tview.TextView
	commandView       *tview.TextView
	runOutputView     *tview.TextView
	filterInput       *tview.InputField
	allServices       []string
	selectedServices  map[string]bool
}

func NewApp() *Darius {
	services := []string{"MySQL", "Postgres", "Redis", "Mongo"}
	darius := &Darius{
		app:              tview.NewApplication(),
		allServices:      services,
		selectedServices: make(map[string]bool),
	}
	darius.init()
	return darius
}

func (d *Darius) init() {
	d.availableServices = d.createServiceList("Available Services", d.allServices, d.toggleServiceSelection)
	d.contextView = d.createTextView("Context", true)
	d.commandView = d.createTextView("Command", true)
	d.commandView.SetText("Select a service & perform operation to see the command to run.")
	d.runOutputView = d.createTextView("Run", true)
	d.filterInput = d.createFilterInput()

	d.setupLayout()
	d.setupKeyBindings()
	d.updateContext()
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

	mainLayout := tview.NewFlex().
		AddItem(leftPane, 0, 1, true).
		AddItem(rightPane, 0, 1, false)

	mainLayout.SetTitle("Helm Manager").SetBorder(true)
	d.app.SetRoot(mainLayout, true)
}

func (d *Darius) setupKeyBindings() {
	d.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEnter:
			d.handleEnterKey()
		case tcell.KeyEscape:
			d.handleEscapeKey()
		case tcell.KeyRune:
			d.handleRuneKey(event.Rune())
		}
		return event
	})
}

func (d *Darius) handleEnterKey() {
	if d.app.GetFocus() == d.availableServices {
		d.toggleServiceSelection(d.availableServices.GetCurrentItem(), "", "", 0)
	} else if d.app.GetFocus() == d.filterInput {
		d.applyFilterToSelection()
	}
}

func (d *Darius) handleEscapeKey() {
	if d.app.GetFocus() == d.filterInput {
		d.filterInput.SetText("")
		d.filterServices("")
		d.app.SetFocus(d.availableServices)
	}
}

func (d *Darius) handleRuneKey(r rune) {
	switch r {
	case 'i', 'I':
		d.executeCommand(CommandSetup)
	case 'c', 'C':
		d.executeCommand(CommandClean)
	case '/':
		d.app.SetFocus(d.filterInput)
	case 'q', 'Q':
		d.app.Stop()
	}
}

func (d *Darius) applyFilterToSelection() {
	filterText := strings.ToLower(d.filterInput.GetText())
	for i := 0; i < d.availableServices.GetItemCount(); i++ {
		service, _ := d.availableServices.GetItemText(i)
		if strings.Contains(strings.ToLower(service), filterText) {
			d.toggleServiceSelection(i, "", "", 0)
		}
	}
}

func (d *Darius) toggleServiceSelection(index int, main, secondary string, shortcut rune) {
	service, _ := d.availableServices.GetItemText(index)
	if d.selectedServices[service] {
		delete(d.selectedServices, service)
	} else {
		d.selectedServices[service] = true
	}
	d.updateContext()
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

func (d *Darius) filterServices(filter string) {
	d.availableServices.Clear()
	for _, service := range d.allServices {
		if strings.Contains(strings.ToLower(service), strings.ToLower(filter)) {
			d.availableServices.AddItem(service, "", 0, nil)
		}
	}
}

func (d *Darius) executeCommand(command string) {
	d.commandView.SetText(CommandPrefix + command)
	cmd := exec.Command("make", command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		d.runOutputView.SetText(fmt.Sprintf("Error: %s\n%s", err, string(output)))
	} else {
		d.runOutputView.SetText(string(output))
	}
}

func (d *Darius) Run() error {
	return d.app.Run()
}
