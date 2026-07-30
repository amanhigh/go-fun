package tools

import (
	"fmt"

	"github.com/bitfield/script"
)

func SendKey(keys string) error {
	_, err := script.Exec(fmt.Sprintf("wtype %v", keys)).String()
	if err != nil {
		return fmt.Errorf("failed to send key: %w", err)
	}
	return nil
}
