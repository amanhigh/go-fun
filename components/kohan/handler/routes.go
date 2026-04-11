package handler

import "github.com/gin-gonic/gin"

// SetupOSRoutes configures OS-related routes
func SetupOSRoutes(os *gin.RouterGroup, osHandler OSHandler) {
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
func SetupStaticRoutes(engine *gin.Engine, assetsPath, imagePath string) {
	engine.Static("/assets", assetsPath)
	engine.Static("/journal-images", imagePath)
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
