package manager

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/components/kohan/clients"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/rs/zerolog/log"
)

type TickerManager interface {
	DownloadTicker(ctx context.Context, ticker string) (err common.HttpError)
}

type TickerManagerImpl struct {
	client    clients.AlphaClient
	downloads string
}

func NewTickerManager(client clients.AlphaClient, downloads string) *TickerManagerImpl {
	return &TickerManagerImpl{
		client:    client,
		downloads: downloads,
	}
}

func (t *TickerManagerImpl) DownloadTicker(ctx context.Context, ticker string) (err common.HttpError) {
	filePath := filepath.Join(t.downloads, ticker+".json")

	// Check if file exists
	if info, statErr := os.Stat(filePath); statErr == nil {
		modTime := info.ModTime().Format("2006-01-02")
		log.Info().Str("Ticker", ticker).Str("ModTime", modTime).Msg("Ticker data already exists")
		return nil
	}

	// Create directory if it doesn't exist
	if err1 := os.MkdirAll(t.downloads, util.DIR_DEFAULT_PERM); err1 != nil {
		return common.NewServerError(err1)
	}

	// Fetch data using AlphaClient
	var data interface{}
	if data, err = t.client.FetchDailyPrices(ctx, ticker); err != nil {
		return err
	}

	// Save data to file
	if jsonData, err1 := json.Marshal(data); err1 == nil {
		if err1 = os.WriteFile(filePath, jsonData, util.DEFAULT_PERM); err1 != nil {
			return common.NewServerError(err1)
		}
		log.Info().Str("Ticker", ticker).Str("Path", filePath).Msg("Ticker data downloaded and saved")
	} else {
		return common.NewServerError(err1)
	}

	return nil
}
