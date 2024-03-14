package core

import (
	"context"
	"errors"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/amanhigh/go-fun/common/tools"
	"github.com/amanhigh/go-fun/common/util"
	"github.com/bitfield/script"
	"github.com/fatih/color"
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
		if err = tools.FocusWindowByTitle("TradingView"); err == nil {
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

func MonitorClipboard(ctx context.Context, capturePath string) {
	color.Green("Monitoring Clipboard")
	ch := clipboard.Watch(ctx, clipboard.FmtText)
	for clipText := range ch {
		ticker := string(clipText)

		// Read OS Environment
		if matcher.MatchString(ticker) {
			color.Green("Recording Ticker: %s", ticker)
			if err := RecordTicker(ticker); err == nil {
				LabelJournal(capturePath, ticker)
			} else {
				color.Red("Open Ticker Failed: %v", err)
			}
		} else {
			// BUG: Fix ActiveWindow and Add check for Label Journal as well.
			windowName, err := tools.GetActiveWindow()
			desktop, err1 := tools.GetDesktop()
			err = errors.Join(err, err1)
			if err != nil {
				color.Red("Active Window/Desktop Detect Failed: %v , Ticker: %s", err, ticker)
				continue
			}

			color.Blue("Detected (W,D,T): %s || %s || %s", windowName, desktop, ticker)
			if strings.Contains(windowName, "trading-tome") {
				color.Green("Opening Ticker: %s", ticker)
				OpenTicker(ticker)
			} else {
				color.Yellow("No Ticker or Window Match: %s || %s", windowName, ticker)
			}
		}
	}
	color.Yellow("Stopping Clipboard Monitor")
}

func LabelJournal(path string, ticker string) {
	files, _ := script.FindFiles(path).Slice()

	for _, file := range files {
		color.Blue("Checking: %s", file)

		if strings.Contains(file, "Screenshot") {
			// Check Age of Files
			info, _ := os.Stat(file)
			diff := time.Now().Sub(info.ModTime())
			color.Blue("File Age: %s %v", file, diff)

			// Age Within Threshold, Perform Rename
			if diff < SCREENSHOT_AGE*2 {
				newName := strings.ReplaceAll(file, "Screenshot", ticker)
				color.Yellow("Renaming: %s", newName)
				os.Rename(file, newName)
			}
		}
	}
}

func restartNetworkManager() {
	script.Exec("sudo systemctl restart NetworkManager").Wait()
}
