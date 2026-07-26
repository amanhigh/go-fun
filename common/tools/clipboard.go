package tools

import (
	"fmt"

	"golang.design/x/clipboard"
)

// ClipCopy writes text to the system clipboard.
// Returns an error if initialization failed or the clipboard is unavailable.
func ClipCopy(text string) error {
	if err := clipboard.Init(); err != nil {
		return fmt.Errorf("clipboard unavailable: %w", err)
	}
	// Write sends data to the clipboard immediately.
	// The returned channel signals when data is overwritten; we discard it.
	ch := clipboard.Write(clipboard.FmtText, []byte(text))
	if ch == nil {
		return fmt.Errorf("clipboard write failed: returned nil channel")
	}
	return nil
}

// ClipPaste reads text from the system clipboard.
// Returns an error if initialization failed, the clipboard is unavailable,
// or no text data is present.
func ClipPaste() (string, error) {
	if err := clipboard.Init(); err != nil {
		return "", fmt.Errorf("clipboard unavailable: %w", err)
	}
	data := clipboard.Read(clipboard.FmtText)
	if data == nil {
		return "", fmt.Errorf("clipboard read failed: no data available")
	}
	return string(data), nil
}
