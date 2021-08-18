package raft

import (
	"os"

	"github.com/sirupsen/logrus"
)

func setupLogging() {
	logLevel := logrus.InfoLevel
	for _, a := range os.Environ() {
		if a == "LOG_LEVEL=debug" {
			//do more stuff here
			logLevel = logrus.DebugLevel
		}
	}
	// parse string, this is built-in feature of logrus
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logLevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors:   false,
		TimestampFormat: "2006-01-02 15:04:05.000",
		FullTimestamp:   true,
	})
}
