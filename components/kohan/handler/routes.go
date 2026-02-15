package handler

import "github.com/gin-gonic/gin"

// SetupImageRoutes configures image-related routes for the given v1 router group
func SetupImageRoutes(v1 *gin.RouterGroup, imageHandler ImageHandler) {
	entries := v1.Group("/journal-entries")
	images := entries.Group("/:id/images")
	{
		images.POST("", imageHandler.HandleCreateImage)
		images.GET("", imageHandler.HandleListImages)
		images.DELETE("/:imageId", imageHandler.HandleDeleteImage)
	}
}

// SetupNoteRoutes configures note-related routes for the given v1 router group
func SetupNoteRoutes(v1 *gin.RouterGroup, noteHandler NoteHandler) {
	entries := v1.Group("/journal-entries")
	notes := entries.Group("/:id/notes")
	{
		notes.POST("", noteHandler.HandleCreateNote)
		notes.GET("", noteHandler.HandleListNotes)
		notes.DELETE("/:noteId", noteHandler.HandleDeleteNote)
	}
}

// SetupTagRoutes configures tag-related routes for the given v1 router group
func SetupTagRoutes(v1 *gin.RouterGroup, tagHandler TagHandler) {
	entries := v1.Group("/journal-entries")
	tags := entries.Group("/:id/tags")
	{
		tags.POST("", tagHandler.HandleCreateTag)
		tags.GET("", tagHandler.HandleListTags)
		tags.DELETE("/:tagId", tagHandler.HandleDeleteTag)
	}
}

// SetupJournalEntryRoutes configures basic journal entry routes
func SetupJournalEntryRoutes(v1 *gin.RouterGroup, journalHandler JournalHandler) {
	entries := v1.Group("/journal-entries")
	{
		entries.GET("", journalHandler.HandleListEntries)
		entries.GET("/:id", journalHandler.HandleGetEntry)
		entries.POST("", journalHandler.HandleCreateEntry)
	}
}
