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

func (l *MigrationLogger) LogInfo(message string, data interface{}) {
	l.writeJSON("INFO", map[string]interface{}{"message": message, "data": data})
}

func (l *MigrationLogger) writeJSON(logType string, data interface{}) {
	jsonData, _ := json.Marshal(data)
	fmt.Fprintf(l.logFile, "%s|%s|%s\n", time.Now().Format(time.RFC3339), logType, string(jsonData))
}

func (l *MigrationLogger) WriteSummary(processed ProcessingCounts, files int) {
	summary := map[string]interface{}{
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
	Ticker      string
	RawTags     []string // Store all raw tags for logging
	Sequence    string
	Type        string
	Status      string
	Direction   string
	Reason      string
	Images      []string
	Note        string
	RawLine     string // Original line for reference
	LineNumber  int    // Line number in file
	RawMarkdown string // Complete raw markdown block for this entry
}

// ============================================================================
// DATE EXTRACTION
// ============================================================================

// extractDateFromFilename extracts date from filename like "2023_06_15.md"
func extractDateFromFilename(filename string) (time.Time, error) {
	base := filepath.Base(filename)
	base = strings.TrimSuffix(base, ".md")

	// Try format: 2023_06_15
	t, err := time.Parse("2006_01_02", base)
	if err == nil {
		return t, nil
	}

	// Try format: 2023-06-15
	t, err = time.Parse("2006-01-02", base)
	if err == nil {
		return t, nil
	}

	return time.Time{}, fmt.Errorf("cannot parse date from filename: %s", filename)
}

// ============================================================================
// PARSING WITH LOGGING
// ============================================================================

// parseLegacyMarkdownWithLogging parses markdown and captures raw content
func parseLegacyMarkdownWithLogging(filePath string, logger *MigrationLogger) ([]LegacyJournalEntry, ProcessingCounts, error) {
	var counts ProcessingCounts

	file, err := os.Open(filePath)
	if err != nil {
		return nil, counts, err
	}
	defer file.Close()

	var entries []LegacyJournalEntry
	var currentEntry *LegacyJournalEntry
	var inCodeBlock bool
	var noteContent strings.Builder
	var rawMarkdown strings.Builder
	var lineNumber int

	scanner := bufio.NewScanner(file)
	// Regex patterns
	entryPattern := regexp.MustCompile(`\|\s*` + "`" + `([A-Z0-9_!]+)` + "`" + `\s*\|(.+)\|`)
	imagePattern := regexp.MustCompile(`!\[.*?\]\(([^)]+)\)`)
	tagPattern := regexp.MustCompile(`#([trm])\.([a-z0-9-]+)`)

	for scanner.Scan() {
		line := scanner.Text()
		lineNumber++

		// Track raw markdown for current entry
		if currentEntry != nil {
			rawMarkdown.WriteString(line + "\n")
		}

		// Handle code blocks for notes
		if strings.Contains(line, "```") {
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
			inCodeBlock = !inCodeBlock
			continue
		}

		if inCodeBlock {
			noteContent.WriteString(line + "\n")
			continue
		}

		// Check for entry line (table row with ticker)
		if matches := entryPattern.FindStringSubmatch(line); matches != nil {
			// Save previous entry if exists
			if currentEntry != nil {
				currentEntry.RawMarkdown = rawMarkdown.String()
				entries = append(entries, *currentEntry)
			}

			ticker := matches[1]
			tagsPart := matches[2]

			entry := LegacyJournalEntry{
				Ticker:     ticker,
				RawLine:    line,
				LineNumber: lineNumber,
			}
			counts.ParsedTickers++

			// Reset raw markdown for new entry
			rawMarkdown.Reset()
			rawMarkdown.WriteString(line + "\n")

			// Parse and store raw tags
			tags := tagPattern.FindAllStringSubmatch(tagsPart, -1)
			for _, tag := range tags {
				entry.RawTags = append(entry.RawTags, fmt.Sprintf("#%s.%s", tag[1], tag[2]))
				prefix := tag[1]
				value := tag[2]

				switch prefix {
				case "t":
					switch value {
					case "mwd", "yr":
						entry.Sequence = strings.ToUpper(value)
					case "wdh":
						entry.Sequence = "WDH" // Store as-is, convert later
					case "rejected":
						entry.Type = "REJECTED"
					case "set":
						entry.Type = "SET"
					case "result":
						entry.Type = "RESULT"
					case "fail":
						entry.Status = "FAIL"
					case "taken":
						entry.Status = "TAKEN"
					case "success":
						entry.Status = "SUCCESS"
					case "running":
						entry.Status = "RUNNING"
					case "broken":
						entry.Status = "BROKEN"
					case "missed":
						entry.Status = "MISSED"
					case "dropped":
						entry.Status = "DROPPED"
					case "trend", "ctrend":
						entry.Direction = value
					}
				case "r":
					entry.Reason = value
				}
			}

			currentEntry = &entry
			continue
		}

		// Check for image line
		if currentEntry != nil {
			if matches := imagePattern.FindStringSubmatch(line); matches != nil {
				currentEntry.Images = append(currentEntry.Images, matches[1])
				counts.ParsedImages++
			}
		}
	}

	// Don't forget the last entry
	if currentEntry != nil {
		currentEntry.RawMarkdown = rawMarkdown.String()
		entries = append(entries, *currentEntry)
	}

	return entries, counts, scanner.Err()
}

// ============================================================================
// CONVERSION WITH LOGGING
// ============================================================================

// sanitizeFileNameWithLogging removes invalid characters and logs changes
func sanitizeFileNameWithLogging(name, file, ticker string, logger *MigrationLogger) string {
	invalidChars := regexp.MustCompile(`[!@#$%^&*()+=\[\]{}|;:'",<>?/\\]`)
	sanitized := invalidChars.ReplaceAllString(name, "_")

	if sanitized != name {
		logger.LogSanitization(file, ticker, "filename", name, sanitized, "invalid_characters_removed")
	}

	return sanitized
}

// sanitizeTickerWithLogging sanitizes ticker symbols and logs changes
func sanitizeTickerWithLogging(ticker, file string, logger *MigrationLogger) string {
	original := ticker

	// Remove trailing ! (futures symbols)
	if strings.HasSuffix(ticker, "!") {
		ticker = strings.TrimSuffix(ticker, "!")
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
func convertToJournalWithLogging(entry LegacyJournalEntry, journalDate time.Time, filePath string, logger *MigrationLogger) barkat.Journal {
	ticker := sanitizeTickerWithLogging(entry.Ticker, filePath, logger)

	// Build images with logging
	images := make([]barkat.Image, 0, len(entry.Images))
	timeframes := []string{"DL", "WK", "MN", "TMN"}

	for i, imgPath := range entry.Images {
		timeframe := timeframes[i%len(timeframes)]
		sanitizedName := sanitizeFileNameWithLogging(filepath.Base(imgPath), filePath, ticker, logger)
		images = append(images, barkat.Image{
			Timeframe: timeframe,
			FileName:  sanitizedName,
			CreatedAt: journalDate,
		})
	}

	originalImageCount := len(images)

	// Log placeholder additions
	if len(images) < 4 {
		logger.LogModification(filePath, ticker, "images", "placeholder_added",
			fmt.Sprintf("added %d placeholder images (had %d, need 4)", 4-len(images), len(images)))
	}
	for len(images) < 4 {
		images = append(images, barkat.Image{
			Timeframe: timeframes[len(images)],
			FileName:  fmt.Sprintf("placeholder_%d.png", len(images)),
			CreatedAt: journalDate,
		})
	}

	// Log truncation
	if len(images) > 16 {
		logger.LogModification(filePath, ticker, "images", "truncated",
			fmt.Sprintf("truncated from %d to 16 images", len(images)))
		images = images[:16]
	}

	// Build tags
	var tags []barkat.Tag
	if entry.Reason != "" {
		tags = append(tags, barkat.Tag{
			Tag:       entry.Reason,
			Type:      "REASON",
			CreatedAt: journalDate,
		})
	}

	// Build notes - include original markdown as first note
	var notes []barkat.Note

	// Determine the note status (must be a valid journal status)
	// Use the entry's status, or derive from type if not set
	noteStatus := entry.Status
	if noteStatus == "" {
		switch entry.Type {
		case "REJECTED":
			noteStatus = "FAIL"
		case "SET":
			noteStatus = "SET"
		case "RESULT":
			noteStatus = "SUCCESS"
		default:
			noteStatus = "FAIL"
		}
	}

	// Add original markdown content as note for preservation
	originalContent := fmt.Sprintf("=== ORIGINAL MARKDOWN ===\n%s", entry.RawMarkdown)
	notes = append(notes, barkat.Note{
		Status:    noteStatus,
		Content:   originalContent,
		Format:    "MARKDOWN",
		CreatedAt: journalDate,
	})

	// Handle sequence mapping with logging
	sequence := entry.Sequence
	if sequence == "WDH" {
		logger.LogSanitization(filePath, ticker, "sequence", "WDH", "MWD", "legacy_sequence_mapped")
		sequence = "MWD"
	}
	if sequence == "" {
		logger.LogModification(filePath, ticker, "sequence", "default_applied", "set to MWD (was empty)")
		sequence = "MWD"
	}

	// Handle status with logging
	status := entry.Status
	if status == "" {
		var defaultStatus string
		switch entry.Type {
		case "REJECTED":
			defaultStatus = "FAIL"
		case "SET":
			defaultStatus = "SET"
		case "RESULT":
			defaultStatus = "SUCCESS"
		default:
			defaultStatus = "FAIL"
		}
		logger.LogModification(filePath, ticker, "status", "default_applied",
			fmt.Sprintf("set to %s based on type %s (was empty)", defaultStatus, entry.Type))
		status = defaultStatus
	}

	// Handle type with logging
	journalType := entry.Type
	if journalType == "" {
		logger.LogModification(filePath, ticker, "type", "default_applied", "set to REJECTED (was empty)")
		journalType = "REJECTED"
	}

	// Log if we had to add placeholders or truncate
	if originalImageCount != len(entry.Images) || len(images) != originalImageCount {
		logger.LogInfo("image_count_change", map[string]interface{}{
			"file":     filepath.Base(filePath),
			"ticker":   ticker,
			"original": len(entry.Images),
			"final":    len(images),
		})
	}

	return barkat.Journal{
		Ticker:    ticker,
		Sequence:  sequence,
		Type:      journalType,
		Status:    status,
		CreatedAt: journalDate,
		Images:    images,
		Tags:      tags,
		Notes:     notes,
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

// ============================================================================
// TEST SUITE
// ============================================================================

var _ = Describe("Barkat Migration Test", func() {
	var (
		client *resty.Client
		logger *MigrationLogger
	)

	BeforeEach(func() {
		client = resty.New()
		client.SetTimeout(10 * time.Second)
		client.SetHeader("Content-Type", "application/json")
		client.SetBaseURL(fmt.Sprintf("http://localhost:%d", testPort))

		// Create logger in test directory
		var err error
		logger, err = NewMigrationLogger("/home/aman/Projects/go-fun/components/kohan/core")
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		if logger != nil {
			logger.Close()
		}
	})

	Context("Single File Migration", func() {
		var (
			testFilePath string
			entries      []LegacyJournalEntry
			parsedCounts ProcessingCounts
		)

		BeforeEach(func() {
			testFilePath = "/home/aman/Projects/go-fun/processed/2023_06_15.md"

			var err error
			entries, parsedCounts, err = parseLegacyMarkdownWithLogging(testFilePath, logger)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should parse legacy markdown file with detailed logging", func() {
			Expect(entries).ToNot(BeEmpty())
			Expect(len(entries)).To(BeNumerically(">=", 6))

			// Log parsing results
			GinkgoWriter.Printf("\n=== Parsing Results ===\n")
			GinkgoWriter.Printf("Parsed tickers: %d\n", parsedCounts.ParsedTickers)
			GinkgoWriter.Printf("Parsed images: %d\n", parsedCounts.ParsedImages)
			GinkgoWriter.Printf("Parsed notes: %d\n", parsedCounts.ParsedNotes)
		})

		It("should migrate entries via POST API with logging", func() {
			journalDate, err := extractDateFromFilename(testFilePath)
			Expect(err).ToNot(HaveOccurred())

			stats := MigrationStats{}
			var migratedCounts ProcessingCounts

			for _, entry := range entries {
				journal := convertToJournalWithLogging(entry, journalDate, testFilePath, logger)

				resp, err := client.R().SetBody(journal).Post(barkat.JournalBase)
				Expect(err).ToNot(HaveOccurred())

				stats.TotalEntries++

				if resp.StatusCode() == http.StatusCreated {
					stats.SuccessCount++

					var envelope common.Envelope[barkat.Journal]
					Expect(json.Unmarshal(resp.Body(), &envelope)).To(Succeed())

					migratedCounts.MigratedJournals++
					migratedCounts.MigratedImages += len(envelope.Data.Images)
					migratedCounts.MigratedNotes += len(envelope.Data.Notes)
					migratedCounts.MigratedTags += len(envelope.Data.Tags)

					logger.LogSuccess(testFilePath, entry.Ticker, envelope.Data.ExternalID,
						len(envelope.Data.Images), len(envelope.Data.Notes) > 0, len(envelope.Data.Tags) > 0)
				} else {
					stats.FailureCount++
					stats.FailedTickers = append(stats.FailedTickers, entry.Ticker)
					stats.FailureMessages = append(stats.FailureMessages, string(resp.Body()))
					logger.LogError(testFilePath, entry.Ticker, entry.LineNumber, string(resp.Body()), entry.RawLine)
				}
			}

			// Report stats
			GinkgoWriter.Printf("\n=== Migration Stats ===\n")
			GinkgoWriter.Printf("Total Entries: %d\n", stats.TotalEntries)
			GinkgoWriter.Printf("Success: %d\n", stats.SuccessCount)
			GinkgoWriter.Printf("Failures: %d\n", stats.FailureCount)
			GinkgoWriter.Printf("Log file: %s\n", logger.GetLogPath())

			Expect(stats.FailureCount).To(Equal(0), "All migrations should succeed")
		})
	})

	Context("Folder Migration with Full Validation", func() {
		var (
			processedFolder string
			allFiles        []string
		)

		BeforeEach(func() {
			processedFolder = "/home/aman/Projects/go-fun/processed"

			files, err := filepath.Glob(filepath.Join(processedFolder, "*.md"))
			Expect(err).ToNot(HaveOccurred())
			sort.Strings(files) // Process in order
			allFiles = files
		})

		It("should migrate all files with detailed logging", func() {
			var totalParsedCounts ProcessingCounts
			var totalMigratedCounts ProcessingCounts
			totalStats := MigrationStats{}

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
					journal := convertToJournalWithLogging(entry, journalDate, filePath, logger)

					resp, err := client.R().SetBody(journal).Post(barkat.JournalBase)
					if err != nil {
						logger.LogError(filePath, entry.Ticker, entry.LineNumber, fmt.Sprintf("http error: %v", err), entry.RawLine)
						totalStats.FailureCount++
						continue
					}

					totalStats.TotalEntries++

					if resp.StatusCode() == http.StatusCreated {
						totalStats.SuccessCount++

						var envelope common.Envelope[barkat.Journal]
						if json.Unmarshal(resp.Body(), &envelope) == nil {
							totalMigratedCounts.MigratedJournals++
							totalMigratedCounts.MigratedImages += len(envelope.Data.Images)
							totalMigratedCounts.MigratedNotes += len(envelope.Data.Notes)
							totalMigratedCounts.MigratedTags += len(envelope.Data.Tags)

							logger.LogSuccess(filePath, entry.Ticker, envelope.Data.ExternalID,
								len(envelope.Data.Images), len(envelope.Data.Notes) > 0, len(envelope.Data.Tags) > 0)
						}
					} else {
						totalStats.FailureCount++
						totalStats.FailedTickers = append(totalStats.FailedTickers,
							fmt.Sprintf("%s:%s", filepath.Base(filePath), entry.Ticker))
						totalStats.FailureMessages = append(totalStats.FailureMessages, string(resp.Body()))
						logger.LogError(filePath, entry.Ticker, entry.LineNumber, string(resp.Body()), entry.RawLine)
					}
				}
			}

			// Write comprehensive summary
			logger.WriteSummary(totalParsedCounts, len(allFiles))

			// Print detailed report
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
			GinkgoWriter.Printf("  Total Entries Processed: %d\n", totalStats.TotalEntries)
			GinkgoWriter.Printf("  Success Rate: %.2f%%\n\n", float64(totalStats.SuccessCount)/float64(totalStats.TotalEntries)*100)

			GinkgoWriter.Printf("--- PROCESSING STATS ---\n")
			GinkgoWriter.Printf("  Total Entries:   %d\n", totalStats.TotalEntries)
			GinkgoWriter.Printf("  Success:         %d\n", totalStats.SuccessCount)
			GinkgoWriter.Printf("  Failures:        %d\n", totalStats.FailureCount)
			GinkgoWriter.Printf("  Sanitizations:   %d\n", logger.GetSanitizationCount())
			GinkgoWriter.Printf("  Modifications:   %d\n\n", logger.GetModificationCount())

			if totalStats.FailureCount > 0 {
				GinkgoWriter.Printf("--- FAILED ENTRIES (first 20) ---\n")
				limit := 20
				if len(totalStats.FailedTickers) < limit {
					limit = len(totalStats.FailedTickers)
				}
				for i := 0; i < limit; i++ {
					GinkgoWriter.Printf("  - %s: %s\n", totalStats.FailedTickers[i], totalStats.FailureMessages[i])
				}
				if len(totalStats.FailedTickers) > 20 {
					GinkgoWriter.Printf("  ... and %d more (see log file)\n", len(totalStats.FailedTickers)-20)
				}
				// HACK: Ensure No Tag is Lost #t, #r, #m or anything else.
				GinkgoWriter.Printf("\n")
			}

			GinkgoWriter.Printf("LOG FILE: %s\n", logger.GetLogPath())
			GinkgoWriter.Printf(strings.Repeat("=", 80) + "\n")

			// Log completion
			GinkgoWriter.Printf("Migration completed with detailed logging\n")

			// Allow up to 5% failure rate as per PRD
			if totalStats.TotalEntries > 0 {
				failureRate := float64(totalStats.FailureCount) / float64(totalStats.TotalEntries)
				Expect(failureRate).To(BeNumerically("<", 0.05), "Failure rate should be < 5%")
			}
		})
	})
})
