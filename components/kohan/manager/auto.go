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
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
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
}

type OSManagerImpl struct {
	screenshotPath string
}

func NewOSManager(screenshotPath string) *OSManagerImpl {
	return &OSManagerImpl{
		screenshotPath: screenshotPath,
	}
}

var _ OSManagerInterface = (*OSManagerImpl)(nil)

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
	if screenshotType == kohan.ScreenshotTypeRegion {
		screenshotErr = tools.NamedRegionScreenshot(dir, fileName)
	} else {
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
	if directoryType == kohan.ScreenshotDirectoryTypeDownload {
		return defaultDownloadsDir()
	}
	now := time.Now()
	return filepath.Join(a.screenshotPath, now.Format("2006"), now.Format("01"))
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
	infoFile := filepath.Join(path, fmt.Sprintf("%s__%s.txt", ticker, time.Now().Format(DATE_FORMAT)))
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
