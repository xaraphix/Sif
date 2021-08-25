package raft

import (
	"os"
	"time"

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
	AckLength     int32
	CommitLength  int32
	SentLength    int32
}

func (rn *RaftNode) LogEvent(event string, details RaftEventDetails) {
	rn.EventLog = append(rn.EventLog, RaftEvent{
		Timestamp: time.Now(),
		Event:     event,
		Details:   details,
	})
}
