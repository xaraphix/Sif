package raft

import (
	"math/rand"
	"time"

	"github.com/sirupsen/logrus"
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

func DestructLeaderHeartbeatMonitor(lhm *LeaderHeartbeatMonitor) {
	lhm = nil
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
				rn.VotedFor == 0 &&
				rn.ElectionInProgress == false &&
				rn.CurrentRole != LEADER {

				logrus.WithFields(logrus.Fields{
					"FollowerId": rn.Id,
					"MyTerm":     rn.CurrentTerm,
				}).Debug("Did not hear leader's heartbeat for quite some time")

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
