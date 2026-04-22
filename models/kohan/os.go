package kohan

// ScreenshotType defines the type of screenshot to capture.
type ScreenshotType string

const (
	ScreenshotTypeFull   ScreenshotType = "FULL"
	ScreenshotTypeRegion ScreenshotType = "REGION"
)

// ScreenshotRequest represents a request to capture and store a single screenshot.
type ScreenshotRequest struct {
	FileName string         `json:"file_name" binding:"required,max=50,image_file"`
	SavePath string         `json:"save_path" binding:"required,max=100,save_path"`
	Type     ScreenshotType `json:"type" binding:"required,oneof=FULL REGION"`
	Window   string         `json:"window" binding:"omitempty,max=30"`
	Notify   bool           `json:"notify"`
}

// ScreenshotResponse represents the response after a successful screenshot capture.
type ScreenshotResponse struct {
	FileName     string `json:"file_name"`
	RelativePath string `json:"relative_path"`
	FullPath     string `json:"full_path"`
}
