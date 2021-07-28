package raft

import (
	"math/rand"
	"time"
)

var LdrHrtbtMontr *LeaderHeartbeatMonitor = &LeaderHeartbeatMonitor{
	Monitor: &Monitor{Stopped: false,
		LastResetAt:     time.Time{},
		TimeoutDuration: time.Duration(rand.Intn(149)+150) * time.Millisecond,
	},
}

var ElctnMontr *ElectionMonitor = &ElectionMonitor{
	Monitor: &Monitor{
		Stopped:         true,
		LastResetAt:     time.Time{},
		TimeoutDuration: time.Duration(rand.Intn(149)+150) * time.Millisecond,
	},
}

type LeaderHeartbeatMonitor struct {
	*Monitor
}

type ElectionMonitor struct {
	*Monitor
}

func (l *LeaderHeartbeatMonitor) Start(rn *RaftNode) {

	l.LastResetAt = time.Now()
	go func(r *RaftNode) {
		for {
			if time.Since(l.LastResetAt) >= l.TimeoutDuration &&
				l.Stopped == false &&
				rn.ElectionInProgress == false {
				rn.ElectionManager.StartElection(rn)
			}

			if l.Stopped {
				break
			}

			l.Sleep()
		}
	}(rn)

}

func (e *ElectionMonitor) Start(rn *RaftNode) {

	e.LastResetAt = time.Now()

	go func(r *RaftNode) {
		for {
			if time.Since(e.LastResetAt) >= e.TimeoutDuration &&
				e.Stopped == false &&
				rn.ElectionInProgress == false {
				rn.ElectionManager.StartElection(rn)
			}

			if e.Stopped {
				break
			}

			e.Sleep()
		}
	}(rn)

}
