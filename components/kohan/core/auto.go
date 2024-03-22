package core

import (
	"context"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/amanhigh/go-fun/common/tools"
	"github.com/amanhigh/go-fun/common/util"
	"github.com/bitfield/script"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.design/x/clipboard"
)

const (
	TICKER_LENGTH  = 15
	TICKER_REGEX   = "^[A-Za-z0-9_]+(\\.[A-Za-z-]+){3,}$"
	SCREENSHOT_AGE = 30 * time.Minute
)

var (
	matcher = regexp.MustCompile(TICKER_REGEX)
)

func OpenTicker(ticker string) (err error) {
	// Check if the length of the ticker is less than 15
	if len(ticker) < TICKER_LENGTH {
		// Focus on the window named "TradingView"
		if err = tools.FocusWindow("TradingView"); err == nil {
			// Focus Input Box
			if err = tools.SendKey("ctrl+asciitilde"); err == nil {
				// Paste the Ticker
				_ = tools.SendKey("Shift+Insert")
				time.Sleep(50 * time.Millisecond)
				// Bang to Open
				err = tools.SendInput("!")
			}
		}
	}
	return
}

func RecordTicker(ticker string) (err error) {
	// Bring Focus Back Lost due to Modal Box
	if err = tools.FocusWindow("TradingView"); err == nil {
		// loop from max to 1
		for i := 4; i > 0; i-- {
			// emulate number key press
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
		tools.Notify(zerolog.InfoLevel, "SCREENSHOTTED....", ticker)
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

func MonitorClipboard(ctx context.Context, capturePath string) {
	log.Info().Msg("Monitoring Clipboard")
	ch := clipboard.Watch(ctx, clipboard.FmtText)
	for clipText := range ch {
		ticker := string(clipText)

		// Read OS Environment
		if matcher.MatchString(ticker) {
			log.Info().Str("Ticker", ticker).Msg("Recording Ticker")
			if err := RecordTicker(ticker); err == nil {
				LabelJournal(capturePath, ticker)
			} else {
				log.Error().Str("Ticker", ticker).Err(err).Msg("Open Ticker Failed")
			}
		} else {
			// BUG: Fix ActiveWindow and Add check for Label Journal as well.
			windowName, err := tools.GetActiveWindow()
			if err != nil {
				log.Error().Str("Ticker", ticker).Err(err).Msg("Active Window/Desktop Detect Failed")
				continue
			}
			subLogger := log.With().Str("Ticker", ticker).Str("Window", windowName).Logger()

			subLogger.Debug().Msg("Window Match")
			if strings.Contains(windowName, "trading-tome") {
				subLogger.Debug().Msg("Open Ticker")
				OpenTicker(ticker)
			} else {
				subLogger.Warn().Msg("No Ticker or Window Match")
			}
		}
	}
	log.Info().Msg("Clipboard Monitor Stopped")
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

func LabelJournal(path string, ticker string) {
	files, _ := script.FindFiles(path).Slice()

	for _, file := range files {
		log.Debug().Str("File", file).Msg("Checking File")

		if strings.Contains(file, "Screenshot") {
			// Check Age of Files
			info, _ := os.Stat(file)
			diff := time.Now().Sub(info.ModTime())
			log.Debug().Dur("Age", diff).Msg("File Age")

			// Age Within Threshold, Perform Rename
			if diff < SCREENSHOT_AGE*2 {
				newName := strings.ReplaceAll(file, "Screenshot", ticker)
				log.Info().Str("Old", file).Str("New", newName).Msg("Rename File")
				os.Rename(file, newName)
			}
		}
	}
}

func restartNetworkManager() {
	script.Exec("sudo systemctl restart NetworkManager").Wait()
}
