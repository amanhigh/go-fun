package util

import (
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

type GracefullShutdown struct {
	quit chan os.Signal
}

func NewGracefulShutdown() *GracefullShutdown {
	return &GracefullShutdown{make(chan os.Signal, 1)}
}

func (self *GracefullShutdown) Wait() {
	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(self.quit, syscall.SIGINT, syscall.SIGTERM)
	sigQuit := <-self.quit
	log.WithField("Signal", sigQuit).Info("Trying Graceful Shutting Signal")

}

func (self *GracefullShutdown) Close() {
	go func() {
		self.quit <- syscall.SIGINT
		log.Info("Trying Graceful Shutting with close")
	}()
}
