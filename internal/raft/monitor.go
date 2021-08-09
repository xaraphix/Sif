package raft

import (
	"math/rand"
	"time"
)

func NewLeaderHeartbeatMonitor(forceNew bool) *LeaderHeartbeatMonitor {
	lhm := &LeaderHeartbeatMonitor{
		Monitor: &Monitor{
			Stopped:         false,
			LastResetAt:     time.Time{},
			TimeoutDuration: time.Duration(rand.Intn(149)+150) * time.Millisecond,
		},
	}
	return lhm
}

type LeaderHeartbeatMonitor struct {
	*Monitor
}

func (l *LeaderHeartbeatMonitor) Start(rn *RaftNode) {

	l.LastResetAt = time.Now()
	l.Stopped = false
	go func(r *RaftNode) {
		for {
			timeExceeded := time.Since(l.LastResetAt) >= l.TimeoutDuration

			if timeExceeded &&
				l.Stopped == false &&
				rn.ElectionInProgress == false &&
				rn.CurrentRole != LEADER {
				l.Stopped = true
				rn.ElectionMgr.StartElection(rn)
			}

			if l.Stopped {
				break
			}

			l.Sleep()

			if l.Stopped {
				break
			}
		}
	}(rn)

}

func (l *LeaderHeartbeatMonitor) Reset() {
	l.LastResetAt = time.Now()
}
