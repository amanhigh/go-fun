package repository

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/rs/zerolog/log" // Added for logging consistency if needed
)

//go:generate mockery --name TuiServiceRepository
type TuiServiceRepository interface {
	LoadAvailableServices(makeDir string) ([]string, error)
	LoadSelectedServices() ([]string, error)
	SaveSelectedServices(services []string) error
	ExecuteMakeCommand(makeDir, file, target string) ([]string, error)
}

type tuiServiceRepositoryImpl struct {
	selectedServicePath string
}

func (r *tuiServiceRepositoryImpl) ExecuteMakeCommand(makeDir, file, target string) ([]string, error) {
	cmd := exec.Command("make", "-s", "-C", makeDir, "-f", file, target)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to execute make command: %w. Output: %s", err, string(output))
	}
	return strings.Split(string(output), "\n"), nil
}

// NewTuiServiceRepository creates a new TuiServiceRepository.
func NewTuiServiceRepository(selectedServicePath string) TuiServiceRepository {
	return &tuiServiceRepositoryImpl{selectedServicePath: selectedServicePath}
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

func getServiceMakeDir(makeDir string) string {
	return makeDir + "/services"
}

func (r *tuiServiceRepositoryImpl) LoadAvailableServices(makeDir string) ([]string, error) {
	var services []string
	// Use the helper function for makeDir, passing the parameter
	lines, err := r.ExecuteMakeCommand(getServiceMakeDir(makeDir), "services.mk", "help")
	if err != nil {
		log.Error().Err(err).Str("makeDir", makeDir).Msg("Failed to load available services via make command")
		// Return error instead of fallback, let caller decide
		return nil, fmt.Errorf("LoadAvailableServices: %w", err)
	}

	ansiRegex := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine == "" {
			continue
		}
		fields := strings.Fields(trimmedLine)
		if len(fields) > 0 {
			service := ansiRegex.ReplaceAllString(fields[0], "")
			if !startsWithExcludedName(service) {
				services = append(services, service)
			}
		}
	}
	if len(services) == 0 {
		log.Warn().Str("makeDir", makeDir).Msg("No available services found or all were excluded.")
		// It's possible no services are found, not necessarily an error.
		// If an empty list is an error state, the caller should handle it.
	}
	return services, nil
}

func (r *tuiServiceRepositoryImpl) LoadSelectedServices() ([]string, error) {
	if _, err := os.Stat(r.selectedServicePath); os.IsNotExist(err) {
		// File not existing is not an error for loading, means no services selected yet.
		// This matches original logic where sm.selectedServices would remain empty.
		return []string{}, nil
	}
	// util.ReadAllLines does not return an error in its signature,
	// but it can panic. It's better to use os.ReadFile and handle errors.
	content, err := os.ReadFile(r.selectedServicePath)
	if err != nil {
		log.Error().Err(err).Str("path", r.selectedServicePath).Msg("Failed to read selected services file")
		return nil, fmt.Errorf("LoadSelectedServices: reading file: %w", err)
	}

	lines := strings.Split(strings.TrimSpace(string(content)), "\n")
	// Filter out empty lines that might result from an empty file or trailing newlines
	var selectedServices []string
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			selectedServices = append(selectedServices, strings.TrimSpace(line))
		}
	}
	return selectedServices, nil
}

func (r *tuiServiceRepositoryImpl) SaveSelectedServices(services []string) error {
	err := util.WriteLines(r.selectedServicePath, services)
	if err != nil {
		log.Error().Err(err).Str("path", r.selectedServicePath).Msg("Failed to save selected services to file")
		return fmt.Errorf("SaveSelectedServices: writing lines: %w", err)
	}
	return nil
}
