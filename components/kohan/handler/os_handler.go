package handler

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/amanhigh/go-fun/common/tools"
	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/components/kohan/manager"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/kohan"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// OSHandler provides HTTP handlers for OS-level system operations.
//
//go:generate mockery --name OSHandler
type OSHandler interface {
	HandleScreenshot(ctx *gin.Context)
	HandleReadClip(ctx *gin.Context)
	HandleRecordTicker(ctx *gin.Context)
	HandleSubmapControl(ctx *gin.Context)
}

type OSHandlerImpl struct {
	capturePath  string
	allowedPaths []string
	autoManager  manager.AutoManagerInterface
}

var _ OSHandler = (*OSHandlerImpl)(nil)

// NewOSHandler creates a new OSHandler.
func NewOSHandler(capturePath string, autoManager manager.AutoManagerInterface, screenshotDirs []string) *OSHandlerImpl {
	allowedPaths := append([]string{capturePath}, screenshotDirs...)
	return &OSHandlerImpl{
		capturePath:  capturePath,
		allowedPaths: allowedPaths,
		autoManager:  autoManager,
	}
}

// HandleScreenshot handles POST /v1/os/screenshot.
func (h *OSHandlerImpl) HandleScreenshot(ctx *gin.Context) {
	var req kohan.ScreenshotRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		httpErr := util.ProcessValidationError(err)
		ctx.JSON(httpErr.Code(), httpErr)
		return
	}

	cleanPath := filepath.Clean(req.SavePath)
	if strings.Contains(cleanPath, "..") {
		ctx.JSON(http.StatusBadRequest, common.NewFailEnvelope(map[string]string{"save_path": "path traversal not allowed"}))
		return
	}

	var resolvedDir string
	if filepath.IsAbs(cleanPath) {
		resolvedDir = cleanPath
	} else {
		resolvedDir = filepath.Join(h.capturePath, cleanPath)
	}

	if !h.isPathAllowed(resolvedDir) {
		ctx.JSON(http.StatusBadRequest, common.NewFailEnvelope(map[string]string{"save_path": "invalid directory"}))
		return
	}
	if _, err := os.Stat(resolvedDir); os.IsNotExist(err) {
		if mkdirErr := os.MkdirAll(resolvedDir, 0o755); mkdirErr != nil {
			ctx.JSON(http.StatusBadRequest, common.NewFailEnvelope(map[string]string{"save_path": "directory not writable"}))
			return
		}
	}

	fullPath := filepath.Join(resolvedDir, req.FileName)
	if err := h.autoManager.Screenshot(ctx, req.Type, req.Window, fullPath); err != nil {
		log.Error().Str("file_name", req.FileName).Str("save_path", req.SavePath).Err(err).Msg("Screenshot failed")
		ctx.JSON(http.StatusInternalServerError, common.NewErrorEnvelope("Unable to capture screenshot", 50001))
		return
	}

	ctx.JSON(http.StatusOK, common.NewEnvelope(kohan.ScreenshotResponse{
		FileName:     req.FileName,
		RelativePath: filepath.Join(req.SavePath, req.FileName),
		FullPath:     fullPath,
	}))
}

// isPathAllowed checks whether the resolved path falls within one of the allowed base directories.
func (h *OSHandlerImpl) isPathAllowed(resolvedPath string) bool {
	for _, allowed := range h.allowedPaths {
		if strings.HasPrefix(resolvedPath, allowed) {
			return true
		}
	}
	return false
}

// HandleReadClip handles GET /v1/clip/
func (h *OSHandlerImpl) HandleReadClip(ctx *gin.Context) {
	text, err := tools.ClipPaste()
	if err == nil {
		ctx.JSON(http.StatusOK, text)
	} else {
		ctx.JSON(http.StatusInternalServerError, err.Error())
	}
}

// HandleRecordTicker handles GET /v1/os/ticker/:ticker/record
func (h *OSHandlerImpl) HandleRecordTicker(ctx *gin.Context) {
	ticker := ctx.Param("ticker")
	if err := h.autoManager.RecordTicker(ctx, ticker, h.capturePath); err == nil {
		ctx.JSON(http.StatusOK, "Success")
	} else {
		log.Error().Str("Ticker", ticker).Err(err).Msg("Record Ticker Failed")
		ctx.JSON(http.StatusInternalServerError, err.Error())
	}
}

// HandleSubmapControl handles POST /v1/submap/:action
func (h *OSHandlerImpl) HandleSubmapControl(ctx *gin.Context) {
	action := ctx.Param("action")

	var request struct {
		Submap string `json:"submap"`
		Ticker string `json:"ticker,omitempty"`
	}

	if err := ctx.ShouldBindJSON(&request); err == nil {
		switch action {
		case "enable":
			err = tools.HyperDispatch("submap " + request.Submap)
		case "disable":
			err = tools.HyperDispatch("submap reset")
		default:
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid action. Use 'enable' or 'disable'"})
			return
		}

		if err == nil {
			log.Info().Str("Action", action).Str("Submap", request.Submap).Msg("Submap Control")
			ctx.JSON(http.StatusOK, gin.H{"status": "success", "action": action, "submap": request.Submap})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}
