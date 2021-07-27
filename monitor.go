package raft

import "time"

type Monitor struct {
	TimeoutDuration time.Duration
	LastResetAt     time.Time
	Stopped         bool
}

type LeaderHeartbeatMonitor struct {
	Monitor
}

type ElectionMonitor struct {
	Monitor
}

func (m *Monitor) StopMonitor() {
	m.Stopped = true
}

func turnOnElectionMonitor(rn *RaftNode) {

	rn.ElectionMonitor.Stopped = false
	rn.ElectionMonitor.LastResetAt = time.Now()
	go func(r *RaftNode) {
		for {
			if time.Since(r.ElectionMonitor.LastResetAt) >= r.ElectionMonitor.TimeoutDuration &&
				rn.ElectionMonitor.Stopped == false {
				rn.StartElection()
			}

			if r.ElectionMonitor.Stopped {
				break
			}

			time.Sleep(r.ElectionMonitor.TimeoutDuration)
		}
	}(rn)

}

func turnOnLeaderHeartbeatMonitor(rn *RaftNode) {

	rn.LeaderHeartbeatMonitor.Stopped = false
	rn.LeaderHeartbeatMonitor.LastResetAt = time.Now()

	go func(r *RaftNode) {
		for {
			if time.Since(r.LeaderHeartbeatMonitor.LastResetAt) >= r.LeaderHeartbeatMonitor.TimeoutDuration &&
				rn.LeaderHeartbeatMonitor.Stopped == false &&
				rn.ElectionInProgress == false {
				rn.StartElection()
			}

			if r.LeaderHeartbeatMonitor.Stopped {
				break
			}

			time.Sleep(r.LeaderHeartbeatMonitor.TimeoutDuration)
		}
	}(rn)

}
