package tui

import (
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
		allServices:      services,
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
