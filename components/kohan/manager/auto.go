package manager

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/amanhigh/go-fun/common/tools"
	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/kohan"
	"github.com/go-co-op/gocron/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	// networkRecoveryCooldown enforces a minimum gap between recovery attempts.
	// It is fixed policy and is not externally configurable.
	networkRecoveryCooldown = 5 * time.Minute

	DATE_FORMAT = "20060102__150405"
	TRADE_INFO  = `
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

type OSManagerInterface interface {
	Screenshot(ctx context.Context, directoryType kohan.ScreenshotDirectoryType, fileName string, screenshotType kohan.ScreenshotType, window string) (string, common.HttpError)
	RecordTicker(ctx context.Context, ticker string) common.HttpError
	MonitorInternetConnection(ctx context.Context)
}

type OSManagerImpl struct {
	wait           time.Duration
	screenshotPath string
	scheduler      gocron.Scheduler

	// consecutiveFailures counts back-to-back unreachable gateway checks.
	consecutiveFailures int
	// lastRecovery records the clock time of the most recent recovery attempt.
	lastRecovery time.Time
}

func NewOSManager(wait time.Duration, screenshotPath string, scheduler gocron.Scheduler) *OSManagerImpl {
	return &OSManagerImpl{
		wait:           wait,
		screenshotPath: screenshotPath,
		scheduler:      scheduler,
	}
}

var _ OSManagerInterface = (*OSManagerImpl)(nil)

// Copy existing implementations preserving comments but as methods
func (a *OSManagerImpl) Screenshot(_ context.Context, directoryType kohan.ScreenshotDirectoryType, fileName string, screenshotType kohan.ScreenshotType, window string) (string, common.HttpError) {
	if window != "" {
		if err := tools.FocusWindow(window); err != nil {
			return "", common.NewServerError(err)
		}
	}

	dir := a.resolveDir(directoryType)
	if err := os.MkdirAll(dir, util.DIR_DEFAULT_PERM); err != nil {
		return "", common.NewServerError(err)
	}
	fullPath := filepath.Join(dir, fileName)

	log.Info().Str("Dir", dir).Str("Name", fileName).Str("Type", string(screenshotType)).Msg("Capturing Screenshot")

	var screenshotErr error
	switch screenshotType {
	case kohan.ScreenshotTypeRegion:
		screenshotErr = tools.NamedRegionScreenshot(dir, fileName)
	case kohan.ScreenshotTypeFull:
		screenshotErr = tools.Screenshot(dir, fileName)
	default:
		screenshotErr = tools.Screenshot(dir, fileName)
	}
	return fullPath, a.mapScreenshotError(screenshotErr)
}

// mapScreenshotError converts screenshot tool errors into the appropriate HttpError.
// User aborts → 409 Conflict; genuine tool/framework failures → 500.
func (a *OSManagerImpl) mapScreenshotError(err error) common.HttpError {
	if err == nil {
		return nil
	}
	if errors.Is(err, tools.ErrScreenshotAborted) {
		return common.NewHttpError(err.Error(), http.StatusConflict)
	}
	return common.NewServerError(err)
}

func (a *OSManagerImpl) resolveDir(directoryType kohan.ScreenshotDirectoryType) string {
	switch directoryType {
	case kohan.ScreenshotDirectoryTypeDownload:
		return defaultDownloadsDir()
	case kohan.ScreenshotDirectoryTypeJournal:
		return filepath.Join(a.screenshotPath, time.Now().Format("2006"), time.Now().Format("01"))
	default:
		return filepath.Join(a.screenshotPath, time.Now().Format("2006"), time.Now().Format("01"))
	}
}

func defaultDownloadsDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(os.TempDir(), "Downloads")
	}
	return filepath.Join(homeDir, "Downloads")
}

func (a *OSManagerImpl) RecordTicker(_ context.Context, ticker string) common.HttpError {
	var err error
	if err = tools.FocusWindow("TradingView"); err == nil {
		log.Info().Str("Ticker", ticker).Msg("Recording Ticker")
		path := a.resolveDir(kohan.ScreenshotDirectoryTypeJournal)
		if mkErr := os.MkdirAll(path, util.DIR_DEFAULT_PERM); mkErr != nil {
			return common.NewServerError(mkErr)
		}
		err = a.takeScreenshots(ticker, path)
		if err == nil && strings.Contains(ticker, ".set") {
			err = a.recordTradeInfo(ticker, path)
		}
		a.sendNotification(ticker)
	}
	if err != nil {
		return common.NewServerError(err)
	}
	return nil
}

func (a *OSManagerImpl) takeScreenshots(ticker, path string) (err error) {
	for i := 4; i > 0; i-- {
		if err = tools.SendKey("-k " + strconv.Itoa(i)); err == nil {
			name := fmt.Sprintf("%s__%s.png", ticker, time.Now().Format(DATE_FORMAT))
			log.Debug().Str("Ticker", ticker).Str("Name", name).Int("Count", i).Msg("Attempting Screenshot")
			time.Sleep(1 * time.Second)
			if err = tools.Screenshot(path, name); err != nil {
				return
			}
		}
	}
	return
}

func (a *OSManagerImpl) recordTradeInfo(ticker, path string) (err error) {
	var tradeInfo string
	infoFile := fmt.Sprintf("%s/%s__%s.txt", path, ticker, time.Now().Format(DATE_FORMAT))
	if tradeInfo, err = tools.PromptText(TRADE_INFO); err == nil {
		if err = os.WriteFile(infoFile, []byte(tradeInfo), util.DEFAULT_PERM); err != nil {
			log.Error().Str("Ticker", ticker).Err(err).Msg("Failed to write trade info")
			return
		}

		// Record Check Screenshot
		checkFile := fmt.Sprintf("%s__%s.png", ticker, time.Now().Format(DATE_FORMAT))
		if checkErr := tools.NamedRegionScreenshot(path, checkFile); checkErr != nil {
			log.Warn().Str("Ticker", ticker).Err(checkErr).Msg("Checklist screenshot not saved (user may have aborted)")
		}
	} else {
		log.Error().Str("Ticker", ticker).Err(err).Msg("Read TradeInfo Failed")
	}
	return
}

func (a *OSManagerImpl) sendNotification(ticker string) {
	if err := tools.Notify(zerolog.InfoLevel, "Recorded", ticker); err != nil {
		log.Error().Err(err).Msg("Failed to send notification")
	}
}

func (a *OSManagerImpl) MonitorInternetConnection(ctx context.Context) {
	// NewJob cannot fail with static hardcoded inputs; Shutdown failing is benign
	_, _ = a.scheduler.NewJob(
		gocron.DurationJob(a.wait),
		gocron.NewTask(func() {
			a.monitorInternetConnection()
		}),
		gocron.WithSingletonMode(gocron.LimitModeReschedule),
	)

	a.scheduler.Start()
	<-ctx.Done()
	_ = a.scheduler.Shutdown()
}

// monitorInternetConnection probes the active connection once and delegates the
// recovery policy and actions to recoverConnection. A resolution failure is
// environmental (not a gateway fault), so it is logged and skipped without
// changing failure state.
func (a *OSManagerImpl) monitorInternetConnection() {
	conn, err := tools.ResolveWiFiConnection()
	if err != nil {
		// Resolution is environmental, not a gateway fault: warn and skip.
		log.Warn().Err(err).Msg("Failed to resolve active Wi-Fi connection; skipping recovery")
		return
	}

	reachable := tools.ConnectionGatewayReachable(conn)
	a.recoverConnection(conn, reachable, tools.ReconnectConnection, tools.RestartNetworkManager)
}

// recoverConnection applies the failure-count, cooldown, and reset policy and,
// when warranted, performs a targeted reconnect (falling back to a restart).
//
// reconnect and restart are injected callbacks so the recovery policy can be
// exercised without invoking nmcli/ping/sudo/systemctl. In production they are
// the real tools.ReconnectConnection and tools.RestartNetworkManager.
func (a *OSManagerImpl) recoverConnection(conn string, reachable bool, reconnect func(string) error, restart func() error) {
	if reachable {
		// Gateway healthy: reset prior failure state.
		a.consecutiveFailures = 0
		return
	}

	a.consecutiveFailures++
	// Policy: wait for two consecutive failures before acting.
	if a.consecutiveFailures < 2 {
		return
	}

	// Throttle recovery to one attempt per cooldown window.
	if !a.lastRecovery.IsZero() && time.Since(a.lastRecovery) < networkRecoveryCooldown {
		return
	}

	// Outside cooldown: record attempt time, then try targeted reconnect first.
	a.lastRecovery = time.Now()
	if err := reconnect(conn); err != nil {
		log.Error().Err(err).Str("Connection", conn).Msg("Targeted reconnect failed; falling back to NetworkManager restart")
		if rErr := restart(); rErr != nil {
			log.Error().Err(rErr).Msg("NetworkManager restart fallback failed")
		}
	}

	// Recovery attempted: start a fresh failure cycle.
	a.consecutiveFailures = 0
}
