package handler

import "github.com/gin-gonic/gin"

// SetupMonitorRoutes configures monitor-related routes
func SetupMonitorRoutes(monitor *gin.RouterGroup, monitorHandler MonitorHandler) {
	monitor.GET("/ticker/:ticker/record", monitorHandler.HandleRecordTicker)
	monitor.GET("/clip/", monitorHandler.HandleReadClip)
	monitor.POST("/submap/:action", monitorHandler.HandleSubmapControl)
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

// SetupJournalEntryRoutes configures basic journal entry routes
func SetupJournalEntryRoutes(journal *gin.RouterGroup, journalHandler JournalHandler) {
	{
		journal.GET("", journalHandler.HandleListEntries)
		journal.GET("/:id", journalHandler.HandleGetEntry)
		journal.POST("", journalHandler.HandleCreateEntry)
		journal.DELETE("/:id", journalHandler.HandleDeleteEntry)
	}
}
