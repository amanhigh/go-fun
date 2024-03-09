package core

import (
	"strconv"
	"time"

	"github.com/amanhigh/go-fun/common/tools"
	"github.com/amanhigh/go-fun/common/util"
	"github.com/bitfield/script"
	"github.com/fatih/color"
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

func MonitorSystem(runCmd string, wait, idle time.Duration) {
	go MonitorIdle(runCmd, idle, wait)
	MonitorInternetConnection(wait)
}

func MonitorInternetConnection(wait time.Duration) {
	util.ScheduleJob(wait, func(_ bool) {
		if tools.CheckInternetConnection() {
			color.Green("Internet UP: %v", time.Now())
		} else {
			color.Red("Internet Outage: %v", time.Now())
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
			color.Yellow("Idle Job Graceful Shutdown: %v", time.Now())
			return
		}

		if ok, err := tools.IsOSIdle(idle); err != nil {
			color.Red("Error Monitoring: %v", err)
			return
		} else if ok && cancel == nil {
			// Start Heavy Program when OS is Idle
			if cancel, err = tools.RunBackgroundProcess(runCmd); err != nil {
				color.Red("Start Program Failed: %v", err)
				return
			}
			color.Green("Heavy Program Started: %v", time.Now())
		} else if !ok && cancel != nil {
			// OS Not idle so Stop Program
			cancel()
			cancel = nil
			color.Red("Heavy Program Stopped: %v", time.Now())
		}
	})
}

func restartNetworkManager() {
	script.Exec("systemctl restart NetworkManager").Wait()
}
