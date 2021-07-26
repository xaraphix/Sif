package raft

import (
	"math/rand"
	"sync"
	"time"
)

const (
	FOLLOWER  = "follower"
	CANDIDATE = "candidate"
)

var (
	once     sync.Once
	raftnode *RaftNode
)

type Node struct {
	Mu sync.Mutex

	Id            int32
	CurrentTerm   int32
	CurrentRole   string
	CurrentLeader int32
	VotedFor      int32
	VotesReceived []int32
	CommitIndex   int32
	CommitLength  int32
	AckedLength   map[int32]int32
	SentLength    map[int32]int32
	LastApplied   int32
	NextIndex     int32
	MatchIndex    int32
	PrevLogIndex  int32
	// Peers         []Peer
	// Log           []*pb.Log

}

type RaftNode struct {
	Node
	LeaderHeartbeatMonitor *LeaderHeartbeatMonitor
	ElectionMonitor        *ElectionMonitor
	Heart                  *Heart
}

type Heart struct {
	DurationBetweenBeats time.Duration
}

type LeaderHeartbeatMonitor struct {
	Monitor
}

type ElectionMonitor struct {
	Monitor
}

type Monitor struct {
	TimeoutDuration time.Duration
	LastResetAt     time.Time
	Stopped         bool
}

func (m *Monitor) StopMonitor() {
	m.Stopped = true
}

func (rn *RaftNode) StartLeaderHeartbeatMonitor() {

	rn.LeaderHeartbeatMonitor.Stopped = false
	go func(r *RaftNode) {
		for {
			if time.Since(r.LeaderHeartbeatMonitor.LastResetAt) >= r.LeaderHeartbeatMonitor.TimeoutDuration {
				rn.StartElection()
			}

			time.Sleep(r.LeaderHeartbeatMonitor.TimeoutDuration)

			if r.LeaderHeartbeatMonitor.Stopped {
				break
			}

		}
	}(rn)
}

func (rn *RaftNode) StartElectionMonitor() {

	rn.ElectionMonitor.Stopped = false
	go func(r *RaftNode) {
		for {
			if time.Since(r.ElectionMonitor.LastResetAt) >= r.ElectionMonitor.TimeoutDuration {
				rn.StartElection()
			}

			time.Sleep(r.ElectionMonitor.TimeoutDuration)

			if r.ElectionMonitor.Stopped {
				break
			}

		}
	}(rn)
}

func (rn *RaftNode) StartElection() {
	rn.Mu.Lock()
	rn.CurrentRole = CANDIDATE
	rn.VotedFor = raftnode.Id
	rn.CurrentTerm = raftnode.CurrentTerm + 1
	rn.Mu.Unlock()
}

func NewRaftNode(returnExistingIfPresent bool) *RaftNode {
	if returnExistingIfPresent {
		once.Do(func() {
			raftnode = &RaftNode{}
			initializeRaftNode(raftnode)
			raftnode.StartLeaderHeartbeatMonitor()
		})
		return raftnode
	} else {
		raftnode = &RaftNode{}
		initializeRaftNode(raftnode)
		raftnode.StartLeaderHeartbeatMonitor()
		return raftnode
	}
}

func initializeRaftNode(rn *RaftNode) {
	rn.CurrentRole = FOLLOWER
	rn.CurrentTerm = 0
	rn.LeaderHeartbeatMonitor = &LeaderHeartbeatMonitor{
		Monitor: Monitor{
			Stopped:         false,
			LastResetAt:     time.Time{},
			TimeoutDuration: time.Duration(rand.Intn(150)+150) * time.Millisecond,
		},
	}
	rn.ElectionMonitor = &ElectionMonitor{
		Monitor: Monitor{
			Stopped:         false,
			LastResetAt:     time.Time{},
			TimeoutDuration: time.Duration(rand.Intn(150)+150) * time.Millisecond,
		},
	}
}
