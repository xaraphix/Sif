package raft

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

var RaftEvents []RaftEvent

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

	RaftEvents = make([]RaftEvent, 0)
}

type RaftEvent struct {
	Timestamp time.Time
	Event     string
	Details   RaftEventDetails
}

type RaftEventDetails struct {
	Id            string
	CurrentRole   string
	CurrentTerm   int32
	CurrentLeader string
	VotedFor      string
	Peer          string
}

func LogEvent(event string, details RaftEventDetails) {
	RaftEvents = append(RaftEvents, RaftEvent{
		Timestamp: time.Now(),
		Event:     event,
		Details:   details,
	})
}
