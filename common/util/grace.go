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
func (self *GracefullShutdown) Wait() (c context.Context) {
	signal.Notify(self.quit, syscall.SIGINT, syscall.SIGTERM)
	sigQuit := <-self.quit
	zerolog.Ctx(self.ctx).Info().Any("Signal", sigQuit).Msg("Trying Graceful Shutting Signal")
	return self.ctx
}

/*
Initiates Graceful Shutdown from within
the Application, without any Signal.
*/
func (self *GracefullShutdown) Stop(c context.Context) {
	zerolog.Ctx(c).Info().Msg("GracefulShutdown Stop Received")
	self.ctx = c

	go func() {
		self.quit <- syscall.SIGINT
	}()
}
