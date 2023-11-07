package util

import (
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
)

/* Graceful Shutdown Handler to aid in Clean Exits.  */
type GracefullShutdown struct {
	quit chan os.Signal
}

/*
Wait waits for an interrupt signal to gracefully shutdown the server.
Application should wait on this function and start cleanup.

kill (no param) default send syscall.SIGTERM
kill -2 is syscall.SIGINT
kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
*/
func (self *GracefullShutdown) Wait() {
	signal.Notify(self.quit, syscall.SIGINT, syscall.SIGTERM)
	sigQuit := <-self.quit
	log.WithField("Signal", sigQuit).Info("Trying Graceful Shutting Signal")
}

/*
Initiates Graceful Shutdown from within
the Application, without any Signal.
*/
func (self *GracefullShutdown) Stop() {
	log.Info("GracefulShutdown Stop Received")

	go func() {
		self.quit <- syscall.SIGINT
	}()
}
