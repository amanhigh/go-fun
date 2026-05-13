package core

import (
	"regexp"
	"time"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/golang-sql/civil"
)

// Pre-compiled regex patterns for validation (PRD Section 3.0)
var (
	// Ticker: uppercase A-Z, digits, dots, underscores, exclamation (e.g., "TCS", "TCS.NS", "GOLD!")
	tickerRegex = regexp.MustCompile(`^[A-Z0-9][A-Z0-9._!]*$`)
	// Alert Symbol: alphanumeric first, then dot, slash, equals (e.g. "USDCAD", "BTC/USD").
	alertSymbolRegex = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9./=]*$`)
	// Alert Name: alphanumeric first, then sanitized display name characters.
	alertNameRegex = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9 .&'()-]*$`)
	// Ticker Exchange: uppercase letters, underscore, dot (source data from exchangeRepo, e.g. FX_IDC).
	tickerExchangeRegex = regexp.MustCompile(`^[A-Z][A-Z._]*$`)
	// Alert Exchange: letters only, mixed-case (source data from pairRepo).
	alertExchangeRegex = regexp.MustCompile(`^[A-Za-z][A-Za-z]*$`)
	// Tag: alphanumeric with hyphens (e.g., "oe", "dep-1")
	tagRegex = regexp.MustCompile(`^[a-zA-Z0-9-]+$`)
	// Override: letters only (e.g., "loc", "abc")
	overrideRegex = regexp.MustCompile(`^[a-zA-Z]*$`)
	// JournalID: jrn_ prefix followed by 8 alphanumeric characters (e.g., "jrn_12345678")
	journalIDRegex = regexp.MustCompile(`^jrn_[a-zA-Z0-9]{8}$`)
	// NoteID: not_ prefix followed by 8 alphanumeric characters (e.g., "not_12345678")
	noteIDRegex = regexp.MustCompile(`^not_[a-zA-Z0-9]{8}$`)
	// TagID: tag_ prefix followed by 8 alphanumeric characters (e.g., "tag_12345678")
	tagIDRegex = regexp.MustCompile(`^tag_[a-zA-Z0-9]{8}$`)
	// ImageID: img_ prefix followed by 8 alphanumeric characters (e.g., "img_12345678")
	imageIDRegex = regexp.MustCompile(`^img_[a-zA-Z0-9]{8}$`)
	// ImageFile: alphanumeric with dots, hyphens, underscores, valid image extensions
	imageFileRegex = regexp.MustCompile(`^[a-zA-Z0-9._-]+\.(png|jpg|jpeg)$`)
)

// RegisterJournalValidators registers custom validators for journal fields.
// Must be called before using Gin binding for journal requests.
func RegisterJournalValidators() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("ticker", TickerValidator)
		_ = v.RegisterValidation("tv_ticker_path", TickerValidator)
		_ = v.RegisterValidation("alert_symbol", AlertSymbolValidator)
		_ = v.RegisterValidation("alert_name", AlertNameValidator)
		_ = v.RegisterValidation("ticker_exchange", TickerExchangeValidator)
		_ = v.RegisterValidation("alert_exchange", AlertExchangeValidator)
		_ = v.RegisterValidation("tag", TagValidator)
		_ = v.RegisterValidation("override", OverrideValidator)
		_ = v.RegisterValidation("image_file", ImageFileValidator)
		_ = v.RegisterValidation("not_future", NotFutureValidator)
		_ = v.RegisterValidation("journal_id", JournalIDValidator)
		_ = v.RegisterValidation("note_id", NoteIDValidator)
		_ = v.RegisterValidation("tag_id", TagIDValidator)
		_ = v.RegisterValidation("image_id", ImageIDValidator)
	}
}

// AlertSymbolValidator validates Alert ticker symbol format using pre-compiled regex
func AlertSymbolValidator(fl validator.FieldLevel) bool {
	field := fl.Field().String()
	return field == "" || alertSymbolRegex.MatchString(field)
}

// AlertNameValidator validates Alert ticker name format using pre-compiled regex
func AlertNameValidator(fl validator.FieldLevel) bool {
	field := fl.Field().String()
	return field == "" || alertNameRegex.MatchString(field)
}

// TickerExchangeValidator validates exchange format using pre-compiled regex
func TickerExchangeValidator(fl validator.FieldLevel) bool {
	field := fl.Field().String()
	return field == "" || tickerExchangeRegex.MatchString(field)
}

// AlertExchangeValidator validates Alert ticker exchange format using pre-compiled regex
func AlertExchangeValidator(fl validator.FieldLevel) bool {
	field := fl.Field().String()
	return field == "" || alertExchangeRegex.MatchString(field)
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

// ImageFileValidator validates image file name format using pre-compiled regex
func ImageFileValidator(fl validator.FieldLevel) bool {
	field := fl.Field().String()
	return field == "" || imageFileRegex.MatchString(field)
}

// NotFutureValidator validates business rule: date should not be in the future
func NotFutureValidator(fl validator.FieldLevel) bool {
	date, ok := fl.Field().Interface().(civil.Date)
	if !ok {
		return false
	}
	now := civil.DateOf(time.Now())
	return date.Before(now) || date.String() == now.String()
}

// JournalIDValidator validates journal ID format using pre-compiled regex
func JournalIDValidator(fl validator.FieldLevel) bool {
	field := fl.Field().String()
	return field == "" || journalIDRegex.MatchString(field)
}

// NoteIDValidator validates note ID format using pre-compiled regex
func NoteIDValidator(fl validator.FieldLevel) bool {
	field := fl.Field().String()
	return field == "" || noteIDRegex.MatchString(field)
}

// TagIDValidator validates tag ID format using pre-compiled regex
func TagIDValidator(fl validator.FieldLevel) bool {
	field := fl.Field().String()
	return field == "" || tagIDRegex.MatchString(field)
}

// ImageIDValidator validates image ID format using pre-compiled regex
func ImageIDValidator(fl validator.FieldLevel) bool {
	field := fl.Field().String()
	return field == "" || imageIDRegex.MatchString(field)
}
