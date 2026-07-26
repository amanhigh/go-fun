package tools

import (
	"fmt"

	"golang.design/x/clipboard"
)

// ClipCopy writes text to the system clipboard.
// The clipboard library must have been initialized (clipboard.Init) before calling this function.
// Returns a non-nil error only when the underlying write channel is nil.
func ClipCopy(text string) error {
	// Write sends data to the clipboard immediately.
	// The returned channel signals when data is overwritten; we discard it.
	ch := clipboard.Write(clipboard.FmtText, []byte(text))
	if ch == nil {
		return fmt.Errorf("clipboard write failed: returned nil channel")
	}
	return nil
}

// ClipPaste reads text from the system clipboard.
// The clipboard library must have been initialized (clipboard.Init) before calling this function.
// Returns an error if no text data is present.
func ClipPaste() (string, error) {
	data := clipboard.Read(clipboard.FmtText)
	if data == nil {
		return "", fmt.Errorf("clipboard read failed: no data available")
	}
	return string(data), nil
}
