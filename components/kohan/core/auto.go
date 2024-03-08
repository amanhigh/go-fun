package core

import "github.com/amanhigh/go-fun/common/tools"

const TICKER_LENGTH = 15

func OpenTicker(ticker string) (err error) {
	// Check if the length of the ticker is less than 15
	if len(ticker) < TICKER_LENGTH {
		// Focus on the window named "TradingView"
		if err = tools.FocusWindowByTitle("TradingView"); err == nil {
			// Focus Input Box
			if err = tools.SendKey("ctrl+asciitilde"); err == nil {
				// Paste the Ticker & Bang to Open
				err = tools.SendInput(ticker + "!")
			}
		}
	}
	return
}
