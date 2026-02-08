package util

import (
	"context"
	"errors"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
)

// WatermillController manages start/stop of a router backed by a gochannel.
type WatermillController interface {
	Start(ctx context.Context)
	Shutdown(ctx context.Context)
}

type watermillController struct {
	router  *message.Router
	channel *gochannel.GoChannel

	cancel context.CancelFunc
}

// NewWatermillController constructs a WatermillLifecycle implementation.
func NewWatermillController(router *message.Router, channel *gochannel.GoChannel) WatermillController {
	return &watermillController{
		router:  router,
		channel: channel,
	}
}

// Start launches the router loop.
func (wc *watermillController) Start(ctx context.Context) {
	if wc.router == nil {
		return
	}
	routerCtx, cancel := context.WithCancel(ctx)
	wc.cancel = cancel

	go func() {
		if err := wc.router.Run(routerCtx); err != nil && !errors.Is(err, context.Canceled) {
			// router logs its own errors
			return
		}
	}()
}

// Shutdown stops router and closes channel.
func (wc *watermillController) Shutdown(_ context.Context) {
	if wc.cancel != nil {
		wc.cancel()
	}
	if wc.router != nil {
		_ = wc.router.Close()
	}
	if wc.channel != nil {
		_ = wc.channel.Close()
	}
}
