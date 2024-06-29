package tui

import (
	"github.com/rivo/tview"
)

const (
	FilterLabel         = "Filter: "
	CommandPrefix       = "Running: make "
	ServiceFetchCommand = "make -C /home/aman/Projects/go-fun/Kubernetes/services/ -f ./services.mk"
	CommandSetup        = "setup"
	CommandClean        = "clean"
	ServiceFilePath     = "/tmp/k8-svc.txt"
	MinikubeStop        = "clean"
	MinikubeReset       = "reset"
)

/*
	TODO: Darius Improvements
	- Vim Like Keys
	- Config Files
	- Clean Selected Services
	- Minikube Control
	- Funapp Verification and Load Test.
	- New Tabs
*/

type Darius struct {
	app               *tview.Application
	availableServices *tview.List
	contextView       *tview.TextView
	commandView       *tview.TextView
	runOutputView     *tview.TextView
	filterInput       *tview.InputField
	allServices       []string
	selectedServices  map[string]bool
	mainFlex          *tview.Flex
	minikubeModal     *tview.Modal
	isMinikubeVisible bool
}

func NewApp() *Darius {
	darius := &Darius{
		app:              tview.NewApplication(),
		selectedServices: make(map[string]bool),
	}
	if err := darius.init(); err != nil {
		panic(err)
	}
	return darius
}

func (d *Darius) init() error {
	services, err := d.fetchAndFilterServices()
	if err != nil {
		return err
	}
	d.allServices = services

	d.availableServices = d.createServiceList("Available Services", d.allServices, d.toggleServiceSelection)
	d.contextView = d.createTextView("Context", true)
	d.commandView = d.createTextView("Command", true)
	d.commandView.SetText("Select a service & perform operation to see the command to run.")
	d.runOutputView = d.createTextView("Run", true)
	d.filterInput = d.createFilterInput()

	d.setupLayout()
	d.setupKeyBindings()

	if err := d.loadSelectedServices(); err != nil {
		return err
	}
	d.updateContext()

	return nil
}

func (d *Darius) Run() error {
	return d.app.Run()
}
