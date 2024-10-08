package tui

import (
	"os/exec"
	"regexp"
	"strings"

	"github.com/amanhigh/go-fun/common/util"
)

type ServiceManager struct {
	allServices         []string
	selectedServices    []string
	filteredServices    []string
	makeDir             string
	selectedServicePath string
}

func (sm *ServiceManager) GetAllServices() []string {
	return sm.allServices
}

func (sm *ServiceManager) GetSelectedServices() []string {
	return sm.selectedServices
}

func (sm *ServiceManager) IsServiceSelected(service string) bool {
	for _, s := range sm.selectedServices {
		if s == service {
			return true
		}
	}
	return false
}

func (sm *ServiceManager) ToggleServiceSelection(service string) {
	if sm.IsServiceSelected(service) {
		sm.removeService(service)
	} else {
		sm.selectedServices = append(sm.selectedServices, service)
	}
	sm.saveSelectedServices()
}

func (sm *ServiceManager) removeService(service string) {
	for i, s := range sm.selectedServices {
		if s == service {
			sm.selectedServices = append(sm.selectedServices[:i], sm.selectedServices[i+1:]...)
			break
		}
	}
}

func (sm *ServiceManager) FilterServices(keyword string) {
	sm.filteredServices = []string{}
	lowerKeyword := strings.ToLower(keyword)
	for _, service := range sm.allServices {
		if strings.Contains(strings.ToLower(service), lowerKeyword) {
			sm.filteredServices = append(sm.filteredServices, service)
		}
	}
}

func (sm *ServiceManager) GetFilteredServices() []string {
	return sm.filteredServices
}

func (sm *ServiceManager) ToggleFilteredServices() {
	for _, service := range sm.filteredServices {
		sm.ToggleServiceSelection(service)
	}
}

func (sm *ServiceManager) saveSelectedServices() {
	util.WriteLines(sm.selectedServicePath, sm.selectedServices)
}

func (sm *ServiceManager) loadSelectedServices() {
	sm.selectedServices = util.ReadAllLines(sm.selectedServicePath)
}

func (sm *ServiceManager) ClearSelectedServices() {
	sm.selectedServices = []string{}
	sm.saveSelectedServices()
}

func (sm *ServiceManager) loadAvailableServices() {
	var services []string
	lines, err := executeMakeCommand(sm.getServiceMakeDir(), "services.mk", "help")
	if err != nil {
		services = []string{"dummy"} // Fallback or error handling
	}
	ansiRegex := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	for _, line := range lines {
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		service := ansiRegex.ReplaceAllString(fields[0], "")
		if !startsWithExcludedName(service) {
			services = append(services, service)
		}
	}
	sm.allServices = services
}

var excludedNames = []string{"make", "help", "[First"}

func startsWithExcludedName(line string) bool {
	for _, name := range excludedNames {
		if strings.HasPrefix(line, name) {
			return true
		}
	}
	return false
}

func (sm *ServiceManager) getServiceMakeDir() string {
	return sm.makeDir + "/services"
}

func executeMakeCommand(dirPath, file, target string) ([]string, error) {
	cmd := exec.Command("make", "-s", "-C", dirPath, "-f", file, target)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	return strings.Split(string(output), "\n"), nil
}

func (sm *ServiceManager) CleanServices() (string, error) {
	output, err := executeMakeCommand(sm.getServiceMakeDir(), "Makefile", "clean")
	if err != nil {
		return "", err
	}
	return strings.Join(output, "\n"), nil
}

func (sm *ServiceManager) SetupServices() (string, error) {
	output, err := executeMakeCommand(sm.getServiceMakeDir(), "Makefile", "setup")
	if err != nil {
		return "", err
	}
	return strings.Join(output, "\n"), nil
}

func (sm *ServiceManager) UpdateServices() (string, error) {
	output, err := executeMakeCommand(sm.getServiceMakeDir(), "Makefile", "update")
	if err != nil {
		return "", err
	}
	return strings.Join(output, "\n"), nil
}
