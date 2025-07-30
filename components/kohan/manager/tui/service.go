package tui

import (
	"strings"

	"github.com/amanhigh/go-fun/components/kohan/repository"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/rs/zerolog/log"
)

// ServiceFilterer provides filtering capabilities for Kubernetes services
type ServiceFilterer interface {
	// Filter Operations

	// FilterServices filters available services based on given keyword
	// Services containing keyword (case-insensitive) are added to filtered list
	FilterServices(keyword string)

	// GetFilteredServices returns services matching current filter
	GetFilteredServices() []string

	// ToggleFilteredServices toggles selection state of all filtered services
	// This affects all services currently in filtered list
	ToggleFilteredServices()
}

// ServiceManager provides management capabilities for Kubernetes services including
// selection, filtering and operations like setup, clean and update
type ServiceManager interface {
	ServiceFilterer // Embed the new interface

	// Service Management Operations

	// GetAllServices returns all available services discovered from the makefiles
	GetAllServices() []string

	// GetSelectedServices returns currently selected services for operations
	GetSelectedServices() []string

	// IsServiceSelected checks if given service is currently selected
	// Returns true if service is selected, false otherwise
	IsServiceSelected(service string) bool

	// ToggleServiceSelection toggles selection state of given service
	// If service is selected, it will be unselected and vice versa
	// Selected services are persisted to configured file path
	ToggleServiceSelection(service string)

	// ClearSelectedServices removes all services from selection
	// and persists empty selection to file
	ClearSelectedServices()

	// Service Operations

	// SetupServices runs setup make target on currently selected services
	// Returns output of make command and any error encountered
	SetupServices() (string, common.HttpError)

	// CleanServices runs clean make target on currently selected services
	// Returns output of make command and any error encountered
	CleanServices() (string, common.HttpError)

	// UpdateServices runs update make target on currently selected services
	// Returns output of make command and any error encountered
	UpdateServices() (string, common.HttpError)
}

func NewServiceManager(makeDir string, repo repository.TuiServiceRepository) *ServiceManagerImpl {
	manager := &ServiceManagerImpl{
		allServices:      []string{},
		selectedServices: []string{},
		makeDir:          makeDir,
		repo:             repo,
	}
	manager.loadAvailableServices()
	manager.loadSelectedServices()
	return manager
}

type ServiceManagerImpl struct {
	allServices      []string
	selectedServices []string
	filteredServices []string
	makeDir          string
	repo             repository.TuiServiceRepository
}

func (sm *ServiceManagerImpl) GetAllServices() []string {
	return sm.allServices
}

func (sm *ServiceManagerImpl) GetSelectedServices() []string {
	return sm.selectedServices
}

func (sm *ServiceManagerImpl) IsServiceSelected(service string) bool {
	for _, s := range sm.selectedServices {
		if s == service {
			return true
		}
	}
	return false
}

func (sm *ServiceManagerImpl) ToggleServiceSelection(service string) {
	if sm.IsServiceSelected(service) {
		sm.removeService(service)
	} else {
		sm.selectedServices = append(sm.selectedServices, service)
	}
	if err := sm.saveSelectedServices(); err != nil {
		log.Error().Err(err).Msg("Failed to save selected services")
	}
}

func (sm *ServiceManagerImpl) removeService(service string) {
	for i, s := range sm.selectedServices {
		if s == service {
			sm.selectedServices = append(sm.selectedServices[:i], sm.selectedServices[i+1:]...)
			break
		}
	}
}

func (sm *ServiceManagerImpl) FilterServices(keyword string) {
	sm.filteredServices = []string{}
	lowerKeyword := strings.ToLower(keyword)
	for _, service := range sm.allServices {
		if strings.Contains(strings.ToLower(service), lowerKeyword) {
			sm.filteredServices = append(sm.filteredServices, service)
		}
	}
}

func (sm *ServiceManagerImpl) GetFilteredServices() []string {
	return sm.filteredServices
}

func (sm *ServiceManagerImpl) ToggleFilteredServices() {
	for _, service := range sm.filteredServices {
		sm.ToggleServiceSelection(service)
	}
}

func (sm *ServiceManagerImpl) saveSelectedServices() common.HttpError {
	return sm.repo.SaveSelectedServices(sm.selectedServices)
}

func (sm *ServiceManagerImpl) loadSelectedServices() {
	services, err := sm.repo.LoadSelectedServices()
	if err != nil {
		log.Error().Err(err).Msg("Failed to load selected services from repository")
		sm.selectedServices = []string{} // Default to empty on error
		return
	}
	sm.selectedServices = services
}

func (sm *ServiceManagerImpl) ClearSelectedServices() {
	sm.selectedServices = []string{}
	if err := sm.saveSelectedServices(); err != nil {
		log.Error().Err(err).Msg("Failed to clear selected services")
	}
}

func (sm *ServiceManagerImpl) loadAvailableServices() {
	services, err := sm.repo.LoadAvailableServices(sm.makeDir)
	if err != nil {
		log.Error().Err(err).Str("makeDir", sm.makeDir).Msg("Failed to load available services from repository")
		sm.allServices = []string{"dummy"} // Fallback as per plan
		return
	}
	sm.allServices = services
}

func (sm *ServiceManagerImpl) getServiceMakeDir() string {
	return sm.makeDir + "/services"
}

func (sm *ServiceManagerImpl) CleanServices() (string, common.HttpError) {
	output, err := sm.repo.ExecuteMakeCommand(sm.getServiceMakeDir(), "Makefile", "clean")
	if err != nil {
		return "", err
	}
	return strings.Join(output, "\n"), nil
}

func (sm *ServiceManagerImpl) SetupServices() (string, common.HttpError) {
	output, err := sm.repo.ExecuteMakeCommand(sm.getServiceMakeDir(), "Makefile", "setup")
	if err != nil {
		return "", err
	}
	return strings.Join(output, "\n"), nil
}

func (sm *ServiceManagerImpl) UpdateServices() (string, common.HttpError) {
	output, err := sm.repo.ExecuteMakeCommand(sm.getServiceMakeDir(), "Makefile", "update")
	if err != nil {
		return "", err
	}
	return strings.Join(output, "\n"), nil
}
