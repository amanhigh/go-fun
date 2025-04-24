package core

import (
	"time"

	"github.com/amanhigh/go-fun/components/kohan/clients"
	"github.com/amanhigh/go-fun/components/kohan/manager"
	"github.com/amanhigh/go-fun/components/kohan/manager/tui"
	"github.com/amanhigh/go-fun/components/kohan/repository"
	"github.com/amanhigh/go-fun/models/config"
	taxmodels "github.com/amanhigh/go-fun/models/tax"
	"github.com/go-resty/resty/v2"

	"github.com/golobby/container/v3"
	"github.com/rivo/tview"
)

// Interface and implementation in same file
type KohanInterface interface {
	GetDariusApp(cfg config.DariusConfig) (*DariusV1, error)
	// Add new method
	GetAutoManager(wait time.Duration, capturePath string) manager.AutoManagerInterface
	GetTaxManager() (manager.TaxManager, error) // Added method
}

// Private singleton instance
var globalInjector *KohanInjector

type KohanInjector struct {
	di     container.Container
	config config.KohanConfig
}

func SetupKohanInjector(config config.KohanConfig) {
	globalInjector = &KohanInjector{
		di:     container.New(),
		config: config,
	}
}

func (ki *KohanInjector) provideAlphaClient(client *resty.Client) clients.AlphaClient {
	return clients.NewAlphaClient(client, ki.config.Tax.AlphaBaseURL, ki.config.Tax.AlphaAPIKey)
}

func (ki *KohanInjector) provideSBIClient(client *resty.Client) clients.SBIClient {
	return clients.NewSBIClient(client, ki.config.Tax.SBIBaseURL)
}

func (ki *KohanInjector) provideTickerManager(client clients.AlphaClient) *manager.TickerManagerImpl {
	return manager.NewTickerManager(client, ki.config.Tax.DownloadsDir)
}

func (ki *KohanInjector) provideExchangeRepository() repository.ExchangeRepository {
	return repository.NewExchangeRepository(ki.config.Tax.SBIFilePath)
}

func (ki *KohanInjector) provideSBIManager(client clients.SBIClient, exchangeRepo repository.ExchangeRepository) manager.SBIManager {
	return manager.NewSBIManager(client, ki.config.Tax.SBIFilePath, exchangeRepo)
}

func (ki *KohanInjector) provideExchangeManager(sbiManager manager.SBIManager) manager.ExchangeManager {
	return manager.NewExchangeManager(sbiManager)
}

func (ki *KohanInjector) provideAccountRepository() repository.AccountRepository {
	return repository.NewAccountRepository(ki.config.Tax.AccountFilePath)
}

func (ki *KohanInjector) provideAccountManager(accountRepo repository.AccountRepository) manager.AccountManager {
	return manager.NewAccountManager(accountRepo)
}

func (ki *KohanInjector) provideTaxValuationManager(exchangeManager manager.ExchangeManager) manager.TaxValuationManager {
	return manager.NewTaxValuationManager(exchangeManager)
}

func (ki *KohanInjector) provideGainsRepository() repository.GainsRepository {
	return repository.NewGainsRepository(ki.config.Tax.GainsFilePath)
}

func (ki *KohanInjector) provideFinancialYearManagerGains() manager.FinancialYearManager[taxmodels.Gains] {
	return manager.NewFinancialYearManager[taxmodels.Gains]()
}

func (ki *KohanInjector) provideCapitalGainManager(
	exchangeMgr manager.ExchangeManager,
	gainsRepo repository.GainsRepository,
	fyMgr manager.FinancialYearManager[taxmodels.Gains],
) manager.CapitalGainManager {
	return manager.NewCapitalGainManager(exchangeMgr, gainsRepo, fyMgr)
}

func (ki *KohanInjector) provideTaxManager(gainMgr manager.CapitalGainManager) manager.TaxManager {
	return manager.NewTaxManager(gainMgr)
}

// Public singleton access - returns interface only
func GetKohanInterface() KohanInterface {
	return globalInjector
}

func (ki *KohanInjector) GetAutoManager(wait time.Duration, capturePath string) manager.AutoManagerInterface {
	// HACK: Move to Provider Based Build ?
	return manager.NewAutoManager(wait, capturePath)
}

func (ki *KohanInjector) GetTaxManager() (manager.TaxManager, error) {
	// Ensure all dependencies for TaxManager are registered
	container.MustSingleton(ki.di, resty.New)
	container.MustSingleton(ki.di, ki.provideSBIClient)

	// Register Repositories
	container.MustSingleton(ki.di, ki.provideExchangeRepository)
	container.MustSingleton(ki.di, ki.provideGainsRepository)

	// Register Managers (dependencies of TaxManager)
	container.MustSingleton(ki.di, ki.provideSBIManager)
	container.MustSingleton(ki.di, ki.provideExchangeManager)
	container.MustSingleton(ki.di, ki.provideFinancialYearManagerGains)
	container.MustSingleton(ki.di, ki.provideCapitalGainManager)

	// Register TaxManager itself
	container.MustSingleton(ki.di, ki.provideTaxManager)

	// Resolve and return TaxManager
	var taxManager manager.TaxManager
	err := ki.di.Resolve(&taxManager)
	return taxManager, err
}

func (ki *KohanInjector) GetDariusApp(cfg config.DariusConfig) (*DariusV1, error) {
	// Register config for this specific build
	container.MustSingleton(ki.di, func() config.DariusConfig {
		return cfg
	})

	// Register other dependencies
	container.MustSingleton(ki.di, tview.NewApplication)

	// Client
	container.MustSingleton(ki.di, resty.New)
	container.MustSingleton(ki.di, ki.provideAlphaClient)
	container.MustSingleton(ki.di, ki.provideSBIClient)

	// Repo
	container.MustSingleton(ki.di, ki.provideExchangeRepository)
	container.MustSingleton(ki.di, ki.provideAccountRepository)

	// Manager
	container.MustSingleton(ki.di, provideServiceManager)
	container.MustSingleton(ki.di, provideUIManager)
	container.MustSingleton(ki.di, provideHotkeyManager)
	container.MustSingleton(ki.di, ki.provideTickerManager)
	container.MustSingleton(ki.di, ki.provideAccountManager)
	container.MustSingleton(ki.di, ki.provideSBIManager)
	container.MustSingleton(ki.di, ki.provideExchangeManager)
	container.MustSingleton(ki.di, ki.provideTaxValuationManager)

	// Build app
	app := &DariusV1{}
	err := ki.di.Fill(app)
	return app, err
}

func provideServiceManager(cfg config.DariusConfig) (serviceManager tui.ServiceManager) {
	return tui.NewServiceManager(cfg.MakeDir, cfg.SelectedServiceFile)
}

func provideUIManager(app *tview.Application, svcManager tui.ServiceManager) tui.UIManager {
	return tui.NewUIManager(app, svcManager)
}

func provideHotkeyManager(uiManager tui.UIManager, serviceManager tui.ServiceManager) tui.HotkeyManager {
	return tui.NewHotkeyManager(uiManager, serviceManager)
}
