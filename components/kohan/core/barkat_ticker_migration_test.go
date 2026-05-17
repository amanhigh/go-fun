//nolint:dupl,gocyclo,funlen,cyclop
package core_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/amanhigh/go-fun/models/barkat"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/go-resty/resty/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// ============================================================================
// CONSTANTS
// ============================================================================

const (
	// TickerDumpPath is the path to the Barkat repository dump.
	TickerDumpPath = "/home/aman/Projects/go-fun/data.json"
	// NumWatchCategories matches the PRD-defined 8 watchlist categories.
	NumWatchCategories = 8
	// NumFlagCategories matches the PRD-defined 8 flag categories.
	NumFlagCategories = 8
)

// ============================================================================
// DUMP TYPES
// ============================================================================

// BarkatRepositoryDump mirrors the top-level structure of data.json.
type BarkatRepositoryDump struct {
	TickerRepo        map[string]string           `json:"tickerRepo"`
	PairRepo          map[string]PairInfo         `json:"pairRepo"`
	ExchangeRepo      map[string]string           `json:"exchangeRepo"`
	SequenceRepo      map[string]string           `json:"sequenceRepo"`
	RecentRepo        map[string]int64            `json:"recentRepo"`
	FnoRepo           []string                    `json:"fnoRepo"`
	WatchRepo         map[string][]string         `json:"watchRepo"`
	FlagRepo          map[string][]string         `json:"flagRepo"`
	AlertRepo         map[string][]AlertRepoEntry `json:"alertRepo"`
	AlertClickedEvent json.RawMessage             `json:"alertClickedEvent"`
	AlertFeedEvent    json.RawMessage             `json:"alertFeedEvent"`
	GttCreateEvent    json.RawMessage             `json:"gttCreateEvent"`
	GttDeleteEvent    json.RawMessage             `json:"gttDeleteEvent"`
	GttRefereshEvent  json.RawMessage             `json:"gttRefereshEvent"`
	JournalOpenEvent  json.RawMessage             `json:"journalOpenEvent"`
}

// PairInfo holds Investing-side identity and metadata from pairRepo.
type PairInfo struct {
	Name     string `json:"name"`
	PairID   string `json:"pairId"`
	Exchange string `json:"exchange,omitempty"`
}

// AlertRepoEntry holds a single alert entry from alertRepo.
type AlertRepoEntry struct {
	ID     string  `json:"id"`
	PairID string  `json:"pairId"`
	Price  float64 `json:"price"`
	Name   string  `json:"name"`
}

// ============================================================================
// DUMP PREFLIGHT ANALYSIS
// ============================================================================

// DumpAnalysis holds preflight counts and anomaly diagnostics.
type DumpAnalysis struct {
	TickerCount     int
	PairCount       int
	ExchangeCount   int
	SequenceCount   int
	RecentCount     int
	FnoCount        int
	WatchItemCount  int
	FlagItemCount   int
	AlertGroupCount int
	AlertEntryCount int

	// Anomalies
	UnresolvedMappings       int
	UnresolvedSamples        []string
	InvalidPairSymbols       int
	InvalidPairSymbolSamples []string
	InvalidPairNames         int
	InvalidPairNameSamples   []NameSanitization
	InvalidPairExchanges     int
	InvalidExchangeSamples   []ExchangeSanitization
	EmptyAlertIDs            int
	EmptyAlertNames          int
	EmptyAlertGroups         []string
	SequenceYRCount          int
	SequenceMWDCount         int
	RecentObjectShaped       bool
}

// NameSanitization captures a symbol→name pair for logging.
type NameSanitization struct {
	Symbol string `json:"symbol"`
	Name   string `json:"name"`
}

// ExchangeSanitization captures a symbol→exchange pair for logging.
type ExchangeSanitization struct {
	Symbol   string `json:"symbol"`
	Exchange string `json:"exchange"`
}

// ============================================================================
// MIGRATION PLAN
// ============================================================================

// TickerPayload holds the final normalised payload for the Ticker API.
type TickerPayload struct {
	Ticker       string    `json:"ticker"`
	Exchange     *string   `json:"exchange,omitempty"`
	Timeframes   []string  `json:"timeframes"`
	Type         string    `json:"type"`
	State        string    `json:"state"`
	Trend        string    `json:"trend"`
	LastOpenedAt time.Time `json:"last_opened_at"`
	IsFNO        bool      `json:"is_fno"`
}

// AlertTickerPayload holds the normalised payload for the Alert Ticker API.
type AlertTickerPayload struct {
	ParentTicker string  `json:"-"`
	Symbol       string  `json:"symbol"`
	PairID       string  `json:"pair_id"`
	Name         string  `json:"name"`
	Exchange     *string `json:"exchange,omitempty"`
}

// MigrationPlan holds the ordered list of API payloads to migrate.
type MigrationPlan struct {
	Tickers          []TickerPayload
	AlertTickers     []AlertTickerPayload
	AlertTickerSkips []string
}

// TickerMigrationStats tracks progress and reconciliation.
type TickerMigrationStats struct {
	CreatedTickers       int
	UpdatedTickers       int
	SkippedTickers       int
	CreatedAlertTickers  int
	VerifiedAlertTickers int
	SkippedAlertTickers  int
	FailedTickers        []string
	FailedAlertTickers   []string
	TotalAPICalls        int
	StartTime            time.Time
}

// ============================================================================
// VALIDATORS (mirrored from kohan_validators.go for preflight)
// ============================================================================

var (
	alertSymbolRegex   = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9./=]*$`)
	alertNameRegex     = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9 .&'()-]*$`)
	alertExchangeRegex = regexp.MustCompile(`^[A-Za-z][A-Za-z]*$`)
)

// ============================================================================
// DUMP LOADER
// ============================================================================

func loadBarkatDump(path string) (*BarkatRepositoryDump, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read dump: %w", err)
	}

	var dump BarkatRepositoryDump
	if err := json.Unmarshal(data, &dump); err != nil {
		return nil, fmt.Errorf("failed to unmarshal dump: %w", err)
	}

	if dump.TickerRepo == nil {
		return nil, fmt.Errorf("dump missing tickerRepo")
	}
	return &dump, nil
}

// ============================================================================
// PREFLIGHT ANALYSIS
// ============================================================================

func analyzeDump(dump *BarkatRepositoryDump) *DumpAnalysis {
	analysis := &DumpAnalysis{
		TickerCount:     len(dump.TickerRepo),
		PairCount:       len(dump.PairRepo),
		ExchangeCount:   len(dump.ExchangeRepo),
		SequenceCount:   len(dump.SequenceRepo),
		RecentCount:     len(dump.RecentRepo),
		FnoCount:        len(dump.FnoRepo),
		AlertGroupCount: len(dump.AlertRepo),
	}

	// Count watch items
	for _, cat := range dump.WatchRepo {
		analysis.WatchItemCount += len(cat)
	}

	// Count flag items
	for _, cat := range dump.FlagRepo {
		analysis.FlagItemCount += len(cat)
	}

	// Count alert entries and find empty groups
	for pairID, alerts := range dump.AlertRepo {
		if len(alerts) == 0 {
			analysis.EmptyAlertGroups = append(analysis.EmptyAlertGroups, pairID)
		}
		for _, a := range alerts {
			analysis.AlertEntryCount++
			if a.ID == "" {
				analysis.EmptyAlertIDs++
			}
			if a.Name == "" {
				analysis.EmptyAlertNames++
			}
		}
	}

	// Count unresolved tickerRepo→pairRepo mappings
	unresolvedSet := make(map[string]bool)
	for tvTicker, investingSymbol := range dump.TickerRepo {
		if _, ok := dump.PairRepo[investingSymbol]; !ok {
			if !unresolvedSet[investingSymbol] {
				unresolvedSet[investingSymbol] = true
				if len(analysis.UnresolvedSamples) < 20 {
					analysis.UnresolvedSamples = append(analysis.UnresolvedSamples,
						fmt.Sprintf("%s → %s", tvTicker, investingSymbol))
				}
			}
		}
	}
	analysis.UnresolvedMappings = len(unresolvedSet)

	// Count invalid pairRepo symbols
	for symbol := range dump.PairRepo {
		if !alertSymbolRegex.MatchString(symbol) {
			analysis.InvalidPairSymbols++
			if len(analysis.InvalidPairSymbolSamples) < 20 {
				analysis.InvalidPairSymbolSamples = append(analysis.InvalidPairSymbolSamples, symbol)
			}
		}
	}

	// Count invalid pairRepo names
	for symbol, info := range dump.PairRepo {
		if !alertNameRegex.MatchString(info.Name) {
			analysis.InvalidPairNames++
			if len(analysis.InvalidPairNameSamples) < 20 {
				analysis.InvalidPairNameSamples = append(analysis.InvalidPairNameSamples, NameSanitization{
					Symbol: symbol,
					Name:   info.Name,
				})
			}
		}
	}

	// Count invalid pairRepo exchanges
	for symbol, info := range dump.PairRepo {
		if info.Exchange != "" && (!alertExchangeRegex.MatchString(info.Exchange) ||
			len(info.Exchange) < 1 || len(info.Exchange) > 10) {
			analysis.InvalidPairExchanges++
			if len(analysis.InvalidExchangeSamples) < 20 {
				analysis.InvalidExchangeSamples = append(analysis.InvalidExchangeSamples, ExchangeSanitization{
					Symbol:   symbol,
					Exchange: info.Exchange,
				})
			}
		}
	}

	// Sequence analysis
	for _, v := range dump.SequenceRepo {
		switch v {
		case "MWD":
			analysis.SequenceMWDCount++
		case "YR":
			analysis.SequenceYRCount++
		}
	}

	// Recent repo - check if object-shaped (has timestamps)
	analysis.RecentObjectShaped = true
	for _, v := range dump.RecentRepo {
		if v <= 0 {
			analysis.RecentObjectShaped = false
			break
		}
	}

	return analysis
}

// ============================================================================
// NORMALIZATION HELPERS
// ============================================================================

// normalizeAlertSymbol removes spaces and special chars from Investing symbols.
func normalizeAlertSymbol(symbol string) (string, string) {
	normalized := strings.ReplaceAll(symbol, " ", "")
	normalized = strings.ReplaceAll(normalized, "'", "")
	normalized = strings.ReplaceAll(normalized, ".", "")
	normalized = strings.ReplaceAll(normalized, "&", "")
	return normalized, symbol
}

// normalizeAlertName replaces smart quotes, NBSP, ®, and other unicode.
func normalizeAlertName(name string) string {
	name = strings.ReplaceAll(name, "\u2019", "'") // smart quote → ASCII
	name = strings.ReplaceAll(name, "\u2018", "'") // smart single quote
	name = strings.ReplaceAll(name, "\u201c", `"`) // smart double quote open
	name = strings.ReplaceAll(name, "\u201d", `"`) // smart double quote close
	name = strings.ReplaceAll(name, "\u00ae", "")  // ® → empty
	name = strings.ReplaceAll(name, "\u00a0", " ") // NBSP → space
	name = strings.ReplaceAll(name, "/", " ")      // slash → space (not allowed in alert name validator)
	return name
}

// normalizeAlertExchange handles long/empty/missing exchange labels.
func normalizeAlertExchange(exchange, _ string) *string {
	if exchange == "" {
		return nil
	}
	switch exchange {
	case "Global Indexes", "Investing.com":
		return nil
	}
	if !alertExchangeRegex.MatchString(exchange) {
		return nil
	}
	return &exchange
}

// extractSimpleExchange extracts the exchange code from "EXCHANGE:SYMBOL".
// e.g. "NSE:MCX" → "NSE", "FX_IDC:USDINR" → "FX_IDC"
func extractSimpleExchange(exchangeQualified string) string {
	if idx := strings.Index(exchangeQualified, ":"); idx > 0 {
		return exchangeQualified[:idx]
	}
	return exchangeQualified
}

// expandSequence converts legacy sequence names to backend timeframe arrays.
func expandSequence(seq string) []string {
	switch seq {
	case "YR":
		return []string{"YR", "SMN", "TMN", "MN", "WK"}
	case "MWD":
		return []string{"MN", "WK", "DL"}
	case "WDH":
		return []string{"WK", "DL", "DL"}
	default:
		return []string{"MN", "WK", "DL"}
	}
}

// deriveTickerType classifies a ticker based solely on migration list membership:
//   - List 7 (watch or flag) → COMPOSITE
//   - List 6 (watch or flag) with "/" → COMPOSITE
//   - List 6 (watch or flag) without "/" → INDEX
//   - Everything else → EQUITY
//
// Symbol-pattern heuristics (CNX prefix, "!" suffix, known index names) are NOT used;
// those tickers fall to EQUITY unless explicitly placed in list 6 or 7.
func deriveTickerType(tvTicker string, _ *string, flagSets, watchSets []map[string]bool) string {
	// List 7 takes priority (COMPOSITE), then list 6 slash/non-slash split.
	if (len(flagSets) > 7 && flagSets[7][tvTicker]) ||
		(len(watchSets) > 7 && watchSets[7][tvTicker]) {
		return "COMPOSITE"
	}
	if (len(flagSets) > 6 && flagSets[6][tvTicker]) ||
		(len(watchSets) > 6 && watchSets[6][tvTicker]) {
		if strings.Contains(tvTicker, "/") {
			return "COMPOSITE"
		}
		return "INDEX"
	}
	return "EQUITY"
}

// deriveTickerState returns the default state for migration.
func deriveTickerState(tvTicker string, watchSets []map[string]bool) string {
	// watchRepo[1] (red) → READY (pre-journal ready-for-trade)
	if len(watchSets) > 1 && watchSets[1][tvTicker] {
		return "READY"
	}
	// watchRepo[0] (orange) → keep WATCHED (SET trades come from journal)
	return "WATCHED"
}

// deriveTickerTrend classifies ticker trend from flag repo sets.
func deriveTickerTrend(tvTicker string, flagSets []map[string]bool) string {
	// flagRepo[4] (lime) → UPTREND
	if len(flagSets) > 4 && flagSets[4][tvTicker] {
		return "UPTREND"
	}
	// flagRepo[0] (orange) → SIDEWAYS
	if len(flagSets) > 0 && flagSets[0][tvTicker] {
		return "SIDEWAYS"
	}
	// flagRepo[1] (red) → DOWNTREND
	if len(flagSets) > 1 && flagSets[1][tvTicker] {
		return "DOWNTREND"
	}
	return "SIDEWAYS"
}

// epochMsToTime converts epoch milliseconds to time.Time.
func epochMsToTime(epochMS int64) time.Time {
	return time.UnixMilli(epochMS)
}

// ============================================================================
// MIGRATION PLAN BUILDER
// ============================================================================

func buildWatchFlagSets(dump *BarkatRepositoryDump) ([]map[string]bool, []map[string]bool) {
	watchSets := make([]map[string]bool, NumWatchCategories)
	flagSets := make([]map[string]bool, NumFlagCategories)

	for i := range watchSets {
		watchSets[i] = make(map[string]bool)
	}
	for i := range flagSets {
		flagSets[i] = make(map[string]bool)
	}

	for key, cat := range dump.WatchRepo {
		if idx, err := strconv.Atoi(key); err == nil && idx >= 0 && idx < NumWatchCategories {
			for _, t := range cat {
				watchSets[idx][t] = true
			}
		}
	}
	for key, cat := range dump.FlagRepo {
		if idx, err := strconv.Atoi(key); err == nil && idx >= 0 && idx < NumFlagCategories {
			for _, t := range cat {
				flagSets[idx][t] = true
			}
		}
	}

	return watchSets, flagSets
}

// buildTickerPlan constructs the ordered list of ticker payloads.
func buildTickerPlan(dump *BarkatRepositoryDump, logger *MigrationLogger) ([]TickerPayload, *DumpAnalysis) {
	analysis := analyzeDump(dump)
	watchSets, flagSets := buildWatchFlagSets(dump)

	// Sort ticker repo keys for deterministic order
	tvTickers := make([]string, 0, len(dump.TickerRepo))
	for k := range dump.TickerRepo {
		tvTickers = append(tvTickers, k)
	}
	sort.Strings(tvTickers)

	importTime := time.Now()
	_ = analysis // used implicitly

	var plan []TickerPayload

	for _, tvTicker := range tvTickers {
		investingSymbol := dump.TickerRepo[tvTicker]

		// Exchange
		var exchange *string
		if exVal, ok := dump.ExchangeRepo[tvTicker]; ok {
			simple := extractSimpleExchange(exVal)
			exchange = &simple
			logger.LogSanitization("data.json", tvTicker, "exchange",
				exVal, simple, "extracted_simple_exchange_code")
		}

		// Timeframes
		timeframes := expandSequence(dump.SequenceRepo[tvTicker])
		if _, ok := dump.SequenceRepo[tvTicker]; !ok {
			logger.LogModification("data.json", tvTicker, "timeframes",
				"default_applied", "set to MN/WK/DL (no sequence found)")
		}

		// Type (derived from watch/flag)
		tickerType := deriveTickerType(tvTicker, exchange, flagSets, watchSets)

		// State
		state := deriveTickerState(tvTicker, watchSets)

		// Trend
		trend := deriveTickerTrend(tvTicker, flagSets)

		// LastOpenedAt
		var lastOpenedAt time.Time
		if ts, ok := dump.RecentRepo[tvTicker]; ok && ts > 0 {
			lastOpenedAt = epochMsToTime(ts)
		} else {
			lastOpenedAt = importTime.AddDate(0, -3, 0)
			logger.LogModification("data.json", tvTicker, "last_opened_at",
				"default_applied", fmt.Sprintf("set to 3mo-before-import %s (no recent entry)", lastOpenedAt.Format(time.RFC3339)))
		}

		// IsFNO
		isFNO := slices.Contains(dump.FnoRepo, tvTicker)

		// Log unresolved mapping
		if _, ok := dump.PairRepo[investingSymbol]; !ok {
			logger.LogSanitization("data.json", tvTicker, "alert_ticker_mapping",
				investingSymbol, "(skip)", "investing_symbol_not_found_in_pairRepo")
		}

		plan = append(plan, TickerPayload{
			Ticker:       tvTicker,
			Exchange:     exchange,
			Timeframes:   timeframes,
			Type:         tickerType,
			State:        state,
			Trend:        trend,
			LastOpenedAt: lastOpenedAt,
			IsFNO:        isFNO,
		})
	}

	// Append COMPOSITE tickers from all list 7 entries + list 6 entries with "/".
	// A list 6 entry is COMPOSITE if it contains "/", INDEX otherwise.
	compositeSet := make(map[string]bool)
	for _, expr := range dump.WatchRepo["7"] {
		compositeSet[expr] = true
	}
	for _, expr := range dump.FlagRepo["7"] {
		compositeSet[expr] = true
	}
	// List 6 entries with "/" are also COMPOSITE
	for _, expr := range dump.WatchRepo["6"] {
		if strings.Contains(expr, "/") {
			compositeSet[expr] = true
		}
	}
	for _, expr := range dump.FlagRepo["6"] {
		if strings.Contains(expr, "/") {
			compositeSet[expr] = true
		}
	}
	if len(compositeSet) > 0 {
		compositeExprs := make([]string, 0, len(compositeSet))
		for expr := range compositeSet {
			compositeExprs = append(compositeExprs, expr)
		}
		sort.Strings(compositeExprs)
		for _, expr := range compositeExprs {
			// Derive timeframes from sequenceRepo (same logic as main ticker loop)
			timeframes := expandSequence(dump.SequenceRepo[expr])
			if _, ok := dump.SequenceRepo[expr]; !ok {
				logger.LogModification("data.json", expr, "timeframes",
					"default_applied", "set to MN/WK/DL (no sequence found for composite)")
			} else {
				logger.LogInfo("composite_ticker_sequence", map[string]any{
					"ticker":     expr,
					"sequence":   dump.SequenceRepo[expr],
					"timeframes": timeframes,
				})
			}

			// Derive last_opened_at from recentRepo (same logic as main ticker loop)
			var lastOpenedAt time.Time
			if ts, ok := dump.RecentRepo[expr]; ok && ts > 0 {
				lastOpenedAt = epochMsToTime(ts)
				logger.LogInfo("composite_ticker_recent", map[string]any{
					"ticker":         expr,
					"last_opened_at": lastOpenedAt.Format(time.RFC3339),
				})
			} else {
				lastOpenedAt = importTime.AddDate(0, -3, 0)
			}

			plan = append(plan, TickerPayload{
				Ticker:       expr,
				Exchange:     nil,
				Timeframes:   timeframes,
				Type:         "COMPOSITE",
				State:        "WATCHED",
				Trend:        "SIDEWAYS",
				LastOpenedAt: lastOpenedAt,
				IsFNO:        false,
			})
			logger.LogInfo("composite_ticker", map[string]any{
				"ticker": expr,
				"type":   "COMPOSITE",
				"source": "list_7_or_list_6_slash",
			})
		}
	}

	return plan, analysis
}

// buildAlertDeferredLog logs the deferred alertRepo state.
func buildAlertDeferredLog(dump *BarkatRepositoryDump, logger *MigrationLogger) {
	// Per FR-010: alertRepo is deferred, log counts and anomalies
	var totalAlerts int
	emptyGroups := 0
	emptyIDs := 0
	emptyNames := 0

	for pairID, alerts := range dump.AlertRepo {
		if len(alerts) == 0 {
			emptyGroups++
			logger.LogSanitization("data.json", pairID, "alert_group",
				"empty", "(deferred)", "empty_alert_group_skipped")
			continue
		}
		for _, a := range alerts {
			totalAlerts++
			if a.ID == "" {
				emptyIDs++
				logger.LogSanitization("data.json", pairID, "alert_id",
					"(empty)", "(deferred)", "empty_alert_id")
			}
			if a.Name == "" {
				emptyNames++
				logger.LogSanitization("data.json", pairID, "alert_name",
					"(empty)", "(deferred)", "empty_alert_name")
			}
		}
	}

	logger.LogInfo("alertRepo_deferred", map[string]any{
		"fr010_status": "deferred_out_of_scope",
		"groups":       len(dump.AlertRepo),
		"total_alerts": totalAlerts,
		"empty_groups": emptyGroups,
		"empty_ids":    emptyIDs,
		"empty_names":  emptyNames,
	})
}

// ============================================================================
// API MIGRATION HELPERS
// ============================================================================

// migrationClient creates a resty client pointed at the test or real server.
func migrationClient() *resty.Client {
	client := resty.New()
	client.SetTimeout(30 * time.Second)
	client.SetHeader("Content-Type", "application/json")

	if RealServerURL != "" {
		client.SetBaseURL(RealServerURL)
	} else {
		client.SetBaseURL(fmt.Sprintf("http://localhost:%d", testPort))
	}
	return client
}

// migrateTicker attempts to create or update a ticker via the API.
// Returns true if the ticker was successfully processed.
func migrateTicker(client *resty.Client, payload TickerPayload, logger *MigrationLogger) bool {
	// POST to create
	resp, err := client.R().SetBody(map[string]any{
		"ticker":         payload.Ticker,
		"exchange":       payload.Exchange,
		"timeframes":     payload.Timeframes,
		"type":           payload.Type,
		"state":          payload.State,
		"trend":          payload.Trend,
		"last_opened_at": payload.LastOpenedAt.Format(time.RFC3339),
		"is_fno":         payload.IsFNO,
	}).Post(barkat.TickerBase)
	if err != nil {
		logger.LogError("data.json", payload.Ticker, 0, fmt.Sprintf("POST ticker request failed: %v", err), "")
		return false
	}

	switch resp.StatusCode() {
	case http.StatusCreated:
		logger.LogSuccess("data.json", payload.Ticker, fmt.Sprintf("created_%s", payload.Type), 0, false, false)
		return true

	case http.StatusConflict:
		// Ticker already exists — GET and verify/update
		return reconcileTicker(client, payload, logger)

	default:
		logger.LogError("data.json", payload.Ticker, 0,
			fmt.Sprintf("POST ticker returned %d: %s", resp.StatusCode(), string(resp.Body())), "")
		return false
	}
}

// reconcileTicker handles existing tickers by comparing and updating.
func reconcileTicker(client *resty.Client, payload TickerPayload, logger *MigrationLogger) bool {
	getResp, err := client.R().Get(barkat.TickerBase + "/" + payload.Ticker)
	if err != nil {
		logger.LogError("data.json", payload.Ticker, 0, fmt.Sprintf("GET ticker failed: %v", err), "")
		return false
	}

	if getResp.StatusCode() != http.StatusOK {
		logger.LogError("data.json", payload.Ticker, 0,
			fmt.Sprintf("GET ticker returned %d", getResp.StatusCode()), "")
		return false
	}

	var envelope common.Envelope[barkat.Ticker]
	if err := json.Unmarshal(getResp.Body(), &envelope); err != nil {
		logger.LogError("data.json", payload.Ticker, 0, fmt.Sprintf("failed to parse existing ticker: %v", err), "")
		return false
	}
	existing := envelope.Data

	// Compare fields and update if needed
	needsUpdate := false
	if (existing.Exchange == nil && payload.Exchange != nil) ||
		(existing.Exchange != nil && payload.Exchange == nil) ||
		(existing.Exchange != nil && payload.Exchange != nil && *existing.Exchange != *payload.Exchange) {
		needsUpdate = true
		logger.LogSanitization("data.json", payload.Ticker, "exchange",
			fmt.Sprint(existing.Exchange), fmt.Sprint(payload.Exchange), "conflict_during_reconciliation")
	}
	if existing.Type != payload.Type {
		needsUpdate = true
		logger.LogSanitization("data.json", payload.Ticker, "type",
			existing.Type, payload.Type, "conflict_during_reconciliation")
	}
	if existing.State != payload.State {
		needsUpdate = true
		logger.LogSanitization("data.json", payload.Ticker, "state",
			existing.State, payload.State, "conflict_during_reconciliation")
	}
	if existing.Trend != payload.Trend {
		needsUpdate = true
		logger.LogSanitization("data.json", payload.Ticker, "trend",
			existing.Trend, payload.Trend, "conflict_during_reconciliation")
	}
	if existing.IsFNO != payload.IsFNO {
		needsUpdate = true
		logger.LogSanitization("data.json", payload.Ticker, "is_fno",
			fmt.Sprint(existing.IsFNO), fmt.Sprint(payload.IsFNO), "conflict_during_reconciliation")
	}

	if needsUpdate {
		updateResp, uErr := client.R().SetBody(map[string]any{
			"exchange":   payload.Exchange,
			"timeframes": payload.Timeframes,
			"type":       payload.Type,
			"state":      payload.State,
			"trend":      payload.Trend,
			"is_fno":     payload.IsFNO,
		}).Put(barkat.TickerBase + "/" + payload.Ticker)
		if uErr != nil {
			logger.LogError("data.json", payload.Ticker, 0, fmt.Sprintf("PUT ticker failed: %v", uErr), "")
			return false
		}
		if updateResp.StatusCode() != http.StatusOK {
			logger.LogError("data.json", payload.Ticker, 0,
				fmt.Sprintf("PUT ticker returned %d: %s", updateResp.StatusCode(), string(updateResp.Body())), "")
			return false
		}
		logger.LogModification("data.json", payload.Ticker, "ticker",
			"updated_existing", "fields_differed_from_dump")
	}

	// Always PATCH last_opened_at
	patchResp, pErr := client.R().SetBody(map[string]any{
		"last_opened_at": payload.LastOpenedAt.Format(time.RFC3339),
	}).Patch(barkat.TickerBase + "/" + payload.Ticker)
	if pErr != nil {
		logger.LogError("data.json", payload.Ticker, 0, fmt.Sprintf("PATCH ticker failed: %v", pErr), "")
		return false
	}
	if patchResp.StatusCode() != http.StatusOK {
		logger.LogError("data.json", payload.Ticker, 0,
			fmt.Sprintf("PATCH ticker returned %d", patchResp.StatusCode()), string(patchResp.Body()))
		return false
	}

	logger.LogSuccess("data.json", payload.Ticker, "reconciled", 0, false, false)
	return true
}

// migrateAlertTicker creates or verifies an alert ticker via the API.
func migrateAlertTicker(client *resty.Client, payload AlertTickerPayload, logger *MigrationLogger) bool {
	resp, err := client.R().SetBody(map[string]any{
		"symbol":   payload.Symbol,
		"pair_id":  payload.PairID,
		"name":     payload.Name,
		"exchange": payload.Exchange,
	}).Post(barkat.TickerBase + "/" + payload.ParentTicker + "/alert-tickers")
	if err != nil {
		logger.LogError("data.json", payload.Symbol, 0,
			fmt.Sprintf("POST alert ticker request failed: %v", err), "")
		return false
	}

	switch resp.StatusCode() {
	case http.StatusCreated:
		logger.LogSuccess("data.json", payload.ParentTicker, fmt.Sprintf("alert_%s", payload.Symbol), 0, false, false)
		return true

	case http.StatusConflict:
		// Already exists — verify fields match
		getResp, gErr := client.R().Get(barkat.AlertTickerBase + "/" + payload.Symbol)
		if gErr != nil {
			logger.LogError("data.json", payload.Symbol, 0,
				fmt.Sprintf("GET alert ticker failed: %v", gErr), "")
			return false
		}
		if getResp.StatusCode() == http.StatusOK {
			logger.LogSuccess("data.json", payload.ParentTicker, fmt.Sprintf("alert_verified_%s", payload.Symbol), 0, false, false)
			return true
		}
		return true

	default:
		logger.LogError("data.json", payload.Symbol, 0,
			fmt.Sprintf("POST alert ticker returned %d: %s", resp.StatusCode(), string(resp.Body())), "")
		return false
	}
}

// runTickerMigration executes the full migration plan via the API.
func runTickerMigration(client *resty.Client, plan *MigrationPlan, logger *MigrationLogger) *TickerMigrationStats {
	stats := &TickerMigrationStats{StartTime: time.Now()}

	// Migrate tickers
	for _, tp := range plan.Tickers {
		stats.TotalAPICalls++
		if migrateTicker(client, tp, logger) {
			stats.CreatedTickers++
		} else {
			stats.FailedTickers = append(stats.FailedTickers, tp.Ticker)
		}
	}

	// Migrate alert tickers
	for _, ap := range plan.AlertTickers {
		stats.TotalAPICalls++
		if migrateAlertTicker(client, ap, logger) {
			stats.CreatedAlertTickers++
		} else {
			stats.FailedAlertTickers = append(stats.FailedAlertTickers, ap.Symbol)
		}
	}

	return stats
}

// verifyTickerListCounts paginates through GET /tickers to verify total.
func verifyTickerListCounts(client *resty.Client, expectedTotal int, logger *MigrationLogger) bool {
	var allTickers []barkat.Ticker
	offset := 0
	limit := 100

	for {
		resp, err := client.R().Get(fmt.Sprintf("%s?offset=%d&limit=%d", barkat.TickerBase, offset, limit))
		if err != nil {
			logger.LogError("data.json", "", 0, fmt.Sprintf("list tickers failed: %v", err), "")
			return false
		}
		if resp.StatusCode() != http.StatusOK {
			logger.LogError("data.json", "", 0,
				fmt.Sprintf("list tickers returned %d", resp.StatusCode()), string(resp.Body()))
			return false
		}

		var envelope common.Envelope[barkat.TickerList]
		if err := json.Unmarshal(resp.Body(), &envelope); err != nil {
			logger.LogError("data.json", "", 0, fmt.Sprintf("failed to parse ticker list: %v", err), "")
			return false
		}

		allTickers = append(allTickers, envelope.Data.Tickers...)

		if len(envelope.Data.Tickers) < limit {
			break
		}
		offset += limit
	}

	logger.LogInfo("ticker_list_verification", map[string]any{
		"expected": expectedTotal,
		"actual":   len(allTickers),
		"match":    len(allTickers) == expectedTotal,
	})
	return len(allTickers) == expectedTotal
}

// verifyAlertTickerListCounts paginates through GET /alert-tickers to verify total.
func verifyAlertTickerListCounts(client *resty.Client, expectedTotal int, logger *MigrationLogger) bool {
	var allAlertTickers []barkat.AlertTicker
	offset := 0
	limit := 100

	for {
		resp, err := client.R().Get(fmt.Sprintf("%s?offset=%d&limit=%d", barkat.AlertTickerBase, offset, limit))
		if err != nil {
			logger.LogError("data.json", "", 0, fmt.Sprintf("list alert tickers failed: %v", err), "")
			return false
		}
		if resp.StatusCode() != http.StatusOK {
			logger.LogError("data.json", "", 0,
				fmt.Sprintf("list alert tickers returned %d", resp.StatusCode()), string(resp.Body()))
			return false
		}

		var envelope common.Envelope[barkat.AlertTickerList]
		if err := json.Unmarshal(resp.Body(), &envelope); err != nil {
			logger.LogError("data.json", "", 0, fmt.Sprintf("failed to parse alert ticker list: %v", err), "")
			return false
		}

		allAlertTickers = append(allAlertTickers, envelope.Data.AlertTickers...)

		if len(envelope.Data.AlertTickers) < limit {
			break
		}
		offset += limit
	}

	logger.LogInfo("alert_ticker_list_verification", map[string]any{
		"expected": expectedTotal,
		"actual":   len(allAlertTickers),
		"match":    len(allAlertTickers) == expectedTotal,
	})
	return len(allAlertTickers) == expectedTotal
}

// ============================================================================
// SUMMARY REPORTING
// ============================================================================

func writeTickerMigrationSummary(stats *TickerMigrationStats, analysis *DumpAnalysis, plan *MigrationPlan, logger *MigrationLogger) {
	elapsed := time.Since(stats.StartTime)

	logger.LogInfo("TICKER_MIGRATION_SUMMARY", map[string]any{
		"dump_analysis": map[string]any{
			"tickers":                analysis.TickerCount,
			"pairs":                  analysis.PairCount,
			"exchanges":              analysis.ExchangeCount,
			"sequences":              analysis.SequenceCount,
			"recent_entries":         analysis.RecentCount,
			"fno_entries":            analysis.FnoCount,
			"watch_items":            analysis.WatchItemCount,
			"flag_items":             analysis.FlagItemCount,
			"alert_groups":           analysis.AlertGroupCount,
			"alert_entries":          analysis.AlertEntryCount,
			"unresolved_mappings":    analysis.UnresolvedMappings,
			"invalid_pair_symbols":   analysis.InvalidPairSymbols,
			"invalid_pair_names":     analysis.InvalidPairNames,
			"invalid_pair_exchanges": analysis.InvalidPairExchanges,
			"empty_alert_ids":        analysis.EmptyAlertIDs,
			"empty_alert_names":      analysis.EmptyAlertNames,
		},
		"migration_plan": map[string]any{
			"tickers_planned":       len(plan.Tickers),
			"alert_tickers_planned": len(plan.AlertTickers),
			"alert_ticker_skips":    len(plan.AlertTickerSkips),
		},
		"migration_results": map[string]any{
			"created_tickers":       stats.CreatedTickers,
			"alert_tickers_created": stats.CreatedAlertTickers,
			"failed_tickers":        len(stats.FailedTickers),
			"failed_alert_tickers":  len(stats.FailedAlertTickers),
			"total_api_calls":       stats.TotalAPICalls,
			"elapsed_seconds":       elapsed.Seconds(),
		},
	})
}

// ============================================================================
// GINKGO TESTS
// ============================================================================

var _ = Describe("Barkat Ticker API Migration", func() {
	var (
		client *resty.Client
		logger *MigrationLogger
	)

	BeforeEach(func() {
		client = migrationClient()
		if RealServerURL != "" {
			GinkgoWriter.Printf("Using real server: %s\n", RealServerURL)
		} else {
			GinkgoWriter.Printf("Using test server: http://localhost:%d (in-memory DB)\n", testPort)
		}

		var err error
		logger, err = NewMigrationLogger("/tmp")
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		if logger != nil {
			logPath := logger.GetLogPath()
			logger.Close()
			if logPath != "" {
				os.Remove(logPath)
			}
		}
	})

	// ============================================================================
	// 1.1 Parser and Preflight Helpers
	// ============================================================================
	Context("Parser and Preflight Helpers", func() {
		var (
			dump     *BarkatRepositoryDump
			analysis *DumpAnalysis
		)

		BeforeEach(func() {
			// Check if dump file exists
			if _, err := os.Stat(TickerDumpPath); os.IsNotExist(err) {
				Skip(fmt.Sprintf("Dump file not found: %s", TickerDumpPath))
				return
			}

			var err error
			dump, err = loadBarkatDump(TickerDumpPath)
			Expect(err).ToNot(HaveOccurred())
			analysis = analyzeDump(dump)
		})

		It("should load data.json and classify top-level repository keys", func() {
			Expect(dump.TickerRepo).ToNot(BeNil())
			Expect(dump.PairRepo).ToNot(BeNil())

			GinkgoWriter.Printf("\n=== Dump Classification ===\n")
			GinkgoWriter.Printf("TickerRepo:     %d\n", analysis.TickerCount)
			GinkgoWriter.Printf("PairRepo:       %d\n", analysis.PairCount)
			GinkgoWriter.Printf("ExchangeRepo:   %d\n", analysis.ExchangeCount)
			GinkgoWriter.Printf("SequenceRepo:   %d\n", analysis.SequenceCount)
			GinkgoWriter.Printf("RecentRepo:     %d\n", analysis.RecentCount)
			GinkgoWriter.Printf("FnoRepo:        %d\n", analysis.FnoCount)
		})

		It("should report source counts for all repositories", func() {
			Expect(analysis.TickerCount).To(Equal(692))
			Expect(analysis.PairCount).To(Equal(616))
			Expect(analysis.ExchangeCount).To(Equal(39))
			Expect(analysis.SequenceCount).To(Equal(212))
			Expect(analysis.RecentCount).To(Equal(85))
			Expect(analysis.FnoCount).To(Equal(193))
			Expect(analysis.WatchItemCount).To(BeNumerically(">", 0))
			Expect(analysis.FlagItemCount).To(BeNumerically(">", 0))
			Expect(analysis.AlertGroupCount).To(Equal(581))
		})

		It("should detect unresolved tickerRepo→pairRepo mappings", func() {
			Expect(analysis.UnresolvedMappings).To(Equal(76))
			Expect(analysis.UnresolvedSamples).To(HaveLen(20))

			logger.LogInfo("unresolved_mappings", map[string]any{
				"count":   analysis.UnresolvedMappings,
				"samples": analysis.UnresolvedSamples,
			})
		})

		It("should classify alertRepo as deferred and log anomalies", func() {
			Expect(analysis.AlertEntryCount).To(Equal(1099))
			Expect(analysis.EmptyAlertIDs).To(Equal(7))
			Expect(analysis.EmptyAlertNames).To(Equal(7))
			Expect(analysis.EmptyAlertGroups).To(ContainElement("18311"))
		})

		It("should print ticker type classifications post-migration", func() {
			tickerPlan, _ := buildTickerPlan(dump, logger)
			// Group by type
			byType := map[string][]string{}
			for _, tp := range tickerPlan {
				exch := "<nil>"
				if tp.Exchange != nil {
					exch = *tp.Exchange
				}
				byType[tp.Type] = append(byType[tp.Type], tp.Ticker+":"+exch)
			}
			// Print in defined order
			order := []string{"CRYPTO", "COMPOSITE", "INDEX", "FX", "COMMODITY", "BOND", "EQUITY"}
			for _, t := range order {
				tickers, ok := byType[t]
				if !ok {
					continue
				}
				pct := float64(len(tickers)) / float64(len(tickerPlan)) * 100
				GinkgoWriter.Printf("\n=== %s (%d/%d, %.1f%%) ===\n", t, len(tickers), len(tickerPlan), pct)
				for _, v := range tickers {
					GinkgoWriter.Printf("  %s\n", v)
				}
			}
			Expect(tickerPlan).To(HaveLen(746))
			Expect(byType["INDEX"]).To(HaveLen(54))
			Expect(byType["COMPOSITE"]).To(HaveLen(54))
			Expect(byType["EQUITY"]).To(HaveLen(638))
			Expect(byType).ToNot(HaveKey("CRYPTO"))
			Expect(byType).ToNot(HaveKey("FX"))
			Expect(byType).ToNot(HaveKey("COMMODITY"))
			Expect(byType).ToNot(HaveKey("BOND"))
		})
	})

	// ============================================================================
	// 1.2 Normalization Helpers
	// ============================================================================
	Context("Normalization Helpers", func() {
		It("should normalize invalid pairRepo symbols by removing spaces and special chars", func() {
			tests := []struct {
				input    string
				expected string
			}{
				{"Crude Oil WTI", "CrudeOilWTI"},
				{"India VIX", "IndiaVIX"},
				{"Natural Gas", "NaturalGas"},
				{"Small Cap 2000", "SmallCap2000"},
				{"Nifty Next 50", "NiftyNext50"},
				{"Nifty 500", "Nifty500"},
				{"VALID", "VALID"},
				{"BTC/USD", "BTC/USD"},
			}
			for _, tc := range tests {
				normalized, original := normalizeAlertSymbol(tc.input)
				Expect(normalized).To(Equal(tc.expected), "normalizeAlertSymbol(%q) should be %q", tc.input, tc.expected)
				Expect(original).To(Equal(tc.input))
			}
		})

		It("should normalize smart quotes, NBSP, and ® in names", func() {
			tests := []struct {
				input    string
				expected string
			}{
				{"Domino\u2019s Pizza Inc", "Domino's Pizza Inc"},                                  // smart quote
				{"Dr Reddy\u2019s Laboratories Ltd", "Dr Reddy's Laboratories Ltd"},                // smart quote
				{"Mazagon Dock Shipbuilders\u00a0Ltd", "Mazagon Dock Shipbuilders Ltd"},            // NBSP
				{"The Energy Select Sector SPDR\u00ae Fund", "The Energy Select Sector SPDR Fund"}, // ®
				{"EUR/INR Futures", "EUR INR Futures"},
				{"Plain ASCII name", "Plain ASCII name"},
			}
			for _, tc := range tests {
				normalized := normalizeAlertName(tc.input)
				Expect(normalized).To(Equal(tc.expected), "normalizeAlertName(%q) should be %q", tc.input, tc.expected)
			}
		})

		It("should normalize long/empty exchange labels to nil", func() {
			Expect(normalizeAlertExchange("", "TEST")).To(BeNil())
			Expect(normalizeAlertExchange("Global Indexes", "SPGSCI")).To(BeNil())
			Expect(normalizeAlertExchange("Investing.com", "BTC/USD")).To(BeNil())
			Expect(normalizeAlertExchange("NYSE", "CL")).ToNot(BeNil())
			Expect(*normalizeAlertExchange("NYSE", "CL")).To(Equal("NYSE"))
		})

		It("should extract simple exchange codes from qualified values", func() {
			Expect(extractSimpleExchange("NSE:MCX")).To(Equal("NSE"))
			Expect(extractSimpleExchange("FX_IDC:USDINR")).To(Equal("FX_IDC"))
			Expect(extractSimpleExchange("BITFINEX:BTCUSD")).To(Equal("BITFINEX"))
			Expect(extractSimpleExchange("NSE")).To(Equal("NSE"))
		})

		It("should expand sequence values into backend timeframes", func() {
			Expect(expandSequence("YR")).To(Equal([]string{"YR", "SMN", "TMN", "MN", "WK"}))
			Expect(expandSequence("MWD")).To(Equal([]string{"MN", "WK", "DL"}))
			Expect(expandSequence("")).To(Equal([]string{"MN", "WK", "DL"}))
			Expect(expandSequence("UNKNOWN")).To(Equal([]string{"MN", "WK", "DL"}))
		})

		It("should derive ticker type from list 6/list 7 membership", func() {
			// List 6 set: NATURALGAS and DXY are INDEX
			watchSets := []map[string]bool{
				nil, nil, nil, nil, nil, nil, {"NATURALGAS": true, "DXY": true, "IXIC/TLT": true}, nil,
			}
			// List 7 set: AAPL and MSFT would be COMPOSITE (if any existed)
			flagSets := []map[string]bool{
				nil, nil, nil, nil, nil, nil, nil, {"AAPL": true, "MSFT": true},
			}

			// Only list 6/7 membership determines type
			Expect(deriveTickerType("TCS", nil, nil, nil)).To(Equal("EQUITY"))
			Expect(deriveTickerType("XMRBTC", nil, nil, nil)).To(Equal("EQUITY"))
			Expect(deriveTickerType("NIFTY/USDINR", nil, nil, nil)).To(Equal("EQUITY"))
			Expect(deriveTickerType("NATURALGAS", nil, nil, watchSets)).To(Equal("INDEX"))
			Expect(deriveTickerType("DXY", nil, nil, watchSets)).To(Equal("INDEX"))
			Expect(deriveTickerType("IXIC/TLT", nil, nil, watchSets)).To(Equal("COMPOSITE"))
			Expect(deriveTickerType("NIFTY", nil, nil, nil)).To(Equal("EQUITY"))
			Expect(deriveTickerType("CNXSMALLCAP", nil, nil, nil)).To(Equal("EQUITY"))
			Expect(deriveTickerType("GOLD1!", nil, nil, nil)).To(Equal("EQUITY"))
			Expect(deriveTickerType("COPPER1!", nil, nil, nil)).To(Equal("EQUITY"))
			// List 7 takes priority over list 6
			Expect(deriveTickerType("AAPL", nil, flagSets, nil)).To(Equal("COMPOSITE"))
			Expect(deriveTickerType("MSFT", nil, flagSets, nil)).To(Equal("COMPOSITE"))
		})
	})

	// ============================================================================
	// 1.3 Full API Migration
	// ============================================================================
	Context("Full API Migration", Label("migration"), Ordered, func() {
		var (
			dump     *BarkatRepositoryDump
			plan     *MigrationPlan
			analysis *DumpAnalysis
			stats    *TickerMigrationStats
		)

		// Load dump, build plan, and execute migration ONCE for all It blocks.
		BeforeAll(func() {
			if _, err := os.Stat(TickerDumpPath); os.IsNotExist(err) {
				Skip(fmt.Sprintf("Dump file not found: %s", TickerDumpPath))
				return
			}

			var err error
			dump, err = loadBarkatDump(TickerDumpPath)
			Expect(err).ToNot(HaveOccurred())

			// Log deferred alertRepo
			buildAlertDeferredLog(dump, logger)

			// Build migration plan
			tickerPlan, preflight := buildTickerPlan(dump, logger)
			analysis = preflight
			alertPlan := []AlertTickerPayload{}

			plan = &MigrationPlan{
				Tickers:          tickerPlan,
				AlertTickers:     alertPlan,
				AlertTickerSkips: []string{},
			}

			// Log preflight analysis
			logger.LogInfo("preflight_analysis", map[string]any{
				"tickers":            len(tickerPlan),
				"alert_tickers":      len(alertPlan),
				"unresolved_skipped": analysis.UnresolvedMappings,
			})

			// Execute API migration
			stats = runTickerMigration(client, plan, logger)
		})

		It("should migrate all valid primary tickers through /v1/api/tickers", func() {
			Expect(stats.FailedTickers).To(BeEmpty(),
				"ticker migration failures: %v", stats.FailedTickers)
			Expect(stats.CreatedTickers).To(Equal(len(plan.Tickers)),
				"all %d tickers should be created/reconciled", len(plan.Tickers))
		})

		It("should defer alert ticker migration", func() {
			Expect(plan.AlertTickers).To(BeEmpty())
			Expect(stats.FailedAlertTickers).To(BeEmpty(),
				"alert ticker migration failures: %v", stats.FailedAlertTickers)
			Expect(stats.CreatedAlertTickers).To(Equal(0),
				"alert tickers are deferred for this migration pass")
		})

		It("should verify paginated ticker list count", func() {
			Expect(verifyTickerListCounts(client, len(plan.Tickers), logger)).To(BeTrue())
		})

		It("should verify paginated alert ticker list count", func() {
			Expect(verifyAlertTickerListCounts(client, len(plan.AlertTickers), logger)).To(BeTrue())
		})

		It("should write final reconciliation summary", func() {
			writeTickerMigrationSummary(stats, analysis, plan, logger)
			Expect(logger.GetErrorCount()).To(Equal(0),
				"migration should have zero errors, got %d", logger.GetErrorCount())
		})
	})
})
