package raft

import (
	"math/rand"
	"time"
)

var (
	leaderHeartbeatMonitor *LeaderHeartbeatMonitor
)

func NewLeaderHeartbeatMonitor(forceNew bool) *LeaderHeartbeatMonitor {
	lhm := &LeaderHeartbeatMonitor{
		Monitor: &Monitor{
			Stopped:         false,
			LastResetAt:     time.Time{},
			TimeoutDuration: time.Duration(rand.Intn(149)+150) * time.Millisecond,
		},
	}
	if forceNew {
		return lhm
	} else {
		if leaderHeartbeatMonitor == nil {
			leaderHeartbeatMonitor = lhm
		}
		return leaderHeartbeatMonitor
	}
}

type LeaderHeartbeatMonitor struct {
	*Monitor
}

func (l *LeaderHeartbeatMonitor) Start(rn *RaftNode, electionUpdates chan ElectionUpdates) {

	l.LastResetAt = time.Now()
	l.Stopped = false
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

			if l.Stopped {
				break
			}
		}
	}(rn)

}
