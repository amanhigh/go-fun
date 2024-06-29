package tui

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

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
