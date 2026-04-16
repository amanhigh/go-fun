package core_test

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
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
	// ProcessedFolder is the base path for journal files
	ProcessedFolder = "/home/aman/Projects/go-fun/processed"

	// TestFile is the specific test file used for single file migration
	TestFile = "2023_06_15.md"
)

var (
	// RealServerURL - set BARKAT_REAL_SERVER_URL to test against live database.
	// If empty, tests run against the in-memory test server.
	RealServerURL = os.Getenv("BARKAT_REAL_SERVER_URL")
)

var (
	entryPattern                     = regexp.MustCompile(`\|\s*` + "`" + `([A-Z0-9_!]+)` + "`" + `\s*\|(.+)\|`)
	imagePattern                     = regexp.MustCompile(`!\[.*?\]\(([^)]+)\)`)
	tagPattern                       = regexp.MustCompile(`#([trm])\.([a-z0-9-]+)`)
	fileNameSanitizer                = regexp.MustCompile(`[!@#$%^&*()+=\[\]{}|;:'\",<>?/\\]`)
	imageDashDatePattern             = regexp.MustCompile(`--(\d{4}-\d{2}-\d{2})-\d+\.(?i:jpg|jpeg|png)$`)
	imageDoubleUnderscoreDatePattern = regexp.MustCompile(`__(\d{8})__\d{6}\.(?i:jpg|jpeg|png)$`)
	imageSingleUnderscoreDatePattern = regexp.MustCompile(`_(\d{8})_\d{6}\.(?i:jpg|jpeg|png)$`)
	statusTagMap                     = map[string]string{
		"fail":     "FAIL",
		"success":  "SUCCESS",
		"running":  "RUNNING",
		"broken":   "BROKEN",
		"missed":   "MISSED",
		"miss":     "MISSED",
		"justloss": "JUST_LOSS",
	}
)

// ============================================================================
// LOGGING INFRASTRUCTURE
// ============================================================================

// MigrationLogger handles all logging for migration process
type MigrationLogger struct {
	logFile       *os.File
	logPath       string
	sanitizations []SanitizationLog
	modifications []ModificationLog
	errors        []ErrorLog
	successes     []SuccessLog
}

type SanitizationLog struct {
	File     string `json:"file"`
	Ticker   string `json:"ticker"`
	Field    string `json:"field"`
	OldValue string `json:"old_value"`
	NewValue string `json:"new_value"`
	Reason   string `json:"reason"`
}

type ModificationLog struct {
	File    string `json:"file"`
	Ticker  string `json:"ticker"`
	Field   string `json:"field"`
	Action  string `json:"action"`
	Details string `json:"details"`
}

type ErrorLog struct {
	File    string `json:"file"`
	Ticker  string `json:"ticker"`
	Line    int    `json:"line,omitempty"`
	Error   string `json:"error"`
	RawData string `json:"raw_data,omitempty"`
}

type SuccessLog struct {
	File       string `json:"file"`
	Ticker     string `json:"ticker"`
	JournalID  string `json:"journal_id"`
	ImageCount int    `json:"image_count"`
	HasNote    bool   `json:"has_note"`
	HasTag     bool   `json:"has_tag"`
}

// ProcessingCounts for migration tracking
type ProcessingCounts struct {
	// From parsing
	ParsedTickers int `json:"parsed_tickers"`
	ParsedImages  int `json:"parsed_images"`
	ParsedNotes   int `json:"parsed_notes"`

	// From API migration
	MigratedJournals int `json:"migrated_journals"`
	MigratedImages   int `json:"migrated_images"`
	MigratedNotes    int `json:"migrated_notes"`
	MigratedTags     int `json:"migrated_tags"`
}

type ImageBuildStats struct {
	SourceRefs             int `json:"source_refs"`
	MigratedSourceRefs     int `json:"migrated_source_refs"`
	SkippedClipboardRefs   int `json:"skipped_clipboard_refs"`
	TruncatedSourceRefs    int `json:"truncated_source_refs"`
	PlaceholderImagesAdded int `json:"placeholder_images_added"`
	FinalImageCount        int `json:"final_image_count"`
}

type ConvertedJournal struct {
	Journal    barkat.Journal
	ImageStats ImageBuildStats
}

type MigrationAccounting struct {
	SourceEntries            int
	SourceImageRefs          int
	MigratedSourceImageRefs  int
	SkippedClipboardRefs     int
	TruncatedSourceImageRefs int
	PlaceholderImagesAdded   int
	FinalImageRecords        int
	ProjectedNotes           int
	ProjectedTags            int
}

func (a *MigrationAccounting) RecordConversion(converted ConvertedJournal) {
	a.SourceEntries++
	a.SourceImageRefs += converted.ImageStats.SourceRefs
	a.MigratedSourceImageRefs += converted.ImageStats.MigratedSourceRefs
	a.SkippedClipboardRefs += converted.ImageStats.SkippedClipboardRefs
	a.TruncatedSourceImageRefs += converted.ImageStats.TruncatedSourceRefs
	a.PlaceholderImagesAdded += converted.ImageStats.PlaceholderImagesAdded
	a.FinalImageRecords += converted.ImageStats.FinalImageCount
	a.ProjectedNotes += len(converted.Journal.Notes)
	a.ProjectedTags += len(converted.Journal.Tags)
}

func (a MigrationAccounting) VerifySourceImageAccounting() {
	Expect(a.MigratedSourceImageRefs+a.SkippedClipboardRefs+a.TruncatedSourceImageRefs).To(Equal(a.SourceImageRefs),
		"source image refs should reconcile to migrated + skipped + truncated refs")
}

func (a MigrationAccounting) VerifyMigratedCounts(migratedCounts ProcessingCounts) {
	a.VerifySourceImageAccounting()
	Expect(migratedCounts.MigratedJournals).To(Equal(a.SourceEntries), "every parsed entry should migrate")
	Expect(migratedCounts.MigratedImages).To(Equal(a.FinalImageRecords), "persisted image records should match projection")
	Expect(migratedCounts.MigratedNotes).To(Equal(a.ProjectedNotes), "persisted note records should match projection")
	Expect(migratedCounts.MigratedTags).To(Equal(a.ProjectedTags), "persisted tag records should match projection")
}

func NewMigrationLogger(logDir string) (*MigrationLogger, error) {
	timestamp := time.Now().Format("20060102_150405")
	logPath := filepath.Join(logDir, fmt.Sprintf("migration_%s.log", timestamp))

	logFile, err := os.Create(logPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create log file: %w", err)
	}

	return &MigrationLogger{
		logFile: logFile,
		logPath: logPath,
	}, nil
}

func (l *MigrationLogger) LogSanitization(file, ticker, field, oldVal, newVal, reason string) {
	entry := SanitizationLog{
		File:     filepath.Base(file),
		Ticker:   ticker,
		Field:    field,
		OldValue: oldVal,
		NewValue: newVal,
		Reason:   reason,
	}
	l.sanitizations = append(l.sanitizations, entry)
	l.writeJSON("SANITIZATION", entry)
}

func (l *MigrationLogger) LogModification(file, ticker, field, action, details string) {
	entry := ModificationLog{
		File:    filepath.Base(file),
		Ticker:  ticker,
		Field:   field,
		Action:  action,
		Details: details,
	}
	l.modifications = append(l.modifications, entry)
	l.writeJSON("MODIFICATION", entry)
}

func (l *MigrationLogger) LogError(file, ticker string, line int, errMsg, rawData string) {
	entry := ErrorLog{
		File:    filepath.Base(file),
		Ticker:  ticker,
		Line:    line,
		Error:   errMsg,
		RawData: rawData,
	}
	l.errors = append(l.errors, entry)
	l.writeJSON("ERROR", entry)
}

func (l *MigrationLogger) LogSuccess(file, ticker, journalID string, imageCount int, hasNote, hasTag bool) {
	entry := SuccessLog{
		File:       filepath.Base(file),
		Ticker:     ticker,
		JournalID:  journalID,
		ImageCount: imageCount,
		HasNote:    hasNote,
		HasTag:     hasTag,
	}
	l.successes = append(l.successes, entry)
	l.writeJSON("SUCCESS", entry)
}

func (l *MigrationLogger) LogInfo(message string, data any) {
	l.writeJSON("INFO", map[string]any{"message": message, "data": data})
}

func (l *MigrationLogger) writeJSON(logType string, data any) {
	jsonData, _ := json.Marshal(data)
	fmt.Fprintf(l.logFile, "%s|%s|%s\n", time.Now().Format(time.RFC3339), logType, string(jsonData))
}

func (l *MigrationLogger) WriteSummary(processed ProcessingCounts, files int) {
	summary := map[string]any{
		"files_processed": files,
		"processing": map[string]int{
			"parsed_tickers":    processed.ParsedTickers,
			"parsed_images":     processed.ParsedImages,
			"parsed_notes":      processed.ParsedNotes,
			"migrated_journals": processed.MigratedJournals,
			"migrated_images":   processed.MigratedImages,
			"migrated_notes":    processed.MigratedNotes,
			"migrated_tags":     processed.MigratedTags,
		},
		"totals": map[string]int{
			"sanitizations": len(l.sanitizations),
			"modifications": len(l.modifications),
			"errors":        len(l.errors),
			"successes":     len(l.successes),
		},
	}
	l.writeJSON("SUMMARY", summary)
}

func (l *MigrationLogger) Close() {
	if l.logFile != nil {
		l.logFile.Close()
	}
}

func (l *MigrationLogger) GetLogPath() string {
	return l.logPath
}

func (l *MigrationLogger) GetErrorCount() int {
	return len(l.errors)
}

func (l *MigrationLogger) GetSanitizationCount() int {
	return len(l.sanitizations)
}

func (l *MigrationLogger) GetModificationCount() int {
	return len(l.modifications)
}

// ============================================================================
// LEGACY ENTRY STRUCTURES
// ============================================================================

// LegacyJournalEntry represents a parsed entry from Logseq markdown
// Stores raw values as-is from markdown, conversion happens separately
type LegacyJournalEntry struct {
	Ticker         string
	RawTags        []string // Store all raw tags for logging
	Sequence       string
	Type           string
	Status         string
	Direction      string   // trend or ctrend -> maps to DIRECTION tag
	ReasonTags     []string // All #r.* tags with their full values (e.g., "dep-loc")
	ManagementTags []string // All #m.* tags (e.g., "ntr", "enl")
	IsImportant    bool     // #important tag present
	Images         []string
	Note           string   // Code block note (Plan notes)
	SimpleNotes    []string // Simple notes outside code blocks (review comments)
	RawLine        string   // Original line for reference
	LineNumber     int      // Line number in file
	RawMarkdown    string   // Complete raw markdown block for this entry
}

// ============================================================================
// DATE EXTRACTION
// ============================================================================

// extractDateFromFilename extracts date from filename like "2023_06_15.md"
func extractDateFromFilename(filename string) (time.Time, error) {
	base := filepath.Base(filename)
	base = strings.TrimSuffix(base, ".md")

	for _, format := range []string{"2006_01_02", "2006-01-02"} {
		t, err := time.Parse(format, base)
		if err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("cannot parse date from filename: %s", filename)
}

// ============================================================================
// PARSING WITH LOGGING
// ============================================================================

// parseLegacyMarkdownWithLogging parses markdown and captures raw content
func parseLegacyMarkdownWithLogging(filePath string, _ *MigrationLogger) ([]LegacyJournalEntry, ProcessingCounts, error) {
	var counts ProcessingCounts

	file, err := os.Open(filePath)
	if err != nil {
		return nil, counts, err
	}
	defer file.Close()

	var entries []LegacyJournalEntry
	var inCodeBlock bool
	var noteContent strings.Builder
	var rawMarkdown strings.Builder
	var lineNumber int

	scanner := bufio.NewScanner(file)

	result := processMarkdownLines(markdownProcessContext{
		scanner:      scanner,
		lineNumber:   lineNumber,
		currentEntry: nil,
		entries:      entries,
		counts:       counts,
		inCodeBlock:  inCodeBlock,
		rawMarkdown:  &rawMarkdown,
		noteContent:  &noteContent,
		entryPattern: entryPattern,
		imagePattern: imagePattern,
	})
	entries = result.entries
	counts = result.counts

	return entries, counts, scanner.Err()
}

// markdownProcessResult holds the result of processing markdown lines
type markdownProcessResult struct {
	entries      []LegacyJournalEntry
	counts       ProcessingCounts
	currentEntry *LegacyJournalEntry
	inCodeBlock  bool
}

// markdownProcessContext holds parameters for markdown processing
type markdownProcessContext struct {
	scanner      *bufio.Scanner
	lineNumber   int
	currentEntry *LegacyJournalEntry
	entries      []LegacyJournalEntry
	counts       ProcessingCounts
	inCodeBlock  bool
	rawMarkdown  *strings.Builder
	noteContent  *strings.Builder
	entryPattern *regexp.Regexp
	imagePattern *regexp.Regexp
}

// processMarkdownLines processes all lines in the markdown file
func processMarkdownLines(ctx markdownProcessContext) markdownProcessResult {
	for ctx.scanner.Scan() {
		line := ctx.scanner.Text()
		ctx.lineNumber++

		result := processLine(lineProcessContext{
			line:         line,
			lineNumber:   ctx.lineNumber,
			currentEntry: ctx.currentEntry,
			entries:      ctx.entries,
			counts:       ctx.counts,
			inCodeBlock:  ctx.inCodeBlock,
			rawMarkdown:  ctx.rawMarkdown,
			noteContent:  ctx.noteContent,
			entryPattern: ctx.entryPattern,
			imagePattern: ctx.imagePattern,
		})
		ctx.currentEntry = result.currentEntry
		ctx.entries = result.entries
		ctx.counts = result.counts
		ctx.inCodeBlock = result.inCodeBlock
	}

	// Don't forget the last entry
	if ctx.currentEntry != nil {
		ctx.currentEntry.RawMarkdown = ctx.rawMarkdown.String()
		ctx.entries = append(ctx.entries, *ctx.currentEntry)
	}

	return markdownProcessResult{ctx.entries, ctx.counts, ctx.currentEntry, ctx.inCodeBlock}
}

// lineProcessResult holds the result of processing a line
type lineProcessResult struct {
	currentEntry *LegacyJournalEntry
	entries      []LegacyJournalEntry
	counts       ProcessingCounts
	inCodeBlock  bool
}

// lineProcessContext holds parameters for line processing
type lineProcessContext struct {
	line         string
	lineNumber   int
	currentEntry *LegacyJournalEntry
	entries      []LegacyJournalEntry
	counts       ProcessingCounts
	inCodeBlock  bool
	rawMarkdown  *strings.Builder
	noteContent  *strings.Builder
	entryPattern *regexp.Regexp
	imagePattern *regexp.Regexp
}

// processLine handles individual line processing during parsing
func processLine(ctx lineProcessContext) lineProcessResult {
	// Check for entry line (table row with ticker)
	if matches := ctx.entryPattern.FindStringSubmatch(ctx.line); matches != nil {
		return handleEntryLine(ctx, matches)
	}

	// Track raw markdown for current entry
	updateRawMarkdown(ctx)

	// Handle code blocks for notes
	if strings.Contains(ctx.line, "```") {
		return handleCodeBlockLine(ctx)
	}

	if ctx.inCodeBlock {
		return handleInCodeBlockLine(ctx)
	}

	// Check for image line
	return handleImageLine(ctx)
}

// updateRawMarkdown tracks raw markdown for current entry
func updateRawMarkdown(ctx lineProcessContext) {
	if ctx.currentEntry != nil {
		ctx.rawMarkdown.WriteString(ctx.line + "\n")
	}
}

// handleCodeBlockLine processes code block markers
func handleCodeBlockLine(ctx lineProcessContext) lineProcessResult {
	ctx.currentEntry, ctx.counts = handleCodeBlock(ctx.line, ctx.inCodeBlock, ctx.currentEntry, ctx.noteContent, ctx.counts)
	ctx.inCodeBlock = !ctx.inCodeBlock
	return lineProcessResult{ctx.currentEntry, ctx.entries, ctx.counts, ctx.inCodeBlock}
}

// handleInCodeBlockLine collects note content inside code blocks
func handleInCodeBlockLine(ctx lineProcessContext) lineProcessResult {
	if ctx.currentEntry != nil {
		ctx.noteContent.WriteString(ctx.line + "\n")
	}
	return lineProcessResult{ctx.currentEntry, ctx.entries, ctx.counts, ctx.inCodeBlock}
}

// handleEntryLine processes journal entry lines (table rows)
func handleEntryLine(ctx lineProcessContext, matches []string) lineProcessResult {
	// Save previous entry if exists
	if ctx.currentEntry != nil {
		ctx.currentEntry.RawMarkdown = ctx.rawMarkdown.String()
		ctx.entries = append(ctx.entries, *ctx.currentEntry)
	}

	ctx.currentEntry = parseEntryLine(matches, ctx.line, ctx.lineNumber)
	ctx.counts.ParsedTickers++

	// Reset raw markdown for new entry
	ctx.rawMarkdown.Reset()
	ctx.noteContent.Reset()
	ctx.rawMarkdown.WriteString(ctx.line + "\n")

	// Parse and store raw tags
	ctx.currentEntry = parseLegacyTags(matches[2], ctx.currentEntry)
	return lineProcessResult{ctx.currentEntry, ctx.entries, ctx.counts, ctx.inCodeBlock}
}

// handleImageLine processes image lines and simple notes
func handleImageLine(ctx lineProcessContext) lineProcessResult {
	if ctx.currentEntry != nil {
		// Check for image
		if matches := ctx.imagePattern.FindStringSubmatch(ctx.line); matches != nil {
			ctx.currentEntry.Images = append(ctx.currentEntry.Images, matches[1])
			ctx.counts.ParsedImages++
			return lineProcessResult{ctx.currentEntry, ctx.entries, ctx.counts, ctx.inCodeBlock}
		}

		// Check for simple note (starts with - but not an image, not empty)
		lineStripped := strings.TrimSpace(ctx.line)
		if after, ok := strings.CutPrefix(lineStripped, "-"); ok {
			noteContent := strings.TrimSpace(after)
			// Skip empty dashes and logseq properties
			if noteContent != "" && !strings.Contains(noteContent, "::") {
				// This is a simple note - capture it!
				ctx.currentEntry.SimpleNotes = append(ctx.currentEntry.SimpleNotes, noteContent)
				ctx.counts.ParsedNotes++
			}
		}
	}
	return lineProcessResult{ctx.currentEntry, ctx.entries, ctx.counts, ctx.inCodeBlock}
}

// handleCodeBlock processes code block start/end markers
func handleCodeBlock(_ string, inCodeBlock bool, currentEntry *LegacyJournalEntry, noteContent *strings.Builder, counts ProcessingCounts) (*LegacyJournalEntry, ProcessingCounts) {
	if inCodeBlock {
		// End of code block
		if currentEntry != nil {
			currentEntry.Note = strings.TrimSpace(noteContent.String())
			if currentEntry.Note != "" {
				counts.ParsedNotes++
			}
		}
		noteContent.Reset()
	}
	return currentEntry, counts
}

// parseEntryLine creates a basic LegacyJournalEntry from a table row
func parseEntryLine(matches []string, line string, lineNumber int) *LegacyJournalEntry {
	ticker := matches[1]
	return &LegacyJournalEntry{
		Ticker:     ticker,
		RawLine:    line,
		LineNumber: lineNumber,
	}
}

// parseLegacyTags extracts and categorizes legacy tags from tag string
func parseLegacyTags(tagsPart string, entry *LegacyJournalEntry) *LegacyJournalEntry {
	tags := tagPattern.FindAllStringSubmatch(tagsPart, -1)

	for _, tag := range tags {
		entry.RawTags = append(entry.RawTags, fmt.Sprintf("#%s.%s", tag[1], tag[2]))
		prefix := tag[1]
		value := tag[2]

		switch prefix {
		case "t":
			entry = processTradeTag(value, entry)
		case "r":
			// Store full reason tag value (e.g., "dep-loc", "nca-egf")
			entry.ReasonTags = append(entry.ReasonTags, value)
		case "m":
			// Store management tags (e.g., "ntr", "enl", "slt")
			entry.ManagementTags = append(entry.ManagementTags, value)
		}
	}

	// Check for #important tag (PRD 4.8.6.3 - must be captured)
	if strings.Contains(tagsPart, "#important") {
		entry.IsImportant = true
	}

	return entry
}

// processTradeTag handles trade tag mappings
func processTradeTag(value string, entry *LegacyJournalEntry) *LegacyJournalEntry {
	switch value {
	case "mwd", "yr":
		entry.Sequence = strings.ToUpper(value)
	case "wdh":
		entry.Sequence = "WDH"
	case "rejected":
		entry.Type = chooseTradeType(entry.Type, "REJECTED")
	case "taken":
		entry.Type = chooseTradeType(entry.Type, "TAKEN")
	case "set":
		entry.Status = "SET"
	case "fail", "success", "running", "broken", "missed", "miss", "justloss":
		entry.Status = statusTagMap[value]
	case "trend", "ctrend":
		entry.Direction = value
	}
	return entry
}

func chooseTradeType(current, candidate string) string {
	priority := map[string]int{
		"":         0,
		"REJECTED": 1,
		"TAKEN":    2,
	}
	if priority[candidate] >= priority[current] {
		return candidate
	}
	return current
}

// ============================================================================
// CONVERSION WITH LOGGING
// ============================================================================

// sanitizeFileNameWithLogging removes invalid characters and logs changes
func sanitizeFileNameWithLogging(name, file, ticker string, logger *MigrationLogger) string {
	sanitized := fileNameSanitizer.ReplaceAllString(name, "_")

	if sanitized != name {
		logger.LogSanitization(file, ticker, "filename", name, sanitized, "invalid_characters_removed")
	}

	return sanitized
}

// sanitizeTickerWithLogging sanitizes ticker symbols and logs changes
func sanitizeTickerWithLogging(ticker, file string, logger *MigrationLogger) string {
	original := ticker

	// Remove trailing ! (futures symbols)
	if before, ok := strings.CutSuffix(ticker, "!"); ok {
		ticker = before
		logger.LogSanitization(file, original, "ticker", original, ticker, "removed_trailing_exclamation")
	}

	// Replace underscores with nothing (NSE tickers don't have underscores)
	// M_MFIN -> MMFIN, MCDOWELL_N -> MCDOWELLN, M_M -> MM
	if strings.Contains(ticker, "_") {
		newTicker := strings.ReplaceAll(ticker, "_", "")
		logger.LogSanitization(file, original, "ticker", ticker, newTicker, "removed_underscores")
		ticker = newTicker
	}

	return ticker
}

// convertToJournalWithLogging converts legacy entry with full logging
func convertToJournalWithLogging(entry LegacyJournalEntry, journalDate time.Time, filePath string, logger *MigrationLogger) ConvertedJournal {
	ticker := sanitizeTickerWithLogging(entry.Ticker, filePath, logger)

	// Build images with logging
	images, imageStats := buildImagesWithLogging(entry.Images, ticker, journalDate, filePath, logger)

	// Build tags per PRD 4.8.6.3
	tags := buildTagsFromLegacy(entry, journalDate)

	// Build notes - include original markdown as first note
	notes := buildNotesFromLegacy(entry, journalDate)

	// Handle defaults with logging
	sequence, status, journalType := applyDefaultsWithLogging(entry, ticker, filePath, logger)

	// Log image count changes
	logImageCountChanges(imageStats.SourceRefs, len(entry.Images), len(images), ticker, filePath, logger)

	return ConvertedJournal{
		Journal: barkat.Journal{
			Ticker:    ticker,
			Sequence:  sequence,
			Type:      journalType,
			Status:    status,
			CreatedAt: journalDate,
			Images:    images,
			Tags:      tags,
			Notes:     notes,
		},
		ImageStats: imageStats,
	}
}

// extractImageDateFromFilename extracts the date from an image filename.
// Returns journalDate if no date is found (fallback).
func extractImageDateFromFilename(basename string, journalDate time.Time) time.Time {
	for _, candidate := range []struct {
		pattern *regexp.Regexp
		format  string
	}{
		{pattern: imageDashDatePattern, format: "2006-01-02"},
		{pattern: imageDoubleUnderscoreDatePattern, format: "20060102"},
		{pattern: imageSingleUnderscoreDatePattern, format: "20060102"},
	} {
		matches := candidate.pattern.FindStringSubmatch(basename)
		if len(matches) < 2 {
			continue
		}
		parsed, err := time.Parse(candidate.format, matches[1])
		if err == nil {
			return parsed
		}
	}

	return journalDate
}

func finalizeImagesWithLogging(images []barkat.Image, ticker, filePath string, logger *MigrationLogger, journalDate time.Time, stats ImageBuildStats) ([]barkat.Image, ImageBuildStats) {
	if len(images) < 4 {
		logger.LogModification(filePath, ticker, "images", "placeholder_added",
			fmt.Sprintf("added %d placeholder images (had %d valid, need 4)", 4-len(images), len(images)))
	}
	for len(images) < 4 {
		images = append(images, barkat.Image{
			Timeframe: []string{"DL", "WK", "MN", "TMN"}[len(images)],
			FileName:  fmt.Sprintf("placeholder_%d.png", len(images)),
			CreatedAt: journalDate,
		})
		stats.PlaceholderImagesAdded++
	}

	stats.FinalImageCount = len(images)
	return images, stats
}

// buildImagesWithLogging creates and validates image list
func buildImagesWithLogging(imagePaths []string, ticker string, journalDate time.Time, filePath string, logger *MigrationLogger) ([]barkat.Image, ImageBuildStats) {
	images := make([]barkat.Image, 0, len(imagePaths))
	timeframes := []string{"DL", "WK", "MN", "TMN"}
	stats := ImageBuildStats{SourceRefs: len(imagePaths)}

	for _, imgPath := range imagePaths {
		basename := filepath.Base(imgPath)

		// Skip Logseq clipboard paste images — they reference files that were never on disk.
		// These have names like "image_XXXXXXXXXXXXX_X.png" and are stored in Logseq's
		// assets folder, not in the trading image directory. Creating DB records for them
		// produces broken image references in the UI ("Cannot read clipboard" error).
		if strings.HasPrefix(basename, "image_") {
			logger.LogModification(filePath, ticker, "images", "skipped_clipboard_paste",
				fmt.Sprintf("skipped clipboard paste image: %s (file never existed on disk)", basename))
			stats.SkippedClipboardRefs++
			continue
		}

		if len(images) >= 16 {
			stats.TruncatedSourceRefs++
			continue
		}

		timeframe := timeframes[len(images)%len(timeframes)]
		sanitizedName := sanitizeFileNameWithLogging(basename, filePath, ticker, logger)
		imageDate := extractImageDateFromFilename(basename, journalDate)
		images = append(images, barkat.Image{
			Timeframe: timeframe,
			FileName:  sanitizedName,
			CreatedAt: imageDate,
		})
		stats.MigratedSourceRefs++
	}

	if stats.TruncatedSourceRefs > 0 {
		logger.LogModification(filePath, ticker, "images", "truncated",
			fmt.Sprintf("skipped %d source images beyond 16-image limit", stats.TruncatedSourceRefs))
	}

	return finalizeImagesWithLogging(images, ticker, filePath, logger, journalDate, stats)
}

func verifyMigratedJournalMatchesProjection(actual, expected barkat.Journal) {
	Expect(actual.Ticker).To(Equal(expected.Ticker))
	Expect(actual.Sequence).To(Equal(expected.Sequence))
	Expect(actual.Type).To(Equal(expected.Type))
	Expect(actual.Status).To(Equal(expected.Status))
	Expect(actual.CreatedAt.Format(time.DateOnly)).To(Equal(expected.CreatedAt.Format(time.DateOnly)))
	Expect(actual.Images).To(HaveLen(len(expected.Images)))
	Expect(actual.Tags).To(HaveLen(len(expected.Tags)))
	Expect(actual.Notes).To(HaveLen(len(expected.Notes)))
	verifyProjectedImages(actual.Images, expected.Images)
	verifyProjectedTags(actual.Tags, expected.Tags)
	verifyProjectedNotes(actual.Notes, expected.Notes)
}

func verifyProjectedImages(actualImages, expectedImages []barkat.Image) {
	for i, expectedImage := range expectedImages {
		actualImage := actualImages[i]
		Expect(actualImage.ExternalID).To(HavePrefix("img_"))
		Expect(actualImage.Timeframe).To(Equal(expectedImage.Timeframe))
		Expect(actualImage.FileName).To(Equal(expectedImage.FileName))
		Expect(actualImage.CreatedAt.Format(time.DateOnly)).To(Equal(expectedImage.CreatedAt.Format(time.DateOnly)))
	}
}

func verifyProjectedTags(actualTags, expectedTags []barkat.Tag) {
	for i, expectedTag := range expectedTags {
		actualTag := actualTags[i]
		Expect(actualTag.ExternalID).To(HavePrefix("tag_"))
		Expect(actualTag.Tag).To(Equal(expectedTag.Tag))
		Expect(actualTag.Type).To(Equal(expectedTag.Type))
		if expectedTag.Override == nil {
			Expect(actualTag.Override).To(BeNil())
		} else {
			Expect(actualTag.Override).ToNot(BeNil())
			Expect(*actualTag.Override).To(Equal(*expectedTag.Override))
		}
	}
}

func verifyProjectedNotes(actualNotes, expectedNotes []barkat.Note) {
	for i, expectedNote := range expectedNotes {
		actualNote := actualNotes[i]
		Expect(actualNote.ExternalID).To(HavePrefix("not_"))
		Expect(actualNote.Status).To(Equal(expectedNote.Status))
		Expect(actualNote.Format).To(Equal(expectedNote.Format))
		Expect(actualNote.Content).To(Equal(expectedNote.Content))
	}
}

// buildTagsFromLegacy converts legacy tags to barkat.Tags
func buildTagsFromLegacy(entry LegacyJournalEntry, journalDate time.Time) []barkat.Tag {
	var tags []barkat.Tag

	// Add direction tag (trend/ctrend -> DIRECTION type)
	if entry.Direction != "" {
		tags = append(tags, buildDirectionTag(entry.Direction, journalDate))
	}

	// Add reason tags (PRD 4.8.6.3.2)
	tags = append(tags, buildReasonTags(entry.ReasonTags, journalDate)...)

	// Add management tags (PRD 4.8.6.3.3)
	tags = append(tags, buildManagementTags(entry.ManagementTags, journalDate)...)

	// Add #important tag if present (PRD 4.8.6.3 - must be captured)
	if entry.IsImportant {
		tags = append(tags, barkat.Tag{Tag: "important", Type: "MANAGEMENT", CreatedAt: journalDate})
	}

	return tags
}

// buildDirectionTag creates a DIRECTION tag
func buildDirectionTag(direction string, journalDate time.Time) barkat.Tag {
	return barkat.Tag{Tag: direction, Type: "DIRECTION", CreatedAt: journalDate}
}

// buildReasonTags creates REASON tags with optional overrides
func buildReasonTags(reasonTags []string, journalDate time.Time) []barkat.Tag {
	var tags []barkat.Tag
	for _, reasonTag := range reasonTags {
		parts := strings.SplitN(reasonTag, "-", 2)
		tag := barkat.Tag{Tag: parts[0], Type: "REASON", CreatedAt: journalDate}
		if len(parts) > 1 {
			override := parts[1]
			tag.Override = &override
		}
		tags = append(tags, tag)
	}
	return tags
}

// buildManagementTags creates MANAGEMENT tags
func buildManagementTags(mgmtTags []string, journalDate time.Time) []barkat.Tag {
	var tags []barkat.Tag
	for _, mgmtTag := range mgmtTags {
		tags = append(tags, barkat.Tag{Tag: mgmtTag, Type: "MANAGEMENT", CreatedAt: journalDate})
	}
	return tags
}

// buildNotesFromLegacy creates notes from legacy entry
// Plan notes (code blocks) have status SET (when trade is set)
// Simple notes (review comments) have FINAL status (success/fail outcome)
// Note: Model allows max=1 note, so we combine all content into one note
func buildNotesFromLegacy(entry LegacyJournalEntry, journalDate time.Time) []barkat.Note {
	sections := buildNoteSections(entry)
	if len(sections) == 0 {
		return nil
	}

	return []barkat.Note{{
		Status:    finalStatusForLegacyNote(entry),
		Content:   strings.Join(sections, "\n\n"),
		Format:    "MARKDOWN",
		CreatedAt: journalDate,
	}}
}

func buildNoteSections(entry LegacyJournalEntry) []string {
	rawMarkdown := strings.TrimSpace(entry.RawMarkdown)
	planNote := strings.TrimSpace(entry.Note)
	sections := buildPrimaryNoteSections(rawMarkdown, planNote)

	if reviewSection := buildReviewNoteSection(entry.SimpleNotes); reviewSection != "" {
		sections = append(sections, reviewSection)
	}

	return sections
}

func buildPrimaryNoteSections(rawMarkdown, planNote string) []string {
	if rawMarkdown != "" {
		return []string{"=== ORIGINAL MARKDOWN ===\n" + rawMarkdown}
	}
	if planNote != "" {
		return []string{"=== PLAN NOTES ===\n" + planNote}
	}
	return nil
}

func buildReviewNoteSection(simpleNotes []string) string {
	if len(simpleNotes) == 0 {
		return ""
	}

	var reviewBuilder strings.Builder
	reviewBuilder.WriteString("=== REVIEW NOTES ===\n")
	for _, simpleNote := range simpleNotes {
		reviewBuilder.WriteString("- ")
		reviewBuilder.WriteString(simpleNote)
		reviewBuilder.WriteString("\n")
	}

	return strings.TrimSpace(reviewBuilder.String())
}

func finalStatusForLegacyNote(entry LegacyJournalEntry) string {
	if entry.Status != "" {
		return entry.Status
	}

	return deriveStatusFromType(entry.Type)
}

// deriveStatusFromType maps journal type to default status
func deriveStatusFromType(journalType string) string {
	if journalType == "TAKEN" {
		return "SET"
	}
	return "FAIL"
}

// applyDefaultsWithLogging applies default values and logs changes
func applyDefaultsWithLogging(entry LegacyJournalEntry, ticker, filePath string, logger *MigrationLogger) (string, string, string) {
	// Handle sequence with logging (WDH is now a valid sequence per PRD 4.8.6.3.1)
	sequence := entry.Sequence
	if sequence == "" {
		logger.LogModification(filePath, ticker, "sequence", "default_applied", "set to MWD (was empty)")
		sequence = "MWD"
	}

	// Handle type with logging
	journalType := entry.Type
	if journalType == "" {
		logger.LogModification(filePath, ticker, "type", "default_applied", "set to REJECTED (was empty)")
		journalType = "REJECTED"
	}

	// Handle status with logging
	status := entry.Status
	if status == "" {
		defaultStatus := deriveStatusFromType(journalType)
		logger.LogModification(filePath, ticker, "status", "default_applied",
			fmt.Sprintf("set to %s based on type %s (was empty)", defaultStatus, journalType))
		status = defaultStatus
	}

	return sequence, status, journalType
}

// logImageCountChanges logs image count modifications
func logImageCountChanges(originalCount, entryImageCount, finalCount int, ticker, filePath string, logger *MigrationLogger) {
	if originalCount != entryImageCount || finalCount != originalCount {
		logger.LogInfo("image_count_change", map[string]any{
			"file":     filepath.Base(filePath),
			"ticker":   ticker,
			"original": entryImageCount,
			"final":    finalCount,
		})
	}
}

// ============================================================================
// MIGRATION STATS
// ============================================================================

type MigrationStats struct {
	TotalEntries    int
	SuccessCount    int
	FailureCount    int
	FailedTickers   []string
	FailureMessages []string
}

type migrationRunContext struct {
	client         *resty.Client
	logger         *MigrationLogger
	stats          *MigrationStats
	migratedCounts *ProcessingCounts
	accounting     *MigrationAccounting
}

func assertCreatedJournal(resp *resty.Response) common.Envelope[barkat.Journal] {
	Expect(resp.StatusCode()).To(Equal(http.StatusCreated))

	var envelope common.Envelope[barkat.Journal]
	Expect(json.Unmarshal(resp.Body(), &envelope)).To(Succeed())

	return envelope
}

func migrateJournalWithLogging(ctx migrationRunContext, entry LegacyJournalEntry, converted ConvertedJournal, filePath string) {
	journal := converted.Journal
	ctx.stats.TotalEntries++
	ctx.accounting.RecordConversion(converted)

	resp, err := ctx.client.R().SetBody(journal).Post(barkat.JournalBase)
	Expect(err).ToNot(HaveOccurred())

	if resp.StatusCode() == http.StatusCreated {
		ctx.stats.SuccessCount++

		envelope := assertCreatedJournal(resp)
		ctx.migratedCounts.MigratedJournals++
		ctx.migratedCounts.MigratedImages += len(envelope.Data.Images)
		ctx.migratedCounts.MigratedNotes += len(envelope.Data.Notes)
		ctx.migratedCounts.MigratedTags += len(envelope.Data.Tags)
		verifyMigratedJournalMatchesProjection(envelope.Data, journal)

		ctx.logger.LogSuccess(filePath, entry.Ticker, envelope.Data.ExternalID,
			len(envelope.Data.Images), len(envelope.Data.Notes) > 0, len(envelope.Data.Tags) > 0)
	} else {
		ctx.stats.FailureCount++
		ctx.stats.FailedTickers = append(ctx.stats.FailedTickers, entry.Ticker)
		ctx.stats.FailureMessages = append(ctx.stats.FailureMessages, string(resp.Body()))
		ctx.logger.LogError(filePath, entry.Ticker, entry.LineNumber, string(resp.Body()), entry.RawLine)
	}

	Expect(resp.StatusCode()).To(Equal(http.StatusCreated),
		"Migration should succeed for %s: status=%d body=%s",
		entry.Ticker, resp.StatusCode(), string(resp.Body()))
}

// ============================================================================
// TEST SUITE
// ============================================================================

var _ = Describe("Barkat Migration Test", func() {
	var (
		client *resty.Client
		logger *MigrationLogger
	)

	newMigrationTestClient := func(baseURL string) *resty.Client {
		baseClient := resty.New()
		baseClient.SetTimeout(10 * time.Second)
		baseClient.SetHeader("Content-Type", "application/json")
		baseClient.SetBaseURL(baseURL)
		return baseClient
	}

	printParsedCounts := func(counts ProcessingCounts) {
		GinkgoWriter.Printf("Parsed tickers: %d\n", counts.ParsedTickers)
		GinkgoWriter.Printf("Parsed images: %d\n", counts.ParsedImages)
		GinkgoWriter.Printf("Parsed notes: %d\n", counts.ParsedNotes)
	}

	BeforeEach(func() {
		if RealServerURL != "" {
			client = newMigrationTestClient(RealServerURL)
			GinkgoWriter.Printf("Using real server: %s\n", RealServerURL)
		} else {
			client = newMigrationTestClient(fmt.Sprintf("http://localhost:%d", testPort))
			GinkgoWriter.Printf("Using test server: http://localhost:%d (in-memory DB)\n", testPort)
		}

		// Create logger in temp directory for cleanup
		var err error
		logger, err = NewMigrationLogger("/tmp")
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		if logger != nil {
			logPath := logger.GetLogPath()
			logger.Close()

			// Clean up migration log file from temp directory
			if logPath != "" {
				os.Remove(logPath) // Ignore errors, cleanup is best effort
			}
		}
	})

	Context("Parser helpers", func() {
		It("should preserve raw markdown and code block note content for a parsed entry", func() {
			tempDir := GinkgoT().TempDir()
			filePath := filepath.Join(tempDir, "2023_06_15.md")
			content := strings.Join([]string{
				"| `SBIN` | #t.set #t.running #t.full #r.dep-loc #m.ntr |",
				"- review note",
				"```md",
				"plan line 1",
				"plan line 2",
				"```",
				"![chart](images/sbin.png)",
				"",
			}, "\n")

			Expect(os.WriteFile(filePath, []byte(content), 0o600)).To(Succeed())

			entries, parsedCounts, err := parseLegacyMarkdownWithLogging(filePath, logger)
			Expect(err).ToNot(HaveOccurred())
			Expect(entries).To(HaveLen(1))
			Expect(parsedCounts.ParsedTickers).To(Equal(1))
			Expect(parsedCounts.ParsedImages).To(Equal(1))
			Expect(parsedCounts.ParsedNotes).To(Equal(2))

			entry := entries[0]
			Expect(entry.Note).To(Equal("plan line 1\nplan line 2"))
			Expect(entry.SimpleNotes).To(Equal([]string{"review note"}))
			Expect(entry.RawMarkdown).To(ContainSubstring("| `SBIN` | #t.set #t.running #t.full #r.dep-loc #m.ntr |"))
			Expect(entry.RawMarkdown).To(ContainSubstring("- review note"))
			Expect(entry.RawMarkdown).To(ContainSubstring("plan line 1"))
			Expect(entry.RawMarkdown).To(ContainSubstring("![chart](images/sbin.png)"))
			Expect(entry.RawMarkdown).To(ContainSubstring("#t.full"))
		})

		It("should map legacy trade status aliases without changing unsupported tags", func() {
			entry := &LegacyJournalEntry{}

			entry = processTradeTag("rejected", entry)
			entry = processTradeTag("taken", entry)
			Expect(entry.Type).To(Equal("TAKEN"))
			Expect(entry.Status).To(BeEmpty())

			entry = processTradeTag("miss", entry)
			Expect(entry.Status).To(Equal("MISSED"))

			entry = processTradeTag("justloss", entry)
			Expect(entry.Status).To(Equal("JUST_LOSS"))

			entry = processTradeTag("full", entry)
			Expect(entry.Status).To(Equal("JUST_LOSS"))
		})

		It("should skip creating placeholder notes when the entry has no note content", func() {
			notes := buildNotesFromLegacy(LegacyJournalEntry{Type: "TAKEN"}, time.Date(2023, time.June, 15, 0, 0, 0, 0, time.UTC))
			Expect(notes).To(BeEmpty())
		})

		It("should extract image dates from both jpg and png filename patterns", func() {
			journalDate := time.Date(2023, time.June, 15, 0, 0, 0, 0, time.UTC)

			Expect(extractImageDateFromFilename("SNAG--2023-06-12-014.jpg", journalDate).Format(time.DateOnly)).To(Equal("2023-06-12"))
			Expect(extractImageDateFromFilename("AETHER.mwd.trend.result__20230901__223855.png", journalDate).Format(time.DateOnly)).To(Equal("2023-09-01"))
			Expect(extractImageDateFromFilename("JSL.mwd.trend.rejected.ooa_20230726_224243.png", journalDate).Format(time.DateOnly)).To(Equal("2023-07-26"))
			Expect(extractImageDateFromFilename("unknown.png", journalDate)).To(Equal(journalDate))
		})

		It("should reconcile skipped clipboard and placeholder image accounting", func() {
			journalDate := time.Date(2023, time.June, 15, 0, 0, 0, 0, time.UTC)

			images, stats := buildImagesWithLogging([]string{
				"../assets/image_1686845892831_0.png",
				"../assets/trading/2023/09/AETHER.mwd.trend.result__20230901__223855.png",
			}, "AETHER", journalDate, TestFile, logger)

			Expect(stats.SourceRefs).To(Equal(2))
			Expect(stats.SkippedClipboardRefs).To(Equal(1))
			Expect(stats.MigratedSourceRefs).To(Equal(1))
			Expect(stats.PlaceholderImagesAdded).To(Equal(3))
			Expect(stats.TruncatedSourceRefs).To(Equal(0))
			Expect(stats.FinalImageCount).To(Equal(4))
			Expect(images).To(HaveLen(4))
			Expect(images[0].FileName).To(Equal("AETHER.mwd.trend.result__20230901__223855.png"))
			Expect(images[0].CreatedAt.Format(time.DateOnly)).To(Equal("2023-09-01"))
		})
	})

	Context("Single File Migration", func() {
		var (
			testFilePath = filepath.Join(ProcessedFolder, TestFile)
			entries      []LegacyJournalEntry
			parsedCounts ProcessingCounts
		)

		BeforeEach(func() {
			// Check if file exists before proceeding
			if _, err := os.Stat(testFilePath); os.IsNotExist(err) {
				Skip(fmt.Sprintf("Test file not found: %s", testFilePath))
				return
			}

			var err error
			entries, parsedCounts, err = parseLegacyMarkdownWithLogging(testFilePath, logger)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should parse legacy markdown file with detailed logging", func() {
			Expect(entries).ToNot(BeEmpty())
			Expect(len(entries)).To(BeNumerically(">=", 6))

			// Log parsing results
			GinkgoWriter.Printf("\n=== Parsing Results ===\n")
			printParsedCounts(parsedCounts)
		})

		It("should migrate entries via POST API with logging", func() {
			journalDate, err := extractDateFromFilename(testFilePath)
			Expect(err).ToNot(HaveOccurred())

			stats := MigrationStats{}
			var migratedCounts ProcessingCounts
			var accounting MigrationAccounting
			migrationCtx := migrationRunContext{client: client, logger: logger, stats: &stats, migratedCounts: &migratedCounts, accounting: &accounting}

			defer func() {
				GinkgoWriter.Printf("\n=== Migration Stats ===\n")
				GinkgoWriter.Printf("Total Entries: %d\n", stats.TotalEntries)
				GinkgoWriter.Printf("Success: %d\n", stats.SuccessCount)
				GinkgoWriter.Printf("Failures: %d\n", stats.FailureCount)
				GinkgoWriter.Printf("Log file: %s\n", logger.GetLogPath())
			}()

			for _, entry := range entries {
				converted := convertToJournalWithLogging(entry, journalDate, testFilePath, logger)
				migrateJournalWithLogging(migrationCtx, entry, converted, testFilePath)

			}

			GinkgoWriter.Printf("\n=== Parsing Results ===\n")
			printParsedCounts(parsedCounts)

			GinkgoWriter.Printf("\n=== Migration Results ===\n")
			GinkgoWriter.Printf("Migrated journals: %d\n", migratedCounts.MigratedJournals)
			GinkgoWriter.Printf("Migrated images: %d\n", migratedCounts.MigratedImages)
			GinkgoWriter.Printf("Migrated notes: %d\n", migratedCounts.MigratedNotes)
			GinkgoWriter.Printf("Migrated tags: %d\n", migratedCounts.MigratedTags)

			Expect(accounting.SourceEntries).To(Equal(parsedCounts.ParsedTickers))
			Expect(accounting.SourceImageRefs).To(Equal(parsedCounts.ParsedImages))
			accounting.VerifyMigratedCounts(migratedCounts)
		})

		It("should keep sampled migrated note content meaningful", func() {
			journalDate, err := extractDateFromFilename(testFilePath)
			Expect(err).ToNot(HaveOccurred())

			sampled := 0
			for _, entry := range entries {
				converted := convertToJournalWithLogging(entry, journalDate, testFilePath, logger)
				journal := converted.Journal
				if len(journal.Notes) == 0 {
					continue
				}

				resp, reqErr := client.R().SetBody(journal).Post(barkat.JournalBase)
				Expect(reqErr).ToNot(HaveOccurred())
				Expect(resp.StatusCode()).To(Equal(http.StatusCreated))

				var envelope common.Envelope[barkat.Journal]
				Expect(json.Unmarshal(resp.Body(), &envelope)).To(Succeed())
				Expect(envelope.Data.Notes).ToNot(BeEmpty())

				noteContent := strings.TrimSpace(envelope.Data.Notes[0].Content)
				Expect(noteContent).ToNot(Equal("=== ORIGINAL MARKDOWN ==="))
				Expect(noteContent).ToNot(Equal("=== PLAN NOTES ==="))

				if after, ok := strings.CutPrefix(noteContent, "=== ORIGINAL MARKDOWN ==="); ok {
					Expect(strings.TrimSpace(after)).ToNot(BeEmpty())
				}

				sampled++
				if sampled == 3 {
					break
				}
			}

			Expect(sampled).To(Equal(3))
		})
	})

	Context("Folder Migration with Full Validation", func() {
		var (
			testFolder = ProcessedFolder
			allFiles   []string
		)

		BeforeEach(func() {
			files, err := filepath.Glob(filepath.Join(testFolder, "*.md"))
			if err != nil {
				// If folder doesn't exist or is inaccessible, skip this test
				Skip(fmt.Sprintf("Processed folder not accessible: %v", err))
				return
			}

			if len(files) == 0 {
				// If no markdown files found, skip this test
				Skip("No markdown files found in processed folder")
				return
			}

			sort.Strings(files) // Process in order
			allFiles = files
		})

		It("should migrate all files with detailed logging", func() {
			var totalParsedCounts ProcessingCounts
			var totalMigratedCounts ProcessingCounts
			totalStats := MigrationStats{}
			var accounting MigrationAccounting
			migrationCtx := migrationRunContext{client: client, logger: logger, stats: &totalStats, migratedCounts: &totalMigratedCounts, accounting: &accounting}

			defer func() {
				logger.WriteSummary(totalParsedCounts, len(allFiles))

				GinkgoWriter.Printf("\n" + strings.Repeat("=", 80) + "\n")
				GinkgoWriter.Printf("MIGRATION SUMMARY REPORT\n")
				GinkgoWriter.Printf(strings.Repeat("=", 80) + "\n\n")

				GinkgoWriter.Printf("FILES PROCESSED: %d\n\n", len(allFiles))

				GinkgoWriter.Printf("--- PROCESSING SUMMARY ---\n")
				GinkgoWriter.Printf("  Files Processed:     %d\n", len(allFiles))

				GinkgoWriter.Printf("--- PARSING RESULTS ---\n")
				GinkgoWriter.Printf("  Parsed Tickers: %d\n", totalParsedCounts.ParsedTickers)
				GinkgoWriter.Printf("  Parsed Images:  %d\n", totalParsedCounts.ParsedImages)
				GinkgoWriter.Printf("  Parsed Notes:   %d\n\n", totalParsedCounts.ParsedNotes)

				GinkgoWriter.Printf("--- MIGRATION RESULTS ---\n")
				GinkgoWriter.Printf("  Migrated Journals: %d\n", totalMigratedCounts.MigratedJournals)
				GinkgoWriter.Printf("  Migrated Images:   %d\n", totalMigratedCounts.MigratedImages)
				GinkgoWriter.Printf("  Migrated Notes:    %d\n", totalMigratedCounts.MigratedNotes)
				GinkgoWriter.Printf("  Migrated Tags:     %d\n\n", totalMigratedCounts.MigratedTags)

				GinkgoWriter.Printf("--- PROCESSING STATUS ---\n")
				if totalStats.TotalEntries > 0 {
					GinkgoWriter.Printf("  Total Entries Processed: %d\n", totalStats.TotalEntries)
					GinkgoWriter.Printf("  Success Rate: %.2f%%\n\n", float64(totalStats.SuccessCount)/float64(totalStats.TotalEntries)*100)
				}

				GinkgoWriter.Printf("--- PROCESSING STATS ---\n")
				GinkgoWriter.Printf("  Total Entries:   %d\n", totalStats.TotalEntries)
				GinkgoWriter.Printf("  Success:         %d\n", totalStats.SuccessCount)
				GinkgoWriter.Printf("  Failures:        %d\n", totalStats.FailureCount)
				GinkgoWriter.Printf("  Sanitizations:   %d\n", logger.GetSanitizationCount())
				GinkgoWriter.Printf("  Modifications:   %d\n\n", logger.GetModificationCount())

				if totalStats.FailureCount > 0 {
					GinkgoWriter.Printf("--- FAILED ENTRIES (first 20) ---\n")
					limit := min(len(totalStats.FailedTickers), 20)
					for i := range limit {
						GinkgoWriter.Printf("  - %s: %s\n", totalStats.FailedTickers[i], totalStats.FailureMessages[i])
					}
					if len(totalStats.FailedTickers) > 20 {
						GinkgoWriter.Printf("  ... and %d more (see log file)\n", len(totalStats.FailedTickers)-20)
					}
					GinkgoWriter.Printf("\n")
				}

				GinkgoWriter.Printf("LOG FILE: %s\n", logger.GetLogPath())
				GinkgoWriter.Printf(strings.Repeat("=", 80) + "\n")
				GinkgoWriter.Printf("Migration completed with detailed logging\n")
			}()

			for _, filePath := range allFiles {
				// Parse with logging
				entries, parsedCounts, err := parseLegacyMarkdownWithLogging(filePath, logger)
				if err != nil {
					logger.LogError(filePath, "", 0, fmt.Sprintf("parse error: %v", err), "")
					continue
				}

				totalParsedCounts.ParsedTickers += parsedCounts.ParsedTickers
				totalParsedCounts.ParsedImages += parsedCounts.ParsedImages
				totalParsedCounts.ParsedNotes += parsedCounts.ParsedNotes

				// Extract date from filename
				journalDate, err := extractDateFromFilename(filePath)
				if err != nil {
					logger.LogError(filePath, "", 0, fmt.Sprintf("date extraction error: %v", err), "")
					// Use file modification time as fallback
					info, _ := os.Stat(filePath)
					if info != nil {
						journalDate = info.ModTime()
					} else {
						journalDate = time.Now()
					}
				}

				// Migrate entries
				for _, entry := range entries {
					converted := convertToJournalWithLogging(entry, journalDate, filePath, logger)
					migrateJournalWithLogging(migrationCtx, entry, converted, filePath)
				}
			}

			Expect(accounting.SourceEntries).To(Equal(totalParsedCounts.ParsedTickers))
			Expect(accounting.SourceImageRefs).To(Equal(totalParsedCounts.ParsedImages))
			accounting.VerifyMigratedCounts(totalMigratedCounts)
		})
	})
})
