package core

import (
	"fmt"
	"time"

	"github.com/amanhigh/go-fun/components/kohan/handler"
	"github.com/amanhigh/go-fun/components/kohan/manager"
	"github.com/amanhigh/go-fun/models/config"

	"github.com/golobby/container/v3"
)

// Interface and implementation in same file
type KohanInterface interface {
	GetDariusApp(cfg config.DariusConfig) (*DariusV1, error)
	GetAutoManager(wait time.Duration, capturePath string) manager.AutoManagerInterface
	GetTaxManager() (manager.TaxManager, error)
	GetBrokerageManager() (manager.BrokerageManager, error)
	// HACK: Change it to GetKohanServer and introduce/link to a CLI Command.
	GetJournalHandler() (*handler.JournalHandler, error)
}

// Private singleton instance
var globalInjector *KohanInjector

type KohanInjector struct {
	di     container.Container
	config config.KohanConfig
}

// ---- KohanInjector Methods ----

func SetupKohanInjector(config config.KohanConfig) {
	globalInjector = &KohanInjector{
		di:     container.New(),
		config: config,
	}
}

// Public singleton access - returns interface only
func GetKohanInterface() KohanInterface {
	return globalInjector
}

func (ki *KohanInjector) GetAutoManager(wait time.Duration, capturePath string) manager.AutoManagerInterface {
	return manager.NewAutoManager(wait, capturePath)
}

func (ki *KohanInjector) GetJournalHandler() (*handler.JournalHandler, error) {
	if err := ki.registerJournalDependencies(); err != nil {
		return nil, err
	}

	var journalHandler *handler.JournalHandler
	if err := ki.di.Resolve(&journalHandler); err != nil {
		return nil, fmt.Errorf("failed to resolve journal handler: %w", err)
	}
	return journalHandler, nil
}

func (ki *KohanInjector) GetTaxManager() (manager.TaxManager, error) {
	ki.registerTaxDependencies()

	// Resolve and return TaxManager
	var taxManager manager.TaxManager
	err := ki.di.Resolve(&taxManager)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve tax manager: %w", err)
	}
	return taxManager, nil
}

func (ki *KohanInjector) GetBrokerageManager() (manager.BrokerageManager, error) {
	ki.registerTaxDependencies()

	var brokerageManager manager.BrokerageManager
	err := ki.di.Resolve(&brokerageManager)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve brokerage manager: %w", err)
	}
	return brokerageManager, nil
}

func (ki *KohanInjector) GetDariusApp(cfg config.DariusConfig) (*DariusV1, error) {
	ki.registerDariusDependencies(cfg)

	// Build app
	app := &DariusV1{}
	err := ki.di.Fill(app)
	if err != nil {
		return nil, fmt.Errorf("failed to fill darius app: %w", err)
	}
	return app, nil
}
