package main

import (
	"github.com/Sirupsen/logrus"
	"os"
)

func main() {
	logger := logrus.New()
	loggerStdOut := logrus.New()
	loggerStdOut.Out = os.Stdout
	if f, err := os.OpenFile("/tmp/testl", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666); err == nil {
		logger.Out = f
	}

	logger.Info("I am on File")
	loggerStdOut.Info("I am on StdOut")
}
