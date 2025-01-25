package core

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/amanhigh/go-fun/common/tools"
	"github.com/amanhigh/go-fun/common/util"
	"github.com/bitfield/script"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	TICKER_LENGTH  = 15
	SIDE_MONITOR   = 1
	MAIL_WORKSPACE = "2"
	DATE_FORMAT    = "20060102__150405"
	LOGSEQ_CLASS   = "Logseq"
	TRADE_INFO     = `
Trends
HTF - Up
MTF - Up
TTF - Up

Plan: Longs @ TTF DZ

Obstacles:
-

Support:
-`
)

// FIXME: #C Extract AutoManager
func OpenTicker(ticker string) (err error) {
	// Focus on the window named "TradingView"
	log.Debug().Str("Ticker", ticker).Msg("OpenTicker")
	if err = tools.FocusWindow("TradingView"); err == nil {
		// Focus Input Box
		if err = tools.SendKey("-M Ctrl b -m Ctrl"); err == nil {
			// TASK: Copy Ticker once Clipboard Library is Fixed
			// Copy runs into doom loop with wl-paste Watch
			if err = tools.SendKey("-M Ctrl v -m Ctrl"); err == nil {
				time.Sleep(50 * time.Millisecond)
				// Bang ! to Open
				err = tools.SendInput("xox")
				// Return Focus Back
				if focusErr := tools.FocusLastWindow(); focusErr != nil {
					log.Error().Err(focusErr).Msg("Failed to return focus")
				}
			}
		}
	}
	return
}

func RecordTicker(ticker, path string) (err error) {
	if err = tools.FocusWindow("TradingView"); err == nil {
		log.Info().Str("Ticker", ticker).Msg("Recording Ticker")
		err = takeScreenshots(ticker, path)
		if err == nil && strings.Contains(ticker, ".set") {
			err = recordTradeInfo(ticker, path)
		}
		sendNotification(ticker)
	}
	return
}

func takeScreenshots(ticker, path string) (err error) {
	for i := 4; i > 0; i-- {
		if err = tools.SendKey("-k " + strconv.Itoa(i)); err == nil {
			name := fmt.Sprintf("%s__%s.png", ticker, time.Now().Format(DATE_FORMAT))
			log.Debug().Str("Ticker", ticker).Str("Name", name).Int("Count", i).Msg("Attempting Screenshot")
			time.Sleep(1 * time.Second)
			if err = tools.NamedScreenshot(path, name); err != nil {
				return
			}
		}
	}
	return
}

func recordTradeInfo(ticker, path string) (err error) {
	var tradeInfo string
	infoFile := fmt.Sprintf("%s/%s__%s.txt", path, ticker, time.Now().Format(DATE_FORMAT))
	if tradeInfo, err = tools.PromptText(TRADE_INFO); err == nil {
		if err = os.WriteFile(infoFile, []byte(tradeInfo), util.DEFAULT_PERM); err != nil {
			log.Error().Str("Ticker", ticker).Err(err).Msg("Failed to write trade info")
			return
		}

		// Record Check Screenshot
		checkFile := fmt.Sprintf("%s__%s.png", ticker, time.Now().Format(DATE_FORMAT))
		_ = tools.NamedRegionScreenshot(path, checkFile)
	} else {
		log.Error().Str("Ticker", ticker).Err(err).Msg("Read TradeInfo Failed")
	}
	return
}

func sendNotification(ticker string) {
	if err := tools.Notify(zerolog.InfoLevel, "Recorded", ticker); err != nil {
		log.Error().Err(err).Msg("Failed to send notification")
	}
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

func TryOpenTicker(ticker string) {
	window, err := tools.GetHyperWindow()
	if err == nil && window.Class == LOGSEQ_CLASS && window.Monitor == SIDE_MONITOR && window.Workspace.Name == MAIL_WORKSPACE {
		if openErr := OpenTicker(ticker); openErr != nil {
			log.Error().Err(err).Str("Ticker", ticker).Msg("Failed to open ticker")
		} else {
			log.Info().Str("Ticker", ticker).Msg("Opening Ticker")
		}
	} else {
		if err != nil {
			log.Error().Err(err).Msg("OpenTicker: GetHyperWindow Failed")
			return
		}
		log.Debug().Str("Ticker", ticker).Str("Class", window.Class).Int("Monitor", window.Monitor).Str("Workspace", window.Workspace.Name).Str("Window", window.Title).Msg("OpenTicker: Logseq Not Active")
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
