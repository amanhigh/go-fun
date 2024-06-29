package tui

type ServiceManager struct {
	allServices      []string
	selectedServices []string
}

func NewServiceManager(services []string) *ServiceManager {
	return &ServiceManager{
		allServices:      services,
		selectedServices: []string{},
	}
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
}

func (sm *ServiceManager) removeService(service string) {
	for i, s := range sm.selectedServices {
		if s == service {
			sm.selectedServices = append(sm.selectedServices[:i], sm.selectedServices[i+1:]...)
			break
		}
	}
}
