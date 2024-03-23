package core

import (
	"fmt"
	"strconv"
	"time"

	"github.com/amanhigh/go-fun/common/tools"
	"github.com/amanhigh/go-fun/common/util"
	"github.com/bitfield/script"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	TICKER_LENGTH = 15
)

func OpenTicker(ticker string) (err error) {
	// Check if the length of the ticker is less than 15
	if len(ticker) < TICKER_LENGTH {
		// Focus on the window named "TradingView"
		log.Debug().Str("Ticker", ticker).Msg("OpenTicker")
		if err = tools.FocusWindow("brave-browser"); err == nil {
			// Focus Input Box
			if err = tools.SendKey("-M Ctrl b"); err == nil {
				// Copy Ticker
				if err = tools.ClipCopy(ticker); err == nil {
					// Paste the Ticker
					time.Sleep(3 * time.Second)
					if err = tools.SendKey("-M Ctrl v"); err == nil {
						time.Sleep(500 * time.Millisecond)
						// Code (xox) to Open
						err = tools.SendInput("xox")
					}
				}
			}
		}
	}
	return
}

func RecordTicker(ticker string) (err error) {
	// Bring Focus Back Lost due to Modal Box
	if err = tools.FocusWindow("TradingView"); err == nil {
		log.Info().Str("Ticker", ticker).Msg("Recording Ticker")
		// loop from max to 1
		for i := 4; i > 0; i-- {
			// emulate number key press
			if err = tools.SendKey("-k " + strconv.Itoa(i)); err == nil {
				// File Name POWERINDIA.mwd.trend.rejected.nca_20240321_193916.png
				name := fmt.Sprintf("%s__%s.png", ticker, time.Now().Format("20060102__150405"))
				log.Debug().Str("Ticker", ticker).Str("Name", name).Int("Count", i).Msg("Attempting Screenshot")

				// Wait
				time.Sleep(1 * time.Second)

				// Take Screenshot
				if err = tools.NamedScreenshot(name); err != nil {
					return
				}
			}
		}

		// send desktop notification
		tools.Notify(zerolog.InfoLevel, "Recorded", ticker)
	}

	return
}

func MonitorInternetConnection(wait time.Duration) {
	util.ScheduleJob(wait, func(_ bool) {
		if tools.CheckInternetConnection() {
			log.Info().Msg("Internet UP")
		} else {
			log.Warn().Msg("Internet DOWN")
			restartNetworkManager()
			//Extra Wait for Network Manager
			time.Sleep(5 * time.Second)
		}
	})
}

func MonitorIdle(runCmd string, wait, idle time.Duration) {
	var cancel util.CancelFunc

	// Start the monitoring
	util.ScheduleJob(wait, func(exit bool) {
		//Handle Graceful Shutdown
		if exit && cancel != nil {
			cancel()
			log.Info().Msg("Heavy Program Graceful Shutdown")
			return
		}

		if ok, err := tools.IsOSIdle(idle); err != nil {
			log.Error().Err(err).Msg("Idle Check failed")
			return
		} else if ok && cancel == nil {
			// Start Heavy Program when OS is Idle
			if cancel, err = tools.RunBackgroundProcess(runCmd); err != nil {
				log.Error().Str("Cmd", runCmd).Err(err).Msg("Start Program Failed")
				return
			}
			log.Info().Msg("Heavy Program Started")
		} else if !ok && cancel != nil {
			// OS Not idle so Stop Program
			cancel()
			cancel = nil
			log.Warn().Msg("Heavy Program Stopped")
		}
	})
}

func ProcessTicker(ticker string) {
	// BUG: Connect to Clipboard
	ok, err := tools.IsWindowFocused("trading-tome")
	if err != nil {
		log.Error().Err(err).Msg("OpenTicker: IsWindowFocused Failed")
		return
	}

	if ok {
		OpenTicker(ticker)
		log.Info().Str("Ticker", ticker).Msg("OpenTicker: Window Focused")
	} else {
		log.Debug().Str("Ticker", ticker).Msg("OpenTicker: Window Not Focused")
	}
}

func MonitorSubmap() {
	wait := time.Second
	util.ScheduleJob(wait, func(_ bool) {
		err := tools.ActivateSubmap("swiftkeys", "SwiftKeys")
		if err != nil {
			log.Error().Err(err).Msg("Activate Submap Failed")
		}
	})
}

func restartNetworkManager() {
	script.Exec("sudo systemctl restart NetworkManager").Wait()
}
