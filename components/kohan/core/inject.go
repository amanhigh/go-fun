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

// Add this new provider function
func (ki *KohanInjector) provideValuationManager(
	tickerManager manager.TickerManager,
	accountManager manager.AccountManager,
	tradeRepository repository.TradeRepository,
) manager.ValuationManager {
	return manager.NewValuationManager(tickerManager, accountManager, tradeRepository)
}

func (ki *KohanInjector) provideTaxValuationManager(
	exchangeManager manager.ExchangeManager,
	valuationManager manager.ValuationManager, // Added parameter
) manager.TaxValuationManager {
	return manager.NewTaxValuationManager(exchangeManager, valuationManager) // Pass new dependency
}

func (ki *KohanInjector) provideGainsRepository() repository.GainsRepository {
	return repository.NewGainsRepository(ki.config.Tax.GainsFilePath)
}

func (ki *KohanInjector) provideFinancialYearManagerGains() manager.FinancialYearManager[taxmodels.Gains] {
	return manager.NewFinancialYearManager[taxmodels.Gains]()
}

// Add new provider for FinancialYearManager[tax.Interest]
func (ki *KohanInjector) provideFinancialYearManagerInterest() manager.FinancialYearManager[taxmodels.Interest] {
	return manager.NewFinancialYearManager[taxmodels.Interest]()
}

func (ki *KohanInjector) provideCapitalGainManager(
	exchangeMgr manager.ExchangeManager,
	gainsRepo repository.GainsRepository,
	fyMgr manager.FinancialYearManager[taxmodels.Gains],
) manager.CapitalGainManager {
	return manager.NewCapitalGainManager(exchangeMgr, gainsRepo, fyMgr)
}

// Add/Update provider for InterestManager
func (ki *KohanInjector) provideInterestManager(
	exchangeMgr manager.ExchangeManager,
	fyMgr manager.FinancialYearManager[taxmodels.Interest], // Added parameter
	interestRepo repository.InterestRepository, // Added parameter
) manager.InterestManager {
	// Call updated constructor
	return manager.NewInterestManager(exchangeMgr, fyMgr, interestRepo)
}

func (ki *KohanInjector) provideFinancialYearManagerDividends() manager.FinancialYearManager[taxmodels.Dividend] {
	return manager.NewFinancialYearManager[taxmodels.Dividend]()
}

func (ki *KohanInjector) provideDividendRepository() repository.DividendRepository {
	return repository.NewDividendRepository(ki.config.Tax.DividendFilePath)
}

// Add new provider for InterestRepository (Adjust config path if needed)
func (ki *KohanInjector) provideInterestRepository() repository.InterestRepository {
	// Ensure ki.config.Tax.InterestFilePath is the correct config field
	return repository.NewInterestRepository(ki.config.Tax.InterestFilePath)
}

// Add new provider function
func (ki *KohanInjector) provideTradeRepository() repository.TradeRepository {
	// Ensure ki.config.Tax.BrokerStatementPath is the correct config field
	return repository.NewTradeRepository(ki.config.Tax.BrokerStatementPath)
}

// Update provider for TaxManager to accept TaxValuationManager
func (ki *KohanInjector) provideTaxManager(
	gainMgr manager.CapitalGainManager,
	dividendManager manager.DividendManager,
	interestManager manager.InterestManager,
	taxValuationManager manager.TaxValuationManager, // Added parameter
) manager.TaxManager {
	return manager.NewTaxManager(gainMgr, dividendManager, interestManager, taxValuationManager) // Pass new dependency
}

// Public singleton access - returns interface only
func GetKohanInterface() KohanInterface {
	return globalInjector
}

func (ki *KohanInjector) GetAutoManager(wait time.Duration, capturePath string) manager.AutoManagerInterface {
	// HACK: Move to Provider Based Build ?
	return manager.NewAutoManager(wait, capturePath)
}

func (ki *KohanInjector) registerBaseDependencies() {
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

func (ki *KohanInjector) registerRepositories() {
	container.MustSingleton(ki.di, ki.provideExchangeRepository)
	container.MustSingleton(ki.di, ki.provideGainsRepository)
	container.MustSingleton(ki.di, ki.provideDividendRepository)
	container.MustSingleton(ki.di, ki.provideInterestRepository)
	container.MustSingleton(ki.di, ki.provideAccountRepository)
	container.MustSingleton(ki.di, ki.provideTradeRepository)
}

func (ki *KohanInjector) registerCoreManagers() {
	// First register managers that don't have dependencies
	container.MustSingleton(ki.di, ki.provideSBIManager)
	container.MustSingleton(ki.di, ki.provideExchangeManager)
	container.MustSingleton(ki.di, ki.provideAccountManager)

	// Then register TickerManager with resolved AlphaClient
	var alphaClient clients.AlphaClient
	container.MustResolve(ki.di, &alphaClient)
	container.MustSingleton(ki.di, func() manager.TickerManager {
		return ki.provideTickerManager(alphaClient)
	})

	// Finally register managers that depend on TickerManager
	container.MustSingleton(ki.di, ki.provideValuationManager)
	container.MustSingleton(ki.di, ki.provideTaxValuationManager)
}

func (ki *KohanInjector) registerFinancialYearManagers() {
	container.MustSingleton(ki.di, ki.provideFinancialYearManagerGains)
	container.MustSingleton(ki.di, ki.provideFinancialYearManagerDividends)
	container.MustSingleton(ki.di, ki.provideFinancialYearManagerInterest)
}

func (ki *KohanInjector) registerTaxComponents() {
	container.MustSingleton(ki.di, ki.provideCapitalGainManager)
	container.MustSingleton(ki.di, ki.provideDividendManager)
	container.MustSingleton(ki.di, ki.provideInterestManager)
	container.MustSingleton(ki.di, ki.provideTaxManager)
}

func (ki *KohanInjector) registerTaxDependencies() {
	ki.registerBaseDependencies()
	ki.registerRepositories()
	ki.registerCoreManagers()
	ki.registerFinancialYearManagers()
	ki.registerTaxComponents()
}

func (ki *KohanInjector) GetTaxManager() (manager.TaxManager, error) {
	ki.registerTaxDependencies()

	// Resolve and return TaxManager
	var taxManager manager.TaxManager
	err := ki.di.Resolve(&taxManager)
	return taxManager, err
}

func (ki *KohanInjector) GetDariusApp(cfg config.DariusConfig) (*DariusV1, error) {
	ki.registerDariusDependencies(cfg)

	// Build app
	app := &DariusV1{}
	err := ki.di.Fill(app)
	return app, err
}

func (ki *KohanInjector) registerDariusDependencies(cfg config.DariusConfig) {
	ki.registerDariusConfig(cfg)
	ki.registerDariusClientsAndRepos()
	ki.registerDariusManagers()
}

func (ki *KohanInjector) registerDariusConfig(cfg config.DariusConfig) {
	// Register config for this specific build
	container.MustSingleton(ki.di, func() config.DariusConfig {
		return cfg
	})

	// Register other dependencies
	container.MustSingleton(ki.di, tview.NewApplication)
}

func (ki *KohanInjector) registerDariusClientsAndRepos() {
	// Client
	container.MustSingleton(ki.di, resty.New)
	container.MustSingleton(ki.di, ki.provideAlphaClient)
	container.MustSingleton(ki.di, ki.provideSBIClient)

	// Repo
	container.MustSingleton(ki.di, ki.provideExchangeRepository)
	container.MustSingleton(ki.di, ki.provideAccountRepository)
	container.MustSingleton(ki.di, ki.provideDividendRepository)
	container.MustSingleton(ki.di, ki.provideInterestRepository)
}

func (ki *KohanInjector) registerDariusManagers() {
	// Manager
	// FIXME: Remove Duplicate Container Registration
	// FIXME: Remove Duplicate Container Registration
	container.MustSingleton(ki.di, provideServiceManager)
	container.MustSingleton(ki.di, provideUIManager)
	container.MustSingleton(ki.di, provideHotkeyManager)
	container.MustSingleton(ki.di, ki.provideTickerManager)
	container.MustSingleton(ki.di, ki.provideAccountManager)
	container.MustSingleton(ki.di, ki.provideSBIManager)
	container.MustSingleton(ki.di, ki.provideExchangeManager)
	container.MustSingleton(ki.di, ki.provideTaxValuationManager)
	container.MustSingleton(ki.di, ki.provideFinancialYearManagerGains) // Ensure this exists
	container.MustSingleton(ki.di, ki.provideFinancialYearManagerDividends)
	container.MustSingleton(ki.di, ki.provideFinancialYearManagerInterest) // Added or ensure exists
	container.MustSingleton(ki.di, ki.provideCapitalGainManager)           // Ensure this exists
	container.MustSingleton(ki.di, ki.provideDividendManager)
	container.MustSingleton(ki.di, ki.provideInterestManager) // Ensure this uses the updated provider
	// Ensure TaxManager registration (if used) uses the updated provider
	// container.MustSingleton(ki.di, ki.provideTaxManager)
}

// Update provider for DividendManager
func (ki *KohanInjector) provideDividendManager(
	exchangeMgr manager.ExchangeManager,
	fyMgr manager.FinancialYearManager[taxmodels.Dividend], // Added parameter
	dividendRepo repository.DividendRepository, // Added parameter
) manager.DividendManager {
	// Call updated constructor
	return manager.NewDividendManager(exchangeMgr, fyMgr, dividendRepo)
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
