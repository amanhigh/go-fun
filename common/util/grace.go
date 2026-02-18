package util

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
)

type Shutdown interface {
	Wait() context.Context
	Stop(c context.Context)
}

/* Graceful Shutdown Handler to aid in Clean Exits.  */
type GracefullShutdown struct {
	quit chan os.Signal
	ctx  context.Context
}

func NewGracefulShutdown() Shutdown {
	return &GracefullShutdown{make(chan os.Signal, 1), context.Background()}
}

/*
Wait waits for an interrupt signal to gracefully shutdown the server.
Application should wait on this function and start cleanup.

kill (no param) default send syscall.SIGTERM
kill -2 is syscall.SIGINT
kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
*/
func (gs *GracefullShutdown) Wait() (c context.Context) {
	signal.Notify(gs.quit, syscall.SIGINT, syscall.SIGTERM)
	sigQuit := <-gs.quit
	zerolog.Ctx(gs.ctx).Info().Any("Signal", sigQuit).Msg("Trying Graceful Shutting Signal")
	return gs.ctx
}

/*
Initiates Graceful Shutdown from within
the Application, without any Signal.
*/
func (gs *GracefullShutdown) Stop(c context.Context) {
	zerolog.Ctx(c).Info().Msg("GracefulShutdown Stop Received")
	gs.ctx = c

	go func() {
		gs.quit <- syscall.SIGINT
	}()
}
