package tui

import (
	"os/exec"
	"regexp"
	"strings"

	"github.com/amanhigh/go-fun/common/util"
)

type ServiceManager struct {
	allServices      []string
	selectedServices []string
	filteredServices []string
	serviceFilePath  string
}

func NewServiceManager(services []string) (sm *ServiceManager) {
	sm = &ServiceManager{
		allServices:      fetchServices(),
		selectedServices: []string{},
		serviceFilePath:  "/tmp/k8-svc.txt",
	}
	sm.loadSelectedServices()
	return
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
	util.WriteLines(sm.serviceFilePath, sm.selectedServices)
}

func (sm *ServiceManager) loadSelectedServices() {
	sm.selectedServices = util.ReadAllLines(sm.serviceFilePath)
}

func (sm *ServiceManager) ClearSelectedServices() {
	sm.selectedServices = []string{}
	sm.saveSelectedServices()
}

func fetchServices() []string {
	lines, err := executeMakeCommand("/home/aman/Projects/go-fun/Kubernetes/services", "services.mk", "help")
	if err != nil {
		return []string{"Error fetching services"} // Fallback or error handling
	}
	var services []string
	ansiRegex := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	for _, line := range lines {
		// Skip empty lines and lines starting with "make" or "make:" as they are not services
		if line == "" || strings.HasPrefix(line, "make") {
			continue
		}

		fields := strings.Fields(line)
		service := ansiRegex.ReplaceAllString(fields[0], "")
		services = append(services, service)
	}
	return services
}

func executeMakeCommand(dirPath, file, target string) ([]string, error) {
	cmd := exec.Command("make", "-s", "-C", dirPath, "-f", file, target)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	return strings.Split(string(output), "\n"), nil
}
