package core

import (
	"strconv"
	"time"

	"github.com/amanhigh/go-fun/common/tools"
)

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

func RecordTicker(ticker string) (err error) {
	// Bring Focus Back Lost due to Modal Box
	if err = tools.FocusWindowByTitle("TradingView"); err == nil {
		// loop from max to 1
		for i := 4; i > 0; i-- {
			// emulate number key press with xdotool
			if err = tools.SendKey(strconv.Itoa(i)); err == nil {
				// Wait
				time.Sleep(1 * time.Second)

				// Take Screenshot
				if err = tools.Screenshot(); err != nil {
					return
				}
			}
		}

		// send desktop notification
		tools.Notify("SCREENSHOTTED....", ticker)
	}

	return
}
