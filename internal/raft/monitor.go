package raft

import (
	"math/rand"
	"time"
)

var (
	electionMonitor        *ElectionMonitor
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

func NewElectionMonitor(forceNew bool) *ElectionMonitor {
	em := &ElectionMonitor{
		Monitor: &Monitor{
			Stopped:         true,
			LastResetAt:     time.Time{},
			TimeoutDuration: time.Duration(rand.Intn(149)+150) * time.Millisecond,
		},
	}
	if forceNew {
		return em
	} else {
		if electionMonitor == nil {
			electionMonitor = em
		}
		return electionMonitor
	}
}

type LeaderHeartbeatMonitor struct {
	*Monitor
}

type ElectionMonitor struct {
	*Monitor
}

func (l *LeaderHeartbeatMonitor) Start(rn *RaftNode) {

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
		}
	}(rn)

}

func (e *ElectionMonitor) Start(rn *RaftNode) {

	e.LastResetAt = time.Now()
	e.Stopped = false
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

//func (e *ElectionMonitor) Stop() {
//	e.Stopped = true
//}
//
//func (e *ElectionMonitor) GetLastResetAt() time.Time {
//	return e.LastResetAt
//}
//
//func (e *ElectionMonitor) Sleep() {
//	time.Sleep(e.TimeoutDuration)
//}
//
//func (e *ElectionMonitor) IsAutoStartOn() bool {
//	return e.AutoStart
//}
//
//func (l *LeaderHeartbeatMonitor) Stop() {
//	l.Stopped = true
//}
//
//func (l *LeaderHeartbeatMonitor) GetLastResetAt() time.Time {
//	return l.LastResetAt
//}
//
//func (l *LeaderHeartbeatMonitor) Sleep() {
//	time.Sleep(l.TimeoutDuration)
//}
//
//func (l *LeaderHeartbeatMonitor) IsAutoStartOn() bool {
//	return l.AutoStart
//}
