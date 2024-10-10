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

func OpenTicker(ticker string) (err error) {
	// Focus on the window named "TradingView"
	log.Debug().Str("Ticker", ticker).Msg("OpenTicker")

	if err = tools.FocusWindow("TradingView"); err == nil {
		// Focus Input Box
		if err = tools.SendKey("-M Ctrl b -m Ctrl"); err == nil {
			// HACK: Copy Ticker once Clipboard Library is Fixed
			// Copy runs into doom loop with wl-paste Watch
			if err = tools.SendKey("-M Ctrl v -m Ctrl"); err == nil {
				time.Sleep(50 * time.Millisecond)
				// Bang ! to Open
				err = tools.SendInput("xox")
				// Return Focus Back
				tools.FocusLastWindow()
			}
		}
	}
	return
}

func RecordTicker(ticker, path string) (err error) {
	var tradeInfo string
	// Bring Focus Back Lost due to Modal Box
	if err = tools.FocusWindow("TradingView"); err == nil {
		log.Info().Str("Ticker", ticker).Msg("Recording Ticker")
		// loop from max to 1
		for i := 4; i > 0; i-- {
			// emulate number key press
			if err = tools.SendKey("-k " + strconv.Itoa(i)); err == nil {
				// File Name POWERINDIA.mwd.trend.rejected.nca_20240321_193916.png
				name := fmt.Sprintf("%s__%s.png", ticker, time.Now().Format(DATE_FORMAT))
				log.Debug().Str("Ticker", ticker).Str("Name", name).Int("Count", i).Msg("Attempting Screenshot")

				// Wait
				time.Sleep(1 * time.Second)

				// Take Screenshot
				if err = tools.NamedScreenshot(path, name); err != nil {
					return
				}
			}
		}

		// Read Trade Info for Set Trades
		if strings.Contains(ticker, ".set") {
			// Generate File name so it is ordered in Journal
			infoFile := fmt.Sprintf("%s/%s__%s.txt", path, ticker, time.Now().Format(DATE_FORMAT))
			if tradeInfo, err = tools.PromptText(TRADE_INFO); err == nil {
				os.WriteFile(infoFile, []byte(tradeInfo), util.DEFAULT_PERM)

				// Record Check Screenshot
				checkFile := fmt.Sprintf("%s__%s.png", ticker, time.Now().Format(DATE_FORMAT))
				err = tools.NamedRegionScreenshot(path, checkFile)

			} else {
				log.Error().Str("Ticker", ticker).Err(err).Msg("Read TradeInfo Failed")
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

func TryOpenTicker(ticker string) {
	window, err := tools.GetHyperWindow()
	if err == nil && window.Class == LOGSEQ_CLASS && window.Monitor == SIDE_MONITOR && window.Workspace.Name == MAIL_WORKSPACE {
		OpenTicker(ticker)
		log.Info().Str("Ticker", ticker).Msg("Opening Ticker")
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
