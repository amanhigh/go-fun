package handler

import (
	"net/http"

	kohanassets "github.com/amanhigh/go-fun/components/kohan/assets"
	"github.com/gin-gonic/gin"
)

// SetupOSRoutes configures OS-related routes
func SetupOSRoutes(os *gin.RouterGroup, osHandler OSHandler) {
	os.POST("/screenshot", osHandler.HandleScreenshot)
	os.GET("/ticker/:ticker/record", osHandler.HandleRecordTicker)
	os.GET("/clip/", osHandler.HandleReadClip)
	os.POST("/submap/:action", osHandler.HandleSubmapControl)
}

// SetupImageRoutes configures image-related routes for the given journal router group
func SetupImageRoutes(journal *gin.RouterGroup, imageHandler ImageHandler) {
	images := journal.Group("/:id/images")
	{
		images.POST("", imageHandler.HandleCreateImage)
		images.GET("", imageHandler.HandleListImages)
		images.DELETE("/:imageId", imageHandler.HandleDeleteImage)
	}
}

// SetupNoteRoutes configures note-related routes for the given journal router group
func SetupNoteRoutes(journal *gin.RouterGroup, noteHandler NoteHandler) {
	notes := journal.Group("/:id/notes")
	{
		notes.POST("", noteHandler.HandleCreateNote)
		notes.GET("", noteHandler.HandleListNotes)
		notes.DELETE("/:noteId", noteHandler.HandleDeleteNote)
	}
}

// SetupTagRoutes configures tag-related routes for the given journal router group
func SetupTagRoutes(journal *gin.RouterGroup, tagHandler TagHandler) {
	tags := journal.Group("/:id/tags")
	{
		tags.POST("", tagHandler.HandleCreateTag)
		tags.GET("", tagHandler.HandleListTags)
		tags.DELETE("/:tagId", tagHandler.HandleDeleteTag)
	}
}

// SetupStaticRoutes configures all static file serving.
func SetupStaticRoutes(engine *gin.Engine, imagePath string) {
	engine.StaticFS("/assets", http.FS(kohanassets.FS))
	engine.Static("/journal/images", imagePath)
}

// SetupJournalRoutes configures basic journal routes
func SetupJournalRoutes(journal *gin.RouterGroup, journalHandler JournalHandler) {
	{
		journal.GET("", journalHandler.HandleListJournals)
		journal.GET("/:id", journalHandler.HandleGetJournal)
		journal.POST("", journalHandler.HandleCreateJournal)
		journal.PATCH("/:id", journalHandler.HandleUpdateReviewStatus)
		journal.DELETE("/:id", journalHandler.HandleDeleteJournal)
	}
}

// SetupTickerRoutes configures Barkat ticker routes.
func SetupTickerRoutes(ticker *gin.RouterGroup, tickerHandler TickerHandler) {
	{
		ticker.GET("", tickerHandler.HandleListTickers)
		ticker.GET("/:ticker", tickerHandler.HandleGetTicker)
		ticker.POST("", tickerHandler.HandleCreateTicker)
		ticker.PUT("/:ticker", tickerHandler.HandleUpdateTicker)
		ticker.PATCH("/:ticker", tickerHandler.HandlePatchTickerLastOpened)
		ticker.DELETE("/:ticker", tickerHandler.HandleDeleteTicker)
	}
}

// HACK: Merge Route Functions.
// SetupTickerAlertRoutes configures nested Alert ticker routes under primary tickers.
func SetupTickerAlertRoutes(ticker *gin.RouterGroup, alertTickerHandler AlertTickerHandler) {
	{
		ticker.POST("/:ticker/alert-tickers", alertTickerHandler.HandleCreateAlertTicker)
	}
}

// SetupTickerPriceAlertRoutes configures nested price alert routes under primary tickers.
func SetupTickerPriceAlertRoutes(ticker *gin.RouterGroup, priceAlertHandler PriceAlertHandler) {
	{
		ticker.POST("/:ticker/alerts", priceAlertHandler.HandleCreatePendingPriceAlert)
	}
}

// SetupAlertTickerRoutes configures top-level Alert ticker routes.
func SetupAlertTickerRoutes(alertTicker *gin.RouterGroup, alertTickerHandler AlertTickerHandler) {
	{
		alertTicker.GET("", alertTickerHandler.HandleListAlertTickers)
		alertTicker.GET("/:symbol", alertTickerHandler.HandleGetAlertTicker)
		alertTicker.DELETE("/:symbol", alertTickerHandler.HandleDeleteAlertTicker)
	}
}

// SetupPriceAlertRoutes configures top-level price alert routes.
func SetupPriceAlertRoutes(alert *gin.RouterGroup, priceAlertHandler PriceAlertHandler) {
	{
		alert.GET("", priceAlertHandler.HandleListPriceAlerts)
		alert.PUT("", priceAlertHandler.HandleReplacePriceAlerts)
		alert.DELETE("/:alert-id", priceAlertHandler.HandleDeletePriceAlert)
	}
}

// SetupAuditRoutes configures Barkat audit routes.
func SetupAuditRoutes(audits *gin.RouterGroup, auditHandler AuditHandler) {
	{
		audits.GET("", auditHandler.HandleListAudits)
		audits.GET("/:audit-id/results", auditHandler.HandleExecuteAudit)
	}
}
