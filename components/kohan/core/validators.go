package core

import (
	"regexp"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Pre-compiled regex patterns for validation (PRD Section 3.0)
var (
	// Ticker: uppercase A-Z, digits, dots, hyphens, ampersands (e.g., "TCS", "TCS.NS", "M&M")
	tickerRegex = regexp.MustCompile(`^[A-Z0-9][A-Z0-9.\-&]*$`)
	// Tag: alphanumeric with hyphens (e.g., "oe", "dep-1")
	tagRegex = regexp.MustCompile(`^[a-zA-Z0-9-]+$`)
	// Override: letters only (e.g., "loc", "abc")
	overrideRegex = regexp.MustCompile(`^[a-zA-Z]*$`)
	// FileName: alphanumeric with dots, hyphens, underscores, valid image extensions
	fileNameRegex = regexp.MustCompile(`^[a-zA-Z0-9._-]+\.(png|jpg|jpeg)$`)
)

// RegisterJournalValidators registers custom validators for journal fields.
// Must be called before using Gin binding for journal requests.
func RegisterJournalValidators() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("ticker", TickerValidator)
		_ = v.RegisterValidation("tag", TagValidator)
		_ = v.RegisterValidation("override", OverrideValidator)
		_ = v.RegisterValidation("file_name", FileNameValidator)
	}
}

// TickerValidator validates ticker format using pre-compiled regex
func TickerValidator(fl validator.FieldLevel) bool {
	field := fl.Field().String()
	return field == "" || tickerRegex.MatchString(field)
}

// TagValidator validates tag format using pre-compiled regex
func TagValidator(fl validator.FieldLevel) bool {
	field := fl.Field().String()
	return field == "" || tagRegex.MatchString(field)
}

// OverrideValidator validates override format using pre-compiled regex
func OverrideValidator(fl validator.FieldLevel) bool {
	field := fl.Field().String()
	return field == "" || overrideRegex.MatchString(field)
}

// FileNameValidator validates file name format using pre-compiled regex
func FileNameValidator(fl validator.FieldLevel) bool {
	field := fl.Field().String()
	return field == "" || fileNameRegex.MatchString(field)
}
