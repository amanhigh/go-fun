package kohan

// ScreenshotType defines the type of screenshot to capture.
type ScreenshotType string

const (
	ScreenshotTypeFull   ScreenshotType = "FULL"
	ScreenshotTypeRegion ScreenshotType = "REGION"
)

// ScreenshotDirectoryType defines the base directory where the screenshot is stored.
type ScreenshotDirectoryType string

const (
	ScreenshotDirectoryTypeJournal  ScreenshotDirectoryType = "JOURNAL"
	ScreenshotDirectoryTypeDownload ScreenshotDirectoryType = "DOWNLOAD"
)

// ScreenshotRequest represents a request to capture and store a single screenshot.
type ScreenshotRequest struct {
	FileName      string                  `json:"file_name" binding:"required,max=50,image_file"`
	DirectoryType ScreenshotDirectoryType `json:"directory_type" binding:"required,oneof=JOURNAL DOWNLOAD"`
	Type          ScreenshotType          `json:"type" binding:"required,oneof=FULL REGION"`
	Window        string                  `json:"window" binding:"omitempty,max=30"`
	Notify        bool                    `json:"notify"`
}

// ScreenshotResponse represents the response after a successful screenshot capture.
type ScreenshotResponse struct {
	FileName string `json:"file_name"`
	FullPath string `json:"full_path"`
}
