package core

import (
	"fmt"
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
	GetAutoManager(wait time.Duration, capturePath string) manager.AutoManagerInterface
	GetTaxManager() (manager.TaxManager, error)
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

// ---- Client Providers ----
func (ki *KohanInjector) provideAlphaClient(client *resty.Client) clients.AlphaClient {
	return clients.NewAlphaClient(client, ki.config.Tax.AlphaBaseURL, ki.config.Tax.AlphaAPIKey)
}

func (ki *KohanInjector) provideSBIClient(client *resty.Client) clients.SBIClient {
	return clients.NewSBIClient(client, ki.config.Tax.SBIBaseURL)
}

// ---- Repository Providers ----
func (ki *KohanInjector) provideExchangeRepository() repository.ExchangeRepository {
	return repository.NewExchangeRepository(ki.config.Tax.SBIFilePath)
}

func (ki *KohanInjector) provideGainsRepository() repository.GainsRepository {
	return repository.NewGainsRepository(ki.config.Tax.GainsFilePath)
}

func (ki *KohanInjector) provideDividendRepository() repository.DividendRepository {
	return repository.NewDividendRepository(ki.config.Tax.DividendFilePath)
}

func (ki *KohanInjector) provideInterestRepository() repository.InterestRepository {
	return repository.NewInterestRepository(ki.config.Tax.InterestFilePath)
}

func (ki *KohanInjector) provideAccountRepository() repository.AccountRepository {
	return repository.NewAccountRepository(ki.config.Tax.AccountFilePath)
}

func (ki *KohanInjector) provideTradeRepository() repository.TradeRepository {
	return repository.NewTradeRepository(ki.config.Tax.BrokerStatementPath)
}

// ---- Manager Providers ----
func (ki *KohanInjector) provideTickerManager(client clients.AlphaClient) *manager.TickerManagerImpl {
	return manager.NewTickerManager(client, ki.config.Tax.DownloadsDir)
}

func (ki *KohanInjector) provideSBIManager(client clients.SBIClient, exchangeRepo repository.ExchangeRepository) manager.SBIManager {
	return manager.NewSBIManager(client, ki.config.Tax.SBIFilePath, exchangeRepo)
}

func (ki *KohanInjector) provideExchangeManager(sbiManager manager.SBIManager) manager.ExchangeManager {
	return manager.NewExchangeManager(sbiManager)
}

func (ki *KohanInjector) provideAccountManager(accountRepo repository.AccountRepository) manager.AccountManager {
	return manager.NewAccountManager(accountRepo)
}

func (ki *KohanInjector) provideValuationManager(
	tickerManager manager.TickerManager,
	accountManager manager.AccountManager,
	tradeRepository repository.TradeRepository,
	fyManager manager.FinancialYearManager[taxmodels.Trade],
) manager.ValuationManager {
	return manager.NewValuationManager(tickerManager, accountManager, tradeRepository, fyManager)
}

func (ki *KohanInjector) provideTaxValuationManager(
	exchangeManager manager.ExchangeManager,
	valuationManager manager.ValuationManager,
) manager.TaxValuationManager {
	return manager.NewTaxValuationManager(exchangeManager, valuationManager)
}

func (ki *KohanInjector) provideFinancialYearManagerGains() manager.FinancialYearManager[taxmodels.Gains] {
	return manager.NewFinancialYearManager[taxmodels.Gains]()
}

func (ki *KohanInjector) provideFinancialYearManagerInterest() manager.FinancialYearManager[taxmodels.Interest] {
	return manager.NewFinancialYearManager[taxmodels.Interest]()
}

func (ki *KohanInjector) provideFinancialYearManagerDividends() manager.FinancialYearManager[taxmodels.Dividend] {
	return manager.NewFinancialYearManager[taxmodels.Dividend]()
}

func (ki *KohanInjector) provideFinancialYearManagerTrade() manager.FinancialYearManager[taxmodels.Trade] {
	return manager.NewFinancialYearManager[taxmodels.Trade]()
}

func (ki *KohanInjector) provideCapitalGainManager(
	exchangeMgr manager.ExchangeManager,
	gainsRepo repository.GainsRepository,
	fyMgr manager.FinancialYearManager[taxmodels.Gains],
) manager.CapitalGainManager {
	return manager.NewCapitalGainManager(exchangeMgr, gainsRepo, fyMgr)
}

func (ki *KohanInjector) provideInterestManager(
	exchangeMgr manager.ExchangeManager,
	fyMgr manager.FinancialYearManager[taxmodels.Interest],
	interestRepo repository.InterestRepository,
) manager.InterestManager {
	return manager.NewInterestManager(exchangeMgr, fyMgr, interestRepo)
}

func (ki *KohanInjector) provideDividendManager(
	exchangeMgr manager.ExchangeManager,
	fyMgr manager.FinancialYearManager[taxmodels.Dividend],
	dividendRepo repository.DividendRepository,
) manager.DividendManager {
	return manager.NewDividendManager(exchangeMgr, fyMgr, dividendRepo)
}

func (ki *KohanInjector) provideExcelManager() manager.ExcelManager {
	// Assuming ki.config.Tax.YearlySummaryExcelPath is available and validated.
	return manager.NewExcelManager(ki.config.Tax.YearlySummaryPath)
}

func (ki *KohanInjector) provideTaxManager(
	gainMgr manager.CapitalGainManager,
	dividendManager manager.DividendManager,
	interestManager manager.InterestManager,
	taxValuationManager manager.TaxValuationManager,
	excelMgr manager.ExcelManager,
) manager.TaxManager {
	return manager.NewTaxManager(gainMgr, dividendManager, interestManager, taxValuationManager, excelMgr)
}

func provideTuiServiceRepository(cfg config.DariusConfig) repository.TuiServiceRepository {
	return repository.NewTuiServiceRepository(cfg.SelectedServiceFile)
}

func provideServiceManager(cfg config.DariusConfig, repo repository.TuiServiceRepository) tui.ServiceManager {
	return tui.NewServiceManager(cfg.MakeDir, repo)
}

func provideUIManager(app *tview.Application, svcManager tui.ServiceManager) tui.UIManager {
	return tui.NewUIManager(app, svcManager)
}

func provideHotkeyManager(uiManager tui.UIManager, serviceManager tui.ServiceManager) tui.HotkeyManager {
	return tui.NewHotkeyManager(uiManager, serviceManager)
}

// ---- Dependency Registration ----

// registerClients registers core clients like REST client and API clients.
func (ki *KohanInjector) registerClients() {
	// First register the REST client
	container.MustSingleton(ki.di, resty.New)

	// Then register clients that depend on REST client
	var client *resty.Client
	container.MustResolve(ki.di, &client)

	container.MustSingleton(ki.di, func() clients.AlphaClient {
		return ki.provideAlphaClient(client)
	})
	container.MustSingleton(ki.di, func() clients.SBIClient {
		return ki.provideSBIClient(client)
	})
}

// registerRepositories registers all repository implementations.
func (ki *KohanInjector) registerRepositories() {
	container.MustSingleton(ki.di, ki.provideExchangeRepository)
	container.MustSingleton(ki.di, ki.provideGainsRepository)
	container.MustSingleton(ki.di, ki.provideDividendRepository)
	container.MustSingleton(ki.di, ki.provideInterestRepository)
	container.MustSingleton(ki.di, ki.provideAccountRepository)
	container.MustSingleton(ki.di, ki.provideTradeRepository)
}

// registerCoreManagers registers managers that depend on clients or repositories.
func (ki *KohanInjector) registerCoreManagers() {
	// Register managers that depend on clients or repositories
	container.MustSingleton(ki.di, ki.provideSBIManager)
	container.MustSingleton(ki.di, ki.provideExchangeManager)
	container.MustSingleton(ki.di, ki.provideAccountManager)

	// Register TickerManager (depends on AlphaClient)
	var alphaClient clients.AlphaClient
	container.MustResolve(ki.di, &alphaClient)
	container.MustSingleton(ki.di, func() manager.TickerManager {
		return ki.provideTickerManager(alphaClient)
	})

	// Register managers that depend on TickerManager and/or AccountManager/TradeRepository
	container.MustSingleton(ki.di, ki.provideValuationManager)
	container.MustSingleton(ki.di, ki.provideTaxValuationManager)
}

// registerFinancialYearManagers registers managers for handling financial year data.
func (ki *KohanInjector) registerFinancialYearManagers() {
	container.MustSingleton(ki.di, ki.provideFinancialYearManagerGains)
	container.MustSingleton(ki.di, ki.provideFinancialYearManagerDividends)
	container.MustSingleton(ki.di, ki.provideFinancialYearManagerInterest)
	container.MustSingleton(ki.di, ki.provideFinancialYearManagerTrade)
}

// registerTaxComponents registers managers specifically for tax calculations.
func (ki *KohanInjector) registerTaxComponents() {
	container.MustSingleton(ki.di, ki.provideCapitalGainManager)
	container.MustSingleton(ki.di, ki.provideDividendManager)
	container.MustSingleton(ki.di, ki.provideInterestManager)
	// Register ExcelManager first since TaxManager depends on it
	container.MustSingleton(ki.di, ki.provideExcelManager)
	container.MustSingleton(ki.di, ki.provideTaxManager)
}

// registerTaxDependencies registers all dependencies required for tax calculations.
func (ki *KohanInjector) registerTaxDependencies() {
	ki.registerClients()
	ki.registerRepositories()
	ki.registerFinancialYearManagers() // Register FY Managers first
	ki.registerCoreManagers()          // Then register managers that might depend on them
	ki.registerTaxComponents()
}

// registerDariusConfig registers configuration specific to the Darius application.
func (ki *KohanInjector) registerDariusConfig(cfg config.DariusConfig) {
	// Register config for this specific build
	container.MustSingleton(ki.di, func() config.DariusConfig {
		return cfg
	})

	// Register other dependencies
	container.MustSingleton(ki.di, tview.NewApplication)
}

// registerDariusTuiManagers registers managers specific to the Darius TUI.
// Other required managers (Ticker, Account, SBI, Exchange, Tax, etc.) are registered
// by calling registerBaseDependencies, registerRepositories, registerCoreManagers, etc.
// within registerDariusDependencies.
func (ki *KohanInjector) registerDariusTuiManagers() {
	// Register TuiServiceRepository first as ServiceManager depends on it.
	container.MustSingleton(ki.di, provideTuiServiceRepository)

	// Register ServiceManager, UIManager, HotkeyManager
	container.MustSingleton(ki.di, provideServiceManager)
	container.MustSingleton(ki.di, provideUIManager)
	container.MustSingleton(ki.di, provideHotkeyManager)
}

// registerDariusDependencies registers all dependencies required for the Darius application.
func (ki *KohanInjector) registerDariusDependencies(cfg config.DariusConfig) {
	// Register base dependencies, repositories, and all managers needed for tax calculations
	// This ensures all components potentially used by Darius are available.
	ki.registerTaxDependencies() // Includes CapitalGain, Dividend, Interest, Tax managers

	// Register Darius specific configurations and TUI managers
	ki.registerDariusConfig(cfg)
	ki.registerDariusTuiManagers()
}
