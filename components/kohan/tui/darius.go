package tui

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

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

func (d *Darius) fetchAndFilterServices() ([]string, error) {
	cmd := exec.Command("bash", "-c", ServiceFetchCommand)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	// Remove ANSI escape codes
	ansiEscape := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	cleanedOutput := ansiEscape.ReplaceAllString(string(output), "")

	// Process the output to remove "help" and "make"
	lines := strings.Split(cleanedOutput, "\n")
	var services []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "help") || strings.HasPrefix(line, "make") {
			continue
		}
		// Extract the service name (the first word before the spaces)
		service := strings.Fields(line)[0]
		services = append(services, service)
	}
	return services, nil
}

func (d *Darius) saveSelectedServices() error {
	file, err := os.Create(ServiceFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for service := range d.selectedServices {
		_, err := writer.WriteString(service + "\n")
		if err != nil {
			return err
		}
	}
	return writer.Flush()
}

func (d *Darius) loadSelectedServices() error {
	file, err := os.Open(ServiceFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // No file to load
		}
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		service := scanner.Text()
		d.selectedServices[service] = true
	}
	return scanner.Err()
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

	// Save the selected services to file
	if err := d.saveSelectedServices(); err != nil {
		d.runOutputView.SetText(fmt.Sprintf("Error saving selected services: %s", err))
	}
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

func (d *Darius) executeCommand(command string) {
	cmdString := fmt.Sprintf("make -C /home/aman/Projects/go-fun/Kubernetes %s", command)
	d.commandView.SetText(cmdString)

	// Split the command and its arguments
	cmd := exec.Command("make", "-C", "/home/aman/Projects/go-fun/Kubernetes/services", command)
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
