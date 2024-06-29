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
	cmd := exec.Command("make", "-C", "/home/aman/Projects/go-fun/Kubernetes/services", "-f", "services.mk")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return []string{"Error fetching services"} // Fallback or error handling
	}
	lines := strings.Split(string(output), "\n")
	var services []string
	ansiRegex := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	for _, line := range lines {
		if strings.Contains(line, "make: Leaving directory") {
			break
		}
		if trimmed := strings.TrimLeft(line, " \t"); len(trimmed) > 0 && !strings.HasPrefix(trimmed, "make") && !strings.HasPrefix(trimmed, "[") {
			fields := strings.Fields(trimmed)
			if len(fields) > 0 {
				service := ansiRegex.ReplaceAllString(fields[0], "")
				services = append(services, service)
			}
		}
	}
	return services
}
