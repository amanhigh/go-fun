package tui

import (
	"os/exec"
	"regexp"
	"strings"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/rs/zerolog/log"
)

// ServiceManager provides management capabilities for Kubernetes services including
// selection, filtering and operations like setup, clean and update
type ServiceManager interface {
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
	SetupServices() (string, error)

	// CleanServices runs clean make target on currently selected services
	// Returns output of make command and any error encountered
	CleanServices() (string, error)

	// UpdateServices runs update make target on currently selected services
	// Returns output of make command and any error encountered
	UpdateServices() (string, error)

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

// TODO: Move to Repository Layer ?
func NewServiceManager(makeDir, serviceFile string) *ServiceManagerImpl {
	manager := &ServiceManagerImpl{
		allServices:         []string{},
		selectedServices:    []string{},
		makeDir:             makeDir,
		selectedServicePath: serviceFile,
	}
	manager.loadAvailableServices()
	manager.loadSelectedServices()
	return manager
}

type ServiceManagerImpl struct {
	allServices         []string
	selectedServices    []string
	filteredServices    []string
	makeDir             string
	selectedServicePath string
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

func (sm *ServiceManagerImpl) saveSelectedServices() error {
	return util.WriteLines(sm.selectedServicePath, sm.selectedServices)
}

func (sm *ServiceManagerImpl) loadSelectedServices() {
	sm.selectedServices = util.ReadAllLines(sm.selectedServicePath)
}

func (sm *ServiceManagerImpl) ClearSelectedServices() {
	sm.selectedServices = []string{}
	if err := sm.saveSelectedServices(); err != nil {
		log.Error().Err(err).Msg("Failed to clear selected services")
	}
}

func (sm *ServiceManagerImpl) loadAvailableServices() {
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

func (sm *ServiceManagerImpl) getServiceMakeDir() string {
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

func (sm *ServiceManagerImpl) CleanServices() (string, error) {
	output, err := executeMakeCommand(sm.getServiceMakeDir(), "Makefile", "clean")
	if err != nil {
		return "", err
	}
	return strings.Join(output, "\n"), nil
}

func (sm *ServiceManagerImpl) SetupServices() (string, error) {
	output, err := executeMakeCommand(sm.getServiceMakeDir(), "Makefile", "setup")
	if err != nil {
		return "", err
	}
	return strings.Join(output, "\n"), nil
}

func (sm *ServiceManagerImpl) UpdateServices() (string, error) {
	output, err := executeMakeCommand(sm.getServiceMakeDir(), "Makefile", "update")
	if err != nil {
		return "", err
	}
	return strings.Join(output, "\n"), nil
}
